package handlers

import (
	"math/big"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/p2p"
	"FichainCore/p2p/message"
)

// simple handler for client

type ClientHandler struct{}

func (h *ClientHandler) Handlers() map[string]func(p2p.Peer, *message.Message) error {
	return map[string]func(p2p.Peer, *message.Message) error{
		message.MessageBalance: h.Balance,
	}
}

func (h *ClientHandler) Balance(peer p2p.Peer, msg *message.Message) error {
	logger.Info("Received balance: ",
		new(big.Int).SetBytes(msg.Payload.(*message.BytesMessage).Data).String(),
	)
	return nil
}
