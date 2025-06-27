package handlers

import (
	"FichainCore/block"
	"FichainCore/common"
)

type Blockchain interface {
	CurrentBlock() *block.Block
}

type Node interface {
	Address() common.Address
}
