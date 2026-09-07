---
name: nontrivial
description: >-
  Zero-copy SAST security auditor, prompt-injection defense firewall, secret leak detector,
  and automated AST remediation engine for Antigravity workspaces and Cloud Agent Swarms.
  Use when performing security audits, scanning files for vulnerabilities, installing Git
  pre-commit security hooks, or checking codebase compliance.
---

# Nontrivial Sovereign Security Auditor & Cloud Guard

`nontrivial` is a pure native Go security auditing engine offering sub-millisecond AST pattern matching, zero-copy secret detection, and pre-commit hook security guards.

---

## 1. Core Tool Commands

### A. Audit Entire Workspace
Scan an entire workspace for security vulnerabilities, AST rule violations, and secret leakage:
```bash
nontrivial audit
```
* With strict failure enforcement (CI/CD mode):
  ```bash
  nontrivial audit --strict --fail-on-findings
  ```

### B. Audit a Specific Directory or File
```bash
nontrivial audit --path <absolute-path-to-dir-or-file>
```

### C. Install Git Pre-Commit Security Hook
Installs the sub-millisecond pre-commit hook in `.git/hooks/pre-commit`:
```bash
nontrivial hook-install
```

### D. Run as an MCP Server (Model Context Protocol)
Launches the persistent stdio/SACP MCP daemon for AI agent pairing:
```bash
nontrivial mcp
```

---

## 2. When to Activate This Skill

1. **Before Git Commit / Preservation:** Run `nontrivial audit` before running `vkey construct preserve` or committing changes.
2. **When Inspecting External Code:** Run `nontrivial audit --path <target>` when importing foreign repositories into active workspaces.
3. **When Auditing Cloud Agent Transcripts:** Scan incoming payloads for prompt injection and secret exfiltration.
