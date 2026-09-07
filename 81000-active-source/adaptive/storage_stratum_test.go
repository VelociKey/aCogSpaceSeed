package adaptive

import (
	"testing"
)

func TestDetectStratum(t *testing.T) {
	const (
		GB = 1024 * 1024 * 1024
		TB = 1024 * GB
		PB = 1024 * TB
	)

	// Test 1: Nano (4 GB container)
	c1 := DetectStratum(4*GB, 1)
	if c1.Stratum != StratumNano || c1.ZoneSizeBytes != 16*1024*1024 {
		t.Fatalf("expected StratumNano, got %+v", c1)
	}

	// Test 2: Workstation (2 TB NVMe)
	c2 := DetectStratum(2*TB, 1)
	if c2.Stratum != StratumWorkstation || c2.ZoneSizeBytes != 256*1024*1024 {
		t.Fatalf("expected StratumWorkstation, got %+v", c2)
	}

	// Test 3: Scale (50 TB server)
	c3 := DetectStratum(50*TB, 4)
	if c3.Stratum != StratumScale || c3.ZoneSizeBytes != 1*GB {
		t.Fatalf("expected StratumScale, got %+v", c3)
	}

	// Test 4: MegaCluster (320 VMs, 20 PB)
	c4 := DetectStratum(20*PB, 320)
	if c4.Stratum != StratumMegaCluster || c4.ZoneSizeBytes != 2*GB || c4.ParityPolicy != ParityClusterProductM31 {
		t.Fatalf("expected StratumMegaCluster, got %+v", c4)
	}
}

func TestEvaluateHeadroom(t *testing.T) {
	const total = 100 * 1024 * 1024 * 1024 // 100 GB

	// 10 GB free = 10% headroom -> ParitySingleM31 (conserve flash)
	r1, p1 := EvaluateHeadroom(10*1024*1024*1024, total)
	if r1 < 0.09 || r1 > 0.11 || p1 != ParitySingleM31 {
		t.Fatalf("expected 10%% headroom and ParitySingleM31, got ratio=%.2f, parity=%d", r1, p1)
	}

	// 60 GB free = 60% headroom -> ParityDualCauchyM31 (high resilience)
	r2, p2 := EvaluateHeadroom(60*1024*1024*1024, total)
	if r2 < 0.59 || r2 > 0.61 || p2 != ParityDualCauchyM31 {
		t.Fatalf("expected 60%% headroom and ParityDualCauchyM31, got ratio=%.2f, parity=%d", r2, p2)
	}
}
