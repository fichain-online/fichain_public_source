package handlers

import (
	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/consensus/poa_consensus"
	"FichainCore/p2p"
	"FichainCore/p2p/message"
	"FichainCore/p2p/message_sender"
	"FichainCore/transaction"
	"FichainCore/transaction_pool"
	"FichainCore/transaction_validator"
)

type TransactionHandler struct {
	transactionValidator *transaction_validator.TransactionValidator
	proposerSchedule     *poa_consensus.ProposerSchedule
	txPool               *transaction_pool.TransactionPool
	node                 Node
	bc                   Blockchain
	sender               *message_sender.MessageSender
}

func NewTransactionHandler(
	validator *transaction_validator.TransactionValidator,
	schedule *poa_consensus.ProposerSchedule,
	pool *transaction_pool.TransactionPool,
	node Node,
	blockchain Blockchain,
	sender *message_sender.MessageSender,
) *TransactionHandler {
	return &TransactionHandler{
		transactionValidator: validator,
		proposerSchedule:     schedule,
		txPool:               pool,
		node:                 node,
		bc:                   blockchain,
		sender:               sender,
	}
}

func (h *TransactionHandler) Handlers() map[string]func(p2p.Peer, *message.Message) error {
	return map[string]func(p2p.Peer, *message.Message) error{
		message.MessageSendTransaction: h.Transaction,
	}
}

func (h *TransactionHandler) Transaction(peer p2p.Peer, msg *message.Message) error {
	tx := msg.Payload.(*transaction.Transaction)
	logger.Info("[TransactionHandler] receive transaction", tx)
	err := h.transactionValidator.QuickVerify(*tx)
	if err != nil {
		logger.Error("Error when quick verify transaction", tx, err)
		return err
	}
	currentBlockNumber := h.bc.CurrentBlock().Header.Height

	// forward to next proposer or if is next proposer then import to pool
	nextProposer, err := h.proposerSchedule.GetProposer(
		currentBlockNumber + 2,
	) // + 1 is node producing, so it will be +2
	if err != nil {
		logger.Error("error when get next proposer", err)
		return err
	}

	if nextProposer == h.node.Address() {
		h.txPool.AddTransaction(tx)
	} else {
		err := h.sender.SendMessageToAddress(
			nextProposer,
			message.MessageSendTransaction,
			tx,
		)
		if err != nil {
			logger.Warn("error when send tx to proposer", err, "adding transaction to pool")
			h.txPool.AddTransaction(tx)
		}
	}

	return nil
}
