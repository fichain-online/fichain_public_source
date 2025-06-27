package handlers

import (
	"fmt"
	"time"

	logger "github.com/hieuphanuit/golang-simple-logger"

	"FichainCore/block"
	"FichainCore/cmd/explorer/database"
	"FichainCore/cmd/explorer/models"
	"FichainCore/log"
	"FichainCore/p2p"
	"FichainCore/p2p/message"
	"FichainCore/p2p/message_sender"
	"FichainCore/params"
	"FichainCore/receipt"
	"FichainCore/transaction"
	"FichainCore/types"
)

type ChainEventHandler struct {
	Sender *message_sender.MessageSender
}

func NewChainEventHandler(
	Sender *message_sender.MessageSender,
) *ChainEventHandler {
	return &ChainEventHandler{
		Sender: Sender,
	}
}

func (h *ChainEventHandler) Handlers() map[string]func(p2p.Peer, *message.Message) error {
	return map[string]func(p2p.Peer, *message.Message) error{
		message.MessageChainEvent: h.ChainEvent,
		message.MessageReceipts:   h.Receipts,
	}
}

// blockToModel converts a block.Block to a models.BlockDB for database insertion.
func blockToModel(b *block.Block) *models.BlockDB {
	// Defensively handle nil inputs to prevent panics
	if b == nil || b.Header == nil {
		return nil
	}

	blockHash := b.Hash()

	// Convert the list of block transactions to database models
	txsDb := make([]*models.TransactionDB, len(b.Transactions))
	for i, tx := range b.Transactions {
		txsDb[i] = transactionToModel(tx, i, b)
	}

	bDb := &models.BlockDB{
		// Core Identifiers
		Hash:   blockHash.Bytes(),
		Height: b.Header.Height,

		// Header Fields - Direct Mappings
		ParentHash:       b.Header.ParentHash.Bytes(),
		StateRoot:        b.Header.StateRoot.Bytes(),
		TransactionsRoot: b.Header.TransactionsRoot.Bytes(),
		ReceiptRoot:      b.Header.ReceiptRoot.Bytes(),
		UncleHash:        b.Header.UncleHash.Bytes(),
		Proposer:         b.Header.Proposer.Bytes(),
		Prevrandao:       b.Header.Prevrandao.Bytes(),
		GasUsed:          b.Header.GasUsed,
		ExtraData:        b.Header.ExtraData,
		Signature:        b.Header.Signature,

		// Fields requiring conversion or calculation
		Timestamp: time.Unix(int64(b.Header.Timestamp), 0).
			UTC(),
		// Convert uint64 epoch to time.Time
		Bloom: b.Header.Bloom.Bytes(), // Use the 256-byte bloom slice

		// Denormalized count
		TransactionCount: len(b.Transactions),

		// Relationships
		Transactions: txsDb,
	}

	return bDb
}

func transactionToModel(
	tx *transaction.Transaction,
	idx int,
	bl *block.Block,
) *models.TransactionDB {
	// pbData := tx.Proto().(*pb.Transaction)
	from, _ := tx.From(params.TempChainId)
	to := tx.To()

	tDb := &models.TransactionDB{
		Hash:             tx.Hash().Bytes(),
		BlockHash:        bl.Hash().Bytes(),
		BlockHeight:      bl.Header.Height,
		TransactionIndex: uint32(idx),

		// Core Transaction Data
		FromAddress: from.Bytes(),
		ToAddress:   to.Bytes(),
		Nonce:       tx.Nonce(),
		Amount:      tx.Amount().Bytes(),
		GasLimit:    tx.Gas(), // Note: proto field is `Gas`, DB column is `gas_limit`
		GasPrice:    tx.GasPrice().Bytes(),
		Data:        tx.Data(),
		Message:     tx.Message(),

		// Signature components
		Signature: tx.Sign().Bytes(),
	}

	return tDb
}

