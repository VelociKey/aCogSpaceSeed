# Guideline: Sovereign Environment Preparation & Toolchains

## 1. Environment Preparation
1. Ensure `antigravity-cli (agy)` is installed and running.
2. **Redirect Metadata:** To ensure `agy`, Go, and Bazel do not use the user home directory or OS toolchains for temporary state/execution, set the following environment variables:
   - `AGY_HOME=C:\aCogSpaceSeed\c0990-ephemeral-scratch`
   - `GOCACHE=C:\aCogSpaceSeed\00flow\hydrationcache\c0990-ephemeral-scratch\go-cache`
   - `GOTMPDIR=C:\aCogSpaceSeed\00flow\hydrationcache\c0990-ephemeral-scratch\go-tmp`
   - `JAVA_HOME=C:\aCogSpaceSeed\00flow\forge\92000-external-toolchains\jdk`
   - Configure `PATH` hermetically: do NOT inherit or prepend to the host OS `%PATH%`. Construct it strictly using only the paths to the sovereign toolchains inside `forge` (e.g. `jdk/bin`, `go/bin`) and core system directories (e.g. `C:\Windows\system32`).
3. Run `agy install` inside the terminal to link with the Antigravity IDE.
4. Use the **Agent Manager** panel in Antigravity to interact with these roles.
5. **Hermetic Credentials Isolation:** Sandboxes and task execution runners MUST NEVER bridge, copy, or read host OS user credentials, configurations, or profiles (e.g., `.gitconfig`, `.ssh`, `gcloud`, or other credential caches). All sandboxed executions must run strictly decoupled from the host user's secrets to prevent accidental leakage into generated or distributed codebases/binaries.

## 2. Authoritative Sovereign Toolchain Paths
To prevent background command and filesystem search bottlenecks, all agents must utilize the following paths directly:
- **VelociKey Deterministic Toolchain (vkey℠):** `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkey.exe`
  - *Filesystem Polyfills & Operations:* `vkey.exe inspect bash <cmd>` (e.g. `rm`, `cp`, `mv`, `ls`, `mkdir`, `cat`)
  - *Git Versioning & Scaffolding:* `vkey.exe construct <subcommand>` (e.g. `preserve`, `add`, `status`, `both`)
  - *Codebase AST Inspection:* `vkey.exe inspect <subcommand>`
- **Sovereign Go Executable:** `C:\aCogSpaceSeed\00flow\forge\92000-external-toolchains\go\bin\go.exe`
- **Sovereign Bazel Universal Builder:** `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\int-rehydrator.exe`
- **Sovereign Workspace Builder & Tester (realize):** `C:\aCogSpaceSeed\00flow\forge\96000-internal-executables\realize.exe`
- **Formal Prover Engine (Lean 4 Microkernel):** `C:\aCogSpaceSeed\00flow\forge\92000-external-toolchains\lean4\bin\lean.exe`
- **Sovereign Ephemeral Scratch:** `C:\aCogSpaceSeed\c0990-ephemeral-scratch`
- **Cognitive Layer Silo (000ALL):** `C:\aCogSpaceSeed\000ALL`
- **External Model Catalog:** `C:\aCogSpaceSeed\00flow\hydration\84500-external-models-source\models.manifest.webnf`
- **External Toolchain Catalog:** `C:\aCogSpaceSeed\00flow\hydration\400-registry\toolchains.manifest.webnf`
