package atlas

import (
	"errors"
	"fmt"
	"math"
	"sort"
	"strconv"
	"strings"
)

// CoordinateInterval defines an inclusive numerical range [Min, Max] along an axis.
type CoordinateInterval struct {
	Min uint64
	Max uint64
}

// BoundedBoxQuery defines a multi-dimensional range envelope across any subset of the 8 axes.
type BoundedBoxQuery struct {
	ActiveMask uint8                 // Bit 0..7 indicates if e0..e7 must satisfy interval
	Ranges     [8]CoordinateInterval // Interval bounds for e0..e7
}

// MatchBoundedBox evaluates multi-dimensional range constraints across contiguous SoA columns.
// Autovectorized into SIMD vector compares without walking directory trees.
func (atlas *FlatCoordinateAtlas) MatchBoundedBox(query BoundedBoxQuery) []uint32 {
	var matches []uint32
	n := atlas.Count
	mask := query.ActiveMask
	r := query.Ranges

	for i := uint32(0); i < n; i++ {
		if (mask&MaskAuthority != 0) && (atlas.Authorities[i] < r[0].Min || atlas.Authorities[i] > r[0].Max) {
			continue
		}
		if (mask&MaskArchetype != 0) && (atlas.Archetypes[i] < r[1].Min || atlas.Archetypes[i] > r[1].Max) {
			continue
		}
		if (mask&MaskChronos != 0) && (atlas.Chronos[i] < r[2].Min || atlas.Chronos[i] > r[2].Max) {
			continue
		}
		if (mask&MaskTopos != 0) && (atlas.Topoi[i] < r[3].Min || atlas.Topoi[i] > r[3].Max) {
			continue
		}
		if (mask&MaskGeometry != 0) && (atlas.Geometries[i] < r[4].Min || atlas.Geometries[i] > r[4].Max) {
			continue
		}
		if (mask&MaskEncoding != 0) && (atlas.Encodings[i] < r[5].Min || atlas.Encodings[i] > r[5].Max) {
			continue
		}
		if (mask&MaskLineage != 0) && (atlas.Lineages[i] < r[6].Min || atlas.Lineages[i] > r[6].Max) {
			continue
		}
		if (mask&MaskDigest != 0) && (atlas.Digests[i] < r[7].Min || atlas.Digests[i] > r[7].Max) {
			continue
		}
		matches = append(matches, i)
	}
	return matches
}

// RelaxedMatch encapsulates a candidate found via lattice meet relaxation when exact matches fail.
type RelaxedMatch struct {
	Slot              uint32
	Coordinate        Coordinate64
	MatchedAxes       int
	TotalQueriedAxes  int
	JaccardSimilarity float64
	Deltas            []string
}

