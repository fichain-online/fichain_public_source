package block

import (
	"encoding/hex"
	"fmt"
	"math/big"

	logger "github.com/HendrickPhan/golang-simple-logger"
	"google.golang.org/protobuf/proto"

	"FichainCore/common"
	"FichainCore/crypto"
	pb "FichainCore/proto"
)

type BlockHeader struct {
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
	Signature        []byte
	ExtraData        []byte
	GasUsed          uint64 // <-- New field
}

// --- BlockHeader conversion functions ---

func (h *BlockHeader) Proto() *pb.BlockHeader {
	if h.Bloom == nil {
		h.Bloom = big.NewInt(0)
	}
	return &pb.BlockHeader{
		Height:           h.Height,
		ParentHash:       h.ParentHash.Bytes(),
		StateRoot:        h.StateRoot.Bytes(),
		TransactionsRoot: h.TransactionsRoot.Bytes(),
		ReceiptRoot:      h.ReceiptRoot.Bytes(),
		UncleHash:        h.UncleHash.Bytes(),
		Bloom:            h.Bloom.Bytes(),
		Timestamp:        h.Timestamp,
		Prevrandao:       h.Prevrandao.Bytes(), // <-- New
		Proposer:         h.Proposer.Bytes(),
		Signature:        h.Signature,
		ExtraData:        h.ExtraData,
	}
}

func (h *BlockHeader) FromProto(pbHeader *pb.BlockHeader) {
	h.Height = pbHeader.Height
	h.ParentHash = common.BytesToHash(pbHeader.ParentHash)
	h.StateRoot = common.BytesToHash(pbHeader.StateRoot)
	h.TransactionsRoot = common.BytesToHash(pbHeader.TransactionsRoot)
	h.ReceiptRoot = common.BytesToHash(pbHeader.ReceiptRoot)
	h.UncleHash = common.BytesToHash(pbHeader.UncleHash)
	h.Bloom = new(big.Int).SetBytes(pbHeader.Bloom)
	h.Timestamp = pbHeader.Timestamp
	h.Prevrandao = common.BytesToHash(pbHeader.Prevrandao) // <-- New
	h.Proposer = common.BytesToAddress(pbHeader.Proposer)
	h.Signature = pbHeader.Signature
	h.ExtraData = pbHeader.ExtraData
}

// ToHashData converts BlockHeader to BlockHashData (used for hash calculation)
func (h *BlockHeader) ToHashData() *BlockHashData {
	var bloom *big.Int
	if h.Bloom == nil {
		bloom = big.NewInt(0)
	} else {
		bloom = h.Bloom
	}
	return &BlockHashData{
		Height:           h.Height,
		ParentHash:       h.ParentHash,
		StateRoot:        h.StateRoot,
		TransactionsRoot: h.TransactionsRoot,
		ReceiptRoot:      h.ReceiptRoot,
		UncleHash:        h.UncleHash,
		Bloom:            bloom,
		Timestamp:        h.Timestamp,
		Prevrandao:       h.Prevrandao, // <-- New
		Proposer:         h.Proposer,
		ExtraData:        h.ExtraData,
	}
}

func (b *BlockHeader) Hash() common.Hash {
	hashDataBytes, err := proto.Marshal(b.ToHashData().Proto())
	if err != nil {
		logger.Error("Error marshaling block header for hash", "err", err)
		return common.Hash{}
	}
	return crypto.Keccak256Hash(hashDataBytes)
}

func (h *BlockHeader) String() string {
	return fmt.Sprintf(
		`BlockHeader {
  Height:           %d
  ParentHash:       %s
  StateRoot:        %s
  TransactionsRoot: %s
  ReceiptRoot:      %s
  UncleHash:        %s
  Bloom:            %s
  Timestamp:        %d
  Prevrandao:       %s
  Proposer:         %s
  Signature:        %s
  ExtraData:        %s
  GasUsed:          %d
}`,
		h.Height,
		h.ParentHash.String(),
		h.StateRoot.String(),
		h.TransactionsRoot.String(),
		h.ReceiptRoot.String(),
		h.UncleHash.String(),
		h.Bloom.Text(16),
		h.Timestamp,
		h.Prevrandao.String(),
		h.Proposer.String(),
		hex.EncodeToString(h.Signature),
		hex.EncodeToString(h.ExtraData),
		h.GasUsed,
	)
}
