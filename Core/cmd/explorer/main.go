package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	logger "github.com/HendrickPhan/golang-simple-logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"FichainCore/cmd/explorer/config"
	"FichainCore/cmd/explorer/controllers"
	"FichainCore/cmd/explorer/database"
	explorer_handlers "FichainCore/cmd/explorer/handlers"
	"FichainCore/cmd/explorer/routes"
	"FichainCore/common"
	core_config "FichainCore/config"
	"FichainCore/handlers"
	"FichainCore/p2p/client"
	"FichainCore/p2p/message_sender"
	"FichainCore/p2p/router"
	"FichainCore/signer"
	"FichainCore/types"
)

func main() {
	// set logger config
	loggerConfig := &logger.LoggerConfig{
		Flag:    logger.FLAG_DEBUG,
		Outputs: []*os.File{os.Stdout},
	}
	logger.SetConfig(loggerConfig)

	// Parse command-line flags
	configFile := flag.String("config", "", "Config file path")
	flag.Parse()

	if *configFile == "" {
		log.Fatal("You must provide a config file --config")
	}
	config.InitConfig(*configFile)
	core_config.SetConfig(&config.GetConfig().Core)

	database.Init(config.GetConfig().Database)

	// Setup router with Ping handler
	router := router.NewRouter()
	pingPongHandler := handlers.NewPingPongHandler()
	sender := message_sender.NewMessageSender(nil)
	chainEventHandler := explorer_handlers.NewChainEventHandler(
		sender,
	)
	for i, v := range pingPongHandler.Handlers() {
		router.RegisterHandler(i, v)
	}
	for i, v := range chainEventHandler.Handlers() {
		router.RegisterHandler(i, v)
	}

	signerInstance := signer.NewSigner(
		types.PrivateKeyFromBytes(
			common.FromHex(config.GetConfig().Core.PrivateKey),
		),
	)

	if config.GetConfig().Core.BootAddress == "" {
		logger.Error("missing boot address")
		panic("err")
	}

	cl := client.NewTCPClient(
		30*time.Second,
		signerInstance,
	) // timeout 30s
	peer, err := cl.Dial(config.GetConfig().Core.BootAddress)
	if err != nil {
		logger.Error(fmt.Sprintf(
			"[%s] Failed to connect to peer %s: %v",
			config.GetConfig().Core.NodeID,
			config.GetConfig().Core.BootAddress,
			err,
		))
		panic("err")
	}

	// Start reading loop
	go peer.ReadLoop(router)

	// API

	// Initialize controllers
	transactionController := controllers.NewTransactionController(database.DB)
	// Initialize the Gin router
	r := gin.Default()
	// Initialize cors config
	corsConfig := cors.DefaultConfig()
	corsConfig.AllowOrigins = []string{"*"}
	corsConfig.AllowHeaders = []string{"*"}
	corsConfig.AllowCredentials = true

	r.Use(cors.New(corsConfig))
	// Setup the user API routes
	apiRouter := r.Group("/api")
	// routes.i(apiRouter, hands) // Initialize user routes
	routes.SetupTransactionRoute(apiRouter, transactionController)
	// Run the server
	go func() {
		err = r.Run(config.GetConfig().APIAddress)
		if err != nil {
			log.Fatal("Failed to start the server")
		}
	}()
	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Printf("[%s] Shutting down...", config.GetConfig().Core.NodeID)
	// n.Stop()
}
