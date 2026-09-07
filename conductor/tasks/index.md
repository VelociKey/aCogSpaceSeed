# Project Task Index: sTelemetry, Value Vector Studio & Ledger Integration

This index lists the modular phases mapping the platform monetization engine, identity security, and ledger integration workflows. All code promotions target silo `00floo`. Phases are ordered by dependency: each phase builds on the outputs of the prior.

## Phase Pipeline Overview

```
[Phase 0] Workspace Alignment
      │
      ▼
[Phase 1] Double-Entry Ledger Backend  ←─────────────────────┐
      │                                                        │
      ▼                                                        │
[Phase 2] sTelemetry qAPC Receiver                            │
      │                                                        │
      ▼                                                        │
[Phase 3] Offering Receptor Rules Engine                       │
      │                                                        │
      ▼                                                        │
[Phase 4] Concurrent Ordering Workflow ─── posts to ──────────┘
      │
      ▼
[Phase 5] Studio & Portal Frontends (consumes Phases 1–4 APIs)
      │
      ▼
[Phase 6] E2E Integration Suite (validates Phases 1–5 end-to-end)

[Phase 7] FLOW-Track (parallel; sovereign tooling, not a pipeline dependency)
```

## Phases & Tasks

| # | Task | Workspace | Status |
|---|------|-----------|--------|
| **[Phase 0](file:///C:/aCogSpaceSeed/conductor/tasks/c0000-phase0-alignment.md)** | Workspace and Title Alignment | `o-afflume`, `x-meter`, `x-wrapstripe`, `x-flowtrack` | ✅ Complete |
| **[Phase 1](file:///C:/aCogSpaceSeed/conductor/tasks/c0001-phase1-ledger.md)** | Core Double-Entry Ledger & Ingestor Backend | `o-afflume/pkg/generalledger` | ✅ Complete |
| **[Phase 2](file:///C:/aCogSpaceSeed/conductor/tasks/c0002-phase2-stelemetry.md)** | sTelemetry & qAPC ValueVector Receiver | `meter`, `afflume/qapc` | ✅ Complete |
| **[Phase 3](file:///C:/aCogSpaceSeed/conductor/tasks/c0003-phase3-receptor.md)** | Offering Receptor Rules Engine | `afflume/offering` | ✅ Complete |
| **[Phase 4](file:///C:/aCogSpaceSeed/conductor/tasks/c0004-phase4-ordering.md)** | Concurrent Ordering Workflow | `o-afflume/pkg/workflow`, `x-wrapstripe` | 📋 Planned |
| **[Phase 5](file:///C:/aCogSpaceSeed/conductor/tasks/c0005-phase5-studios.md)** | Symmetrical Studio & Portal Frontends | `o-afflume-realization` | 📋 Planned |
| **[Phase 6](file:///C:/aCogSpaceSeed/conductor/tasks/c0006-phase6-integration.md)** | E2E Integration Suite & Verification | `o-afflume/cmd/o-afflume` | 📋 Planned |
| **[Phase 7](file:///C:/aCogSpaceSeed/conductor/tasks/c0007-phase7-flowtrack.md)** | Sovereign Task & Time Tracker (FLOW-Track) | `x-flowtrack` | 🔄 Active (more complete than initially documented) |

## Key Precision Standards (All Phases)
- **Monetary amounts:** 9-digit integer scaling ($10^9$ nano-units). No floating-point at any layer.
- **Chain integrity:** All append-only `.webnf` files use Blake3 hash-linking.
- **Grammar compliance:** All `.webnf` programs must validate against their `.wag` grammar.
- **Logging:** All diagnostic output uses `log/slog` structured logging.
- **Build:** All compilation via `int-rehydrator.exe -local-only`. No raw `go build`.
- **Concurrency:** All shared maps use `adaptivelock.Map`; no `sync.Mutex`-wrapped standard maps.
