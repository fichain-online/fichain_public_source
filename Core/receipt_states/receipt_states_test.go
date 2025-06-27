package receipt_states_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"FichainCore/common"
	"FichainCore/log"
	"FichainCore/mpt"
	"FichainCore/receipt"
	"FichainCore/receipt_states"
	"FichainCore/storage"
)

func TestReceiptStates(t *testing.T) {
	db := storage.NewMemoryDb()
	trie, err := mpt.New(common.Hash{}, db)
	assert.NoError(t, err)

	rs := receipt_states.NewReceiptStates(trie)

	// Create a fake receipt
	index := uint32(0)
	toAddr := common.HexToAddress("0x2222222222222222222222222222222222222222")
	contractAddr := common.HexToAddress("0x8888888888888888888888888888888888888888")

	receipt1 := &receipt.Receipt{
		TxHash:            common.HexToHash("0xabc"),
		BlockHash:         common.HexToHash("0xdef"),
		BlockNumber:       123,
		TxIndex:           index,
		From:              common.HexToAddress("0x1111111111111111111111111111111111111111"),
		To:                &toAddr,
		CumulativeGasUsed: 21000,
		GasUsed:           21000,
		ContractAddress:   &contractAddr,
		Logs: []*log.Log{
			{
				Address: toAddr,
				Topics:  []common.Hash{common.HexToHash("0xbeef")},
				Data:    []byte("log data"),
			},
		},
		LogsBloom: make([]byte, 256),
		Status:    1,
	}

	// SET
	err = rs.SetReceipt(index, receipt1.Data())
	assert.NoError(t, err)

	// GET
	got, err := rs.GetReceipt(index)
	assert.NoError(t, err)
	assert.Equal(t, len(receipt1.Logs), len(got.Logs))
	assert.Equal(t, receipt1.Logs[0].Data, got.Logs[0].Data)

	// DELETE
	err = rs.DeleteReceipt(index)
	assert.NoError(t, err)

	// Try to get after delete
	got, err = rs.GetReceipt(index)
	// assert.Error(t, err)
	assert.Nil(t, got)

	// Set again for commit test
	err = rs.SetReceipt(index, receipt1.Data())
	assert.NoError(t, err)

	// COMMIT
	rootHash, err := rs.Commit(true, db)
	assert.NoError(t, err)
	assert.NotEqual(t, common.Hash{}, rootHash)

	// Ensure trie is still working after commit
	got2, err := rs.GetReceipt(index)
	assert.NoError(t, err)
	assert.Equal(t, receipt1.Data(), got2)
}
