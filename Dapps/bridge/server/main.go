package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	"FichainCore/client"
	"FichainCore/common"
	core_config "FichainCore/config"

	logger "github.com/HendrickPhan/golang-simple-logger"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"

	"FichainBridge/config"
	"FichainBridge/controllers"
	"FichainBridge/database"
	"FichainBridge/minter"
	"FichainBridge/routes"
	"FichainBridge/scanner"
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

	if config.GetConfig().Core.BootAddress == "" {
		logger.Error("missing boot address")
		panic("err")
	}

	// create minter
	txClient := client.NewClient(&config.GetConfig().Core)
	fichainTokenMap := map[string]common.Address{}
	for i, v := range config.GetConfig().FichainTokenMap {
		fichainTokenMap[i] = common.HexToAddress(v)
	}
	minter := minter.NewMinter(
		database.DB,
		txClient,
		fichainTokenMap,
	)

	// API
	scanner, err := scanner.NewBlockScannerService(
		database.DB,
		config.GetConfig().TokenMap,
		config.GetConfig().NetworkConnectionString,
		minter,
		"last_block.json",
	)
	if err != nil {
		panic("failed to create scanner service: " + err.Error())
	}

	go scanner.SubscribeAndProcess()
	// Initialize controllers
	depositWalletController := controllers.NewDepositWalletController(database.DB, scanner)
	depositLogController := controllers.NewDepositLogController(database.DB)
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
	routes.SetupDepositWalletRoute(apiRouter, depositWalletController)
	routes.SetupDepositLogRoute(apiRouter, depositLogController)
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
