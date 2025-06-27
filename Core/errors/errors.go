package errors

import (
	"errors"
)

var (
	// evm error
	ErrOutOfGas                 = errors.New("out of gas")
	ErrCodeStoreOutOfGas        = errors.New("contract creation code storage out of gas")
	ErrDepth                    = errors.New("max call depth exceeded")
	ErrTraceLimitReached        = errors.New("the number of logs reached the specified limit")
	ErrInsufficientBalance      = errors.New("insufficient balance for transfer")
	ErrContractAddressCollision = errors.New("contract address collision")
	//
	ErrGasLimitReached = errors.New("gas limit reached")

	//
	ErrInsufficientBalanceForGas = errors.New("insufficient balance to pay for gas")

	// ErrKnownBlock is returned when a block to import is already known locally.
	ErrKnownBlock = errors.New("block already known")
	// ErrBlacklistedHash is returned if a block to import is on the blacklist.
	ErrBlacklistedHash = errors.New("blacklisted hash")
	// ErrNonceTooHigh is returned if the nonce of a transaction is higher than the
	// next one expected based on the local chain.
	ErrNonceTooHigh = errors.New("nonce too high")
	// ErrInvalidSender is returned if the transaction contains an invalid signature.
	ErrInvalidSender = errors.New("invalid sender")
	// ErrNonceTooLow is returned if the nonce of a transaction is lower than the
	// one present in the local chain.
	ErrNonceTooLow = errors.New("nonce too low")
	// ErrUnderpriced is returned if a transaction's gas price is below the minimum
	// configured for the transaction pool.
	ErrUnderpriced = errors.New("transaction underpriced")
	// ErrReplaceUnderpriced is returned if a transaction is attempted to be replaced
	// with a different one without the required price bump.
	ErrReplaceUnderpriced = errors.New("replacement transaction underpriced")
	// ErrInsufficientFunds is returned if the total cost of executing a transaction
	// is higher than the balance of the user's account.
	ErrInsufficientFunds = errors.New("insufficient funds for gas * price + value")
	// ErrIntrinsicGas is returned if the transaction is specified to use less gas
	// than required to start the invocation.
	ErrIntrinsicGas = errors.New("intrinsic gas too low")
	// ErrGasLimit is returned if a transaction's requested gas limit exceeds the
	// maximum allowance of the current block.
	ErrGasLimit = errors.New("exceeds block gas limit")
	// ErrNegativeValue is a sanity error to ensure noone is able to specify a
	// transaction with a negative value.
	ErrNegativeValue = errors.New("negative value")
	// ErrOversizedData is returned if the input data of a transaction is greater
	// than some meaningful limit a user might use. This is not a consensus error
	// making the transaction invalid, rather a DOS protection.
	ErrOversizedData = errors.New("oversized data")

	// evm
	ErrWriteProtection       = errors.New("evm: write protection")
	ErrReturnDataOutOfBounds = errors.New("evm: return data out of bounds")
	ErrExecutionReverted     = errors.New("evm: execution reverted")
	ErrMaxCodeSizeExceeded   = errors.New("evm: max code size exceeded")

	ErrAborted = errors.New("aborted")

	ErrNoGenesis = errors.New("Genesis not found in chain")

	// consensus

	ErrUnknownAncestor = errors.New("unknown ancestor")

	// ErrPrunedAncestor is returned when validating a block requires an ancestor
	// that is known, but the state of which is not available.
	ErrPrunedAncestor = errors.New("pruned ancestor")

	// ErrFutureBlock is returned when a block's timestamp is in the future according
	// to the current node.
	ErrFutureBlock = errors.New("block in the future")

	// ErrInvalidNumber is returned if a block's number doesn't equal it's parent's
	// plus one.
	ErrInvalidNumber = errors.New("invalid block number")
)
