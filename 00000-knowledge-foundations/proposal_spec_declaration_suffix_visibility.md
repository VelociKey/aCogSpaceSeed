# Proposal: spec: declaration-suffix visibility modifiers (`public`, `private`, `package`)

**Author:** Sovereign Architecture & Antigravity Team  
**Date:** September 2026  
**Discussion / Issue:** golang/go#issue-proposal-declaration-suffix-visibility  

---

## 1. Abstract

We propose introducing optional declaration-tail visibility modifiers—`public`, `private`, and `package`—to the Go specification for top-level declarations (`func`, `type`, `var`, `const`) and struct fields. 

By placing the visibility modifier at the **tail** of declarations (e.g., `func parseConfig() (Config, error) package { ... }`, `type token struct { ... } private`, `Field int public`), this proposal:
1. **Decouples semantic encapsulation from typography**: Eliminates naming friction with acronyms (`url`, `json`, `id`), Unicode scripts without casing, and awkward renaming.
2. **Eliminates directory-as-grammar (`internal/`) leakage**: Provides explicit, source-level package/subsystem encapsulation without forcing artificial filesystem hierarchy or path re-nesting.
3. **Preserves standard Go parser & AST lookahead**: Avoids prefix keywords (`pub`, `private`) which break existing Go tooling, parsers (`go/parser`), linters, DWARF symbols, and LSP engines (`gopls`).
4. **Guarantees 100% Go 1 Backward Compatibility**: Any declaration omitting a visibility suffix defaults to standard Go identifier capitalization semantics.

---

## 2. Background & Problem Statement

### 2.1 Conflating Typography with Semantic Visibility
Since its inception, Go has used identifier capitalization to dictate visibility: uppercase identifiers are exported; lowercase identifiers are unexported (package-scoped). While concise, this design has accumulated significant ergonomic and architectural drawbacks:

- **Acronym and Casing Clashes**: Field names or identifiers that naturally start with an acronym must either violate idiomatic naming conventions or distort casing (e.g., `url` vs `URL`, `jsonEncoder` vs `JSONEncoder`).
- **Non-Latin Scripts**: In many international scripts (CJK, Arabic, Hebrew, Devanagari), case distinction does not exist. Casing-based export creates awkward disparities for non-Latin identifiers.
- **Refactoring Breakage**: Changing an unexported symbol to exported changes its lexical name, requiring cascading find-and-replace throughout call sites rather than adjusting an access modifier in place.
- **No True File or Module Scope**: Go has no mechanism to make a helper private to a single `.go` source file; all unexported declarations are visible to the entire package.

### 2.2 Directory Structure as a Pseudo-Grammar (`internal/`)
To achieve encapsulation across multiple packages without exposing APIs publicly, Go introduced the `internal/` directory convention (Go 1.4). 

However, making the filesystem layout an intrinsic semantic element of the language grammar violates the principle of AST self-containment:
- Code encapsulation should be defined declaratively in the source AST, not dictated by physical directory placement on disk.
- Splitting cohesive packages merely to satisfy `internal/` boundaries introduces module fragmentation, circular import traps, and cross-package allocation overhead.

### 2.3 Why Prefix Modifiers Fail in the Go Ecosystem
Proposals in other languages often suggest leading keywords (e.g., `pub func Foo()`, `private type Bar struct`). In Go, prefix modifiers present critical problems:
- **Parser Lookahead & Token Ambiguity**: `go/parser` relies on standard declaration start tokens (`func`, `type`, `var`, `const`, `import`, `package`). Introducing prefix modifiers complicates grammar lookahead and introduces ambiguity with package names or custom types.
- **Tooling Breakage**: Code indexers, symbol extractors, syntax highlighters, and AST analyzers look for the canonical declaration keyword at the start of a statement. A prefix keyword breaks downstream AST assumptions across the entire Go ecosystem.

---

## 3. Proposed Solution: Declaration-Suffix Modifiers

We propose placing optional visibility keywords at the **declaration tail**, immediately preceding the block body or assignment.

### 3.1 Syntax Overview

#### Functions and Methods
```go
// Exported outside package, despite lowercase name
func generateUUID() string public {
    return uuid.NewString()
}

// Package-private helper (explicit, even if capitalized for domain reasons)
func ValidateHeader(h Header) error package {
    ...
}

// Strictly private to this compilation unit / source file
func internalBufferPool() *sync.Pool private {
    ...
}
```

