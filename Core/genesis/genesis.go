package genesis

import (
	"fmt"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/block"
	"FichainCore/block_chain"
	"FichainCore/common"
	"FichainCore/common/hexutil"
	"FichainCore/common/math"
	"FichainCore/consensus/poa_consensus"
	"FichainCore/database"
	"FichainCore/params"
	"FichainCore/state"
)

type GenesisAuthority struct {
	Validators map[common.Address]*math.HexOrDecimal256 `json:"validators,omitempty"`
	Observers  map[common.Address][]common.Address      `json:"observers,omitempty"`
}

type GenesisFiatReserve struct {
	Balances map[common.Address]*math.HexOrDecimal256 `json:"balances"`
}

type GenesisAccount struct {
	Code       hexutil.Bytes               `json:"code,omitempty"`
	Storage    map[common.Hash]common.Hash `json:"storage,omitempty"`
	Balance    *math.HexOrDecimal256       `json:"balance"             gencodec:"required"`
	Nonce      math.HexOrDecimal64         `json:"nonce,omitempty"`
	PrivateKey hexutil.Bytes               `json:"secretKey,omitempty"`
}

type Genesis struct {
	Config     *params.ChainConfig                         `json:"config"`
	Nonce      math.HexOrDecimal64                         `json:"nonce"`
	Timestamp  math.HexOrDecimal64                         `json:"timestamp"`
	ExtraData  hexutil.Bytes                               `json:"extraData"`
	GasLimit   math.HexOrDecimal64                         `json:"gasLimit"   gencodec:"required"`
	Difficulty *math.HexOrDecimal256                       `json:"difficulty" gencodec:"required"`
	Mixhash    common.Hash                                 `json:"mixHash"`
	Proposer   common.Address                              `json:"proposer"`
	Alloc      map[common.UnprefixedAddress]GenesisAccount `json:"alloc"      gencodec:"required"`
	Number     math.HexOrDecimal64                         `json:"number"`
	GasUsed    math.HexOrDecimal64                         `json:"gasUsed"`
	ParentHash common.Hash                                 `json:"parentHash"`

	FiatReserve GenesisFiatReserve `json:"fiat_reserve"`
	Authority   GenesisAuthority   `json:"authority"`
}

func (g *Genesis) ToBlock(db database.Database) *block.Block {
	if db == nil {
		db, _ = database.NewMemDatabase()
	}
	statedb, _ := state.New(common.Hash{}, state.NewDatabase(db))
	for addr, account := range g.Alloc {
		address := common.Address(addr)
		statedb.AddBalance(address, account.Balance.ToBig())
		statedb.SetCode(address, account.Code)
		statedb.SetNonce(address, uint64(account.Nonce))
		for key, value := range account.Storage {
			statedb.SetState(address, key, value)
		}
	}
	root := statedb.IntermediateRoot(true)
	head := &block.BlockHeader{
		Height:     uint64(g.Number),
		Timestamp:  uint64(g.Timestamp),
		ParentHash: g.ParentHash,
		ExtraData:  g.ExtraData,
		GasUsed:    uint64(g.GasUsed),
		Prevrandao: common.BigToHash(g.Difficulty.ToBig()),
		StateRoot:  root,
		Proposer:   g.Proposer,
	}

	if g.Difficulty == nil {
		head.Prevrandao = common.BigToHash(params.GenesisDifficulty)
	}
	statedb.Commit(false)
	statedb.Database().TrieDB().Commit(root, true)

	return block.NewBlock(head, nil, nil, nil)
}

func (g *Genesis) MustCommit(db database.Database) *block.Block {
	block, err := g.Commit(db)
	if err != nil {
		panic(err)
	}
	return block
}

// Commit writes the block and state of a genesis specification to the database.
// The block is committed as the canonical head block.
func (g *Genesis) Commit(db database.Database) (*block.Block, error) {
	block := g.ToBlock(db)
	logger.DebugP("Commiting block", block.Header.Hash().String())
	if block.Header.Height != 0 {
		return nil, fmt.Errorf("can't commit genesis block with number > 0")
	}
	if err := block_chain.WriteTd(db, block.Hash(), block.Header.Height, g.Difficulty.ToBig()); err != nil {
		return nil, err
	}
	if err := block_chain.WriteBlock(db, block); err != nil {
		return nil, err
	}
	if err := block_chain.WriteBlockReceipts(db, block.Hash(), block.Header.Height, nil); err != nil {
		return nil, err
	}
	if err := block_chain.WriteCanonicalHash(db, block.Hash(), block.Header.Height); err != nil {
		return nil, err
	}
	if err := block_chain.WriteHeadBlockHash(db, block.Hash()); err != nil {
		return nil, err
	}
	if err := block_chain.WriteHeadHeaderHash(db, block.Hash()); err != nil {
		return nil, err
	}
	config := g.Config
	if config == nil {
		config = params.AllEthashProtocolChanges
	}
	return block, block_chain.WriteChainConfig(db, block.Hash(), config)
}

func (g *Genesis) CommitAuthorities(
	validatorDb database.Database,
	observerDb database.Database,
	fiatReserveDb database.Database,
) error {
	authority := poa_consensus.NewAuthority()
	fiatReserve := poa_consensus.NewFiatReserve()
	for i, v := range g.Authority.Validators {
		authority.AddValidator(i, v.ToBig())
	}
	for i, v := range g.Authority.Observers {
		authority.AddObserver(i, v)
	}
	for i, v := range g.FiatReserve.Balances {
		fiatReserve.Deposit(i, v.ToBig())
	}
	err := authority.CommitToStorage(validatorDb, observerDb)
	if err != nil {
		return err
	}

	return fiatReserve.CommitToStorage(fiatReserveDb)
}
