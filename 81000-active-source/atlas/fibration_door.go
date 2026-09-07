package atlas

import (
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// LensKind specifies the projective lens through which the 8D topos is viewed as a 1D tree.
type LensKind int

const (
	LensSemantic LensKind = iota // /by-model/{archetype}/{model_id}/epoch_{n}/{enc}_{digest}.tensor
	LensHardware                // /by-hardware/node_{topos}/{encoding}/{geom}_{digest}.tensor
	LensChronos                 // /by-date/epoch_{chronos}/auth_{authority}/{enc}_{digest}.tensor
	LensFlat                    // /auth_{authority}/{enc}_{digest}.tensor
)

func (l LensKind) String() string {
	switch l {
	case LensSemantic:
		return "by-model"
	case LensHardware:
		return "by-hardware"
	case LensChronos:
		return "by-date"
	case LensFlat:
		return "flat"
	default:
		return "custom"
	}
}

// VirtualDirEntry represents an ephemeral directory or file entry presented to legacy 1D tools.
// Zero disk inodes or directory nodes are allocated on physical storage media.
type VirtualDirEntry struct {
	Name         string       // Name of child directory or file
	IsDir        bool         // True if intermediate virtual directory
	Size         uint64       // Size in bytes if file
	Coordinate   Coordinate64 // 8D coordinate if file
	PayloadRef   uint64       // Storage extent offset / chunk ID
	White3Digest uint64       // White3-64 invariant state witness (e7)
}

// FibrationDoor manages bijective projection (pi) and lifting (s) between the 8D topos and 1D paths.
type FibrationDoor struct {
	Atlas *FlatCoordinateAtlas
}

// NewFibrationDoor initializes a fibration door connected to a flat coordinate atlas.
func NewFibrationDoor(atlas *FlatCoordinateAtlas) *FibrationDoor {
	return &FibrationDoor{
		Atlas: atlas,
	}
}

// Project computes the forward projection pi: O^8 -> T_1D for a given lens.
func (d *FibrationDoor) Project(c Coordinate64, lens LensKind) string {
	switch lens {
	case LensSemantic:
		archName := ArchetypeName(c.Archetype)
		encName := EncodingName(c.Encoding)
		return fmt.Sprintf("/by-model/%s/model_%08x/epoch_%d/%s_%016x.tensor",
			archName,
			c.Authority,
			c.Chronos,
			encName,
			c.Digest,
		)
	case LensHardware:
		encName := EncodingName(c.Encoding)
		return fmt.Sprintf("/by-hardware/node_%d/%s/rank_%d_%016x.tensor",
			c.Topos,
			encName,
			c.Geometry,
			c.Digest,
		)
	case LensChronos:
		encName := EncodingName(c.Encoding)
		return fmt.Sprintf("/by-date/epoch_%d/auth_%08x/%s_%016x.tensor",
			c.Chronos,
			c.Authority,
			encName,
			c.Digest,
		)
	case LensFlat:
		fallthrough
	default:
		encName := EncodingName(c.Encoding)
		return fmt.Sprintf("/auth_%08x/%s_%016x.tensor",
			c.Authority,
			encName,
			c.Digest,
		)
	}
}

// Lift computes the reverse lifting s: T_1D -> O^8 for a given lens.
// Parses 1D virtual path strings back into authentic 64-byte coordinates.
func (d *FibrationDoor) Lift(virtualPath string, lens LensKind) (Coordinate64, error) {
	clean := strings.Trim(virtualPath, "/")
	parts := strings.Split(clean, "/")

	var c Coordinate64

	switch lens {
	case LensSemantic:
		// Format: by-model/{archetype}/{model_id}/epoch_{n}/{enc}_{digest}.tensor
		if len(parts) != 5 || parts[0] != "by-model" {
			return c, fmt.Errorf("invalid semantic path format: %s", virtualPath)
		}
		c.Archetype = parseArchetype(parts[1])
		if strings.HasPrefix(parts[2], "model_") {
			v, err := strconv.ParseUint(strings.TrimPrefix(parts[2], "model_"), 16, 64)
			if err == nil {
				c.Authority = v
			}
		}
		if strings.HasPrefix(parts[3], "epoch_") {
			v, err := strconv.ParseUint(strings.TrimPrefix(parts[3], "epoch_"), 10, 64)
			if err == nil {
				c.Chronos = v
			}
		}
		c.Encoding, c.Digest = parseFileTerminal(parts[4])
		return c, nil

	case LensHardware:
		// Format: by-hardware/node_{topos}/{encoding}/rank_{geom}_{digest}.tensor
		if len(parts) != 4 || parts[0] != "by-hardware" {
			return c, fmt.Errorf("invalid hardware path format: %s", virtualPath)
		}
		if strings.HasPrefix(parts[1], "node_") {
			v, err := strconv.ParseUint(strings.TrimPrefix(parts[1], "node_"), 10, 64)
			if err == nil {
				c.Topos = v
			}
		}
		c.Encoding = parseEncoding(parts[2])
		fileTerm := parts[3]
		fileClean := strings.TrimSuffix(fileTerm, ".tensor")
		if strings.HasPrefix(fileClean, "rank_") {
			rankParts := strings.Split(strings.TrimPrefix(fileClean, "rank_"), "_")
			if len(rankParts) >= 2 {
				g, _ := strconv.ParseUint(rankParts[0], 10, 64)
				c.Geometry = g
				dVal, _ := strconv.ParseUint(rankParts[1], 16, 64)
				c.Digest = dVal
			}
		}
		return c, nil

	case LensChronos:
		// Format: by-date/epoch_{chronos}/auth_{authority}/{enc}_{digest}.tensor
		if len(parts) != 4 || parts[0] != "by-date" {
			return c, fmt.Errorf("invalid chronos path format: %s", virtualPath)
		}
		if strings.HasPrefix(parts[1], "epoch_") {
			v, err := strconv.ParseUint(strings.TrimPrefix(parts[1], "epoch_"), 10, 64)
			if err == nil {
				c.Chronos = v
			}
		}
		if strings.HasPrefix(parts[2], "auth_") {
			v, err := strconv.ParseUint(strings.TrimPrefix(parts[2], "auth_"), 16, 64)
			if err == nil {
				c.Authority = v
			}
		}
		c.Encoding, c.Digest = parseFileTerminal(parts[3])
		return c, nil

	case LensFlat:
		// Format: auth_{authority}/{enc}_{digest}.tensor
		if len(parts) != 2 || !strings.HasPrefix(parts[0], "auth_") {
			return c, fmt.Errorf("invalid flat path format: %s", virtualPath)
		}
		v, err := strconv.ParseUint(strings.TrimPrefix(parts[0], "auth_"), 16, 64)
		if err == nil {
			c.Authority = v
		}
		c.Encoding, c.Digest = parseFileTerminal(parts[1])
		return c, nil

	default:
		return c, errors.New("unsupported lens kind")
	}
}

// ListVirtualDirectory enumerates child virtual directories and files under a prefix.
// Completely zero disk inodes; runs in microseconds by evaluating SoA column projections.
func (d *FibrationDoor) ListVirtualDirectory(virtualDir string, lens LensKind) ([]VirtualDirEntry, error) {
	if d.Atlas == nil {
		return nil, errors.New("fibration door has uninitialized atlas")
	}

	normDir := "/" + strings.Trim(virtualDir, "/")
	if normDir == "/" {
		normDir = ""
	}

	n := d.Atlas.Count
	seenDirs := make(map[string]bool)
	var entries []VirtualDirEntry

	for i := uint32(0); i < n; i++ {
		c, payloadRef, err := d.Atlas.Get(i)
		if err != nil {
			continue
		}

		projPath := d.Project(c, lens)
		if normDir != "" && !strings.HasPrefix(projPath, normDir+"/") {
			continue
		}

		// Strip prefix
		sub := strings.TrimPrefix(projPath, normDir+"/")
		parts := strings.Split(sub, "/")

		if len(parts) == 1 {
			// Terminal file entry
			entries = append(entries, VirtualDirEntry{
				Name:         parts[0],
				IsDir:        false,
				Size:         4096, // default page quantum
				Coordinate:   c,
				PayloadRef:   payloadRef,
				White3Digest: c.Digest,
			})
		} else {
			// Subdirectory entry
			dirName := parts[0]
			if !seenDirs[dirName] {
				seenDirs[dirName] = true
				entries = append(entries, VirtualDirEntry{
					Name:  dirName,
					IsDir: true,
				})
			}
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].IsDir != entries[j].IsDir {
			return entries[i].IsDir // directories first
		}
		return entries[i].Name < entries[j].Name
	})

	return entries, nil
}

