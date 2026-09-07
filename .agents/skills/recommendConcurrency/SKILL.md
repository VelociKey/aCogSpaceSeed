---
name: recommendConcurrency
description: Analyzes Go workspaces or files for concurrency issues, lock contention, cacheline false-sharing, and provides architectural lock-free recommendations based on applyconcurrency patterns.
user_invocable: true
arguments:
  - name: targets
    description: Comma-separated or space-separated list of workspace names or directories to analyze (e.g. 'hydration, adaptivelock')
    required: true
  - name: files
    description: Optional specific list of Go file paths to restrict analysis to
    required: false
  - name: interactive
    description: Set to true if workspaces interact and should be analyzed for cross-workspace locking dependencies
    required: false
    default: "false"
---

# Concurrency Recommendations Skill: High-Performance Architectural Advisor

Use this skill to audit Go code across the sovereign fleet for concurrency correctness, lock contention, cacheline false-sharing, and architectural patterns (lock-free rings, thread-local arenas, adaptive locks, channel topologies). The analysis leverages the `applyconcurrency` selection rules, `adaptivelock` fit criteria, and the **`AGENTS.md` Rule #6 Double-Performance Standard**.

---

## Usage Patterns

- `/recommendConcurrency --targets evolution-scaffold`
- `/recommendConcurrency --targets "usegrammar, latentlingua" --interactive true`
- `/recommendConcurrency --targets sacp --files "sacp_session.go, sacp_mux.go"`

---

## Sovereign Concurrency Invariant Rules

| Invariant Rule | Name | Check & Violation Threshold | Remediation |
|---|---|---|---|
| **Rule V10** | **Lock-Order Graph (LOG) Cycles** | Nested mutex acquisitions ($A \to B \to A$) that create static deadlock cycles. | Enforce flat lock hierarchy ($D_{\text{lock}} = 1$) or atomic state machine. |
| **Rule V12** | **Unbounded Select Blocking** | Channel `select` blocks without a `default:`, `case <-ctx.Done():`, or `case <-time.After():`. | Add drop-safe `default:` clause or cancellation context. |
| **Rule V14** | **Cacheline False Sharing** | Contiguous atomic counters (`atomic.Uint64`) or mutexes sharing the same 64-byte L1 line. | Add 64-byte isolation padding (`[64]byte` pads) between hot counters. |
| **Rule V16** | **Locks over Blocking I/O** | Holding a mutex across disk reads/writes, network sockets, or subprocess calls. | Snapshot data under lock, release lock before performing I/O. |
| **Rule V18** | **Hot-Path Allocations** | Memory allocations inside concurrency loops or channel transfers exceeding $0\text{ B/op}$. | Use contiguous byte arenas (`ADPHByteArena`) or pre-allocated ring buffers. |

---

## Primitive Selection & Architecture Matrix

Assess the workload pattern for each concurrency site and map it to the authoritative `applyconcurrency` taxonomy:

| Workload Pattern | Data vs Control Plane | Access Characteristics | Recommended Sovereign Primitive |
| :--- | :---: | :--- | :--- |
| **Ultra High-Frequency Streaming** | Data Plane | $>10^6\text{ ops/sec}$, latency $<10\text{ ns}$ | **Lock-Free SPSC/MPSC Ring Buffer** (`ShadowRingBuffer` + 64B padding) |
| **Worker Local Accumulation** | Data Plane | Thread-local symbol/AST building | **Thread-Local Byte Arena** (`ADPHByteArena`, zero locks, zero GC) |
| **Batch Append Ingestion** | Data Plane | Multiple workers append, read once | **Lock-Free FanInCollector** (`chan T` with worker fan-in) |
| **Read-Heavy Global State** | Control Plane | $>99\%$ reads, session lifetime, $<200$ items | **`AdaptiveMap`** (`adaptivelock.Map`) |
| **Request-Scoped Lookup** | Control Plane | Request-lifetime, temporary map | **`MutexMap`** (`sync.Mutex` map) |
| **One-Shot Gate / Breaker** | Control Plane | Initialized once, checked continuously | **`AtomicFlag`** (`atomic.Bool` / `atomic.Uint32`) |
| **Fine-Grained Partitioning** | Control Plane | Locking distinct items independently | **`PerResourceMutex`** / Sharded Lock Ring |
| **State Mutation Serializer** | Control Plane | High-contention map updates | **`CommandQueue`** (Actor channel serialization) |

---

## Execution Protocol for the Agent

When this skill is triggered, perform the following 5-phase analysis:

### 1. Code Scan & Discovery
* Scan target workspaces/files using ripgrep (`grep_search`) or AST analysis tools.
* Identify all instances of:
  * Mutexes: `sync.Mutex`, `sync.RWMutex`
  * Atomics: `sync/atomic` types (`atomic.Uint64`, `atomic.Bool`, `atomic.Pointer`)
  * Channels: `chan`, `select`, `<-`, `close`
  * Goroutines: `go ` invocations

### 2. Deadlock & False-Sharing Auditing
* Map out nested lock acquisitions (Rule V10).
* Verify timeout and default guards on all `select` blocks (Rule V12).
* Audit struct field layouts for adjacent atomics without 64-byte cacheline isolation (Rule V14).
* Ensure zero locks are held across blocking I/O (Rule V16).

### 3. Lock-Free & Double-Performance Opportunities
* Proactively suggest replacing mutex-guarded slices with `FanInCollector` or `ShadowRingBuffer`.
* Suggest thread-local arena buffers (`ADPHByteArena`) to eliminate heap allocations ($0\text{ B/op}$).

### 4. Deliverable: Markdown Concurrency Audit Report
Produce a structured markdown audit report in the artifact directory (`concurrency_analysis.md`) containing:
* **Findings Matrix:** Scanned files, lock sites, primitives, and workload classification.
* **Invariant Compliance Table:** V10, V12, V14, V16, and V18 verification status.
* **Refactoring Blueprint:** Concrete code diffs showing how to refactor locks into lock-free or adaptive equivalents.

### 5. Deliverable: Declarative `.concurrency.webnf` Contract
Emit a formal `.concurrency.webnf` specification for the workspace:

```webnf
; ==========================================================================
; SOVEREIGN CONCURRENCY DECLARATION CONTRACT (weBNF)
; ==========================================================================

CONTRACT_ID          = "CONTRACT-CONCURRENCY-AUTODETECT-001" ;
APPLY_CONCURRENCY    = "ENABLED" ;

; MEMORY & CACHELINE POLICIES:
ISOLATION_PADDING    = "64_BYTE_CACHELINE_PADS" ;
DATA_PLANE_STREAM    = "LOCK_FREE_SPSC_OR_ARENA" ;
CONTROL_PLANE_CHAN   = "BUFFERED_SELECT_WITH_GUARD" ;

; INVARIANT VERIFICATION:
MAX_LOCK_DEPTH       = 1 ;
ALLOW_LOCKS_OVER_IO  = false ;
SELECT_TIMEOUT_GUARD = true ;
STATUS               = "CERTIFIED_SOVEREIGN_CONCURRENCY" ;
```
