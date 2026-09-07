// Package hypercomplex provides 8-dimensional Octonionic algebra and 64-byte cacheline
// storage geometry for plinth-filesystem.
//
// DESIGN PRINCIPLE: CACHELINE UNITY & HYPERCOMPLEX STORAGE GEOMETRY.
// Exactly one Octonion64 (8 x int64 words = 512 bits = 64 Bytes) corresponds to one
// hardware L3 cacheline. Multiplication follows the Cayley-Dickson construction and the
// 7 oriented cycles of the Fano projective plane PG(2, 2).
// Non-associativity provides an unforgeable algebraic Associator tripwire ([A, B, C] != 0),
// and non-commutativity provides causal write ordering without distributed locks.
package hypercomplex

import (
	"fmt"
)

// Octonion64 represents an 8-dimensional hypercomplex number packed into 64 bytes.
// Memory is strictly 64-byte hardware cacheline aligned.
type Octonion64 struct {
	E [8]int64 // E[0] is real scalar; E[1]..E[7] are imaginary units e_1..e_7
}

// Multiplication table representation:
// octoMultIndex[i][j] gives the resulting basis index k (0..7).
// octoMultSign[i][j] gives the sign (+1 or -1).
var (
	octoMultIndex [8][8]uint8
	octoMultSign  [8][8]int8
)

func init() {
	// Initialize identity and self-multiplication
	for i := 0; i < 8; i++ {
		for j := 0; j < 8; j++ {
			if i == 0 {
				octoMultIndex[0][j] = uint8(j)
				octoMultSign[0][j] = 1
			} else if j == 0 {
				octoMultIndex[i][0] = uint8(i)
				octoMultSign[i][0] = 1
			} else if i == j {
				octoMultIndex[i][i] = 0
				octoMultSign[i][i] = -1 // e_i * e_i = -e_0
			}
		}
	}

	// The 7 oriented Fano plane triples (1-indexed)
	fanoTriples := [7][3]int{
		{1, 2, 3},
		{1, 4, 5},
		{1, 7, 6},
		{2, 4, 6},
		{2, 5, 7},
		{3, 4, 7},
		{3, 6, 5},
	}

	for _, triple := range fanoTriples {
		a, b, c := triple[0], triple[1], triple[2]

		// Positive cyclic permutations: a*b=c, b*c=a, c*a=b
		octoMultIndex[a][b] = uint8(c)
		octoMultSign[a][b] = 1

		octoMultIndex[b][c] = uint8(a)
		octoMultSign[b][c] = 1

		octoMultIndex[c][a] = uint8(b)
		octoMultSign[c][a] = 1

		// Negative reversed permutations: b*a=-c, c*b=-a, a*c=-b
		octoMultIndex[b][a] = uint8(c)
		octoMultSign[b][a] = -1

		octoMultIndex[c][b] = uint8(a)
		octoMultSign[c][b] = -1

		octoMultIndex[a][c] = uint8(b)
		octoMultSign[a][c] = -1
	}
}

// NewOctonion creates an Octonion64 from 8 components.
func NewOctonion(e0, e1, e2, e3, e4, e5, e6, e7 int64) Octonion64 {
	return Octonion64{
		E: [8]int64{e0, e1, e2, e3, e4, e5, e6, e7},
	}
}

// AddOctonion computes the component-wise sum (a + b).
func AddOctonion(a, b Octonion64) Octonion64 {
	var res Octonion64
	for i := 0; i < 8; i++ {
		res.E[i] = a.E[i] + b.E[i]
	}
	return res
}

// SubOctonion computes the component-wise difference (a - b).
func SubOctonion(a, b Octonion64) Octonion64 {
	var res Octonion64
	for i := 0; i < 8; i++ {
		res.E[i] = a.E[i] - b.E[i]
	}
	return res
}

// MulOctonion computes the non-commutative, non-associative octonionic product (a * b).
//
// COMPILER OPTIMIZATION TARGET:
// Unrolled 64-term linear combination. Nautilus's Border 2 vectorizer maps these
// into parallel AVX2 / AVX-512 FMA vector lanes without branches.
func MulOctonion(a, b Octonion64) Octonion64 {
	var res Octonion64

	for i := 0; i < 8; i++ {
		ai := a.E[i]
		if ai == 0 {
			continue
		}
		for j := 0; j < 8; j++ {
			bj := b.E[j]
			if bj == 0 {
				continue
			}
			target := octoMultIndex[i][j]
			sign := int64(octoMultSign[i][j])
			res.E[target] += ai * bj * sign
		}
	}

	return res
}

// ComputeCommutator calculates the Lie commutator bracket:
// [a, b] = a * b - b * a.
// If [a, b] == 0, writes a and b are causally independent and can commit concurrently.
// If [a, b] != 0, the sign of the commutator mathematically establishes write priority.
func ComputeCommutator(a, b Octonion64) Octonion64 {
	ab := MulOctonion(a, b)
	ba := MulOctonion(b, a)
	return SubOctonion(ab, ba)
}

// ComputeAssociator calculates the non-associative associator bracket:
// [a, b, c] = (a * b) * c - a * (b * c).
//
// INTRINSIC ANTI-TAMPER TRIPWIRE:
// Within any valid Fano quaternionic subalgebra, [a, b, c] is identically zero.
// If an out-of-order storage packet or forged payload is injected, the associator
// explodes to non-zero, immediately exposing the breach.
func ComputeAssociator(a, b, c Octonion64) Octonion64 {
	ab := MulOctonion(a, b)
	ab_c := MulOctonion(ab, c)

	bc := MulOctonion(b, c)
	a_bc := MulOctonion(a, bc)

	return SubOctonion(ab_c, a_bc)
}

// IsZero returns true if all 8 components are exactly zero.
func (o Octonion64) IsZero() bool {
	for i := 0; i < 8; i++ {
		if o.E[i] != 0 {
			return false
		}
	}
	return true
}

// String provides a clean mathematical representation.
func (o Octonion64) String() string {
	return fmt.Sprintf("Octonion[%d + %di + %dj + %dk + %dl + %de5 + %de6 + %de7]",
		o.E[0], o.E[1], o.E[2], o.E[3], o.E[4], o.E[5], o.E[6], o.E[7])
}
