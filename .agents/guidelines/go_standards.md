# Go Standards Family Registry

This registry serves as the entry point for all Go programming conventions, optimization guidelines, and compiler compliance rules. Refer to the specific sub-guideline matching your task profile:

---

## 1. Concurrency & Memory Safety
* **Trigger:** Writing or modifying Go files containing synchronization primitives (`sync.Mutex`, `sync.RWMutex`, `sync.Map`, `chan`, `go`).
* **Guideline:** [go_concurrency.md](file:///c:/aCogSpaceSeed/.agents/guidelines/go_concurrency.md)

---

## 2. AST Compliance & Remediations (Shift-Left)
* **Trigger:** Generating new code or refactoring structures to conform with the `goenhancer` AST rules.
* **Guideline:** [goenhancer_preventative.md](file:///c:/aCogSpaceSeed/.agents/guidelines/goenhancer_preventative.md)

---

## 3. Maturity, Dead-Stores & Structured Logging
* **Trigger:** Setting up logging standards, error handling chains, sentinel errors, or process exit/panic guards.
* **Guideline:** [go_maturity_standards.md](file:///C:/aCogSpaceSeed/000ALL/00000-knowledge-foundations/120-system-architecture/language-maturity/go_maturity_standards.md)

---

## 4. Modern Go 1.26.4+ Best Practices Standard
* **Trigger:** Writing, refactoring, or compiling Go code across sovereign fleet workspaces.
* **Mandates:**
  1. **WebNF & WebNF-Materialized JSON Telemetry Standard:** All logging, benchmark recording, and execution tracing MUST use **native `.webnf` DSL program records** OR **JSON materialized directly from `.webnf` records** via `log/slog`. Unstructured text logging (`fmt.Println`, `log.Println`) is strictly prohibited.
  2. **Zero-Allocation String Builders:** Use `strconv.Itoa`, `strconv.FormatFloat`, `strconv.FormatUint`, or `strings.Builder` with `Grow()` capacity pre-allocation. Prohibit `fmt.Sprintf` in hot loops.
  3. **Go 1.24+/1.26+ `omitzero` Tag Annotations:** Use `json:"field,omitzero"` for WebNF-materialized struct field marshaling.
  4. **Go 1.24+/1.26+ `weak.Pointer` Zero-Leak Caching:** Use `weak.Make(ptr)` and `weakRef.Value()` for temporary caches.
  5. **Pure Go / Zero CGo:** Avoid CGo to eliminate ~80ns context-switching overhead and ensure cross-platform portability.
