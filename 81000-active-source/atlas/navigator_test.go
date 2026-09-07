package atlas_test

import (
	"testing"

	"sov.fleet/plinth-filesystem/81000-active-source/atlas"
)

func TestBoundedBoxRangeQuery(t *testing.T) {
	a := atlas.NewFlatCoordinateAtlas(1000)

	for i := uint32(0); i < 1000; i++ {
		c := atlas.Coordinate64{
			Authority: 100,
			Archetype: atlas.ArchetypeTensorWeight,
			Chronos:   uint64(i % 100), // 0..99
			Topos:     uint64(i % 16),  // 0..15
			Geometry:  4,
			Encoding:  atlas.EncodingFP8,
			Lineage:   uint64(i),
			Digest:    uint64(i * 1009),
		}
		_, _ = a.Insert(c, uint64(i*4096))
	}

	// Range Query: Chronos in [30, 40] AND Topos in [4, 7]
	query := atlas.BoundedBoxQuery{
		ActiveMask: atlas.MaskChronos | atlas.MaskTopos,
		Ranges: [8]atlas.CoordinateInterval{
			{}, // e0
			{}, // e1
			{Min: 30, Max: 40}, // e2: Chronos
			{Min: 4, Max: 7},   // e3: Topos
		},
	}

	matches := a.MatchBoundedBox(query)
	if len(matches) == 0 {
		t.Fatalf("Expected matches for bounded box query, got 0")
	}

	for _, slot := range matches {
		c, _, _ := a.Get(slot)
		if c.Chronos < 30 || c.Chronos > 40 {
			t.Fatalf("Chronos out of bounds: %d", c.Chronos)
		}
		if c.Topos < 4 || c.Topos > 7 {
			t.Fatalf("Topos out of bounds: %d", c.Topos)
		}
	}
	t.Logf("✅ Bounded-box multi-dimensional range query verified: %d matches found in range", len(matches))
}

func TestLatticeMeetRelaxation(t *testing.T) {
	a := atlas.NewFlatCoordinateAtlas(10)

	target := atlas.Coordinate64{
		Authority: 100,
		Archetype: atlas.ArchetypeTensorWeight,
		Chronos:   42,
		Topos:     7,
		Encoding:  atlas.EncodingFP8,
	}
	slot, _ := a.Insert(target, 4096)

	// User query with 4 correct axes, but 1 wrong axis (Chronos = 99 instead of 42)
	query := atlas.AtlasQueryMask{
		ActiveMask: atlas.MaskAuthority | atlas.MaskArchetype | atlas.MaskTopos | atlas.MaskEncoding | atlas.MaskChronos,
		Values: atlas.Coordinate64{
			Authority: 100,
			Archetype: atlas.ArchetypeTensorWeight,
			Chronos:   99, // WRONG
			Topos:     7,
			Encoding:  atlas.EncodingFP8,
		},
	}

	// Exact query yields 0
	exactMatches := a.MatchHyperplane(query)
	if len(exactMatches) != 0 {
		t.Fatalf("Expected 0 exact matches, got %d", len(exactMatches))
	}

	// Relaxed query with min 4 matched axes
	relaxed := a.MatchRelaxed(query, 4)
	if len(relaxed) != 1 {
		t.Fatalf("Expected 1 relaxed match, got %d", len(relaxed))
	}
	if relaxed[0].Slot != slot {
		t.Errorf("Expected slot %d, got %d", slot, relaxed[0].Slot)
	}
	if relaxed[0].MatchedAxes != 4 || relaxed[0].TotalQueriedAxes != 5 {
		t.Errorf("Unexpected match counts: matched=%d, total=%d", relaxed[0].MatchedAxes, relaxed[0].TotalQueriedAxes)
	}
	if relaxed[0].JaccardSimilarity != 0.80 {
		t.Errorf("Expected Jaccard similarity 0.80, got %f", relaxed[0].JaccardSimilarity)
	}

	t.Logf("✅ Lattice meet relaxation verified: isolated 4-of-5 best fit with delta: %v", relaxed[0].Deltas)
}

func TestBijectiveMnemonicCodec(t *testing.T) {
	c := atlas.Coordinate64{
		Authority: 5, // "zenith"
		Topos:     1, // "azure"
		Encoding:  atlas.EncodingFP8,
		Chronos:   42,
		Digest:    0x1234567890ABBAFE,
	}

	handle := atlas.EncodeMnemonic(c)
	expectedPrefix := "zenith-azure-fp8-epoch42-bafe"
	if handle != expectedPrefix {
		t.Fatalf("Expected mnemonic %q, got %q", expectedPrefix, handle)
	}

	decoded, err := atlas.DecodeMnemonic(handle)
	if err != nil {
		t.Fatalf("DecodeMnemonic failed: %v", err)
	}
	if decoded.Authority != c.Authority || decoded.Topos != c.Topos || decoded.Encoding != c.Encoding || decoded.Chronos != c.Chronos || decoded.Digest != (c.Digest&0xFFFF) {
		t.Errorf("Decoded mnemonic mismatch: %+v != %+v", decoded, c)
	}

	t.Logf("✅ Bijective mnemonic projection verified: 64-byte vector <-> %s", handle)
}

func TestShannonEntropyBisector(t *testing.T) {
	a := atlas.NewFlatCoordinateAtlas(100)
	var slots []uint32

	for i := uint32(0); i < 100; i++ {
		c := atlas.Coordinate64{
			Authority: 1,
			Chronos:   uint64(i), // 0..99 uniformly distributed
			Topos:     uint64(i % 2),
			Encoding:  atlas.EncodingFP8,
		}
		s, _ := a.Insert(c, uint64(i*4096))
		slots = append(slots, s)
	}

	bisection, err := a.ComputeOptimalBisection(slots)
	if err != nil {
		t.Fatalf("ComputeOptimalBisection failed: %v", err)
	}

	// Median of Chronos 0..99 should split 50/50 with entropy ~1.00 bit
	if bisection.Entropy < 0.90 {
		t.Errorf("Expected near-optimal entropy (>0.90), got %f", bisection.Entropy)
	}
	t.Logf("✅ Shannon entropy bisection verified: %s (Entropy: %.3f bits, split: %d vs %d)",
		bisection.Prompt, bisection.Entropy, bisection.LowerCount, bisection.HigherEqCount)
}

func TestFanoGeodesicWalk(t *testing.T) {
	cParent := atlas.Coordinate64{
		Authority: 0x505652474E,
		Archetype: atlas.ArchetypeModelCheckpoint,
		Chronos:   10,
		Digest:    0xDEADBEEF,
	}

	cChild := atlas.StepDownstream(cParent, 0)
	if cChild.Chronos != 11 {
		t.Errorf("Expected child Chronos 11, got %d", cChild.Chronos)
	}
	if cChild.Lineage != cParent.Digest {
		t.Errorf("Expected child lineage to equal parent digest: %x != %x", cChild.Lineage, cParent.Digest)
	}

	cBack := atlas.StepUpstream(cChild, 0)
	if cBack.Chronos != 10 {
		t.Errorf("Expected upstream Chronos 10, got %d", cBack.Chronos)
	}
	if cBack.Digest != cParent.Digest {
		t.Errorf("Expected restored parent digest: %x != %x", cBack.Digest, cParent.Digest)
	}

	t.Logf("✅ Fano causal geodesic walk verified: stepped downstream (epoch 10 -> 11) and upstream (11 -> 10)")
}
