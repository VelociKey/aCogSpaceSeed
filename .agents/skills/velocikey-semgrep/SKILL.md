---
name: velocikey-semgrep
description: Zero-copy SAST monorepo security auditor and pre-commit hook security guard derived from semgrep-seerci. Offers sub-millisecond AST pattern matching and Git pre-commit enforcement.
---

# VelociKey Semgrep Extension (`semgrep-seerci`)

The `velocikey-semgrep` extension integrates zero-copy static security analysis and pre-commit hook guards directly into Antigravity IDE and AGY CLI.

## Injected Slash Commands

- `/seerci-scan`: Fast zero-copy static security scan over active codebase.
- `/seerci-audit`: Deep monorepo SAST audit and rule attestation.
- `/seerci-guard`: Install and verify sub-1ms Git pre-commit security guards.

## Injected Capabilities

- Sub-millisecond pre-commit hook guard (`git_hook_guard.go`).
- Monorepo AST pattern matching and rule evaluation engine.
- Quantum SAST vulnerability scanning.
