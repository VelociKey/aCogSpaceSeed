# Phase 0: Workspace and Title Alignment

* **Status:** Complete
* **Target Promotion Silo:** `00floo`
* **Primary Workspace:** [o-afflume](file:///c:/aCogSpaceSeed/00floo/o-afflume)
* **Supporting Workspaces:** [x-meter](file:///c:/aCogSpaceSeed/00xper/x-meter), [x-wrapstripe](file:///c:/aCogSpaceSeed/00xper/x-wrapstripe), [x-flowtrack](file:///c:/aCogSpaceSeed/00xper/x-flowtrack)

## Objectives
Establish a clean, taxonomy-conforming foundation for all platform monetization workspaces before any logic is built. This phase ensures that all module namespaces, `go.work` bindings, and directory structures are correctly aligned so downstream phases never encounter import-path or build-graph drift.

## Checklist

### Workspace & Module Namespace Alignment
- [x] Verify directory layouts and package bindings under [go.work](file:///c:/aCogSpaceSeed/go.work)
- [x] Rename `00xper/x-meter` module package namespace to `x-meter-ingestion`
- [x] Rename `00xper/x-wrapstripe` module package namespace to `x-wrapstripe-ledger-adapter`
- [x] Ensure `o-afflume` is registered under `00floo/o-afflume` with the correct `sov.fleet/o-afflume` module path

### Taxonomy Conformance
- [x] Confirm `o-afflume.swdt.webnf` registers workspace against the SWDT taxonomy
- [x] Confirm `x-meter.swdt.webnf`, `x-wrapstripe.swdt.webnf`, and `x-flowtrack.swdt.webnf` are present and valid
- [x] Validate AGENTS.md files are present in all active workspaces

### Sovereign Executable Attestation
- [x] Build initial `o-afflume.exe` binary and record its [Blake3 veracity seal](file:///c:/aCogSpaceSeed/00floo/o-afflume/o-afflume.exe.blake3)
- [x] Build initial `x-meter.exe` binary and record its [Blake3 veracity seal](file:///c:/aCogSpaceSeed/00xper/x-meter/x-meter.exe.blake3)
- [x] Build initial `x-wrapstripe.exe` binary and record its [Blake3 veracity seal](file:///c:/aCogSpaceSeed/00xper/x-wrapstripe/x-wrapstripe.exe.blake3)
- [x] Build initial `x-flowtrack.exe` binary and record its [Blake3 veracity seal](file:///c:/aCogSpaceSeed/00xper/x-flowtrack/x-flowtrack.exe.blake3)

## Verification
- All workspaces resolve under a single `go.work` with zero import-cycle errors.
- SWDT `.webnf` registry files are present for every workspace in scope.
- `int-rehydrator.exe -local-only` builds all workspace binaries without errors.
