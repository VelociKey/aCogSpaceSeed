// Package slab provides flattened, rolled-out data-oriented storage structures
// for plinth-filesystem.
//
// flat_extent_stream.go provides contiguous, zero-pointer linear payload storage
// matching physical NVMe Zoned Namespaces (ZNS) append semantics and Single-Level Store (SLS)
// memory mappings.
package slab

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"sync"
)

// Hardware Alignment Constants
const (
	CachelineAlignBytes uint64 = 64
	PageAlignBytes      uint64 = 4096
)

var (
	ErrExtentStreamFull = errors.New("flat extent stream capacity exhausted")
	ErrExtentOutOfBounds = errors.New("extent range exceeds stream bounds")
)

func alignUp(val, align uint64) uint64 {
	if align == 0 || val%align == 0 {
		return val
	}
	return ((val / align) + 1) * align
}

// FlatExtentStream represents a contiguous, zero-pointer linear memory/storage stream.
// Reads are 0-copy direct slice windows; writes are sequential append-only bursts (WAF = 1.0).
type FlatExtentStream struct {
	mu sync.RWMutex

	// Contiguous Linear Memory Buffer
	Capacity  uint64
	WriteTail uint64
	Buffer    []byte
}

// NewFlatExtentStream allocates a contiguous byte slab of specified capacity.
func NewFlatExtentStream(capacity uint64) *FlatExtentStream {
	if capacity == 0 {
		capacity = 64 * 1024 * 1024 // 64 MB default slab
	}
	// Round up to 4KB page boundary
	capacity = alignUp(capacity, PageAlignBytes)
	return &FlatExtentStream{
		Capacity:  capacity,
		WriteTail: 0,
		Buffer:    make([]byte, capacity),
	}
}

// Append writes data sequentially to the stream, enforcing 64-byte cacheline alignment.
// Returns the allocated linear byte offset and the SHA-256 / White3 seal.
func (s *FlatExtentStream) Append(data []byte) (offset uint64, length uint32, seal [32]byte, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	dataLen := uint64(len(data))
	startOffset := alignUp(s.WriteTail, CachelineAlignBytes)
	endOffset := startOffset + dataLen

	if endOffset > s.Capacity {
		return 0, 0, [32]byte{}, ErrExtentStreamFull
	}

	copy(s.Buffer[startOffset:endOffset], data)
	s.WriteTail = endOffset

	seal = sha256.Sum256(data)
	return startOffset, uint32(dataLen), seal, nil
}

// ReadWindow returns a direct zero-copy slice window into the contiguous buffer.
// Zero heap allocations, zero copying.
func (s *FlatExtentStream) ReadWindow(offset uint64, length uint32) ([]byte, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	end := offset + uint64(length)
	if end > s.WriteTail || end > s.Capacity {
		return nil, ErrExtentOutOfBounds
	}
	return s.Buffer[offset:end], nil
}

// AllocatedBytes returns the total bytes committed to the stream.
func (s *FlatExtentStream) AllocatedBytes() uint64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.WriteTail
}

// String summarizes the stream utilization.
func (s *FlatExtentStream) String() string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return fmt.Sprintf("FlatExtentStream[tail=%d/%d bytes (%.2f%% used)]",
		s.WriteTail, s.Capacity, float64(s.WriteTail)/float64(s.Capacity)*100)
}
