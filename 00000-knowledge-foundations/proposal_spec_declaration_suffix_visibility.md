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

## 6. Upstream Implementation Architecture & Transformation of the Capitalization Format Hack

An inspection of the Go compiler (`cmd/compile`) and runtime ABI (`internal/abi`, `reflect`) reveals that **Go's runtime and ABI do not use rune casing**. They already represent visibility as a **1-bit boolean flag**. The capitalization rule is an artifact of the compiler frontend, making the transformation to declaration-suffix modifiers straightforward.

### 6.1 The Anatomy of the Current Capitalization Hack
In upstream Go, the capitalization rule is distributed across four distinct layers:

1. **Lexical / Token Layer (`go/token/token.go`, `cmd/compile/internal/types/sym.go`)**:
   ```go
   // cmd/compile/internal/types/sym.go
   func IsExported(name string) bool {
       if r := name[0]; r < utf8.RuneSelf {
           return 'A' <= r && r <= 'Z'
       }
       r, _ := utf8.DecodeRuneInString(name)
       return unicode.IsUpper(r)
   }
   ```
   Every visibility check triggers rune decoding and Unicode uppercase table lookups.

2. **AST & Typechecker Layer (`go/types/object.go`)**:
   ```go
   func (obj *object) Exported() bool {
       return isExported(obj.name) // Re-decodes string runes on every invocation!
   }
   ```
   Because AST nodes (`ast.FuncDecl`, `ast.Field`) store no visibility metadata, `obj.Exported()` is forced to re-parse the string characters on demand.

3. **Selector Disambiguation (`go/types/object.go`)**:
   ```go
   func Id(pkg *Package, name string) string {
       if isExported(name) {
           return name
       }
       return pkg.path + "." + name
   }
   ```
   Unexported identifiers are qualified with `package.path` so that unexported symbols across packages do not collide in selector resolution and method sets.

4. **Runtime ABI & Reflection Layer (`cmd/compile/internal/reflectdata/reflect.go`, `internal/abi/type.go`)**:
   ```go
   // cmd/compile/internal/reflectdata/reflect.go
   var bits byte
   if exported {
       bits |= 1 << 0 // Bit 0 is the EXPORTED bit
   }
   ```
   ```go
   // internal/abi/type.go
   func (n Name) IsExported() bool {
       return (*n.Bytes)&(1<<0) != 0
   }
   ```
   At runtime, `reflect.StructField.IsExported()` and interface method dispatches query `bit 0`. **The compiled runtime binary never inspects string casing.**

### 6.2 Step-by-Step Compiler Transformation

```
                  ┌──────────────────────────────────────────────────────────┐
                  │                 CURRENT UPSTREAM HACK                    │
                  │ func worker() string                                     │
                  │   -> name[0] = 'w' -> unicode.IsUpper('w') -> FALSE      │
                  └──────────────────────────────────────────────────────────┘
                                                │
                                                ▼
                  ┌──────────────────────────────────────────────────────────┐
                  │               TRANSFORMED COMPILER PIPELINE              │
                  │ func worker() string public { ... }                      │
                  │   -> Parser consumes optional `public` / `private`       │
                  │   -> AST stores VisPublic (1-byte enum)                  │
                  │   -> obj.Exported() queries VisPublic -> TRUE            │
                  │   -> reflectdata sets Bit 0 = 1 (Identical ABI!)         │
                  └──────────────────────────────────────────────────────────┘
```

#### Step 1: Parser Ingestion (`cmd/compile/internal/syntax/parser.go`)
Directly in `funcDeclOrNil()`, consume the optional suffix before the block `{`:
```go
// Immediately after parsing signature:
f.TParamList, f.Type = p.funcType("")

// NEW: Consume declaration-tail suffix with 1-token lookahead:
if p.got(_Public) {
    f.Visibility = VisPublic
} else if p.got(_Private) {
    f.Visibility = VisPrivate
}

// Then parse body block as usual:
if p.tok == _Lbrace {
    f.Body = p.funcBody()
}
```
Because `public` and `private` are contextual keywords parsed immediately before `{`, lookahead is exactly 1 token with zero ambiguity.

