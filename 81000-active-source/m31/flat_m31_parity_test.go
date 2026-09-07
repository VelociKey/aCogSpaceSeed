package m31

import (
	"math/rand"
	"testing"
)

func TestM31_FieldArithmetic(t *testing.T) {
	// Identity: a + 0 = a
	if AddM31(12345, 0) != 12345 {
		t.Fatalf("AddM31 identity failed")
	}

	// Wrap around at M31Prime
	if AddM31(M31Prime-1, 2) != 1 {
		t.Fatalf("AddM31 wrap failed: got %d", AddM31(M31Prime-1, 2))
	}

	// Subtraction identity
	if SubM31(10, 10) != 0 {
		t.Fatalf("SubM31 identity failed")
	}
	// Subtraction wrap
	if SubM31(5, 10) != M31Prime-5 {
		t.Fatalf("SubM31 wrap failed: got %d", SubM31(5, 10))
	}

	// Multiplication & Modular Inverse: a * inv(a) = 1
	testValues := []uint32{2, 3, 7, 1024, 65537, M31Prime - 2}
	for _, a := range testValues {
		inv := InvM31(a)
		prod := MulM31(a, inv)
		if prod != 1 {
			t.Fatalf("modular inverse failed for a=%d: inv=%d, prod=%d", a, inv, prod)
		}
	}
}

func TestM31_ParityAndReconstruction(t *testing.T) {
	const k = 4
	const words = WordsPerPage // 1024 words = 4KB

	r := rand.New(rand.NewSource(42))

	// Generate K data slabs
	dataSlabs := make([][]uint32, k)
	for j := 0; j < k; j++ {
		dataSlabs[j] = make([]uint32, words)
		for i := 0; i < words; i++ {
			dataSlabs[j][i] = uint32(r.Int31n(int32(M31Prime)))
		}
	}

	// Assign Cauchy/Vandermonde coefficients (non-zero)
	coeffs := []uint32{1, 3, 5, 7}

	// Compute Parity Slab
	paritySlab := make([]uint32, words)
	if err := ComputeParitySlab(dataSlabs, coeffs, paritySlab); err != nil {
		t.Fatalf("ComputeParitySlab failed: %v", err)
	}

	// Verify Authentic Parity
	if !VerifyParitySlab(dataSlabs, coeffs, paritySlab) {
		t.Fatalf("expected authentic parity to pass verification")
	}

	// Inject Bit-Rot into Data Slab 2 (word 500)
	originalWord := dataSlabs[2][500]
	dataSlabs[2][500] ^= 0x00000001 // single bit flip

	if VerifyParitySlab(dataSlabs, coeffs, paritySlab) {
		t.Fatalf("expected parity verification to FAIL on bit-rot, but passed")
	}

	// Restore word for reconstruction test
	dataSlabs[2][500] = originalWord

	// Test Shard Reconstruction: Suppose Shard 2 was completely lost
	targetIdx := 2
	targetCoeff := coeffs[targetIdx]
	targetOriginal := make([]uint32, words)
	copy(targetOriginal, dataSlabs[targetIdx])

	// Build surviving sets
	survivingSlabs := make([][]uint32, 0, k-1)
	survivingCoeffs := make([]uint32, 0, k-1)
	for j := 0; j < k; j++ {
		if j != targetIdx {
			survivingSlabs = append(survivingSlabs, dataSlabs[j])
			survivingCoeffs = append(survivingCoeffs, coeffs[j])
		}
	}

	// Reconstruct into fresh target slab
	reconstructedSlab := make([]uint32, words)
	err := ReconstructMissingShard(survivingSlabs, survivingCoeffs, targetCoeff, paritySlab, reconstructedSlab)
	if err != nil {
		t.Fatalf("ReconstructMissingShard failed: %v", err)
	}

	// Assert 100% word-for-word equality
	for i := 0; i < words; i++ {
		if reconstructedSlab[i] != targetOriginal[i] {
			t.Fatalf("reconstruction mismatch at word %d: got %d, expected %d",
				i, reconstructedSlab[i], targetOriginal[i])
		}
	}
}
