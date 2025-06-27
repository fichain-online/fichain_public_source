package poa_consensus

import (
	"fmt"
	"math/big"
	"sort"

	logger "github.com/HendrickPhan/golang-simple-logger"
	"golang.org/x/crypto/sha3"

	"FichainCore/common"
	"FichainCore/params"
)

type ProposerSchedule struct {
	schedule map[uint64]common.Address
}

func NewProposerSchedule(
	salt []byte,
	fromBlock uint64,
	authoritiesWithWeight map[common.Address]*big.Int,
) *ProposerSchedule {
	s := &ProposerSchedule{
		schedule: map[uint64]common.Address{},
	}
	s.UpdateSchedule(salt, fromBlock, authoritiesWithWeight)
	logger.Debug("length of proposer schedule", len(s.schedule))

	return s
}

// GetProposer returns the proposer for a given block height.
func (ps *ProposerSchedule) GetProposer(blockHeight uint64) (common.Address, error) {
	proposer, exists := ps.schedule[blockHeight]
	if !exists {
		return common.Address{}, fmt.Errorf("no proposer found for block height %d", blockHeight)
	}
	return proposer, nil
}

func (ps *ProposerSchedule) UpdateSchedule(
	salt []byte,
	fromBlock uint64,
	authoritiesWithWeight map[common.Address]*big.Int,
) error {
	// Step 1: Sort authorities by weight (descending order)
	type authority struct {
		address common.Address
		weight  *big.Int
	}

	var sortedAuthorities []authority
	for address, weight := range authoritiesWithWeight {
		sortedAuthorities = append(sortedAuthorities, authority{
			address: address,
			weight:  weight,
		})
	}

	// Sort in descending order of weight
	sort.SliceStable(sortedAuthorities, func(i, j int) bool {
		return sortedAuthorities[i].weight.Cmp(sortedAuthorities[j].weight) > 0
	})

	// Step 2: Calculate the total weight
	totalWeight := big.NewInt(0)
	for _, auth := range sortedAuthorities {
		totalWeight.Add(totalWeight, auth.weight)
	}

	// Step 3: Generate a deterministic proposer schedule based on salt and weight distribution
	ps.schedule = make(map[uint64]common.Address)

	// Initialize variables for block assignments
	authoritiesAssigned := make(map[common.Address]uint64)

	// Use the salt to ensure deterministic results
	hash := sha3.NewLegacyKeccak256()
	hash.Write(salt)
	hash.Write(totalWeight.Bytes()) // Adding total weight for determinism

	// Step 4: Assign blocks to authorities based on their weight
	for blockHeight := uint64(fromBlock); blockHeight < fromBlock+params.TempEpochLength; blockHeight++ { // Example for the first 100 blocks
		// Hash block height with the base hash
		blockHash := sha3.NewLegacyKeccak256()
		blockHash.Write(hash.Sum(nil))
		blockHash.Write([]byte(fmt.Sprintf("%d", blockHeight)))

		// Select an authority based on weight
		for _, auth := range sortedAuthorities {
			weight := auth.weight
			// Calculate how many blocks this authority should control
			// Proportional assignment: authority weight / total weight * total blocks
			assignedBlocks := new(
				big.Int,
			).Mul(weight, big.NewInt(params.TempEpochLength))
			// Total of 100 blocks
			assignedBlocks.Div(assignedBlocks, totalWeight)

			// Check if this authority should be assigned to this block
			if authoritiesAssigned[auth.address] < assignedBlocks.Uint64() {
				ps.schedule[blockHeight] = auth.address
				authoritiesAssigned[auth.address]++
				break
			}
		}
	}

	return nil
}
