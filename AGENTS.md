# Sovereign Agent Registry & Fleet Governance

> [!IMPORTANT]
> **TOP 5 ZERO-TOLERANCE HARD INVARIANTS FOR AI AGENTS (GEMINI & CLAUDE):**
> 1. **NO RAW GIT:** Never execute raw `git init`, `git add`, `git commit`, or `git push` directly inside silo workspaces. Always use sovereign `vkey construct preserve/add/status/both` toolchain helpers.
> 2. **NO POWERSHELL:** Never write or execute PowerShell (`pwsh`) or `.ps1` scripts. Always use built-in IDE agent tools (`view_file`, `write_to_file`, `replace_file_content`, `grep_search`, `list_dir`) and Go-native `vkey` toolchain commands for all filesystem operations.
> 3. **NO ROOT README.md:** Never create `README.md` at workspace roots. All documentation must reside strictly in taxonomy directories under `000ALL/cognition/`.
> 4. **PURE NATIVE GO & NO CGO:** Avoid CGo (C-bindings) and npm/Node.js packages. Use pure Go to eliminate context-switching latency (~80ns) and maintain cross-platform portability.
> 5. **PREFIX-FREE TAXONOMY & ROLE ENTRY POINTS:** Workspace names must be strictly prefix-free (e.g. `pubsub`, `firestore`, `neuro-wholecell`). Active source code lives under `81000-active-source/` using role-based entry point names (`<role_name>.<ext>`), explicitly deprecating generic `main.<ext>`.
> 6. **DECLARATIVE FLUTTER & NO RAW CSS/HTML:** Never edit raw `.css` or `.html` files directly in `web_release/` or deployment directories. All UI must be authored strictly as declarative Dart ($\text{UI} = f(\text{State})$) in `81000-active-source/ui/` and bound to the Go Hub via `dart:js_interop` / SACP (see [flutter_gohub_declarative_pipeline.md](file:///c:/aCogSpaceSeed/.agents/guidelines/flutter_gohub_declarative_pipeline.md)).
> 7. **SOVEREIGN AGENT SANDBOXING (com.velocikey.vkeygit):** When performing multi-file refactoring, experimental optimizations, or risky AST mutations, agents MUST spin up an isolated scratch worktree using the sovereign skill `com.velocikey.vkeygit` (`C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkey.exe sandbox create` or `vkeygit.exe`). Perform all speculative edits and test verifications inside `c0990-ephemeral-scratch/sandboxes/` before merging back, ensuring zero active workspace corruption.

> [!CAUTION]
> **PROHIBITED EXECUTABLES & HOOK GUARD POLICY (ANTIGRAVITY v2.7.1):**
> - **NO PYTHON / SCRIPT INTERPRETERS:** Agents must NEVER run `python`, `python3`, or `.py` scripts autonomously. (Rare explicit user override: Require user token `SOVEREIGN_ALLOW_PYTHON=1` or `vkey inspect exec-python`).
> - **NO DIRECT GO / RAW TOOLCHAINS:** Agents must NEVER run raw `go build`, `go test`, or `go run` outside of sovereign `realize.exe`.
> - **TOOL ACTION DECISION MATRIX:**
>   - *Action: Build / Test Workspace* $\rightarrow$ `C:\aCogSpaceSeed\00flow\forge\96000-internal-executables\realize.exe -workspace <name>`
>   - *Action: Ephemeral Scratch Sandboxing / Memory* $\rightarrow$ `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkey.exe sandbox / memory` (or `vkeygit.exe`)
>   - *Action: Codebase Search / AST Audit* $\rightarrow$ Built-in `grep_search` OR `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkey.exe inspect`
>   - *Action: Git / Versioning & Workspace Scaffolding* $\rightarrow$ `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkey.exe construct <subcommand>`
>   - *Action: Formal Prover & Verification* $\rightarrow$ `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkey.exe prove` (`com.velocikey.formalprover`)
>   - *Hook Enforcement* $\rightarrow$ Enforced at runtime by [`.agents/hooks.json`](file:///C:/aCogSpaceSeed/.agents/hooks.json) via `vkey inspect guard`.

This workspace is orchestrated using the **Antigravity Agent Manager** and the **Gemini CLI**.

---

## 1. The Cognitive Linguistic Foundation of 00flow

Our entire software and agent development platform leads with formal grammars, strict semantic/lexical taxonomies, and Domain Specific Languages (DSLs) to create a perfect **neuro-symbolic resonance** with Large Language Models (LLMs and SLMs):

1. **Deterministic Cognitive Rails:** System configurations, build harnesses, and workspace promotions are restricted to non-recursive, flat-block DSL grammars (`.wag` files). This eliminates AI agent hallucination and state-tracking errors by representing the absolute mathematical boundary of valid intent.
2. **LatentLingua (Linguistic Blood Stream):** weBNF (Witte's / White's Evolved Backus-Naur Form grammar) leverages the pre-existing training weights of all LLMs/SLMs (which natively understand EBNF via compiler specifications and RFCs). Agents ingest `.wag` grammars and produce `.webnf` programs with zero-shot compliance, replacing imprecise natural language documentation with compiler-validated contracts. Core grammars reside in **Die** containers for system storage and **Mold** containers for licensee distribution (see [grammar_storage_containers.md](file:///c:/aCogSpaceSeed/.agents/guidelines/grammar_storage_containers.md)).
3. **First-Principles Semantic Abstraction Naming Rule:** All Domain-Specific Language (DSL) hints, kernel primitives, and mathematical contracts MUST be named using abstract, first-principles domain-semantic terminology (*e.g., `TOPOLOGICAL_MORTON_PARTITION`, `SPARSE_PRECISION_TILED_SOLVER`, `ALGEBRAIC_REGISTER_FUSION`, `ZERO_COPY_WAVELET_BEAMFORMING`*) rather than legacy software packages or commercial tool names (*avoiding legacy names like GROMACS, LAMMPS, AUTODOCK, RELION*). This guarantees zero semantic drift, enables zero-shot cross-fertilization across scientific domains, and enforces pure mathematical lowering in compilers like `wSTMulc`. Application workspace names MUST NEVER be prepended with artificial single-letter prefixes (such as `q`) unless `q` is part of the original upstream application name (e.g., `qiskit`). Always use strictly prefix-free domain semantic names.
4. **Lexical Directory Taxonomy:** Our standardized 5-digit/4-digit directory structure acts as a physical lexical grammar, enabling LLMs to navigate the fleet with $O(1)$ search latency and zero semantic drift.
5. **Syntactic Supply-Chain Auditing:** Formal language contracts enable specialized agents to audit system architecture, cryptoseals, and dependencies purely at the syntactic and AST level, guaranteeing air-gapped zero-trust execution.
6. **LatentFlow Neuro-Symbolic Grammar Standard (STAGED / INACTIVE):** LatentFlow is STAGED and INACTIVE by default. AI agents MUST NOT emit `nsle.webnf` program records unless explicitly requested by the user. When active, agents WILL utilize the `LatentFlow` WSN grammar ([nsle.wag](file:///c:/aCogSpaceSeed/00xper/latentflow/nsle.wag)) and `wSTMulc` universal compiler engine to offload heavy calculations via compact `webnf.sn` program records (`<name>-nsle.webnf`). Refer to [latentflow_usage_guideline.md](file:///c:/aCogSpaceSeed/.agents/guidelines/latentflow_usage_guideline.md) for operational rules.

---

## 2. Agent Roles & Core Protocols

### Agent Roles

#### Conductor
* **Role:** Orchestrator of multi-step workflows.
* **Responsibility:** Manages task delegation, progress tracking, and cross-agent communication.
* **Context:** Workflow definitions and state reside in the `/conductor` directory.

#### Antigravity (IDE Agent)
* **Role:** Resident pair-programming assistant within the Antigravity IDE.
* **Responsibility:** Code editing, terminal management, local file operations, and real-time pair programming.
* **Connection:** Linked to the Gemini CLI via `/ide install`.

#### Jules (Engineering Agent)
* **Role:** Autonomous software engineering agent.
* **Responsibility:** End-to-end task execution including writing/running unit tests, algorithmic research, implementation refinement, and sovereign architectural critique.
* **Function:** Operates in the 94xxx Cognitive Layer (Actors) via the Jules CLI.
* **Tools & Polyfills:** Equipped with the Go-native `inspect` tool twin (from `vkey inspect` / `fab-aides`) to perform in-process filesystem scans and regex pattern matching without subprocess overhead.

#### NATVS Engine (Orchestration Actor)
* **Role:** Fleet-aware execution engine.
* **Responsibility:** Manages the NATVS lifecycle (Negotiation, Assimilation, Transformation, Verification, Synthesis).
* **Function:** Orchestrates multi-silo tasks using the Jules Agent and sovereign toolchains.
* **Interface:** Invoked via Gemini-CLI or the Antigravity Agent Manager.

### Core Protocols

#### NATVS Protocol
All agentic interactions within the fleet MUST follow the 5-phase NATVS lifecycle:
1. **N - Negotiation:** Handshake on the Whisper Bus (Capability / Price / Time).
2. **A - Assimilation:** Aggregate context ingestion from multiple silos into the agent's reasoning plane.
3. **T - Transformation:** Execute computational work (Critic, Creator, Fixer).
4. **V - Verification:** Validate results modularly (Test/Pass, Complexity delta).
5. **S - Synthesis:** Merge fleet state and metabolically prune ephemeral artifacts.

#### The Jules Repoless Loop
When integrating the Jules Agent:
1. **No Git Remotes:** Do not use `git push` or `git pull` for Jules context.
2. **Repoless Snapshot:** Use `jules remote new --repo .` to upload ephemeral context.
3. **Local Sync:** Always use `jules remote pull` to merge changes back to local disk.
4. **Autonomous Healing & Bailout:** Run the `Fixer` loop (Test $\to$ Fix $\to$ Test) with a hard limit of **3 attempts**. If objectives are not met or if performance stalls, stop and report findings immediately.
   * **Success Criteria:** ALL of: (a) `go build ./...` succeeds, (b) `go test ./... -count=1` passes with 0 failures, (c) `go vet ./...` reports no issues, (d) no new `panic` calls introduced in library packages.
   * **Partial Success:** If attempt 3 achieves (a)+(b) but fails (c)+(d), report as *"builds and tests pass with warnings"* — do NOT continue fixing.
5. **Escape Analysis Tooling:** Run `go build -gcflags="-m"` or `go test -gcflags="-m" ./...` on modified packages to identify heap allocations empirically.

---

## 3. Platform Structure & Silo Architecture

### 000ALL (The Cognitive Silo)
* **Cognition:** Central authority for architectural blueprints, master area indexes, and agent narratives.

### 00flow (SDLC Platform Silo)
* **aCogSpaceSeed:** Bootstrap workspace for platform initialization and taxonomy management.
* **Forge:** Central repository for internal binaries and executables (`forge/96000-internal-executables`). For the licensee binary distribution model (**Foundry**), see [grammar_storage_containers.md](file:///c:/aCogSpaceSeed/.agents/guidelines/grammar_storage_containers.md).
* **LatentLingua:** The "Grammar Studio" for developing DSLs using simplified eBNF/WSN.

### 00AAIF (Agent Architecture & Integration Framework)
* Modernizes and standardizes agentic specifications across the fleet.

### 00floo (Operations & Domain Logic Silo)
* Contains domain-specific logic and operations workspaces using strictly **prefix-free domain semantic names** (e.g., `firestore`, `pubsub`, `afflume`).

### 00flon (Network & Routing Silo)
* Contains network routing, SACP driver, and transport workspaces (e.g., `ingress`, `sovereign-ingress` alias).

### boundary-conduit (Integration Gateway Class)
* Serves as the boundary layer between external protocol fits and internal core engines (e.g., network translation process `janus` and REST/webhook API receptors).

### 86SREF (Reference Silo)
* Contains read-only reference clones from the **VelociKey** GitHub organization.
* **Mandate:** This silo is strictly **read-only**. Never modify files here; migrate code to an active workspace before adaptation.

---

## 4. Technology Stack & Sovereign Toolchain Paths

### Technology Stack Constraints
* **UI Layer:** Dart & Flutter (Material Design 3). Web interaction with the Go backend MUST flow through the **GO Hub** using modern `dart:js_interop` / `package:web` and raw `golang.org/x/net/quic` / SACP transport connections, bypassing HTTP REST polling and raw CSS/HTML editing entirely. UI is strictly declarative Dart ($\text{UI} = f(\text{State})$). Refer to [flutter_gohub_declarative_pipeline.md](file:///C:/aCogSpaceSeed/.agents/guidelines/flutter_gohub_declarative_pipeline.md) and [flutter_web_constraints.md](file:///C:/aCogSpaceSeed/000ALL/00000-knowledge-foundations/120-system-architecture/compilation-and-builds/flutter_web_constraints.md).
* **Backend Layer:** Pure Go (without CGo context-switching latency).
* **Networking Layer:** Native Pure Go QUIC (`golang.org/x/net/quic`), qAPC (QUIC Actor Procedure Call), and WebTransport over physical UDP sockets, or zero-copy browser shared memory. The entire fleet is graduated to native `golang.org/x/net/quic` for zero-allocation connection migration, multi-stream multiplexing, and sub-microsecond actor procedure calls.
* **Logic Layer:** Domain Specific Languages (DSLs) based on Wirth Syntax Notation.
* **Compilation & Build Engine:** **Nautilus Sovereign Three-Border Compiler** (`nautilus.exe`), replacing legacy `wSTMulc`/`wtulc`. Orchestrated by `realize.exe`. Run workspace tests using `realize.exe -workspace <path>`.
* **Sovereign Reanimation Gateway:** **`seerSTMci`** (`seerstmci_driver`), the authoritative Scientific Transformational Mathematics engine for soul extraction, evolution, and continuous reanimation of foreign/upstream code into zero-allocation pure Go actor workspaces via 6-Pillar Rosetta Functors.

### Authoritative Sovereign Toolchain Paths
Always invoke sovereign toolchains directly using absolute workspace paths:
* **VelociKey Toolchain Engine (vkey℠):** `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkey.exe`
* **Sovereign Three-Border Compiler (Nautilus):** `C:\aCogSpaceSeed\00flow\forge\96000-internal-executables\nautilus.exe`
* **Sovereign Reanimation Gateway (seerSTMci):** `C:\aCogSpaceSeed\00flow\forge\96000-internal-executables\seerstmci`
* **Realize Workspace Builder:** `C:\aCogSpaceSeed\00flow\forge\96000-internal-executables\realize.exe`
* **Sovereign Go Executable:** `C:\aCogSpaceSeed\00flow\forge\92000-external-toolchains\go\bin\go.exe`
* **Formal Prover Engine (Lean 4):** `C:\aCogSpaceSeed\00flow\forge\92000-external-toolchains\lean4\bin\lean.exe`
* **Sovereign Ephemeral Scratch:** `C:\aCogSpaceSeed\c0990-ephemeral-scratch`

---

## 5. Workflow Conventions & Zero-Trust Layout Rules

### Workflow Conventions
* **Task Tracking:** Document multi-step tasks under `/conductor/tasks/`.
* **Taxonomy Provisioning:** Use `vkey construct both -name <name> -silo <silo>` to provision a new workspace, Git repository, and remote GitHub twin under `VelociKey`. Align taxonomy using `lingualattice align -WorkspacePath <path>`.
* **Git Operations:** Never execute raw Git commands (`git init`, `git add`, `git commit`, `git push`). Always use `vkey construct` helpers:
  * Commit and push tracking changes: `vkey construct preserve -silo <silo>`
  * Stage local changes: `vkey construct add`
  * Inspect workspace status: `vkey construct status`
* **Reactive Execution:** When executing background tasks, do not poll or use sleep timers. Yield control (call no more tools) and wait for system notification.
* **Workspace Boundary Discovery:** Workspace boundaries and roots are governed dynamically by **`go.work`** and **`.workspace-dna.webnf`** (Antigravity v2.10.0+). Legacy `.gitroot` is permanently retired.
* **Silo Promotion Policy:** Promoting a workspace from experimental (`00xper`) to production (`00flow`) requires explicit, independent user approval.

### Zero-Trust Taxonomy & Active Source Layout Rules
* **Anti-Standard Layout Rule:** Never create standard root files (e.g. `README.md`).
* **Root Exception List:** Only `go.work`, `go.work.sum`, `go.mod`, `go.sum`, `AGENTS.md`, `.gitignore`, and `.workspace-dna.webnf` are authorized at workspace roots.
* **Language-Agnostic Active Source Layout:** Active source code lives under `81000-active-source/` using a flat layout:
  * **Role-Based Entry Points:** Executable entry points reside directly at the root of `81000-active-source/` named after their executable role, e.g. `<role_name>.<ext>` (e.g., `pubsub_supervisor_gateway.go`). Generic `main.<ext>` names are deprecated.
  * **Domain Package Subdirectories:** Modules reside directly under `81000-active-source/` in subdirectories named strictly after their architectural domains (e.g., `/transport`, `/calculus`).
  * **Prohibited Directories:** Generic boilerplate folders (`cmd/`, `pkg/`, `internal/`, `src/`, `lib/`) are prohibited.
* **Diagram Standard:** Agents MUST render architectural and sequence flow diagrams using GitHub-Flavored Markdown **Mermaid Diagrams** (`mermaid`).

---

## 6. Optimization Standards & Protocol Governance

### Optimization & Code Maturity Standards
* Refer to [go_maturity_standards.md](file:///C:/aCogSpaceSeed/000ALL/00000-knowledge-foundations/120-system-architecture/language-maturity/go_maturity_standards.md) for concurrency lock selections (`adaptivelock`) and compilation optimizations (`discard`).
* Refer to [go_standards.md](file:///c:/aCogSpaceSeed/.agents/guidelines/go_standards.md) when writing or modifying Go files containing concurrency primitives (`sync.Mutex`, `sync.RWMutex`, `sync.Map`, `chan`, `go`).
* **Automatic Double-Performance Standard:** Analyze generated or modified code to double execution speed and reduce GC allocations by $\ge 50\%$. Achieve this via stack allocation, buffer pooling, pre-allocated capacities, and eliminating reflection, `fmt.Sprintf` hex formatting, and intermediate string-to-byte casting.
  * **Measurement Protocol:** Run `go test -bench=. -benchmem` before and after changes on modified packages.
  * **Exception Rule:** If the modified code is already allocation-free or operates at $< 100\text{ ns/op}$, escape analysis review (`go build -gcflags="-m"`) alone satisfies the standard.

### Hybrid Adaptive SACP Protocol Standard
All client/server transactions MUST use the **Hybrid Adaptive SACP Protocol**:
* **Browser-Local (Same-Memory):** Routes frames directly through zero-copy shared memory offsets utilizing Wasm BI/IR (Bidirectional Intercept / Inbound-Outbound Routing) to execute the server module natively in the client process.
* **Networked (Socket):** Routes frames over physical SACP UDP/QUIC network sockets to remote host servers.
* Deploy builds directly to `/96000-internal-executables/web_release` without introducing legacy directory structures.