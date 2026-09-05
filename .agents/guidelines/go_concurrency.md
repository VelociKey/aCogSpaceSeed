# Guideline: Go Concurrency & Lock-Free Design Rails

When writing, scaffolding, or modifying Go code that contains synchronization primitives (`sync.Mutex`, `sync.RWMutex`, `sync.Map`, `sync.Once`), channels (`chan`), or asynchronous goroutines (`go func()`):

5: 1. **Invoke the Concurrency Advisor Skill:** 
6:    Proactively run the `recommendConcurrency` skill on the target workspace to analyze the access patterns, lock contention, and deadlock risks.
2. **Grammar-Driven Concurrency Contracts:**
   Do NOT implement raw, ad-hoc locks directly in the Go source. First declare all concurrent states (locks, channels, lookup maps) inside the local `.concurrency.webnf` contract file.
3. **Compile & Generate Context:**
   Run the `applyconcurrency -gen -workspace <path>` compiler CLI to generate type-safe wrapper structs (`concurrency_gen.go`). Implement the Go code against these generated, optimized structures rather than standard lock primitives.
4. **Enforce Static Verification:**
   Before completing any task, execute the validator using `applyconcurrency -validate -workspace <path>` to statically verify:
   * **Rule V10 (LOG cycle checks):** Proof of zero deadlock cycles.
   * **Rule V12 (Unbounded Select):** Ensuring all selects have safety timeouts or context Done receivers.
