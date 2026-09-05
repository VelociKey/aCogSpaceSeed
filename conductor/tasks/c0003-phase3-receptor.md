# Phase 3: Offering Receptor Rules Engine

* **Status:** Complete
* **Target Promotion Silo:** `00floo`
* **Primary Workspace:** [afflume](file:///c:/aCogSpaceSeed/00floo/afflume)
* **Key Packages:** [offering](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/offering), [afflume](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/afflume)

## Objectives
Build the VPU rules pricing resolver that consumes `ValueVector` batches (from Phase 2) and computes billable invoice components per workspace/customer. The resolver maps raw metric counts through the `VIPRules` weight table and the `OfferingRules` subscription tier to produce itemized billing line items for the checkout workflow (Phase 4).

## Current State (Audit)
- [pricing.go](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/offering/pricing.go) — `OfferingRules`, `VIPRule`, `ResolveVPUBillable()`, `CalculateTotalVPUs()`, `ResolveForWorkspace()` implemented with 5 tiers and 9 VIP metrics
- [catalog.go](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/offering/catalog.go) — product catalog definitions for all 5 tiers and feature SKUs
- [current_offering_rules.webnf](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/offering/current_offering_rules.webnf) — dynamic .webnf rules program
- [pricing_test.go](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/offering/pricing_test.go) — unit tests for VPU resolution across 5 tiers
- [vpu_receptor.go](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/offering/vpu_receptor.go) — `VPUReceptor` implementation for `vpu_billing_resolved` events

## Checklist

### 3.1 Offering DSL Grammar
- [x] Define `current_offering_rules.webnf` covering:
  - `TierSKU`, `IncludedVPUs`, `OverageSKU`, `OveragePriceUnits` (nano-unit int64)
  - `VIPRule` mapping across 9 VIP execution metrics (`vpu.units`, `tokens.llm`, `quic_stream`, `sacp_message`, `wasm_compilation`, `actor_spawn`, `grammar_validation`, `blake3_attestation`, `spool_flush`)
  - Effective date range blocks for versioned pricing changes
- [x] Implement `LoadOfferingRulesFromWebnf(path string) error` in `pricing.go` to load live rules from `current_offering_rules.webnf`
- [x] Validate loaded rules on startup: catalog fallback and feature SKUs defined for all VIP rules

### 3.2 VPU Accumulation & Receptor Resolution
- [x] `CalculateTotalVPUs(metrics map[string]int64) int64` — implemented
- [x] `ResolveVPUBillable(ctx, tierSKU, consumedVPUs) ([]billing.ProductChoice, int64, error)` — implemented
- [x] Implement `ResolveForWorkspace(ctx context.Context, wv ValueVector, tierSKU string) (*BillableInvoice, error)`
- [x] Define `BillableInvoice` struct with 9-digit integer scaling ($10^9$ nano-units)
- [x] Wire `ResolveForWorkspace` as consumer of `ValueVector` channel

### 3.3 Receptor Pipeline Integration
- [x] Implement `VPUReceptor` in [offering/vpu_receptor.go](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/offering/vpu_receptor.go) handling `vpu_billing_resolved` events
- [x] Implement RBAC guard in `afflume/rbac.go`
- [x] Add circuit breaker in `afflume/circuit.go`

### 3.4 Multi-Tier Pricing Completeness
- [x] Expand `ActiveOfferingRules` to cover all 5 subscription tiers:
  | Tier | SKU | Included VPUs | Overage Price |
  |---|---|---|---|
  | Pre-Seed | `afflume-preseed` | 700 | TBD |
  | Seed/Growth | `afflume-seed` | 2,400 | TBD |
  | Series A/B/C | `afflume-series` | 8,000 | TBD |
  | Enterprise | `afflume-enterprise` | custom | negotiated |
- [ ] Add overage VPU metric types to `ActiveVIPRules` for the full feature surface: `quic_stream`, `sacp_message`, `wasm_compilation`, `actor_spawn`, `grammar_validation`

### 3.5 Attestations
- [ ] Add `90000-authority/900-attestations/offering_attestation_test.go` covering:
  - VPU resolution for each of the 5 subscription tiers
  - Overage calculation correctness at boundary values (consumed = included, consumed = included + 1)
  - RBAC rejection for unauthorized roles
  - Circuit breaker opens after 3 failures; closes after successful resolution
  - Rule hot-reload: modify `current_offering_rules.webnf` and verify in-process reload picks up new prices

## Verification
- All `pricing_test.go` tests pass
- `ResolveForWorkspace` produces balanced billable invoices (VPUs × rate = subtotal) for all 5 tiers
- Circuit breaker attestation passes
- `int-rehydrator.exe -local-only` passes for `o-afflume`
