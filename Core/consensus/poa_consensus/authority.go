package poa_consensus

import (
	"encoding/json"
	"fmt"
	"math/big"
	"sync"

	logger "github.com/HendrickPhan/golang-simple-logger"

	"FichainCore/common"
	"FichainCore/database"
)

type NodeRole string

const (
	RoleValidator    NodeRole = "validator"
	RoleObserverNode NodeRole = "observer"
)

type Authority struct {
	validators map[common.Address]*big.Int         // address -> weight
	observers  map[common.Address][]common.Address // address -> accessible node list

	sync.RWMutex
}

func NewAuthority() *Authority {
	return &Authority{
		validators: make(map[common.Address]*big.Int),
		observers:  make(map[common.Address][]common.Address),
	}
}

// --------------------------- Load Functions ---------------------------

func (a *Authority) LoadFromDB(
	validatorDB, observerDB database.Database,
) error {
	a.Lock()
	defer a.Unlock()

	// Load validators

	err := validatorDB.IterateKeys(func(key, value []byte) error {
		// fmt.Printf("Key: %s, Value: %s\n", key, value)
		address := common.BytesToAddress(key)
		weight := new(big.Int).SetBytes(value)
		a.validators[address] = weight
		return nil
	})
	if err != nil {
		return fmt.Errorf("error loading validators: %w", err)
	}

	// Load observers
	err = observerDB.IterateKeys(func(key, value []byte) error {
		address := common.BytesToAddress(key)
		var accessList []common.Address
		if err := json.Unmarshal(value, &accessList); err != nil {
			logger.Warn(
				"unable to decode observer access list",
				"key",
				string(key),
				"err",
				err,
			)
		}
		a.observers[address] = accessList
		return nil
	})
	if err != nil {
		return fmt.Errorf("error loading validators: %w", err)
	}

	return nil
}

// --------------------------- Commit Functions ---------------------------

func (a *Authority) CommitToStorage(
	validatorDB, observerDB database.Database,
) error {
	a.RLock()
	defer a.RUnlock()

	// Commit validators
	valBatch := validatorDB.NewBatch()
	for address, weight := range a.validators {
		valBatch.Put(address.Bytes(), weight.Bytes())
	}
	if err := valBatch.Write(); err != nil {
		return fmt.Errorf("failed to commit validators: %w", err)
	}

	// Commit observers
	obsBatch := observerDB.NewBatch()
	for address, accessList := range a.observers {
		bytes, err := json.Marshal(accessList)
		if err != nil {
			return fmt.Errorf("failed to marshal observer list: %w", err)
		}
		obsBatch.Put(address.Bytes(), bytes)
	}
	if err := obsBatch.Write(); err != nil {
		return fmt.Errorf("failed to commit observers: %w", err)
	}

	return nil
}

// --------------------------- Validator APIs ---------------------------

func (a *Authority) AddValidator(address common.Address, weight *big.Int) error {
	a.Lock()
	defer a.Unlock()

	if _, exists := a.validators[address]; exists {
		return fmt.Errorf("validator already exists: %s", address.Hex())
	}
	a.validators[address] = weight
	return nil
}

func (a *Authority) RemoveValidator(address common.Address) error {
	a.Lock()
	defer a.Unlock()

	if _, exists := a.validators[address]; !exists {
		return fmt.Errorf("validator not found: %s", address.Hex())
	}
	delete(a.validators, address)
	return nil
}

func (a *Authority) GetValidatorWeight(address common.Address) (*big.Int, bool) {
	a.RLock()
	defer a.RUnlock()

	weight, ok := a.validators[address]
	return weight, ok
}

func (a *Authority) ListValidators() map[common.Address]*big.Int {
	a.RLock()
	defer a.RUnlock()

	// return a copy
	result := make(map[common.Address]*big.Int)
	for addr, w := range a.validators {
		result[addr] = new(big.Int).Set(w)
	}
	return result
}

// --------------------------- Observer APIs ---------------------------

func (a *Authority) AddObserver(address common.Address, accessList []common.Address) error {
	a.Lock()
	defer a.Unlock()

	if _, exists := a.observers[address]; exists {
		return fmt.Errorf("observer already exists: %s", address.Hex())
	}
	a.observers[address] = accessList
	return nil
}

func (a *Authority) RemoveObserver(address common.Address) error {
	a.Lock()
	defer a.Unlock()

	if _, exists := a.observers[address]; !exists {
		return fmt.Errorf("observer not found: %s", address.Hex())
	}
	delete(a.observers, address)
	return nil
}

func (a *Authority) GetObserverAccessList(address common.Address) ([]common.Address, bool) {
	a.RLock()
	defer a.RUnlock()

	list, ok := a.observers[address]
	return list, ok
}

func (a *Authority) ListObservers() map[common.Address][]common.Address {
	a.RLock()
	defer a.RUnlock()

	// return a copy
	result := make(map[common.Address][]common.Address)
	for addr, list := range a.observers {
		copied := make([]common.Address, len(list))
		copy(copied, list)
		result[addr] = copied
	}
	return result
}
