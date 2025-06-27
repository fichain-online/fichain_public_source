package node

import (
	"math"
	"time"

	logger "github.com/hieuphanuit/golang-simple-logger"

	"FichainCore/block"
	"FichainCore/block_builder"
	"FichainCore/block_chain"
	"FichainCore/common"
	"FichainCore/config"
	"FichainCore/consensus/poa_consensus"
	"FichainCore/database"
	"FichainCore/evm"
	"FichainCore/handlers"
	"FichainCore/notifier"
	"FichainCore/p2p"
	"FichainCore/p2p/lookup_table"
	"FichainCore/p2p/message_sender"
	"FichainCore/p2p/router"
	"FichainCore/p2p/server"
	ws_server "FichainCore/p2p/ws/server"
	"FichainCore/params"
	"FichainCore/signer"
	"FichainCore/state"
	"FichainCore/transaction_pool"
	"FichainCore/transaction_validator"
	"FichainCore/types"
)

type Node struct {
	config *config.Config
	// boot node

	// networks
	server        *server.TCPServer
	wsServer      *ws_server.WebSocketServer
	lookupTable   *lookup_table.LookupTable
	messageSender *message_sender.MessageSender
	router        *router.Router

	// blockchain
	// ---- genesis
	// ---- storages
	// ---- statedb
	// ---- evm
	bc                   *block_chain.BlockChain
	authority            *poa_consensus.Authority
	transactionValidator *transaction_validator.TransactionValidator
	proposerSchedule     *poa_consensus.ProposerSchedule
	transactionPool      *transaction_pool.TransactionPool
	signer               *signer.Signer
	database             database.Database
	stateDB              *state.StateDB

	blockBuilder *block_builder.BlockBuilder

	clientNotifier   *notifier.ClientNotifier
	explorerNotifier *notifier.ExplorerNotifier

	// handler
	// ---- ping pong
	pingPongHandler    *handlers.PingPongHandler
	transactionHandler *handlers.TransactionHandler
	blockHandler       *handlers.BlockHandler
	stateHandler       *handlers.StateHandler
	receiptHandler     *handlers.ReceiptHandler

	// ---- consensus
}

func New() *Node {
	n := &Node{}
	n.initBlockchain()
	n.initNetwork()
	n.initHandlers()
	n.initOthers()
	return n
}

// inits
func (n *Node) initNetwork() {
	// Setup config
	cfg := p2p.DefaultConfig()
	cfg.ListenAddress = config.GetConfig().TCPServerAddress
	cfg.Debug = true
	//
	n.router = router.NewRouter()
	//
	n.lookupTable = lookup_table.NewLookupTable()
	//
	n.messageSender = message_sender.NewMessageSender(n.lookupTable)
	// Start server
	n.server = server.NewTCPServer(
		cfg.ListenAddress,
		n.router,
		n.lookupTable,
		n.signer,
	)
	n.wsServer = ws_server.NewWebSocketServer(
		// config.GetConfig().WsServerAddress,
		"/ws",
		n.router,
		n.lookupTable,
		n.signer,
	)
	logger.Info("Inited network")
}

func (n *Node) initHandlers() {
	n.pingPongHandler = handlers.NewPingPongHandler()
	n.transactionHandler = handlers.NewTransactionHandler(
		n.transactionValidator,
		n.proposerSchedule,
		n.transactionPool,
		// n, // todo
		n,
		n.bc,
		n.messageSender,
	)
	n.blockHandler = &handlers.BlockHandler{}
	n.stateHandler = handlers.NewStateHandler(
		n.stateDB,
		n.messageSender,
		n.bc,
	)
	n.receiptHandler = handlers.NewReceiptHandler(
		n.stateDB,
		n.database,
		n.messageSender,
		n.bc,
	)

	// register to router
	n.router.RegisterHanlders(n.pingPongHandler.Handlers())
	n.router.RegisterHanlders(n.transactionHandler.Handlers())
	n.router.RegisterHanlders(n.stateHandler.Handlers())
	n.router.RegisterHanlders(n.transactionHandler.Handlers())
	n.router.RegisterHanlders(n.receiptHandler.Handlers())

	logger.Info("Inited handlers")
}

