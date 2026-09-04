package overlay

import (
	"fmt"
	"os"
	"strings"

	"sov.fleet/extendgo/81000-active-source/spec_ingest"
)

// ExtensionOverlay defines sovereign additions, modifications, and constraints to standard Go EBNF.
type ExtensionOverlay struct {
	TargetGrammar string
	TargetVersion string
	ScalarTypes   []string
	DecimalRule   string
	ExtraKeywords []string
}

// DefaultGoExtensions returns the authoritative sovereign Go language extensions specification.
func DefaultGoExtensions() *ExtensionOverlay {
	return &ExtensionOverlay{
		TargetGrammar: "go",
		TargetVersion: "1.27",
		ScalarTypes: []string{
			`"float2048"`, `"float1024"`, `"float512"`, `"float256"`, `"float128"`, `"float64"`, `"float32"`,
			`"complex4096"`, `"complex2048"`, `"complex1024"`, `"complex512"`, `"complex256"`, `"complex128"`, `"complex64"`,
			`"int2048"`, `"int1024"`, `"int512"`, `"int256"`, `"int128"`, `"int64"`, `"int32"`, `"int16"`, `"int8"`, `"int"`,
			`"uint2048"`, `"uint1024"`, `"uint512"`, `"uint256"`, `"uint128"`, `"uint64"`, `"uint32"`, `"uint16"`, `"uint8"`, `"uint"`,
			`DecimalType`,
			`GenericType`,
			`"string"`, `"bool"`, `"byte"`, `"rune"`, `"error"`, `Identifier`,
		},
		DecimalRule:   `"decimal" "[" ( Number / Identifier ) "]"`,
		ExtraKeywords: []string{`"decimal"`},
	}
}

// RuleOverlayFunctor applies categorical rule morphisms from the overlay onto base Go rules.
type RuleOverlayFunctor struct {
	Overlay               *ExtensionOverlay
	Spec                  *ExtensionSpec
	LastSubsumptionReport *SubsumptionReport
}

// NewRuleOverlayFunctor creates an instance initialized with the given overlay.
func NewRuleOverlayFunctor(ov *ExtensionOverlay) *RuleOverlayFunctor {
	if ov == nil {
		ov = DefaultGoExtensions()
	}
	functor := &RuleOverlayFunctor{Overlay: ov}
	if spec, err := LoadDefaultGoExtensionSpec(); err == nil && spec != nil {
		functor.Spec = spec
	}
	return functor
}

// LoadDefaultGoExtensionSpec discovers and parses the authoritative go_extensions.wag DSL file.
func LoadDefaultGoExtensionSpec() (*ExtensionSpec, error) {
	candidates := []string{
		`C:\aCogSpaceSeed\00flow\extendgo\81000-active-source\overlay\go_extensions.webnf`,
		`C:\aCogSpaceSeed\00flow\extendgo\81000-active-source\overlay\go_extensions.wag`,
		`81000-active-source/overlay/go_extensions.webnf`,
		`81000-active-source/overlay/go_extensions.wag`,
		`00flow/extendgo/81000-active-source/overlay/go_extensions.webnf`,
		`00flow/extendgo/81000-active-source/overlay/go_extensions.wag`,
		`../81000-active-source/overlay/go_extensions.webnf`,
		`../81000-active-source/overlay/go_extensions.wag`,
	}
	for _, c := range candidates {
		if data, err := os.ReadFile(c); err == nil {
			return ParseExtensionSpec(data)
		}
	}
	return nil, fmt.Errorf("no go_extensions.wag found in search candidates")
}

