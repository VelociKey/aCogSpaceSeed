---
name: adaptive-orchestrator
description: Sovereign Adaptive Orchestrator & Definition-of-Need (.wag / .webnf) Skill. Transforms base unitary computations (nanosecond SIMD kernels like UUID or millisecond service/actor workflows like construct preserve) into 4-regime adaptive pipelines with flat SoA slabs, slot-indexed multi-caller returns, priority lanes, and chunk-boundary cost/deadline interrupts using the cadence-orchestrator factory.
---

# Sovereign Skill: `adaptive-orchestrator`

> **Factory Workspace:** [`00xper/corealgorithms/cadence-orchestrator`](file:///c:/aCogSpaceSeed/00xper/corealgorithms/cadence-orchestrator)  
> **Formal Grammar (`weBNF.sn` / `.wag`):** [`adaptive_orchestrator.wag`](file:///c:/aCogSpaceSeed/00xper/corealgorithms/cadence-orchestrator/81000-active-source/adaptive_orchestrator.wag) (`:com:velocikey:orchestration:adaptive`)  
> **Reference Implementation:** [`00xper/corealgorithms/uuid`](file:///c:/aCogSpaceSeed/00xper/corealgorithms/uuid)

---

## 1. When to Activate This Skill

Activate `adaptive-orchestrator` whenever a user or workflow needs to:
1. Take a **base unitary computation** (deterministic math/hash/SIMD kernel OR service/actor/I/O call such as `vkey construct preserve`) and scale it across **single-item direct calls, preallocated batch slabs, concurrent callers, and multi-core worker pools**.
2. Compute optimal inline-vs-parallel thresholds, chunk sizes, and ring capacities from **Unitary Execution Time ($T_{\text{unit}}$)** and **Response-Time Budget ($B_{\text{SLA}}$)**.
3. Orchestrate workloads requiring **Priority Bands (`critical`, `high`, `normal`, `background`)**, **Interruptibility by Cost Budget / SLA Deadline / Priority Preemption**, or **Leader-Coalesced Deduplication** across concurrent callers.
4. Generate **Structure-of-Arrays (SoA) preallocated slabs**, **Slot-Indexed Multi-Caller Return Rings (`CallerIDs[slot]`)**, and **B-Tree-Free Sentinel-Remapped Flat Tables** that comply 100% with the **Nautilus Three-Border Compiler Optimization Principle** (@c:\aCogSpaceSeed\.agents\rules\compiler_nautilus.md).

---

## 2. The 4-Step Engineering Workflow

### Step 1: Write & Benchmark the Base Unitary Computation
1. Author the clean, single-item unitary function first (e.g., `ComputeTwoStageUnit` or `PreserveWorkspaceUnit`).
2. Flatten any avoidable heap allocations in the unitary path (use stack buffers `[128]byte`, lookup tables `[256]uint32`, or split variable-length keys into fixed `(uint64, uint64)` stages).
3. Measure or establish two core timing numbers:
   - **`unitary_ns` ($T_{\text{unit}}$):** Time in nanoseconds to execute 1 item (e.g., `36 ns` for Two-Stage UUID; `50,000,000 ns` (`50 ms`) for `vkey-git` workspace preserve).
   - **`response_budget_ns` ($B_{\text{SLA}}$):** Target maximum response latency budget in nanoseconds (e.g., `2,000 ns` for in-memory token minting; `5,000,000,000 ns` (`5 s`) for CLI batch preservation).

### Step 2: Author the Definition of Need (`.webnf` / `.wag`)
Express the workload's input/output schema, timing budget, and orchestration policy using the `:com:velocikey:orchestration:adaptive` grammar defined in [`adaptive_orchestrator.wag`](file:///c:/aCogSpaceSeed/00xper/corealgorithms/cadence-orchestrator/81000-active-source/adaptive_orchestrator.wag):

```webnf
:com:velocikey:orchestration:adaptive

OrchestrationNeed <DomainName> {
    target {
        package = "<pkg>" ;
        unitary_func = "<UnitaryFuncName>" ;
        workload_kind = "deterministic_simd" ; # or "bounded_service_io" / "hybrid_pipeline"
    }
    pace_budget {
        unitary_ns = 36 ;
        response_budget_ns = 2000 ;
        handoff_ns = 1200 ;
        simd_speedup_x100 = 200 ;
        max_workers = 8 ;
    }
    policy {
        orchestrator_mode = "cadence" ; # "cadence" | "priority" | "interruptible" | "leader_coalesced" | "composite"
        priority_bands = [ "critical", "high", "normal", "background" ] ;
        interrupt_on = [ "cost_budget", "deadline_sla", "priority_preempt" ] ;
        interrupt_granularity = "chunk_boundary" ;
        max_cost_budget = 4096 ;
        unit_cost = 1 ;
    }
    slab_schema {
        inputs = [ DomainID : "uint64", KeySeed : "uint64", Tag : "uint8" ] ;
        outputs = [ Stage1 : "uint64", Stage2 : "uint64", Checksum16 : "uint16" ] ;
        caller_return = "slot_indexed" ;
    }
    flat_storage {
        enabled = true ;
        primary_key = "Stage1" ;
        verify_key = "Stage2" ;
        sentinel_remap = "zero_empty_one_tombstone" ;
    }
}
```

### Step 3: Materialize via `cadence-orchestrator` Factory
Use [`cadenceorchestrator.ParseNeedDSL`](file:///c:/aCogSpaceSeed/00xper/corealgorithms/cadence-orchestrator/81000-active-source/need_dsl.go) and [`cadenceorchestrator.SynthesizeAdaptiveOrchestratorGo`](file:///c:/aCogSpaceSeed/00xper/corealgorithms/cadence-orchestrator/81000-active-source/code_factory.go) (or bind directly to [`cadenceorchestrator.NewCadenceOrchestrator[SlabT]`](file:///c:/aCogSpaceSeed/00xper/corealgorithms/cadence-orchestrator/81000-active-source/cadence_orchestrator.go)):

1. **Mathematical Regime Tuning (`spec.ComputePacePlan()`):**
   - **Fast Deterministic SIMD (`unitary_ns < handoff_ns`):**
     $$N_{\text{inline}} = \text{align}_{16}\left(\left\lceil \frac{T_{\text{handoff}}}{T_{\text{SIMD}}} \right\rceil\right)$$
     Batches $\le N_{\text{inline}}$ stay 100% single-threaded inline (`0 B/op`, zero channel overhead).
   - **Service / Actor / I/O (`unitary_ns >= handoff_ns`, e.g., `vkey construct preserve`):**
     $N_{\text{inline}} = 1, \text{ChunkSize} = 1$. A single workspace runs directly inline; $\ge 2$ workspaces immediately fan out across the bounded worker pool!
2. **Nautilus-Clean Linear Range Kernel (`RangeKernel[SlabT]`):**
   - Always write `for i := start; i < end; i++` over contiguous SoA slices (`slab.<Col>In[i]`, `slab.<Col>Out[i]`).
   - **NEVER** hardcode SIMD lane counts (`for n-i >= 8`), **NEVER** manually unroll loops, and **NEVER** put per-item interrupt/cost branches inside the `for i := start; i < end; i++` loop.
3. **Chunk-Boundary Interrupts & Priority Preemption:**
   - Cost budgets (`MaxCostBudget`), SLA deadlines (`ResponseBudgetNs`), cancellation flags (`CancelFlag`), and higher-priority preemption (`YieldOnHigherPriority`) are evaluated **strictly between `[start, end)` chunks** in `ExecuteInterruptible`.
4. **Zero-Routing Multi-Caller Slot Returns:**
   - Concurrent callers claim slot $i$ via `slab.ClaimCallerSlot(callerID, priority, ...)` and read their own output directly from `slab.ReadCallerResult(slot, callerID, epoch)` in $O(1)$ without per-caller channels or maps.

### Step 4: Verify with a 4-Regime Attestation Benchmark
Always verify the generated orchestrator across all 4 regimes using `realize.exe`:
- **Regime 1 (Direct Unitary):** `N = 1` inline latency (`0 B/op`).
- **Regime 2 (Preallocated SoA Slab):** Medium batch (`N <= BaseInlineThreshold`) single-core SIMD linear sweep (`0 B/op`).
- **Regime 3 (Adaptive Multi-Core Fanout):** Large batch (`N >> BaseInlineThreshold`) chunked across preallocated worker tickets (`0 B/op`).
- **Regime 4 (Multi-Caller Coalesced Ring / Leader Coalescer):** Concurrent callers claiming slots in a shared ring slab or deduplicating via `LeaderCoalescer`.
