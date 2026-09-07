# Task: Afflume Three-Party System & Plinth GCS Loader Implementation Plan

* **Status:** Ready for Execution (Post-Restart Resume Point)
* **Target Silo:** `00floo`
* **Primary Workspace:** [`00floo/afflume`](file:///c:/aCogSpaceSeed/00floo/afflume)
* **Dependent Workspaces:** [`00floo/afflumeattestation`](file:///c:/aCogSpaceSeed/00floo/afflumeattestation), [`00flow/adph`](file:///c:/aCogSpaceSeed/00flow/adph), [`00flow/financialkernel`](file:///c:/aCogSpaceSeed/00flow/financialkernel)

---

## 1. Architectural Blueprint: The Three-Party System + Plinth Loader

```
┌────────────────────────────────────────────────────────────────────────────────────────┐
│                        CLOUD RUN v2 CONTAINER INSTANCE                                 │
│                                                                                        │
│  [PLINTH GCS LOADER] (Phase 3)                                                         │
│  1. Boots container with CLOUD_RUN_TASK_INDEX                                          │
│  2. Streams `cohort_{task_index}.soa.bin` from same-region GCS ($0.00 egress)          │
│  3. Verifies White3 / Blake3 cryptoseal                                                │
│  4. Maps raw bytes directly into Contiguous SoA RAM                                    │
│                                                                                        │
│  ┌──────────────────────────────────────────────────────────────────────────────────┐  │
│  │ PARTY 3: CLIENT AUTH & GATE SENTRY (< 10 ns) (Phase 2)                           │  │
│  │ • Client probes: CustomerUUID, ADPH Table Version                                │  │
│  │ • Branchless read: ServiceGatesOpen[c] & MacroStages[c]                          │  │
│  │ • Response: Proceed / Denied (OverdraftCapped / DisputeLock) + ADPH Handshake     │  │
│  └──────────────────────────────────────────────────────────────────────────────────┘  │
│                                                                                        │
│  ┌──────────────────────────────────────────────────────────────────────────────────┐  │
│  │ PARTY 1: IMMUTABLE EVENT LOG (Dual Ledger A/B) (Phase 1)                         │  │
│  │ • Raw Wire Append: Ledger A                                                      │  │
│  │ • Certified Transmutation: Ledger B (White3 Sealed)                              │  │
│  │ • Talos Binary Transmuter: Ingests []ADPHSlotCount batches                       │  │
│  └────────────────────────┬─────────────────────────────────────────────────────────┘  │
│                           │                                                            │
│  ┌────────────────────────▼─────────────────────────────────────────────────────────┐  │
│  │ PARTY 2: FINANCIAL LIFECYCLE & VPU ARITHMETIC (Already Implemented Core)         │  │
│  │ • Paying Event: Adds (+) VPU balance to Customer (Stripe / ACH / Wire settled)   │  │
│  │ • Talos Event: Evaluates Hybrid Rates & subtracts (-) consumed VPUs              │  │
│  │ • CoalescedLifecyclePass: Single L1 sweep fuses debt, taxes, grace, and gates    │  │
│  └──────────────────────────────────────────────────────────────────────────────────┘  │
│                                                                                        │
│  [PLINTH CHECKPOINT & PERSIST]                                                         │
│  5. Plinth writes updated `cohort_{task_index}.soa.bin` to GCS (If-Generation-Match)   │
│  6. Plinth whispers M31 state digest across Plinth-to-Plinth Gossip Bus                │
└────────────────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Phase-by-Phase Execution Plan

### Phase 1: Party 1 Talos Non-String Telemetry Ingestion Bridge
* **Goal:** Connect inbound non-string ADPH telemetry from running customer workloads directly into Ledger A/B and pipe to Party 2 for VPU consumption deduction.
* **Files to Create / Modify:**
  1. [`00floo/afflume/81000-active-source/eventconduit/types.go`](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/eventconduit/types.go):
     - Register `SourceTalosMetric InboundSource = "adapter:talos-metric-reporting"`.
     - Register `EventTalosUsageReported EventType = "EVENT_TALOS_USAGE_REPORTED"`.
  2. [`00floo/afflume/81000-active-source/eventconduit/talos_binary_transmuter.go`](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/eventconduit/talos_binary_transmuter.go):
     - Implements `EventTransmuter` for `SourceTalosMetric`.
     - Ingests packed 8-byte `[]ADPHSlotCount` wire stream.
     - Commits raw frame to Ledger A (append-only wire log) in $< 10\text{ ns}$.
     - Calls `matrix.EvaluateHybridConsumption(&contract, epochRegistry, batch)`.
     - Commits certified VPU consumption event to Ledger B.
     - Debits customer VPU balance.
  3. [`00floo/afflume/81000-active-source/eventconduit/talos_binary_transmuter_test.go`](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/eventconduit/talos_binary_transmuter_test.go):
     - Unit test verifying raw wire logging, rate resolution, bonus VPU deduction, and net VPU debit.

---

### Phase 2: Party 3 Client Authentication & Gate Sentry Endpoint
* **Goal:** Implement the ultra-fast gate sentry procedure call that clients query to authenticate and check service gate authorization before running jobs.
* **Files to Create / Modify:**
  1. [`00floo/afflume/81000-active-source/sentry/gate_sentry.go`](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/sentry/gate_sentry.go):
     - `GateSentry`: High-throughput, lock-free sentry actor.
     - Holds atomic pointer to active `FleetFinancialLifecycleMatrix` and `VIPSymbolTable`.
     - `CheckIn(customerIdx int, clientADPHVersion uint32) CheckInResponse`:
       - Branchless read of `matrix.ServiceGatesOpen[customerIdx]`.
       - Reads `matrix.MacroStages[customerIdx]` and `matrix.MicroStateBitmask[customerIdx]`.
       - Invokes `vipTable.CheckClientVersion(clientADPHVersion)`.
       - Returns `CheckInResponse{ Authorized: bool, Stage: uint16, Handshake: HandshakeResult }` in $< 10\text{ ns}$.
  2. [`00floo/afflume/81000-active-source/sentry/gate_sentry_test.go`](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/sentry/gate_sentry_test.go):
     - Tests verifying:
       - Active subscribed customer $\to$ `Authorized: true`.
       - Overdraft capped customer $\to$ `Authorized: false`.
       - Wire in transit customer $\to$ `Authorized: true` (protective override).
       - Dispute locked customer $\to$ `Authorized: false`.
       - Outdated ADPH client $\to$ `StatusUpdateRequired` with dense slot list.

---

### Phase 3: Plinth GCS Binary Slab Streaming Hydrator (Loader)
* **Goal:** Enable Cloud Run v2 tasks to stream contiguous SoA slabs straight into memory from same-region GCS with zero deserialization overhead and zero heap allocations.
* **Files to Create / Modify:**
  1. [`00floo/afflume/81000-active-source/matrix/matrix_binary_codec.go`](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/matrix/matrix_binary_codec.go):
     - `WriteTo(w io.Writer) (int64, error)`: Streams `FleetFinancialLifecycleMatrix` contiguous SoA arrays into binary format prefixed by White3 cryptoseal.
     - `ReadFrom(r io.Reader) (int64, error)`: Zero-copy hydration straight into allocated memory slices, verifying cryptoseal integrity.
  2. [`00floo/afflume/81000-active-source/matrix/matrix_binary_codec_test.go`](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/matrix/matrix_binary_codec_test.go):
     - Tests round-trip streaming serialization/deserialization for 100,000 customers.
     - Verifies byte-exact parity and benchmark throughput ($> 1\text{ GB/sec}$).

---

### Phase 4: End-to-End Integration Loop Test
* **Goal:** Test the complete 3-party loop in an end-to-end integration test.
* **Files to Create / Modify:**
  1. [`00floo/afflume/81000-active-source/three_party_integration_test.go`](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/three_party_integration_test.go):
     - Scenario:
       1. Customer registers $\to$ Plinth loads matrix.
       2. Party 3 authorizes client check-in (`ServiceGatesOpen == true`).
       3. Client runs heavy workload $\to$ emits Talos ADPH metric batch.
       4. Party 1 receives wire batch $\to$ logs to Ledger A/B $\to$ Party 2 computes consumption and subtracts VPUs.
       5. Customer exhausts VPU balance and breaches overdraft ceiling $\to$ `CoalescedLifecyclePass()` sets `OverdraftCapped` and transitions stage to `StageGraceTier1`.
       6. Service gate locks (`ServiceGatesOpen == false`).
       7. Party 3 check-in now rejects client (`Authorized: false`).
       8. Customer pays via wire $\to$ Paying event arrives $\to$ Party 2 credits VPU balance and sets `FlagWireInTransit`.
       9. Service gate immediately re-opens.

---

### Phase 5: Build, Notarize, and Preserve
* Compile & verify: `realize.exe -workspace C:\aCogSpaceSeed\00floo\afflume`
* Audit: Zero SAST vulnerabilities, 0 compiler warnings.
* Preserve & Push: `vkey construct preserve -silo 00floo`

---

## 3. Immediate Resume Directive (Upon Restart)

When restarting Antigravity, begin execution immediately with **Phase 1**:
1. Open [`conductor/tasks/afflume-three-party-plinth-execution-plan.md`](file:///c:/aCogSpaceSeed/conductor/tasks/afflume-three-party-plinth-execution-plan.md).
2. Wire `SourceTalosMetric` and `talos_binary_transmuter.go` in `eventconduit/`.
3. Proceed sequentially through Phases 1 to 5.
