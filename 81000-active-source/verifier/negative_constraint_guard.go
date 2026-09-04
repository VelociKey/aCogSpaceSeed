package verifier

import (
	"fmt"
	"strings"

	"sov.fleet/extendgo/81000-active-source/spec_ingest"
)

// VerificationReport holds the audit results of the synthesized grammar.
type VerificationReport struct {
	Valid             bool
	CheckedRulesCount int
	Violations        []string
}

// VerifyGrammarConstraints audits the rule set against Go language negative constraints
// and verifies that all sovereign scalar extensions are present.
func VerifyGrammarConstraints(rules []spec_ingest.RawRule) *VerificationReport {
	report := &VerificationReport{
		Valid:             true,
		CheckedRulesCount: len(rules),
	}

	ruleMap := make(map[string]string)
	for _, r := range rules {
		ruleMap[r.Name] = r.Body
	}

	// 1. Negative Constraint Guard: InterfaceMethod must NOT have TypeParamList
	if imBody, ok := ruleMap["InterfaceMethod"]; ok {
		if strings.Contains(imBody, "TypeParamList") || strings.Contains(imBody, "[ TypeParamList ]") {
			report.Valid = false
			report.Violations = append(report.Violations, "Negative constraint violation: InterfaceMethod must not permit type parameter clauses [ TypeParamList ]")
		}
	} else {
		report.Valid = false
		report.Violations = append(report.Violations, "Missing critical production: InterfaceMethod")
	}

	// 2. Sovereign Scalar Precision Types Audit
	requiredScalars := []string{
		`"float2048"`, `"float1024"`, `"float512"`, `"float256"`, `"float128"`,
		`"complex4096"`, `"complex2048"`, `"complex1024"`, `"complex512"`, `"complex256"`,
		`"int2048"`, `"int1024"`, `"int512"`, `"int256"`, `"int128"`,
		`"uint2048"`, `"uint1024"`, `"uint512"`, `"uint256"`, `"uint128"`,
		`DecimalType`,
	}

	if btBody, ok := ruleMap["BaseType"]; ok {
		for _, req := range requiredScalars {
			if !strings.Contains(btBody, req) {
				report.Valid = false
				report.Violations = append(report.Violations, fmt.Sprintf("Missing required sovereign scalar in BaseType: %s", req))
			}
		}
	} else {
		report.Valid = false
		report.Violations = append(report.Violations, "Missing critical production: BaseType")
	}

	// 3. Go 1.27 Generic Method Receiver Check
	if fdBody, ok := ruleMap["FuncDecl"]; ok {
		if !strings.Contains(fdBody, "Receiver") || !strings.Contains(fdBody, "TypeParamList") {
			report.Valid = false
			report.Violations = append(report.Violations, "FuncDecl must support both [ Receiver ] and [ TypeParamList ] for Go 1.27 generic methods")
		}
	}

	// 4. Go 1.27 Chained Selector Check
	if selBody, ok := ruleMap["Selector"]; ok {
		if !strings.Contains(selBody, "TypeArgs") {
			report.Valid = false
			report.Violations = append(report.Violations, "Selector must support [ TypeArgs ] for Go 1.27 generic method invocations")
		}
	}

	return report
}
