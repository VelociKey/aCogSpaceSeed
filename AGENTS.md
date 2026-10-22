# Sovereign Agent Registry & Fleet Governance

> [!IMPORTANT]
> **Authoritative Platform Governance for AI Agents (Gemini & Claude)**  
> This specification defines the absolute mathematical boundaries and operational rails for the VelociKey sovereign software fleet.

---

<hard_invariants>
## Top 10 Zero-Tolerance Hard Invariants

1. **NO RAW GIT:** Never execute raw `git init`, `git add`, `git commit`, or `git push` directly inside silo workspaces. Always use sovereign `vkey construct preserve/add/status/both` toolchain helpers.
2. **NO POWERSHELL & NO SERIAL BATCH-EDITING LOOPS:** Never write or execute PowerShell (`pwsh`) or `.ps1` scripts. For single-file surgical logic edits, use built-in IDE tools (`view_file`, `write_to_file`, `replace_file_content`). For repetitive pattern migrations or multi-file transforms (>= 3 sites), agents are STRICTLY PROHIBITED from running serial `replace_file_content` loops. Always use Go-native batch polyfills (`vkey inspect bash sed`, `vkey inspect scan-replace`) to execute whole-workspace updates in a single round-trip.
3. **NO ROOT README.md:** Never create `README.md` at workspace roots. All documentation must reside strictly in taxonomy directories under `000ALL/cognition/`.
4. **PURE NATIVE GO & NO CGO:** Avoid CGo (C-bindings) and npm/Node.js packages. Use pure Go to eliminate context-switching latency (~80ns) and maintain cross-platform portability.
5. **PREFIX-FREE TAXONOMY & ROLE ENTRY POINTS:** Workspace names must be strictly prefix-free (e.g. `nautilus`, `ingress`, `commerce`). Active source code lives under `81000-active-source/` using role-based entry point names (`<role_name>.go`), explicitly deprecating generic `main.go`.
6. **DECLARATIVE FLUTTER & NO RAW CSS/HTML:** Never edit raw `.css` or `.html` files directly in `web_release/` or deployment directories. All UI must be authored strictly as declarative Dart ($\text{UI} = f(\text{State})$) in `81000-active-source/ui/` and bound to the Go Hub via `dart:js_interop` / SACP (see @.agents/rules/flutter_gohub.md).
7. **SOVEREIGN AGENT SANDBOXING (com.velocikey.vkeygit):** When performing multi-file refactoring, experimental optimizations, or risky AST mutations, agents MUST spin up an isolated scratch worktree using the sovereign skill `com.velocikey.vkeygit` (`C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkey.exe sandbox create` or `vkeygit.exe`). Perform all speculative edits inside `00flow/hydrationcache/c0990-ephemeral-scratch/sandboxes/` before merging back.
8. **NO BUILT-IN SEARCH & ZERO-TRAVERSAL ON ROOT/BINARIES:** Main agents and autonomous subagents are STRICTLY PROHIBITED from calling built-in IDE search tools (`find_by_name`, `grep_search`). Furthermore, broad un-scoped searches across the workspace root (`c:\aCogSpaceSeed`, `.`), `.git/`, `hydrationcache/`, `c0990-ephemeral-scratch/`, or Forge binary directories are forbidden. All codebase navigation and code search must strictly use sovereign `00flow/fab-aides` inspection tools (`vkey inspect findgo`, `vkey inspect bash grep`, `vkey inspect bash find`). Enforced deterministically at runtime by `.agents/hooks.json` via `vkey inspect guard`.
9. **THE NAUTILUS COMPILER OPTIMIZATION PRINCIPLE:** In all generated Go code, agents are strictly forbidden from implementing manual platform-specific SIMD assembly, unrolled intrinsics, or artificial struct padding hacks (`_ [64]byte` pads). Write clean, idiomatic Go data structures and linear memory indexing over raw Plinth memory slices (`as.Bytes()`). The Nautilus Sovereign Three-Border Compiler (`nautilus.exe`) automatically detects host chip architecture (Zen 4 AVX-512, Intel AVX2, ARM NEON) and emits optimal vectorization (see @.agents/rules/compiler_nautilus.md).
10. **NON-INTERACTIVE CLI INVARIANT:** AI agents and subagents must append non-interactive flags (`-q`, `--force`, `-y`) to all executed commands and CLI binaries to prevent hanging in interactive confirmation loops.
</hard_invariants>

---

<prohibited_executables>
## Prohibited Executables & Hook Guard Policy

