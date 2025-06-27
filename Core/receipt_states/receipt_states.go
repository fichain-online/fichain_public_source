package receipt_states

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"

	logger "github.com/HendrickPhan/golang-simple-logger"
	"google.golang.org/protobuf/proto"

	"FichainCore/common"
	"FichainCore/mpt"
	pb "FichainCore/proto"
	"FichainCore/receipt"
	"FichainCore/storage"
)

type ReceiptStates struct {
	trie *mpt.MerklePatriciaTrie
}

func NewReceiptStates(trie *mpt.MerklePatriciaTrie) *ReceiptStates {
	return &ReceiptStates{trie: trie}
}

func (rs *ReceiptStates) GetReceipt(txIndex uint32) (*receipt.ReceiptData, error) {
	key := encodeTxIndexKey(txIndex)
	value, err := rs.trie.Get(key)
	if err != nil || value == nil {
		return nil, err
	}

	var pbReceipt pb.ReceiptData
	if err := proto.Unmarshal(value, &pbReceipt); err != nil {
		return nil, err
	}

	r := &receipt.ReceiptData{}
	r.FromProto(&pbReceipt)
	return r, nil
}

func (rs *ReceiptStates) SetReceipt(txIndex uint32, r *receipt.ReceiptData) error {
	key := encodeTxIndexKey(txIndex)
	data, err := proto.Marshal(r.Proto())
	if err != nil {
		return err
	}
	return rs.trie.Update(key, data)
}

func (rs *ReceiptStates) DeleteReceipt(txIndex uint32) error {
	key := encodeTxIndexKey(txIndex)
	return rs.trie.Delete(key)
}

func (rs *ReceiptStates) Commit(collectLeaf bool, db storage.Storage) (common.Hash, error) {
	rootHash, nodeSet, oldKeys, err := rs.trie.Commit(collectLeaf)
	if err != nil {
		return common.Hash{}, fmt.Errorf("trie commit error: %w", err)
	}

	// Step 3: Persist nodeSet to database
	if nodeSet != nil && db != nil {
		batch := [][2][]byte{}
		for _, node := range nodeSet.Nodes {
			logger.DebugP("saving node", node.Hash)
			batch = append(batch, [2][]byte{node.Hash.Bytes(), node.Blob})
		}
		if err := db.BatchPut(batch); err != nil {
			return common.Hash{}, fmt.Errorf("batch put error: %w", err)
		}
	}

	// Optional: delete oldKeys (for pruning/archive)
	for _, key := range oldKeys {
		logger.DebugP("Old keys", hex.EncodeToString(key))
		if err := db.Delete(key); err != nil {
			return common.Hash{}, fmt.Errorf("failed to delete old key: %w", err)
		}
	}

	// Step 5: Replace old trie with new trie instance
	newTrie, err := mpt.New(rootHash, db)
	if err != nil {
		return common.Hash{}, fmt.Errorf("failed to create new trie: %w", err)
	}
	rs.trie = newTrie

	return rootHash, nil
}

func (rs *ReceiptStates) Hash() common.Hash {
	return rs.trie.Hash()
}

func encodeTxIndexKey(index uint32) []byte {
	// Use 4-byte big endian encoding
	key := make([]byte, 4)
	binary.BigEndian.PutUint32(key, index)
	return key
}