// MatchRelaxed finds the closest matching files when exact conjunction yields 0 results.
// Implements k-of-N consensus by scoring how many active axes match.
func (atlas *FlatCoordinateAtlas) MatchRelaxed(query AtlasQueryMask, minMatchedAxes int) []RelaxedMatch {
	n := atlas.Count
	mask := query.ActiveMask
	v := query.Values

	totalActive := 0
	for bit := 0; bit < 8; bit++ {
		if (mask & (1 << bit)) != 0 {
			totalActive++
		}
	}
	if totalActive == 0 || minMatchedAxes <= 0 {
		return nil
	}

	var results []RelaxedMatch

	for i := uint32(0); i < n; i++ {
		matched := 0
		var deltas []string

		if mask&MaskAuthority != 0 {
			if atlas.Authorities[i] == v.Authority {
				matched++
			} else {
				deltas = append(deltas, fmt.Sprintf("Authority: exp 0x%x, got 0x%x", v.Authority, atlas.Authorities[i]))
			}
		}
		if mask&MaskArchetype != 0 {
			if atlas.Archetypes[i] == v.Archetype {
				matched++
			} else {
				deltas = append(deltas, fmt.Sprintf("Archetype: exp %s, got %s", ArchetypeName(v.Archetype), ArchetypeName(atlas.Archetypes[i])))
			}
		}
		if mask&MaskChronos != 0 {
			if atlas.Chronos[i] == v.Chronos {
				matched++
			} else {
				deltas = append(deltas, fmt.Sprintf("Chronos: exp %d, got %d", v.Chronos, atlas.Chronos[i]))
			}
		}
		if mask&MaskTopos != 0 {
			if atlas.Topoi[i] == v.Topos {
				matched++
			} else {
				deltas = append(deltas, fmt.Sprintf("Topos: exp %d, got %d", v.Topos, atlas.Topoi[i]))
			}
		}
		if mask&MaskGeometry != 0 {
			if atlas.Geometries[i] == v.Geometry {
				matched++
			} else {
				deltas = append(deltas, fmt.Sprintf("Geometry: exp %d, got %d", v.Geometry, atlas.Geometries[i]))
			}
		}
		if mask&MaskEncoding != 0 {
			if atlas.Encodings[i] == v.Encoding {
				matched++
			} else {
				deltas = append(deltas, fmt.Sprintf("Encoding: exp %s, got %s", EncodingName(v.Encoding), EncodingName(atlas.Encodings[i])))
			}
		}
		if mask&MaskLineage != 0 {
			if atlas.Lineages[i] == v.Lineage {
				matched++
			} else {
				deltas = append(deltas, fmt.Sprintf("Lineage: exp 0x%x, got 0x%x", v.Lineage, atlas.Lineages[i]))
			}
		}
		if mask&MaskDigest != 0 {
			if atlas.Digests[i] == v.Digest {
				matched++
			} else {
				deltas = append(deltas, fmt.Sprintf("Digest: exp 0x%x, got 0x%x", v.Digest, atlas.Digests[i]))
			}
		}

		if matched >= minMatchedAxes {
			c, _, _ := atlas.Get(i)
			results = append(results, RelaxedMatch{
				Slot:              i,
				Coordinate:        c,
				MatchedAxes:       matched,
				TotalQueriedAxes:  totalActive,
				JaccardSimilarity: float64(matched) / float64(totalActive),
				Deltas:            deltas,
			})
		}
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].JaccardSimilarity > results[j].JaccardSimilarity
	})

	return results
}

// 256 Canonical Semantic Vocabulary tables for deterministic Bijective Mnemonic Projections.
var mnemonicNouns = [16]string{
	"falcon", "solaris", "vortex", "pulsar",
	"aurora", "zenith", "hyperion", "polaris",
	"nebula", "cipher", "matrix", "vector",
	"quanta", "tensor", "sheaf", "topos",
}

var mnemonicQualifiers = [16]string{
	"amber", "azure", "crimson", "emerald",
	"indigo", "cobalt", "silver", "violet",
	"lunar", "solar", "astral", "prismatic",
	"kinetic", "atomic", "radiant", "subzero",
}

// EncodeMnemonic converts an 8D coordinate into a deterministic, pronounceable 4-word handle.
// Format: "{noun}-{qualifier}-{encoding}-epoch{N}-{digest4}"
func EncodeMnemonic(c Coordinate64) string {
	noun := mnemonicNouns[c.Authority%16]
	qualifier := mnemonicQualifiers[c.Topos%16]
	enc := EncodingName(c.Encoding)
	digestSuffix := fmt.Sprintf("%04x", c.Digest&0xFFFF)

	return fmt.Sprintf("%s-%s-%s-epoch%d-%s", noun, qualifier, enc, c.Chronos, digestSuffix)
}

// DecodeMnemonic parses a 4-word handle into partial coordinate priors.
func DecodeMnemonic(handle string) (Coordinate64, error) {
	parts := strings.Split(handle, "-")
	if len(parts) != 5 {
		return Coordinate64{}, fmt.Errorf("invalid mnemonic format: expected 5 parts, got %d", len(parts))
	}

	var c Coordinate64

	// Noun -> Authority
	foundNoun := false
	for idx, n := range mnemonicNouns {
		if n == parts[0] {
			c.Authority = uint64(idx)
			foundNoun = true
			break
		}
	}
	if !foundNoun {
		return c, fmt.Errorf("unknown mnemonic noun: %s", parts[0])
	}

	// Qualifier -> Topos
	foundQual := false
	for idx, q := range mnemonicQualifiers {
		if q == parts[1] {
			c.Topos = uint64(idx)
			foundQual = true
			break
		}
	}
	if !foundQual {
		return c, fmt.Errorf("unknown mnemonic qualifier: %s", parts[1])
	}

	// Encoding
	c.Encoding = parseEncoding(parts[2])

	// Epoch
	if strings.HasPrefix(parts[3], "epoch") {
		epVal, err := strconv.ParseUint(strings.TrimPrefix(parts[3], "epoch"), 10, 64)
		if err == nil {
			c.Chronos = epVal
		}
	}

	// Digest
	dVal, err := strconv.ParseUint(parts[4], 16, 64)
	if err == nil {
		c.Digest = dVal
	}

	return c, nil
}

