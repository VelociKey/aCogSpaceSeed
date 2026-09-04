package spec_ingest

import (
	"fmt"
	"strings"

	languagegrammar "sov.fleet/languagegrammar/81000-active-source"
	macroir "sov.fleet/macroir/81000-active-source"
)

// EBNFASTNodeKind classifies an AST node in the Go EBNF specification tree.
type EBNFASTNodeKind uint8

const (
	NodeProduction  EBNFASTNodeKind = iota // Production rule definition
	NodeAlternation                        // Alternation expression ('|' / '/')
	NodeSequence                           // Sequence of terms
	NodeOption                             // Optional expression ('[' ... ']')
	NodeRepetition                         // Repetition expression ('{' ... '}')
	NodeGroup                              // Grouped expression ('(' ... ')')
	NodeTerminal                           // Quoted terminal token
	NodeNonTerminal                        // Identifier / rule reference
	NodeOperator                           // Operator (range '...', etc.)
)

func (k EBNFASTNodeKind) String() string {
	switch k {
	case NodeProduction:
		return "Production"
	case NodeAlternation:
		return "Alternation"
	case NodeSequence:
		return "Sequence"
	case NodeOption:
		return "Option"
	case NodeRepetition:
		return "Repetition"
	case NodeGroup:
		return "Group"
	case NodeTerminal:
		return "Terminal"
	case NodeNonTerminal:
		return "NonTerminal"
	case NodeOperator:
		return "Operator"
	default:
		return "Unknown"
	}
}

// EBNFASTNode represents an AST node mapped directly to a 128-bit Macro-IR opcode.
type EBNFASTNode struct {
	Kind     EBNFASTNodeKind
	Opcode   macroir.MacroOpcode
	Text     string
	Children []*EBNFASTNode
}

// ParseEBNFWithScanner tokenizes and parses Go specification EBNF using the
// sovereign GoEBNFSpecScanner from languagegrammar.DefaultRegistry().
// It returns both legacy-compatible RawRule structures and the 128-bit Macro-IR opcode AST.
func ParseEBNFWithScanner(src []byte) ([]RawRule, []*EBNFASTNode, error) {
	reg := languagegrammar.DefaultRegistry()
	scanner := reg.ForLanguage("go_ebnf_spec")
	if scanner == nil {
		scanner = languagegrammar.NewGoEBNFSpecScanner()
	}

	tokens := scanner.ScanTokenize(src)
	if len(tokens) == 0 {
		return nil, nil, fmt.Errorf("no tokens extracted from specification source")
	}

	// Filter tokens if wrapped in HTML <pre class="ebnf"> or <pre class="grammar"> tags
	filteredTokens := filterHTMLGrammarBlocks(tokens)
	if len(filteredTokens) == 0 {
		filteredTokens = tokens
	}

	var rawRules []RawRule
	var astRoots []*EBNFASTNode

	i := 0
	n := len(filteredTokens)

	for i < n {
		// Production definition: Identifier "=" ... "."
		if i+1 < n && filteredTokens[i].Kind == languagegrammar.TokIdentifier && filteredTokens[i+1].Text == "=" {
			prodName := filteredTokens[i].Text
			i += 2 // skip identifier and "="

			startBody := i
			for i < n && !isRuleTerminator(filteredTokens[i]) {
				i++
			}
			prodTokens := filteredTokens[startBody:i]
			if i < n && isRuleTerminator(filteredTokens[i]) {
				i++ // skip "."
			}

			// Format body string for RawRule
			var bodyParts []string
			for _, tok := range prodTokens {
				if tok.Kind == languagegrammar.TokString {
					bodyParts = append(bodyParts, fmt.Sprintf("%q", tok.Text))
				} else {
					bodyParts = append(bodyParts, tok.Text)
				}
			}
			bodyStr := strings.Join(bodyParts, " ")

			rawRules = append(rawRules, RawRule{
				Name: prodName,
				Body: bodyStr,
			})

			// Build Macro-IR AST Node for this production
			prodAST := &EBNFASTNode{
				Kind:     NodeProduction,
				Opcode:   macroir.OpAllocSlot,
				Text:     prodName,
				Children: lowerTokensToAST(prodTokens),
			}
			astRoots = append(astRoots, prodAST)
		} else {
			i++
		}
	}

	if len(rawRules) == 0 {
		return nil, nil, fmt.Errorf("no valid EBNF production rules found in token stream")
	}

	return rawRules, astRoots, nil
}

