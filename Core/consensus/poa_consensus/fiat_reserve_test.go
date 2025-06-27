package poa_consensus

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"FichainCore/common"
	"FichainCore/database"
)

func TestFiatReserve(t *testing.T) {
	// Create an in-memory database for testing
	memDB, _ := database.NewMemDatabase()

	// Create a new FiatReserve
	fiatReserve := NewFiatReserve()

	// Define an amount
	amount := &big.Int{}
	amount.SetString("10000000000000000000000000000", 10)

	// Deposit amounts for two addresses
	err := fiatReserve.Deposit(common.Address{0x01}, amount)
	assert.NoError(t, err, "Deposit should not return an error")

	err = fiatReserve.Deposit(common.Address{0x02}, amount)
	assert.NoError(t, err, "Deposit should not return an error")

	// Commit the fiatReserve state to storage
	err = fiatReserve.CommitToStorage(memDB)
	assert.NoError(t, err, "CommitToStorage should not return an error")

	// Create a new FiatReserve and load the data from storage
	fiatReserve2 := NewFiatReserve()
	err = fiatReserve2.LoadFromStorage(memDB)
	assert.NoError(t, err, "LoadFromStorage should not return an error")

	// Check the balance for address 0x01
	balance, err := fiatReserve2.GetBalance(common.Address{0x01})
	assert.NoError(t, err, "GetBalance should not return an error")
	assert.Equal(
		t,
		amount,
		balance,
		"The balance for address 0x01 should be equal to the deposited amount",
	)

	// Check the balance for address 0x02
	balance2, err := fiatReserve2.GetBalance(common.Address{0x02})
	assert.NoError(t, err, "GetBalance should not return an error")
	assert.Equal(
		t,
		amount,
		balance2,
		"The balance for address 0x02 should be equal to the deposited amount",
	)
}