func logToModel(log *log.Log, logIdx uint32, b *block.Block) *models.LogDB {
	lDb := &models.LogDB{
		// Composite Primary Key
		BlockHash: b.Hash().Bytes(),
		LogIndex:  logIdx,
		// Foreign key to transaction
		TransactionHash: b.Transactions[log.TxIndex].Hash().Bytes(),

		// Log Data
		EmitterAddress: log.Address.Bytes(),
		Data:           log.Data,
		Removed:        log.Removed,
	}

	// Safely populate topics by checking the slice length
	if len(log.Topics) > 0 {
		lDb.Topic0 = log.Topics[0].Bytes()
	}
	if len(log.Topics) > 1 {
		lDb.Topic1 = log.Topics[1].Bytes()
	}
	if len(log.Topics) > 2 {
		lDb.Topic2 = log.Topics[2].Bytes()
	}
	if len(log.Topics) > 3 {
		lDb.Topic3 = log.Topics[3].Bytes()
	}

	return lDb
}

func receiptToModel(receipt *receipt.Receipt) *models.ReceiptDB {
	rDb := &models.ReceiptDB{
		TransactionHash:   receipt.TxHash.Bytes(),
		Status:            receipt.Status,
		CumulativeGasUsed: receipt.CumulativeGasUsed,
		GasUsed:           receipt.GasUsed,
		ContractAddress:   receipt.ContractAddress.Bytes(),
		LogsBloom:         receipt.LogsBloom,
	}
	return rDb
}

func (h *ChainEventHandler) ChainEvent(peer p2p.Peer, msg *message.Message) error {
	// let save block, and transaction and logs
	chainEvent := msg.Payload.(*types.PbChainEventWrap)
	bl := &block.Block{}
	bl.FromProto(chainEvent.PbEvent.Block)
	bDb := blockToModel(bl)
	transactionsToSave := bDb.Transactions
	bDb.Transactions = nil
	// save to db
	rs := database.DB.Save(bDb)
	if rs.Error != nil {
		logger.Error("Error when save block to database", rs.Error)
		return rs.Error
	}
	logger.Info(fmt.Sprintf("Block number %v saved to db", bl.Header.Height))

	if len(transactionsToSave) > 0 {
		rs = database.DB.Save(transactionsToSave)
		if rs.Error != nil {
			logger.Error("Error when save block to database", rs.Error)
			return rs.Error
		}
		logger.Info(fmt.Sprintf("%v transactions saved to db", len(transactionsToSave)))
	}

	if len(chainEvent.PbEvent.GetLogs()) > 0 {
		logs := make([]*models.LogDB, len(chainEvent.PbEvent.GetLogs()))
		for i, v := range chainEvent.PbEvent.GetLogs() {
			l := &log.Log{}
			l.FromProto(v)
			logs[i] = logToModel(l, uint32(i), bl)
		}
		// save to db
		rs := database.DB.Save(logs)
		if rs.Error != nil {
			logger.Error("Error when save logs to database", rs.Error)
			return rs.Error
		}
		logger.Info(fmt.Sprintf("%v logs saved to db", len(logs)))
	}

	if len(bl.Transactions) > 0 {
		// send get receipts
		h.Sender.SendMessageToPeer(
			peer,
			message.MessageGetReceipts,
			&message.BytesMessage{
				Data: bl.Hash().Bytes(),
			},
		)
	}
	return nil
}

func (h *ChainEventHandler) Receipts(peer p2p.Peer, msg *message.Message) error {
	// let save block, and transaction and logs
	receipts := msg.Payload.(*receipt.Receipts)
	// get receipts
	receiptsDb := []*models.ReceiptDB{}
	for _, v := range receipts.Receipts() {
		receiptsDb = append(receiptsDb, receiptToModel(v))
	}
	if len(receiptsDb) > 0 {
		rs := database.DB.Save(receiptsDb)
		if rs.Error != nil {
			logger.Error("Error when save block to database", rs.Error)
			return rs.Error
		}
		logger.Info(fmt.Sprintf("%v receipt saved to db", len(receiptsDb)))
	}
	return nil
}
