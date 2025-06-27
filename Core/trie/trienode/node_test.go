package trienode

import (
	"crypto/rand"
	"testing"

	"github.com/ethereum/go-ethereum/crypto"

	"FichainCore/common"
)

func BenchmarkMerge(b *testing.B) {
	b.Run("1K", func(b *testing.B) {
		benchmarkMerge(b, 1000)
	})
	b.Run("10K", func(b *testing.B) {
		benchmarkMerge(b, 10_000)
	})
}

func benchmarkMerge(b *testing.B, count int) {
	x := NewNodeSet(common.Hash{})
	y := NewNodeSet(common.Hash{})
	addNode := func(s *NodeSet) {
		path := make([]byte, 4)
		rand.Read(path)
		blob := make([]byte, 32)
		rand.Read(blob)
		hash := crypto.Keccak256Hash(blob)
		s.AddNode(path, New(hash, blob))
	}
	for i := 0; i < count; i++ {
		// Random path of 4 nibbles
		addNode(x)
		addNode(y)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		// Store set x into a backup
		z := NewNodeSet(common.Hash{})
		z.Merge(common.Hash{}, x.Nodes)
		// Merge y into x
		x.Merge(common.Hash{}, y.Nodes)
		x = z
	}
}
