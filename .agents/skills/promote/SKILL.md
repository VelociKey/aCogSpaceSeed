---
name: promote
description: Promotes a qualified workspace from the experimental silo (00xper) to the production/stable silo (00flow), renaming and updating all module dependencies.
user_invocable: true
arguments:
  - name: workspace
    description: Current folder name of the workspace to promote (e.g. emulator-stripetax)
    required: true
  - name: target_silo
    description: Silo directory to promote to (defaults to 00flow)
    required: false
    default: "00flow"
  - name: source_silo
    description: Current directory silo of the workspace (defaults to 00xper)
    required: false
    default: "00xper"
  - name: new_name
    description: New folder name of the workspace (defaults to stripping the 'x-' experimental prefix if present)
    required: false
---

# Promote Skill: Workspace Promotion and Silo Migration

Use this skill to promote a qualified workspace from an experimental silo (typically `00xper`) to a production/stable silo (typically `00flow`), automatically renaming files, updating reference paths, updating Go module/work files, and building via the rehydrator.

## Usage Patterns

- `use promote on <workspace>`: Promotes the workspace from `00xper` to `00flow`, stripping any `x-` experimental prefix (e.g., `x-emulator-stripetax` becomes `emulator-stripetax`).
- `use promote on <workspace> to <target_silo>`: Promotes the workspace to a custom target silo.
- `use promote on <workspace> --new_name <new_name>`: Promotes and renames the workspace to a custom name.

## Migration Steps Performed by the Agent

When this skill is invoked, you should execute the following sequence:

1. **Calculate Target Names:**
   * Determine the target folder name. If `new_name` is not provided, strip any `x-` experimental prefix (e.g., `x-emulator-stripetax` becomes `emulator-stripetax`). If the workspace has no `x-` prefix, use the name as-is.
2. **Directory Relocation:**
   * Move the project folder from `C:\aCogSpaceSeed\<source_silo>\<workspace>` to `C:\aCogSpaceSeed\<target_silo>\<new_name>`.
3. **Go Module Updates:**
   * In `C:\aCogSpaceSeed\<target_silo>\<new_name>\go.mod`, update the module path definition to use the new module name (e.g., `module sov.fleet/emulator-stripetax` instead of `module sov.fleet/x-emulator-stripetax`).
4. **Rewrite Import References:**
   * Search for all Go source files (`*.go`) inside the newly moved workspace, and replace any import paths referencing the old module path (e.g., `sov.fleet/x-emulator-stripetax`) with the new module path (e.g., `sov.fleet/emulator-stripetax`).
   * Do the same across any downstream workspaces that depend on or import this workspace.
5. **Update Root Workspace (`go.work`):**
   * Edit `C:\aCogSpaceSeed\go.work` (or standard multi-workspace config files) to remove the old path `./<source_silo>/<workspace>` and add the new path `./<target_silo>/<new_name>`.
6. **Update Emulator Registry:**
   * If the workspace is an emulator, locate and update the registry mapping in [`registry.emulator.webnf`](file:///C:/aCogSpaceSeed/00flow/okf-store/00000-knowledge-foundations/registry.emulator.webnf) to point to the new path and workspace name.
7. **Rehydrate & Seal:**
   * Run the primary builder `C:\aCogSpaceSeed\00flow\hydration\int-rehydrator.exe -workspace <new_name> -local-only` from its new home to compile and verify all test targets, sealing the final binary into `forge`.
