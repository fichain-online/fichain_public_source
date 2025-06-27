package fake_consensus

import (
	"math/big"

	"FichainCore/block"
	"FichainCore/common"
	"FichainCore/consensus"
	"FichainCore/receipt"
	"FichainCore/state"
	"FichainCore/transaction"
)

type FakeConsensus struct{}

func NewFakeConsensus() *FakeConsensus {
	return &FakeConsensus{}
}

func (fc *FakeConsensus) Author(header *block.BlockHeader) (common.Address, error) {
	// Return the coinbase as the author for simplicity
	return header.Proposer, nil
}

func (fc *FakeConsensus) VerifyHeader(
	chain consensus.ChainReader,
	header *block.BlockHeader,
	seal bool,
) error {
	// Skip verification
	return nil
}

func (fc *FakeConsensus) VerifyHeaders(
	chain consensus.ChainReader,
	headers []*block.BlockHeader,
	seals []bool,
) (chan<- struct{}, <-chan error) {
	abort := make(chan struct{})
	results := make(chan error, len(headers))

	go func() {
		defer close(results)
		for range headers {
			select {
			case <-abort:
				return
			default:
				results <- nil // always succeed
			}
		}
	}()

	return abort, results
}

func (fc *FakeConsensus) VerifyUncles(chain consensus.ChainReader, blk *block.Block) error {
	// Skip verification
	return nil
}

func (fc *FakeConsensus) VerifySeal(chain consensus.ChainReader, header *block.BlockHeader) error {
	// Skip seal verification
	return nil
}

func (fc *FakeConsensus) Prepare(chain consensus.ChainReader, header *block.BlockHeader) error {
	// Set dummy difficulty for testing
	header.Prevrandao = common.Hash{0x01}
	return nil
}

func (fc *FakeConsensus) Finalize(
	chain consensus.ChainReader,
	header *block.BlockHeader,
	statedb *state.StateDB,
	txs []*transaction.Transaction,
	uncles []*block.BlockHeader,
	receipts []*receipt.Receipt,
) (*block.Block, error) {
	// Just construct and return the block with all inputs
	return block.NewBlock(header, txs, uncles, receipts), nil
}

func (fc *FakeConsensus) Seal(
	chain consensus.ChainReader,
	blk *block.Block,
	stop <-chan struct{},
) (*block.Block, error) {
	// Return the block without actual sealing
	return blk, nil
}

func (fc *FakeConsensus) CalcDifficulty(
	chain consensus.ChainReader,
	time uint64,
	parent *block.BlockHeader,
) *big.Int {
	// Always return constant difficulty
	return big.NewInt(1)
}
