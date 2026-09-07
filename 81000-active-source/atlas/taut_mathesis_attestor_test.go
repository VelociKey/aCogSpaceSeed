package atlas_test

import (
	"testing"
	"unsafe"

	"sov.fleet/plinth-filesystem/81000-active-source/atlas"
)

func TestSelfProvingProofHeaderPhysicalAlignment(t *testing.T) {
	var hdr atlas.SelfProvingProofHeader
	size := unsafe.Sizeof(hdr)
	if size != 64 {
		t.Fatalf("Expected SelfProvingProofHeader to be exactly 64 bytes, got %d", size)
	}
	t.Logf("✅ SelfProvingProofHeader memory size verified: 64 bytes (512-bit register width)")
}

func TestTautMathesisFanoAttestation(t *testing.T) {
	c1, c2, c3 := atlas.GenerateValidFanoTriple(0x505652474E, 0)

	hdr, ok := atlas.AttestFanoCollinearity(c1, c2, c3)
	if !ok {
		t.Fatalf("Expected formal attestation to pass for authentic Fano triple")
	}

	if string(hdr.MagicHeader[:]) != "MATHESIS" {
		t.Errorf("Magic header mismatch: %s", string(hdr.MagicHeader[:]))
	}
	if hdr.White3Witness == 0 {
		t.Errorf("White3Witness must be non-zero")
	}
	if hdr.White3Seal == [32]byte{} {
		t.Errorf("White3Seal must be non-zero")
	}
	if (hdr.InvariantBitmask & atlas.InvariantFanoAssociatorZero) == 0 {
		t.Errorf("Bitmask missing InvariantFanoAssociatorZero flag")
	}

	// Tampered triple must fail attestation
	tampered := c3
	tampered.Geometry = 999 // off-plane perturbation
	_, tamperedOk := atlas.AttestFanoCollinearity(c1, c2, tampered)
	if tamperedOk {
		t.Fatalf("Attestation unexpectedly passed for tampered off-plane coordinate")
	}

	t.Logf("✅ In-process Taut-Mathesis Fano attestation certified with 256-bit White3 seal: %x...", hdr.White3Seal[:8])
}

func TestTautMathesisSoAAndFibrationAttestation(t *testing.T) {
	a := atlas.NewFlatCoordinateAtlas(1024)
	door := atlas.NewFibrationDoor(a)

	hdrSoA, okSoA := atlas.AttestZeroAllocationSoA(a)
	if !okSoA || (hdrSoA.InvariantBitmask&atlas.InvariantZeroHeapAllocation) == 0 {
		t.Fatalf("SoA zero-allocation attestation failed")
	}

	hdrFib, okFib := atlas.AttestFibrationSoundness(door)
	if !okFib || (hdrFib.InvariantBitmask&atlas.InvariantFibrationBijective) == 0 {
		t.Fatalf("Fibration door soundness attestation failed")
	}

	leanScript := atlas.GenerateLean4Script("fano_lineage_associator_zero", hdrSoA)
	if len(leanScript) == 0 {
		t.Fatalf("Generated Lean 4 script is empty")
	}

	t.Logf("✅ Taut-Mathesis zero-allocation and fibration soundness certified with White3 seals")
	t.Logf("✅ Companion Lean 4 dual-audit script emitted successfully")
}
