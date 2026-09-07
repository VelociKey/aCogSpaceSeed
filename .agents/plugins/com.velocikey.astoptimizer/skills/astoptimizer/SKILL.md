---
name: com.velocikey.astoptimizer
displayName: vkey-optimizer
description: >-
  Automated Go and Dart double-performance refactoring engine powered by Qwen-2.5-Coder-32B
  and Codestral-22B. Enforces stack allocation, arena buffer pooling, lock-free atomics,
  and GC allocation reduction >= 50% with verified pre- and post-benchmarking.
user_invocable: true
arguments:
  - name: workspace
    description: Target workspace name or directory path (e.g. 'tensorcore', 'transport')
    required: true
  - name: benchmark
    description: Run automated before/after benchmark comparison to prove speedup
    required: false
    default: "true"
---

# vkey℠ AST Optimizer: Double-Performance Standard Refactoring

`vkey-optimizer` enforces the **Double-Performance Standard** (AGENTS.md Section 6): reducing GC allocations by $\ge 50\%$ and doubling throughput ($2\times$) via zero-allocation Go/Dart AST transformations.

---

## Injected Slash Commands

- `/vkey-optimizer --workspace <name>`: Optimize all packages in a workspace.
- `/vkey-optimizer-bench <pkg_path>`: Run before/after escape analysis and benchmark comparison.

---

## Agent Operational Protocol

1. **Pre-Optimization Benchmark Baseline:**
   * Run `go test -bench=. -benchmem ./...` inside the target workspace and record `ns/op` and `B/op`.

2. **AST Transformation & Infill:**
   * Invoke `qwen-2.5-coder-32b` and `codestral-22b` (FIM) to replace:
     * `fmt.Sprintf` hex formatting $\to$ byte bit-shifting.
     * `sync.Mutex` contention hotspots $\to$ lock-free `sync/atomic`.
     * Intermediate string-to-byte casting $\to$ zero-copy `unsafe.StringData` / arena pools.

3. **Post-Optimization Verification & Nontrivial Audit:**
   * Re-run `go test -bench=. -benchmem` to mathematically verify the speedup delta.
   * Run `nontrivial.exe scan` to certify zero memory safety vulnerabilities.
