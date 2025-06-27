package block

import (
	"fmt"
	"log/slog"
	"time"

	"google.golang.org/protobuf/proto"

	"FichainCore/bloom"
	"FichainCore/common"
	"FichainCore/crypto"
	"FichainCore/params"
	pb "FichainCore/proto"
	"FichainCore/receipt"
	"FichainCore/transaction"
	"FichainCore/trie"
	"FichainCore/types"
)

type Block struct {
	Header       *BlockHeader
	Transactions []*transaction.Transaction
	Uncles       []*BlockHeader
}

func NewBlock(
	header *BlockHeader,
	txs []*transaction.Transaction,
	uncles []*BlockHeader,
	receipts []*receipt.Receipt,
) *Block {
	b := &Block{
		Header:       header,
		Transactions: txs,
	}

	if len(txs) == 0 {
		b.Header.TransactionsRoot = types.EmptyRootHash
	} else {
		var err error
		b.Header.TransactionsRoot, err = trie.DeriveSha(txs)
		if err != nil {
			panic(err.Error)
		}
		b.Transactions = make([]*transaction.Transaction, len(txs))
		copy(b.Transactions, txs)
	}

	if len(receipts) == 0 {
		b.Header.ReceiptRoot = types.EmptyRootHash
	} else {
		var err error
		b.Header.ReceiptRoot, err = trie.DeriveSha(receipts)
		if err != nil {
			panic(err)
		}
		b.Header.Bloom = bloom.CreateBloom(receipts).Big()
	}

	b.Header.UncleHash = CalcUncleHash(uncles)
	b.Uncles = uncles
	return b
}

func (b *Block) Hash() common.Hash {
	hashDataBytes, err := proto.Marshal(b.Header.ToHashData().Proto())
	if err != nil {
		slog.Error("Error marshaling block header for hash", "err", err)
		return common.Hash{}
	}
	return crypto.Keccak256Hash(hashDataBytes)
}

// Marshal serializes the block to protobuf bytes
func (b *Block) Marshal() ([]byte, error) {
	return proto.Marshal(b.Proto())
}

// Unmarshal deserializes protobuf bytes into a block
func (b *Block) Unmarshal(data []byte) error {
	var pbBlock pb.Block
	err := proto.Unmarshal(data, &pbBlock)
	if err != nil {
		return err
	}
	b.FromProto(&pbBlock)
	return nil
}

// Proto converts Block to protobuf message
func (b *Block) Proto() *pb.Block {
	pbTxs := make([]*pb.Transaction, len(b.Transactions))
	for i, tx := range b.Transactions {
		pbTxs[i] = tx.Proto().(*pb.Transaction)
	}
	return &pb.Block{
		Header: b.Header.Proto(),
		Txns:   pbTxs,
	}
}

// FromProto populates the Block from protobuf message
func (b *Block) FromProto(pbBlock *pb.Block) {
	header := &BlockHeader{}
	header.FromProto(pbBlock.Header)

	txs := make([]*transaction.Transaction, len(pbBlock.Txns))
	for i, pbTx := range pbBlock.Txns {
		tx := &transaction.Transaction{}
		tx.FromProto(pbTx)
		txs[i] = tx
	}

	b.Header = header
	b.Transactions = txs
}

func (b *Block) GasLimit() uint64 {
	// improve later
	return params.TempGasLimit
}

func (b *Block) Body() *Body {
	return &Body{
		Transactions: b.Transactions,
		Uncles:       b.Uncles,
	}
}

// --- Stringer (optional) ---
func (b *Block) String() string {
	return fmt.Sprintf("Block %d | TxCount: %d | Hash: %s | Time: %s",
		b.Header.Height,
		len(b.Transactions),
		b.Hash().Hex(),
		time.Unix(int64(b.Header.Timestamp), 0).Format(time.RFC3339),
	)
}
