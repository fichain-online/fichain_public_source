package gas_pool

import (
	"fmt"
	"math"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/errors"
)

// GasPool tracks the amount of gas available during execution of the transactions
// in a block. The zero value is a pool with zero gas available.
type GasPool uint64

// AddGas makes gas available for execution.
func (gp *GasPool) AddGas(amount uint64) *GasPool {
	if uint64(*gp) > math.MaxUint64-amount {
		panic("gas pool pushed above uint64")
	}
	*(*uint64)(gp) += amount
	return gp
}

// SubGas deducts the given amount from the pool if enough gas is
// available and returns an error otherwise.
func (gp *GasPool) SubGas(amount uint64) error {
	logger.Warn("sub gas", amount)
	if uint64(*gp) < amount {
		return errors.ErrGasLimitReached
	}
	*(*uint64)(gp) -= amount
	return nil
}

// Gas returns the amount of gas remaining in the pool.
func (gp *GasPool) Gas() uint64 {
	return uint64(*gp)
}

func (gp *GasPool) String() string {
	return fmt.Sprintf("%d", *gp)
}
