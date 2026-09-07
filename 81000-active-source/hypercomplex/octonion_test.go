package hypercomplex

import (
	"testing"
)

func TestOctonion_BasisMultiplication(t *testing.T) {
	e0 := NewOctonion(1, 0, 0, 0, 0, 0, 0, 0)
	e1 := NewOctonion(0, 1, 0, 0, 0, 0, 0, 0)
	e2 := NewOctonion(0, 0, 1, 0, 0, 0, 0, 0)
	e3 := NewOctonion(0, 0, 0, 1, 0, 0, 0, 0)
	e4 := NewOctonion(0, 0, 0, 0, 1, 0, 0, 0)

	// 1. Identity: e0 * a = a
	prod01 := MulOctonion(e0, e1)
	if prod01 != e1 {
		t.Fatalf("expected e0 * e1 = e1, got %v", prod01)
	}

	// 2. Square of imaginary units: e1 * e1 = -e0
	e1_sq := MulOctonion(e1, e1)
	expected_neg_e0 := NewOctonion(-1, 0, 0, 0, 0, 0, 0, 0)
	if e1_sq != expected_neg_e0 {
		t.Fatalf("expected e1^2 = -e0, got %v", e1_sq)
	}

	// 3. Fano cycle 1: e1 * e2 = e3, e2 * e1 = -e3
	e1_e2 := MulOctonion(e1, e2)
	if e1_e2 != e3 {
		t.Fatalf("expected e1 * e2 = e3, got %v", e1_e2)
	}
	e2_e1 := MulOctonion(e2, e1)
	expected_neg_e3 := NewOctonion(0, 0, 0, -1, 0, 0, 0, 0)
	if e2_e1 != expected_neg_e3 {
		t.Fatalf("expected e2 * e1 = -e3, got %v", e2_e1)
	}

	// 4. Non-commutativity / Commutator: [e1, e2] = 2*e3
	comm := ComputeCommutator(e1, e2)
	expectedComm := NewOctonion(0, 0, 0, 2, 0, 0, 0, 0)
	if comm != expectedComm {
		t.Fatalf("expected [e1, e2] = 2e3, got %v", comm)
	}

	// Scalar commutes: [e0, e1] = 0
	comm0 := ComputeCommutator(e0, e1)
	if !comm0.IsZero() {
		t.Fatalf("expected [e0, e1] = 0, got %v", comm0)
	}

	// 5. Associator test:
	// A) Collinear units (e1, e2, e3 lie on line L1) form a Quaternionic subalgebra:
	// Associator [e1, e2, e3] MUST be identically zero!
	assocQuat := ComputeAssociator(e1, e2, e3)
	if !assocQuat.IsZero() {
		t.Fatalf("expected associator of quaternionic line [e1, e2, e3] == 0, got %v", assocQuat)
	}

	// B) Non-collinear units (e1, e2, e4 do not form a Fano line):
	// Associator [e1, e2, e4] MUST BE NON-ZERO! (Proof of non-associativity)
	assocNonQuat := ComputeAssociator(e1, e2, e4)
	if assocNonQuat.IsZero() {
		t.Fatalf("expected non-associative associator [e1, e2, e4] != 0, but got zero")
	}
}
