package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/config"
	"FichainCore/node"
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
	n := node.New()
	n.Start()

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Printf("[%s] Shutting down...", config.GetConfig().NodeID)
	n.Stop()
}
