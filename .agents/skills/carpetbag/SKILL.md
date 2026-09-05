---
name: carpetbag
description: Sovereign Universal Carpetbag & Worktree Probe Skill for Antigravity. Automatically packages workspace worktrees (via vkeygit), web research, and prompt directives into a cryptosealed .webnf Carpetbag, dispatches to tensorcore-umc for in-slab execution, and streams certified results.
user_invocable: true
arguments:
  - name: prompt
    description: The diagnostic, verification, or mathematical prompt to evaluate
    required: true
  - name: model
    description: Target model ID (default kimi-k3, gemma-4-mini, deepseek-r1)
    required: false
    default: "kimi-k3"
  - name: worktrees
    description: Target workspace paths or comma-separated list of workspaces
    required: false
    default: "."
---

# Carpetbag Skill: Sovereign Worktree & Multi-Workspace Probe Engine

Use this skill to evaluate complex architectural, mathematical, or concurrency probes against workspace worktrees using **TensorCore UMC** and **Post-Riemannian MCT weights**.

The skill automatically traverses the target worktree(s), packages the files and directives into a cryptosealed **`CARPETBAG-VERKLE-V1` `.webnf`** payload, streams it to `tensorcore-core`'s in-slab `VerkleReceiver`, and returns clean response text with the single-line badge trailer.

## Usage Patterns

- `probe <model> on <workspace> "<prompt>"`: Evaluates the prompt against the designated workspace worktree using the specified model.
- `probe <workspace> "<prompt>"`: Defaults to `kimi-k3` for sovereign reasoning.
- `probe multi-workspace "00flow/tensorcore-umc,00flow/tensorcore-core" "<prompt>"`: Bundles multiple workspaces into a single Verkle Carpetbag.

## Execution Command

```bash
C:\aCogSpaceSeed\00flow\forge\99000-internal-actors\tensorcore-umc.exe probe {{model}} -worktrees "{{worktrees}}" "{{prompt}}"
```
