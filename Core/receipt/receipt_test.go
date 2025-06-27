package receipt

import (
	"bytes"
	"encoding/hex"
	"testing"

	"FichainCore/common"
	"FichainCore/log"
)

func TestReceiptHash(t *testing.T) {
	// Sample log
	testLog := &log.Log{
		Address: common.HexToAddress("0x000000000000000000000000000000000000dead"),
		Topics: []common.Hash{
			common.HexToHash("0x01"),
			common.HexToHash("0x02"),
		},
		Data:     []byte("example event data"),
		LogIndex: 0,
	}

	// Receipt with all fields
	toAddress := common.HexToAddress("0x000000000000000000000000000000000000beef")
	contractAddress := common.HexToAddress("0x000000000000000000000000000000000000cafe")

	r := &Receipt{
		TxHash: common.HexToHash(
			"0xaaaabbbbccccddddeeeeffff0000111122223333444455556666777788889999",
		),
		BlockHash: common.HexToHash(
			"0xbbbbccccddddeeeeffff0000111122223333444455556666777788889999aaaa",
		),
		BlockNumber:       123456,
		TxIndex:           5,
		From:              common.HexToAddress("0x000000000000000000000000000000000000face"),
		To:                &toAddress,
		CumulativeGasUsed: 21000,
		GasUsed:           21000,
		ContractAddress:   &contractAddress,
		Logs:              []*log.Log{testLog},
		LogsBloom:         make([]byte, 256),
		Status:            1,
	}

	hash1 := r.Hash()
	hash2 := r.Hash() // Hash must be deterministic

	if !bytes.Equal(hash1.Bytes(), hash2.Bytes()) {
		t.Errorf("Expected deterministic hash, got different results:\nHash1: %s\nHash2: %s",
			hex.EncodeToString(hash1.Bytes()), hex.EncodeToString(hash2.Bytes()))
	}

	// Optionally log the hash
	t.Logf("Receipt Hash: %s", hash1.Hex())
}
