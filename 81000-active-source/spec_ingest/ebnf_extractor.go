package spec_ingest

import (
	"fmt"
	"os"
	"regexp"
	"strings"
)

var preEbnfRegex = regexp.MustCompile(`(?s)<pre class="ebnf">(.*?)</pre>`)

// ExtractedEBNF holds the raw EBNF text extracted from the specification.
type ExtractedEBNF struct {
	SourcePath string
	RuleBlocks []string
	FullText   string
}

// ExtractEBNFFromFile reads a Go language specification (HTML or text) and extracts all EBNF blocks.
func ExtractEBNFFromFile(path string) (*ExtractedEBNF, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read spec file %s: %w", path, err)
	}
	content := string(data)

	res := &ExtractedEBNF{SourcePath: path}

	// Check if it's HTML with <pre class="ebnf">
	if strings.Contains(content, `<pre class="ebnf">`) {
		matches := preEbnfRegex.FindAllStringSubmatch(content, -1)
		for _, m := range matches {
			if len(m) > 1 {
				cleaned := strings.TrimSpace(m[1])
				if cleaned != "" {
					res.RuleBlocks = append(res.RuleBlocks, cleaned)
				}
			}
		}
		res.FullText = strings.Join(res.RuleBlocks, "\n\n")
		return res, nil
	}

	// Plain text EBNF fallback
	res.FullText = content
	res.RuleBlocks = []string{content}
	return res, nil
}
