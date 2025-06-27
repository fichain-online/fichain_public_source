package cli

import (
	"bufio"
	"crypto/ecdsa"
	"encoding/hex"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strings"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/call_data"
	"FichainCore/cmd/client/handlers"
	"FichainCore/common"
	"FichainCore/common/hexutil"
	"FichainCore/crypto"
	"FichainCore/p2p"
	"FichainCore/p2p/message"
	"FichainCore/p2p/message_sender"
	"FichainCore/params"
	"FichainCore/signer"
	"FichainCore/transaction"
)

var (
	ErrorGetAccountStateTimedOut = errors.New("get account state timed out")
	ErrorInvalidAction           = errors.New("invalid action")
)

type Cli struct {
	// clientContext *client_context.ClientContext

	//
	stop     bool
	commands map[int]string
	reader   *bufio.Reader

	sender   *message_sender.MessageSender
	nodePeer p2p.Peer

	clientHandler *handlers.ClientHandler
	signer        *signer.Signer
	// transactionController client_types.TransactionController
	// accountStateChan      chan types.AccountState
	// defaultRelatedAddress map[common.Address][][]byte
}

func NewCli(
	sender *message_sender.MessageSender,
	nodePeer p2p.Peer,
	clientHandler *handlers.ClientHandler,
	signer *signer.Signer,
) *Cli {
	commands := map[int]string{
		0: "Exit",
		1: "Send transaction",
		2: "Call smart contract",
		3: "Create account",
		4: "Get balance",
		5: "Get stats",
		6: "Change log level",
	}
	return &Cli{
		sender:        sender,
		nodePeer:      nodePeer,
		clientHandler: clientHandler,
		signer:        signer,

		stop:     false,
		commands: commands,
	}
}

func (cli *Cli) Start() {
	cli.reader = bufio.NewReader(os.Stdin)
	for {
		if cli.stop {
			return
		}
		cli.PrintCommands()

		command := cli.ReadInput()
		switch command {
		case "0":
		// TODO
		case "1":
			err := cli.SendTransaction()
			if err != nil {
				logger.Warn("err", err)
			}
		case "2":
			cli.CallSmartContract()
		case "3":
			cli.CreateAccount()
		case "4":
			cli.PrintMessage("Enter address: ", "")
			cli.GetBalance(cli.ReadInputAddress())
			// TODO4
		case "5":
			cli.Subscribe()
		case "6":
			cli.PrintMessage("Enter address: ", "")
			// cli.StakeState(cli.ReadInputAddress())
		case "7":
		case "8":
			cli.GetStats()
		case "9":
			cli.ChangeLogLevel()
		}
	}
}

func (cli *Cli) Subscribe() {
	panic("TODO")
}

// TODE Cli stop
func (cli *Cli) Stop() {
}

func (cli *Cli) PrintCommands() {
	str := logger.Cyan + "======= Commands =======\n" + logger.Purple
	for i := 0; i < len(cli.commands); i++ {
		str += fmt.Sprintf("%v: %v\n", i, cli.commands[i])
	}
	str += logger.Reset
	fmt.Print(str)
}

