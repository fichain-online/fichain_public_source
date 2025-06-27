package poa_consensus

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"FichainCore/common"
	"FichainCore/database"
)

func TestAuthority(t *testing.T) {
	// Initialize in-memory databases for validator and observer storage
	validatorDB, _ := database.NewMemDatabase()
	observerDB, _ := database.NewMemDatabase()

	// Create a new Authority instance
	authority := NewAuthority()

	// Add validators with a weight of 1
	err := authority.AddValidator(common.Address{0x01}, big.NewInt(1))
	assert.NoError(t, err, "Adding validator should not return an error")

	err = authority.AddValidator(common.Address{0x02}, big.NewInt(1))
	assert.NoError(t, err, "Adding validator should not return an error")

	// Add observers with access to addresses 0x01 and 0x02
	err = authority.AddObserver(common.Address{0x01}, []common.Address{{0x01}, {0x02}})
	assert.NoError(t, err, "Adding observer should not return an error")

	err = authority.AddObserver(common.Address{0x02}, []common.Address{{0x01}, {0x02}})
	assert.NoError(t, err, "Adding observer should not return an error")

	// Commit the state of authority to the respective storages
	err = authority.CommitToStorage(validatorDB, observerDB)
	assert.NoError(t, err, "Committing to storage should not return an error")

	// Create a new Authority instance to load data from storage
	authority2 := NewAuthority()
	err = authority2.LoadFromDB(validatorDB, observerDB)
	assert.NoError(t, err, "Loading from storage should not return an error")

	// Verify that the list of validators is correct
	validators := authority.ListValidators()
	assert.Len(t, validators, 2, "There should be 2 validators in authority1")
	assert.Contains(t, validators, common.Address{0x01}, "Validator 0x01 should be in the list")
	assert.Contains(t, validators, common.Address{0x02}, "Validator 0x02 should be in the list")

	// Verify that the list of observers is correct
	observer := authority.ListObservers()
	assert.Len(t, observer, 2, "There should be 2 observers in authority1")
	assert.Contains(t, observer, common.Address{0x01}, "Observer 0x01 should be in the list")
	assert.Contains(t, observer, common.Address{0x02}, "Observer 0x02 should be in the list")

	// Verify that the list of validators is correct after loading from storage
	validators2 := authority2.ListValidators()
	assert.Len(t, validators2, 2, "There should be 2 validators in authority2")
	assert.Contains(t, validators2, common.Address{0x01}, "Validator 0x01 should be in the list")
	assert.Contains(t, validators2, common.Address{0x02}, "Validator 0x02 should be in the list")

	// Verify that the list of observers is correct after loading from storage
	observer2 := authority2.ListObservers()
	assert.Len(t, observer2, 2, "There should be 2 observers in authority2")
	assert.Contains(t, observer2, common.Address{0x01}, "Observer 0x01 should be in the list")
	assert.Contains(t, observer2, common.Address{0x02}, "Observer 0x02 should be in the list")
}
