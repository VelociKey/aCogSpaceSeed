# Rule: VCS Engine Switching & Routing Governance

> [!IMPORTANT]
> **VCS ENGINE SELECTION PROTOCOL (com.velocikey.vkeygit)**
> 1. **Conversational Directives:**
>    * When the user requests *"Use vkeygit for version control"*, *"Enable vkeygit"*, or `/vkeygit on`:
>      - Execute `vkeygit config set-engine vkeygit`.
>      - Route all future version control, status, branch, and commit operations to `vkeygit`.
>    * When the user requests *"Switch back to git"*, *"Use standard git"*, or `/vkeygit off`:
>      - Execute `vkeygit config set-engine git`.
>      - Revert version control operations to standard `git` CLI subprocesses.
> 2. **Persistent Engine State:**
>    - The active engine preference is persisted in `.git/vkeygit_state.webnf` at the repository root.
> 3. **Default Behavior:**
>    - If unspecified, `vkeygit` is the default sovereign engine, providing sub-millisecond execution, isolated scratch sandboxes, and White3-PQC cognitive DAG memory.
