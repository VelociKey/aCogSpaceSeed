---
name: com.velocikey.vkeygit
displayName: vkeygit
description: >-
  Sovereign ephemeral scratch sandboxing, topological multi-repo merge-back, and quantum-safe cognitive DAG memory for Antigravity V2 LLM Agents (Gemini & Claude).
  Use to create isolated scratch worktrees for speculative coding and testing, merge back verified changes, or maintain branching cognitive memory without context token bloat.
---

# vkeygit: Ephemeral Sandboxing & Cognitive Memory Guide

`vkeygit` provides Antigravity V2 agents with **sub-millisecond virtual worktree sandboxes** and **content-addressed cognitive DAG memory**.

---

## Autonomous CLI Execution Rules (CRITICAL FOR ALL AGENTS)

> [!IMPORTANT]
> **NEVER ask the user to type `/vkeygit` or `/vkeygit-sandbox-create` in the chat UI.**
> Those slash commands are user-facing triggers; **agents must always invoke the `vkeygit.exe` binary directly** via `run_command` using the full path `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkeygit.exe`.
>
> **3 Mandatory Invocation Rules on Windows (`pwsh`):**
> 1. **Do NOT quote the executable path:** Never run `"C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkeygit.exe" sandbox ...` (in `pwsh`, a quoted string at the start of a line without `&` causes `ParserError: Unexpected token 'sandbox'`). Always pass the unquoted path `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkeygit.exe`.
> 2. **Always pass `-dir <workspace_repo_root>`:** `c:\aCogSpaceSeed` is a multi-repo root; the `.git` directory resides inside individual workspaces (e.g. `c:\aCogSpaceSeed\00flow\unilangunichip\nautilus`). Always pass `-dir <workspace_repo_root>`.
> 3. **Always pass `-scope` relative to `-dir`:** Scope paths are matched relative to `-dir` (e.g. `-scope 81000-active-source`, NOT `-scope c:\aCogSpaceSeed\...`).

---

## Agent Operational Workflows

### 1. Ephemeral Sandbox Protocol (Safe Refactoring & Execution)

When the user asks for multi-file refactoring, compiler changes, or speculative testing:

1. **Create Sandbox:**
   ```bash
   C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkeygit.exe sandbox create -dir c:\aCogSpaceSeed\00flow\unilangunichip\nautilus -name laso-ws1 -scope 81000-active-source
   ```
2. **Execute Autonomous Work Inside Sandbox:**
   * Make edits directly inside the sandbox path returned (`C:\aCogSpaceSeed\00flow\hydrationcache\c0990-ephemeral-scratch\sandboxes\<name>`).
   * If tests fail, iterate in the sandbox without polluting the user's primary workspace.
3. **Verify Status & Partitioning:**
   ```bash
   C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkeygit.exe sandbox status -dir C:\aCogSpaceSeed\00flow\hydrationcache\c0990-ephemeral-scratch\sandboxes\laso-ws1
   ```
4. **Merge Back to Workspace:**
   ```bash
   C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkeygit.exe sandbox merge -dir C:\aCogSpaceSeed\00flow\hydrationcache\c0990-ephemeral-scratch\sandboxes\laso-ws1 --direct-merge
   ```
5. **Discard Sandbox:**
   ```bash
   C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkeygit.exe sandbox discard -dir C:\aCogSpaceSeed\00flow\hydrationcache\c0990-ephemeral-scratch\sandboxes\laso-ws1
   ```

---

### 2. Cognitive Memory Protocol (Branching Reasoning)

To record key milestones or branch alternative hypotheses:

1. **Save Reasoning Checkpoint:**
   ```bash
   C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkeygit.exe memory save -dir c:\aCogSpaceSeed\00flow\unilangunichip\nautilus -key "hypothesis-simd-v1" -prompt "AVX-512 tensor optimization" -tags "simd,speedup" -content "Plan and test results..."
   ```
2. **Recall Checkpoint:**
   ```bash
   C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkeygit.exe memory recall -dir c:\aCogSpaceSeed\00flow\unilangunichip\nautilus -key "hypothesis-simd-v1"
   ```
3. **List Active Cognitive DAG:**
   ```bash
   C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkeygit.exe memory list -dir c:\aCogSpaceSeed\00flow\unilangunichip\nautilus
   ```
