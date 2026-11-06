---
name: tensor-swarm
description: Dispatches high-speed parallel subagent sweeps across sovereign open models (vkalpaca, deepseek-r1, qwen-coder, gemma-2, vkgem) using UMOG 16-worker dispatch and 16-slot token matrix.
user_invocable: true
arguments:
  - name: tasks
    description: List or description of parallel tasks/prompts to evaluate
    required: true
  - name: model
    description: "Target model: vkalpaca, deepseek-r1, qwen-coder, gemma-2, vkgem (default: vkalpaca)"
    required: false
    default: vkalpaca
---

# Tensor Swarm Skill: Sovereign Open Model Multi-Agent Engine

Use this skill when Antigravity needs ultra-fast, local, deterministic, and parallel subagent execution across open models complementing cloud models (Gemini & Claude).

## Architecture & Subagent Flywheel

The sovereign `tensor` family open model engine executes directly against native host hardware with:
1. **UMOG 16-Worker Parallel Subagent Dispatcher**: Asynchronous non-blocking subagent queues executing in pure Go.
2. **Dedicated 16-Slot Multi-Agent Token Pool**: Cacheline-isolated (natural 64-bit word struct alignment) pinned memory buffers delivering $0\text{ B/op}$ steady-state allocations.
3. **Plinth MTC4 In-Situ Evaluation**: Evaluates compressed 8-bit/4-bit Morton tiles directly into registers, bypassing DRAM dequantization bottlenecks.
4. **Float128 RMSNorm & Octonionic Conformal RoPE**: Supports context scaling up to 2,000,000 tokens with certified numerical precision (< 1e-4 relative error).

## 5 First-Class Tool Calls Available via MCP (`com.velocikey.tensor`)

- `tensor_model_infer`: Unified single-model inference, `vkgem` 256K multilingual translation (`target_language` / `source_text`), `qwen-coder` test synthesis (`test_type` / `source_code` / `file_path`), 5-mode `PersonaVault` & CAA steering (`lora_adapter: "persona:..."` / `"steer:..."`), and Kronecker `slerp:` / `ties:` merged models.
- `tensor_swarm_dispatch`: Concurrent execution of up to 16 subagents in parallel with Tropical $(\min,+)$ geodesic routing and Galois $\text{GF}(2)$ bitwise deduplication.
- `tensor_speculative_critic`: Unified speculative code critique + deep AST security scan (`slog`, CGo, raw `git`, loop allocations) + `TracIn` influence attribution + $<1\ \mu\text{s}$ $O(1)$ `KVCacheEpochManager` rollback on violation.
- `tensor_carpetbag_pack`: First-class worktree and context bundler (`CARPETBAG-VERKLE-V1` + White3-256 cryptoseal + VerkleReceiver roundtrip verification + direct mounting into the unified 16-EB `MountModel` aperture).
- `tensor_catalog_inspect`: Unified model catalog + `LithicModelDatabase` workstation runnability + 16-EB `UnifiedMountReceipt` facet status (`SOVEREIGN_SECRET_VAULT` vs. `OPEN_ZERO_COPY_APERTURE`).

## Model Specializations & Complementary Roles

- **`vkalpaca`** (405B Doubled LLaMA 3.1): Deep synthesis, sovereign historical knowledge bases, high-context domain attestation.
- **`deepseek-r1`** (MoE Reasoning): Chain-of-thought mathematical proofing, formal verification, and branch exploration.
- **`qwen-coder`** (AST Specialist): High-accuracy code generation, signature adaptation, and zero-allocation test synthesis.
- **`gemma-2`** (Dense Speculator): Compact 2B dense representations for instant AST auditing and pre-flight critique.
- **`vkgem`** (Doubled Gemma 4 Multilingual Specialist): **Language Translation & Cross-Lingual Localization**. Powered by a massive 256,000-token SentencePiece multilingual vocabulary and 1,000,000-token Octonionic $S^7$ RoPE context, eliminating byte-fallback fragmentation and preserving semantic fidelity across global languages (French, German, Spanish, Italian, Japanese, Chinese, Latin, and classical dialects).

## Command Line Invocation

Run the sovereign MCP daemon via Stdio pipe:
`C:\aCogSpaceSeed\00flow\forge\99000-internal-actors\vkeymcp.exe -stdio`

Or run multi-agent benchmarks via family attestation:
`C:\aCogSpaceSeed\00flow\forge\99000-internal-actors\tensor-attestation.exe`
