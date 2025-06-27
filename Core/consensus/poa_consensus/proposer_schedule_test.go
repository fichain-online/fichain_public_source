package poa_consensus

import (
	"math/big"
	"testing"

	"github.com/stretchr/testify/assert"

	"FichainCore/common"
)

func TestProposerSchedule(t *testing.T) {
	// Sample authorities with weights
	authorities := map[common.Address]*big.Int{
		common.HexToAddress("0xabc123"): big.NewInt(50),
		common.HexToAddress("0xdef456"): big.NewInt(100),
		common.HexToAddress("0x789012"): big.NewInt(200),
	}

	// Salt for schedule update (it can be a fixed value or dynamically generated)
	salt := []byte("test_salt")

	// Create a new ProposerSchedule
	ps := &ProposerSchedule{}

	// Update the schedule with authorities and salt
	err := ps.UpdateSchedule(salt, 10, authorities)
	assert.NoError(t, err, "Error should be nil when updating schedule")

	// Test that the same input produces the same result (deterministic)
	t.Run("TestDeterministicProposerGeneration", func(t *testing.T) {
		ps2 := &ProposerSchedule{}
		err := ps2.UpdateSchedule(salt, 10, authorities)
		assert.NoError(t, err, "Error should be nil when updating schedule")

		// Fetch proposer for block height 0 from both schedules
		proposer0a, err := ps.GetProposer(10)
		assert.NoError(t, err, "Error should be nil when fetching proposer for block height 0")
		proposer0b, err := ps2.GetProposer(10)
		assert.NoError(
			t,
			err,
			"Error should be nil when fetching proposer for block height 0 from second schedule",
		)

		// Verify that both schedules produce the same proposer for block height 0
		assert.Equal(
			t,
			proposer0a,
			proposer0b,
			"Proposer for block height 0 should be the same across both schedules",
		)

		// Fetch proposer for block height 1 from both schedules
		proposer1a, err := ps.GetProposer(11)
		assert.NoError(t, err, "Error should be nil when fetching proposer for block height 1")
		proposer1b, err := ps2.GetProposer(11)
		assert.NoError(
			t,
			err,
			"Error should be nil when fetching proposer for block height 1 from second schedule",
		)

		// Verify that both schedules produce the same proposer for block height 1
		assert.Equal(
			t,
			proposer1a,
			proposer1b,
			"Proposer for block height 1 should be the same across both schedules",
		)
	})

	// Test for invalid block height (no proposer for block height 1000)
	t.Run("TestProposerForInvalidBlockHeight", func(t *testing.T) {
		// Fetch proposer for an invalid block height (e.g., 10011)
		proposer, err := ps.GetProposer(10011)
		assert.Error(t, err, "Error should occur when fetching proposer for block height 1000")
		assert.Equal(
			t,
			common.Address{},
			proposer,
			"Proposer for block height 1000 should be an empty address",
		)
	})
}
