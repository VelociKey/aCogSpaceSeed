package overlay

import (
	"fmt"
	"regexp"
	"strings"

	"sov.fleet/extendgo/81000-active-source/spec_ingest"
)

const (
	LifecycleActive     = "ACTIVE"
	LifecycleSubsumed   = "SUBSUMED"
	LifecycleDeprecated = "DEPRECATED"
)

// SubsumedFeature records details about a sovereign extension feature that has been
// natively adopted and subsumed by upstream Go, rendering it no longer needed.
type SubsumedFeature struct {
	ExtensionName string // Identifier of the extension declaration
	Feature       string // Semantic feature token (e.g., GENERIC_RECEIVER_METHODS)
	TargetRule    string // Rule in Go specification targeted by the extension
	SubsumedIf    string // Predicate tested against upstream Go spec
	SubsumedIn    string // Go version where native adoption occurred
	Reason        string // Human-readable explanation of subsumption
	UpstreamMatch string // Upstream rule snippet confirming convergence
}

// SubsumptionReport summarizes the evaluation of general sovereign extensions against
// an upstream Go version specification.
type SubsumptionReport struct {
	GoVersion        string
	TotalExtensions  int
	ActiveCount      int
	SubsumedCount    int
	ActiveExtensions []TargetedExtension
	SubsumedFeatures []SubsumedFeature
}

var versionRegex = regexp.MustCompile(`(?i)(?:Language version go|version go|go1\.)([0-9]+\.[0-9]+)`)

// DetectGoVersionFromSpec extracts the Go version string from specification content.
func DetectGoVersionFromSpec(content string, fallback string) string {
	matches := versionRegex.FindStringSubmatch(content)
	if len(matches) > 1 {
		return matches[1]
	}
	if fallback != "" {
		return fallback
	}
	return "1.27"
}

