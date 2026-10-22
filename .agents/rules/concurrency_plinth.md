# Sovereign Rule: Concurrency, Memory Arenas & Substrate Contracts

> **Scope:** Core storage engines, low-level data substrates, memory arenas, and concurrent actors.  
> **Primary Workspaces:** `00ulib/substrate`, `00flow/plinth`, `00flow/adaptivelock`.

---

<hard_invariants>
## 1. Zero-Allocation & Substrate Contracts
1. **TWO-STAGE HARDWARE & SUBSTRATE CONTRACTS:** When creating or modernizing core storage logic, agents must formally separate implementation into:
   - **Stage A:** Interface contract + In-Memory Test Shim (enables deterministic, dependency-free unit testing).
   - **Stage B:** Physical Plinth `FileBackedAddressSpace` memory-mapped extent wiring (delivers zero-copy persistence).
2. **CONCURRENT TEST ADAPTER SYNCHRONIZATION:** When mutating an interface or constructor signature, agents MUST simultaneously provide backwards-compatible adapters (e.g. `OSFileWALWriter`) in existing `_test.go` suites. Agents are forbidden from leaving test suites broken while progressing to subsequent phases.
3. **PRE-FLIGHT AST COMPLIANCE GATE:** Before concluding synthesis, agents must run AST conformance audits (`vkey inspect slog-check`, `vkey inspect audit-imports`) to purge residual struct tags, unused imports, or non-compliant print calls.
4. **NATURAL CACHELINE ALIGNMENT:** Rely on linear indexing over raw Plinth memory slices (`as.Bytes()`) and 64-bit natural struct ordering. Never introduce artificial `_ [64]byte` struct pads.
</hard_invariants>

---

<verification_protocol>
## 2. Escape Analysis & Allocation Tooling
- Verify zero steady-state heap allocations on hot paths:
  `go build -gcflags="-m"` or `realize.exe -workspace <name> -test`
- Maintain $0\text{ B/op}$ memory profiles in all substrate codec loops.
</verification_protocol>
