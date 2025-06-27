package consensus

import (
	"math/big"

	"FichainCore/block"
	"FichainCore/common"
	"FichainCore/params"
	"FichainCore/receipt"
	"FichainCore/state"
	"FichainCore/transaction"
)

// ChainReader defines a small collection of methods needed to access the local
// blockchain during header and/or uncle verification.
type ChainReader interface {
	// Config retrieves the blockchain's chain configuration.
	Config() *params.ChainConfig

	// CurrentHeader retrieves the current header from the local chain.
	CurrentHeader() *block.BlockHeader

	// GetHeader retrieves a block header from the database by hash and number.
	GetHeader(hash common.Hash, number uint64) *block.BlockHeader

	// GetHeaderByNumber retrieves a block header from the database by number.
	GetHeaderByNumber(number uint64) *block.BlockHeader

	// GetHeaderByHash retrieves a block header from the database by its hash.
	GetHeaderByHash(hash common.Hash) *block.BlockHeader

	// GetBlock retrieves a block from the database by hash and number.
	GetBlock(hash common.Hash, number uint64) *block.Block
}

// Engine is an algorithm agnostic consensus engine.
type Engine interface {
	// Author retrieves the Ethereum address of the account that minted the given
	// block, which may be different from the header's coinbase if a consensus
	// engine is based on signatures.
	Author(header *block.BlockHeader) (common.Address, error)

	// VerifyHeader checks whether a header conforms to the consensus rules of a
	// given engine. Verifying the seal may be done optionally here, or explicitly
	// via the VerifySeal method.
	VerifyHeader(chain ChainReader, header *block.BlockHeader, seal bool) error

	// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
	// concurrently. The method returns a quit channel to abort the operations and
	// a results channel to retrieve the async verifications (the order is that of
	// the input slice).
	VerifyHeaders(
		chain ChainReader,
		headers []*block.BlockHeader,
		seals []bool,
	) (chan<- struct{}, <-chan error)

	// VerifyUncles verifies that the given block's uncles conform to the consensus
	// rules of a given engine.
	VerifyUncles(chain ChainReader, block *block.Block) error

	// VerifySeal checks whether the crypto seal on a header is valid according to
	// the consensus rules of the given engine.
	VerifySeal(chain ChainReader, header *block.BlockHeader) error

	// Prepare initializes the consensus fields of a block header according to the
	// rules of a particular engine. The changes are executed inline.
	Prepare(chain ChainReader, header *block.BlockHeader) error

	// Finalize runs any post-transaction state modifications (e.g. block rewards)
	// and assembles the final block.
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at finalization (e.g. block rewards).
	Finalize(
		chain ChainReader,
		header *block.BlockHeader,
		state *state.StateDB,
		txs []*transaction.Transaction,
		uncles []*block.BlockHeader,
		receipts []*receipt.Receipt,
	) (*block.Block, error)

	// Seal generates a new block for the given input block with the local miner's
	// seal place on top.
	Seal(chain ChainReader, block *block.Block, stop <-chan struct{}) (*block.Block, error)

	// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
	// that a new block should have.
	CalcDifficulty(chain ChainReader, time uint64, parent *block.BlockHeader) *big.Int

	// APIs returns the RPC APIs this consensus engine provides.
	// APIs(chain ChainReader) []rpc.API
}
