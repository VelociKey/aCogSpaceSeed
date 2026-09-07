# Sovereign Agent Sandboxing & Cognitive Memory Governance Rule

> [!IMPORTANT]
> **RULE: SPECULATIVE MODIFICATIONS MUST USE EPHEMERAL SANDBOXES (com.velocikey.vkeygit)**
> 1. **Zero Active Workspace Blast Radius:** When executing speculative architectural refactors, experimental compiler optimizations, or multi-file mutations, agents MUST spin up an ephemeral sandbox via `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkeygit.exe sandbox create` in `c0990-ephemeral-scratch/sandboxes/`.
> 2. **Verification Before Merge:** Agents MUST run test suites inside the sandbox before executing `vkeygit sandbox merge` back onto the primary workspace.
> 3. **Clean Teardown:** After successful merge-back or upon abandoning a failed attempt, agents MUST call `vkeygit sandbox discard` to prune the scratch directory.
> 4. **Cognitive Memory Preservation:** Key milestone plans, architectural decisions, and benchmark scorecards MUST be sealed into the cognitive DAG via `vkeygit memory save`.