#### Types
```go
// Exported type with lowercase identifier
type config struct {
    Host string public
    Port int    public
    salt string private
} public

// Package-private type
type connectionPool struct {
    ...
} package
```

#### Variables and Constants
```go
var DefaultTimeout = 30 * time.Second public
const maxRetries = 5 package
var localScratch = make([]byte, 1024) private
```

#### Struct Fields
```go
type Request struct {
    // Explicit public field matching lowercase JSON convention
    id        int64  public  `json:"id"`
    authToken string private // Strictly invisible outside the declaring file/type
    status    string package // Visible to package, invisible to consumers
}
```

---

## 4. Formal Grammar Specification (EBNF / WSN)

The formal grammar additions are flat, non-recursive, and integrate cleanly into the existing Go specification:

```ebnf
VisibilitySuffix = "public" | "private" | "package" ;

FuncDecl = "func" [ Receiver ] Identifier [ TypeParamList ] Signature [ VisibilitySuffix ] [ Block ] ;

TypeSpec = Identifier [ TypeParamList ] [ VisibilitySuffix ] [ "=" ] Type ;

FieldDecl = ( IdentifierList Type | EmbeddedField ) [ VisibilitySuffix ] [ Tag ] ;

VarSpec = IdentifierList ( ( Type [ VisibilitySuffix ] [ "=" ExpressionList ] ) | ( [ VisibilitySuffix ] "=" ExpressionList ) ) ;

ConstSpec = IdentifierList [ [ Type ] [ VisibilitySuffix ] ] "=" ExpressionList ;
```

### Parser Invariant
The parser encounters `func`, `type`, `var`, or `const` in standard position. The identifier is parsed unambiguously. When completing the signature or type specification, the parser checks for an optional terminal token `public`, `private`, or `package`.

---

## 5. Semantic Rules

1. **Precedence**:
   - When a `VisibilitySuffix` is explicitly present, it **overrides** the capitalization of the identifier.
   - Example: `func auditLog() public` is exported. `func SystemReset() private` is unexported and private to the compilation file.
2. **Default Fallback (Zero Breaking Changes)**:
   - When `VisibilitySuffix` is omitted, the symbol's visibility defaults strictly to Go 1 capitalization rules (`Uppercase` = exported, `lowercase` = package-scoped).
3. **Levels of Visibility**:
   - `public`: Accessible to any importing package.
   - `package`: Accessible anywhere within the same Go package (replaces the need for arbitrary `internal/` directory packaging).
   - `private`: Accessible only within the declaring source file / lexical scope.

---

## 6. Ecosystem & Tooling Impact

- **`go/parser` & `go/ast`**: Add an optional `Visibility *ast.Ident` or `VisibilityToken token.Pos` to `ast.FuncDecl`, `ast.TypeSpec`, and `ast.Field`.
- **`go/types`**: During object resolution, query `obj.Exported()` via the explicit visibility specifier if present, falling back to `token.IsExported(obj.Name())`.
- **`gopls` & Autocompletion**: Autocompletion engines can offer completions based on semantic access rather than casing heuristics.
- **Reflection (`reflect`)**: Struct field reflection will respect `Field.IsExported()` using the field's explicit visibility.

---

## 7. Compatibility Guarantee

This proposal strictly satisfies the **Go 1 Compatibility Promise**:
- Every existing valid Go program remains valid with identical semantic behavior.
- `public`, `private`, and `package` are contextually recognized as visibility modifiers at declaration tail positions, avoiding breaking code where these words are used as local variable names.

---

## 8. Summary Comparison

| Aspect | Go 1 (Status Quo) | Prefix Modifiers (`pub func`) | Proposed: Declaration-Suffix (`func ... public`) |
| :--- | :--- | :--- | :--- |
| **Export Mechanism** | Identifier Capitalization | Leading Keyword | Trailing Declaration Suffix |
| **Parser Disruption** | None | High (breaks statement start) | None (canonical leading keyword preserved) |
| **Acronym Ergonomics** | Poor (`json` vs `JSON`) | Good | Excellent |
| **Directory Coupling** | Relies on `internal/` | Relies on `internal/` | Purely Source-Level Encapsulation |
| **File-Level Privacy** | Not supported | Varied | Supported (`private`) |
| **Go 1 Compatible** | Baseline | Breaking or complex lookahead | 100% Backwards-Compatible |