#### Step 2: AST & Symbol Representation (`go/ast`, `cmd/compile/internal/types`)
Add an explicit 2-bit visibility enum to AST nodes and symbols:
```go
type Visibility uint8

const (
    VisDefault Visibility = iota // 0: fallback to capitalization
    VisPublic                    // 1: explicitly exported
    VisPrivate                   // 2: explicitly private (file/lexical)
)
```

#### Step 3: Replace the Rune Inspection Hack in `obj.Exported()`
```go
// Transformed: Direct O(1) integer check with backwards-compatibility fallback
func (obj *object) Exported() bool {
    switch obj.vis {
    case VisPublic:
        return true
    case VisPrivate:
        return false
    default:
        // 100% Go 1 backwards compatibility fallback:
        ch, _ := utf8.DecodeRuneInString(obj.name)
        return unicode.IsUpper(ch)
    }
}
```
Replaces thousands of redundant `utf8.DecodeRuneInString` and `unicode.IsUpper` table queries with a single integer branch.

#### Step 4: Selector Disambiguation (`Id`)
Update `Id()` to key off semantic export status rather than string casing:
```go
func Id(pkg *Package, name string, exported bool) string {
    if exported {
        return name
    }
    return pkg.path + "." + name
}
```

#### Step 5: Zero Runtime ABI & Reflection Breakage
In `cmd/compile/internal/reflectdata/reflect.go`:
```go
// Replace:
nsym := dname(ft.Sym.Name, ft.Note, nil, types.IsExported(ft.Sym.Name), ...)

// With:
nsym := dname(ft.Sym.Name, ft.Note, nil, ft.Sym.Exported(), ...)
```
Because the runtime ABI already uses bit 0 of the name descriptor, this transformation requires **zero changes** to:
- `internal/abi/type.go`
- `reflect.Type` and `reflect.Value`
- DWARF debug symbols
- Garbage collector type shape descriptors

---

## 7. Ecosystem & Tooling Impact

- **`go/parser` & `go/ast`**: Add an optional `Visibility *ast.Ident` or `VisibilityToken token.Pos` to `ast.FuncDecl`, `ast.TypeSpec`, and `ast.Field`.
- **`go/types`**: During object resolution, query `obj.Exported()` via the explicit visibility specifier if present, falling back to `token.IsExported(obj.Name())`.
- **`gopls` & Autocompletion**: Autocompletion engines can offer completions based on semantic access rather than casing heuristics.
- **Reflection (`reflect`)**: Struct field reflection will respect `Field.IsExported()` using the field's explicit visibility.

---

## 8. Compatibility Guarantee

This proposal strictly satisfies the **Go 1 Compatibility Promise**:
- Every existing valid Go program remains valid with identical semantic behavior.
- `public` and `private` are contextually recognized as visibility modifiers at declaration tail positions, avoiding breaking code where these words are used as local variable names.

---

## 9. Summary Comparison

| Aspect | Go 1 (Status Quo) | Prefix Modifiers (`pub func`) | Proposed: Declaration-Suffix (`func ... public`) |
| :--- | :--- | :--- | :--- |
| **Export Mechanism** | Identifier Capitalization | Leading Keyword | Trailing Declaration Suffix |
| **Parser Disruption** | None | High (breaks statement start) | None (canonical leading keyword preserved) |
| **Acronym Ergonomics** | Poor (`json` vs `JSON`) | Good | Excellent |
| **Directory Coupling** | Relies on `internal/` | Relies on `internal/` | Purely Source-Level Encapsulation |
| **File-Level Privacy** | Not supported | Varied | Supported (`private`) |
| **Runtime ABI Impact** | Baseline (Bit 0) | May alter name symbols | Zero changes (Bit 0 preserved) |
| **Go Minimalism** | High | Low (introduces new syntax blocks) | High (only `public` and `private`) |
| **Go 1 Compatible** | Baseline | Breaking or complex lookahead | 100% Backwards-Compatible |
