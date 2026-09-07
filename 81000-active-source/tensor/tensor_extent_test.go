package tensor

import (
	"testing"
)

func TestFlatTensorExtent_RowMajorAndSlicing(t *testing.T) {
	// 4D Tensor: [Batch=2, Heads=8, Seq=64, Hidden=128] with BF16 (2 bytes each)
	dims := []uint32{2, 8, 64, 128}
	zoneID := uint32(1)
	zoneOffset := uint64(1024 * 1024) // 1 MB base offset

	extent, err := NewFlatTensorExtent(zoneID, zoneOffset, DTypeBF16, dims, 4096)
	if err != nil {
		t.Fatalf("failed to create tensor extent: %v", err)
	}

	if extent.Rank != 4 {
		t.Fatalf("expected rank 4, got %d", extent.Rank)
	}

	// Verify strides:
	// Strides[3] = 2 (elem size)
	// Strides[2] = 128 * 2 = 256
	// Strides[1] = 64 * 256 = 16384
	// Strides[0] = 8 * 16384 = 131072
	if extent.Strides[3] != 2 {
		t.Fatalf("expected stride[3]=2, got %d", extent.Strides[3])
	}
	if extent.Strides[2] != 256 {
		t.Fatalf("expected stride[2]=256, got %d", extent.Strides[2])
	}
	if extent.Strides[1] != 16384 {
		t.Fatalf("expected stride[1]=16384, got %d", extent.Strides[1])
	}
	if extent.Strides[0] != 131072 {
		t.Fatalf("expected stride[0]=131072, got %d", extent.Strides[0])
	}

	// Compute coordinate: [1, 3, 10, 5]
	// Expected offset: zoneOffset + 1*131072 + 3*16384 + 10*256 + 5*2
	// = 1048576 + 131072 + 49152 + 2560 + 10 = 1231370
	coords := []uint32{1, 3, 10, 5}
	off, err := extent.ComputeByteOffset(coords)
	if err != nil {
		t.Fatalf("failed to compute offset: %v", err)
	}
	expectedOffset := zoneOffset + 1*131072 + 3*16384 + 10*256 + 5*2
	if off != expectedOffset {
		t.Fatalf("offset mismatch: got %d, expected %d", off, expectedOffset)
	}

	// Test Zero-Copy Sub-Tensor Slicing: Slice Head #3 (dim 1, start 3, length 1)
	sliced, err := extent.SliceSubTensor(1, 3, 1)
	if err != nil {
		t.Fatalf("failed to slice sub-tensor: %v", err)
	}

	// Sliced base offset should equal zoneOffset + 3*16384 = 1048576 + 49152 = 1097728
	expectedSubBase := zoneOffset + 3*16384
	if sliced.ZoneOffset != expectedSubBase {
		t.Fatalf("sliced base mismatch: got %d, expected %d", sliced.ZoneOffset, expectedSubBase)
	}
	if sliced.Dims[1] != 1 {
		t.Fatalf("expected sliced dimension 1 length=1, got %d", sliced.Dims[1])
	}

	// Accessing element [0, 0, 10, 5] in the sliced sub-tensor should map to the exact same physical byte!
	subCoords := []uint32{1, 0, 10, 5}
	subOff, err := sliced.ComputeByteOffset(subCoords)
	if err != nil {
		t.Fatalf("failed to compute sub-tensor offset: %v", err)
	}
	if subOff != off {
		t.Fatalf("sub-tensor offset %d does not match parent offset %d (zero-copy violated)", subOff, off)
	}
}
