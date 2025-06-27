package state_processor

import (
	"math"
	"math/big"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/common"
	"FichainCore/errors"
	"FichainCore/evm"
	"FichainCore/gas_pool"
	"FichainCore/params"
	"FichainCore/transaction"
)

/*
The State Transitioning Model

A state transition is a change made when a transaction is applied to the current world state
The state transitioning model does all all the necessary work to work out a valid new state root.

1) Nonce handling
2) Pre pay gas
3) Create a new state object if the recipient is \0*32
4) Value transfer
== If contract creation ==

	4a) Attempt to run transaction data
	4b) If valid, use result as code for the new state object

== end ==
5) Run Script section
6) Derive new state root
*/
type StateTransition struct {
	gp         *gas_pool.GasPool
	tx         *transaction.Transaction
	gas        uint64
	gasPrice   *big.Int
	initialGas uint64
	value      *big.Int
	data       []byte
	state      evm.StateDB
	evm        *evm.EVM
}

// IntrinsicGas computes the 'intrinsic gas' for a message with the given data.
func IntrinsicGas(data []byte, contractCreation bool) (uint64, error) {
	// Set the starting gas for the raw transaction
	var gas uint64
	if contractCreation {
		gas = params.TxGasContractCreation
	} else {
		gas = params.TxGas
	}
	logger.DebugP("IntrinsicGas 1", gas)
	// Bump the required gas by the amount of transactional data
	if len(data) > 0 {
		// Zero and non-zero bytes are priced differently
		var nz uint64
		for _, byt := range data {
			if byt != 0 {
				nz++
			}
		}
		// Make sure we don't exceed uint64 for all data combinations
		if (math.MaxUint64-gas)/params.TxDataNonZeroGas < nz {
			return 0, errors.ErrOutOfGas
		}
		gas += nz * params.TxDataNonZeroGas

		z := uint64(len(data)) - nz
		if (math.MaxUint64-gas)/params.TxDataZeroGas < z {
			return 0, errors.ErrOutOfGas
		}
		gas += z * params.TxDataZeroGas
	}

	logger.DebugP("IntrinsicGas 2", gas)
	return gas, nil
}

// NewStateTransition initialises and returns a new state transition object.
func NewStateTransition(
	evm *evm.EVM,
	tx *transaction.Transaction,
	gp *gas_pool.GasPool,
) *StateTransition {
	return &StateTransition{
		gp:       gp,
		evm:      evm,
		tx:       tx,
		gasPrice: tx.GasPrice(),
		value:    tx.Amount(),
		data:     tx.Data(),
		state:    evm.StateDB,
	}
}

// ApplyMessage computes the new state by applying the given message
// against the old state within the environment.
//
// ApplyMessage returns the bytes returned by any EVM execution (if it took place),
// the gas used (which includes gas refunds) and an error if it failed. An error always
// indicates a core error meaning that the message would always fail for that particular
// state and would never be accepted within a block.
func ApplyMessage(
	evm *evm.EVM,
	tx *transaction.Transaction,
	gp *gas_pool.GasPool,
) ([]byte, uint64, bool, error) {
	return NewStateTransition(evm, tx, gp).TransitionDb()
}

func (st *StateTransition) from() evm.AccountRef {
	f, _ := st.tx.From(params.TempChainId)
	if !st.state.Exist(f) {
		st.state.CreateAccount(f)
	}
	return evm.AccountRef(f)
}

func (st *StateTransition) to() evm.AccountRef {
	if st.tx == nil {
		return evm.AccountRef{}
	}
	to := st.tx.To()
	if (to == common.Address{}) {
		return evm.AccountRef{} // contract creation
	}

	reference := evm.AccountRef(to)
	if !st.state.Exist(to) {
		st.state.CreateAccount(to)
	}
	return reference
}

func (st *StateTransition) useGas(amount uint64) error {
	if st.gas < amount {
		return errors.ErrOutOfGas
	}
	st.gas -= amount

	return nil
}

