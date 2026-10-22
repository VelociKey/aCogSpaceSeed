# Phase 4: Concurrent Ordering Workflow

* **Status:** ✅ 100% Complete & Verified
* **Target Promotion Silo:** `00floo`
* **Primary Workspace:** [`00floo/afflume`](file:///c:/aCogSpaceSeed/00floo/afflume) (module: `sov.fleet/afflume`)
* **Key Packages:**
  - Core Workflow: [`00floo/afflume/81000-active-source/workflow`](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/workflow)
  - Tax Calculation: [`00floo/afflume/81000-active-source/taxinvoicer`](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/taxinvoicer)
  - Emulators & Concurrency: [`00floo/afflume/81000-active-source/emulators`](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/emulators), [`sov.fleet/substrate/concurrency/81000-active-source/adaptivelock`](file:///c:/aCogSpaceSeed/00ulib/substrate/concurrency/81000-active-source/adaptivelock)
  - Stripe Sandbox / Emulator: [`00flow/emulator-stripetax`](file:///c:/aCogSpaceSeed/00flow/emulator-stripetax)

## Objectives
Implement the production-ready concurrent checkout state machine that takes a `BillableInvoice` (from Phase 3), runs it through tax calculation, Stripe invoice creation/finalization, payment confirmation, and atomic ledger commit — all with full concurrency-safe session management using `adaptivelock.Map` and emulator-backed integration tests.

## Current State (Workspace Audit)
- [`checkout_workflow.go`](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/workflow/checkout_workflow.go) — Contains initial scaffold of `CheckoutSession`, `Registry`, and `WorkflowEngine`.
  - *Identified Migration Requirement:* Migrate `Registry.sessions` from `sync.RWMutex + map[string]*CheckoutSession` to `*adaptivelock.Map[string, *CheckoutSession]` per sovereign concurrency invariant.
  - *Identified FSM Alignment:* Collapse interim pending states into the canonical 7-state FSM.
- [`checkout_workflow_test.go`](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/workflow/checkout_workflow_test.go) — Unit tests for FSM transitions, E2E workflow, concurrent checkouts, and journal replay.
- [`emulators/cloudtasks.go`](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/emulators/cloudtasks.go) & [`emulators/pubsub.go`](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/emulators/pubsub.go) — Fully implemented using `adaptivelock.Map`.
- [`taxinvoicer/checkout.go`](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/taxinvoicer/checkout.go) — `CheckoutOrchestrator` coordinating AvaTax estimates, Stripe draft items, and finalization.

## Checklist

### 4.1 State Machine Completion (✅ Complete & Verified)
- [x] Canonical 7-State Enum `CheckoutState`:
  ```go
  StateInitiated          CheckoutState = "INITIATED"
  StateTaxCalculated      CheckoutState = "TAX_CALCULATED"
  StateStripeDraftCreated CheckoutState = "STRIPE_DRAFT_CREATED"
  StateStripeFinalized    CheckoutState = "STRIPE_FINALIZED"
  StatePaymentConfirmed   CheckoutState = "PAYMENT_CONFIRMED"
  StateLedgerCommitted    CheckoutState = "LEDGER_COMMITTED" // Terminal
  StateFailed             CheckoutState = "FAILED"           // Terminal
  ```
- [x] Thread-safe `CheckoutSession` struct with mutex-protected field updates and `sync.Pool` recycling.
- [x] Migrate `Registry` from `sync.RWMutex` + `map[string]*CheckoutSession` to lock-free `*adaptivelock.Map[string, *CheckoutSession]`.
- [x] Implement canonical `TransitionTo(ctx context.Context, invoiceID string, next CheckoutState, errArgs ...any) error`:
  - Enforces the strict transition table:
    ```
    INITIATED            → TAX_CALCULATED | FAILED
    TAX_CALCULATED       → STRIPE_DRAFT_CREATED | FAILED
    STRIPE_DRAFT_CREATED → STRIPE_FINALIZED | FAILED
    STRIPE_FINALIZED     → PAYMENT_CONFIRMED | FAILED
    PAYMENT_CONFIRMED    → LEDGER_COMMITTED | FAILED
    LEDGER_COMMITTED     → (terminal: invalid transition error)
    FAILED               → (terminal: invalid transition error)
    ```
  - Emits `slog.Info("checkout_transition", "invoice_id", invoiceID, "from", curr, "to", next)`.
  - Appends transition record to `sessions.workflow.webnf` journal via non-blocking channel or direct write.
- [x] Implement `FailSession(ctx context.Context, invoiceID string, err error) error`:
  - Transitions session to `StateFailed`, records `LastError`, and emits `slog.Error("checkout_transition_failed", ...)`.
- [x] Implement `RecoverFailedSessions(ctx context.Context)`:
  - Iterates over registry via `sessions.Range(...)`, logging operator triage warnings for all sessions in `StateFailed`.

### 4.2 Tax Calculation Step
- [x] Integrate [`00floo/afflume/81000-active-source/taxinvoicer`](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/taxinvoicer) into `CalculateTax(ctx context.Context, session *CheckoutSession) error`:
  - Evaluates tax amounts in 9-digit integer scaling ($10^9$ nano-units).
  - Updates `session.TaxAmount` and transitions to `StateTaxCalculated` or `StateFailed`.

### 4.3 Stripe API Integration (Emulator-Backed)
- [x] Implement `CreateStripeDraftInvoice(ctx context.Context, session *CheckoutSession) error`:
  - Dispatches draft invoice creation to Stripe emulator / `StripeTaxClient`.
  - Writes `StripeDraftID` to session; transitions to `StateStripeDraftCreated`.
- [x] Implement `FinalizeStripeInvoice(ctx context.Context, session *CheckoutSession) error`:
  - Finalizes draft invoice; receives hosted payment URL.
  - Writes `PaymentURL` to session; transitions to `StateStripeFinalized`.
- [x] Implement `ConfirmPayment(ctx context.Context, session *CheckoutSession, webhookEvent webhooktransmuter.NormalizedEvent) error`:
  - Triggered upon receiving payment confirmation webhook.
  - Validates `StripeDraftID` matches session; transitions to `StatePaymentConfirmed`.

### 4.4 Atomic Ledger Commit Step
- [x] Implement `CommitToLedger(ctx context.Context, session *CheckoutSession) error`:
  - Dispatches balanced double-entry postings to [`generalledger`](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/generalledger).
  - Enforces idempotency (if already in `StateLedgerCommitted`, return `nil` immediately).
  - Transitions to `StateLedgerCommitted`.

### 4.5 Session Persistence & Recovery
- [x] Maintain append-only session journal in `sessions.workflow.webnf`.
- [x] Verify `LoadActiveSessions(path string) ([]*CheckoutSession, error)` faithfully reconstructs active in-flight sessions on daemon startup.

### 4.6 Concurrent Load Handling
- [x] Implement `RunConcurrentCheckouts(ctx context.Context, invoices []*offering.BillableInvoice, customerUUID string, concurrencyLimit int) error`.
- [x] Enforce backpressure check: reject and emit `slog.Warn("checkout_backpressure_active")` when active sessions $\ge 500$.

### 4.7 Verification
- [x] `realize.exe -workspace afflume -local-only` compiles cleanly and passes all test suites.
- [x] `checkout_workflow_test.go` verifies all valid and invalid transitions with 100% state coverage.
- [x] Concurrent 50-session benchmark completes with 0 race conditions and 0 heap leaks.
