package qec

import (
	"math"
	"testing"
)

// TestLogicalCodewords verifies that logical |0_L> and |1_L> are properly normalized
// and strictly orthogonal in C^128.
func TestLogicalCodewords(t *testing.T) {
	psi0 := EncodeLogicalZero()
	psi1 := EncodeLogicalOne()

	// Norm of |0_L>
	fid0 := psi0.Fidelity(psi0)
	if math.Abs(fid0-1.0) > 1e-5 {
		t.Fatalf("expected |0_L> norm == 1.0, got %f", fid0)
	}

	// Norm of |1_L>
	fid1 := psi1.Fidelity(psi1)
	if math.Abs(fid1-1.0) > 1e-5 {
		t.Fatalf("expected |1_L> norm == 1.0, got %f", fid1)
	}

	// Orthogonality: <0_L | 1_L> == 0
	fidOverlap := psi0.Fidelity(psi1)
	if fidOverlap > 1e-7 {
		t.Fatalf("expected <0_L|1_L> == 0, got fidelity %e", fidOverlap)
	}
}

// TestAllSingleQubitPauliErrors exhaustively tests all 21 single-qubit errors
// (7 qubits x 3 Pauli errors X, Y, Z) on the Pauli frame engine.
// All 21 cases MUST be detected with exact syndromes and recovered to identity.
func TestAllSingleQubitPauliErrors(t *testing.T) {
	ops := []PauliOp{PauliX, PauliY, PauliZ}

	for q := 0; q < 7; q++ {
		expectedIndex := uint8(q + 1) // 1-based index (1..7)

		for _, op := range ops {
			t.Run(op.String()+"_on_q"+string(rune('0'+q)), func(t *testing.T) {
				frame := NewPauliFrame7()
				frame.ApplyOp(q, op)

				if frame.Weight() != 1 {
					t.Fatalf("expected frame weight 1, got %d", frame.Weight())
				}

				synX, synZ := ExtractSyndrome(frame)

				switch op {
				case PauliX:
					if synX != expectedIndex {
						t.Fatalf("expected synX == %d, got %d", expectedIndex, synX)
					}
					if synZ != 0 {
						t.Fatalf("expected synZ == 0, got %d", synZ)
					}
				case PauliZ:
					if synX != 0 {
						t.Fatalf("expected synX == 0, got %d", synX)
					}
					if synZ != expectedIndex {
						t.Fatalf("expected synZ == %d, got %d", expectedIndex, synZ)
					}
				case PauliY:
					if synX != expectedIndex {
						t.Fatalf("expected synX == %d, got %d", expectedIndex, synX)
					}
					if synZ != expectedIndex {
						t.Fatalf("expected synZ == %d, got %d", expectedIndex, synZ)
					}
				}

				// Perform recovery
				corrected := CorrectPauliFrame(frame, synX, synZ)
				if !corrected.IsIdentity() {
					t.Fatalf("correction failed: remaining frame %s (weight %d)", corrected, corrected.Weight())
				}
			})
		}
	}
}

// TestStateVectorFidelityRecovery tests full 128-dimensional wavefunctions.
// An error is injected into |0_L>, syndromes are extracted, the correction
// is applied to the wave function, and fidelity with initial |0_L> must be 1.000000.
func TestStateVectorFidelityRecovery(t *testing.T) {
	psiInitial := EncodeLogicalZero()

	for q := 0; q < 7; q++ {
		for _, op := range []PauliOp{PauliX, PauliY, PauliZ} {
			psiCorrupt := psiInitial
			psiCorrupt.ApplyPauliToState(q, op)

			// Construct tracking frame to extract syndrome
			frame := NewPauliFrame7()
			frame.ApplyOp(q, op)
			synX, synZ := ExtractSyndrome(frame)

			// Apply recovery operator to the wavefunction
			psiRecovered := psiCorrupt
			if synX > 0 {
				psiRecovered.ApplyPauliToState(int(synX-1), PauliX)
			}
			if synZ > 0 {
				psiRecovered.ApplyPauliToState(int(synZ-1), PauliZ)
			}

			fid := psiInitial.Fidelity(psiRecovered)
			if math.Abs(fid-1.0) > 1e-4 {
				t.Fatalf("statevector recovery failed for %s on qubit %d: fidelity = %f", op, q, fid)
			}
		}
	}
}

// TestFanoAssociatorTripwire verifies that the 7 canonical lines in the Fano plane
// have exactly zero octonionic associator ([e_a, e_b, e_c] == 0), whereas non-collinear
// triples explode to non-zero associators.
func TestFanoAssociatorTripwire(t *testing.T) {
	canonicalLines := [7][3]int{
		{1, 2, 3},
		{1, 4, 5},
		{1, 7, 6},
		{2, 4, 6},
		{2, 5, 7},
		{3, 4, 7},
		{3, 6, 5},
	}

	// 1. All 7 lines must have zero associators (associative quaternionic subalgebras)
	for i, line := range canonicalLines {
		isZero, assoc := CheckFanoAssociatorTripwire(line[0], line[1], line[2])
		if !isZero {
			t.Fatalf("line %d (%v) failed: associator %s is non-zero", i, line, assoc)
		}
	}

	// 2. Non-collinear triples must be non-zero (non-associative tripwire)
	nonCollinear := [][3]int{
		{1, 2, 4},
		{1, 3, 5},
		{2, 3, 4},
		{4, 5, 6},
		{5, 6, 7},
	}

	for _, triple := range nonCollinear {
		isZero, assoc := CheckFanoAssociatorTripwire(triple[0], triple[1], triple[2])
		if isZero {
			t.Fatalf("expected non-collinear triple %v to trip associator, but got zero (%s)", triple, assoc)
		}
	}
}

// TestMonteCarloSimulation verifies the statistical performance of the QEC engine.
func TestMonteCarloSimulation(t *testing.T) {
	shots := 20000
	noise := 0.01 // 1% physical error rate per qubit
	telem := RunMonteCarloBatch(shots, noise)

	if telem.TotalShots != shots {
		t.Fatalf("expected %d shots, got %d", shots, telem.TotalShots)
	}

	if telem.TotalErrorsInjected == 0 {
		t.Fatalf("expected errors to be injected at 1%% noise rate")
	}

	if telem.SingleErrorsCorrected == 0 {
		t.Fatalf("expected single errors to be corrected")
	}

	// At 1% noise, single errors dominate (~7%), and 100% of single errors are corrected.
	// Fidelity should be > 0.99
	if telem.FidelityMean < 0.99 {
		t.Fatalf("expected fidelity > 0.99, got %f", telem.FidelityMean)
	}
}

// BenchmarkExtractAndCorrect benchmarks zero-lookup syndrome decoding latency.
func BenchmarkExtractAndCorrect(b *testing.B) {
	frame := NewPauliFrame7()
	frame.ApplyOp(4, PauliY) // corrupt qubit 4 with Y error

	b.ResetTimer()
	b.ReportAllocs()

	for i := 0; i < b.N; i++ {
		synX, synZ := ExtractSyndrome(frame)
		_ = CorrectPauliFrame(frame, synX, synZ)
	}
}
