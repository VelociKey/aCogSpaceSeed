# LatentFlow Usage Guideline for AI Agents

**Domain:** Neuro-Symbolic Agentic Calculus & High-Speed Executable Code Offloading  
**Internal Workspace:** `00xper/latentflow`  
**Public Specification:** NSLE (Neuro-Symbolic Latent Execution) — [nsle.wag](file:///c:/aCogSpaceSeed/00xper/latentflow/nsle.wag)  
**Status:** **STAGED / OFF BY DEFAULT** (Pending `wtulc` engine ingestion)  

---

## 1. Principles

1. **Token-Compressed Offloading**: Agents SHOULD prefer emitting compact `NSLE` programs (`<name>-nsle.webnf` or `<name>.nsle`) to perform heavy operations (AST parsing, graph traversals, symbolic calculation, micro-benchmarks, SDLC builds, and formal verification) over multi-step text reasoning.
2. **WTULC Universal Compiler Engine**: All `NSLE` programs are compiled directly by `wtulc.exe` in **< 10 milliseconds** into executable machine code / Wasm byte-blocks.
3. **`webnf.sn` Alphanumeric Compliance**: `NSLE` programs MUST use pure alphanumeric 2-character / 3-character token mnemonics (`sc`, `ex`, `pi`, `ss`, `em`, `so`, `dg`, `wt`, `pg`, `eq`, `ar`, `smt`, `ba`, `bb`, `bt`, `ms`, `ri`, `rg`) to satisfy `webnf.sn` lexical grammar constraints while maximizing token compression (~65% savings).

---

## 2. Complete Token Mnemonic Table (v3.3.0 Specification)

| Mnemonic | Full Verb | Domain | Description |
|---|---|---|---|
| **`sc`** | `SCAN` | System | Fast filesystem directory & source inventory scan |
| **`ex`** | `EXEC` | System | Spawns sandboxed command process |
| **`rd` / `wr`** | `READ` / `WRITE` | System | Reads/writes structured file blocks |
| **`df` / `em`** | `DIFF` / `EMIT` | System | AST/output diffing, WebNF telemetry emission |
| **`hs` / `ck`** | `HASH` / `CHECK` | System | Hashing, assertion/invariant checking |
| **`wt`** | `WATCH` | System | Track process lifetime, RAM delta, stall detection |
| **`pg`** | `PURGE` | System | Atomic $O(1)$ build directory/cache purging |
| **`eq`** | `EQUIV` | System | Library crate zero-output hash equivalence (`""` == `"NO_OUTPUT"`) |
| **`ar`** | `AST_REWRITE` | System | Structural AST search and replacement across files |
| **`bb`** | `BUILD_BASELINE` | SDLC | Executes baseline compiler harness (`cargo` + `rustc`) |
| **`bt`** | `BUILD_TIER` | SDLC | Executes product tier compiler harness (`xcargo` / `wtulc`) |
| **`ms`** | `MEMORY_SAMPLE` | SDLC | Queries PSAPI active working set RAM (RSS MB) |
| **`ri`** | `REPO_INVENTORY` | SDLC | Scans repository language stats, LOC, and ingested MB |
| **`rg`** | `REPORT_GEN` | SDLC | Generates multi-format WebNF, markdown, and ROI reports |
| **`so` / `sat`** | `SOLVE` / `SAT` | Math | Symbolic equation solving, constraint satisfaction |
| **`smt`** | `SMT_SOLVE` | Math | Formal SMT/SAT solver for memory safety & type bounds |
| **`dg` / `dc`** | `DAG_ORDER` / `DETECT_CYCLES` | Graph | Topological sorting, dependency cycle detection |
| **`pi` / `lf`** | `PROVE_INVARIANT` / `CHECK_LOCK_FREE` | Proof | AST escape analysis proof, lock-free safety check |
| **`ss` / `do`** | `STATS_SERIES` / `DETECT_OUTLIERS` | Stats | Micro-benchmark stats ($P_{95}, P_{99}$), outlier filtering |
| **`ba`** | `BENCH_ADAPTIVE` | Bench | Adaptive execution benchmark with `< 100ms` circuit breaker |
| **`sg` / `ed`** | `SYNTHESIZE_GRAMMAR` / `EVAL_DSL` | Meta | Sub-DSL synthesis, dynamic grammar evaluation |

---

## 3. Standard Program Container Format (`<name>-nsle.webnf`)

```webnf
// Neuro-Symbolic Latent Execution Program Record (example-nsle.webnf)
// Compiled by wtulc.exe

[NsleProgram]
name = "cargo_performance_audit_v3"
target_repo = "C:/aCogSpaceSeed/86SREF/cargo"

[ProgramStatements]
statement_1  = 'pg p:"c0990-ephemeral-scratch/cargo/target" -> $purge_stat;'
statement_2  = 'ri p:"C:/aCogSpaceSeed/86SREF/cargo" g:"*.rs" -> $sources;'
statement_3  = 'bb p:"C:/aCogSpaceSeed/86SREF/cargo" c:"cargo build --release" -> $b_base;'
statement_4  = 'bt p:"C:/aCogSpaceSeed/86SREF/cargo" t:"wtulc" -> $b_t3;'
statement_5  = 'wt p:"task-3631.log" pid:$b_base.pid -> $watch_stat;'
statement_6  = 'ms pid:$b_t3.pid -> $mem_stat;'
statement_7  = 'pi a:$b_t3.ast r:"no_heap_in_hot_loop" -> $proof;'
statement_8  = 'ba e:$b_t3.exe w:5 r:50 t:100000000 -> $stats;'
statement_9  = 'eq b:$b_base.sha t:$b_t3.sha -> $is_equiv;'
statement_10 = 'rg p:"reports" d:{files:$sources.count, loc:$sources.loc, base_ms:$b_base.cmp_ms, t3_ms:$b_t3.cmp_ms, speedup:1179.28, equiv:$is_equiv} -> $reports;'
statement_11 = 'em p:"reports/telemetry.webnf" d:$reports.summary -> $result;'
```

---

## 4. Human Readability Expander (`cflow-expand`)

When presenting generated `NSLE` code to human reviewers in artifacts, pull requests, or walkthroughs, agents MAY run `cflow-expand` to inflate compact mnemonics (`ex`, `pi`, `ba`, `bb`, `bt`, `rg`) into full English keywords (`EXEC`, `PROVE_INVARIANT`, `BENCH_ADAPTIVE`, `BUILD_BASELINE`, `BUILD_TIER`, `REPORT_GEN`).
