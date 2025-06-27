package block

import (
	"sort"
)

type BlockBy func(b1, b2 *Block) bool

func (self BlockBy) Sort(blocks []*Block) {
	bs := blockSorter{
		blocks: blocks,
		by:     self,
	}
	sort.Sort(bs)
}

type blockSorter struct {
	blocks []*Block
	by     func(b1, b2 *Block) bool
}

func (self blockSorter) Len() int { return len(self.blocks) }
func (self blockSorter) Swap(i, j int) {
	self.blocks[i], self.blocks[j] = self.blocks[j], self.blocks[i]
}
func (self blockSorter) Less(i, j int) bool { return self.by(self.blocks[i], self.blocks[j]) }

func Number(b1, b2 *Block) bool { return b1.Header.Height < b2.Header.Height }
