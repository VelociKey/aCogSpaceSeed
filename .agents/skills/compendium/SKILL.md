---
name: compendium
description: "Compile and synthesize workspace compendiums and silo summaries into 000all/cognition/c0411-compendium/"
user_invocable: true
---

# Compendium Builder Skill
This skill synthesizes concise, high-signal silo summaries under `000all/cognition/c0411-compendium/` detailing all workspaces and their architectural purposes for LLM chat ingestion.

## Command Execution
Run the sovereign `silo_summarizer.exe` tool to generate the individual silo status documents and the unified master summary:
```bash
# Generate concise summaries per silo
C:/aCogSpaceSeed/00flow/forge/96000-internal-executables/silo_summarizer.exe -silo 00flow
C:/aCogSpaceSeed/00flow/forge/96000-internal-executables/silo_summarizer.exe -silo 00floo
C:/aCogSpaceSeed/00flow/forge/96000-internal-executables/silo_summarizer.exe -silo 00flon
C:/aCogSpaceSeed/00flow/forge/96000-internal-executables/silo_summarizer.exe -silo 00xper

# Generate the unified master summary across all silos
C:/aCogSpaceSeed/00flow/forge/96000-internal-executables/silo_summarizer.exe -silo all
```
