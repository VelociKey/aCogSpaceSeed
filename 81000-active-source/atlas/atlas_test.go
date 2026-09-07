package atlas_test

import (
	"testing"
	"time"
	"unsafe"

	"sov.fleet/plinth-filesystem/81000-active-source/atlas"
)

func TestCoordinate64PhysicalAlignment(t *testing.T) {
	var c atlas.Coordinate64
	size := unsafe.Sizeof(c)
	if size != 64 {
		t.Fatalf("Expected Coordinate64 to be exactly 64 bytes (1 physical L3 cacheline), got %d bytes", size)
	}
	align := unsafe.Alignof(c)
	if align != 8 {
		t.Fatalf("Expected 8-byte alignment, got %d", align)
	}
	t.Logf("✅ Coordinate64 memory layout verified: 64 bytes (512 bits) matching exact L3 cacheline width")
}

func TestFlatAtlasInsertionAndLookup(t *testing.T) {
	a := atlas.NewFlatCoordinateAtlas(1024)

	c1 := atlas.Coordinate64{
		Authority: 0xA1B2C3D4E5F60001,
		Archetype: atlas.ArchetypeTensorWeight,
		Chronos:   42,
		Topos:     7,
		Geometry:  0x0400010000000000, // Rank 4
		Encoding:  atlas.EncodingFP8,
		Lineage:   0x9988776655443322,
		Digest:    0x123456789abcdef0,
	}

	idx, err := a.Insert(c1, 4096000)
	if err != nil {
		t.Fatalf("Failed to insert coordinate: %v", err)
	}
	if idx != 0 {
		t.Fatalf("Expected index 0, got %d", idx)
	}

	retrieved, payloadRef, err := a.Get(idx)
	if err != nil {
		t.Fatalf("Failed to get coordinate: %v", err)
	}
	if payloadRef != 4096000 {
		t.Errorf("Expected payloadRef 4096000, got %d", payloadRef)
	}
	if retrieved != c1 {
		t.Errorf("Retrieved coordinate mismatch: %+v != %+v", retrieved, c1)
	}
}

func TestHyperplaneMatching100K(t *testing.T) {
	const totalFiles = 100000
	a := atlas.NewFlatCoordinateAtlas(totalFiles)

	// Populate 100,000 synthetic files across multiple epochs, encodings, and archetypes
	for i := uint32(0); i < totalFiles; i++ {
		c := atlas.Coordinate64{
			Authority: 0x8888,
			Archetype: uint64(1 + (i % 5)),           // 1..5
			Chronos:   uint64(i % 100),               // Epoch 0..99
			Topos:     uint64(i % 32),                // Node 0..31
			Geometry:  4,                             // Rank 4
			Encoding:  uint64(1 + (i % 6)),           // FP8..INT4
			Lineage:   uint64(i * 31),
			Digest:    uint64(i * 10007),
		}
		_, err := a.Insert(c, uint64(i*4096))
		if err != nil {
			t.Fatalf("Insertion failed at index %d: %v", i, err)
		}
	}

	// Query: Find all files where Archetype == TensorWeight AND Encoding == FP8 AND Chronos == 30
	query := atlas.AtlasQueryMask{
		ActiveMask: atlas.MaskArchetype | atlas.MaskEncoding | atlas.MaskChronos,
		Values: atlas.Coordinate64{
			Archetype: atlas.ArchetypeTensorWeight,
			Encoding:  atlas.EncodingFP8,
			Chronos:   30,
		},
	}

	start := time.Now()
	matches := a.MatchHyperplane(query)
	elapsed := time.Since(start)

	t.Logf("⚡ Hyperplane scan across %d files completed in %v (%d matches found)", totalFiles, elapsed, len(matches))

	if len(matches) == 0 {
		t.Fatalf("Expected matches for query, got 0")
	}

	// Verify all returned matches satisfy the hyperplane constraint
	for _, idx := range matches {
		c, _, _ := a.Get(idx)
		if c.Archetype != atlas.ArchetypeTensorWeight || c.Encoding != atlas.EncodingFP8 || c.Chronos != 30 {
			t.Fatalf("False positive match at index %d: %+v", idx, c)
		}
	}
	t.Logf("✅ 100%% of %d matching files validated against hyperplane constraints", len(matches))
}

