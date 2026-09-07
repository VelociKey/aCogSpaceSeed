# Active Task Goal: SYS- Area Benchmark Governance & Maturity Promotion Gates

**Task ID:** `00xper-TASK-SYS-BENCHMARK-GOVERNANCE`  
**Sponsor:** Conductor / NATVS Engine  
**Status:** ACTIVE  
**Target Workspaces:** `00xper/ms-benchmarks`, `00xper/ms-warehouse`

---

## 🎯 Global Goal Mandate

All SYS- area benchmark execution runs that provision GCP Spot VMs (cost range $0.98–$32.00/hr) must follow the NATVS 5-phase lifecycle. Maturity promotions (M5→M6) must be gated behind automated validation of B5 hardware tile speedup evidence. No hardcoded promotions.

---

## 📋 Directives & Invariant Constraints

1. **Budget Authorization Invariant (Negotiation):** Before any Spot VM is provisioned for a benchmark run, the estimated cost ceiling must be calculated and logged. Runs exceeding $50 cumulative must require explicit user confirmation via the Non-Spot Governor double-permission gate pattern.
2. **Execution Ledger Invariant (Synthesis):** Every benchmark run must produce a WebNF ledger record appended to the workstation ledger or streamed to GCS via the QUIC UDP streamer. Runs without ledger records are considered void.
3. **Promotion Gate Invariant:** M5→M6 promotion requires:
   - $B_0 \to B_5 \ge 10^5\times$ measured speedup (not projected)
   - B5 `Status` == `"MEASURED"` or `"VERIFIED"` (not `"PROJECTED"`)
   - Export compliance flag is `GREEN`
   - System screener compatibility result is `COMPATIBLE` for at least 1 GCP chip target
4. **Maturity Signal Invariant:** When a benchmark run validates a new B-stage for a SYS- area, the corresponding maturity advancement signal must be emitted. The ms-warehouse `AdvanceMaturity()` function must be called to update the fleet maturity state.

---

## 🚀 Active Roadmap

### Phase 1: Promotion Gate Implementation ✅
- [x] Create `promotion_gate.go` in `ms-benchmarks/81000-active-source/runner/`
- [x] Implement `ValidatePromotionGate(rec SystemBenchmarkRecord) PromotionResult`
- [x] Enforce speedup threshold, measurement status, and export compliance checks
- [x] Add unit tests for promotion gate validation

### Phase 2: Maturity Signal Derivation ✅
- [x] Create `maturity_signal.go` in `ms-benchmarks/81000-active-source/runner/`
- [x] Implement `DeriveMaturityFromBenchmark(rec SystemBenchmarkRecord) MaturitySignal`
- [x] Map B-stage completion to M-stage advancement
- [x] Add unit tests for maturity signal derivation

### Phase 3: SuperPrimitive Coverage Expansion ✅
- [x] Bind 5 unbound SuperPrimitives to anchor SYS- areas in `ms-warehouse`
- [x] Verify cross-domain discovery reports full 10/10 coverage

### Phase 4: Performance Optimization ✅
- [x] Rewrite `ToWebNFFormat()` to use `strconv.Append*` instead of `fmt.Fprintf`
- [x] Eliminate reflection overhead and reduce GC allocations

### Phase 5: End-to-End Integration (FUTURE)
- [ ] Wire `ms-benchmarks` maturity signals to `ms-warehouse` `AdvanceMaturity()` in a shared orchestrator
- [ ] Add conductor cron job to reconcile maturity state weekly
- [ ] Integrate benchmark run cost tracking into fleet compendium

---

## NATVS Protocol Mapping

| NATVS Phase | Benchmark Governance Activity |
|:--|:--|
| **N (Negotiation)** | Budget estimation, chip target selection, spot price threshold check |
| **A (Assimilation)** | Ingest SYS- area specs from ms-warehouse, screen hardware compatibility |
| **T (Transformation)** | Execute B0→B5 progression benchmark suite on target chip |
| **V (Verification)** | Validate promotion gate, derive maturity signal, append WebNF ledger |
| **S (Synthesis)** | Stream records to GCS, advance maturity state, reconcile fleet scoreboard |
