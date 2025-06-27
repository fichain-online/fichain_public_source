package block_validator

import (
	"fmt"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/block"
	"FichainCore/bloom"
	"FichainCore/consensus"
	"FichainCore/errors"
	"FichainCore/params"
	"FichainCore/receipt"
	"FichainCore/state"
	"FichainCore/trie"
)

// BlockValidator is responsible for validating block headers, uncles and
// processed state.
//
// BlockValidator implements Validator.
type BlockValidator struct {
	config *params.ChainConfig // Chain configuration options
	bc     BlockChain          // Canonical block chain
	engine consensus.Engine    // Consensus engine used for validating
}

// NewBlockValidator returns a new block validator which is safe for re-use
func NewBlockValidator(
	config *params.ChainConfig,
	blockchain BlockChain,
	engine consensus.Engine,
) *BlockValidator {
	validator := &BlockValidator{
		config: config,
		engine: engine,
		bc:     blockchain,
	}
	return validator
}

// ValidateBody validates the given block's uncles and verifies the the block
// header's transaction and uncle roots. The headers are assumed to be already
// validated at this point.
func (v *BlockValidator) ValidateBody(bl *block.Block) error {
	// Check whether the block's known, and if not, that it's linkable
	if v.bc.HasBlockAndState(bl.Hash(), uint64(bl.Header.Height)) {
		return errors.ErrKnownBlock
	}
	if !v.bc.HasBlockAndState(bl.Header.ParentHash, uint64(bl.Header.Height-1)) {
		if !v.bc.HasBlock(bl.Header.ParentHash, uint64(bl.Header.Height-1)) {
			return errors.ErrUnknownAncestor
		}
		return errors.ErrPrunedAncestor
	}
	// Header validity is known at this point, check the uncles and transactions
	header := bl.Header
	if err := v.engine.VerifyUncles(v.bc, bl); err != nil {
		return err
	}
	if hash := block.CalcUncleHash(bl.Uncles); hash != header.UncleHash {
		return fmt.Errorf("uncle root hash mismatch: have %x, want %x", hash, header.UncleHash)
	}
	hash, err := trie.DeriveSha(bl.Transactions)
	if err != nil {
		return err
	}
	if hash != header.TransactionsRoot {
		return fmt.Errorf(
			"transaction root hash mismatch: have %x, want %x",
			hash,
			header.TransactionsRoot,
		)
	}
	return nil
}

// ValidateState validates the various changes that happen after a state
// transition, such as amount of used gas, the receipt roots and the state root
// itself. ValidateState returns a database batch if the validation was a success
// otherwise nil and an error is returned.
func (v *BlockValidator) ValidateState(
	block, parent *block.Block,
	statedb *state.StateDB,
	receipts []*receipt.Receipt,
	usedGas uint64,
) error {
	header := block.Header
	logger.DebugP("Header in validate state", header)
	if header.GasUsed != usedGas {
		return fmt.Errorf("invalid gas used (remote: %d local: %d)", header.GasUsed, usedGas)
	}
	// Validate the received block's bloom with the one derived from the generated receipts.
	// For valid blocks this should always validate to true.
	rbloom := bloom.CreateBloom(receipts)
	if rbloom != bloom.BytesToBloom(header.Bloom.Bytes()) {
		return fmt.Errorf("invalid bloom (remote: %x  local: %x)", header.Bloom, rbloom)
	}
	// Tre receipt Trie's root (R = (Tr [[H1, R1], ... [Hn, R1]]))
	receiptSha, err := trie.DeriveSha(receipts)
	if err != nil {
		return err
	}
	if receiptSha != header.ReceiptRoot {
		return fmt.Errorf(
			"invalid receipt root hash (remote: %x local: %x)",
			header.ReceiptRoot,
			receiptSha,
		)
	}
	// Validate the state root against the received state root and throw
	// an error if they don't match.
	if root := statedb.IntermediateRoot(true); header.StateRoot != root {
		return fmt.Errorf("invalid merkle root (remote: %x local: %x)", header.StateRoot, root)
	}
	return nil
}

// CalcGasLimit computes the gas limit of the next block after parent.
// This is miner strategy, not consensus protocol.
func CalcGasLimit(parent *block.Block) uint64 {
	// contrib = (parentGasUsed * 3 / 2) / 1024
	contrib := (parent.Header.GasUsed + parent.Header.GasUsed/2) / params.GasLimitBoundDivisor

	// decay = parentGasLimit / 1024 -1
	decay := parent.GasLimit()/params.GasLimitBoundDivisor - 1

	/*
		strategy: gasLimit of block-to-mine is set based on parent's
		gasUsed value.  if parentGasUsed > parentGasLimit * (2/3) then we
		increase it, otherwise lower it (or leave it unchanged if it's right
		at that usage) the amount increased/decreased depends on how far away
		from parentGasLimit * (2/3) parentGasUsed is.
	*/
	limit := parent.GasLimit() - decay + contrib
	if limit < params.MinGasLimit {
		limit = params.MinGasLimit
	}
	// however, if we're now below the target (TargetGasLimit) we increase the
	// limit as much as we can (parentGasLimit / 1024 -1)
	if limit < params.TargetGasLimit {
		limit = parent.GasLimit() + decay
		if limit > params.TargetGasLimit {
			limit = params.TargetGasLimit
		}
	}
	return limit
}
