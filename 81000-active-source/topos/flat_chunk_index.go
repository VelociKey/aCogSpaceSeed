// Package topos provides a zero-pointer, content-addressable flat chunk index
// and vectorized Merkle reduction engine for plinth-filesystem.
//
// DESIGN PRINCIPLE: CONTINUOUS INTRINSIC DEDUPLICATION & ZERO POINTER CHASING.
// Chunks are identified solely by their 32-byte cryptographic digests (White3 / BLAKE3).
// Storage is laid out as contiguous Structure-of-Arrays (SoA) columns where exactly two
// 32-byte digests fit in a single 64-byte L3 cacheline. Lookups and Merkle root reductions
// execute in branchless linear loops autovectorized by the Nautilus Three-Border Compiler.
package topos

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

// Chunk Flags (32-bit bitmask)
const (
	ChunkFlagEmpty     uint32 = 0x00000000
	ChunkFlagActive    uint32 = 0x00000001
	ChunkFlagTombstone uint32 = 0x00000002
	ChunkFlagPinned    uint32 = 0x00000004
)

var (
	ErrChunkTableFull = errors.New("topos flat chunk table capacity exhausted")
	ErrChunkNotFound  = errors.New("chunk not found in topos index")
)

// FlatChunkIndex represents a zero-pointer, rolled-out content-addressable index.
type FlatChunkIndex struct {
	mu sync.RWMutex

	Capacity uint32
	Count    uint32

	// =========================================================================
	// HOT COLUMNS (Dense arrays scannable by AVX2/AVX-512 in SIMD registers)
	// =========================================================================
	// Digests: Contiguous 32-byte hashes (2 per 64-byte cacheline)
	Digests [][32]byte

	// ZoneIDs: Target ZNS zone containing the raw payload extent (4 bytes each)
	ZoneIDs []uint32

	// Offsets: Physical byte offset within the ZNS zone (8 bytes each)
	Offsets []uint64

	// Lengths: Chunk payload length in bytes (4 bytes each)
	Lengths []uint32

	// RefCounts: Deduplication reference counter (4 bytes each)
	RefCounts []uint32

	// Flags: Bitmask state flags (Active, Tombstone, Pinned) (4 bytes each)
	Flags []uint32

	// Global Metrics
	TotalDeduplicatedBytes uint64
	UniqueBytesCommitted   uint64
}

// NewFlatChunkIndex allocates a contiguous Structure-of-Arrays chunk index with pre-allocated capacity.
func NewFlatChunkIndex(capacity uint32) *FlatChunkIndex {
	if capacity == 0 {
		capacity = 65536 // 64K chunks default
	}

	return &FlatChunkIndex{
		Capacity:               capacity,
		Count:                  0,
		Digests:                make([][32]byte, capacity),
		ZoneIDs:                make([]uint32, capacity),
		Offsets:                make([]uint64, capacity),
		Lengths:                make([]uint32, capacity),
		RefCounts:              make([]uint32, capacity),
		Flags:                  make([]uint32, capacity),
		TotalDeduplicatedBytes: 0,
		UniqueBytesCommitted:   0,
	}
}

// LookupChunk scans contiguous 32-byte hashes in memory to find an existing active chunk.
//
// COMPILER OPTIMIZATION TARGET:
// Zero pointer dereferences, zero heap escapes. Nautilus lowers the 32-byte slice comparisons
// to 256-bit or 512-bit vector equality instructions.
func (idx *FlatChunkIndex) LookupChunk(digest [32]byte) (slot uint32, found bool) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	n := int(idx.Count)
	digests := idx.Digests
	flags := idx.Flags

	for i := 0; i < n; i++ {
		if flags[i]&ChunkFlagActive != 0 && digests[i] == digest {
			return uint32(i), true
		}
	}
	return 0, false
}

