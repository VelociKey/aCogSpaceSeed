# Proposal: spec: declaration-suffix visibility modifiers (`public`, `private`)

**Author:** Sovereign Architecture & Antigravity Team  
**Date:** September 2026  
**Discussion / Issue:** golang/go#issue-proposal-declaration-suffix-visibility  

---

## 1. Abstract

We propose introducing optional declaration-tail visibility modifiers—`public` and `private`—to the Go language specification for top-level declarations (`func`, `type`, `var`, `const`) and struct fields. Package-level visibility remains the natural language default when modifiers are omitted.

By placing the visibility modifier at the **tail** of declarations (e.g., `func parseConfig() (Config, error) public { ... }`, `type token struct { ... } private`, `field int public`), this proposal:
1. **Decouples semantic encapsulation from typography**: Eliminates naming friction with acronyms (`url`, `json`, `id`), Unicode scripts without casing, and awkward renaming.
2. **Eliminates directory-as-grammar (`internal/`) leakage**: Provides explicit, source-level file/package encapsulation without forcing artificial filesystem hierarchy or directory re-nesting.
3. **Preserves standard Go parser & AST lookahead**: Avoids prefix keywords (`pub`, `private`) which break existing Go tooling, parsers (`go/parser`), linters, DWARF symbols, and LSP engines (`gopls`).
4. **Adheres to Go's Minimalist Aesthetic**: Omits any redundant `package` keyword—package-level scope remains the universal default for unexported declarations.
5. **Guarantees 100% Go 1 Backward Compatibility**: Any declaration omitting a visibility suffix defaults to standard Go identifier capitalization semantics.

---

## 2. Background & Problem Statement

### 2.1 Conflating Typography with Semantic Visibility
Since its inception, Go has used identifier capitalization to dictate visibility: uppercase identifiers are exported; lowercase identifiers are unexported (package-scoped). While concise, this design has accumulated significant ergonomic and architectural drawbacks:

- **Acronym and Casing Clashes**: Field names or identifiers that naturally start with an acronym must either violate idiomatic naming conventions or distort casing (e.g., `url` vs `URL`, `jsonEncoder` vs `JSONEncoder`).
- **Non-Latin Scripts**: In many international scripts (CJK, Arabic, Hebrew, Devanagari), case distinction does not exist. Casing-based export creates awkward disparities for non-Latin identifiers.
- **Refactoring Breakage**: Changing an unexported symbol to exported changes its lexical name, requiring cascading find-and-replace throughout call sites rather than adjusting an access modifier in place.
- **No True File-Level Scope**: Go has no mechanism to make a helper private to a single `.go` source file; all unexported declarations are unconditionally visible across the entire package.

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

We propose placing optional visibility keywords (`public`, `private`) at the **declaration tail**, immediately preceding the block body or assignment.

### 3.1 Syntax Overview

#### Functions and Methods
```go
// Exported outside package, despite lowercase name
func generateUUID() string public {
    return uuid.NewString()
}

// Strictly private to this compilation unit / source file
func internalBufferPool() *sync.Pool private {
    ...
}

// Unmodified: Standard Go default (package-scoped because lowercase)
func validateHeader(h Header) error {
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

// Strictly private type (file-scoped)
type connectionPool struct {
    ...
} private

// Unmodified: Standard Go default (package-scoped)
type sessionState struct {
    ...
}
```

#### Variables and Constants
```go
var DefaultTimeout = 30 * time.Second public
var localScratch = make([]byte, 1024) private

// Unmodified: Standard Go defaults
var retryLimit = 5
const defaultPort = 8080
```

#### Struct Fields
```go
type Request struct {
    // Explicit public field matching lowercase JSON convention
    id        int64  public  `json:"id"`
    authToken string private // Strictly invisible outside declaring file/type
    status    string         // Default package visibility
}
```

---

## 4. Formal Grammar Specification (EBNF / WSN)

The formal grammar additions are flat, non-recursive, and integrate cleanly into the existing Go specification:

```ebnf
VisibilitySuffix = "public" | "private" ;

FuncDecl = "func" [ Receiver ] Identifier [ TypeParamList ] Signature [ VisibilitySuffix ] [ Block ] ;

TypeSpec = Identifier [ TypeParamList ] [ VisibilitySuffix ] [ "=" ] Type ;

FieldDecl = ( IdentifierList Type | EmbeddedField ) [ VisibilitySuffix ] [ Tag ] ;

VarSpec = IdentifierList ( ( Type [ VisibilitySuffix ] [ "=" ExpressionList ] ) | ( [ VisibilitySuffix ] "=" ExpressionList ) ) ;

ConstSpec = IdentifierList [ [ Type ] [ VisibilitySuffix ] ] "=" ExpressionList ;
```

### Parser Invariant
The parser encounters `func`, `type`, `var`, or `const` in standard position. The identifier is parsed unambiguously. When completing the signature or type specification, the parser checks for an optional terminal token `public` or `private`.

---

## 5. Semantic Rules

1. **Precedence**:
   - When a `VisibilitySuffix` is explicitly present, it **overrides** the capitalization of the identifier.
   - Example: `func auditLog() public` is exported outside the package. `func SystemReset() private` is unexported and private to the compilation file.
2. **Default Fallback (Zero Breaking Changes)**:
   - When `VisibilitySuffix` is omitted, the symbol's visibility defaults strictly to Go 1 capitalization rules (`Uppercase` = exported, `lowercase` = package-scoped).
3. **Levels of Visibility**:
   - `public`: Accessible to any importing package.
   - `private`: Accessible only within the declaring source file (`.go` compilation unit).
   - *(omitted)*: Standard Go behavior (package-scoped if lowercase; exported if uppercase).

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
- `public` and `private` are contextually recognized as visibility modifiers at declaration tail positions, avoiding breaking code where these words are used as local variable names.

---

## 8. Summary Comparison

| Aspect | Go 1 (Status Quo) | Prefix Modifiers (`pub func`) | Proposed: Declaration-Suffix (`func ... public`) |
| :--- | :--- | :--- | :--- |
| **Export Mechanism** | Identifier Capitalization | Leading Keyword | Trailing Declaration Suffix |
| **Parser Disruption** | None | High (breaks statement start) | None (canonical leading keyword preserved) |
| **Acronym Ergonomics** | Poor (`json` vs `JSON`) | Good | Excellent |
| **Directory Coupling** | Relies on `internal/` | Relies on `internal/` | Purely Source-Level Encapsulation |
| **File-Level Privacy** | Not supported | Varied | Supported (`private`) |
| **Go Minimalism** | High | Low (introduces new syntax blocks) | High (only `public` and `private`) |
| **Go 1 Compatible** | Baseline | Breaking or complex lookahead | 100% Backwards-Compatible |
