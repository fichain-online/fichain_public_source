package types

import (
	"github.com/ethereum/go-ethereum/crypto"

	"FichainCore/common"
)

var (
	// EmptyRootHash is the known root hash of an empty merkle trie.
	EmptyRootHash = common.HexToHash(
		"56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
	)
	// TODO
	EmptyUncleHash = common.Hash{}

	// EmptyCodeHash is the known hash of the empty EVM bytecode.
	EmptyCodeHash = crypto.Keccak256Hash(
		nil,
	) // c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470

)
