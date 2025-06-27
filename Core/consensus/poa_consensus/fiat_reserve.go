package poa_consensus

import (
	"errors"
	"math/big"
	"sync"

	"FichainCore/common"
	"FichainCore/database"
)

type FiatReserve struct {
	balances map[common.Address]*big.Int
	sync.RWMutex
}

func NewFiatReserve() *FiatReserve {
	return &FiatReserve{
		balances: make(map[common.Address]*big.Int),
	}
}

func (r *FiatReserve) Deposit(address common.Address, amount *big.Int) error {
	if amount.Sign() <= 0 {
		return errors.New("deposit amount must be positive")
	}

	r.Lock()
	defer r.Unlock()

	if _, ok := r.balances[address]; !ok {
		r.balances[address] = new(big.Int)
	}
	r.balances[address].Add(r.balances[address], amount)
	return nil
}

func (r *FiatReserve) Withdraw(address common.Address, amount *big.Int) error {
	if amount.Sign() <= 0 {
		return errors.New("withdrawal amount must be positive")
	}

	r.Lock()
	defer r.Unlock()

	balance, ok := r.balances[address]
	if !ok || balance.Cmp(amount) < 0 {
		return errors.New("insufficient balance")
	}

	balance.Sub(balance, amount)
	return nil
}

func (r *FiatReserve) GetBalance(address common.Address) (*big.Int, error) {
	r.RLock()
	defer r.RUnlock()

	balance, ok := r.balances[address]
	if !ok {
		return big.NewInt(0), nil
	}

	// return a copy to avoid external mutation
	return new(big.Int).Set(balance), nil
}

func (r *FiatReserve) SyncFromBank() error {
	// Simulated external sync (stub)
	// In real implementation, you would query an API or read from file/db
	return nil
}

func (r *FiatReserve) LoadFromStorage(db database.Database) error {
	err := db.IterateKeys(func(key, value []byte) error {
		address := common.BytesToAddress(key) // directly use key as address string
		amount := new(big.Int).SetBytes(value)
		r.balances[address] = amount
		return nil
	})
	return err
}

func (r *FiatReserve) CommitToStorage(db database.Database) error {
	r.RLock()
	defer r.RUnlock()

	batch := db.NewBatch()
	for address, amount := range r.balances {
		key := address.Bytes()  // key is just address string
		value := amount.Bytes() // store big.Int as raw bytes
		batch.Put(key, value)
	}
	return batch.Write()
}
