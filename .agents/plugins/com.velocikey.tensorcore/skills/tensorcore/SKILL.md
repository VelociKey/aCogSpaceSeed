---
name: tensorcore
description: Sovereign TensorCore High-Performance Inference, Virtual MoE, Flat DAG Pipelines, and Multi-Document Codec Engine for Antigravity Agents.
---

# TensorCore Sovereign Inference & Model Orchestration Skill

Activate this skill when you need to:
1. **Enforce 100% Zero-Failure Schema Output**: Guarantee generated JSON matches exact TypeScript/JSON schemas using Double-Array Trie (DAT) token-level logit masking.
2. **Process Multi-Document Batches**: Extract facts from heterogeneous sets of PDFs (`%PDF-`), CSV tables, clinical XMLs (`<?xml`), and TXT files without Python/Poppler dependencies.
3. **Execute Virtual MoE & Concurrence of Experts (CoE)**: Deliberate across specialist models (Qwen-Medical, DeepSeek-R1, LLaMA-3.3, Mistral) using low-rank Kronecker adapters ($\Delta \mathbf{W} = \mathbf{A} \otimes \mathbf{B}$) and SIMD logit blending.
4. **Run Flat Rolled-Out Model Pipelines**: Execute non-recursive, staged model chains (`flat_pipeline`) with zero-copy in-memory pointer buses.
5. **Ingest Private Models & Continual Learning**: Ingest SafeTensors/GGUF checkpoints into `.model.webnf` and extract post-run adaptation deltas into cryptosealed `.delta.webnf` records.

---

## 🛠️ Key Capabilities & Available Tools

* **`tensorcore_transpile_schema`**: Transpiles raw TypeScript interfaces or JSON Schemas into formal `.model.webnf` contracts.
* **`tensorcore_batch_extract`**: Runs pure-Go parallel document parsing across heterogeneous batches.
* **`tensorcore_virtual_moe_deliberate`**: Reconciles multi-specialist finding consensus under `WEIGHTED_ENTROPY` or `LOGIT_BLENDING` strategies.
* **`tensorcore_execute_flat_dag`**: Executes flat rolled-out stage sequences without academic "Dragon Book" AST recursion.

---

## 📜 Declarative Flat Pipeline Syntax (`flat_pipeline`)

```webnf
flat_pipeline IntakeAndAuditFlow {
    stage 0 IngestStage {
        task extract_facts : model="gemma-4-31b", schema="ClinicalFactsSchema" ;
    }
    stage 1 ParallelAuditStage {
        task oncology_audit : model="qwen-2.5-72b", adapter="oncology_v2", input="extract_facts" ;
        task pharma_audit   : model="deepseek-r1", input="extract_facts" ;
    }
    stage 2 ConcurrenceStage {
        task arbitrate : strategy="WEIGHTED_ENTROPY", inputs=["oncology_audit", "pharma_audit"] ;
    }
    stage 3 PresentationStage {
        task flutter_ui : model="gemma-4-31b", action="dart_m3_generation", input="arbitrate" ;
    }
}
```

---

## 🛡️ Telemetry & Provenance Invariants
* All model runs emit `.telemetry.webnf` tracking `model_name`, `chip_architecture`, `platform: "plinth"`, and a cumulative **BLAKE3 cryptoseal**.
* Never emit raw cloud provider strings (e.g. `aws-p5`, `gcp-a2`).
