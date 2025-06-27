package poa_consensus

import (
	"math/big"

	"FichainCore/block"
	"FichainCore/common"
	"FichainCore/consensus"
	"FichainCore/receipt"
	"FichainCore/state"
	"FichainCore/transaction"
)

type POAConsensus struct{}

func (c *POAConsensus) Author(header *block.BlockHeader) (common.Address, error) {
	// TODO
	return header.Proposer, nil
}

// VerifyHeader checks whether a header conforms to the consensus rules of a
// given engine. Verifying the seal may be done optionally here, or explicitly
// via the VerifySeal method.
func (c *POAConsensus) VerifyHeader(
	chain consensus.ChainReader,
	header *block.BlockHeader,
	seal bool,
) error {
	// TODO
	return nil
}

// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
// concurrently. The method returns a quit channel to abort the operations and
// a results channel to retrieve the async verifications (the order is that of
// the input slice).
func (c *POAConsensus) VerifyHeaders(
	chain consensus.ChainReader,
	headers []*block.BlockHeader,
	seals []bool,
) (chan<- struct{}, <-chan error) {
	// TODO
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

// VerifyUncles verifies that the given block's uncles conform to the consensus
// rules of a given engine.
func (c *POAConsensus) VerifyUncles(chain consensus.ChainReader, block *block.Block) error {
	// TODO
	return nil
}

// VerifySeal checks whether the crypto seal on a header is valid according to
// the consensus rules of the given engine.
func (c *POAConsensus) VerifySeal(chain consensus.ChainReader, header *block.BlockHeader) error {
	// TODO
	return nil
}

// Prepare initializes the consensus fields of a block header according to the
// rules of a particular engine. The changes are executed inline.
func (c *POAConsensus) Prepare(chain consensus.ChainReader, header *block.BlockHeader) error {
	// TODO
	return nil
}

// Finalize runs any post-transaction state modifications (e.g. block rewards)
// and assembles the final block.
// Note: The block header and state database might be updated to reflect any
// consensus rules that happen at finalization (e.g. block rewards).
func (c *POAConsensus) Finalize(
	chain consensus.ChainReader,
	header *block.BlockHeader,
	state *state.StateDB,
	txs []*transaction.Transaction,
	uncles []*block.BlockHeader,
	receipts []*receipt.Receipt,
) (*block.Block, error) {
	// TODO
	return block.NewBlock(header, txs, uncles, receipts), nil
}

// Seal generates a new block for the given input block with the local miner's
// seal place on top.
func (c *POAConsensus) Seal(
	chain consensus.ChainReader,
	block *block.Block,
	stop <-chan struct{},
) (*block.Block, error) {
	// TODO
	return nil, nil
}

// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
// that a new block should have.
func (c *POAConsensus) CalcDifficulty(
	chain consensus.ChainReader,
	time uint64,
	parent *block.BlockHeader,
) *big.Int {
	// TODO
	return big.NewInt(0)
}
