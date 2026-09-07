---
name: com.velocikey.formalprover
displayName: vkey-prover
description: >-
  Deterministic mathematical theorem prover and invariant verifier powered by the
  Lean 4 microkernel and DeepSeek-R1 Olympiad reasoning. Use to formally certify
  hash collision-resistance, algebraic ring commutativity, and finite-field Galois arithmetic.
user_invocable: true
arguments:
  - name: theorem
    description: Theorem statement or invariant to prove (e.g. 'White3 hash collision invariance')
    required: true
  - name: auto
    description: Attempt automatic proof closure via omega/linarith/ring/simp decision procedures
    required: false
    default: "true"
---

# vkey℠ Formal Prover: Zero-Hallucination Mathematical Verification

`vkey-prover` links Antigravity with the **Lean 4 microkernel** and **DeepSeek-R1** to produce formally certified mathematical proofs with **0% hallucination**.

---

## Injected Slash Commands

- `/vkey-prover --theorem "<statement>"`: Formally decompose and prove a mathematical theorem or invariant.
- `/vkey-prover-auto "<name>" "<signature>"`: Test if a theorem closes automatically with decision procedures (`omega`, `linarith`, `ring`, `simp`).
- `/vkey-prover-status`: Inspect Lean 4 microkernel availability and version.

---

## Execution Protocol Performed by the Agent

When this skill is invoked, execute the following neuro-symbolic sequence:

1. **Auto-Decision Procedure Fast-Path:**
   * Run `vkey prove auto <theorem_name> "<signature>"`:
     ```bash
     C:\aCogSpaceSeed\00flow\forge\97000-internal-toolchains\vkey.exe prove auto <name> "<signature>"
     ```
   * If closed in $< 50\text{ ms}$, return immediate Q.E.D. attestation.

2. **DeepSeek-R1 Theorem Decomposition:**
   * If complex, prompt `deepseek-r1` to structure the proof into a lemma DAG with explicit Lean 4 tactic blocks (`intro`, `apply`, `rw`, `exact`).

3. **In-Memory Microkernel Verification:**
   * Pipe the proof into the sovereign Lean 4 driver (`vkey prove theorem ...`).
   * If Lean emits compiler errors or unsolved goals (`⊢`), feed the exact diagnostics back to `deepseek-r1` to repair the tactic line. Hard limit of 3 repair iterations.

4. **Cryptographic Notarization:**
   * Record the verified proof artifact and Blake3 hash in `000ALL/00000-knowledge-foundations/900-attestations/`.
