# Sovereign Ancillary Discoveries & Future Work Registry

> [!IMPORTANT]
> **Purpose**:
> This registry serves as the authoritative ledger for capturing ancillary technical discoveries, compiler optimizations, and architectural enhancements discovered during active development across the fleet. It allows immediate logging of valuable cross-cutting ideas without derailing the main trunk of active development.

---

## Active Ancillary Discoveries & Planned Improvements

| ID | Subsystem | Discovery & Enhancement | Target Workspaces | Specification Document | Status |
| :---: | :--- | :--- | :--- | :--- | :---: |
| **ANC-001** | **Nautilus Border 1** | **Ghost Tokens & Trivia Sidecar**: Dual-stream preprocessing to preserve comments, formal invariants, and zero-drift line/col diagnostic error pointers without polluting 128-bit Macro-IR opcodes. | `00flow/nautilus`<br>`00flow/languagegrammar`<br>`00flow/macroir` | [`future-nautilus-border1-ghost-tokens.md`](file:///c:/aCogSpaceSeed/conductor/tasks/future-nautilus-border1-ghost-tokens.md) | ⏳ Deferred |
| **ANC-002** | **Plinth SLS / Substrate** | **Docstring Entropy Header**: Storing 16-bit docstring richness score inside 64-byte `FlatSymbol.Reserved` to enable server-side confluence doc selection with zero wire payload leakage. | `00skim/skim/manifold`<br>`00flow/logiclibrary` | [`comment_handling_with_flatsymbol_wire_protocol.md`](file:///C:/Users/josep/.gemini/antigravity/brain/fbe0ea05-ca3a-48af-a0c3-14c644c8322f/comment_handling_with_flatsymbol_wire_protocol.md) | ⏳ Deferred |
| **ANC-003** | **lingua-transpiler** | **Macro-IR to WIR Direct Adapter**: Implementing `NormalizeMacroIRToWIR` to bypass textual intermediate WAG parsing and feed binary `Instruction128` streams straight into `SynthesizeGoFromWIR`. | `00flow/lingua-transpiler`<br>`00flow/macroir` | [`macro_ir_to_source_language_proposal_architecture.md`](file:///C:/Users/josep/.gemini/antigravity/brain/fbe0ea05-ca3a-48af-a0c3-14c644c8322f/macro_ir_to_source_language_proposal_architecture.md) | ⏳ Deferred |
| **ANC-004** | **producefiledifference** | **AST Zhang-Shasha Near-Clone Diffing**: Utilizing Tree Edit Distance similarity ($\ge 85\%$) to flag semantic drift and unpatched bug discrepancies across near-congruent files. | `00flow/logiclibrary`<br>`00skim/skim/manifold` | [`skim_coalesce_research_and_decision_framework.md`](file:///C:/Users/josep/.gemini/antigravity/brain/fbe0ea05-ca3a-48af-a0c3-14c644c8322f/skim_coalesce_research_and_decision_framework.md) | ⏳ Deferred |

---

## Protocol for Logging New Ancillary Discoveries

When an agent or engineer discovers an optimization, structural defect, or toolchain enhancement during active development:
1. Do **not** stop the main trunk to perform premature speculative refactoring.
2. Add a new row to this registry with a unique `ANC-nnn` ID.
3. Document the technical rationale and proposed architecture in a task specification under `conductor/tasks/`.
4. Resume active execution on the primary milestone.
