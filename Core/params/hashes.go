package params

import (
	"FichainCore/common"
	"FichainCore/crypto"
)

var (
	EmptyRootHash = common.HexToHash(
		"56e81f171bcc55a6ff8345e692c0f86e5b48e01b996cadc001622fb5e363b421",
	)
	EmptyCodeHash = crypto.Keccak256Hash(
		nil,
	) // c5d2460186f7233c927e7db2dcc703c0e500b653ca82273b7bfad8045d85a470

)
