package zns

import (
	"testing"
)

func TestFlatZoneManager_BasicLifecycle(t *testing.T) {
	// 4 zones of 1MB each
	zoneSize := uint64(1024 * 1024)
	mgr := NewFlatZoneManager(4, zoneSize)

	if mgr.TotalZones != 4 {
		t.Fatalf("expected 4 zones, got %d", mgr.TotalZones)
	}

	// Allocate 100KB extent
	reqLen := uint64(100 * 1024)
	zoneID, offset, alignedLen, err := mgr.AllocateExtent(reqLen)
	if err != nil {
		t.Fatalf("unexpected allocation error: %v", err)
	}
	if zoneID != 0 {
		t.Fatalf("expected zone 0, got %d", zoneID)
	}
	if offset != 0 {
		t.Fatalf("expected offset 0, got %d", offset)
	}
	if alignedLen < reqLen || alignedLen%CachelineAlignBytes != 0 {
		t.Fatalf("expected 64B cacheline alignment, got %d", alignedLen)
	}

	// Verify zone 0 state transitioned to ZoneOpen
	base, cap, wp, state, erases, err := mgr.GetZoneInfo(0)
	if err != nil {
		t.Fatalf("failed to get zone info: %v", err)
	}
	if state != ZoneOpen {
		t.Fatalf("expected state ZoneOpen (%d), got %d", ZoneOpen, state)
	}
	if wp != base+alignedLen {
		t.Fatalf("expected wp=%d, got %d", base+alignedLen, wp)
	}
	if erases != 0 {
		t.Fatalf("expected 0 erases, got %d", erases)
	}
	if cap != zoneSize {
		t.Fatalf("expected cap=%d, got %d", zoneSize, cap)
	}

	// Allocate second extent in same zone
	zoneID2, offset2, alignedLen2, err := mgr.AllocateExtent(200 * 1024)
	if err != nil {
		t.Fatalf("failed second allocation: %v", err)
	}
	if zoneID2 != 0 {
		t.Fatalf("expected second extent in zone 0, got %d", zoneID2)
	}
	if offset2 != alignedLen || alignedLen2 == 0 {
		t.Fatalf("expected offset %d and non-zero alignedLen2, got offset=%d, len=%d", alignedLen, offset2, alignedLen2)
	}

	// Reset Zone 0
	if err := mgr.ResetZone(0); err != nil {
		t.Fatalf("failed to reset zone 0: %v", err)
	}

	_, _, wpReset, stateReset, erasesReset, _ := mgr.GetZoneInfo(0)
	if stateReset != ZoneEmpty {
		t.Fatalf("expected ZoneEmpty after reset, got %d", stateReset)
	}
	if wpReset != base {
		t.Fatalf("expected wp to reset to base %d, got %d", base, wpReset)
	}
	if erasesReset != 1 {
		t.Fatalf("expected 1 erase count, got %d", erasesReset)
	}
}

func TestFlatZoneManager_ZoneExhaustion(t *testing.T) {
	// 2 zones of 64KB each
	zoneSize := uint64(64 * 1024)
	mgr := NewFlatZoneManager(2, zoneSize)

	// Fill zone 0
	_, _, _, err := mgr.AllocateExtent(64 * 1024)
	if err != nil {
		t.Fatalf("failed to fill zone 0: %v", err)
	}

	_, _, _, state0, _, _ := mgr.GetZoneInfo(0)
	if state0 != ZoneFull {
		t.Fatalf("expected zone 0 to be ZoneFull (%d), got %d", ZoneFull, state0)
	}

	// Fill zone 1
	_, _, _, err = mgr.AllocateExtent(64 * 1024)
	if err != nil {
		t.Fatalf("failed to fill zone 1: %v", err)
	}

	_, _, _, state1, _, _ := mgr.GetZoneInfo(1)
	if state1 != ZoneFull {
		t.Fatalf("expected zone 1 to be ZoneFull (%d), got %d", ZoneFull, state1)
	}

	// Attempting another allocation should fail with ErrNoAvailableZones
	_, _, _, err = mgr.AllocateExtent(1024)
	if err != ErrNoAvailableZones {
		t.Fatalf("expected ErrNoAvailableZones, got %v", err)
	}
}
