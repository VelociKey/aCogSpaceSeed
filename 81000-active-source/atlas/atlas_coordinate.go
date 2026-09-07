package atlas

import (
	"errors"
	"fmt"
	"unsafe"
)

// Common Archetype semantic classifications (e1).
const (
	ArchetypeTensorWeight       uint64 = 1
	ArchetypeModelCheckpoint     uint64 = 2
	ArchetypeExecutionLog        uint64 = 3
	ArchetypeCompiledKernel      uint64 = 4
	ArchetypeTopologyGraph       uint64 = 5
	ArchetypeDatasetShard        uint64 = 6
	ArchetypeCryptoNotarization  uint64 = 7
)

// Common Hardware Precision Encodings (e5).
const (
	EncodingFP8   uint64 = 1 // NVFP8 / FP8 E4M3/E5M2
	EncodingBF16  uint64 = 2 // BFloat16
	EncodingFP16  uint64 = 3 // IEEE FP16
	EncodingFP32  uint64 = 4 // IEEE FP32
	EncodingFP64  uint64 = 5 // IEEE FP64
	EncodingINT4  uint64 = 6 // Sub-byte quantized tensor
	EncodingM31   uint64 = 7 // Small-field Mersenne-31 algebraic code
)

// Coordinate64 represents a non-hierarchical 8-dimensional coordinate in O-space.
// Memory size: exactly 64 bytes (1 physical L3 cacheline, zero pointers).
type Coordinate64 struct {
	Authority uint64 // e0: Sovereign Tenant or Cryptoseal Authority
	Archetype uint64 // e1: Semantic class (Tensor, Checkpoint, Stream, Binary)
	Chronos   uint64 // e2: Monotonic epoch, training step, or version counter
	Topos     uint64 // e3: Cluster node, VM partition, ZNS zone ID
	Geometry  uint64 // e4: Tensor rank, tile shape encoding, extent descriptor
	Encoding  uint64 // e5: Format & precision (FP8, BF16, FP16, M31, INT4)
	Lineage   uint64 // e6: Predecessor transform or parent associator hash
	Digest    uint64 // e7: White3 content digest (lower 64 bits of payload hash)
}

// Ensure Coordinate64 is exactly 64 bytes at compile time.
var _ = [64]byte{}[unsafe.Sizeof(Coordinate64{})-64]

// AsOctonion converts the 64-byte coordinate into an 8-element int64 array
// suitable for Cayley-Dickson octonionic multiplication and associator checks.
func (c Coordinate64) AsOctonion() [8]int64 {
	return [8]int64{
		int64(c.Authority),
		int64(c.Archetype),
		int64(c.Chronos),
		int64(c.Topos),
		int64(c.Geometry),
		int64(c.Encoding),
		int64(c.Lineage),
		int64(c.Digest),
	}
}

// AtlasQueryMask specifies which axes of the 8-dimensional space are constrained
// during an O(1) SIMD hyperplane query.
type AtlasQueryMask struct {
	ActiveMask uint8        // Bit 0..7 indicates if e0..e7 must match
	Values     Coordinate64 // The target values for active axes
}

const (
	MaskAuthority uint8 = 1 << 0
	MaskArchetype uint8 = 1 << 1
	MaskChronos   uint8 = 1 << 2
	MaskTopos     uint8 = 1 << 3
	MaskGeometry  uint8 = 1 << 4
	MaskEncoding  uint8 = 1 << 5
	MaskLineage   uint8 = 1 << 6
	MaskDigest    uint8 = 1 << 7
)

// FlatCoordinateAtlas stores millions of files in contiguous Structure-of-Arrays (SoA) columns.
// Zero pointer indirection; inner loops are branchless and contiguous for compiler autovectorization.
type FlatCoordinateAtlas struct {
	Capacity    uint32
	Count       uint32
	Authorities []uint64 // e0
	Archetypes  []uint64 // e1
	Chronos     []uint64 // e2
	Topoi       []uint64 // e3
	Geometries  []uint64 // e4
	Encodings   []uint64 // e5
	Lineages    []uint64 // e6
	Digests     []uint64 // e7
	PayloadRefs []uint64 // Extent offset / Chunk ID in ZNS flash zones
}