// AnalyzeSubsumption tests general sovereign extensions against the upstream Go rules
// for a target Go version, categorizing features into ACTIVE vs SUBSUMED.
func AnalyzeSubsumption(spec *ExtensionSpec, baseRules []spec_ingest.RawRule, targetVersion string) *SubsumptionReport {
	if targetVersion == "" {
		targetVersion = "1.27"
	}

	report := &SubsumptionReport{
		GoVersion:       targetVersion,
		TotalExtensions: len(spec.Extensions),
	}

	// Index upstream rules by name for fast predicate lookup
	ruleMap := make(map[string]string)
	for _, r := range baseRules {
		ruleMap[r.Name] = r.Body
	}

	for _, ext := range spec.Extensions {
		// Keyword registrations
		if ext.Action == "REGISTER_KEYWORD" {
			// Check if keyword is already standard in upstream Go
			kwSubsumed := false
			for _, kw := range ext.Keywords {
				kwPattern := `"` + kw + `"`
				if kwRule, ok := ruleMap["Keyword"]; ok && strings.Contains(kwRule, kwPattern) {
					kwSubsumed = true
					report.SubsumedFeatures = append(report.SubsumedFeatures, SubsumedFeature{
						ExtensionName: ext.Name,
						Feature:       "KEYWORD_" + strings.ToUpper(kw),
						TargetRule:    "Keyword",
						SubsumedIf:    kwPattern,
						SubsumedIn:    targetVersion,
						Reason:        fmt.Sprintf("Keyword %q natively reserved in Go %s", kw, targetVersion),
						UpstreamMatch: kwRule,
					})
				}
			}
			if !kwSubsumed {
				report.ActiveExtensions = append(report.ActiveExtensions, ext)
			}
			continue
		}

		// Evaluate subsumption predicate if present
		isSubsumed := false
		if ext.SubsumedIf != "" {
			// 1. Direct rule match
			if targetBody, ok := ruleMap[ext.TargetRule]; ok && strings.Contains(targetBody, ext.SubsumedIf) {
				isSubsumed = true
				subVer := ext.SubsumedIn
				if subVer == "" {
					subVer = targetVersion
				}
				report.SubsumedFeatures = append(report.SubsumedFeatures, SubsumedFeature{
					ExtensionName: ext.Name,
					Feature:       ext.Feature,
					TargetRule:    ext.TargetRule,
					SubsumedIf:    ext.SubsumedIf,
					SubsumedIn:    subVer,
					Reason:        fmt.Sprintf("Upstream Go %s rule %s matches subsumption predicate %q", targetVersion, ext.TargetRule, ext.SubsumedIf),
					UpstreamMatch: targetBody,
				})
			} else {
				// 2. Related rule matches (e.g. FuncDecl vs FunctionDecl / MethodDecl)
				relatedRules := []string{"MethodDecl", "FunctionDecl", "FuncDecl", "TypeDecl", "TypeDef", "Selector"}
				for _, rName := range relatedRules {
					if rBody, ok := ruleMap[rName]; ok && strings.Contains(rBody, ext.SubsumedIf) {
						isSubsumed = true
						subVer := ext.SubsumedIn
						if subVer == "" {
							subVer = targetVersion
						}
						report.SubsumedFeatures = append(report.SubsumedFeatures, SubsumedFeature{
							ExtensionName: ext.Name,
							Feature:       ext.Feature,
							TargetRule:    ext.TargetRule,
							SubsumedIf:    ext.SubsumedIf,
							SubsumedIn:    subVer,
							Reason:        fmt.Sprintf("Upstream Go %s rule %s matches subsumption predicate %q (subsuming %s)", targetVersion, rName, ext.SubsumedIf, ext.Name),
							UpstreamMatch: rBody,
						})
						break
					}
				}
			}
		}

		if isSubsumed {
			continue
		}

		// Check explicit lifecycle tag
		if ext.Lifecycle == LifecycleSubsumed {
			subVer := ext.SubsumedIn
			if subVer == "" {
				subVer = targetVersion
			}
			report.SubsumedFeatures = append(report.SubsumedFeatures, SubsumedFeature{
				ExtensionName: ext.Name,
				Feature:       ext.Feature,
				TargetRule:    ext.TargetRule,
				SubsumedIf:    ext.SubsumedIf,
				SubsumedIn:    subVer,
				Reason:        fmt.Sprintf("Feature explicitly tagged as SUBSUMED as of Go %s", subVer),
			})
			continue
		}

		// Feature remains active
		report.ActiveExtensions = append(report.ActiveExtensions, ext)
	}

	report.ActiveCount = len(report.ActiveExtensions)
	report.SubsumedCount = len(report.SubsumedFeatures)
	return report
}

