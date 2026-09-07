package topos

import (
	"crypto/sha256"
	"fmt"
	"testing"
)

func TestFlatChunkIndex_Deduplication(t *testing.T) {
	idx := NewFlatChunkIndex(1024)

	payloadA := []byte("DATA_CHUNK_ALPHA")
	digestA := sha256.Sum256(payloadA)

	// First registration: unique chunk
	slotA, dedup1, err := idx.RegisterChunk(digestA, 1, 0, uint32(len(payloadA)))
	if err != nil {
		t.Fatalf("unexpected registration error: %v", err)
	}
	if dedup1 {
		t.Fatalf("expected dedup1=false on first insert")
	}
	if slotA != 0 {
		t.Fatalf("expected slot 0, got %d", slotA)
	}

	// Second registration of exact same data: should deduplicate
	slotA2, dedup2, err := idx.RegisterChunk(digestA, 1, 0, uint32(len(payloadA)))
	if err != nil {
		t.Fatalf("unexpected registration error on dedup: %v", err)
	}
	if !dedup2 {
		t.Fatalf("expected dedup2=true on duplicate insert")
	}
	if slotA2 != slotA {
		t.Fatalf("expected duplicate to return original slot %d, got %d", slotA, slotA2)
	}

	// Verify refCount = 2
	_, _, _, _, refCount, err := idx.GetChunkInfo(slotA)
	if err != nil {
		t.Fatalf("failed to get chunk info: %v", err)
	}
	if refCount != 2 {
		t.Fatalf("expected refCount=2, got %d", refCount)
	}

	// Insert distinct chunk
	payloadB := []byte("DATA_CHUNK_BETA")
	digestB := sha256.Sum256(payloadB)
	slotB, dedupB, err := idx.RegisterChunk(digestB, 1, 4096, uint32(len(payloadB)))
	if err != nil || dedupB || slotB != 1 {
		t.Fatalf("failed distinct insert: slot=%d, dedup=%v, err=%v", slotB, dedupB, err)
	}

	// Compute Merkle Root across slots [0, 1]
	root := idx.ComputeMerkleRoot([]uint32{slotA, slotB})
	if root == [32]byte{} {
		t.Fatalf("expected non-empty Merkle root")
	}

	// Verify root is deterministic
	root2 := idx.ComputeMerkleRoot([]uint32{slotA, slotB})
	if root != root2 {
		t.Fatalf("expected identical Merkle roots: %x != %x", root, root2)
	}

	// Release chunk A (first decrement)
	freed, err := idx.ReleaseChunk(slotA)
	if err != nil || freed {
		t.Fatalf("expected freed=false on first release (refCount 2->1), got freed=%v, err=%v", freed, err)
	}

	// Release chunk A (second decrement -> tombstone)
	freed, err = idx.ReleaseChunk(slotA)
	if err != nil || !freed {
		t.Fatalf("expected freed=true on second release, got freed=%v, err=%v", freed, err)
	}

	// Lookup for freed chunk should now fail
	_, found := idx.LookupChunk(digestA)
	if found {
		t.Fatalf("expected tombstoned chunk to not be found")
	}
}

func TestFlatChunkIndex_LookupScan(t *testing.T) {
	idx := NewFlatChunkIndex(256)

	var targetDigest [32]byte
	for i := 0; i < 100; i++ {
		d := sha256.Sum256([]byte(fmt.Sprintf("CHUNK_%04d", i)))
		slot, _, err := idx.RegisterChunk(d, 0, uint64(i*4096), 4096)
		if err != nil {
			t.Fatalf("failed inserting chunk %d: %v", i, err)
		}
		if i == 42 {
			targetDigest = d
			if slot != 42 {
				t.Fatalf("expected slot 42, got %d", slot)
			}
		}
	}

	slotFound, found := idx.LookupChunk(targetDigest)
	if !found || slotFound != 42 {
		t.Fatalf("expected lookup to find slot 42, got slot=%d, found=%v", slotFound, found)
	}
}