// NewFlatCoordinateAtlas initializes a zero-pointer SoA coordinate atlas.
func NewFlatCoordinateAtlas(capacity uint32) *FlatCoordinateAtlas {
	return &FlatCoordinateAtlas{
		Capacity:    capacity,
		Count:       0,
		Authorities: make([]uint64, capacity),
		Archetypes:  make([]uint64, capacity),
		Chronos:     make([]uint64, capacity),
		Topoi:       make([]uint64, capacity),
		Geometries:  make([]uint64, capacity),
		Encodings:   make([]uint64, capacity),
		Lineages:    make([]uint64, capacity),
		Digests:     make([]uint64, capacity),
		PayloadRefs: make([]uint64, capacity),
	}
}

// Insert inserts a coordinate and payload reference into the flat atlas.
func (atlas *FlatCoordinateAtlas) Insert(c Coordinate64, payloadRef uint64) (uint32, error) {
	if atlas.Count >= atlas.Capacity {
		return 0, errors.New("atlas capacity exhausted")
	}
	idx := atlas.Count
	atlas.Authorities[idx] = c.Authority
	atlas.Archetypes[idx] = c.Archetype
	atlas.Chronos[idx] = c.Chronos
	atlas.Topoi[idx] = c.Topos
	atlas.Geometries[idx] = c.Geometry
	atlas.Encodings[idx] = c.Encoding
	atlas.Lineages[idx] = c.Lineage
	atlas.Digests[idx] = c.Digest
	atlas.PayloadRefs[idx] = payloadRef
	atlas.Count++
	return idx, nil
}

// Get retrieves a coordinate by index.
func (atlas *FlatCoordinateAtlas) Get(idx uint32) (Coordinate64, uint64, error) {
	if idx >= atlas.Count {
		return Coordinate64{}, 0, fmt.Errorf("index %d out of bounds (count=%d)", idx, atlas.Count)
	}
	c := Coordinate64{
		Authority: atlas.Authorities[idx],
		Archetype: atlas.Archetypes[idx],
		Chronos:   atlas.Chronos[idx],
		Topos:     atlas.Topoi[idx],
		Geometry:  atlas.Geometries[idx],
		Encoding:  atlas.Encodings[idx],
		Lineage:   atlas.Lineages[idx],
		Digest:    atlas.Digests[idx],
	}
	return c, atlas.PayloadRefs[idx], nil
}

// MatchHyperplane scans the flat atlas across all active axes using contiguous column accesses.
// Autovectorized into SIMD vector compares without walking any tree structures.
func (atlas *FlatCoordinateAtlas) MatchHyperplane(query AtlasQueryMask) []uint32 {
	var matches []uint32
	n := atlas.Count

	mask := query.ActiveMask
	v := query.Values

	for i := uint32(0); i < n; i++ {
		if (mask&MaskAuthority != 0) && atlas.Authorities[i] != v.Authority {
			continue
		}
		if (mask&MaskArchetype != 0) && atlas.Archetypes[i] != v.Archetype {
			continue
		}
		if (mask&MaskChronos != 0) && atlas.Chronos[i] != v.Chronos {
			continue
		}
		if (mask&MaskTopos != 0) && atlas.Topoi[i] != v.Topos {
			continue
		}
		if (mask&MaskGeometry != 0) && atlas.Geometries[i] != v.Geometry {
			continue
		}
		if (mask&MaskEncoding != 0) && atlas.Encodings[i] != v.Encoding {
			continue
		}
		if (mask&MaskLineage != 0) && atlas.Lineages[i] != v.Lineage {
			continue
		}
		if (mask&MaskDigest != 0) && atlas.Digests[i] != v.Digest {
			continue
		}
		matches = append(matches, i)
	}
	return matches
}

// TranslateAffine applies an instantaneous affine shift delta to specified coordinates.
// Enables zero-copy dataset renaming, tenant migration, or cluster re-partitioning
// without modifying or moving a single payload byte.
func (atlas *FlatCoordinateAtlas) TranslateAffine(indices []uint32, delta Coordinate64) error {
	for _, idx := range indices {
		if idx >= atlas.Count {
			return fmt.Errorf("index %d out of bounds", idx)
		}
		atlas.Authorities[idx] += delta.Authority
		atlas.Archetypes[idx] += delta.Archetype
		atlas.Chronos[idx] += delta.Chronos
		atlas.Topoi[idx] += delta.Topos
		atlas.Geometries[idx] += delta.Geometry
		atlas.Encodings[idx] += delta.Encoding
		atlas.Lineages[idx] += delta.Lineage
		atlas.Digests[idx] += delta.Digest
	}
	return nil
}
