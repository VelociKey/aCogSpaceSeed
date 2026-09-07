---
name: change-agent
description: Execute a Jules/NATVS workflow with a specific WebNF markdown goal or quoted string for a workspace.
user_invocable: true
arguments:
  - name: goal
    description: Raw goal string or path to the WebNF markdown goal file
    required: true
  - name: workspace
    description: Path to the target workspace directory
    required: true
  - name: direct
    description: Execute directly in-process instead of delegating to the background daemon
    required: false
    type: boolean
    default: false
---

# Change Agent Skill

Use this skill to assign a new goal task to the Jules agent or queue/run it using the `natvs-engine` orchestrator.

## Usage Patterns

- `use change-agent on <goal> for <workspace>`: Queues the task and runs it via the background coordination daemon.
- `use change-agent on <goal> for <workspace> --direct`: Runs the task synchronously in-process, bypassing the background daemon queue.

## Parameters

- **goal:** The raw goal string (in quotes) or the absolute/relative path to the WebNF markdown file containing the goal to execute.
- **workspace:** The target workspace path to run the goal on (e.g. `00flow/s-natives`).
- **direct:** Bypasses queue-watcher daemon check and runs directly in-process.

## Command Mappings

### Enqueued Daemon Mode (Default)
Spawns the detached background queue watcher daemon (if not already running) and delegates the task to `tasks.queue.webnf`:
```powershell
& "C:\aCogSpaceSeed\00flow\s-natives\natvs-engine.exe" "jules-run" "{{goal}}" "{{workspace}}"
```

### Direct Synchronous Mode
Executes the task directly and synchronously in the current terminal process:
```powershell
& "C:\aCogSpaceSeed\00flow\s-natives\natvs-engine.exe" "jules-run" "{{goal}}" "{{workspace}}" --direct
```
