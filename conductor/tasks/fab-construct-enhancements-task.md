# Active Task Goal: Sovereign Workspace Construction & Semantic Stage/Publish Engine

**Task ID:** `00flow-TASK-FAB-CONSTRUCT-ENHANCEMENT`  
**Sponsor:** Conductor / NATVS Engine  
**Status:** PROPOSED (Negotiation & Assimilation Phase)  
**Target Tracks:** **[fab_construct_enhancement_track.md](file:///c:/aCogSpaceSeed/conductor/tracks/fab_construct_enhancement_track.md)**

> [!NOTE]
> `conductor/` changes (including this file) are committed to the root `VelociKey/aCogSpaceSeed` repo.
> That repo is not covered by any `fab-construct preserve -silo` pattern.
> `fab-construct add` also hangs at the root (`nontrivial` scans the entire fleet).
> Until Phase 5 below is implemented, root-repo commits must be handled manually.

---

## 🎯 Global Goal Mandate

The **NATVS Engine** or **Gemini-CLI** must enhance the Go-based **`fab-construct`** utility inside the `s-fab-aides` workspace. The tool must support modular workspace/GitHub twin construction and implement high-level semantic Git operations (`add` and `publish`) backed by native security, vulnerability, and license compliance auditing (using **`nontrivial`**).

---

## 📋 Directives & Invariant Constraints

1.  **Twin Separability Invariant:** Users/Agents must be able to invoke `fab-construct workspace` to build local layouts only, `fab-construct github` to provision GitHub twins only, and `fab-construct both` for unified initialization.
2.  **License Compliance Invariant:** The semantic `add` command must automatically scan all project dependencies using **`nontrivial`** to assert that no viral/copyleft licenses (e.g., GPLv3) are introduced. Staging must halt immediately if non-compliant licenses are detected.
3.  **Security & Maturity Gates:** The semantic `publish` (push) command must execute vulnerability and secret scanning (via **`nontrivial`**) and assess software maturity metrics (documentation presence, test coverage) before pushing. It must physically block any remote push to the `zenith` branch if any safety assertion fails.

---

## 🚀 Active Roadmap to Execute

*   [ ] **Phase 1: Refactor CLI Subcommands**
    *   Enhance subcommands to handle isolation between local workspace directory hydration and GitHub API remote instantiation.
*   [ ] **Phase 2: Implement Semantic Staging (`add` subcommand)**
    *   Integrate `nontrivial` license compliance scan (library call where possible, CLI wrapper as fallback), scoped to the target workspace only.
    *   Add automated Software Bill of Materials (SBOM) generation and mapping.
    *   Verify commit message schema compliance before calling local Git commit.
*   [ ] **Phase 3: Implement Semantic Publishing (`publish` subcommand)**
    *   Integrate `nontrivial` vulnerability and secret scanner.
    *   Develop a Software Maturity Assessment routine (checking for test presence and `AGENTS.md` compliance).
    *   Execute `git push origin zenith` only upon $100\%$ gate clearance.

*   [ ] **Phase 4: Implement `sync` subcommand — Remote Rebase/Fast-Forward**
    *   **Capability Gap Discovered:** `2026-07-14` — `applyconcurrency` remote `zenith` had ahead commits after a concurrent external push. No sovereign toolchain path (`fab-construct`, `fab-inspect bash`) could perform a pull/fetch/rebase. Agents were blocked.
    *   **Requirement:** Add `fab-construct sync -workspace <path>` (or `-silo <silo>`) that:
        1. Fetches the remote `zenith` branch (`git fetch origin zenith`).
        2. Rebases local commits on top (`git rebase origin/zenith`), or fast-forwards if no local commits exist.
        3. Validates the result still compiles (`go build ./...`) before returning success.
        4. On rebase conflict: aborts, restores prior state, and reports the conflicting files — never leaves a workspace in a dirty mid-rebase state.
    *   **Invariant:** `preserve` should grow a `--sync-first` flag that automatically invokes `sync` before pushing, eliminating the "fetch first" rejection class entirely.
    *   **Priority:** HIGH — currently requires an exception to the raw-git prohibition each time a remote workspace diverges.

*   [ ] **Phase 5: Root Repo Coverage (`VelociKey/aCogSpaceSeed`)**
    *   **Capability Gap Discovered:** `2026-07-14` — `conductor/` lives in the root `aCogSpaceSeed` git repo (`VelociKey/aCogSpaceSeed`). No `fab-construct preserve -silo` pattern covers it (silo pattern requires `NNxxxx` or `000ALL`). `fab-construct add` at the root hangs because `nontrivial` traverses the entire fleet (~100 workspaces) when invoked from the go.work root.
    *   **Requirement A:** Extend `preserve` to accept `-silo root` (or `-workspace .`) that targets the root repo directly.
    *   **Requirement B:** `fab-construct add` must scope the `nontrivial` compliance scan to only the changed files/modules in the target workspace, not the entire go.work tree. The current fleet-wide scan makes `add` unusable at the root.
    *   **Priority:** HIGH — `conductor/` task files, tracks, and state cannot be pushed to remote without manual intervention.
