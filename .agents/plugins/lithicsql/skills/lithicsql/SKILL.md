---
name: lithicsql
description: Antigravity plugin skill for deploying, managing, benchmarking, autonomous server self-deployment, and Antigravity subscription billing for LithicSQL.
---

# LithicSQL Plugin Management & Self-Deployment Skill

Use this skill when interacting with, self-deploying to remote servers, or managing metered billing for **LithicSQL (`sqllithic`)** via the Antigravity Plugin Ecosystem.

## Self-Deployment & Billing Architecture

- **Autonomous Server Self-Deployment:** 1-Click push of `lithicsql.exe` (`2.47 MB`) to remote server instances via SSH/QUIC SACP.
- **Antigravity Metered Billing:** Billing flows directly through the developer's/enterprise's existing Antigravity Subscription based on Virtual Processing Units (VPUs) and Volume Performance (VIP).

## Common Commands

### 1. Self-Deploy LithicSQL to Remote Server
```powershell
lithicsql-deploy --target=user@remote-server-host --port=5433 --billing=antigravity-subscription
```

### 2. Execute Benchmark Suite
```powershell
$env:GOWORK='off'; C:\aCogSpaceSeed\00flow\forge\92000-external-toolchains\go\bin\go.exe test ./00xper/lithicattestation/81000-active-source/... -v
```

### 3. Check Peak Memory Footprint
```powershell
$env:GOWORK='off'; C:\aCogSpaceSeed\00flow\forge\92000-external-toolchains\go\bin\go.exe test ./00xper/sqllithic/81000-active-source/sqllithic_mem_profile_test.go -v
```
