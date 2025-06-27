package handlers

import (
	"encoding/binary"
	"math/big"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/call_data"
	"FichainCore/p2p"
	"FichainCore/p2p/message"
	"FichainCore/p2p/message_sender"
	"FichainCore/receipt"
	"FichainCore/transaction"
)

// simple handler for client

type ClientHandler struct {
	NonceChan        chan uint64
	ReceiptChan      chan *receipt.Receipt
	CallResponseChan chan *call_data.CallSmartContractResponse
	Sender           *message_sender.MessageSender
}

func NewClientHandler(
	Sender *message_sender.MessageSender,
) *ClientHandler {
	return &ClientHandler{
		NonceChan:        make(chan uint64),
		ReceiptChan:      make(chan *receipt.Receipt),
		CallResponseChan: make(chan *call_data.CallSmartContractResponse),
		Sender:           Sender,
	}
}

func (h *ClientHandler) Handlers() map[string]func(p2p.Peer, *message.Message) error {
	return map[string]func(p2p.Peer, *message.Message) error{
		message.MessageBalance:    h.Balance,
		message.MessageNonce:      h.Nonce,
		message.MessageTxMined:    h.TxMined,
		message.MessageReceipt:    h.TxReceipt,
		message.MessageCallResult: h.CallResult,
	}
}

func (h *ClientHandler) Balance(peer p2p.Peer, msg *message.Message) error {
	logger.Info("Received balance: ",
		new(big.Int).SetBytes(msg.Payload.(*message.BytesMessage).Data).String(),
	)
	return nil
}

func (h *ClientHandler) Nonce(peer p2p.Peer, msg *message.Message) error {
	num := binary.BigEndian.Uint64(msg.Payload.(*message.BytesMessage).Data)
	h.NonceChan <- num
	return nil
}

func (h *ClientHandler) CallResult(peer p2p.Peer, msg *message.Message) error {
	logger.Info("Received call result: ", msg.Payload.(*call_data.CallSmartContractResponse))
	h.CallResponseChan <- msg.Payload.(*call_data.CallSmartContractResponse)
	return nil
}

func (h *ClientHandler) TxMined(peer p2p.Peer, msg *message.Message) error {
	tx := msg.Payload.(*transaction.Transaction)
	// send get receipt
	err := h.Sender.SendMessageToPeer(
		peer,
		message.MessageGetReceipt,
		&message.BytesMessage{
			Data: tx.Hash().Bytes(),
		},
	)
	if err != nil {
		logger.Error("Error when send get receipt to node")
	}
	return nil
}

func (h *ClientHandler) TxReceipt(peer p2p.Peer, msg *message.Message) error {
	logger.Info("Received receipt: ", msg.Payload.(*receipt.Receipt))
	h.ReceiptChan <- msg.Payload.(*receipt.Receipt)
	return nil
}
