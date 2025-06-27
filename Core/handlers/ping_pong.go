package handlers

import (
	"time"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/config"
	"FichainCore/p2p"
	"FichainCore/p2p/message"
)

type PingPongHandler struct{}

func NewPingPongHandler() *PingPongHandler {
	return &PingPongHandler{}
}

func (h *PingPongHandler) Handlers() map[string]func(p2p.Peer, *message.Message) error {
	return map[string]func(p2p.Peer, *message.Message) error{
		message.MessagePing: h.Ping,
		message.MessagePong: h.Pong,
	}
}

func (h *PingPongHandler) Ping(peer p2p.Peer, msg *message.Message) error {
	logger.Info("Received ping from ", peer.ID())
	// Send pong response
	// pingMsg := &message.NewPong(*nodeID)
	pongMsg := &message.Pong{
		NodeID:    config.GetConfig().NodeID,
		Timestamp: time.Now().Unix(),
	}
	fmtMsg := &message.Message{
		Header: &message.Header{
			Version:     1,
			SenderID:    config.GetConfig().NodeID,
			MessageType: "pong",
			Timestamp:   time.Now().Unix(),
			Signature:   []byte{},
		},
		Payload: pongMsg,
	}

	err := peer.Send(fmtMsg)
	if err != nil {
		return err
	} else {
		logger.Info("Pong message sent to ", peer.ID())
	}
	return nil
}

func (h *PingPongHandler) Pong(peer p2p.Peer, msg *message.Message) error {
	logger.Info("Received pong from ", peer.ID())
	return nil
}