func (n *Node) initBlockchain() {
	n.signer = signer.NewSigner(
		types.PrivateKeyFromBytes(
			common.FromHex(config.GetConfig().PrivateKey),
		),
	)

	// load storage, statedb
	bdb, err := database.NewBadgerDB(config.GetConfig().StatesDBPath)
	if err != nil {
		panic(err)
	}
	n.database = bdb
	headHeaderHash := block_chain.GetHeadHeaderHash(bdb)
	blockNumber := block_chain.GetBlockNumber(bdb, headHeaderHash)
	header := block_chain.GetHeader(bdb, headHeaderHash, blockNumber)
	logger.Info("Head block header", header)
	db := state.NewDatabase(bdb)
	n.stateDB, err = state.New(header.StateRoot, db)

	n.transactionValidator = &transaction_validator.TransactionValidator{
		EkycApiUrl:  config.GetConfig().EkycApiUrl,
		ChainConfig: params.TestChainConfig,
	}

	n.bc, err = block_chain.NewBlockChain(
		bdb,
		&block_chain.CacheConfig{
			Disabled: true, // TODO: check if this needed
		},
		params.TestChainConfig,        // TODO: check if this needed
		&poa_consensus.POAConsensus{}, // TODO: check if this needed
		evm.Config{},                  // TODO: check if this needed
	)
	if err != nil {
		logger.Error("error when init block chain ", err)
		panic(err.Error())
	}

	n.authority = poa_consensus.NewAuthority()
	validatorDB, err := database.NewBadgerDB(config.GetConfig().AuthorityValidatorDBPath)
	observerDB, err := database.NewBadgerDB(config.GetConfig().AuthorityObserverDBPath)
	n.authority.LoadFromDB(validatorDB, observerDB)

	scheduleBlockNum := uint64(params.TempEpochLength * math.Floor(
		float64(blockNumber)/float64(params.TempEpochLength),
	))
	logger.DebugP("[Node] scheduleBlockNum", scheduleBlockNum)
	salt := block_chain.GetCanonicalHash(bdb, scheduleBlockNum)

	n.proposerSchedule = poa_consensus.NewProposerSchedule(
		salt.Bytes(),
		scheduleBlockNum,
		n.authority.ListValidators(), // todo: need to pass block number to it so it con get correct list of validator at that time
	)

	n.transactionPool = transaction_pool.NewTransactionPool()

	address := n.Address()
	n.blockBuilder = block_builder.NewBlockBuilder(
		n.transactionPool,
		n.transactionValidator,
		&address,
		n.bc,
		n.stateDB,
	)

	logger.Info("Inited blockchain")
}

func (n *Node) initOthers() {
	n.clientNotifier = notifier.NewClientNotifier(
		n.messageSender,
	)
	n.clientNotifier.SubscribeChanEvent(n.bc)

	n.explorerNotifier = notifier.NewExplorerNotifier(
		n.messageSender,
		n.lookupTable,
	)
	n.explorerNotifier.SubscribeChanEvent(n.bc)
}

// run
func (n *Node) Start() {
	// listen tcp
	go func() {
		if err := n.server.Listen(); err != nil {
			logger.Error("[%s] Failed to start server: %v", config.GetConfig().NodeID, err)
			panic("Error")
		}
	}()
	// listen ws
	go func() {
		if err := n.wsServer.Listen(config.GetConfig().WsServerAddress); err != nil {
			logger.Error("[%s] Failed to start WS server: %v", config.GetConfig().NodeID, err)
			panic("Error")
		}
	}()

	// let test generate block
	go func() {
		for {
			bl, err := n.blockBuilder.GenerateBlock(
				n.bc.CurrentHeader(),
			)
			if err != nil {
				logger.Error("error when generate block", err)
				continue
			}
			// add block to chain
			_, err = n.bc.InsertChain([]*block.Block{bl})
			time.Sleep(1500 * time.Millisecond)
		}
	}()
	//
}

func (n *Node) Stop() {
	n.server.Close()
	n.wsServer.Close()
	// TODO: more clean up
}

// getters
func (n *Node) Address() common.Address {
	addr, err := n.signer.WalletAddress()
	if err != nil {
		panic(err)
	}
	return addr
}