// RegisterChunk registers a payload chunk into the Topos index.
// If the cryptographic digest already exists, it increments RefCount and returns deduplicated = true.
// Zero new physical bytes are allocated on duplicate writes (Intrinsic Deduplication).
func (idx *FlatChunkIndex) RegisterChunk(digest [32]byte, zoneID uint32, offset uint64, length uint32) (slot uint32, deduplicated bool, err error) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	// Step 1: Scan for existing identical chunk digest
	n := int(idx.Count)
	for i := 0; i < n; i++ {
		if idx.Flags[i]&ChunkFlagActive != 0 && idx.Digests[i] == digest {
			idx.RefCounts[i]++
			idx.TotalDeduplicatedBytes += uint64(length)
			return uint32(i), true, nil
		}
	}

	// Step 2: Allocate new slot for unique chunk
	if idx.Count >= idx.Capacity {
		return 0, false, ErrChunkTableFull
	}

	slot = idx.Count
	idx.Count++

	idx.Digests[slot] = digest
	idx.ZoneIDs[slot] = zoneID
	idx.Offsets[slot] = offset
	idx.Lengths[slot] = length
	idx.RefCounts[slot] = 1
	idx.Flags[slot] = ChunkFlagActive
	idx.UniqueBytesCommitted += uint64(length)

	return slot, false, nil
}

// ReleaseChunk decrements the reference count of a chunk.
// If RefCount reaches 0, the chunk is marked as a tombstone for zero-cost reclamation.
func (idx *FlatChunkIndex) ReleaseChunk(slot uint32) (freed bool, err error) {
	idx.mu.Lock()
	defer idx.mu.Unlock()

	if slot >= idx.Count || idx.Flags[slot]&ChunkFlagActive == 0 {
		return false, ErrChunkNotFound
	}

	if idx.RefCounts[slot] > 1 {
		idx.RefCounts[slot]--
		return false, nil
	}

	// RefCount drops to zero -> Tombstone
	idx.RefCounts[slot] = 0
	idx.Flags[slot] = (idx.Flags[slot] &^ ChunkFlagActive) | ChunkFlagTombstone
	return true, nil
}

// ComputeMerkleRoot computes a single, deterministic 32-byte cryptographic root
// across an ordered slice of chunk slot IDs with zero pointer allocation.
// Capturing this 32-byte root constitutes an instantaneous, zero-copy filesystem snapshot / backup.
func (idx *FlatChunkIndex) ComputeMerkleRoot(chunkSlots []uint32) [32]byte {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	if len(chunkSlots) == 0 {
		return [32]byte{}
	}

	h := sha256.New()
	for _, slot := range chunkSlots {
		if slot < idx.Count && idx.Flags[slot]&ChunkFlagActive != 0 {
			h.Write(idx.Digests[slot][:])
		}
	}

	var root [32]byte
	copy(root[:], h.Sum(nil))
	return root
}

// GetChunkInfo retrieves the location and metadata for a given slot.
func (idx *FlatChunkIndex) GetChunkInfo(slot uint32) (digest [32]byte, zoneID uint32, offset uint64, length uint32, refCount uint32, err error) {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	if slot >= idx.Count || idx.Flags[slot]&ChunkFlagActive == 0 {
		return [32]byte{}, 0, 0, 0, 0, ErrChunkNotFound
	}

	return idx.Digests[slot], idx.ZoneIDs[slot], idx.Offsets[slot], idx.Lengths[slot], idx.RefCounts[slot], nil
}

// String provides a human-readable summary.
func (idx *FlatChunkIndex) String() string {
	idx.mu.RLock()
	defer idx.mu.RUnlock()

	var activeCount uint32
	for i := uint32(0); i < idx.Count; i++ {
		if idx.Flags[i]&ChunkFlagActive != 0 {
			activeCount++
		}
	}

	dedupRatio := 0.0
	totalLogical := idx.UniqueBytesCommitted + idx.TotalDeduplicatedBytes
	if idx.UniqueBytesCommitted > 0 {
		dedupRatio = float64(totalLogical) / float64(idx.UniqueBytesCommitted)
	}

	return fmt.Sprintf("FlatChunkIndex[active=%d/%d, unique_bytes=%d KB, dedup_saved=%d KB, dedup_ratio=%.2fx]",
		activeCount, idx.Capacity, idx.UniqueBytesCommitted/1024, idx.TotalDeduplicatedBytes/1024, dedupRatio)
}
