# Rule: Sovereign Agent Sandboxing & Cognitive Memory Governance

> [!IMPORTANT]
> **RULE: SPECULATIVE MODIFICATIONS MUST USE EPHEMERAL SANDBOXES**
> 1. **Zero Active Workspace Blast Radius:** When executing speculative architectural refactors, experimental compiler optimizations, or unfamiliar dependency integrations, agents SHOULD spin up an ephemeral sandbox via `vkeygit sandbox create` in `c0990-ephemeral-scratch/sandboxes/`.
> 2. **Verification Before Merge:** Agents MUST run test suites inside the sandbox before executing `vkeygit sandbox merge` back onto the primary workspace.
> 3. **Clean Teardown:** After successful merge-back or upon abandoning a failed attempt, agents MUST call `vkeygit sandbox discard` to prune the scratch directory.
> 4. **Cognitive Memory Preservation:** Key milestone plans, architectural decisions, and benchmark scorecards SHOULD be sealed into the cognitive DAG via `vkeygit memory save`.
