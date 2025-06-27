package main

import (
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"FichainCore/common"
	"FichainCore/config"
	"FichainCore/handlers"
	"FichainCore/p2p"
	"FichainCore/p2p/client"
	"FichainCore/p2p/message"
	"FichainCore/p2p/router"
	"FichainCore/p2p/server"
	"FichainCore/signer"
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
	for i, v := range pingPongHandler.Handlers() {
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
	go func() {
		if err := server.Listen(); err != nil {
			log.Fatalf("[%s] Failed to start server: %v", config.GetConfig().NodeID, err)
		}
	}()

	// Optionally dial a peer
	if config.GetConfig().BootAddress != "" {
		client := client.NewTCPClient(
			30*time.Second,
			signerInstance,
		) // timeout 30s
		peer, err := client.Dial(config.GetConfig().BootAddress)
		if err != nil {
			log.Printf(
				"[%s] Failed to connect to peer %s: %v",
				config.GetConfig().NodeID,
				config.GetConfig().BootAddress,
				err,
			)
		} else {
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

			// Wait for peer to disconnect
			<-peer.Done()
		}
	}

	// Handle graceful shutdown
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop
	log.Printf("[%s] Shutting down...", config.GetConfig().NodeID)
	server.Close()
}
