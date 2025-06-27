package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"math/big"
	"os"
	"os/signal"
	"syscall"
	"time"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/common"
	"FichainCore/config"
	"FichainCore/handlers"
	"FichainCore/p2p"
	"FichainCore/p2p/client"
	"FichainCore/p2p/message"
	"FichainCore/p2p/message_sender"
	"FichainCore/p2p/router"
	"FichainCore/p2p/server"
	"FichainCore/signer"
	"FichainCore/transaction"
	"FichainCore/types"
)

func main() {
	// Parse command-line flags
	configFile := flag.String("config", "", "Config file path")
	flag.Parse()

	if *configFile == "" {
		log.Fatal("You must provide a config file --config")
	}
	config.InitConfig(*configFile)

	// Setup config
	cfg := p2p.DefaultConfig()
	cfg.ListenAddress = config.GetConfig().TCPServerAddress
	cfg.Debug = true

	// Setup router with Ping handler
	router := router.NewRouter()
	pingPongHandler := handlers.NewPingPongHandler()
	clientHandler := &handlers.ClientHandler{}
	for i, v := range pingPongHandler.Handlers() {
		router.RegisterHandler(i, v)
	}
	for i, v := range clientHandler.Handlers() {
		router.RegisterHandler(i, v)
	}

	signerInstance := signer.NewSigner(
		types.PrivateKeyFromBytes(
			common.FromHex(config.GetConfig().PrivateKey),
		),
	)

	// Start server
	server := server.NewTCPServer(
		cfg.ListenAddress,
		router,
		signerInstance,
	)
	// client dont have to listen
	// go func() {
	// 	if err := server.Listen(); err != nil {
	// 		log.Fatalf("[%s] Failed to start server: %v", config.GetConfig().NodeID, err)
	// 	}
	// }()

	// Optionally dial a peer
	if config.GetConfig().BootAddress == "" {
		logger.Error("missing boot address")
		panic("err")
	}

	client := client.NewTCPClient(
		30*time.Second,
		signerInstance,
	) // timeout 30s
	peer, err := client.Dial(config.GetConfig().BootAddress)
	if err != nil {
		logger.Error(
			"[%s] Failed to connect to peer %s: %v",
			config.GetConfig().NodeID,
			config.GetConfig().BootAddress,
			err,
		)
		panic("err")
	}

	log.Printf("[%s] Connected to peer: %s", config.GetConfig().NodeID, peer.ID())
	// Register the peer to the server so it can receive messages
	server.RegisterPeer(peer)
	// Send ping after connecting

	pingMsg := &message.Ping{
		NodeID:    config.GetConfig().NodeID,
		Timestamp: time.Now().Unix(),
	}
	fmtMsg := &message.Message{
		Header: &message.Header{
			Version:     1,
			SenderID:    config.GetConfig().NodeID,
			MessageType: "ping",
			Timestamp:   time.Now().Unix(),
			Signature:   []byte{},
		},
		Payload: pingMsg,
	}

	err = peer.Send(fmtMsg)
	if err != nil {
		slog.Info(fmt.Sprintf("Error when send message %v", err))
	}

	// create sender for ez send message
	sender := message_sender.NewMessageSender(
		nil,
	)

	// Test send get balance
	err = sender.SendMessageToPeer(peer, message.MessageGetBalance, &message.BytesMessage{
		Data: common.FromHex("0000000000000000000000000000000000000001"),
	})
	if err != nil {
		logger.Error("Error when send get balance", err)
	}
	logger.Info("Sent get balance")

	// Test send get transaction
	toAddress := common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678")
	nonce := uint64(1)
	amount := big.NewInt(1000000000000000000) // 1 ETH (assuming 18 decimals)
	data := []byte("example data")
	gas := uint64(21000)
	gasPrice := big.NewInt(20000000000) // 20 Gwei
	txMessage := "transfer ETH"

	// Create a new transaction
	tx := transaction.NewTransaction(toAddress, nonce, amount, data, gas, gasPrice, txMessage)
	err = sender.SendMessageToPeer(peer, message.MessageSendTransaction, tx)
	if err != nil {
		logger.Error("Error when send get balance", err)
	}
	logger.Info("Sent transaction")

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Printf("[%s] Shutting down...", config.GetConfig().NodeID)
	server.Close()
}
