package event

import (
	"FichainCore/block"
	"FichainCore/common"
	"FichainCore/log"
)

type ChainEvent struct {
	Block *block.Block
	Hash  common.Hash
	Logs  []*log.Log
}

type ChainSideEvent struct {
	Block *block.Block
}

type ChainHeadEvent struct {
	Block *block.Block
}

// RemovedLogsEvent is posted when a reorg happens
type RemovedLogsEvent struct{ Logs []*log.Log }
