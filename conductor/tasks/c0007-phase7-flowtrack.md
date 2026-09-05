# Phase 7: Sovereign Task & Time Tracker (FLOW-Track)

* **Status:** Active (substantially implemented; gaps remain in grammar, IDE integration, and HTML report)
* **Promotion Target:** `00flow` (from `00xper/x-flowtrack`)
* **Primary Workspace:** [x-flowtrack](file:///c:/aCogSpaceSeed/00xper/x-flowtrack)
* **CLI Binary:** `flow-track.exe` (produced by `cmd/flow-track/main.go`)
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
- `export` command — computes CapEx/OpEx/Research hour summaries per role, founder sweat-equity valuation via BLS wage rates, JSON output for DCAA cost allocation

**Existing CLI Commands:**
| Command | Usage | Status |
|---|---|---|
| `start <task-id>` | Begin timing session (role, SOC, state flags) | ✅ Implemented |
| `stop` | Seal session with description + classification | ✅ Implemented |
| `status` | Show active session + elapsed time | ✅ Implemented |
| `verify` | Validate Blake3 hash chain integrity | ✅ Implemented |
| `export` | Emit DCAA JSON cost allocation summary | ✅ Implemented |

## Checklist

### 7.1 Grammar Formalization
- [ ] Define `flowtrack.wag` grammar in [x-flowtrack/34000-dsl-programs/](file:///c:/aCogSpaceSeed/00xper/x-flowtrack) (or `00000-knowledge-foundations/`) to formally specify the `currenttasks.flowtrack.webnf` record schema:
  - `task` record: `id`, `title`, `status`, `category`, `project`
  - `log` record: `id`, `developer`, `role`, `soc_code`, `state_code`, `start`, `end`, `classification`, `description`, `hash`
  - `soc` record (BLS wage table): `soc_code`, `state_code`, `hourly_rate`
  - Header block: `version`, `genesis_hash`
- [ ] Validate `currenttasks.flowtrack.webnf` (if it exists in `04000-daily-chronicles/`) against `flowtrack.wag` via `fab-inspect.exe`
- [ ] Replace regex-based parsing in [store.go](file:///c:/aCogSpaceSeed/00xper/x-flowtrack/81000-active-source/pkg/tracker/store.go) with codec-generated parser from `flowtrack.wag` (eliminates fragile regex drift as grammar evolves)

### 7.2 Task Management Commands
- [ ] Implement `flow-track task add --title "..." --category capex_feature --project tedco` — appends a new `TaskDefinition` record to the database
- [ ] Implement `flow-track task list [--status planned|active|completed]` — prints a formatted table of tasks filtered by status
- [ ] Implement `flow-track task close <task-id>` — marks a task as `completed` with a timestamp annotation
- [ ] Implement `flow-track task abandon <task-id> --reason "..."` — marks as `abandoned` with mandatory reason field

### 7.3 Period Reporting
- [ ] Implement `flow-track report --period 2026-06 [--format json|table|csv]` — generates a period summary filtered by month:
  - Total hours by classification (CapEx, OpEx, Research) per developer
  - Founder sweat-equity hours and dollar valuation (from BLS rates)
  - R&D capitalization ratio: `research_hours / total_hours`
  - DCAA overhead allocation ratio: `opex_hours / (capex_hours + opex_hours)`
- [ ] Implement `flow-track report --ytd` — same as above but for year-to-date
- [ ] Write period summary to `04000-daily-chronicles/<year>/<month>/flowtrack_summary.json` automatically on each `report` run

### 7.4 Auditor-Facing HTML Report Generator
- [ ] Implement `flow-track html --output <path>` — generates a self-contained HTML report file from the `.webnf` database:
  - Summary table: developer × classification breakdown (hours + dollar amounts)
  - Session log table: date, task, description, duration, classification, role
  - Chain integrity status badge (green/red) with hash verification results
  - Founder sweat-equity section with BLS wage citation
  - Generated timestamp and database path in footer
- [ ] Use only Go's `html/template` package (zero external dependencies)
- [ ] Sign the generated HTML file with a Blake3 hash and write a corresponding `.blake3` attestation file alongside it

### 7.5 IDE Integration (Antigravity)
- [ ] Implement `flow-track ide-hook --event branch-switch --branch <name>` — called by the Antigravity IDE when the user switches Git branches; auto-stops any running timer and starts a new session for the branch's associated task
- [ ] Implement `flow-track ide-hook --event file-save --file <path>` — records a lightweight "pulse" event (not a full session) that logs which files were touched during the active session (written as a `pulse` record type in the `.webnf` database)
- [ ] Document the Antigravity IDE hook integration points in `00000-knowledge-foundations/ide_integration_guide.md`

### 7.6 BLS Wage Rate Management
- [ ] Implement `flow-track rates import --source bls_oes_api` — fetches the latest BLS OES wage data for configured SOC codes and states, converting to nano-units and writing to `bls_rates.webnf`
- [ ] Implement `flow-track rates list` — prints the loaded BLS wage table for review
- [ ] Implement annual rate refresh reminder: on `export` or `report`, if `bls_rates.webnf` is older than 365 days, emit a warning `slog.Warn("bls_rates_stale", "age_days", n)`

### 7.7 Promotion Readiness
- [ ] Update [x-flowtrack.swdt.webnf](file:///c:/aCogSpaceSeed/00xper/x-flowtrack/x-flowtrack.swdt.webnf) — set `promotion_ready = "yes"` and `promotion_silo = "00floo"` once all checklist items are complete
- [ ] Run `promote` skill to move `x-flowtrack` → `00floo/s-flowtrack` (or `o-flowtrack`) with updated module namespace
- [ ] Integrate `flow-track.exe` binary build into the `71000-build-harness` sovereign build target
- [ ] Register the promoted workspace in `00flow/s-seed/00200-workspace-taxonomy/taxonomy.txt`

### 7.8 Attestations
- [ ] Expand [tracker_test.go](file:///c:/aCogSpaceSeed/00xper/x-flowtrack/81000-active-source/pkg/tracker/tracker_test.go) to cover:
  - `start` → `stop` round-trip: verify serialized record passes `VerifyChainIntegrity`
  - Tamper simulation: mutate a log record byte and assert `VerifyChainIntegrity` returns false
  - `export` JSON output: verify CapEx/OpEx hour totals match known fixture data
  - Concurrent `start` rejection: two goroutines calling `StartSession` simultaneously, second should fail with "session already active" error
  - Period report accuracy: sessions spanning month boundaries correctly apportioned

## Verification
- `flow-track start`, `stop`, `status`, `verify`, `export` all pass on the live `currenttasks.flowtrack.webnf` database
- `flow-track verify` reports `Success` on the current database
- `flow-track html` generates a valid HTML report with chain integrity badge
- All `tracker_test.go` tests pass including tamper detection
- `int-rehydrator.exe -local-only` passes for `x-flowtrack`
- Grammar validator accepts `currenttasks.flowtrack.webnf` against `flowtrack.wag` with zero violations
