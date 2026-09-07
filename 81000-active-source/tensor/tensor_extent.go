// Package tensor provides native multi-dimensional Tensor-Strided Extent descriptors
// and zero-copy tile slicing for plinth-filesystem.
//
// DESIGN PRINCIPLE: DATA IS A TENSOR, NOT A 1D TAPE.
// Replaces flat byte-stream offsets with N-dimensional tensor geometry matching NVIDIA Blackwell
// Tensor Core GEMM tiling, MoE expert slicing, and NVFP4/FP8/BF16 microscaling memory formats.
package tensor

import (
	"errors"
	"fmt"
)

// Tensor Data Types (32-bit enum)
const (
	DTypeUnknown uint32 = 0
	DTypeFP4     uint32 = 1 // 4-bit Floating Point (Blackwell Microscaling)
	DTypeFP8     uint32 = 2 // 8-bit Floating Point (E4M3 / E5M2)
	DTypeBF16    uint32 = 3 // 16-bit Brain Floating Point
	DTypeFP16    uint32 = 4 // 16-bit Standard IEEE Float
	DTypeFP32    uint32 = 5 // 32-bit Standard IEEE Float
	DTypeINT8    uint32 = 6 // 8-bit Integer Quantization
	DTypeM31     uint32 = 7 // 31-bit Small-Field Mersenne Prime Field
)

var (
	ErrRankExceedsMax    = errors.New("tensor rank exceeds maximum supported dimensions (8)")
	ErrCoordinatesLength = errors.New("coordinate slice length must match tensor rank")
	ErrIndexOutOfBounds  = errors.New("tensor coordinate index exceeds dimension bounds")
	ErrSliceOutOfBounds  = errors.New("slice range exceeds target dimension length")
)

// ElementSizeBytes returns the element footprint in bytes (rounding sub-byte types up to byte fraction).
func ElementSizeBytes(dtype uint32) float64 {
	switch dtype {
	case DTypeFP4:
		return 0.5 // 2 elements per byte
	case DTypeFP8, DTypeINT8:
		return 1.0
	case DTypeBF16, DTypeFP16:
		return 2.0
	case DTypeFP32, DTypeM31:
		return 4.0
	default:
		return 1.0
	}
}

// FlatTensorExtent defines a multi-dimensional tensor projected onto a physical ZNS zone.
type FlatTensorExtent struct {
	DType          uint32    // Quantization / precision format
	Rank           uint32    // Number of dimensions (1..8)
	Dims           [8]uint32 // Shape per dimension
	Strides        [8]uint64 // Byte stride per dimension
	ZoneID         uint32    // Target ZNS flash zone
	ZoneOffset     uint64    // Physical base byte offset
	TileAlignBytes uint32    // Hardware alignment (64B cacheline or 4KB GEMM tile)
	TotalBytes     uint64    // Total physical allocation
}

func alignUp(val, align uint64) uint64 {
	if align == 0 || val%align == 0 {
		return val
	}
	return ((val / align) + 1) * align
}

// NewFlatTensorExtent creates an N-dimensional tensor extent with automatic row-major stride computation
// and hardware tile alignment.
func NewFlatTensorExtent(zoneID uint32, zoneOffset uint64, dtype uint32, dims []uint32, tileAlign uint32) (*FlatTensorExtent, error) {
	rank := len(dims)
	if rank == 0 || rank > 8 {
		return nil, ErrRankExceedsMax
	}
	if tileAlign == 0 {
		tileAlign = 64 // 64-byte default L3 cacheline
	}

	extent := &FlatTensorExtent{
		DType:          dtype,
		Rank:           uint32(rank),
		ZoneID:         zoneID,
		ZoneOffset:     zoneOffset,
		TileAlignBytes: tileAlign,
	}

	for i := 0; i < rank; i++ {
		extent.Dims[i] = dims[i]
	}

	// Compute row-major byte strides backwards
	elemSize := ElementSizeBytes(dtype)
	var currentStride uint64 = uint64(elemSize)
	if elemSize < 1.0 {
		currentStride = 1 // Pack sub-byte types
	}

	for i := rank - 1; i >= 0; i-- {
		extent.Strides[i] = currentStride
		currentStride *= uint64(dims[i])
	}

	extent.TotalBytes = alignUp(currentStride, uint64(tileAlign))
	return extent, nil
}

// ComputeByteOffset calculates the exact linear flash byte offset for an N-dimensional coordinate.
// Executes in O(1) time with ZERO heap allocations.
func (t *FlatTensorExtent) ComputeByteOffset(coords []uint32) (uint64, error) {
	if uint32(len(coords)) != t.Rank {
		return 0, ErrCoordinatesLength
	}

	var offset uint64 = t.ZoneOffset
	for i := uint32(0); i < t.Rank; i++ {
		c := coords[i]
		if c >= t.Dims[i] {
			return 0, ErrIndexOutOfBounds
		}
		offset += uint64(c) * t.Strides[i]
	}

	return offset, nil
}

// SliceSubTensor performs zero-copy sub-tensor slicing (e.g. slicing 1 attention head out of 32,
// or slicing 1 MoE expert out of 8) directly in flash memory.
// Returns a new FlatTensorExtent pointing to the exact physical sub-region with adjusted shape.
func (t *FlatTensorExtent) SliceSubTensor(dimIndex uint32, start uint32, length uint32) (*FlatTensorExtent, error) {
	if dimIndex >= t.Rank {
		return nil, ErrIndexOutOfBounds
	}
	if start+length > t.Dims[dimIndex] {
		return nil, ErrSliceOutOfBounds
	}

	// Calculate new physical base offset
	newOffset := t.ZoneOffset + uint64(start)*t.Strides[dimIndex]

	sub := &FlatTensorExtent{
		DType:          t.DType,
		Rank:           t.Rank,
		ZoneID:         t.ZoneID,
		ZoneOffset:     newOffset,
		TileAlignBytes: t.TileAlignBytes,
		Dims:           t.Dims,
		Strides:        t.Strides,
	}

	// Update sliced dimension
	sub.Dims[dimIndex] = length
	sub.TotalBytes = uint64(length) * t.Strides[dimIndex]

	return sub, nil
}

// String provides a human-readable summary of the tensor geometry.
func (t *FlatTensorExtent) String() string {
	return fmt.Sprintf("FlatTensorExtent[DType=%d, Rank=%d, Dims=%v, Zone=%d, BaseOffset=%d, TotalBytes=%d]",
		t.DType, t.Rank, t.Dims[:t.Rank], t.ZoneID, t.ZoneOffset, t.TotalBytes)
}
