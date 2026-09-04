package ghost

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"regexp"
	"sort"
	"unicode"
	"unicode/utf8"
)

// VisibilityKind represents the semantic access level of a symbol.
type VisibilityKind uint8

const (
	VisUnspecified VisibilityKind = iota
	VisPublic                     // Exported outside package
	VisPrivate                    // Package-internal or file-scoped
)

func (v VisibilityKind) String() string {
	switch v {
	case VisPublic:
		return "public"
	case VisPrivate:
		return "private"
	default:
		return "unspecified"
	}
}

// VisibilityAttr records the effective visibility and whether it was ghost-materialized.
type VisibilityAttr struct {
	Kind    VisibilityKind `json:"kind"`
	IsGhost bool           `json:"is_ghost"`
	Keyword string         `json:"keyword"`
}

// DeclarationKind identifies the syntactical role of a declaration.
type DeclarationKind string

const (
	DeclFunc  DeclarationKind = "func"
	DeclType  DeclarationKind = "type"
	DeclField DeclarationKind = "field"
	DeclVar   DeclarationKind = "var"
	DeclConst DeclarationKind = "const"
)

// MaterializedDeclaration represents an analyzed declaration in the AST.
type MaterializedDeclaration struct {
	Name           string          `json:"name"`
	Kind           DeclarationKind `json:"kind"`
	Visibility     VisibilityAttr  `json:"visibility"`
	OriginalCasing string          `json:"original_casing"` // "upper" or "lower"
	StartOffset    int             `json:"start_offset"`
	EndOffset      int             `json:"end_offset"`
	InsertOffset   int             `json:"insert_offset"` // Where ghost suffix should be injected
}

// MaterializeReport summarizes the ghost-materialization pass over a compilation unit.
type MaterializeReport struct {
	TotalDecls      int                       `json:"total_decls"`
	GhostPublic     int                       `json:"ghost_public"`
	GhostPrivate    int                       `json:"ghost_private"`
	PhysicalPublic  int                       `json:"physical_public"`
	PhysicalPrivate int                       `json:"physical_private"`
	Declarations    []MaterializedDeclaration `json:"declarations"`
}

type physicalToken struct {
	offset  int
	keyword string
}

var funcSuffixRegex = regexp.MustCompile(`\b(public|private)\s*(\{)`)

func sanitizePhysicalSuffixes(src []byte) ([]byte, []physicalToken) {
	sanitized := make([]byte, len(src))
	copy(sanitized, src)
	var tokens []physicalToken

	matches := funcSuffixRegex.FindAllSubmatchIndex(src, -1)
	for _, m := range matches {
		kwStart, kwEnd := m[2], m[3]
		kw := string(src[kwStart:kwEnd])
		tokens = append(tokens, physicalToken{offset: kwStart, keyword: kw})
		for i := kwStart; i < kwEnd; i++ {
			sanitized[i] = ' '
		}
	}
	return sanitized, tokens
}

// InferVisibilityFromCasing inspects the first rune of an identifier and returns
// a ghost-materialized VisibilityAttr reflecting the author's casing intent.
func InferVisibilityFromCasing(ident string) VisibilityAttr {
	if len(ident) == 0 {
		return VisibilityAttr{Kind: VisPrivate, IsGhost: true, Keyword: "private"}
	}
	r, _ := utf8.DecodeRuneInString(ident)
	if unicode.IsUpper(r) {
		return VisibilityAttr{Kind: VisPublic, IsGhost: true, Keyword: "public"}
	}
	return VisibilityAttr{Kind: VisPrivate, IsGhost: true, Keyword: "private"}
}

