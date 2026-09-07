---
name: builder
description: Build workspaces, silos, or all components using the internal rehydrator.
user_invocable: true
arguments:
  - name: target
    description: The workspace name, a comma-separated list of workspaces, a silo name (e.g. 00flow), or "all"
    required: true
  - name: force
    description: Force execution and overwrite existing files
    required: false
    type: boolean
    default: false
---

# Builder Skill: Workspace Rehydration & Compilation

Use this skill to orchestrate and execute Go package builds using the Bazel-emulated internal hydration tool (`int-rehydrator.exe`).

## Usage Patterns

- `use builder on <workspace>`: Builds a single workspace (e.g. `s-logiclibrary`).
- `use builder on <workspace1, workspace2, ...>`: Builds a specified list of workspaces sequentially.
- `use builder on <silo>`: Builds all workspaces under a specific silo (e.g. `00flow`).
- `use builder on all`: Builds all workspaces in the `go.work` configuration under `00flow`.
- `use builder on <target> --force`: Bypasses checks and forces build synthesis.
- `use builder on <target> --help`: Displays the help and syntax guide.

## Ordering & Optimization Constraints

> [!IMPORTANT]
> **Foundational Requirement:** `s-logiclibrary` must **always** be compiled and processed before any other workspace that imports it. If the target selection includes `s-logiclibrary` (or when running `all` / `00flow`), the builder must execute the build of `s-logiclibrary` first.

> [!TIP]
> **Composite BLAKE3 Incremental Build Skip Standard:**
> - To optimize build execution speed, `realize.exe` computes a composite 256-bit BLAKE3 hash:
>   $$\text{Hash}_{\text{Composite}} = \text{BLAKE3}\Big(\text{BLAKE3}(\text{Workspace Source } 81000\text{-active-source/}) \;\parallel\; \text{BLAKE3}(nautilus\text{.exe}) \;\parallel\; \text{BLAKE3}(\text{WAG Grammars})\Big)$$
> - **Zero Rebuild Skip Rule:** If no source files in the target workspace have changed AND `nautilus.exe` has not changed (matching the `.exe.blake3` seal in `forge/96000-internal-executables/`), `realize` MUST **SKIP** the build for that workspace with status `[SKIP: ZERO_CHANGES_DETECTED]`.
> - Bypassed only if `--force` is specified.

## Mappings

For each workspace, the builder maps the build target to its harness file:

- **Harness Path:** `C:\aCogSpaceSeed\00flow\<workspace>\71000-build-harness\workspace.harness`
- **Fallback/Generation:** If the `workspace.harness` file does not exist, generate it dynamically at that path with the following schema:
  ```webnf
  ; Declarative Workspace Harness
  workspace_harness {
      name : "<workspace>" ;
      targets {
          target {
              mode : DYNAMIC ;
              engine : "nautilus-compiler" ;
              src : "<source_path_relative_to_workspace_root>" ;
              out : "<binary_name>.exe" ;
              category : "<category>" ;
          }
      }
      promote {
          realm : PLATFORM ;
          origin : INTERNAL ;
          category : <category> ;
      }
  }
  ```

## Command Execution

Executes the build command:
```powershell
& "C:\aCogSpaceSeed\00flow\s-hydration\int-rehydrator.exe" -harness "C:\aCogSpaceSeed\00flow\<workspace>\71000-build-harness\workspace.harness"
```
If `--force` is specified, also pass the appropriate flag where supported by the underlying hydrator runtime.