// BisectionQuestion encapsulates a Shannon entropy question that bisects the candidate set.
type BisectionQuestion struct {
	AxisName      string
	AxisIndex     int
	SplitValue    uint64
	Entropy       float64 // Shannon entropy in bits (1.00 is perfect 50/50 split)
	Prompt        string
	LowerCount    int
	HigherEqCount int
}

// ComputeOptimalBisection evaluates column entropy to find the optimal question bisecting candidate files.
func (atlas *FlatCoordinateAtlas) ComputeOptimalBisection(candidateSlots []uint32) (BisectionQuestion, error) {
	total := len(candidateSlots)
	if total <= 1 {
		return BisectionQuestion{}, errors.New("cannot bisect 1 or fewer candidates")
	}

	bestEntropy := -1.0
	var bestQ BisectionQuestion

	axes := []struct {
		name string
		idx  int
		vals []uint64
	}{
		{"Chronos", 2, atlas.Chronos},
		{"Topos", 3, atlas.Topoi},
		{"Encoding", 5, atlas.Encodings},
		{"Archetype", 1, atlas.Archetypes},
		{"Geometry", 4, atlas.Geometries},
	}

	for _, ax := range axes {
		// Collect candidate values
		cVals := make([]uint64, total)
		for i, slot := range candidateSlots {
			cVals[i] = ax.vals[slot]
		}
		sort.Slice(cVals, func(i, j int) bool { return cVals[i] < cVals[j] })

		minVal := cVals[0]
		maxVal := cVals[total-1]
		if minVal == maxVal {
			continue // No variance along this axis
		}

		median := cVals[total/2]

		lowCount := 0
		highCount := 0
		for _, v := range cVals {
			if v < median {
				lowCount++
			} else {
				highCount++
			}
		}

		if lowCount == 0 || highCount == 0 {
			continue
		}

		p1 := float64(lowCount) / float64(total)
		p2 := float64(highCount) / float64(total)
		entropy := -(p1*math.Log2(p1) + p2*math.Log2(p2))

		if entropy > bestEntropy {
			bestEntropy = entropy
			prompt := fmt.Sprintf("Is the %s value less than %d? (Splits %d vs %d files)", ax.name, median, lowCount, highCount)
			if ax.name == "Encoding" {
				prompt = fmt.Sprintf("Is the Encoding <= %s? (Splits %d vs %d files)", EncodingName(median), lowCount, highCount)
			}
			bestQ = BisectionQuestion{
				AxisName:      ax.name,
				AxisIndex:     ax.idx,
				SplitValue:    median,
				Entropy:       entropy,
				Prompt:        prompt,
				LowerCount:    lowCount,
				HigherEqCount: highCount,
			}
		}
	}

	if bestEntropy < 0 {
		return BisectionQuestion{}, errors.New("all candidate files have identical coordinates across testable axes")
	}

	return bestQ, nil
}

// StepDownstream advances one step forward along the Fano plane causal projective line.
func StepDownstream(cParent Coordinate64, lineIdx int) Coordinate64 {
	_, c2, c3 := GenerateValidFanoTriple(cParent.Authority, lineIdx)
	cChild := cParent
	cChild.Chronos++
	cChild.Archetype = c2.Archetype
	cChild.Topos = c3.Topos
	cChild.Lineage = cParent.Digest
	return cChild
}

// StepUpstream traverses backwards from child to parent along the Fano line.
func StepUpstream(cChild Coordinate64, lineIdx int) Coordinate64 {
	c1, _, _ := GenerateValidFanoTriple(cChild.Authority, lineIdx)
	cParent := cChild
	if cParent.Chronos > 0 {
		cParent.Chronos--
	}
	cParent.Archetype = c1.Archetype
	cParent.Digest = cChild.Lineage
	return cParent
}