// ScanAndMaterializeSource parses Go source code, discovers all declarations,
// detects explicit suffixes or materializes ghost suffixes, and returns a detailed report.
func ScanAndMaterializeSource(src []byte) (*MaterializeReport, error) {
	sanitizedSrc, physicalTokens := sanitizePhysicalSuffixes(src)
	fset := token.NewFileSet()
	file, err := parser.ParseFile(fset, "input.go", sanitizedSrc, parser.ParseComments)
	if err != nil {
		return nil, fmt.Errorf("parsing failed: %w", err)
	}

	report := &MaterializeReport{
		Declarations: make([]MaterializedDeclaration, 0),
	}

	// Traverse top-level declarations
	for _, decl := range file.Decls {
		switch d := decl.(type) {
		case *ast.FuncDecl:
			mat := processFuncDecl(d, fset, src, physicalTokens)
			recordDecl(report, mat)

		case *ast.GenDecl:
			for _, spec := range d.Specs {
				switch s := spec.(type) {
				case *ast.TypeSpec:
					mat := processTypeSpec(s, d, fset, src)
					recordDecl(report, mat)

					// Also scan struct fields if it's a struct type
					if st, ok := s.Type.(*ast.StructType); ok && st.Fields != nil {
						for _, field := range st.Fields.List {
							fieldMats := processField(field, fset, src)
							for _, fm := range fieldMats {
								recordDecl(report, fm)
							}
						}
					}

				case *ast.ValueSpec:
					for _, name := range s.Names {
						kind := DeclVar
						if d.Tok == token.CONST {
							kind = DeclConst
						}
						mat := processValueSpec(name, s, kind, fset, src)
						recordDecl(report, mat)
					}
				}
			}
		}
	}

	return report, nil
}

func recordDecl(rep *MaterializeReport, mat MaterializedDeclaration) {
	rep.Declarations = append(rep.Declarations, mat)
	rep.TotalDecls++

	if mat.Visibility.IsGhost {
		if mat.Visibility.Kind == VisPublic {
			rep.GhostPublic++
		} else {
			rep.GhostPrivate++
		}
	} else {
		if mat.Visibility.Kind == VisPublic {
			rep.PhysicalPublic++
		} else {
			rep.PhysicalPrivate++
		}
	}
}

// processFuncDecl inspects a function declaration for an explicit suffix before Body,
// or materializes a ghost suffix from the function name.
func processFuncDecl(fn *ast.FuncDecl, fset *token.FileSet, src []byte, physicalTokens []physicalToken) MaterializedDeclaration {
	name := fn.Name.Name
	casing := "lower"
	r, _ := utf8.DecodeRuneInString(name)
	if unicode.IsUpper(r) {
		casing = "upper"
	}

	// Calculate insertion offset right before the opening brace of Body, or at end of signature
	insertOffset := -1
	if fn.Body != nil {
		insertOffset = fset.Position(fn.Body.Lbrace).Offset
	} else if fn.Type != nil {
		insertOffset = fset.Position(fn.Type.End()).Offset
	}

	// Check if any physical token falls between signature end and body brace
	startCheck := fset.Position(fn.Type.End()).Offset
	endCheck := insertOffset
	for _, pt := range physicalTokens {
		if pt.offset >= startCheck && pt.offset <= endCheck {
			kind := VisPublic
			if pt.keyword == "private" {
				kind = VisPrivate
			}
			return MaterializedDeclaration{
				Name:           name,
				Kind:           DeclFunc,
				Visibility:     VisibilityAttr{Kind: kind, IsGhost: false, Keyword: pt.keyword},
				OriginalCasing: casing,
				StartOffset:    fset.Position(fn.Pos()).Offset,
				EndOffset:      fset.Position(fn.End()).Offset,
				InsertOffset:   -1,
			}
		}
	}

	// Fallback to ghost materialization from casing
	vis := InferVisibilityFromCasing(name)
	return MaterializedDeclaration{
		Name:           name,
		Kind:           DeclFunc,
		Visibility:     vis,
		OriginalCasing: casing,
		StartOffset:    fset.Position(fn.Pos()).Offset,
		EndOffset:      fset.Position(fn.End()).Offset,
		InsertOffset:   insertOffset,
	}
}

