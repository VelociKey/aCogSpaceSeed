# Phase 4: Concurrent Ordering Workflow

* **Status:** Planned
* **Target Promotion Silo:** `00floo`
* **Primary Workspace:** [o-afflume](file:///c:/aCogSpaceSeed/00floo/o-afflume)
* **Key Package:** [pkg/workflow](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/workflow)
* **Emulator Adapter:** [x-wrapstripe](file:///c:/aCogSpaceSeed/00xper/x-wrapstripe)

## Objectives
Implement the production-ready concurrent checkout state machine that takes a `BillableInvoice` (from Phase 3), runs it through tax calculation, Stripe invoice creation/finalization, payment confirmation, and atomic ledger commit — all with full concurrency-safe session management and emulator-backed integration tests.

## Current State (Audit)
The following have already been scaffolded:
- [checkout_workflow.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/workflow/checkout_workflow.go) — 7-state machine (`CheckoutState`), thread-safe `CheckoutSession`, `WorkflowRegistry` using `adaptivelock`
- [checkout_workflow_test.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/workflow/checkout_workflow_test.go) — basic state transition unit tests
- [pkg/emulators/cloudtasks.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/emulators/cloudtasks.go) — Cloud Tasks emulator adapter
- [pkg/emulators/pubsub.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/emulators/pubsub.go) — Pub/Sub emulator adapter
- [x-wrapstripe/pkg/wrapstripe/sandbox.go](file:///c:/aCogSpaceSeed/00xper/x-wrapstripe/81000-active-source/pkg/wrapstripe/sandbox.go) — Stripe emulator sandbox wrapper

## Checklist

### 4.1 State Machine Completion
- [x] State enum `CheckoutState` with 7 states — implemented
- [x] `CheckoutSession` struct (thread-safe, `sync.Mutex`) — implemented
- [x] `WorkflowRegistry` (`adaptivelock.Map`) — implemented
- [ ] Implement `TransitionTo(ctx context.Context, invoiceID string, next CheckoutState) error` — validates allowed transitions (FSM guard), updates `State` and `UpdatedAt`, emits `slog.Info("checkout_transition", "invoice_id", ..., "from", ..., "to", ...)`
- [ ] Define the allowed FSM transition table:
  ```
  INITIATED         → TAX_CALCULATED | FAILED
  TAX_CALCULATED    → STRIPE_DRAFT_CREATED | FAILED
  STRIPE_DRAFT_CREATED → STRIPE_FINALIZED | FAILED
  STRIPE_FINALIZED  → PAYMENT_CONFIRMED | FAILED
  PAYMENT_CONFIRMED → LEDGER_COMMITTED | FAILED
  LEDGER_COMMITTED  → (terminal)
  FAILED            → (terminal, emit structured error)
  ```
- [ ] Implement `RecoverFailedSessions(ctx context.Context)` — on startup, scan persisted sessions in state `FAILED` and emit `slog.Warn` records with `invoice_id` and `last_error` for operator triage

### 4.2 Tax Calculation Step
- [ ] Integrate [pkg/taxinvoicer](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/taxinvoicer) into the `TAX_CALCULATED` transition
- [ ] Implement `CalculateTax(ctx context.Context, session *CheckoutSession) error`:
  - Calls `taxinvoicer.Calculate(subtotalNano, currencyCode, customerJurisdiction)`
  - Writes `TaxAmount` back to session in nano-units
  - Transitions to `TAX_CALCULATED` or `FAILED`

### 4.3 Stripe API Integration (Emulator-Backed)
- [ ] Implement `CreateStripeDraftInvoice(ctx context.Context, session *CheckoutSession) error`:
  - Calls Stripe (via `x-wrapstripe` sandbox) to create a draft invoice with the line items from the `BillableInvoice`
  - Writes `StripeDraftID` to session; transitions to `STRIPE_DRAFT_CREATED`
- [ ] Implement `FinalizeStripeInvoice(ctx context.Context, session *CheckoutSession) error`:
  - Finalizes the draft invoice; receives hosted payment URL
  - Writes `PaymentURL` to session; transitions to `STRIPE_FINALIZED`
- [ ] Implement `ConfirmPayment(ctx context.Context, session *CheckoutSession, webhookEvent ReceptorPayload) error`:
  - Called by the Stripe webhook adapter (Phase 1) when `invoice.payment_succeeded` arrives
  - Validates `StripeDraftID` matches session; transitions to `PAYMENT_CONFIRMED`

### 4.4 Atomic Ledger Commit Step
- [ ] Implement `CommitToLedger(ctx context.Context, session *CheckoutSession) error`:
  - Calls `LedgerPoster.ResolveEventPostings("invoice.payment_succeeded", ...)` with the session's payload
  - Writes the resulting postings to the append-only ledger chain
  - Transitions to `LEDGER_COMMITTED`
- [ ] Implement idempotency: if `LEDGER_COMMITTED` is already set for a given `InvoiceID`, skip without error

### 4.5 Session Persistence
- [ ] Persist `CheckoutSession` state to a `sessions.workflow.webnf` file (append-only, one record per transition) so that in-flight sessions survive process restarts
- [ ] Define `workflow.wag` grammar for the session record format
- [ ] Implement `LoadActiveSessions(path string) ([]*CheckoutSession, error)` to replay the session log on startup

### 4.6 Concurrent Load Handling
- [ ] Implement `RunConcurrentCheckouts(ctx context.Context, invoices []*BillableInvoice) error` — fans out checkout sessions across a `sync.WaitGroup` with a configurable concurrency limit (default: 10 concurrent sessions)
- [ ] Add back-pressure: if `WorkflowRegistry` exceeds 500 active sessions, block new submissions and emit `slog.Warn("checkout_backpressure_active")`

### 4.7 Attestations
- [ ] Add `90000-authority/900-attestations/workflow_attestation_test.go` covering:
  - All valid FSM transitions
  - All invalid FSM transitions return errors
  - Concurrent 50-session load test: all complete in `LEDGER_COMMITTED` state
  - Session recovery: FAILED sessions are discovered and logged on restart
  - Idempotent re-commit: `CommitToLedger` called twice for same `InvoiceID` is a no-op

## Verification
- `checkout_workflow_test.go` passes with 100% state coverage
- Concurrent 50-session load test completes within 10s against the Stripe emulator
- All committed sessions produce balanced ledger postings (attested in Phase 6)
- `int-rehydrator.exe -local-only` passes for `o-afflume`
