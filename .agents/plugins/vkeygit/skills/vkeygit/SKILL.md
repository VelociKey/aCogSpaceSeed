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

## Injected Slash Commands

- `/vkeygit-sandbox-create <name> [scope]`: Spin up an isolated worktree in scratch memory (`<20ms`).
- `/vkeygit-sandbox-merge <path>`: Decompose and apply verified changes back to active workspace via 3-way AST merge.
- `/vkeygit-sandbox-discard <path>`: Discard ephemeral sandbox without leaving any workspace footprint.
- `/vkeygit-memory-save <key> <content>`: Store plan milestones or reasoning traces in `.git/agent_memory/` with White3-PQC cryptoseals.
- `/vkeygit-memory-recall <key|tag>`: Recall historical thought DAGs without burning context window tokens.

---

## Agent Operational Workflows

### 1. Ephemeral Sandbox Protocol (Safe Refactoring & Execution)

When the user asks for multi-file refactoring, compiler changes, or speculative testing:

1. **Create Sandbox:**
   ```bash
   vkeygit sandbox create -name "task-opt" -scope "src/calculus,src/transport"
   ```
2. **Execute Autonomous Work Inside Sandbox:**
   * Make edits directly inside the sandbox path returned.
   * Run compiler and test harnesses inside the sandbox.
   * If tests fail, iterate in the sandbox without polluting the user's primary workspace.
3. **Verify Status & Partitioning:**
   ```bash
   vkeygit sandbox status -dir "<sandbox_path>"
   ```
4. **Merge Back to Workspace:**
   ```bash
   vkeygit sandbox merge -dir "<sandbox_path>" --direct-merge
   ```
5. **Discard Sandbox:**
   ```bash
   vkeygit sandbox discard -dir "<sandbox_path>"
   ```

---

### 2. Cognitive Memory Protocol (Branching Reasoning)

To record key milestones or branch alternative hypotheses:

1. **Save Reasoning Checkpoint:**
   ```bash
   vkeygit memory save -key "hypothesis-simd-v1" -prompt "AVX-512 tensor optimization" -tags "simd,speedup" -content "Plan and test results..."
   ```
2. **Recall Checkpoint:**
   ```bash
   vkeygit memory recall -key "hypothesis-simd-v1"
   ```
3. **List Active Cognitive DAG:**
   ```bash
   vkeygit memory list
   ```