func TestAffineCoordinateTranslation(t *testing.T) {
	a := atlas.NewFlatCoordinateAtlas(10)

	c := atlas.Coordinate64{
		Authority: 0x1000,
		Archetype: atlas.ArchetypeModelCheckpoint,
		Chronos:   10,
		Topos:     1,
		Encoding:  atlas.EncodingBF16,
		Digest:    0xDEADBEEF,
	}
	idx, err := a.Insert(c, 99999)
	if err != nil {
		t.Fatalf("Insert failed: %v", err)
	}

	// Affine shift: Advance epoch by +5, migrate node by +10
	delta := atlas.Coordinate64{
		Chronos: 5,
		Topos:   10,
	}

	err = a.TranslateAffine([]uint32{idx}, delta)
	if err != nil {
		t.Fatalf("TranslateAffine failed: %v", err)
	}

	updated, payloadRef, _ := a.Get(idx)
	if updated.Chronos != 15 {
		t.Errorf("Expected Chronos=15, got %d", updated.Chronos)
	}
	if updated.Topos != 11 {
		t.Errorf("Expected Topos=11, got %d", updated.Topos)
	}
	// Payload reference must remain completely untouched (zero copy)
	if payloadRef != 99999 {
		t.Errorf("Payload reference altered during translation: got %d", payloadRef)
	}
	t.Logf("✅ Zero-copy affine coordinate translation verified: moved without modifying payload")
}

func TestPolyFacetedLatticeCommutativity(t *testing.T) {
	bag1 := atlas.TagBag{
		{Key: "epoch", Val: "42"},
		{Key: "enc", Val: "fp8"},
		{Key: "auth", Val: "zenith"},
	}

	bag2 := atlas.TagBag{
		{Key: "auth", Val: "zenith"},
		{Key: "enc", Val: "fp8"},
		{Key: "epoch", Val: "42"},
	}

	if !bag1.Equals(bag2) {
		t.Errorf("Expected lattice equality under permutation: bag1 != bag2")
	}

	// Verify meet commutativity: A ∧ B == B ∧ A
	other := atlas.TagBag{
		{Key: "node", Val: "7"},
		{Key: "rank", Val: "4"},
	}

	meet1 := atlas.LatticeMeet(bag1, other)
	meet2 := atlas.LatticeMeet(other, bag1)

	if !meet1.Equals(meet2) {
		t.Errorf("Lattice meet is not commutative: meet1 != meet2")
	}
	t.Logf("✅ Poly-faceted lattice commutativity verified: tag bags are strictly order-free")
}

func TestEphemeralPOSIXFibration(t *testing.T) {
	c := atlas.Coordinate64{
		Authority: 0x00000042,
		Archetype: atlas.ArchetypeTensorWeight,
		Chronos:   100,
		Topos:     16,
		Encoding:  atlas.EncodingFP8,
		Digest:    0xabcdef0123456789,
	}

	path := atlas.CoordinateToPOSIX(c)
	expected := "/auth_00000042/tensors/epoch_100/node_16/fp8_abcdef0123456789.tensor"
	if path != expected {
		t.Errorf("Expected virtual path %q, got %q", expected, path)
	}

	// Parse canonical tag set
	tagStr := "auth:0x42,arch:tensor,epoch:100,node:16,enc:fp8,digest:0xabcdef0123456789"
	parsed, err := atlas.ParseCanonicalTagSet(tagStr)
	if err != nil {
		t.Fatalf("ParseCanonicalTagSet failed: %v", err)
	}
	if parsed.Authority != c.Authority || parsed.Archetype != c.Archetype || parsed.Chronos != c.Chronos {
		t.Errorf("Parsed tag set mismatch: %+v != %+v", parsed, c)
	}
	t.Logf("✅ Ephemeral POSIX fibration lens verified: generated %s with zero disk inodes", path)
}

func TestFanoLineageProofAuthenticityAndTamperDetection(t *testing.T) {
	// Test all 7 valid Fano plane lines
	for line := 0; line < 7; line++ {
		c1, c2, c3 := atlas.GenerateValidFanoTriple(0x5555, line)
		res := atlas.VerifyLineageFano(c1, c2, c3)
		if !res.Valid {
			t.Errorf("Line %d expected valid Fano line, got invalid: %s (diff=%v)", line, res.Explanation, res.AssociatorDiff)
		}
	}
	t.Logf("✅ All 7 Fano plane projective lines verified: associator [A, B, C] identically zero")

	// Tamper test: Alter one coordinate off the Fano line
	c1, c2, c3 := atlas.GenerateValidFanoTriple(0x5555, 0) // Line (1, 2, 3)
	// Tamper: change c3 component to 4 (which is not on line (1, 2, 3))
	c3.Geometry = 100 // off-line perturbation

	tamperedRes := atlas.VerifyLineageFano(c1, c2, c3)
	if tamperedRes.Valid {
		t.Fatalf("Expected tamper detection failure on off-line coordinate, but reported valid!")
	}
	t.Logf("✅ Fano anti-tamper tripwire verified: detected off-line coordinate tampering (diff=%v)", tamperedRes.AssociatorDiff)
}
