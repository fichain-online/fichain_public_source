package main

import (
	"flag"
	"fmt"
	"time"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/cmd/client/cli"
	client_handlers "FichainCore/cmd/client/handlers"
	"FichainCore/common"
	"FichainCore/config"
	"FichainCore/handlers"
	"FichainCore/p2p/client"
	"FichainCore/p2p/message_sender"
	"FichainCore/p2p/router"
	"FichainCore/signer"
	"FichainCore/types"
)

func main() {
	configFile := flag.String("config", "", "Config file path")
	flag.Parse()

	if *configFile == "" {
		logger.Error("You must provide a config file --config")
		panic("")
	}
	config.InitConfig(*configFile)

	// Setup router with Ping handler
	router := router.NewRouter()
	pingPongHandler := handlers.NewPingPongHandler()
	sender := message_sender.NewMessageSender(nil)
	clientHandler := client_handlers.NewClientHandler(sender)

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

	if config.GetConfig().BootAddress == "" {
		logger.Error("missing boot address")
		panic("err")
	}

	cl := client.NewTCPClient(
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
	address, err := signerInstance.WalletAddress()
	if err != nil {
		panic(err)
	}

	logger.Info("Running client with address", address.String())
	c := cli.NewCli(
		sender,
		peer,
		clientHandler,
		signerInstance,
	)
	c.Start()
}
