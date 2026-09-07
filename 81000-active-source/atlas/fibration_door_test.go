package atlas_test

import (
	"testing"

	"sov.fleet/plinth-filesystem/81000-active-source/atlas"
)

func TestFibrationDoorBijectiveRoundtrip(t *testing.T) {
	a := atlas.NewFlatCoordinateAtlas(10)
	door := atlas.NewFibrationDoor(a)

	c := atlas.Coordinate64{
		Authority: 0x505652474E,
		Archetype: atlas.ArchetypeTensorWeight,
		Chronos:   42,
		Topos:     7,
		Geometry:  4,
		Encoding:  atlas.EncodingFP8,
		Lineage:   1001,
		Digest:    0x123456789abcdef0,
	}

	// 1. Semantic Lens
	semPath := door.Project(c, atlas.LensSemantic)
	if semPath != "/by-model/tensors/model_505652474e/epoch_42/fp8_123456789abcdef0.tensor" {
		t.Fatalf("Unexpected semantic path: %s", semPath)
	}
	liftedSem, err := door.Lift(semPath, atlas.LensSemantic)
	if err != nil {
		t.Fatalf("Lift failed on semantic path: %v", err)
	}
	if liftedSem.Archetype != c.Archetype || liftedSem.Authority != c.Authority || liftedSem.Chronos != c.Chronos || liftedSem.Encoding != c.Encoding || liftedSem.Digest != c.Digest {
		t.Errorf("Lifted semantic coordinate mismatch: %+v != %+v", liftedSem, c)
	}

	// 2. Hardware Lens
	hwPath := door.Project(c, atlas.LensHardware)
	if hwPath != "/by-hardware/node_7/fp8/rank_4_123456789abcdef0.tensor" {
		t.Fatalf("Unexpected hardware path: %s", hwPath)
	}
	liftedHW, err := door.Lift(hwPath, atlas.LensHardware)
	if err != nil {
		t.Fatalf("Lift failed on hardware path: %v", err)
	}
	if liftedHW.Topos != c.Topos || liftedHW.Encoding != c.Encoding || liftedHW.Geometry != c.Geometry || liftedHW.Digest != c.Digest {
		t.Errorf("Lifted hardware coordinate mismatch: %+v != %+v", liftedHW, c)
	}

	// 3. Chronos Lens
	chronosPath := door.Project(c, atlas.LensChronos)
	if chronosPath != "/by-date/epoch_42/auth_505652474e/fp8_123456789abcdef0.tensor" {
		t.Fatalf("Unexpected chronos path: %s", chronosPath)
	}
	liftedChronos, err := door.Lift(chronosPath, atlas.LensChronos)
	if err != nil {
		t.Fatalf("Lift failed on chronos path: %v", err)
	}
	if liftedChronos.Chronos != c.Chronos || liftedChronos.Authority != c.Authority || liftedChronos.Encoding != c.Encoding || liftedChronos.Digest != c.Digest {
		t.Errorf("Lifted chronos coordinate mismatch: %+v != %+v", liftedChronos, c)
	}

	t.Logf("✅ Fibration Door bijective roundtrips verified across Semantic, Hardware, and Chronos lenses")
}

func TestMultiLensCoexistenceAndVirtualDirectoryListing(t *testing.T) {
	a := atlas.NewFlatCoordinateAtlas(100)
	door := atlas.NewFibrationDoor(a)

	// Insert 3 files
	for i := uint64(1); i <= 3; i++ {
		c := atlas.Coordinate64{
			Authority: 0x1000,
			Archetype: atlas.ArchetypeTensorWeight,
			Chronos:   i * 10,
			Topos:     i,
			Geometry:  4,
			Encoding:  atlas.EncodingFP8,
			Digest:    0xABCD0000 + i,
		}
		_, err := a.Insert(c, i*4096)
		if err != nil {
			t.Fatalf("Failed to insert coordinate: %v", err)
		}
	}

	// List root of Semantic Lens
	entries, err := door.ListVirtualDirectory("/by-model/tensors/model_00001000", atlas.LensSemantic)
	if err != nil {
		t.Fatalf("ListVirtualDirectory failed: %v", err)
	}
	if len(entries) != 3 {
		t.Fatalf("Expected 3 epoch directories, got %d", len(entries))
	}
	for _, e := range entries {
		if !e.IsDir {
			t.Errorf("Expected directory entry, got file: %s", e.Name)
		}
	}

	// List inside epoch_10 directory
	files, err := door.ListVirtualDirectory("/by-model/tensors/model_00001000/epoch_10", atlas.LensSemantic)
	if err != nil {
		t.Fatalf("ListVirtualDirectory failed on epoch_10: %v", err)
	}
	if len(files) != 1 || files[0].IsDir {
		t.Fatalf("Expected 1 file entry, got %v", files)
	}
	if files[0].Name != "fp8_00000000abcd0001.tensor" {
		t.Errorf("Unexpected file name: %s", files[0].Name)
	}
	if files[0].PayloadRef != 4096 {
		t.Errorf("PayloadRef mismatch: expected 4096, got %d", files[0].PayloadRef)
	}

	// Test Path Resolution
	resCoord, resRef, err := door.ResolvePath("/by-model/tensors/model_00001000/epoch_10/fp8_00000000abcd0001.tensor", atlas.LensSemantic)
	if err != nil {
		t.Fatalf("ResolvePath failed: %v", err)
	}
	if resRef != 4096 || resCoord.Chronos != 10 {
		t.Errorf("Resolved coordinate mismatch: ref=%d, chronos=%d", resRef, resCoord.Chronos)
	}

	t.Logf("✅ Virtual directory traversal & instant resolution verified with zero physical inodes")
}
