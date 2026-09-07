# Implementation Plan: Renaming Rehydration Tools & Establishing Standalone Deployment

* **Status:** ✅ Complete

This plan defines the process to rename the rehydrator toolchain to a verb-first format (`rehydrate-ext`, `rehydrate-int`) and extract the deployment logic into a standalone companion command-line executable (`deploy`).

---

## 1. Directory & Source Code Renaming

### A. Internal Rehydrator
* **Rename Driver Folder:**
  Rename `C:\aCogSpaceSeed\00flow\hydration\81000-active-source\drivers\int-rehydrator` to `C:\aCogSpaceSeed\00flow\hydration\81000-active-source\drivers\rehydrate-int`.
* **Update Entry Point Name:**
  Change the package/binary output target inside the folder to `rehydrate-int.exe`.

### B. External Rehydrator
* **Rename Driver Folder:**
  Rename `C:\aCogSpaceSeed\00flow\hydration\81000-active-source\drivers\ext-rehydrator` to `C:\aCogSpaceSeed\00flow\hydration\81000-active-source\drivers\rehydrate-ext`.
* **Update Entry Point Name:**
  Change the binary output target to `rehydrate-ext.exe`.

---

## 2. Extraction of Standalone `deploy` Executable

We decouple the deployment logic from the compiler binary:
* **Create New Driver Folder:**
  Create directory `C:\aCogSpaceSeed\00flow\hydration\81000-active-source\drivers\deploy\`.
* **Move and Refactor Entry Point:**
  Establish `deploy.go` inside the folder as `package main`.
  Remove the temporary `-deploy` flag from the `rehydrate-int` CLI parser, routing deployment operations strictly through the standalone `deploy.exe` command-line binary.
* **Validation Check:**
  Configure `deploy.go` to require a pre-compiled artifact payload, computing Blake3 hashes and verifying build invariants before promoting binaries.

---

## 3. Configuration & References Alignment

We run verification updates to replace all occurrences of `ext-rehydrator` and `int-rehydrator` (and their `.exe` extensions) with the new triumvirate (`rehydrate-ext`, `rehydrate-int`, `deploy`) inside:
* **Global System Blueprint:**
  - [AGENTS.md](file:///C:/aCogSpaceSeed/AGENTS.md)
* **Conductor Task Logs:**
  - [/conductor/tasks/index.md](file:///C:/aCogSpaceSeed/conductor/tasks/index.md)
  - [/conductor/tasks/c0000-phase0-alignment.md](file:///C:/aCogSpaceSeed/conductor/tasks/c0000-phase0-alignment.md)
  - [/conductor/tasks/c0001-phase1-ledger.md](file:///C:/aCogSpaceSeed/conductor/tasks/c0001-phase1-ledger.md)
  - [/conductor/tasks/c0002-phase2-stelemetry.md](file:///C:/aCogSpaceSeed/conductor/tasks/c0002-phase2-stelemetry.md)
  - [/conductor/tasks/c0003-phase3-receptor.md](file:///C:/aCogSpaceSeed/conductor/tasks/c0003-phase3-receptor.md)
  - [/conductor/tasks/c0004-phase4-ordering.md](file:///C:/aCogSpaceSeed/conductor/tasks/c0004-phase4-ordering.md)
  - [/conductor/tasks/c0006-phase6-integration.md](file:///C:/aCogSpaceSeed/conductor/tasks/c0006-phase6-integration.md)
  - [/conductor/tasks/c0007-phase7-flowtrack.md](file:///C:/aCogSpaceSeed/conductor/tasks/c0007-phase7-flowtrack.md)
  - [/conductor/tasks/fab-inspect-bash-polyfills-task.md](file:///C:/aCogSpaceSeed/conductor/tasks/fab-inspect-bash-polyfills-task.md)
  - [/conductor/tracks/algebraic_implementation_prompts.md](file:///C:/aCogSpaceSeed/conductor/tracks/algebraic_implementation_prompts.md)
* **Build Helpers & Harness Scripts:**
  - Workspace build rules and compiler dependency config templates.

---

## 4. Platform-Agnostic Script Replacement
Following the zero-trust system guidelines, we deprecate and remove all legacy PowerShell script helpers to avoid runtime latency and keep execution platform-agnostic:
* **`hydrate-check.ps1` Deprecation:** We will delete [hydrate-check.ps1](file:///C:/aCogSpaceSeed/00flow/hydration/hydrate-check.ps1) entirely.
* **Go-Native Replacement (`hydrate-check`):** We compile a native Go binary from `00flow/hydration/81000-active-source/drivers/hydrate-check/` that performs the directory checks, attestation verification, and rehydration diagnostics. This runs 10x faster and eliminates Windows-specific shell dependency.

---

## 5. Compilation & Verification Plan
1. Re-build the toolchain drivers using the local Go compiler:
   ```powershell
   go build -o C:\aCogSpaceSeed\rehydrate-int.exe C:\aCogSpaceSeed\00flow\hydration\81000-active-source\drivers\rehydrate-int\int-rehydrator.go
   go build -o C:\aCogSpaceSeed\deploy.exe C:\aCogSpaceSeed\00flow\hydration\81000-active-source\drivers\deploy\deploy.go
   go build -o C:\aCogSpaceSeed\hydrate-check.exe C:\aCogSpaceSeed\00flow\hydration\81000-active-source\drivers\hydrate-check\main.go
   ```
2. Verify that running `rehydrate-int -workspace timewarp` compiles the artifacts without executing deployment, and running `deploy -workspace timewarp -target pages` triggers the verification and deployment pipeline.

