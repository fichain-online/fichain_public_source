package handlers

import (
	"errors"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/block_chain"
	"FichainCore/common"
	"FichainCore/database"
	"FichainCore/p2p"
	"FichainCore/p2p/message"
	"FichainCore/p2p/message_sender"
	"FichainCore/params"
	"FichainCore/receipt"
	"FichainCore/state"
)

type ReceiptHandler struct {
	stateDB  *state.StateDB
	database database.Database

	sender *message_sender.MessageSender
	bc     *block_chain.BlockChain
}

func NewReceiptHandler(
	stateDB *state.StateDB,
	database database.Database,
	sender *message_sender.MessageSender,
	bc *block_chain.BlockChain,
) *ReceiptHandler {
	return &ReceiptHandler{
		stateDB:  stateDB,
		database: database,
		sender:   sender,
		bc:       bc,
	}
}

func (h *ReceiptHandler) Handlers() map[string]func(p2p.Peer, *message.Message) error {
	return map[string]func(p2p.Peer, *message.Message) error{
		message.MessageGetReceipt:  h.GetReceipt,
		message.MessageGetReceipts: h.GetReceipts,
	}
}

func (h *ReceiptHandler) GetReceipts(peer p2p.Peer, msg *message.Message) error {
	// todo chekc only explorer is allow to get all receipt
	blockHash := common.BytesToHash(msg.Payload.(*message.BytesMessage).Data)
	blockNumber := block_chain.GetBlockNumber(h.database, blockHash)
	blockReceipts := block_chain.GetBlockReceipts(h.database, blockHash, blockNumber)
	receipts := receipt.NewReceipts(blockReceipts)
	err := h.sender.SendMessageToPeer(
		peer,
		message.MessageReceipts,
		receipts,
	)
	return err
}

func (h *ReceiptHandler) GetReceipt(peer p2p.Peer, msg *message.Message) error {
	//
	txHash := common.BytesToHash(msg.Payload.(*message.BytesMessage).Data)
	// query block
	blockHash, blockNumber, txIndex := block_chain.GetTxLookupEntry(h.database, txHash)
	// query tx
	body := block_chain.GetBody(h.database, blockHash, blockNumber)
	if txIndex >= uint64(len(body.Transactions)) {
		logger.Error("invalid tx index in body")
		return errors.New("invalid tx index")
	}
	tx := body.Transactions[txIndex]
	// query receipt
	rcpts := block_chain.GetBlockReceipts(h.database, blockHash, blockNumber)
	if txIndex >= uint64(len(rcpts)) {
		logger.Error("invalid tx index in receipts")
		return errors.New("invalid tx index")
	}
	rcpt := rcpts[txIndex]
	from, _ := tx.From(params.TempChainId)
	rcpt.BlockHash = blockHash
	rcpt.BlockNumber = blockNumber
	rcpt.TxIndex = uint32(txIndex)
	rcpt.From = from
	rcpt.To = tx.To()
	rcpt.Amount = tx.Amount()
	err := h.sender.SendMessageToPeer(
		peer,
		message.MessageReceipt,
		rcpt,
	)
	return err
}
