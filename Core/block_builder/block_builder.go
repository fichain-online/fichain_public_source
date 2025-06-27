package block_builder

import (
	"math/big"
	"time"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/block"
	"FichainCore/block_chain"
	"FichainCore/bloom"
	"FichainCore/common"
	"FichainCore/gas_pool"
	"FichainCore/log"
	"FichainCore/params"
	"FichainCore/receipt"
	"FichainCore/state"
	"FichainCore/state_processor"
	"FichainCore/transaction"
	"FichainCore/transaction_pool"
	"FichainCore/transaction_validator"
	"FichainCore/trie"
)

type BlockBuilder struct {
	transactionPool      *transaction_pool.TransactionPool
	transactionValidator *transaction_validator.TransactionValidator
	coinBase             *common.Address

	bc            *block_chain.BlockChain
	statedb       *state.StateDB
	gasPool       *gas_pool.GasPool
	currentHeader *block.BlockHeader

	txs      []*transaction.Transaction
	receipts []*receipt.Receipt
}

func NewBlockBuilder(
	transactionPool *transaction_pool.TransactionPool,
	transactionValidator *transaction_validator.TransactionValidator,
	coinBase *common.Address,
	bc *block_chain.BlockChain,
	statedb *state.StateDB,
) *BlockBuilder {
	return &BlockBuilder{
		transactionPool:      transactionPool,
		transactionValidator: transactionValidator,
		coinBase:             coinBase,
		bc:                   bc,
		statedb:              statedb,
	}
}

func (bb *BlockBuilder) applyTransaction(
	tx *transaction.Transaction,
	txIndex int,
) (error, []*log.Log) {
	snap := bb.statedb.Snapshot()
	logger.Warn("TODO update mainnet chain config")

	bb.statedb.Prepare(tx.Hash(), bb.currentHeader.ParentHash, txIndex)
	logger.Debug("GetVmConfig", bb.bc.GetVmConfig())
	receipt, _, err := state_processor.ApplyTransaction(
		params.MainnetChainConfig,
		bb.bc,
		bb.coinBase,
		bb.gasPool,
		bb.statedb,
		bb.currentHeader,
		tx,
		&bb.currentHeader.GasUsed,
		bb.bc.GetVmConfig(),
	)
	if err != nil {
		bb.statedb.RevertToSnapshot(snap)
		return err, nil
	}
	bb.txs = append(bb.txs, tx)
	bb.receipts = append(bb.receipts, receipt)

	return nil, receipt.Logs
}

func (bb *BlockBuilder) GenerateBlock(
	lastBlockHeader *block.BlockHeader,
) (*block.Block, error) {
	bb.statedb.Reset(lastBlockHeader.StateRoot)

	// TODO: initial block header
	bb.currentHeader = &block.BlockHeader{
		Height:     lastBlockHeader.Height + 1,
		ParentHash: lastBlockHeader.Hash(),
		Proposer:   *bb.coinBase,
		Timestamp:  uint64(time.Now().Unix()),
		Prevrandao: common.BigToHash(big.NewInt(time.Now().Unix())),
	}

	bl := &block.Block{
		Header: bb.currentHeader,
	}
	bb.txs = make([]*transaction.Transaction, 0)
	bb.receipts = make([]*receipt.Receipt, 0)
	gp := gas_pool.GasPool(params.TempGasLimit)
	bb.gasPool = &gp

	// take all transaction from pool by address and nonce
	mAddressTxs := bb.transactionPool.GetPendingTransactions()
	bb.transactionPool.Clear()
	txIndex := 0
	for _, txs := range mAddressTxs {
		// deep verify transaction
		for _, tx := range txs {
			err := bb.transactionValidator.DeepVerify(tx)
			// apply transaction
			if err != nil {
				//
				logger.Warn("error when deep verify transaction", tx, "skiped")
				continue
			}

			err, _ = bb.applyTransaction(
				tx,
				txIndex,
			)
			if err != nil {
				logger.Warn("error when apply transaction", err, tx, "skiped")
				continue
			}
			txIndex++
		}
	}
	bl.Transactions = bb.txs
	unclesHash := block.CalcUncleHash(bl.Uncles)
	bl.Header.UncleHash = unclesHash
	var err error
	//
	bl.Header.TransactionsRoot, err = trie.DeriveSha(bl.Transactions)
	if err != nil {
		logger.DebugP("Error 1")
		return nil, err
	}
	//
	bl.Header.ReceiptRoot, err = trie.DeriveSha(bb.receipts)
	if err != nil {
		logger.DebugP("Error 2")
		return nil, err
	}
	bl.Header.StateRoot = bb.statedb.IntermediateRoot(true)
	//
	rbloom := bloom.CreateBloom(bb.receipts)
	bl.Header.Bloom = new(big.Int).SetBytes(rbloom.Bytes())

	return bl, nil
}
