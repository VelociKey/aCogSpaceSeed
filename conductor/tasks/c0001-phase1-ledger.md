# Phase 1: Core Double-Entry Ledger & Ingestor Backend

* **Status:** Active
* **Target Promotion Silo:** `00floo`
* **Primary Workspace:** [o-afflume](file:///c:/aCogSpaceSeed/00floo/o-afflume)

---

## Architectural Context

Phase 1 is organized around two clearly separated concerns that **must not be conflated**:

### A. The Afflume Receptor Abstraction Pipeline

This is the **core, vendor-neutral** business pipeline. It expresses the full offering-to-ledger lifecycle through abstract interfaces defined in [`pkg/billing/receptors.go`](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/billing/receptors.go). No Stripe types, no HMAC, no vendor SDKs appear here.

```
OfferingProvider         → resolves what is being sold (SKUs, prices)
      │
      ▼
InvoiceProvider          → builds AuthoritativeInvoice (subtotal, line items, tax)
      │
      ▼
InvoiceProcessor         → certifies the draft invoice (FiscalStamp, final TaxAmountUnits)
      │
      ▼
PaymentEngine            → pushes the certified invoice to a payment gateway → PaymentSession
      │
      ▼
BankMachineReceptor      → reconciles payout confirmation from banking rails
      │
      ▼
LedgerPoster             → converts settled PaymentSession events → balanced double-entry postings
```

Each step is an **interface**. Any vendor can plug in by implementing the interface. The pipeline itself never references a vendor.

### B. The Stripe VendorAdapter (a separate concern from Phase 1 core)

The Stripe adapter lives in [`pkg/webhooktransmuter/`](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/webhooktransmuter) and [`pkg/afflume/adapters.go`](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/afflume/adapters.go). It implements:
- `WebhookAdapter.VerifyAndParse()` — HMAC-SHA256 signature verification (Stripe-specific) and raw JSON → `NormalizedEvent` translation
- `VendorAdapter.Translate()` — converts Stripe payload types to the normalized receptor input shapes

**Stripe-specific concerns (HMAC, Stripe event type strings, Stripe SDK types) are fully encapsulated inside the adapter.** The pipeline only ever sees `NormalizedEvent` and the abstract receptor interfaces.

---

## Precision Standards
- **Monetary amounts:** 9-digit integer scaling ($10^9$ nano-units). No floating-point at any layer.
- **Immutability:** The `ledger.webnf` chain is append-only. No record is ever mutated.
- **Multi-Currency:** Variable scaling strategies supported via currency metadata blocks (e.g., JPY zero-decimal override vs. USD standard scaling).

---

## Current State (Audit)

### Fully Implemented
| File | What It Is |
|---|---|
| [billing/receptors.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/billing/receptors.go) | `OfferingProvider`, `InvoiceProvider`, `InvoiceProcessor`, `PaymentEngine`, `BankMachineReceptor` interfaces; `AuthoritativeInvoice`, `CertifiedInvoice`, `DraftInvoice`, `PaymentSession`, `CustomerProfile`, `TaxDetail` types |
| [afflume/adapters.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/afflume/adapters.go) | `VendorAdapter` interface, `AdapterFactory`, stub impls for Stripe, Xero, AvaTax, Relay |
| [afflume/receptor.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/afflume/receptor.go) | `Receptor` interface, lock-free `Registry`, `DunningManager` |
| [taxinvoicer/sovereign.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/taxinvoicer/sovereign.go) | `SovereignInvoicer` implements `billing.InvoiceProvider` — Wayfair nexus tax calculation |
| [taxinvoicer/checkout.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/taxinvoicer/checkout.go) | `CheckoutOrchestrator`, `StateNexusLedger` (economic nexus tracking) |
| [taxinvoicer/avatax.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/taxinvoicer/avatax.go) | AvaTax adapter implementation |
| [taxinvoicer/stripe_tax.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/taxinvoicer/stripe_tax.go) | Stripe Tax adapter implementation |
| [taxinvoicer/webhook.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/taxinvoicer/webhook.go) | `WebhookHandler` for AvaTax commit on payment success |
| [invoiceprocessor/mock_processor.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/invoiceprocessor/mock_processor.go) | `SovereignInvoiceProcessor` implements `billing.InvoiceProcessor` — fiscal stamp + tax certification |
| [webhooktransmuter/factory.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/webhooktransmuter/factory.go) | `WebhookAdapter` interface, `TransmuterFactory`, `NormalizedEvent` (vendor-neutral output type) |
| [webhooktransmuter/parser.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/webhooktransmuter/parser.go) | `ParseEventWebNF()` — parses normalized `.webnf` event records into `IngestedEvent` |
| [generalledger/poster.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/generalledger/poster.go) | `LedgerPoster`, `LedgerRules`, `EntryRule`, `EventRule` |
| [generalledger/journal.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/generalledger/journal.go) | `TransactionJournal` event sourcing model |
| [generalledger/ingestor.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/generalledger/ingestor.go) | Polling-based file watcher → processes journal events → ledger postings |
| [generalledger/mock_ledger.go](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/generalledger/mock_ledger.go) | In-memory ledger for test isolation |
| [currentledgercoa.webnf](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/generalledger/currentledgercoa.webnf) | Active Chart of Accounts program |
| [masterledgercoa.webnf](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/generalledger/masterledgercoa.webnf) | Full Chart of Accounts reference |
| [currentledgerrules.webnf](file:///c:/aCogSpaceSeed/00floo/o-afflume/81000-active-source/pkg/generalledger/currentledgerrules.webnf) | Active posting rules program |

---

## Checklist

### 1.1 Receptor Pipeline — Interface Completeness
- [x] `OfferingProvider` interface — defined
- [x] `InvoiceProvider` interface (`CreateInvoice`) — defined + `SovereignInvoicer` implementation
- [x] `InvoiceProcessor` interface (`CertifyInvoice`) — defined + `SovereignInvoiceProcessor` implementation
- [x] `PaymentEngine` interface (`PushInvoiceForPayment`) — defined
- [x] `BankMachineReceptor` interface (`RegisterPayout`, `ConfirmDeposit`) — defined

### 1.2 Tax Calculation — Expansion & Correctness
- [x] Wayfair nexus rules for SD (4.5%) and MD (6.0%) in `SovereignInvoicer.CreateInvoice` — implemented
- [x] `StateNexusLedger.CheckNexus()` — economic nexus threshold tracking — implemented
- [x] Integer-only precision and validation for financial checks — implemented

### 1.3 Invoice Processor — Certification
- [x] `SovereignInvoiceProcessor.CertifyInvoice` — fiscal stamp + tax calculation — implemented (Upgraded to Blake3)

### 1.4 Stripe Adapter — Webhook Transmuter
- [x] `WebhookAdapter` interface — defined in `webhooktransmuter/factory.go`
- [x] `NormalizedEvent` type — defined (provider-neutral output)
- [x] `TransmuterFactory` — register/retrieve adapters by provider name
- [x] Standardize Zero-Knowledge SOC 2 attestation handler and tests on Blake3 — implemented

---

## Verification
- `int-rehydrator.exe -local-only` passes for `o-afflume`
- Full pipeline attestation: `OfferingProvider` → `InvoiceProvider` → `InvoiceProcessor` → mock `PaymentEngine` → `LedgerPoster` produces balanced postings for all 5 subscription tiers
- Stripe adapter attestation: HMAC rejection and all 5 event types verified
- Chain integrity validator reports zero broken links on a 1,000-entry synthetic chain
- `go test -race ./...` passes — zero data races across the pipeline
