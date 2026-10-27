# Sovereign Agent Sandboxing & Cognitive Memory Governance Rule

> [!IMPORTANT]
> **RULE: SPECULATIVE MODIFICATIONS MUST USE EPHEMERAL SANDBOXES (com.velocikey.vkeygit)**
> 1. **Zero Active Workspace Blast Radius:** When executing speculative architectural refactors, experimental compiler optimizations, or multi-file mutations, agents MUST spin up an ephemeral sandbox via `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkeygit.exe sandbox create -dir <workspace_repo_root> -name <sandbox_name> -scope <relative_scope_paths>` in `c0990-ephemeral-scratch/sandboxes/`.
> 2. **Autonomous CLI Invocation (Never Delegate to User Chat):**
>    - **DO NOT** ask the user to type `/vkeygit on` or `/vkeygit-sandbox-create` in the chat UI. Agents MUST execute `vkeygit.exe` directly via `run_command`.
>    - **NO QUOTES AROUND EXECUTABLE PATH:** In `run_command`, never wrap the `vkeygit.exe` path in quotes (e.g. `"C:\...\vkeygit.exe"` fails with `ParserError: Unexpected token 'sandbox'`). Always invoke `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkeygit.exe ...` directly without quotes.
>    - **ALWAYS PASS `-dir <workspace_repo_root>` AND RELATIVE `-scope`:** Because the workspace root `c:\aCogSpaceSeed` is a multi-repo silo root (not a single `.git` repo), always pass `-dir c:\aCogSpaceSeed\<silo>\<family>\<workspace>` (where `.git` lives) and pass `-scope` relative to `-dir` (e.g., `-scope 81000-active-source`).
> 3. **Verification Before Merge:** Agents MUST run test suites inside the sandbox before executing `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkeygit.exe sandbox merge -dir <sandbox_path> --direct-merge` back onto the primary workspace.
> 4. **Clean Teardown:** After successful merge-back or upon abandoning a failed attempt, agents MUST call `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkeygit.exe sandbox discard -dir <sandbox_path>` to prune the scratch directory.
> 5. **Cognitive Memory Preservation:** Key milestone plans, architectural decisions, and benchmark scorecards MUST be sealed into the cognitive DAG via `C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkeygit.exe memory save -dir <workspace_repo_root> -key <key> -tags <tags> -content <text>`.
