package lookup_table

import (
	"testing"

	"FichainCore/common"
	"FichainCore/p2p/peer"
)

// Create a wrapper to satisfy *peer.TcpPeer type, or adjust as needed if it's an interface or struct in your codebase.
func newMockPeer() *peer.TcpPeer {
	// Cast mock to TcpPeer if TcpPeer is an interface
	mock := &peer.TcpPeer{}
	return (*peer.TcpPeer)(mock)
}

func TestLookupTable(t *testing.T) {
	lt := NewLookupTable()

	addr1 := common.Address{0x01}
	addr2 := common.Address{0x02}

	peer1 := newMockPeer()
	peer2 := newMockPeer()

	// Test Add and Get
	lt.Add(addr1, peer1)
	got, ok := lt.Get(addr1)
	if !ok {
		t.Fatalf("expected peer1 to be found")
	}
	if got != peer1 {
		t.Errorf("expected peer1, got %v", got.ID())
	}

	// Test Has
	if !lt.Has(addr1) {
		t.Errorf("expected Has(addr1) to be true")
	}
	if lt.Has(addr2) {
		t.Errorf("expected Has(addr2) to be false")
	}

	// Test Remove
	lt.Remove(addr1)
	_, ok = lt.Get(addr1)
	if ok {
		t.Errorf("expected peer1 to be removed")
	}

	// Test All
	lt.Add(addr1, peer1)
	lt.Add(addr2, peer2)
	all := lt.All()
	if len(all) != 2 {
		t.Errorf("expected 2 peers in All(), got %d", len(all))
	}
	if all[addr2] != peer2 {
		t.Errorf("expected peer2 in All(), got %v", all[addr2].ID())
	}
}
