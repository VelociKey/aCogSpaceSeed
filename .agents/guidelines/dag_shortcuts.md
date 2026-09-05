# Guideline: DAG & Dependency Shortcuts for LLM Analysis

## Core Intent
To bypass expensive filesystem indexing and raw code searches when analyzing project structure, build orders, or lock propagation, LLM agents (Antigravity, Jules, Sonnet, Opus) should use pre-computed dependency graphs (DAGs).

---

## 1. Available DAG Repositories

### A. Authoritative System DAGs (qdag)
* **Location:** `000ALL/cognition/41000-topology-grammars/`
* **Use Case:** Global build sequencing, workspace dependency boundaries, and promotion order constraints. Use this to determine which workspace to rebuild or modify first when dealing with cascading interface changes.

### B. Per-Workspace Ephemeral DAGs
* **Location:** `<workspace_root>/c4100-topology-grammars/` (e.g. `c:\aCogSpaceSeed\00flow\parallizer\c4100-topology-grammars/`)
* **Use Case:** Locating package imports, internal module calls, and interface structures within a single workspace.

---

## 2. When to Use DAG Shortcuts

Use these shortcuts instead of traversing source directories when asked to perform:

| Task Type | Recommended DAG Shortcut | Why it is Faster |
|:---|:---|:---|
| **Impact Analysis** (e.g. "What breaks if I modify `sacp`?") | Global DAG in `41000-topology-grammars/` | Immediately displays all downstream workspaces importing the target. |
| **Workspace Dependency Checks** (e.g. "How should I structure the Bazel build order?") | Authoritative system DAG | Displays topological sort order without AST scanning. |
| **Lock Import Scopes** (e.g. "Which workspaces use `adaptivelock`?") | Lock-import search map in workspace DAG | Shows where `sync` or `adaptivelock` are imported. |