// ResolvePath locates a virtual path in the atlas and returns its coordinate and payload reference.
func (d *FibrationDoor) ResolvePath(virtualPath string, lens LensKind) (Coordinate64, uint64, error) {
	norm := "/" + strings.Trim(virtualPath, "/")
	n := d.Atlas.Count

	for i := uint32(0); i < n; i++ {
		c, payloadRef, err := d.Atlas.Get(i)
		if err != nil {
			continue
		}
		if d.Project(c, lens) == norm {
			return c, payloadRef, nil
		}
	}

	return Coordinate64{}, 0, fmt.Errorf("virtual path not found in atlas: %s", virtualPath)
}

func parseArchetype(name string) uint64 {
	switch name {
	case "tensors":
		return ArchetypeTensorWeight
	case "checkpoints":
		return ArchetypeModelCheckpoint
	case "logs":
		return ArchetypeExecutionLog
	case "kernels":
		return ArchetypeCompiledKernel
	case "topology":
		return ArchetypeTopologyGraph
	case "datasets":
		return ArchetypeDatasetShard
	case "notarizations":
		return ArchetypeCryptoNotarization
	default:
		return 1
	}
}

func parseEncoding(name string) uint64 {
	switch name {
	case "fp8":
		return EncodingFP8
	case "bf16":
		return EncodingBF16
	case "fp16":
		return EncodingFP16
	case "fp32":
		return EncodingFP32
	case "fp64":
		return EncodingFP64
	case "int4":
		return EncodingINT4
	case "m31":
		return EncodingM31
	default:
		return 1
	}
}

func parseFileTerminal(terminal string) (uint64, uint64) {
	clean := strings.TrimSuffix(terminal, ".tensor")
	parts := strings.Split(clean, "_")
	if len(parts) >= 2 {
		enc := parseEncoding(parts[0])
		digest, _ := strconv.ParseUint(parts[1], 16, 64)
		return enc, digest
	}
	return 1, 0
}
