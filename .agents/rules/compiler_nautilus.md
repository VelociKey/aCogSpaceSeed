# Sovereign Rule: Nautilus Three-Border Compiler & Optimization Principles

> **Scope:** Applicable to all compiler, toolchain, and low-level code generation tasks across the sovereign fleet.  
> **Governed Workspaces:** `00flow/nautilus`, `00flow/chipgrammar`, `00flow/macroir`, `00flow/adph`.

---

<hard_invariants>
## 1. The Nautilus Compiler Optimization Principle (Zero Hand-Rolled SIMD or Padding Hacks)
1. **NO MANUAL PLATFORM SIMD OR INTRINSICS:** In all generated Go code, agents are strictly forbidden from writing hand-rolled platform-specific SIMD assembly, unrolled intrinsics, or architecture-specific vector intrinsics.
2. **NO ARTIFICIAL STRUCT PADDING HACKS:** Hand-rolled `_ [64]byte` padding fields are strictly prohibited. They bloat struct sizes, induce cacheline fragmentation, and distort compiler layout analysis.
3. **LINEAR MEMORY INDEXING OVER PLINTH:** Code must be written with clean, natural, idiomatic Go data structures and linear memory indexing over raw Plinth memory slices (`as.Bytes()`).
4. **COMPILER-DRIVEN OPTIMIZATION:** The Nautilus Sovereign Three-Border Compiler (`nautilus.exe`) automatically detects the host workstation chip architecture (Zen 4 AVX-512, Intel AVX2, ARM NEON), optimizes struct layout and cacheline alignment, analyzes loop dependencies, fuses micro-architectural registers, and emits optimal vectorization.
</hard_invariants>

---

<border_architecture>
## 2. The Three-Border Pipeline Architecture & Epistemic Separation

```
Border 1 (Polyglot Frontend)    --> Macro-IR Serialization (WIR1)   [Language: KNOWN    | Chip: AGNOSTIC]
Border 2 (Nine-Field Optimizer) --> VM-IR Serialization (WVM1)      [Language: AGNOSTIC | Chip: AGNOSTIC]
Border 3 (Backend Scheduler)   --> Native Silicon Machine Code      [Language: AGNOSTIC | Chip: KNOWN]
```

1. **Border 1 — Language Known, Chip Agnostic (`border1/`):** Lowers 1,042+ source dialects into non-recursive, flat-block 128-bit `Instruction128` Macro-IR (`WIR1`). Strictly forbidden from importing chip target metadata (`grammar-chip`, `targets`, `chelm`).
2. **Border 2 — Language Agnostic, Chip Agnostic (`border2/`):** Applies 40 Archimedean Strength-Reduction Laws, Nine-Field Langlands E-Graph equality saturation (pointer-free `ENode`), Dead-WAW register elimination, and Sethi-Ullman / Ershov topological cut-width canonicalization over a zero-pointer `FlatDAG`. Strictly forbidden from importing `grammar-language`, `grammar-chip`, `targets`, or `chelm`.
3. **Border 3 — Language Agnostic, Chip Known (`border3/`):**
   - **Default Engine:** Deterministic Zero-Pointer `FlatDAG` Negative-Late PERT (`NL-PERT`) Instruction Scheduler with **CDC 6000/7000 `FTN` Topological Register-Availability Pressure Control ($R_{\max}$)**, using target `chipgrammar` pipeline latencies, cross-die penalties, and hardware pipe counts (`target.TotalPipes`) for sub-15ms builds.
   - **Deep Tuning Engine (`-opt=max`):** Cycle-accurate genetic autotuner `CHELM` seeded by the `NL-PERT` schedule, with optional SACP GPU VM offloading when arbitrage is positive.
4. **Disaggregated State Preservation:** Intermediate states between borders are serialized via pure Go binary codecs (`wire.SerializeMacroIR` = `WIR1`, `wire.SerializeVMIR` = `WVM1`), enabling border-by-border execution and verification.
</border_architecture>

---

<toolchain_decision_matrix>
## 3. Compiler Execution & Verification Commands

| Compiler Task | Sovereign Command |
| :--- | :--- |
| Build Workspace via Nautilus | `C:\aCogSpaceSeed\00flow\forge\96000-internal-executables\realize.exe -workspace <name>` |
| Staged Non-Stranding Build | `realize.exe -workspace <name> -stage-only` |
| Deep Genetic Tuning | `realize.exe -workspace <name> -opt=max` |
| Differential A/B Benchmark | `realize.exe -workspace <name> -ab` |
| Fixed-Point Self-Hosting | `nautilus2.exe -workspace 00flow/nautilus2` |
</toolchain_decision_matrix>
