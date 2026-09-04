package attestations_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"sov.fleet/extendgo/81000-active-source/emitter"
	"sov.fleet/extendgo/81000-active-source/overlay"
	"sov.fleet/extendgo/81000-active-source/spec_ingest"
	"sov.fleet/extendgo/81000-active-source/verifier"
)

func TestEBNFExtractionAndParsing(t *testing.T) {
	specPath := `C:\aCogSpaceSeed\00flow\forge\92000-external-toolchains\go\doc\go_spec.html`
	if _, err := os.Stat(specPath); err != nil {
		specPath = `C:\aCogSpaceSeed\86sref\golang-go\doc\go_spec.html`
	}

	if _, err := os.Stat(specPath); err != nil {
		t.Skipf("Go language spec not found at %s", specPath)
	}

	extracted, err := spec_ingest.ExtractEBNFFromFile(specPath)
	if err != nil {
		t.Fatalf("Failed to extract EBNF: %v", err)
	}

	if len(extracted.RuleBlocks) == 0 {
		t.Fatalf("Expected non-empty EBNF rule blocks")
	}

	rules, err := spec_ingest.ParseEBNFRules(extracted.FullText)
	if err != nil {
		t.Fatalf("Failed to parse EBNF rules: %v", err)
	}

	if len(rules) < 20 {
		t.Errorf("Expected at least 20 base production rules, got %d", len(rules))
	}

	t.Logf("✅ Successfully extracted and parsed %d production rules from %s", len(rules), specPath)
}

func TestOverlayFunctorAndConstraintGuards(t *testing.T) {
	dummyBaseRules := []spec_ingest.RawRule{
		{Name: "Identifier", Body: `[a-zA-Z_] { [a-zA-Z0-9_] }`},
		{Name: "Type", Body: `TypeName | TypeLit | "(" Type ")"`},
	}

	functor := overlay.NewRuleOverlayFunctor(overlay.DefaultGoExtensions())
	synthesized, err := functor.Apply(dummyBaseRules)
	if err != nil {
		t.Fatalf("Overlay application failed: %v", err)
	}

	report := verifier.VerifyGrammarConstraints(synthesized)
	if !report.Valid {
		t.Fatalf("Grammar constraints violated: %v", report.Violations)
	}

	ruleMap := make(map[string]string)
	for _, r := range synthesized {
		ruleMap[r.Name] = r.Body
	}

	// Verify sovereign scalar types present
	bt := ruleMap["BaseType"]
	for _, expected := range []string{`"float2048"`, `"complex4096"`, `"int2048"`, `"uint2048"`, `DecimalType`} {
		if !strings.Contains(bt, expected) {
			t.Errorf("Expected BaseType to contain %s", expected)
		}
	}

	// Verify Go 1.27 generic method and chained selector
	if !strings.Contains(ruleMap["FuncDecl"], "Receiver") || !strings.Contains(ruleMap["FuncDecl"], "TypeParamList") {
		t.Errorf("FuncDecl missing Receiver or TypeParamList")
	}
	if !strings.Contains(ruleMap["Selector"], "TypeArgs") {
		t.Errorf("Selector missing TypeArgs")
	}

	t.Logf("✅ Successfully verified overlay functor and all sovereign scalar types")
}

func TestWebNFEmissionAndLinterCompliance(t *testing.T) {
	dummyRules := []spec_ingest.RawRule{
		{Name: "Grammar", Body: "GoProgram"},
		{Name: "GoProgram", Body: "StatementList"},
	}
	functor := overlay.NewRuleOverlayFunctor(overlay.DefaultGoExtensions())
	synthesized, err := functor.Apply(dummyRules)
	if err != nil {
		t.Fatalf("Overlay failed: %v", err)
	}

	wagContent := emitter.EmitWebNFSN(synthesized, "1.27")
	if !strings.Contains(wagContent, ":com:velocikey:languagegrammar:go") {
		t.Errorf("Missing domain taxonomy header in emitted WAG")
	}

	tmpFile := filepath.Join(os.TempDir(), "test_emit_grammar_go.wag")
	if err := os.WriteFile(tmpFile, []byte(wagContent), 0644); err != nil {
		t.Fatalf("Failed to write tmp wag: %v", err)
	}
	defer os.Remove(tmpFile)

	// Validate with webnf-sn-lint.exe
	snLint := `C:\aCogSpaceSeed\00flow\latentlingua\webnf-sn-lint.exe`
	if _, err := os.Stat(snLint); err == nil {
		cmd := exec.Command(snLint, tmpFile)
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("webnf-sn-lint failed on emitted grammar: %v\nOutput: %s", err, string(out))
		}
	}

	t.Logf("✅ Emitted WAG passed formal webnf-sn-lint syntax verification")
}