// EmitVersionDeltaWebNF serializes the realized sovereign Go extensions for a specific
// Go version into a clean, grammar-conformant .webnf program document.
func EmitVersionDeltaWebNF(spec *ExtensionSpec, report *SubsumptionReport) string {
	var sb strings.Builder

	sb.WriteString("# ==============================================================================\n")
	sb.WriteString(fmt.Sprintf("# REALIZED SOVEREIGN GO EXTENSIONS SPECIFICATION (extend_go_v%s.webnf)\n", report.GoVersion))
	sb.WriteString(fmt.Sprintf("# Upstream Target: Go v%s\n", report.GoVersion))
	sb.WriteString("# Grammar Conformance: go_extensions.wag (grammar_go_extension_spec.wag)\n")
	sb.WriteString("# Domain: Sovereign Systems, Supercomputing & High-Precision Science\n")
	sb.WriteString(fmt.Sprintf("# Active Extensions: %d | Subsumed Upstream Features: %d\n", report.ActiveCount, report.SubsumedCount))
	sb.WriteString("# ==============================================================================\n\n")

	sb.WriteString(fmt.Sprintf(":com:velocikey:extendgo:go_extensions_v%s\n\n", strings.ReplaceAll(report.GoVersion, ".", "_")))

	// Metadata Block
	sb.WriteString("# --- Realized Extension Metadata ---\n")
	sb.WriteString(fmt.Sprintf("ExtensionMetadata = \"name\" \":\" \"ExtendGo_v%s\"\n", report.GoVersion))
	sb.WriteString("                  / \"target_language\" \":\" \"go\"\n")
	sb.WriteString(fmt.Sprintf("                  / \"base_version\" \":\" %q\n", report.GoVersion))
	sb.WriteString(fmt.Sprintf("                  / \"description\" \":\" \"Realized sovereign extensions active for Go v%s\" ;\n\n", report.GoVersion))

	// Active Extensions
	sb.WriteString("# ==============================================================================\n")
	sb.WriteString(fmt.Sprintf("# 1. ACTIVE SOVEREIGN EXTENSIONS (Required for Go v%s)\n", report.GoVersion))
	sb.WriteString("# ==============================================================================\n\n")

	for idx, ext := range report.ActiveExtensions {
		sb.WriteString(fmt.Sprintf("# --- [%d/%d] Active Extension: %s ---\n", idx+1, report.ActiveCount, ext.Name))
		if ext.Action == "REGISTER_KEYWORD" {
			sb.WriteString(fmt.Sprintf("%s = \"action\" \":\" \"REGISTER_KEYWORD\"\n", ext.Name))
			for _, kw := range ext.Keywords {
				sb.WriteString(fmt.Sprintf("              / %q ;\n", kw))
			}
			sb.WriteString("\n")
			continue
		}

		sb.WriteString(fmt.Sprintf("%s = \"target_rule\" \":\" %q\n", ext.Name, ext.TargetRule))
		sb.WriteString(fmt.Sprintf("                 / \"section\" \":\" %q\n", ext.Section))
		sb.WriteString(fmt.Sprintf("                 / \"action\" \":\" %q\n", ext.Action))
		if ext.Feature != "" {
			sb.WriteString(fmt.Sprintf("                 / \"feature\" \":\" %q\n", ext.Feature))
		}
		sb.WriteString("                 / \"lifecycle\" \":\" \"ACTIVE\"\n")
		if ext.SubsumedIf != "" {
			sb.WriteString(fmt.Sprintf("                 / \"subsumed_if\" \":\" %q\n", ext.SubsumedIf))
		}
		for _, f := range ext.ForbidElements {
			sb.WriteString(fmt.Sprintf("                 / \"forbid\" \":\" %q\n", f))
		}
		if ext.ProductionBody != "" {
			sb.WriteString(fmt.Sprintf("                 / ( %s ) ;\n\n", ext.ProductionBody))
		} else {
			sb.WriteString("                 ;\n\n")
		}
	}

	// Subsumed Upstream Features Audit Trail
	sb.WriteString("# ==============================================================================\n")
	sb.WriteString(fmt.Sprintf("# 2. SUBSUMED UPSTREAM FEATURES (NO LONGER NEEDED as of Go v%s)\n", report.GoVersion))
	sb.WriteString("# ==============================================================================\n")

	if len(report.SubsumedFeatures) == 0 {
		sb.WriteString("# None. All sovereign extensions remain active.\n")
	} else {
		for idx, sf := range report.SubsumedFeatures {
			sb.WriteString(fmt.Sprintf("# [%d] Feature: %s\n", idx+1, sf.Feature))
			sb.WriteString(fmt.Sprintf("#     Extension Decl: %s\n", sf.ExtensionName))
			sb.WriteString(fmt.Sprintf("#     Target Rule   : %s\n", sf.TargetRule))
			sb.WriteString(fmt.Sprintf("#     Subsumed In   : Go v%s (no longer needed as of Go v%s)\n", sf.SubsumedIn, sf.SubsumedIn))
			sb.WriteString(fmt.Sprintf("#     Reason        : %s\n", sf.Reason))
			if sf.UpstreamMatch != "" {
				sb.WriteString(fmt.Sprintf("#     Upstream Match: %s\n", strings.TrimSpace(sf.UpstreamMatch)))
			}
			sb.WriteString("#\n")
		}
	}

	return sb.String()
}
