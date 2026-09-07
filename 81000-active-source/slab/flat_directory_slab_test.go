package slab_test

import (
	"bytes"
	"fmt"
	"testing"

	"sov.fleet/plinth-filesystem/81000-active-source/slab"
)

func TestFlatDirectorySlab_InsertAndFind(t *testing.T) {
	s := slab.NewFlatDirectorySlab(1024)

	var dummySeal [32]byte
	copy(dummySeal[:], "WHITE3_DUMMY_M31_SEAL_TEST_0001")

	// 1. Insert multiple directory entries
	slot1, err := s.Insert(0, "system", 0, 4096, slab.FlagDirectory, dummySeal)
	if err != nil {
		t.Fatalf("failed to insert root child: %v", err)
	}

	slot2, err := s.Insert(slot1, "kernel.sls", 4096, 65536, 0, dummySeal)
	if err != nil {
		t.Fatalf("failed to insert nested file: %v", err)
	}

	slot3, err := s.Insert(slot1, "config.wag", 69632, 1024, slab.FlagReadOnly, dummySeal)
	if err != nil {
		t.Fatalf("failed to insert config: %v", err)
	}

	// 2. Lookup existing files
	foundSlot, ok := s.FindEntry(0, "system")
	if !ok || foundSlot != slot1 {
		t.Errorf("expected to find 'system' at slot %d, got ok=%v, slot=%d", slot1, ok, foundSlot)
	}

	foundNested, ok := s.FindEntry(slot1, "kernel.sls")
	if !ok || foundNested != slot2 {
		t.Errorf("expected to find 'kernel.sls' at slot %d, got ok=%v, slot=%d", slot2, ok, foundNested)
	}

	foundConfig, ok := s.FindEntry(slot1, "config.wag")
	if !ok || foundConfig != slot3 {
		t.Errorf("expected to find 'config.wag' at slot %d, got ok=%v, slot=%d", slot3, ok, foundConfig)
	}

	// 3. Lookup non-existent file
	_, ok = s.FindEntry(slot1, "nonexistent.dat")
	if ok {
		t.Errorf("expected nonexistent.dat to return false")
	}

	// 4. Verify extent metadata
	off, lenBytes, seal, err := s.GetExtentInfo(slot2)
	if err != nil {
		t.Fatalf("GetExtentInfo failed: %v", err)
	}
	if off != 4096 || lenBytes != 65536 || !bytes.Equal(seal[:], dummySeal[:]) {
		t.Errorf("extent mismatch: off=%d, len=%d, seal=%x", off, lenBytes, seal)
	}

	// 5. Soft delete
	if err := s.SoftDelete(slot3); err != nil {
		t.Fatalf("SoftDelete failed: %v", err)
	}
	_, ok = s.FindEntry(slot1, "config.wag")
	if ok {
		t.Errorf("expected soft-deleted config.wag to not be found")
	}

	t.Logf("✅ [Flat Directory Slab Test Passed] %s", s.String())
}

func TestFlatDirectorySlab_ListChildren(t *testing.T) {
	s := slab.NewFlatDirectorySlab(256)
	var dummySeal [32]byte

	parentID := uint32(1)
	fileNames := []string{"alpha.dat", "beta.dat", "gamma.dat", "delta.dat"}

	for i, name := range fileNames {
		_, err := s.Insert(parentID, name, uint64(i*4096), 4096, 0, dummySeal)
		if err != nil {
			t.Fatalf("failed to insert %s: %v", name, err)
		}
	}

	// Another parent directory
	_, err := s.Insert(2, "other.dat", 16384, 1024, 0, dummySeal)
	if err != nil {
		t.Fatalf("failed to insert: %v", err)
	}

	children := s.ListChildren(parentID, nil)
	if len(children) != 4 {
		t.Fatalf("expected 4 children for parent %d, got %d", parentID, len(children))
	}
	t.Logf("✅ [ListChildren Verified] Found %d children for parentID=%d", len(children), parentID)
}

func TestFlatExtentStream_AppendAndReadWindow(t *testing.T) {
	stream := slab.NewFlatExtentStream(1024 * 1024) // 1MB

	payload := []byte("SOVEREIGN_M31_EXTENT_DATA_BLOCK_0123456789")
	off, length, seal, err := stream.Append(payload)
	if err != nil {
		t.Fatalf("Append failed: %v", err)
	}

	if off%64 != 0 {
		t.Errorf("expected 64-byte cacheline alignment, got offset=%d", off)
	}
	if length != uint32(len(payload)) {
		t.Errorf("expected length=%d, got %d", len(payload), length)
	}

	// Read window (zero-copy)
	window, err := stream.ReadWindow(off, length)
	if err != nil {
		t.Fatalf("ReadWindow failed: %v", err)
	}
	if !bytes.Equal(window, payload) {
		t.Errorf("data mismatch in window: got %s, want %s", string(window), string(payload))
	}

	t.Logf("✅ [Flat Extent Stream Verified] %s, seal=%x", stream.String(), seal[:8])
}

// BenchmarkFlatDirectorySlab_LinearScan measures lookup latency over a flat 1,024-entry SoA slab.
func BenchmarkFlatDirectorySlab_LinearScan(b *testing.B) {
	s := slab.NewFlatDirectorySlab(1024)
	var dummySeal [32]byte

	for i := 0; i < 1000; i++ {
		_, _ = s.Insert(1, fmt.Sprintf("dataset_shard_%04d.sls", i), uint64(i*4096), 4096, 0, dummySeal)
	}

	target := "dataset_shard_0999.sls"
	targetHash := slab.FastHash64(target)

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		slot, ok := s.FindByHash(1, targetHash)
		if !ok || slot != 999 {
			b.Fatalf("lookup failed")
		}
	}
}
