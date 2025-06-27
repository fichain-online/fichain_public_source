package handlers

import (
	"FichainCore/p2p"
	"FichainCore/p2p/message"
)

type BlockHandler struct{}

func (h *BlockHandler) Handlers() map[string]func(p2p.Peer, *message.Message) error {
	return map[string]func(p2p.Peer, *message.Message) error{}
}
