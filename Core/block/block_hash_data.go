package block

import (
	"math/big"

	"FichainCore/common"
	pb "FichainCore/proto"
)

type BlockHashData struct {
	Height           uint64
	ParentHash       common.Hash
	StateRoot        common.Hash
	TransactionsRoot common.Hash
	ReceiptRoot      common.Hash
	UncleHash        common.Hash
	Bloom            *big.Int
	Timestamp        uint64
	Prevrandao       common.Hash // <-- New field for PoS randomness
	Proposer         common.Address
	ExtraData        []byte // Optional
}

func (h *BlockHashData) Proto() *pb.BlockHashData {
	return &pb.BlockHashData{
		Height:           h.Height,
		ParentHash:       h.ParentHash.Bytes(),
		StateRoot:        h.StateRoot.Bytes(),
		TransactionsRoot: h.TransactionsRoot.Bytes(),
		ReceiptRoot:      h.ReceiptRoot.Bytes(),
		UncleHash:        h.UncleHash.Bytes(),
		Bloom:            h.Bloom.Bytes(),
		Timestamp:        h.Timestamp,
		Prevrandao:       h.Prevrandao.Bytes(),
		Proposer:         h.Proposer.Bytes(),
		ExtraData:        h.ExtraData,
	}
}

func (h *BlockHashData) FromProto(pbData *pb.BlockHashData) {
	h.Height = pbData.Height
	h.ParentHash = common.BytesToHash(pbData.ParentHash)
	h.StateRoot = common.BytesToHash(pbData.StateRoot)
	h.TransactionsRoot = common.BytesToHash(pbData.TransactionsRoot)
	h.ReceiptRoot = common.BytesToHash(pbData.ReceiptRoot)
	h.Bloom = new(big.Int).SetBytes(pbData.Bloom)
	h.UncleHash = common.BytesToHash(pbData.UncleHash)
	h.Timestamp = pbData.Timestamp
	h.Prevrandao = common.BytesToHash(pbData.Prevrandao)
	h.Proposer = common.BytesToAddress(pbData.Proposer)
	h.ExtraData = pbData.ExtraData
}