// filterHTMLGrammarBlocks extracts tokens situated inside <pre class="ebnf"> or <pre class="grammar">
func filterHTMLGrammarBlocks(tokens []languagegrammar.Token) []languagegrammar.Token {
	var filtered []languagegrammar.Token
	inBlock := false
	i := 0
	n := len(tokens)

	for i < n {
		tok := tokens[i]

		// Check for block closing: < / pre >
		if inBlock && tok.Text == "<" && i+3 < n && tokens[i+1].Text == "/" && tokens[i+2].Text == "pre" && tokens[i+3].Text == ">" {
			inBlock = false
			i += 4
			continue
		}

		// Check for block opening: < pre ... >
		if !inBlock && tok.Text == "<" && i+1 < n && tokens[i+1].Text == "pre" {
			// Scan ahead for matching '>'
			closeIdx := -1
			isGrammar := false
			for j := i + 2; j < n && j < i+12; j++ {
				if tokens[j].Text == ">" {
					closeIdx = j
					break
				}
				if tokens[j].Text == "ebnf" || tokens[j].Text == "grammar" {
					isGrammar = true
				}
			}
			if closeIdx != -1 && isGrammar {
				inBlock = true
				i = closeIdx + 1
				continue
			}
		}

		if inBlock {
			filtered = append(filtered, tok)
		}
		i++
	}

	return filtered
}

// lowerTokensToAST translates raw weBNF tokens into Macro-IR opcode AST nodes
func lowerTokensToAST(tokens []languagegrammar.Token) []*EBNFASTNode {
	var nodes []*EBNFASTNode

	for _, tok := range tokens {
		switch tok.Kind {
		case languagegrammar.TokIdentifier:
			nodes = append(nodes, &EBNFASTNode{
				Kind:   NodeNonTerminal,
				Opcode: macroir.OpBorrowRef,
				Text:   tok.Text,
			})
		case languagegrammar.TokString:
			nodes = append(nodes, &EBNFASTNode{
				Kind:   NodeTerminal,
				Opcode: macroir.OpAllocSlot,
				Text:   tok.Text,
			})
		case languagegrammar.TokOperator:
			op := macroir.OpAluExec
			kind := NodeOperator
			if tok.Text == "|" || tok.Text == "/" {
				op = macroir.OpBranchCMOV
				kind = NodeAlternation
			} else if tok.Text == "..." {
				op = macroir.OpPredScan
			}
			nodes = append(nodes, &EBNFASTNode{
				Kind:   kind,
				Opcode: op,
				Text:   tok.Text,
			})
		case languagegrammar.TokDelimiter:
			op := macroir.OpNop
			kind := NodeGroup
			if tok.Text == "[" || tok.Text == "]" {
				op = macroir.OpBranchCMOV
				kind = NodeOption
			} else if tok.Text == "{" || tok.Text == "}" {
				op = macroir.OpPredScan
				kind = NodeRepetition
			}
			nodes = append(nodes, &EBNFASTNode{
				Kind:   kind,
				Opcode: op,
				Text:   tok.Text,
			})
		default:
			nodes = append(nodes, &EBNFASTNode{
				Kind:   NodeOperator,
				Opcode: macroir.OpNop,
				Text:   tok.Text,
			})
		}
	}

	return nodes
}

func isRuleTerminator(tok languagegrammar.Token) bool {
	return tok.Kind == languagegrammar.TokOperator && tok.Text == "."
}
