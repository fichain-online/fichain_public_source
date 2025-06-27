package lookup_table

import (
	"sync"

	"FichainCore/common"
	"FichainCore/p2p"
)

// LookupTable maps addresses to TcpPeer connections (by pointer)
type LookupTable struct {
	mu    sync.RWMutex
	table map[common.Address]p2p.Peer
}

// NewLookupTable creates and returns a new LookupTable
func NewLookupTable() *LookupTable {
	return &LookupTable{
		table: make(map[common.Address]p2p.Peer),
	}
}

// Add inserts or updates a peer for a given address
func (lt *LookupTable) Add(address common.Address, p p2p.Peer) {
	lt.mu.Lock()
	defer lt.mu.Unlock()
	lt.table[address] = p
}

// Remove deletes a peer by address
func (lt *LookupTable) Remove(address common.Address) {
	lt.mu.Lock()
	defer lt.mu.Unlock()
	delete(lt.table, address)
}

// Get returns the TcpPeer pointer for a given address and a boolean indicating if it was found
func (lt *LookupTable) Get(address common.Address) (p2p.Peer, bool) {
	lt.mu.RLock()
	defer lt.mu.RUnlock()
	p, ok := lt.table[address]
	return p, ok
}

// Has checks if a given address exists in the table
func (lt *LookupTable) Has(address common.Address) bool {
	lt.mu.RLock()
	defer lt.mu.RUnlock()
	_, ok := lt.table[address]
	return ok
}

// All returns a copy of the entire address-to-peer map
func (lt *LookupTable) All() map[common.Address]p2p.Peer {
	lt.mu.RLock()
	defer lt.mu.RUnlock()
	copy := make(map[common.Address]p2p.Peer, len(lt.table))
	for k, v := range lt.table {
		copy[k] = v
	}
	return copy
}