// Apply transforms the ingested rules by injecting sovereign scalar types,
// Go 1.27 generic methods, struct selectors, and negative interface constraints.
func (f *RuleOverlayFunctor) Apply(baseRules []spec_ingest.RawRule) ([]spec_ingest.RawRule, error) {
	ruleMap := make(map[string]string)
	var ruleOrder []string

	for _, r := range baseRules {
		// Convert standard EBNF alternation '|' to weBNF.sn '/'
		body := convertPipesToSlashes(r.Body)
		ruleMap[r.Name] = body
		ruleOrder = append(ruleOrder, r.Name)
	}

	// 1. Ensure core root rules exist
	f.ensureRule(ruleMap, &ruleOrder, "Grammar", "GoProgram")
	f.ensureRule(ruleMap, &ruleOrder, "GoProgram", "StatementList")
	f.ensureRule(ruleMap, &ruleOrder, "StatementList", `{ Statement ";" }`)
	f.ensureRule(ruleMap, &ruleOrder, "Statement", `DeclStmt / ImportStmt / FuncDecl / Comment`)
	f.ensureRule(ruleMap, &ruleOrder, "ImportStmt", `"import" ( StringLit / "(" { StringLit ";" } ")" )`)
	f.ensureRule(ruleMap, &ruleOrder, "DeclStmt", `"type" Identifier [ TypeParamList ] ( StructType / InterfaceType / DataType ) / "var" Identifier [ TypeParamList ] [ DataType ] [ "=" Expr ] / "const" Identifier [ DataType ] [ "=" Expr ] / Identifier [ TypeParamList ] [ ":" DataType ] [ "=" Expr ]`)

	// 2. Go 1.27 Generic Function and Method Declarations
	ruleMap["FuncDecl"] = `"func" [ Receiver ] Identifier [ TypeParamList ] Signature [ Block ]`
	ruleMap["Receiver"] = `"(" [ Identifier ] DataType ")"`
	ruleMap["Signature"] = `"(" [ ParamList ] ")" [ ResultType ]`
	ruleMap["ResultType"] = `DataType / "(" TypeList ")"`
	ruleMap["Block"] = `"{" StatementList "}"`

	// 3. Generic Type Parameters
	ruleMap["TypeParamList"] = `"[" TypeParamDecl { "," TypeParamDecl } "]"`
	ruleMap["TypeParamDecl"] = `IdentifierList TypeConstraint`
	ruleMap["TypeConstraint"] = `"any" / "comparable" / DataType`
	ruleMap["TypeList"] = `DataType { "," DataType }`
	ruleMap["ParamList"] = `Param { "," Param }`
	ruleMap["Param"] = `[ IdentifierList ] DataType`
	ruleMap["IdentifierList"] = `Identifier { "," Identifier }`

	// 4. Go 1.27 Interface Types (Negative constraint: no TypeParamList on InterfaceMethod)
	ruleMap["InterfaceType"] = `"interface" "{" { InterfaceElem ";" } "}"`
	ruleMap["InterfaceElem"] = `InterfaceMethod / TypeConstraint`
	ruleMap["InterfaceMethod"] = `Identifier Signature`

	// 5. Struct Types & Go 1.27 Field Selectors
	ruleMap["StructType"] = `"struct" "{" { FieldDecl ";" } "}"`
	ruleMap["FieldDecl"] = `IdentifierList DataType [ StringLit ] / [ "*" ] Identifier [ StringLit ]`
	ruleMap["DataType"] = `[ "*" / "[]" ] BaseType`

	// 6. Sovereign BaseType with all scientific supercomputing scalar types
	ruleMap["BaseType"] = strings.Join(f.Overlay.ScalarTypes, " / ")
	ruleMap["DecimalType"] = f.Overlay.DecimalRule
	ruleMap["GenericType"] = `Identifier TypeArgs`
	ruleMap["TypeArgs"] = `"[" TypeList "]"`

	// Apply declarative extension overrides from go_extensions.webnf DSL
	if f.Spec != nil {
		report := AnalyzeSubsumption(f.Spec, baseRules, f.Overlay.TargetVersion)
		f.LastSubsumptionReport = report

		for _, ext := range report.ActiveExtensions {
			switch ext.Action {
			case "AUGMENT":
				if ext.TargetRule == "BaseType" && ext.ProductionBody != "" {
					ruleMap["BaseType"] = ext.ProductionBody + ` / "float64" / "float32" / "complex128" / "complex64" / "int64" / "int32" / "int16" / "int8" / "int" / "uint64" / "uint32" / "uint16" / "uint8" / "uint" / GenericType / "string" / "bool" / "byte" / "rune" / "error" / Identifier`
				}
			case "INJECT":
				if ext.ProductionBody != "" {
					ruleMap[ext.TargetRule] = ext.ProductionBody
					f.ensureRule(ruleMap, &ruleOrder, ext.TargetRule, ext.ProductionBody)
				}
			case "REPLACE":
				if ext.ProductionBody != "" {
					ruleMap[ext.TargetRule] = ext.ProductionBody
				}
			case "REGISTER_KEYWORD":
				for _, kw := range ext.Keywords {
					f.Overlay.ExtraKeywords = append(f.Overlay.ExtraKeywords, fmt.Sprintf("%q", kw))
				}
			}
		}
	}

	// 7. Expressions and Go 1.27 chained selectors
	ruleMap["Expr"] = `PrimaryExpr / CompositeLit / Number / StringLit / Identifier`
	ruleMap["PrimaryExpr"] = `( Identifier / "(" Expr ")" ) { Selector / Index / Arguments }`
	ruleMap["Selector"] = `"." Identifier [ TypeArgs ]`
	ruleMap["Index"] = `"[" Expr "]"`
	ruleMap["Arguments"] = `"(" [ ArgList ] ")"`
	ruleMap["ArgList"] = `Expr { "," Expr }`

	// 8. Go 1.27 Composite literals with field selectors
	ruleMap["CompositeLit"] = `( Identifier / GenericType / StructType ) LiteralValue`
	ruleMap["LiteralValue"] = `"{" [ ElementList ] "}"`
	ruleMap["ElementList"] = `KeyedElement { "," KeyedElement } [ "," ]`
	ruleMap["KeyedElement"] = `[ ElementKey ":" ] Expr`
	ruleMap["ElementKey"] = `Identifier { "." Identifier } / Expr`

	// 9. Standard keywords + extra keywords
	standardKeywords := []string{
		`"break"`, `"case"`, `"chan"`, `"const"`, `"continue"`, `"default"`, `"defer"`, `"else"`,
		`"fallthrough"`, `"for"`, `"func"`, `"go"`, `"goto"`, `"if"`, `"import"`, `"interface"`,
		`"map"`, `"package"`, `"range"`, `"return"`, `"select"`, `"struct"`, `"switch"`,
		`"type"`, `"var"`,
	}
	for _, kw := range f.Overlay.ExtraKeywords {
		standardKeywords = append(standardKeywords, kw)
	}
	ruleMap["Keyword"] = strings.Join(standardKeywords, " / ")

	// 10. Terminals & Lexical rules
	ruleMap["Identifier"] = `[a-zA-Z_] { [a-zA-Z0-9_] }`
	ruleMap["Number"] = `[0-9] { [0-9] }`
	ruleMap["StringLit"] = `"\"" { StringChar } "\""`
	ruleMap["StringChar"] = `[a-zA-Z0-9_]`
	ruleMap["Comment"] = `"//" { StringChar }`

	// Canonical rule order for deterministic, stable output
	canonicalOrder := []string{
		"Grammar", "GoProgram", "StatementList", "Statement", "ImportStmt", "DeclStmt",
		"FuncDecl", "Receiver", "Signature", "ResultType", "Block",
		"TypeParamList", "TypeParamDecl", "TypeConstraint", "TypeList", "ParamList", "Param", "IdentifierList",
		"InterfaceType", "InterfaceElem", "InterfaceMethod",
		"StructType", "FieldDecl",
		"DataType", "BaseType", "DecimalType", "GenericType", "TypeArgs",
		"Expr", "PrimaryExpr", "Selector", "Index", "Arguments", "ArgList",
		"CompositeLit", "LiteralValue", "ElementList", "KeyedElement", "ElementKey",
		"Keyword", "Identifier", "Number", "StringLit", "StringChar", "Comment",
	}

	var finalRules []spec_ingest.RawRule
	for _, name := range canonicalOrder {
		if body, ok := ruleMap[name]; ok {
			finalRules = append(finalRules, spec_ingest.RawRule{Name: name, Body: body})
		} else {
			return nil, fmt.Errorf("canonical rule %s missing from synthesized set", name)
		}
	}

	return finalRules, nil
}

func (f *RuleOverlayFunctor) ensureRule(m map[string]string, order *[]string, name, body string) {
	if _, exists := m[name]; !exists {
		m[name] = body
		*order = append(*order, name)
	}
}

func convertPipesToSlashes(s string) string {
	var sb strings.Builder
	inQuote := false
	for i := 0; i < len(s); i++ {
		ch := s[i]
		if ch == '"' {
			inQuote = !inQuote
		}
		if ch == '|' && !inQuote {
			sb.WriteByte('/')
		} else {
			sb.WriteByte(ch)
		}
	}
	return sb.String()
}
