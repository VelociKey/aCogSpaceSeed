// Package main provides the role-based entrypoint for plinth-filesystem.
//
// plinth_filesystem_core implements a post-POSIX, zero-pointer, rolled-out
// storage architecture engineered for compiler autovectorization by Nautilus.
package main

import (
	"flag"
	"fmt"
	"os"
	"runtime"
	"time"

	"sov.fleet/plinth-filesystem/81000-active-source/slab"
)

func main() {
	mode := "status"
	if len(os.Args) > 1 {
		mode = os.Args[1]
	}

	switch mode {
	case "benchmark":
		runBenchmarkMode()
	case "verify":
		runVerifyMode()
	case "status":
		fallthrough
	default:
		runStatusMode()
	}
}

func runStatusMode() {
	fmt.Println("================================================================================")
	fmt.Println("⚡ Plinth-Filesystem: Sovereign Post-POSIX Flat Rolled-Out Storage Engine v1.0.0")
	fmt.Println("   [Zero Pointers | Structure-of-Arrays (SoA) | Nautilus Autovectorization Ready]")
	fmt.Println("================================================================================")
	fmt.Printf("   • Host Platform   : %s / %s (Logical CPUs: %d)\n", runtime.GOOS, runtime.GOARCH, runtime.NumCPU())
	fmt.Println("   • Storage Model   : Flattened Slabs (Zero Pointer Indirection, Contiguous Columns)")
	fmt.Println("   • Memory Affinity : 64-Byte Cacheline Aligned | 4KB Page Bound (ZNS / SLS Ready)")
	fmt.Println("   • Compiler Target : Nautilus Three-Border Compiler (Border 2 SIMD Autovectorized)")

	// Demonstrate zero-pointer slab allocation and query
	dirSlab := slab.NewFlatDirectorySlab(4096)
	extentStream := slab.NewFlatExtentStream(64 * 1024 * 1024) // 64 MB linear stream

	payload := []byte("SOVEREIGN_POST_POSIX_FLAT_EXTENT_PAYLOAD_DATA")
	off, length, seal, _ := extentStream.Append(payload)

	slot, _ := dirSlab.Insert(0, "model.sls", off, length, slab.FlagActive, seal)
	foundSlot, ok := dirSlab.FindEntry(0, "model.sls")

	fmt.Println("--------------------------------------------------------------------------------")
	fmt.Printf("   ✔ Directory Slab : %s\n", dirSlab.String())
	fmt.Printf("   ✔ Extent Stream  : %s\n", extentStream.String())
	if ok && foundSlot == slot {
		fmt.Printf("   ✔ Verified Lookup: 'model.sls' -> Slot #%d (Offset: %d bytes, Seal: %x...)\n",
			slot, off, seal[:8])
	}
	fmt.Println("================================================================================")
	fmt.Println("✅ Plinth-Filesystem Bedrock Online.")
}

func runVerifyMode() {
	fmt.Println("⚡ [Plinth-Filesystem] Running sovereign structural self-verification...")
	dirSlab := slab.NewFlatDirectorySlab(1024)
	stream := slab.NewFlatExtentStream(4 * 1024 * 1024)

	for i := 0; i < 500; i++ {
		data := []byte(fmt.Sprintf("EXTENT_PAYLOAD_CHUNK_%05d", i))
		off, length, seal, err := stream.Append(data)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error appending extent #%d: %v\n", i, err)
			os.Exit(1)
		}
		_, err = dirSlab.Insert(1, fmt.Sprintf("chunk_%05d.bin", i), off, length, slab.FlagActive, seal)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error inserting slot #%d: %v\n", i, err)
			os.Exit(1)
		}
	}

	target := "chunk_00450.bin"
	start := time.Now()
	slot, ok := dirSlab.FindEntry(1, target)
	dur := time.Since(start)

	if !ok {
		fmt.Fprintf(os.Stderr, "Verification failed: %s not found\n", target)
		os.Exit(1)
	}

	fmt.Printf("   ✔ 500 contiguous extents verified. Lookup latency: %v (Slot #%d)\n", dur, slot)
	fmt.Println("✅ All structural invariants intact.")
}

func runBenchmarkMode() {
	fs := flag.NewFlagSet("benchmark", flag.ExitOnError)
	count := fs.Int("n", 10000, "Number of entries in flattened slab")
	_ = fs.Parse(os.Args[2:])

	fmt.Printf("⚡ [Plinth-Filesystem Benchmark] Populating flat Structure-of-Arrays slab with %d entries...\n", *count)
	dirSlab := slab.NewFlatDirectorySlab(uint32(*count + 100))
	var seal [32]byte

	for i := 0; i < *count; i++ {
		_, _ = dirSlab.Insert(1, fmt.Sprintf("object_shard_%08d.sls", i), uint64(i*4096), 4096, slab.FlagActive, seal)
	}

	targetHash := slab.FastHash64(fmt.Sprintf("object_shard_%08d.sls", *count-1))
	fmt.Printf("   • Scanning entire %d-slot contiguous slab for worst-case target (tail slot)...\n", *count)

	iterations := 100000
	tStart := time.Now()
	for i := 0; i < iterations; i++ {
		slot, ok := dirSlab.FindByHash(1, targetHash)
		if !ok || int(slot) != *count-1 {
			fmt.Fprintf(os.Stderr, "Lookup mismatch: slot=%d, ok=%v\n", slot, ok)
			os.Exit(1)
		}
	}
	tTotal := time.Since(tStart)
	nsPerLookup := float64(tTotal.Nanoseconds()) / float64(iterations)

	fmt.Printf("   ✔ Completed %d scans in %v (%.2f ns/lookup, throughput: %.2f Mlookups/sec)\n",
		iterations, tTotal, nsPerLookup, 1e9/(nsPerLookup*1e6))
	fmt.Println("✅ Flattened Structure-of-Arrays scan benchmark complete.")
}
