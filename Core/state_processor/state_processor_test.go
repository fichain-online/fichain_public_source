package state_processor_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"FichainCore/block"
	"FichainCore/block_chain"
	"FichainCore/common"
	"FichainCore/database"
	"FichainCore/evm"
	"FichainCore/params"
	"FichainCore/state"
	"FichainCore/state_processor"
	"FichainCore/transaction"
)

// mockBlock creates a simple block with a single transaction for testing.
func mockBlock() *block.Block {
	header := &block.BlockHeader{}
	tx := &transaction.Transaction{
		// Fill in minimal necessary data
	}
	return &block.Block{
		Header:       header,
		Transactions: []*transaction.Transaction{tx},
	}
}

// mockStateDB creates a dummy StateDB.
func mockStateDB() *state.StateDB {
	// You must return a usable state.StateDB, possibly with a mock backend

	db, _ := database.NewMemDatabase()
	stateDB, _ := state.New(common.Hash{}, state.NewDatabase(db))
	return stateDB
}

func TestStateProcessor_Process_EmptyBlock(t *testing.T) {
	config := &params.ChainConfig{}
	bc := &block_chain.BlockChain{}

	processor := state_processor.NewStateProcessor(config, bc)

	block := &block.Block{
		Transactions: []*transaction.Transaction{},
	}
	stateDB := mockStateDB()
	cfg := evm.Config{}

	receipts, logs, gasUsed, err := processor.Process(block, stateDB, cfg)

	assert.NoError(t, err)
	assert.Equal(t, 0, len(receipts))
	assert.Equal(t, 0, len(logs))
	assert.Equal(t, uint64(0), gasUsed)
}
