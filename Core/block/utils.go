package block

import (
	"github.com/ethereum/go-ethereum/rlp"

	"FichainCore/common"
	"FichainCore/crypto/sha3"
)

func rlpHash(x interface{}) (h common.Hash) {
	hw := sha3.NewKeccak256()
	rlp.Encode(hw, x)
	hw.Sum(h[:0])
	return h
}

func CalcUncleHash(uncles []*BlockHeader) common.Hash {
	return rlpHash(uncles)
}
