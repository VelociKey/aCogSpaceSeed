# Phase 2: sTelemetry & qAPC ValueVector Receiver

* **Status:** Complete
* **Target Promotion Silo:** `00floo`
* **Primary Workspace:** [meter](file:///c:/aCogSpaceSeed/00xper/meter)
* **Integration Target:** [afflume/qapc](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/qapc)

## Objectives
Deploy the `sTelemetry` metrics spooler and the qAPC binary handler so that VPU consumption events from the Antigravity IDE and fleet workspaces are continuously streamed, locally spooled to `.webnf` ValueVector files, and exposed to the receptor pricing engine in Phase 3.

The `MeterReporter` in `meter` acts as the **producer**: it writes spool files. The `qAPC` receiver in `afflume` acts as the **consumer**: it reads those spool files, validates them against the `meterdsl.wag` grammar, and hands off structured `ValueVector` batches to the offering resolver.

## Current State (Audit)
- [meter.go](file:///c:/aCogSpaceSeed/00xper/meter/81000-active-source/pkg/meter/meter.go) — `MeterReporter` with `adaptivelock.Map` buffering, per-session spool file writes in `meter.<sessionUUID>.spool.webnf` format
- [meterdsl.wag](file:///c:/aCogSpaceSeed/00xper/meter/meterdsl.wag) — grammar defining the `sTelemetry` ValueVector record schema
- [sample_telemetry.stelemetry.webnf](file:///c:/aCogSpaceSeed/00xper/meter/sample_telemetry.stelemetry.webnf) — example conforming program
- [qapc/client.go](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/qapc/client.go) — qAPC client stub in `afflume`
- [qapc/receiver.go](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/qapc/receiver.go) — `SpoolWatcher` spool watcher & `ValueVector` aggregator
- [qapc/receiver_test.go](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/qapc/receiver_test.go) — Spool watcher integration test suite

## Checklist

### 2.1 sTelemetry Grammar Finalization
- [x] Review and finalize [meterdsl.wag](file:///c:/aCogSpaceSeed/00xper/meter/meterdsl.wag) to ensure it covers:
  - `workspace_id`, `user_uuid`, `session_uuid`, `metric_name`, `count`, `timestamp` fields
  - Cryptographic `passport` attestation field (links spool record to a Gatekeeper-signed session token)
  - Spool file header block: `version`, `schema_ref`, `genesis_hash`
- [x] Validate `sample_telemetry.stelemetry.webnf` conforms to finalized `meterdsl.wag`
- [x] Define `stelemetry_aggregate.wag` grammar for the rolled-up period ValueVector (hourly/daily aggregates consumed by Phase 3)

### 2.2 MeterReporter Hardening
- [x] In-memory metric buffer with `adaptivelock.Map` — implemented
- [x] Session spool file writer — implemented
- [x] Implement `FlushToSpool(ctx context.Context)` method — serializes the current in-memory buffer to `meter.<sessionUUID>.spool.webnf` with a Blake3 hash linking to the prior flush record (making the spool itself a verifiable chain)
- [x] Implement `AutoFlush(ctx context.Context, interval time.Duration)` — wraps `FlushToSpool` in a background ticker, called on clean session stop
- [x] Implement platform lock file (`meter.<sessionUUID>.lock`) using the platform-native file lock (`lock_windows.go` / `lock_unix.go`) to prevent concurrent spool corruption from multiple IDE processes
- [x] Add `meter_test.go` coverage: concurrent `Increment()` calls, flush correctness, lock collision

### 2.3 qAPC ValueVector Receiver (afflume)
- [x] Implement `qapc/receiver.go` in [afflume](file:///c:/aCogSpaceSeed/00floo/afflume/81000-active-source/qapc) — a directory watcher that monitors the meter spool directory for new `.spool.webnf` files
- [x] Parse each spool file against `meterdsl.wag` grammar via regex / LatentLingua codec; reject and quarantine malformed files
- [x] Aggregate parsed metric records per `workspace_id` + `user_uuid` into a `ValueVector` struct:
  ```go
  type ValueVector struct {
      WorkspaceID string
      UserUUID    string
      PeriodStart time.Time
      PeriodEnd   time.Time
      Metrics     map[string]int64  // metric_name -> total count
  }
  ```
- [x] Emit `ValueVector` batches to an in-process channel consumed by the receptor pricing engine (Phase 3)
- [x] Archive processed spool files to `04000-daily-chronicles/<date>/` after successful parsing; delete originals
- [x] Implement receiver attestation: log `slog.Info` with `spool_file`, `record_count`, `hash`, `parse_elapsed_ms` for each processed file

### 2.4 DSL Grammar Compliance Enforcement
- [x] Wire LatentLingua codec validation from `meterdsl.wag` to produce Go structs consumed by the receiver
- [x] Add a startup self-check: if `meterdsl.wag` schema hash doesn't match the compiled codec hash, abort with a clear `slog.Error` message

### 2.5 Attestations
- [x] Add `qapc/receiver_test.go` covering:
  - Lock collision simulation (two goroutines writing to same spool)
  - Flush → parse round-trip correctness
  - Malformed spool file quarantine behavior
  - ValueVector aggregation accuracy against known fixture data

## Verification
- `int-rehydrator.exe -local-only` passes for `x-meter`
- Sample spool file parses cleanly via `fab-inspect.exe bash grep` scan against `meterdsl.wag` schema
- End-to-end flow: `MeterReporter.Record()` → `FlushToSpool()` → `receiver.go` parses → `ValueVector` emitted with correct counts