func (cli *Cli) SendTransaction() error {
	cli.PrintMessage("Enter to address (leave empty for contract deployment):", "")
	toAddressStr := cli.ReadInput()

	var toAddress common.Address
	var data []byte
	var amount *big.Int
	isContractDeployment := (toAddressStr == "")

	if isContractDeployment {
		// This is a contract deployment
		toAddress = common.Address{} // Use the zero address for contract creation
		var err error
		data, err = cli.getDataForDeploySmartContract()
		if err != nil {
			return fmt.Errorf("could not get contract deployment data: %w", err)
		}

		cli.PrintMessage("Enter amount to send with deployment (default 0):", "")
		amountStr := cli.ReadInput()
		if amountStr == "" {
			amount = big.NewInt(0)
		} else {
			amount = new(big.Int)
			_, success := amount.SetString(amountStr, 10)
			if !success {
				logger.Warn("Invalid amount entered, using 0 as default.", "input", amountStr)
				amount = big.NewInt(0)
			}
		}
	} else {
		// This is a regular transaction
		toAddress = common.HexToAddress(toAddressStr)

		cli.PrintMessage("Enter to amount (default 10*10^18):", "")
		amount = cli.ReadBigInt()

		cli.PrintMessage(`Enter data (hex string, e.g., 0x...):`, "")
		dataStr := cli.ReadInput()
		data = common.FromHex(dataStr)
	}

	// get nonce
	fromAddress, err := cli.signer.WalletAddress()
	if err != nil {
		return err
	}
	cli.sender.SendMessageToPeer(
		cli.nodePeer,
		message.MessageGetNonce,
		&message.BytesMessage{
			Data: fromAddress.Bytes(),
		},
	)
	nonce := <-cli.clientHandler.NonceChan

	// Set appropriate gas limit and message
	var gas uint64
	var txMessage string
	if isContractDeployment {
		gas = uint64(5_000_000) // Higher default gas for deployment
		txMessage = "contract deployment"
	} else {
		gas = uint64(2_000_000) // Default gas for regular tx or contract call
		txMessage = "transfer ETH"
	}

	gasPrice := big.NewInt(20) // 20 Gwei

	// Create a new transaction
	tx := transaction.NewTransaction(toAddress, nonce, amount, data, gas, gasPrice, txMessage)

	hashSign, err := tx.HashSign(params.TempChainId)
	if err != nil {
		return fmt.Errorf("failed to create transaction hash for signing: %w", err)
	}
	s, err := cli.signer.SignHash(hashSign)
	logger.DebugP("signing hash", tx.Hash().String())
	if err != nil {
		return err
	}
	tx.SetSign(s)
	err = cli.sender.SendMessageToPeer(
		cli.nodePeer,
		message.MessageSendTransaction,
		tx,
	)
	if err != nil {
		return err
	}

	logger.Info("Transaction sent", "hash", tx.Hash().String())

	return nil
}

func (cli *Cli) CallSmartContract() {
	cli.PrintMessage("Enter to contract address: ", "")
	toAddress := cli.ReadInputAddress()
	cli.PrintMessage(`Enter data:`, "")
	dataStr := cli.ReadInput()
	data := common.FromHex(dataStr)
	callData := call_data.NewCallSmartContractData(toAddress, data)
	hash, err := callData.HashSign()
	if err != nil {
		logger.Error("error when create call data hash", hash)
		return
	}
	sign, err := cli.signer.SignHash(hash)
	callData.SetSign(sign)
	err = cli.sender.SendMessageToPeer(
		cli.nodePeer,
		message.MessageCallSmartContract,
		callData,
	)
	if err != nil {
		logger.Error("error when send message to peer", err)
		return
	}
	logger.Info("Call sent")
}

func (cli *Cli) CreateAccount() {
	privateKey, _ := crypto.GenerateKey()
	privateKeyHex := hexutil.Encode(crypto.FromECDSA(privateKey))[2:]
	publicKey := privateKey.Public()
	publicKeyECDSA, _ := publicKey.(*ecdsa.PublicKey)
	publicKeyHex := hexutil.Encode(crypto.FromECDSAPub(publicKeyECDSA))
	address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()
	logger.Info(
		fmt.Sprintf(
			"Private key: %v\nPublic key: %v\nAddress: %v\n",
			privateKeyHex,
			publicKeyHex,
			address,
		),
	)
}

func (cli *Cli) ReadInput() string {
	input, err := cli.reader.ReadString('\n')
	if err != nil {
		panic(err)
	}
	input = strings.Replace(input, "\n", "", -1)
	return input
}

