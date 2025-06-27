package client

import (
	"errors"
	"fmt"
	"math/big"
	"time"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/call_data"
	"FichainCore/client/handlers"
	"FichainCore/common"
	"FichainCore/config"
	"FichainCore/p2p"
	tcpClient "FichainCore/p2p/client"
	"FichainCore/p2p/message"
	"FichainCore/p2p/message_sender"
	"FichainCore/p2p/router"
	"FichainCore/params"
	"FichainCore/receipt"
	"FichainCore/signer"
	"FichainCore/transaction"
	"FichainCore/types"
)

/*
The client package simplifies sending blockchain transactions
with a single function call. It handles signing, encoding, and
broadcasting, making it easy to integrate transaction logic into
your app with minimal effort.
*/

type Client struct {
	ServerConnectionPeer p2p.Peer
	Sender               *message_sender.MessageSender
	Signer               *signer.Signer

	Handler *handlers.ClientHandler
}

func NewClient(
	coreConfig *config.Config,
) *Client {
	// init client and connect

	// Setup router with Ping handler
	router := router.NewRouter()
	sender := message_sender.NewMessageSender(nil)
	clientHandler := handlers.NewClientHandler(sender)

	for i, v := range clientHandler.Handlers() {
		router.RegisterHandler(i, v)
	}

	signerInstance := signer.NewSigner(
		types.PrivateKeyFromBytes(
			common.FromHex(coreConfig.PrivateKey),
		),
	)

	if config.GetConfig().BootAddress == "" {
		logger.Error("missing boot address")
		panic("err")
	}

	cl := tcpClient.NewTCPClient(
		30*time.Second,
		signerInstance,
	) // timeout 30s
	peer, err := cl.Dial(config.GetConfig().BootAddress)
	if err != nil {
		logger.Error(fmt.Sprintf(
			"[%s] Failed to connect to peer %s: %v",
			config.GetConfig().NodeID,
			config.GetConfig().BootAddress,
			err,
		))
		panic("err")
	}

	// Start reading loop
	go peer.ReadLoop(router)

	client := &Client{
		ServerConnectionPeer: peer,
		Sender:               sender,
		Signer:               signerInstance,
		Handler:              clientHandler,
	}
	return client
}

func (c *Client) SendTransaction(
	toAddress common.Address,
	amount *big.Int,
	data []byte,
	gas uint64,
	gasPrice *big.Int,
	txMessage string,
) (*receipt.Receipt, error) {
	// get nonce
	fromAddress, err := c.Signer.WalletAddress()
	if err != nil {
		return nil, err
	}
	c.Sender.SendMessageToPeer(
		c.ServerConnectionPeer,
		message.MessageGetNonce,
		&message.BytesMessage{
			Data: fromAddress.Bytes(),
		},
	)
	nonce := <-c.Handler.NonceChan

	// Create a new transaction
	tx := transaction.NewTransaction(toAddress, nonce, amount, data, gas, gasPrice, txMessage)

	hashSign, err := tx.HashSign(params.TempChainId)
	if err != nil {
		return nil, fmt.Errorf("failed to create transaction hash for signing: %w", err)
	}
	s, err := c.Signer.SignHash(hashSign)
	if err != nil {
		return nil, err
	}
	tx.SetSign(s)
	err = c.Sender.SendMessageToPeer(
		c.ServerConnectionPeer,
		message.MessageSendTransaction,
		tx,
	)
	if err != nil {
		return nil, err
	}

	logger.Info("Transaction sent", "hash", tx.Hash().String())
	select {
	case <-time.After(30 * time.Second):
		return nil, errors.New("wait for receipt timeout")
	case r := <-c.Handler.ReceiptChan:
		return r, nil
	}
}

func (c *Client) CallSmartContract(
	toAddress common.Address,
	data []byte,
) (*call_data.CallSmartContractResponse, error) {
	callData := call_data.NewCallSmartContractData(toAddress, data)
	hash, err := callData.HashSign()
	if err != nil {
		logger.Error("error when create call data hash", hash)
		return nil, err
	}
	sign, err := c.Signer.SignHash(hash)
	callData.SetSign(sign)
	err = c.Sender.SendMessageToPeer(
		c.ServerConnectionPeer,
		message.MessageCallSmartContract,
		callData,
	)
	if err != nil {
		logger.Error("error when send message to peer", err)
		return nil, err
	}
	logger.Info("Call sent")
	select {
	case <-time.After(30 * time.Second):
		return nil, errors.New("wait for receipt timeout")
	case r := <-c.Handler.CallResponseChan:
		return r, nil
	}
}
