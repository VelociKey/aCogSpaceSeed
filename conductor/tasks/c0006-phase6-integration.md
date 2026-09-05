# Phase 6: E2E Integration Suite & Verification

* **Status:** Planned
* **Target Promotion Silo:** `00floo`
* **Primary Workspace:** [o-afflume](file:///c:/aCogSpaceSeed/00floo/o-afflume)
* **Test Entry Point:** [cmd/o-afflume/integration_test.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/cmd/o-afflume/integration_test.go)

## Objectives
Build and run a fully automated E2E workflow simulation that exercises the complete payment pipeline: sTelemetry accumulation → VPU pricing → Stripe checkout → webhook receipt → balanced double-entry ledger postings. All steps must be observable via structured `slog` output and require zero manual intervention to run.

## Current State (Audit)
A foundational E2E test already exists:
- [integration_test.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/cmd/o-afflume/integration_test.go) — `TestEndToEndTelemetryToLedgerWorkflow` covers sTelemetry parse → VPU calculate → checkout state machine → ledger posting in a single linear happy-path scenario

**Gaps identified:**
- No parallel/concurrent transaction simulation
- No Stripe emulator webhook delivery in-loop (webhook is simulated as a direct function call, not an actual HTTP POST)
- No chain integrity verification step
- No failure/recovery path tests
- No volume/stress simulation
- No latency or dedup test coverage (files `latency_dedup_test.go` and `volume_test.go` exist but scope is unknown)

## Checklist

### 6.1 Happy-Path Regression Hardening
- [x] `TestEndToEndTelemetryToLedgerWorkflow` — basic linear E2E exists
- [ ] Expand to cover all 5 subscription tiers: for each tier, run the full pipeline with `consumedVPUs = includedVPUs - 1` (no overage) and `consumedVPUs = includedVPUs + 100` (with overage), asserting correct line items and balanced postings
- [ ] Add assertion: `sum(debits) == sum(credits)` verified by calling `MockLedger.AssertBalanced(t)` after each test case
- [ ] Add assertion: `LedgerChain.VerifyChainIntegrity()` passes after all postings are committed

### 6.2 Concurrent Parallel Transaction Simulation
- [ ] Implement `TestConcurrentParallelCheckouts(t *testing.T)`:
  - Launch 25 goroutines, each running a full E2E pipeline with a unique `InvoiceID` and `WorkspaceID`
  - Use `sync.WaitGroup` to wait for all sessions to reach `LEDGER_COMMITTED`
  - Assert all 25 sessions complete without `FAILED` state
  - Assert the ledger contains exactly 25 × N posting records (N = postings per event type)
  - Assert chain integrity passes at the end

### 6.3 Stripe Emulator Webhook Loop
- [ ] Implement `TestStripeWebhookLoop(t *testing.T)`:
  - Start a local HTTP test server (`httptest.NewServer`) simulating the `o-afflume` webhook endpoint
  - Send a realistic `invoice.payment_succeeded` Stripe webhook POST with HMAC-SHA256 signature to the test server
  - Verify: webhook is parsed → `ConfirmPayment` is called → `CommitToLedger` runs → final state is `LEDGER_COMMITTED`
  - Send a `charge.refunded` webhook and verify: correcting journal entry posted, `FAILED` state not triggered

### 6.4 Failure & Recovery Path Tests
- [ ] Implement `TestCheckoutFailureRecovery(t *testing.T)`:
  - Inject a Stripe API error (mocked) during `FinalizeStripeInvoice` → assert session transitions to `FAILED`
  - Restart simulated process (re-load sessions from `sessions.workflow.webnf`) → assert FAILED session is discovered and logged via `RecoverFailedSessions`
  - Manually transition the recovered session back to `INITIATED` and re-run → assert it reaches `LEDGER_COMMITTED`
- [ ] Implement `TestIdempotentWebhookDelivery(t *testing.T)`:
  - Send the same `invoice.payment_succeeded` webhook twice
  - Assert the ledger contains exactly N postings (not 2N) — dedup via ingestor idempotency gate

### 6.5 Volume & Latency Tests
- [ ] Review and expand [latency_dedup_test.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/cmd/o-afflume/latency_dedup_test.go) — establish concrete latency SLA assertions (e.g., P99 ledger posting time < 50ms for 1,000 concurrent events)
- [ ] Review and expand [volume_test.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/cmd/o-afflume/volume_test.go) — simulate 10,000 spool records ingested via the qAPC receiver; assert all parsed and VPU totals match expected aggregate
- [ ] Implement [invoice_processor_test.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/cmd/o-afflume/invoice_processor_test.go) coverage — run the `invoiceprocessor` package with 500 concurrent invoices and assert zero data races (use `-race` flag)

### 6.6 Structured Observability Assertions
- [ ] Implement a `slog.Handler` test spy that captures log records during test runs
- [ ] Assert that each E2E run emits the following structured log events in order:
  1. `sTelemetry_received` with `workspace_id`, `metric_count`
  2. `vpu_resolved` with `total_vpus`, `overage_vpus`, `tier_sku`
  3. `checkout_transition` for each state transition
  4. `webhook_received` with `event_type`, `invoice_id`
  5. `ledger_posted` with `posting_count`, `balanced: true`
  6. `chain_integrity_verified` with `record_count`, `status: ok`

### 6.7 CI/CD Integration
- [ ] Create `71000-build-harness/run_e2e_tests.ps1` — a PowerShell script that:
  1. Runs `int-rehydrator.exe -local-only` to build the binary
  2. Runs `go test -v -race -count=1 ./...` in `cmd/o-afflume`
  3. Captures output to `c0990-ephemeral-scratch/e2e-test.log`
  4. Exits with error code on any failure

## Verification
- All happy-path tests pass for all 5 subscription tiers
- Concurrent 25-session test completes with zero FAILED states and a balanced ledger
- Stripe webhook loop test passes: dedup confirmed, refund correcting entry confirmed
- `go test -race` reports zero data races across all integration tests
- Chain integrity verified after each test suite run
