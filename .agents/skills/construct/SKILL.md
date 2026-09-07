---
name: construct
description: Hydrate a new named workspace under the specified directory silo and automatically provision its GitHub remote twin.
user_invocable: true
arguments:
  - name: name
    description: Name of the workspace (prefix-free semantic name, e.g. latentflow, firestore)
    required: true
  - name: silo
    description: Target directory silo name on disk (defaults to 00flow)
    required: false
    default: "00flow"
  - name: org
    description: Target GitHub organization (default VelociKey)
    required: false
    default: "VelociKey"
---

# Construct Skill: Workspace Taxonomy Hydrator & Provisioner

Use this skill to fabricate a new workspace directory using the standard workspace directory-tree taxonomy and automatically create the corresponding GitHub remote twin at `VelociKey/<workspace_name>` with the branch `zenith`.

## Usage Patterns

- `use construct on <name>`: Synthesizes both local taxonomy structure and GitHub repository under the default `00flow` silo.
- `use construct on <name> for <silo>`: Synthesizes local taxonomy structure under the specified silo directory and provisions the GitHub repository.
- `use construct on <name> --silo <silo>`: Synthesizes both local and remote components under a custom silo directory.

## Parameters

- **name:** Prefix-free semantic name of the repository/workspace (e.g. `latentflow`, `pubsub`, `firestore`).
- **silo:** Parent directory silo on disk (e.g. `00flow`, `00xper`, `000ALL`). Defaults to `00flow`.
- **org:** Target GitHub organization. Defaults to `VelociKey`.

## Command Execution

Executes both local workspace taxonomy hydration and remote Git/GitHub creation:
```bash
C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkey.exe construct both -name "{{name}}" -silo "{{silo}}" -org "{{org}}"
```
