package block_validator

import (
	"FichainCore/block"
	"FichainCore/common"
	"FichainCore/params"
)

type BlockChain interface {
	HasBlockAndState(hash common.Hash, number uint64) bool
	HasBlock(hash common.Hash, number uint64) bool
	Config() *params.ChainConfig
	CurrentHeader() *block.BlockHeader
	GetHeader(hash common.Hash, number uint64) *block.BlockHeader
	GetHeaderByNumber(number uint64) *block.BlockHeader
	GetHeaderByHash(hash common.Hash) *block.BlockHeader
	GetBlock(hash common.Hash, number uint64) *block.Block
}
