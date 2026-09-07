package atlas

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

// PolyFacetedTag represents a single key-value dimension in a meet-semilattice.
type PolyFacetedTag struct {
	Key string
	Val string
}

// TagBag represents an order-free conjunction of semantic facets.
type TagBag []PolyFacetedTag

// Normalize canonicalizes the tag bag by sorting and deduplicating.
// Guarantees mathematical commutativity: A ^ B == B ^ A.
func (tb TagBag) Normalize() TagBag {
	if len(tb) == 0 {
		return tb
	}
	res := make(TagBag, len(tb))
	copy(res, tb)
	sort.Slice(res, func(i, j int) bool {
		if res[i].Key == res[j].Key {
			return res[i].Val < res[j].Val
		}
		return res[i].Key < res[j].Key
	})

	// Deduplicate
	dedup := make(TagBag, 0, len(res))
	for i, t := range res {
		if i == 0 || t.Key != res[i-1].Key || t.Val != res[i-1].Val {
			dedup = append(dedup, t)
		}
	}
	return dedup
}

// LatticeMeet computes the meet (conjunction / specialization) of two tag bags: A ∧ B.
// Returns the union of distinct facets from both bags. If a key conflicts, the more specific tag is retained.
func LatticeMeet(a, b TagBag) TagBag {
	merged := make(TagBag, 0, len(a)+len(b))
	merged = append(merged, a...)
	merged = append(merged, b...)
	return merged.Normalize()
}

// LatticeJoin computes the join (disjunction / generalization) of two tag bags: A ∨ B.
// Returns only the shared facets present in both bags.
func LatticeJoin(a, b TagBag) TagBag {
	na := a.Normalize()
	nb := b.Normalize()
	var common TagBag
	for _, ta := range na {
		for _, tb := range nb {
			if ta.Key == tb.Key && ta.Val == tb.Val {
				common = append(common, ta)
				break
			}
		}
	}
	return common
}

// Equals checks if two tag bags represent identical semantic points in the lattice,
// regardless of input ordering.
func (tb TagBag) Equals(other TagBag) bool {
	na := tb.Normalize()
	nb := other.Normalize()
	if len(na) != len(nb) {
		return false
	}
	for i := range na {
		if na[i].Key != nb[i].Key || na[i].Val != nb[i].Val {
			return false
		}
	}
	return true
}

// ArchetypeName returns the human-readable string for archetype enum.
func ArchetypeName(arch uint64) string {
	switch arch {
	case ArchetypeTensorWeight:
		return "tensors"
	case ArchetypeModelCheckpoint:
		return "checkpoints"
	case ArchetypeExecutionLog:
		return "logs"
	case ArchetypeCompiledKernel:
		return "kernels"
	case ArchetypeTopologyGraph:
		return "topology"
	case ArchetypeDatasetShard:
		return "datasets"
	case ArchetypeCryptoNotarization:
		return "notarizations"
	default:
		return fmt.Sprintf("arch_%d", arch)
	}
}

// EncodingName returns the human-readable string for encoding enum.
func EncodingName(enc uint64) string {
	switch enc {
	case EncodingFP8:
		return "fp8"
	case EncodingBF16:
		return "bf16"
	case EncodingFP16:
		return "fp16"
	case EncodingFP32:
		return "fp32"
	case EncodingFP64:
		return "fp64"
	case EncodingINT4:
		return "int4"
	case EncodingM31:
		return "m31"
	default:
		return fmt.Sprintf("enc_%d", enc)
	}
}

// CoordinateToPOSIX projects a 64-byte coordinate into an ephemeral virtual POSIX path.
// No directory nodes or inodes are allocated or stored on physical media.
func CoordinateToPOSIX(c Coordinate64) string {
	return fmt.Sprintf("/auth_%08x/%s/epoch_%d/node_%d/%s_%016x.tensor",
		c.Authority,
		ArchetypeName(c.Archetype),
		c.Chronos,
		c.Topos,
		EncodingName(c.Encoding),
		c.Digest,
	)
}

// ParseCanonicalTagSet parses a comma-separated attribute bag into a Coordinate64.
// Example: "auth:0x1234,arch:tensor,epoch:42,node:5,enc:fp8,digest:0xabcd"
func ParseCanonicalTagSet(tagStr string) (Coordinate64, error) {
	var c Coordinate64
	pairs := strings.Split(tagStr, ",")
	for _, pair := range pairs {
		kv := strings.SplitN(strings.TrimSpace(pair), ":", 2)
		if len(kv) != 2 {
			continue
		}
		k := strings.ToLower(strings.TrimSpace(kv[0]))
		v := strings.TrimSpace(kv[1])

		switch k {
		case "auth", "authority":
			val, err := parseHexOrDec(v)
			if err != nil {
				return c, err
			}
			c.Authority = val
		case "arch", "archetype":
			switch strings.ToLower(v) {
			case "tensor", "tensorweight":
				c.Archetype = ArchetypeTensorWeight
			case "checkpoint", "model":
				c.Archetype = ArchetypeModelCheckpoint
			case "log":
				c.Archetype = ArchetypeExecutionLog
			case "kernel":
				c.Archetype = ArchetypeCompiledKernel
			default:
				val, err := parseHexOrDec(v)
				if err != nil {
					return c, err
				}
				c.Archetype = val
			}
		case "epoch", "chronos", "step":
			val, err := parseHexOrDec(v)
			if err != nil {
				return c, err
			}
			c.Chronos = val
		case "node", "topos", "zone":
			val, err := parseHexOrDec(v)
			if err != nil {
				return c, err
			}
			c.Topos = val
		case "geom", "geometry", "rank":
			val, err := parseHexOrDec(v)
			if err != nil {
				return c, err
			}
			c.Geometry = val
		case "enc", "encoding", "prec", "precision":
			switch strings.ToLower(v) {
			case "fp8":
				c.Encoding = EncodingFP8
			case "bf16":
				c.Encoding = EncodingBF16
			case "fp16":
				c.Encoding = EncodingFP16
			case "fp32":
				c.Encoding = EncodingFP32
			case "fp64":
				c.Encoding = EncodingFP64
			case "int4":
				c.Encoding = EncodingINT4
			case "m31":
				c.Encoding = EncodingM31
			default:
				val, err := parseHexOrDec(v)
				if err != nil {
					return c, err
				}
				c.Encoding = val
			}
		case "lineage":
			val, err := parseHexOrDec(v)
			if err != nil {
				return c, err
			}
			c.Lineage = val
		case "digest", "hash":
			val, err := parseHexOrDec(v)
			if err != nil {
				return c, err
			}
			c.Digest = val
		}
	}
	return c, nil
}

func parseHexOrDec(s string) (uint64, error) {
	clean := strings.TrimPrefix(strings.ToLower(s), "0x")
	if len(clean) < len(s) {
		return strconv.ParseUint(clean, 16, 64)
	}
	return strconv.ParseUint(s, 10, 64)
}
