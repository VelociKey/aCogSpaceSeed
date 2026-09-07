---
name: compendium
description: "Compile and synthesize workspace compendiums and silo summaries into 000all/cognition/c0411-compendium/"
user_invocable: true
---

# Compendium Builder Skill
This skill compiles workspace files into self-contained text compendiums and generates silo summaries under `000all/cognition/c0411-compendium/`.

## Command Execution
Run the following commands to generate the compendiums and summaries:
```bash
# Remove old summaries if they exist to prevent duplication
rm -f C:/aCogSpaceSeed/000all/cognition/c0411-compendium/silo_all_summary.md

# Loop through all active Go workspaces in go.work and compile their individual compendiums
for relPath in $(grep -oP '(?<=\./)(?:00flow|00floo|00xper|00flon|99cbox|99cogt|51slam)/[^\s]+' C:/aCogSpaceSeed/go.work); do
  absPath="C:/aCogSpaceSeed/$relPath"
  if [ -d "$absPath" ]; then
    echo "Compiling compendium for: $absPath"
    C:/aCogSpaceSeed/00flow/forge/96000-internal-executables/create_compendium.exe -workspace "$absPath"
  fi
done

# Run the silo summarizer to generate silo status documents and master summaries
C:/aCogSpaceSeed/00flow/forge/96000-internal-executables/silo_summarizer.exe -silo 00floo
C:/aCogSpaceSeed/00flow/forge/96000-internal-executables/silo_summarizer.exe -silo 00flow
C:/aCogSpaceSeed/00flow/forge/96000-internal-executables/silo_summarizer.exe -silo 00xper
C:/aCogSpaceSeed/00flow/forge/96000-internal-executables/silo_summarizer.exe -silo all
```