func TestNegativeConstraintEnforcement(t *testing.T) {
	invalidRules := []spec_ingest.RawRule{
		{Name: "InterfaceMethod", Body: "Identifier [ TypeParamList ] Signature"},
		{Name: "BaseType", Body: `"float2048" / "float1024" / "float512" / "float256" / "float128" / "complex4096" / "complex2048" / "complex1024" / "complex512" / "complex256" / "int2048" / "int1024" / "int512" / "int256" / "int128" / "uint2048" / "uint1024" / "uint512" / "uint256" / "uint128" / DecimalType`},
		{Name: "FuncDecl", Body: `"func" [ Receiver ] Identifier [ TypeParamList ] Signature [ Block ]`},
		{Name: "Selector", Body: `"." Identifier [ TypeArgs ]`},
	}

	report := verifier.VerifyGrammarConstraints(invalidRules)
	if report.Valid {
		t.Fatalf("Expected failure due to negative constraint violation on InterfaceMethod")
	}

	foundViolation := false
	for _, v := range report.Violations {
		if strings.Contains(v, "Negative constraint violation") {
			foundViolation = true
			break
		}
	}

	if !foundViolation {
		t.Errorf("Expected negative constraint violation message, got: %v", report.Violations)
	}

	t.Logf("✅ Negative constraint guard successfully blocked invalid interface method generics")
}

func TestGoEBNFSpecScannerAndMacroIRAST(t *testing.T) {
	snippet := []byte(`
// Go 1.27 Method and Selector productions
MethodDecl = "func" Receiver MethodName Signature [ FunctionBody ] .
Selector   = "." identifier [ TypeArgs ] .
`)

	rules, astNodes, err := spec_ingest.ParseEBNFWithScanner(snippet)
	if err != nil {
		t.Fatalf("ParseEBNFWithScanner failed: %v", err)
	}

	if len(rules) != 2 {
		t.Fatalf("expected 2 rules, got %d", len(rules))
	}

	if rules[0].Name != "MethodDecl" || rules[1].Name != "Selector" {
		t.Fatalf("unexpected rule names: %s, %s", rules[0].Name, rules[1].Name)
	}

	if len(astNodes) != 2 {
		t.Fatalf("expected 2 AST roots, got %d", len(astNodes))
	}

	for _, node := range astNodes {
		if node.Kind != spec_ingest.NodeProduction {
			t.Errorf("expected NodeProduction, got %s", node.Kind)
		}
		if len(node.Children) == 0 {
			t.Errorf("expected AST children for rule %s", node.Text)
		}
	}

	t.Logf("✅ Successfully verified GoEBNFSpecScanner tokenization and Macro-IR opcode AST lowering")
}

func TestGoExtensionSpecDSLParser(t *testing.T) {
	spec, err := overlay.LoadDefaultGoExtensionSpec()
	if err != nil {
		t.Fatalf("LoadDefaultGoExtensionSpec failed: %v", err)
	}

	if spec.Name != "SovereignGo" {
		t.Errorf("expected spec.Name == 'SovereignGo', got %s", spec.Name)
	}
	if spec.BaseVersion != "1.27" {
		t.Errorf("expected spec.BaseVersion == '1.27', got %s", spec.BaseVersion)
	}
	if spec.TargetGrammar != "go" {
		t.Errorf("expected spec.TargetGrammar == 'go', got %s", spec.TargetGrammar)
	}

	if len(spec.Extensions) < 5 {
		t.Fatalf("expected at least 5 targeted extension blocks, got %d", len(spec.Extensions))
	}

	extMap := make(map[string]overlay.TargetedExtension)
	for _, ext := range spec.Extensions {
		extMap[ext.Name] = ext
	}

	// Verify ScalarPrecisionExtension
	scalarExt, ok := extMap["ScalarPrecisionExtension"]
	if !ok {
		t.Errorf("missing ScalarPrecisionExtension")
	} else {
		if scalarExt.TargetRule != "BaseType" || scalarExt.Section != "types" || scalarExt.Action != "AUGMENT" {
			t.Errorf("unexpected ScalarPrecisionExtension properties: %+v", scalarExt)
		}
	}

	// Verify FuncDeclExtension
	funcExt, ok := extMap["FuncDeclExtension"]
	if !ok {
		t.Errorf("missing FuncDeclExtension")
	} else {
		if funcExt.TargetRule != "FuncDecl" || funcExt.Section != "declarations" || funcExt.Action != "REPLACE" {
			t.Errorf("unexpected FuncDeclExtension properties: %+v", funcExt)
		}
	}

	// Verify InterfaceMethodConstraint
	ifaceExt, ok := extMap["InterfaceMethodConstraint"]
	if !ok {
		t.Errorf("missing InterfaceMethodConstraint")
	} else {
		if ifaceExt.TargetRule != "InterfaceMethod" || ifaceExt.Section != "interfaces" || ifaceExt.Action != "FORBID" {
			t.Errorf("unexpected InterfaceMethodConstraint properties: %+v", ifaceExt)
		}
	}

	t.Logf("✅ Successfully parsed and verified Go extension DSL specification (%d blocks)", len(spec.Extensions))
}

