package benchmark

import (
	"fmt"
	"testing"

	"sov.fleet/plinth-filesystem/81000-active-source/slab"
)

// PointerInode simulates a traditional POSIX / VFS pointer-based inode node.
type PointerInode struct {
	Parent   *PointerInode
	Name     string
	Offset   uint64
	Length   uint32
	Flags    uint32
	Seal     [32]byte
	Children map[string]*PointerInode
}

func Benchmark_FlatDirectorySlab_Lookup(b *testing.B) {
	const n = 10000
	dirSlab := slab.NewFlatDirectorySlab(n + 10)
	var seal [32]byte

	for i := 0; i < n; i++ {
		name := fmt.Sprintf("dataset_shard_%08d.bin", i)
		_, _ = dirSlab.Insert(1, name, uint64(i*4096), 4096, slab.FlagActive, seal)
	}

	targetName := fmt.Sprintf("dataset_shard_%08d.bin", n-1)
	targetHash := slab.FastHash64(targetName)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		slot, ok := dirSlab.FindByHash(1, targetHash)
		if !ok || int(slot) != n-1 {
			b.Fatalf("lookup failed: slot=%d, ok=%v", slot, ok)
		}
	}
}

func Benchmark_PointerMap_Lookup(b *testing.B) {
	const n = 10000
	root := &PointerInode{
		Name:     "/",
		Children: make(map[string]*PointerInode, n),
	}

	for i := 0; i < n; i++ {
		name := fmt.Sprintf("dataset_shard_%08d.bin", i)
		node := &PointerInode{
			Parent: root,
			Name:   name,
			Offset: uint64(i * 4096),
			Length: 4096,
			Flags:  1,
		}
		root.Children[name] = node
	}

	targetName := fmt.Sprintf("dataset_shard_%08d.bin", n-1)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		node, ok := root.Children[targetName]
		if !ok || node.Offset != uint64((n-1)*4096) {
			b.Fatalf("lookup failed")
		}
	}
}
