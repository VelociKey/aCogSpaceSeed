---
name: create-emulator
description: Automates the generation, scaffolding, and registry setup of a service emulator workspace in 00xper.
user_invocable: true
arguments:
  - name: name
    description: Name of the service to emulate (e.g. pubsub, stripetax)
    required: true
  - name: type
    description: Type of emulator (e.g. taxation, payment, cloud_infrastructure)
    required: true
  - name: spec
    description: Path to OpenAPI spec or proto file (optional)
    required: false
---

# Create Emulator Skill: Automated Workspace Scaffolding

Use this skill to scaffold a new service emulator workspace under `00xper`, configuring Go packages, transpilation slots, registry mappings, and concurrent input queue worker loops.

## Usage Patterns

- `use create-emulator on <name> with type <type>`: Scaffolds a new emulator workspace named `emulator-<name>` under the specified type.
- `use create-emulator on <name> with type <type> using spec <path>`: Scaffolds a new emulator matching endpoints from a spec file.

## Execution Steps

When this skill is invoked, execute the following:

1. **Invoke Scaffolder:**
   * Run the scaffolding program `C:\aCogSpaceSeed\00flow\forge\fab-emulator.go` with target arguments:
     ```bash
     C:/aCogSpaceSeed/00flow/forge/92000-external-toolchains/go/bin/go.exe run C:/aCogSpaceSeed/00flow/forge/fab-emulator.go -name <name> -type <type> -spec <spec>
     ```
2. **Register in go.work:**
   * Edit `C:\aCogSpaceSeed\go.work` and add `./00xper/emulator-<name>` into the `use` block.
3. **Verify and Compile:**
   * Run the rehydrator tool on the generated workspace:
     ```bash
     C:/aCogSpaceSeed/00flow/hydration/int-rehydrator.exe -workspace emulator-<name> -local-only
     ```
