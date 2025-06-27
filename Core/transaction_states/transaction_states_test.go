package transaction_states_test

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"FichainCore/common"
	"FichainCore/mpt"
	"FichainCore/storage"
	"FichainCore/transaction"
	"FichainCore/transaction_states"
)

func mockTransaction() *transaction.Transaction {
	return transaction.NewTransaction(
		common.HexToAddress("0x1234567890abcdef1234567890abcdef12345678"),
		big.NewInt(1),
		big.NewInt(1000),
		[]byte("test-data"),
		21000,
		100,
		"test message",
	)
}

func TestTransactionStates_SetGetDeleteTransaction(t *testing.T) {
	memStorage := storage.NewMemoryDb()
	trie, err := mpt.New(common.Hash{}, memStorage)
	assert.NoError(t, err)

	txState := transaction_states.NewTransactionStates(trie)

	tx := mockTransaction()
	txIndex := uint32(1)

	// Test SetTransaction
	err = txState.SetTransaction(txIndex, tx)
	assert.NoError(t, err, "SetTransaction should not return an error")

	// Test GetTransaction
	gotTx, err := txState.GetTransaction(txIndex)
	assert.NoError(t, err, "GetTransaction should not return an error")
	assert.NotNil(t, gotTx, "Transaction should not be nil")
	assert.Equal(t, tx.To(), gotTx.To())
	assert.Equal(t, tx.Nonce().Cmp(gotTx.Nonce()), 0)
	assert.Equal(t, tx.Amount().Cmp(gotTx.Amount()), 0)
	assert.Equal(t, tx.Message(), gotTx.Message())
	assert.Equal(t, tx.Data(), gotTx.Data())

	// Test DeleteTransaction
	err = txState.DeleteTransaction(txIndex)
	assert.NoError(t, err, "DeleteTransaction should not return an error")

	// Test GetTransaction after delete
	gotTx, err = txState.GetTransaction(txIndex)
	assert.Nil(t, gotTx, "Transaction should be nil after delete")
}

func TestTransactionStates_Commit(t *testing.T) {
	memStorage := storage.NewMemoryDb()
	trie, err := mpt.New(common.Hash{}, memStorage)
	assert.NoError(t, err)

	txState := transaction_states.NewTransactionStates(trie)

	tx := mockTransaction()
	txIndex := uint32(2)
	err = txState.SetTransaction(txIndex, tx)
	assert.NoError(t, err)

	rootHash, err := txState.Commit(true, memStorage)
	assert.NoError(t, err, "Commit should not return an error")
	assert.NotEqual(t, common.Hash{}, rootHash, "Root hash should not be zero")

	// Ensure trie is updated with new root
	assert.Equal(t, rootHash, txState.Hash(), "Hash should match committed root hash")
}
