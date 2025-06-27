package types

import (
	"math/big"

	"FichainCore/common"
	pb "FichainCore/proto"
)

type Transaction interface {
	Hash() common.Hash
	HashSign(chainId *big.Int) (common.Hash, error)
	From(chainId *big.Int) (common.Address, error)
	To() common.Address
	Nonce() uint64
	Amount() *big.Int
	Data() []byte
	MaxGas() uint64
	MaxGasPrice() uint64
	Message() string
	SetSign(Sign)

	Marshal() ([]byte, error)
	Unmarshal([]byte) error
	Proto() *pb.Transaction
	FromProto(*pb.Transaction)
	String() string
}
