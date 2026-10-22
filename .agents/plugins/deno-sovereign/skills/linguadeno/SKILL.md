---
name: linguadeno
displayName: linguadeno
description: Sovereign LinguaDeno & Node.js SACP/qAPC Bridge Engine. Deploys sub-microsecond actor procedure call bridges to legacy Express/Fastify services, executes instant 0.75ms sandboxed TypeScript tools, provisions autonomous BYOE server fleets, and tracks 60/90-day pilot conversions with 100% VPU Acceleration Credit rollover.
---

# Sovereign `linguadeno` Skill for Antigravity Agents

## Overview

The **`linguadeno`** skill equips Antigravity (IDE Agent), Jules (Engineering Agent), and Conductor (Orchestration Agent) with the **Sovereign LinguaDeno Runtime Engine & Node.js SACP/qAPC Bridge**.

It operationalizes the **Stage 4 Node Modernization Bridge** and the **Zero-Free 60/90-Day Convertible Commercial Engine**:
- Transparently elevates legacy Express, Fastify, and NestJS services to sub-microsecond SACP QUIC actor calls on Port 5434.
- Executes sandboxed TypeScript tools and scripts in **`< 1 ms`** with zero heap garbage collection pauses (`0 B/op`).
- Dispatches 1-click autonomous server deployments to customer-owned bare-metal or cloud VMs (BYOE) via SSH or SACP.
- Settles commercial contracts through the sovereign fleet's **`plomb` Multi-Rail Banking Engine** (Mercury / Relay / Stripe), crediting 100% of 60/90-day pilot investments toward annual VPU Acceleration Credits.

```mermaid
flowchart TD
    subgraph LinguaDenoSkill ["SOVEREIGN LINGUADENO & SACP BRIDGE ARCHITECTURE"]
        direction TB
        BRIDGE["Stage 4: @velocikey/node-sacp Bridge\n- Dual-plane Express listener: HTTP (Port 3000) + SACP QUIC (Port 5434)\n- 15x - 30x Latency Drop (< 180µs) with ZERO code rewrites"]
        RUNTIME["Stage 1 & 3: Dynomination (Denomination - Pure Native Go)\n- 2.8 MB static binary, 0.75 ms boot, 0 B/op Plinth stack arenas\n- Drop-in CLI replacement for deno run, test, compile, bench"]
        BYOE["Stage 2: Autonomous BYOE Server Deployer\n- 1-click deploy to customer AWS/GCP/bare-metal servers\n- Zero cloud hosting COGS -> 99.9% Net Profit Margin"]
        SETTLEMENT["plomb Multi-Rail Banking Engine\n- 60/90-day paid pilot conversion with 100% fee rollover as VPU credits\n- Mercury virtual accounts, Relay sub-accounts, and Stripe ledger attestation"]

        BRIDGE --> RUNTIME --> BYOE --> SETTLEMENT
    end
```

---

## 1. Core Agent Capabilities

### 1.1 Injecting the Node.js SACP/qAPC Bridge
Agents can modernize any enterprise Node.js microservice by injecting the `@velocikey/node-sacp` wrapper:
```typescript
import express from "express";
import { enableSovereignSACP } from "@velocikey/node-sacp";

const app = express();
// Existing route handlers remain 100% unchanged...
enableSovereignSACP(app, {
  sacpPort: 5434,
  mode: "ZERO_COPY_qAPC",
  serviceName: "enterprise-api",
  affirmationDomain: "customer.corp"
});
```

### 1.2 Running Sub-Microsecond Standalone Bridge
Run the sovereign bridge executable directly from forge:
```bash
C:\aCogSpaceSeed\00flow\forge\96000-internal-executables\nodesacpbridge.exe -node-port 3000 -sacp-port 5434 -service express-orders -domain customer.corp
```

### 1.3 Autonomous BYOE Server Self-Deployment
Deploy to remote customer infrastructure without Docker or external package managers:
```bash
vkey deploy linguadeno --target=api.internal.customer.corp --profile=edge-cluster
```

### 1.4 Commercial Settlement & VPU Credit Tracking
Verify seat affirmation and pilot conversion via `plomb`:
```bash
C:\aCogSpaceSeed\00flow\forge\96000-internal-executables\purchaseinfo-validation.exe
```

---

## 2. Invariants & Governance

1. **Strict Zero-Free Policy:** Free tiers are prohibited. All evaluations require a paid 60-day or 90-day pilot license.
2. **100% VPU Credit Rollover:** Upon signing an annual license, 100% of pilot fees are credited as Virtual Processing Unit (VPU) Acceleration Credits.
3. **Pure Native Go & No CGo:** All runtime and bridge binaries must be 100% CGo-free Go compiling in `< 2 seconds`.
4. **Sub-Microsecond Latency:** Procedure calls over SACP QUIC must verify $< 180\,\mu\text{s}$ average latency.