// processTypeSpec inspects a type specification.
func processTypeSpec(ts *ast.TypeSpec, gd *ast.GenDecl, fset *token.FileSet, src []byte) MaterializedDeclaration {
	name := ts.Name.Name
	casing := "lower"
	r, _ := utf8.DecodeRuneInString(name)
	if unicode.IsUpper(r) {
		casing = "upper"
	}

	insertOffset := fset.Position(ts.End()).Offset
	vis := InferVisibilityFromCasing(name)

	return MaterializedDeclaration{
		Name:           name,
		Kind:           DeclType,
		Visibility:     vis,
		OriginalCasing: casing,
		StartOffset:    fset.Position(ts.Pos()).Offset,
		EndOffset:      fset.Position(ts.End()).Offset,
		InsertOffset:   insertOffset,
	}
}

// processField inspects a struct field.
func processField(field *ast.Field, fset *token.FileSet, src []byte) []MaterializedDeclaration {
	decls := make([]MaterializedDeclaration, 0, len(field.Names))
	insertOffset := fset.Position(field.End()).Offset

	for _, id := range field.Names {
		casing := "lower"
		r, _ := utf8.DecodeRuneInString(id.Name)
		if unicode.IsUpper(r) {
			casing = "upper"
		}
		vis := InferVisibilityFromCasing(id.Name)
		decls = append(decls, MaterializedDeclaration{
			Name:           id.Name,
			Kind:           DeclField,
			Visibility:     vis,
			OriginalCasing: casing,
			StartOffset:    fset.Position(id.Pos()).Offset,
			EndOffset:      fset.Position(id.End()).Offset,
			InsertOffset:   insertOffset,
		})
	}

	return decls
}

// processValueSpec inspects a var or const spec.
func processValueSpec(name *ast.Ident, vs *ast.ValueSpec, kind DeclarationKind, fset *token.FileSet, src []byte) MaterializedDeclaration {
	casing := "lower"
	r, _ := utf8.DecodeRuneInString(name.Name)
	if unicode.IsUpper(r) {
		casing = "upper"
	}
	vis := InferVisibilityFromCasing(name.Name)
	insertOffset := fset.Position(vs.End()).Offset

	return MaterializedDeclaration{
		Name:           name.Name,
		Kind:           kind,
		Visibility:     vis,
		OriginalCasing: casing,
		StartOffset:    fset.Position(name.Pos()).Offset,
		EndOffset:      fset.Position(name.End()).Offset,
		InsertOffset:   insertOffset,
	}
}

// ModernizeSource rewrites legacy Go source code by materializing ghost visibility tokens
// into explicit physical declaration-tail keywords ("public" or "private").
func ModernizeSource(src []byte) ([]byte, *MaterializeReport, error) {
	report, err := ScanAndMaterializeSource(src)
	if err != nil {
		return nil, nil, err
	}

	// Sort declarations with ghost tokens by insertion offset descending to avoid offset shifting
	ghostDecls := make([]MaterializedDeclaration, 0)
	for _, d := range report.Declarations {
		if d.Visibility.IsGhost && d.InsertOffset > 0 && d.Kind == DeclFunc {
			ghostDecls = append(ghostDecls, d)
		}
	}

	sort.Slice(ghostDecls, func(i, j int) bool {
		return ghostDecls[i].InsertOffset > ghostDecls[j].InsertOffset
	})

	var buf bytes.Buffer
	buf.Write(src)

	// Inject suffixes at insertion points
	res := buf.Bytes()
	for _, d := range ghostDecls {
		if d.InsertOffset < 0 || d.InsertOffset > len(res) {
			continue
		}
		// Insert " <keyword> " before '{'
		suffix := fmt.Sprintf("%s ", d.Visibility.Keyword)
		var newRes []byte
		newRes = append(newRes, res[:d.InsertOffset]...)
		newRes = append(newRes, []byte(suffix)...)
		newRes = append(newRes, res[d.InsertOffset:]...)
		res = newRes
	}

	return res, report, nil
}
