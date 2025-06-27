package block

import (
	"google.golang.org/protobuf/proto"

	pb "FichainCore/proto"
	"FichainCore/transaction"
)

type Body struct {
	Transactions []*transaction.Transaction
	Uncles       []*BlockHeader
}

// Marshal serializes the block to protobuf bytes
func (b *Body) Marshal() ([]byte, error) {
	return proto.Marshal(b.Proto())
}

// Unmarshal deserializes protobuf bytes into a block
func (b *Body) Unmarshal(data []byte) error {
	var pbBody pb.Body
	err := proto.Unmarshal(data, &pbBody)
	if err != nil {
		return err
	}
	b.FromProto(&pbBody)
	return nil
}

// Proto converts Body to protobuf message
func (b *Body) Proto() *pb.Body {
	pbTxs := make([]*pb.Transaction, len(b.Transactions))
	for i, tx := range b.Transactions {
		pbTxs[i] = tx.Proto().(*pb.Transaction)
	}
	pbUncles := make([]*pb.BlockHeader, len(b.Uncles))
	return &pb.Body{
		Txns:   pbTxs,
		Uncles: pbUncles,
	}
}

// FromProto populates the Body from protobuf message
func (b *Body) FromProto(pbBody *pb.Body) {
	txs := make([]*transaction.Transaction, len(pbBody.Txns))
	for i, pbTx := range pbBody.Txns {
		tx := &transaction.Transaction{}
		tx.FromProto(pbTx)
		txs[i] = tx
	}

	uncles := make([]*BlockHeader, len(pbBody.Uncles))
	for i, pbU := range pbBody.Uncles {
		u := &BlockHeader{}
		u.FromProto(pbU)
		uncles[i] = u
	}

	b.Uncles = uncles
	b.Transactions = txs
}
