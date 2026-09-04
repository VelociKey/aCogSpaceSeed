package overlay

import (
	"fmt"
	"strings"

	languagegrammar "sov.fleet/languagegrammar/81000-active-source"
)

// TargetedExtension models a single targeted modification to the Go specification.
type TargetedExtension struct {
	Name           string   // Identifier of the extension block
	TargetRule     string   // Production name in the Go spec being targeted
	Section        string   // Spec section (types, declarations, interfaces, expressions, lexical)
	Action         string   // AUGMENT, REPLACE, INJECT, FORBID, REGISTER_KEYWORD
	Feature        string   // Semantic feature identifier (e.g., GENERIC_RECEIVER_METHODS)
	Lifecycle      string   // ACTIVE, SUBSUMED, DEPRECATED
	SubsumedIn     string   // Go version where this was natively subsumed (e.g., "1.27")
	SubsumedIf     string   // Predicate pattern to test against upstream Go spec
	ForbidElements []string // Elements forbidden by negative constraints
	Keywords       []string // Lexical keywords registered by the extension
	ProductionBody string   // Extended production expression in weBNF.sn
}

// ExtensionSpec represents the complete parsed Go extension specification.
type ExtensionSpec struct {
	Name          string
	BaseVersion   string
	TargetGrammar string
	Description   string
	Extensions    []TargetedExtension
}

// ParseExtensionSpec parses a declarative Go extension specification file (.ext.wag)
// using the sovereign GoExtensionSpecScanner.
func ParseExtensionSpec(src []byte) (*ExtensionSpec, error) {
	reg := languagegrammar.DefaultRegistry()
	scanner := reg.ForLanguage("go_extension_spec")
	if scanner == nil {
		scanner = languagegrammar.NewGoExtensionSpecScanner()
	}

	tokens := scanner.ScanTokenize(src)
	if len(tokens) == 0 {
		return nil, fmt.Errorf("no tokens in extension specification")
	}

	spec := &ExtensionSpec{
		Name:          "SovereignGo",
		BaseVersion:   "1.27",
		TargetGrammar: "go",
	}

	// Split token stream by semicolon ';' delimiter into individual declarations
	var currentDecl []languagegrammar.Token
	var declarations [][]languagegrammar.Token

	for _, tok := range tokens {
		if tok.Kind == languagegrammar.TokSemicolon || tok.Text == ";" {
			if len(currentDecl) > 0 {
				declarations = append(declarations, currentDecl)
				currentDecl = nil
			}
		} else {
			currentDecl = append(currentDecl, tok)
		}
	}
	if len(currentDecl) > 0 {
		declarations = append(declarations, currentDecl)
	}

	for _, decl := range declarations {
		if len(decl) < 3 || decl[0].Kind != languagegrammar.TokIdentifier || decl[1].Text != "=" {
			continue
		}

		declName := decl[0].Text
		bodyTokens := decl[2:]

		if declName == "ExtensionMetadata" {
			parseMetadataDecl(bodyTokens, spec)
			continue
		}

		ext := parseTargetedExtension(declName, bodyTokens)
		spec.Extensions = append(spec.Extensions, ext)
	}

	return spec, nil
}

func parseMetadataDecl(tokens []languagegrammar.Token, spec *ExtensionSpec) {
	for i := 0; i+2 < len(tokens); i++ {
		key := tokens[i].Text
		colon := tokens[i+1].Text
		val := tokens[i+2].Text

		if colon == ":" {
			switch key {
			case "name":
				spec.Name = val
			case "base_version":
				spec.BaseVersion = val
			case "target_grammar", "target_language":
				spec.TargetGrammar = val
			case "description":
				spec.Description = val
			}
		}
	}
}

func parseTargetedExtension(name string, tokens []languagegrammar.Token) TargetedExtension {
	ext := TargetedExtension{Name: name}

	var prodTokens []languagegrammar.Token
	i := 0
	n := len(tokens)

	for i < n {
		// Check for metadata key-value pair: "key" ":" "val"
		if i+2 < n && tokens[i+1].Text == ":" {
			key := tokens[i].Text
			val := tokens[i+2].Text
			switch key {
			case "target_rule":
				ext.TargetRule = val
				i += 3
				if i < n && (tokens[i].Text == "/" || tokens[i].Text == "|") {
					i++
				}
				continue
			case "section":
				ext.Section = val
				i += 3
				if i < n && (tokens[i].Text == "/" || tokens[i].Text == "|") {
					i++
				}
				continue
			case "action":
				ext.Action = val
				i += 3
				if i < n && (tokens[i].Text == "/" || tokens[i].Text == "|") {
					i++
				}
				continue
			case "feature":
				ext.Feature = val
				i += 3
				if i < n && (tokens[i].Text == "/" || tokens[i].Text == "|") {
					i++
				}
				continue
			case "lifecycle":
				ext.Lifecycle = val
				i += 3
				if i < n && (tokens[i].Text == "/" || tokens[i].Text == "|") {
					i++
				}
				continue
			case "subsumed_in":
				ext.SubsumedIn = val
				i += 3
				if i < n && (tokens[i].Text == "/" || tokens[i].Text == "|") {
					i++
				}
				continue
			case "subsumed_if":
				ext.SubsumedIf = val
				i += 3
				if i < n && (tokens[i].Text == "/" || tokens[i].Text == "|") {
					i++
				}
				continue
			case "forbid":
				ext.ForbidElements = append(ext.ForbidElements, val)
				i += 3
				if i < n && (tokens[i].Text == "/" || tokens[i].Text == "|") {
					i++
				}
				continue
			}
		}

		// Check for keyword registrations
		if ext.Action == "REGISTER_KEYWORD" {
			if tokens[i].Kind == languagegrammar.TokString {
				ext.Keywords = append(ext.Keywords, tokens[i].Text)
			}
			i++
			continue
		}

		// Otherwise, accumulate into production tokens
		prodTokens = append(prodTokens, tokens[i])
		i++
	}

	// Unwrap outer grouping parentheses if the entire production is wrapped: ( ... )
	if len(prodTokens) >= 2 && prodTokens[0].Text == "(" && prodTokens[len(prodTokens)-1].Text == ")" {
		prodTokens = prodTokens[1 : len(prodTokens)-1]
	}

	// Reconstruct formatted production body
	var sb strings.Builder
	for idx, tok := range prodTokens {
		if idx > 0 {
			prev := prodTokens[idx-1].Text
			curr := tok.Text
			if curr != "[" && curr != "]" && curr != "(" && curr != ")" && curr != "{" && curr != "}" &&
				prev != "[" && prev != "(" && prev != "{" {
				sb.WriteString(" ")
			} else if curr == "/" || prev == "/" {
				sb.WriteString(" ")
			}
		}
		if tok.Kind == languagegrammar.TokString {
			sb.WriteString(fmt.Sprintf("%q", tok.Text))
		} else {
			sb.WriteString(tok.Text)
		}
	}

	ext.ProductionBody = sb.String()
	return ext
}