func (cli *Cli) ReadInputAddress() common.Address {
	input := cli.ReadInput()
	address := common.HexToAddress(input)
	return address
}

func (cli *Cli) ReadBigInt() *big.Int {
	input := cli.ReadInput()
	if input == "" {
		input = "10000000000000000000"
	}
	bigInt := big.NewInt(0)
	bigInt.SetString(input, 10)
	return big.NewInt(0).SetBytes(bigInt.Bytes())
}

func (cli *Cli) PrintMessage(message string, color string) {
	if color == "" {
		color = logger.Purple
	}
	fmt.Printf(color+"%v\n"+logger.Reset, message)
}

func (cli *Cli) GetBalance(address common.Address) {
	cli.sender.SendMessageToPeer(
		cli.nodePeer,
		message.MessageGetBalance,
		&message.BytesMessage{
			Data: address.Bytes(),
		},
	)
}

func (cli *Cli) getDataForCallSmartContract() ([]byte, error) {
	return nil, nil
}

func (cli *Cli) ReadValidatorAddresses() []common.Address {
	cli.PrintMessage("Enter Validator Addresses: ", "")
	stringvalidatorAddresses := cli.ReadInput()
	if stringvalidatorAddresses == "" {
		return []common.Address{}
	}
	hexValidatorAddresses := strings.Split(stringvalidatorAddresses, ",")
	validatorAddresses := make([]common.Address, len(hexValidatorAddresses))
	for idx, hexAddress := range hexValidatorAddresses {
		address := common.HexToAddress(hexAddress)
		logger.Debug(address)
		validatorAddresses[idx] = address
	}
	return validatorAddresses
}

func (cli *Cli) GetStats() {
	// parentConn := cli.clientContext.ConnectionsManager.ParentConnection()
	// cli.clientContext.MessageSender.SendBytes(
	// 	parentConn,
	// 	command.GetStats,
	// 	[]byte{},
	// )
}

func (cli *Cli) ChangeLogLevel() {
	// parentConn := cli.clientContext.ConnectionsManager.ParentConnection()
	// str := p_common.Cyan + "======= Log level =======\n" + p_common.Purple
	// loglevel := map[int]string{
	// 	0: "DEBUGP",
	// 	1: "ERROR",
	// 	2: "WARN",
	// 	3: "INFO",
	// 	4: "DEBUG",
	// 	5: "TRACE",
	// }
	// for i := 0; i < len(loglevel); i++ {
	// 	str += fmt.Sprintf("%v: %v\n", i, loglevel[i])
	// }
	// fmt.Print(str)
	// level, err := strconv.Atoi(cli.ReadInput())
	// if err != nil {
	// 	logger.Error(err)
	// 	return
	// }
	// cli.clientContext.MessageSender.SendBytes(
	// 	parentConn,
	// 	command.ChangeLogLevel,
	// 	uint256.NewInt(
	// 		uint64(level)).Bytes(),
	// )
}

func (cli *Cli) getDataForDeploySmartContract() ([]byte, error) {
	cli.PrintMessage("Enter path to file containing contract creation bytecode (hex):", "")
	filePath := cli.ReadInput()
	if filePath == "" {
		return nil, errors.New("file path cannot be empty for contract deployment")
	}

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("file not found: %s", filePath)
	}

	fileContent, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("error reading file %s: %w", filePath, err)
	}

	// File content is expected to be a hex string, so trim whitespace and decode it.
	// hexString := strings.TrimSpace(string(fileContent))
	hexString := string(fileContent)
	data := common.FromHex(hexString)

	if len(data) == 0 && hexString != "" && hexString != "0x" {
		return nil, fmt.Errorf("invalid or empty hex data in file '%s'", filePath)
	}
	logger.Warn("Data", hex.EncodeToString(data))

	cli.PrintMessage(
		fmt.Sprintf("Read %d bytes of contract data from %s.", len(data), filePath),
		logger.Green,
	)

	return data, nil
}
