---
name: velocikey-antigravity-sdk-go
description: Ultra-fast Go-native Antigravity SDK extension. Provides sub-microsecond WebNF DSL program record conversions, zero-copy -mem=manual arena buffer pooling, and qAPC (QUIC Actor Procedure Call) messaging.
---

# VelociKey Antigravity SDK for Go Extension

The `velocikey-antigravity-sdk-go` extension exposes Go-native SDK orchestration directly within Antigravity IDE and AGY CLI.

## Injected Slash Commands

- `/agysdk-chat`: Stream sub-microsecond WebNF responses and thought deltas locally.
- `/agysdk-status`: Audit qAPC channel metrics and `-mem=manual` pool memory allocations.

## Injected Capabilities

- WebNF DSL Program Record Conversion (`ToWebNF()`, `ParseWebNF()`).
- High-Throughput qAPC Async Procedure Calling (`DispatchqAPC`).
- Zero-Allocation Arena Memory Pooling (`-mem=manual`).
