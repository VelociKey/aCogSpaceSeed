// Package adaptive provides hardware-aware, dynamic storage stratum detection
// and geometric adaptation for plinth-filesystem.
//
// DESIGN PRINCIPLE: ELASTIC SCALE FROM 1GB CONTAINERS TO 100PB BLACKWELL SUPERCLUSTERS.
// Instead of static formatting, plinth-filesystem evaluates available memory, flash capacity,
// and VM cluster group size to select the optimal Structure-of-Arrays slab geometry, zone sizes,
// and M31 algebraic parity policies.
package adaptive

import (
	"fmt"
)

// StorageStratum defines the operational scale category.
type StorageStratum uint8

const (
	StratumNano        StorageStratum = 1 // < 64 GB (Ephemeral containers, edge VMs)
	StratumWorkstation StorageStratum = 2 // 64 GB - 10 TB (Local NVMe workstation)
	StratumScale       StorageStratum = 3 // 10 TB - 1 PB (Multi-NVMe servers, rack scale)
	StratumMegaCluster StorageStratum = 4 // > 1 PB (32x to 320x Blackwell superclusters)
)

// Parity Policy Constants
const (
	ParityNone              uint32 = 0 // Raw high-speed scratch (WAF = 1.0)
	ParitySingleM31         uint32 = 1 // 4+1 or 8+1 M31 linear parity
	ParityDualCauchyM31     uint32 = 2 // 8+2 M31 Cauchy matrix (dual drive/zone tolerance)
	ParityClusterProductM31 uint32 = 3 // 28+4 Pod-Local / 8+2 Super-Spine Product Code
)

// AdaptiveGeometryConfig encapsulates the automatically tuned filesystem parameters.
type AdaptiveGeometryConfig struct {
	Stratum                  StorageStratum
	StratumName              string
	TotalStorageBytes        uint64
	ZoneSizeBytes            uint64
	InitialDirectoryCapacity uint32
	InitialChunkCapacity     uint32
	ParityPolicy             uint32
	InlineThresholdBytes     uint32
	EstimatedRAMFootprintMB  uint32
}

// DetectStratum calculates the optimal geometry configuration based on physical capacity
// and cluster VM count.
func DetectStratum(totalBytes uint64, numVMs uint32) AdaptiveGeometryConfig {
	const (
		GB = 1024 * 1024 * 1024
		TB = 1024 * GB
		PB = 1024 * TB
	)

	// Mega-Cluster rule: > 1 PB or > 32 VMs
	if totalBytes >= 1*PB || numVMs > 32 {
		return AdaptiveGeometryConfig{
			Stratum:                  StratumMegaCluster,
			StratumName:              "MegaCluster (Blackwell SuperPOD)",
			TotalStorageBytes:        totalBytes,
			ZoneSizeBytes:            2 * GB, // 2 GB ZNS zones
			InitialDirectoryCapacity: 16777216,
			InitialChunkCapacity:     16777216,
			ParityPolicy:             ParityClusterProductM31,
			InlineThresholdBytes:     4096,
			EstimatedRAMFootprintMB:  480,
		}
	}

	// Scale rule: 10 TB to 1 PB
	if totalBytes >= 10*TB {
		return AdaptiveGeometryConfig{
			Stratum:                  StratumScale,
			StratumName:              "Scale (Multi-NVMe Enterprise Rack)",
			TotalStorageBytes:        totalBytes,
			ZoneSizeBytes:            1 * GB,
			InitialDirectoryCapacity: 2097152,
			InitialChunkCapacity:     2097152,
			ParityPolicy:             ParityDualCauchyM31,
			InlineThresholdBytes:     1024,
			EstimatedRAMFootprintMB:  96,
		}
	}

	// Workstation rule: 64 GB to 10 TB
	if totalBytes >= 64*GB {
		return AdaptiveGeometryConfig{
			Stratum:                  StratumWorkstation,
			StratumName:              "Workstation (Local NVMe Bedrock)",
			TotalStorageBytes:        totalBytes,
			ZoneSizeBytes:            256 * 1024 * 1024, // 256 MB
			InitialDirectoryCapacity: 262144,
			InitialChunkCapacity:     262144,
			ParityPolicy:             ParitySingleM31,
			InlineThresholdBytes:     512,
			EstimatedRAMFootprintMB:  12,
		}
	}

	// Default Nano rule: < 64 GB
	return AdaptiveGeometryConfig{
		Stratum:                  StratumNano,
		StratumName:              "Nano (Container / Ephemeral Scratch)",
		TotalStorageBytes:        totalBytes,
		ZoneSizeBytes:            16 * 1024 * 1024, // 16 MB
		InitialDirectoryCapacity: 4096,
		InitialChunkCapacity:     4096,
		ParityPolicy:             ParityNone,
		InlineThresholdBytes:     256,
		EstimatedRAMFootprintMB:  2,
	}
}

// EvaluateHeadroom inspects available capacity and suggests parity density shifts.
func EvaluateHeadroom(availableBytes, totalBytes uint64) (headroomRatio float64, recommendedParity uint32) {
	if totalBytes == 0 {
		return 0, ParityNone
	}
	headroomRatio = float64(availableBytes) / float64(totalBytes)

	switch {
	case headroomRatio < 0.15:
		// Headroom tight: recommend single parity to conserve flash
		return headroomRatio, ParitySingleM31
	case headroomRatio > 0.40:
		// Abundant headroom: recommend dual Cauchy for maximum resilience
		return headroomRatio, ParityDualCauchyM31
	default:
		return headroomRatio, ParitySingleM31
	}
}

// String summarizes the adaptive configuration.
func (c AdaptiveGeometryConfig) String() string {
	return fmt.Sprintf("AdaptiveConfig[%s: zone_size=%d MB, dir_cap=%d, chunk_cap=%d, RAM=~%d MB]",
		c.StratumName, c.ZoneSizeBytes/(1024*1024), c.InitialDirectoryCapacity,
		c.InitialChunkCapacity, c.EstimatedRAMFootprintMB)
}
