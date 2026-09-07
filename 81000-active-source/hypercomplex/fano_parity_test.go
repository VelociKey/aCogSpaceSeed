package hypercomplex

import (
	"testing"
)

func TestFanoStripe7_SelfHealing(t *testing.T) {
	d0 := NewOctonion(10, 20, 30, 40, 50, 60, 70, 80)
	d1 := NewOctonion(11, 21, 31, 41, 51, 61, 71, 81)
	d2 := NewOctonion(12, 22, 32, 42, 52, 62, 72, 82)
	d3 := NewOctonion(13, 23, 33, 43, 53, 63, 73, 83)

	stripe := NewFanoStripe7(d0, d1, d2, d3)

	// 1. Verify authentic stripe
	if !stripe.VerifyStripe() {
		t.Fatalf("expected authentic stripe to pass verification")
	}

	// 2. Test Bit-Rot Detection on Shard 0
	originalS0 := stripe.Shards[0]
	stripe.Shards[0].E[3] ^= 0x01
	if stripe.VerifyStripe() {
		t.Fatalf("expected stripe verification to FAIL on bit-rot in shard 0")
	}
	stripe.Shards[0] = originalS0 // restore

	// 3. Test Recovery of Shards 0, 1, 2
	rec0, err := stripe.ReconstructShard(0, 0)
	if err != nil || rec0 != d0 {
		t.Fatalf("failed reconstructing shard 0: %v, got %v", err, rec0)
	}

	rec1, err := stripe.ReconstructShard(1, 0)
	if err != nil || rec1 != d1 {
		t.Fatalf("failed reconstructing shard 1: %v, got %v", err, rec1)
	}

	rec2, err := stripe.ReconstructShard(2, 0)
	if err != nil || rec2 != d2 {
		t.Fatalf("failed reconstructing shard 2: %v, got %v", err, rec2)
	}

	// 4. Test Multi-Path Recovery of Shard 3 across all 3 intersecting Fano lines!
	for opt := 0; opt < 3; opt++ {
		rec3, err := stripe.ReconstructShard(3, opt)
		if err != nil {
			t.Fatalf("failed reconstructing shard 3 via line option %d: %v", opt, err)
		}
		if rec3 != d3 {
			t.Fatalf("reconstruction mismatch on option %d: got %v, expected %v", opt, rec3, d3)
		}
	}

	// 5. Test Recovery of Parity Shards 4, 5, 6
	for p := 4; p <= 6; p++ {
		recP, err := stripe.ReconstructShard(p, 0)
		if err != nil || recP != stripe.Shards[p] {
			t.Fatalf("failed reconstructing parity shard %d: %v, got %v", p, err, recP)
		}
	}
}
