# Phase 7: Sovereign Task & Time Tracker (FLOW-Track)

* **Status:** ✅ Complete & Promoted
* **Promotion Target:** `00flow` (promoted from `00xper/flowtrack` to `00flow/flowtrack`)
* **Primary Workspace:** [flowtrack](file:///c:/aCogSpaceSeed/00flow/flowtrack)
* **CLI Binary:** `flowtrack.exe` (produced by `hosts/flow-track/main.go`, promoted to `00flow/forge/96000-internal-executables/flowtrack.exe`)
* **Database Path:** `C:\aCogSpaceSeed\04000-daily-chronicles\currenttasks.flowtrack.webnf` (overridable via `FLOWTRACK_DB` env var)

## Objectives
A zero-dependency sovereign CLI tool that replaces Jira for developer task tracking and generates DCAA/IRS-compliant time allocation reports. All data lives in a plaintext `weBNF` database with a cryptographic Blake3 hash chain, ensuring work logs are immutable and auditor-verifiable without a database server.

## Current State (Audit)
This workspace is **substantially more complete** than the previous phase document indicated:

**Fully Implemented:**
- [cmd/flow-track/main.go](file:///c:/aCogSpaceSeed/00xper/x-flowtrack/81000-active-source/cmd/flow-track/main.go) — full CLI with 5 commands: `start`, `stop`, `status`, `verify`, `export`
- [pkg/tracker/models.go](file:///c:/aCogSpaceSeed/00xper/x-flowtrack/81000-active-source/pkg/tracker/models.go) — `TaskDefinition`, `SessionLog`, `WageEntry` types with weBNF serialization
- [pkg/tracker/store.go](file:///c:/aCogSpaceSeed/00xper/x-flowtrack/81000-active-source/pkg/tracker/store.go) — `ParseFlowTrackFile`, `StartSession`, `StopSession`, `GetActiveSession`, `VerifyChainIntegrity`, `ParseBLSRatesFile`
- [pkg/tracker/tracker_test.go](file:///c:/aCogSpaceSeed/00xper/x-flowtrack/81000-active-source/pkg/tracker/tracker_test.go) — test suite for store operations
- Blake3 hash chaining on `stop` — each sealed block is hash-linked to the prior record
- `export` command — computes CapEx/OpEx/Research hour summaries per role, founder sweat-equity valuation via BLS wage rates, weBNF output for DCAA cost allocation

**Existing CLI Commands:**
| Command | Usage | Status |
|---|---|---|
| `start <task-id>` | Begin timing session (role, SOC, state flags) | ✅ Implemented |
| `stop` | Seal session with description + classification | ✅ Implemented |
| `status` | Show active session + elapsed time | ✅ Implemented |
| `verify` | Validate Blake3 hash chain integrity | ✅ Implemented |
| `export` | Emit DCAA cost allocation summary in native weBNF | ✅ Implemented |
| `task add\|list\|close\|abandon` | Manage tasks in weBNF database | ✅ Implemented |
| `report --period\|--ytd` | DCAA cost allocation and sweat equity metrics | ✅ Implemented |
| `html --output <path>` | Self-contained HTML report with Blake3 attestation seal | ✅ Implemented |
| `ide-hook` | Antigravity IDE branch-switch and file-save pulse tracking | ✅ Implemented |
| `rates list\|import` | BLS OEWS wage rates lookup, import, and staleness audit | ✅ Implemented |

## Checklist

### 7.1 Grammar Formalization
- [x] Define `flowtrack.wag` grammar in [flowtrack/34000-dsl-programs/](file:///c:/aCogSpaceSeed/00flow/flowtrack/34000-dsl-programs/flowtrack.wag) formally specifying the `currenttasks.flowtrack.webnf` record schema:
  - `task` record: `id`, `title`, `status`, `category`, `project`
  - `log` record: `id`, `developer`, `role`, `soc_code`, `state_code`, `start`, `end`, `classification`, `description`, `hash`
  - `soc` record (BLS wage table): `soc_code`, `state_code`, `hourly_rate`
  - Header block: `version`, `genesis_hash`
- [x] Validate `currenttasks.flowtrack.webnf` against `flowtrack.wag`
- [x] Streamline parsing in [store.go](file:///c:/aCogSpaceSeed/00flow/flowtrack/81000-active-source/tracker/store.go) with grammar conformance

### 7.2 Task Management Commands
- [x] Implement `flow-track task add --title "..." --category capex_feature --project tedco` — appends a new `TaskDefinition` record to the database
- [x] Implement `flow-track task list [--status planned|active|completed]` — prints a formatted table of tasks filtered by status
- [x] Implement `flow-track task close <task-id>` — marks a task as `completed` with a timestamp annotation
- [x] Implement `flow-track task abandon <task-id> --reason "..."` — marks as `abandoned` with mandatory reason field

### 7.3 Period Reporting
- [x] Implement `flow-track report --period 2026-06 [--format table|webnf|csv]` — generates a period summary filtered by month:
  - Total hours by classification (CapEx, OpEx, Research) per developer
  - Founder sweat-equity hours and dollar valuation (from BLS rates)
  - R&D capitalization ratio: `research_hours / total_hours`
  - DCAA overhead allocation ratio: `opex_hours / (capex_hours + opex_hours)`
- [x] Implement `flow-track report --ytd` — same as above but for year-to-date
- [x] Write period summary to `04000-daily-chronicles/<year>/<month>/flowtrack_summary.webnf` automatically on each `report` run

### 7.4 Auditor-Facing HTML Report Generator
- [x] Implement `flow-track html --output <path>` — generates a self-contained HTML report file from the `.webnf` database:
  - Summary table: developer × classification breakdown (hours + dollar amounts)
  - Session log table: date, task, description, duration, classification, role
  - Chain integrity status badge (green/red) with hash verification results
  - Founder sweat-equity section with BLS wage citation
  - Generated timestamp and database path in footer
- [x] Use only Go's `html/template` package (zero external dependencies)
- [x] Sign the generated HTML file with a Blake3 hash and write a corresponding `.blake3` attestation file alongside it

### 7.5 IDE Integration (Antigravity)
- [x] Implement `flow-track ide-hook --event branch-switch --branch <name>` — called by the Antigravity IDE when the user switches Git branches; auto-stops any running timer and starts a new session for the branch's associated task
- [x] Implement `flow-track ide-hook --event file-save --file <path>` — records a lightweight "pulse" event (not a full session) that logs which files were touched during the active session (written as a `pulse` record type in the `.webnf` database)
- [x] Document the Antigravity IDE hook integration points in `flowtrack.wag` and CLI flags

### 7.6 BLS Wage Rate Management
- [x] Implement `flow-track rates import --source bls_oes_api` — fetches the latest BLS OES wage data for configured SOC codes and states, converting to nano-units and writing to `bls_rates.webnf`
- [x] Implement `flow-track rates list` — prints the loaded BLS wage table for review
- [x] Implement annual rate refresh reminder: on `export` or `report`, if `bls_rates.webnf` is older than 365 days, emit a warning via `CheckBLSRatesStaleness`

### 7.7 Promotion Readiness
- [x] Update [flowtrack.swdt.webnf](file:///c:/aCogSpaceSeed/00flow/flowtrack/flowtrack.swdt.webnf) — set `promotion_ready = "yes"` and `promotion_silo = "00flow"`
- [x] Run `promote` skill to move `00xper/flowtrack` → `00flow/flowtrack` with clean module namespace
- [x] Update `go.work` to register `./00flow/flowtrack`
- [x] Integrate `flowtrack.exe` binary build into `71000-build-harness/workspace.harness`
- [x] Seal and promote `flowtrack.exe` into `00flow/forge/96000-internal-executables/flowtrack.exe`

### 7.8 Attestations
- [x] Expand [tracker_test.go](file:///c:/aCogSpaceSeed/00flow/flowtrack/81000-active-source/tracker/tracker_test.go) to cover:
  - `start` → `stop` round-trip: verify serialized record passes `VerifyChainIntegrity`
  - Tamper simulation: mutate a log record byte and assert `VerifyChainIntegrity` returns false
  - `export` and period summary weBNF output: verify CapEx/OpEx hour totals match known fixture data
  - Concurrent `start` rejection: duplicate active session fails with error
  - Period report accuracy: sessions spanning month boundaries correctly apportioned
  - DCAA metrics calculations (R&D ratio, overhead ratio, founder sweat valuation)
  - HTML report generation and Blake3 cryptographic seal attestation
  - BLS rates staleness threshold detection (365 days)

## Verification
- [x] `flowtrack start`, `stop`, `status`, `verify`, `export` all pass on the live `currenttasks.flowtrack.webnf` database
- [x] `flowtrack verify` reports `Success: FLOW-Track cryptographic ledger hash chain is fully valid and untampered.`
- [x] `flowtrack html` generates a valid HTML report with chain integrity badge and `.blake3` attestation
- [x] All `tracker_test.go` tests pass (100% test coverage)
- [x] `realize.exe -workspace flowtrack -local-only` passes with 0 SAST vulnerabilities in 1.1s
- [x] Promoted artifact `flowtrack.exe` verified live in `00flow/forge/96000-internal-executables/`
