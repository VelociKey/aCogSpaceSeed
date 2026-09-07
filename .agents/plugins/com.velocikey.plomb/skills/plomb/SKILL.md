---
name: plomb
description: Sovereign Plomb Germinate & Plomb Grow software delivery engine for Antigravity, featuring real-time host topology inspection, zero-secrets passwordless seat affirmations, Nautilus bare-metal compilation, and Mercury/Relay settlement verification.
---

# Plomb: Sovereign Antigravity Software Delivery Engine

`plomb` is the official Antigravity plugin and skill for **bare-metal, chip-tailored software delivery** without generic bloated container images or static license keys.

---

## 1. Core Architecture

1. **`plomb germinate` (Host Topology Inspection & Seed Handshake):**
   * Discovers the exact CPU architecture (x86_64 AVX-512/AVX2, ARM64 NEON, Apple Silicon M-series).
   * Discovers accelerator hardware (NVIDIA Hopper SM90, Blackwell SM100, Apple Metal MPS).
   * Connects to the Sovereign SACP Hub to authenticate the developer's corporate email passwordlessly.
   * Acquires a signed, ephemeral `PlombGerminateManifest`.

2. **`plomb grow` (Nautilus Local Bare-Metal Synthesis):**
   * Invokes the local **Nautilus Three-Border Compiler**.
   * Emits exact SIMD compilation flags (`-mavx512f -mavx512bw`, `-framework Metal`, `-gencode arch=compute_90`).
   * Binds zero-copy manual arena buffer pools (`-mem=manual -zero-copy`).
   * Emits a permanent **BLAKE3 Build Invariant Seal** on the grown binary.

3. **Zero-Secrets Passwordless Affirmation:**
   * Eliminates all static API keys and license strings.
   * Developers authenticate with their corporate email (`alex@synapse.ai`).
   * Company Security Officers affirm the seat and bind it to a verified physical rooftop location.

4. **Multi-Rail Banking Verification:**
   * Verifies settled funds from **Mercury** (virtual accounts) or **Relay** (designated sub-accounts) before activating accounts.
   * All financial and settlement events are cryptosealed in the **Permanent Append-Only Transaction Log**.

---

## 2. Available Plugin Tools

* `plomb_germinate`: Inspects host hardware and generates the signed seed manifest.
* `plomb_grow`: Synthesizes and hydrates the bare-metal binary locally via Nautilus.
* `plomb_affirm_seat`: Affirms an employee seat with physical rooftop tax sourcing.
* `plomb_verify_settlement`: Verifies wire/ACH payment settlement in the permanent ledger.

---

## 3. Go Entry Point & Execution

* Binary Entrypoint: `C:\aCogSpaceSeed\00flow\forge\96000-internal-executables\purchaseinfo-validation.exe`
* Workspace: `00xper/purchaseinfo-validation`
* Testing: `realize.exe -workspace 00xper/purchaseinfo-validation`
