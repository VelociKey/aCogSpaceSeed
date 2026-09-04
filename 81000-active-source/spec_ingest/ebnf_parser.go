package spec_ingest

import (
	"fmt"
	"strings"
)

// RawRule represents a single parsed production rule from the Go specification EBNF.
type RawRule struct {
	Name string
	Body string
}

// ParseEBNFRules splits raw EBNF text into individual production rules.
// Go spec EBNF terminates rules with a period '.' at the end of the production.
func ParseEBNFRules(ebnfText string) ([]RawRule, error) {
	var rules []RawRule

	lines := strings.Split(ebnfText, "\n")
	var currentRuleName string
	var currentBody strings.Builder

	flushRule := func() {
		if currentRuleName != "" {
			bodyStr := strings.TrimSpace(currentBody.String())
			bodyStr = strings.TrimSuffix(bodyStr, ".")
			bodyStr = strings.TrimSpace(bodyStr)
			if bodyStr != "" {
				rules = append(rules, RawRule{
					Name: currentRuleName,
					Body: bodyStr,
				})
			}
			currentRuleName = ""
			currentBody.Reset()
		}
	}

	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "//") || strings.HasPrefix(line, "/*") {
			continue
		}

		// Check for new rule start: "RuleName = ..."
		if eqIdx := strings.Index(line, "="); eqIdx > 0 {
			candidateName := strings.TrimSpace(line[:eqIdx])
			// Check if candidateName is a single identifier (no spaces inside)
			if isValidIdentifier(candidateName) {
				flushRule()
				currentRuleName = candidateName
				rest := strings.TrimSpace(line[eqIdx+1:])
				currentBody.WriteString(rest)
				if strings.HasSuffix(rest, ".") {
					flushRule()
				}
				continue
			}
		}

		// Continuation line
		if currentRuleName != "" {
			currentBody.WriteString(" ")
			currentBody.WriteString(line)
			if strings.HasSuffix(line, ".") {
				flushRule()
			}
		}
	}

	flushRule()

	if len(rules) == 0 {
		return nil, fmt.Errorf("no valid EBNF production rules found in input")
	}

	return rules, nil
}

func isValidIdentifier(s string) bool {
	if len(s) == 0 {
		return false
	}
	for i, r := range s {
		if i == 0 {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || r == '_') {
				return false
			}
		} else {
			if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '_') {
				return false
			}
		}
	}
	return true
}