- **NO PYTHON / SCRIPT INTERPRETERS:** Agents must NEVER run `python`, `python3`, or `.py` scripts autonomously (rare override: `SOVEREIGN_ALLOW_PYTHON=1` or `vkey inspect exec-python`).
- **NO DIRECT GO / RAW TOOLCHAINS:** Agents must NEVER run raw `go build`, `go test`, or `go run` outside of sovereign `realize.exe`.
- **NO BUILT-IN FD/RIPGREP TOOLS:** Built-in `find_by_name` and `grep_search` are permanently blocked by runtime hook guards.
</prohibited_executables>

---

<tool_decision_matrix>
## Sovereign Tool Action Decision Matrix

| Action Required | Sovereign Command / Tool |
| :--- | :--- |
| **Build / Test Workspace** | `C:\aCogSpaceSeed\00flow\forge\96000-internal-executables\realize.exe -workspace <name>` |
| **Locate Go Files / Types in $O(1)$** | `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkey.exe inspect findgo -query <name>` |
| **Search Code Across Workspaces** | `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkey.exe inspect bash grep -r <pattern> <dir>` |
| **Multi-File Batch Refactor** | `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkey.exe inspect bash sed -i "s/old/new/g" <files>` |
| **Codebase Conformance / Linter Audit** | `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkey.exe inspect slog-check / audit-imports` |
| **Cross-Platform Shell Operations** | `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkey.exe inspect bash <ls\|cat\|grep\|find\|sed\|head\|tail\|cp\|mv\|rm\|mkdir>` |
| **Git / Versioning Operations** | `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkey.exe construct <preserve\|add\|status\|both>` |
| **Formal Prover & Verification** | `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkey.exe prove` |
</tool_decision_matrix>

---

<brain_and_memory_governance>
## The 3-Tier Sovereign Brain & Hydrationcache Architecture

1. **Tier 1: Runtime Session Brain (`%USERPROFILE%\.gemini\antigravity\brain\<id>\`):**
   - Machine-local session storage containing transcripts (`transcript.jsonl`), background task logs, and user-facing artifacts.
   - All internal temporary scratch scripts and data MUST be routed to `<appDataDir>\brain\<id>\scratch\` with `UserFacing: false`.
   - Indexed in Antigravity v2.21.1 for ⌘+K (Ctrl+K) full-text conversation search.
2. **Tier 2: Quantum-Safe Cognitive DAG Memory (`vkeygit memory`):**
   - Checkpointed AST diffs, task snapshots, and benchmark scorecards cryptosealed into content-addressed memory in `00flow/hydrationcache/c0990-ephemeral-scratch/sandboxes/`.
3. **Tier 3: Institutional Fleet Cognition (`000ALL/cognition/`):**
   - Permanent monorepo repository of blueprints, mathematical contributions, compendiums (`c0411-compendium`), and sovereign narratives committed to Git `zenith` branches.
4. **Disposable Hydrationcache Principle:**
   - `00flow/hydrationcache/c0990-ephemeral-scratch/` is a **purely disposable acceleration cache**.
   - If purged, the fleet suffers **ZERO functional impact**; all source code and attestation integrity remain intact.
</brain_and_memory_governance>

---

<agent_roles>
## Sovereign Agent Roles & Core Protocols

- **Conductor:** Orchestrator of multi-step workflows. Definitions and state reside in `/conductor`.
- **Antigravity (IDE Agent):** Resident pair-programming assistant within Antigravity IDE. For multi-site edits, dispatches Go-native batch polyfills (`vkey inspect bash sed`) rather than serial loops.
- **Jules (Engineering Agent):** Autonomous software engineer operating in the cognitive layer. Employs the **Jules Repoless Loop** with a strict **3-attempt autonomous healing limit**.
- **NATVS Engine:** Fleet-aware lifecycle engine governing Negotiation, Assimilation, Transformation, Verification, and Synthesis.
</agent_roles>

---

<domain_rules_index>
## Modular Domain Rules Index

For specialized domain contracts, consult the authoritative rules:
- **Nautilus Compiler Optimization & Borders:** @.agents/rules/compiler_nautilus.md
- **Family Governance & 12-Step Promotion:** @.agents/rules/family_governance.md
- **Declarative Flutter & Go Hub SACP:** @.agents/rules/flutter_gohub.md
- **Concurrency & Plinth Substrate Arenas:** @.agents/rules/concurrency_plinth.md
- **Ephemeral Scratch Sandboxing & Memory:** @.agents/rules/vkeygit_sandbox.md
</domain_rules_index>