func (st *StateTransition) buyGas() error {
	var (
		state  = st.state
		sender = st.from()
	)
	mgval := new(big.Int).Mul(new(big.Int).SetUint64(st.tx.Gas()), st.gasPrice)
	logger.DebugP("Balance", state.GetBalance(sender.Address()))
	logger.DebugP("Address", sender.Address().String())
	logger.DebugP("MG value", mgval)
	if state.GetBalance(sender.Address()).Cmp(mgval) < 0 {
		return errors.ErrInsufficientBalanceForGas
	}
	if err := st.gp.SubGas(st.tx.Gas()); err != nil {
		return err
	}
	st.gas += st.tx.Gas()

	st.initialGas = st.tx.Gas()
	state.SubBalance(sender.Address(), mgval)
	return nil
}

func (st *StateTransition) preCheck() error {
	sender := st.from()

	nonce := st.state.GetNonce(sender.Address())
	if nonce < st.tx.Nonce() {
		return errors.ErrNonceTooHigh
	} else if nonce > st.tx.Nonce() {
		return errors.ErrNonceTooLow
	}
	return st.buyGas()
}

// TransitionDb will transition the state by applying the current message and
// returning the result including the the used gas. It returns an error if it
// failed. An error indicates a consensus issue.
func (st *StateTransition) TransitionDb() (ret []byte, usedGas uint64, failed bool, err error) {
	if err = st.preCheck(); err != nil {
		return
	}
	sender := st.from() // err checked in preCheck

	contractCreation := st.tx.To() == common.Address{}
	logger.DebugP("contractCreation", contractCreation)

	// Pay intrinsic gas
	gas, err := IntrinsicGas(st.data, contractCreation)
	if err = st.useGas(gas); err != nil {
		return nil, 0, false, err
	}

	var (
		vm = st.evm
		// vm errors do not effect consensus and are therefor
		// not assigned to err, except for insufficient balance
		// error.
		vmerr error
	)
	if contractCreation {
		ret, _, st.gas, vmerr = vm.Create(sender, st.data, st.gas, st.value)
		logger.DebugP("St.gas", st.gas)
		logger.Warn("XX1")
	} else {
		logger.Warn("XX2")
		// Increment the nonce for the next transaction
		st.state.SetNonce(sender.Address(), st.state.GetNonce(sender.Address())+1)
		ret, st.gas, vmerr = vm.Call(sender, st.to().Address(), st.data, st.gas, st.value)
	}
	if vmerr != nil {
		logger.Error("VM returned with error", "err", vmerr, string(ret))
		// The only possible consensus-error would be if there wasn't
		// sufficient balance to make the transfer happen. The first
		// balance transfer may never fail.
		if vmerr == errors.ErrInsufficientBalance {
			return nil, 0, false, vmerr
		}
	}
	st.refundGas()
	st.state.AddBalance(
		st.evm.Coinbase,
		new(big.Int).Mul(new(big.Int).SetUint64(st.gasUsed()), st.gasPrice),
	)

	return ret, st.gasUsed(), vmerr != nil, err
}

func (st *StateTransition) refundGas() {
	// Apply refund counter, capped to half of the used gas.
	refund := st.gasUsed() / 2
	logger.DebugP(
		"XX enter refund gas",
		"refund",
		refund,
		"gasused",
		st.gasUsed(),
		"get refund ",
		st.state.GetRefund(),
	)
	if refund > st.state.GetRefund() {
		refund = st.state.GetRefund()
	}
	st.gas += refund

	// Return ETH for remaining gas, exchanged at the original rate.
	sender := st.from()

	remaining := new(big.Int).Mul(new(big.Int).SetUint64(st.gas), st.gasPrice)
	st.state.AddBalance(sender.Address(), remaining)

	// Also return remaining gas to the block gas counter so it is
	// available for the next transaction.
	st.gp.AddGas(st.gas)
}

// gasUsed returns the amount of gas used up by the state transition.
func (st *StateTransition) gasUsed() uint64 {
	logger.DebugP("StateTransition gasUsed", st.initialGas, st.gas)
	return st.initialGas - st.gas
}
