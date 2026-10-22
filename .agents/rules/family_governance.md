# Sovereign Rule: Family Governance & 12-Step Promotion Lifecycle

> **Scope:** Monorepo-wide architectural governance for all 470+ workspaces.  
> **Hierarchical Stratum:** Silo $\longrightarrow$ Family $\longrightarrow$ Workspace.

---

<hard_invariants>
## 1. Core Structural Invariants
1. **NO ORPHAN WORKSPACES:** Every workspace MUST belong to exactly one family within its silo. Orphan workspaces directly under a silo root are deprecated.
2. **PREFIX-FREE TAXONOMY:** Workspace and family names must be strictly prefix-free domain semantic names (e.g. `nautilus`, `ingress`, `commerce`). No artificial single-letter prefixes (such as `q` or `s`).
3. **SILO LAYERING DIRECTIONALITY:** Dependencies MUST flow downward:
   $$\text{00xper} \longrightarrow \text{00floo / 00flon} \longrightarrow \text{00flow} \longrightarrow \text{00taut}$$
   - `00flow` workspaces MUST NOT depend on `00floo` or `00flon`.
   - `00floo` workspaces MUST NOT depend on `00flon` directly (use `logiclibrary`).
   - Production workspaces MUST NOT depend on `00xper`.
</hard_invariants>

---

<attestation_stratigraphy>
## 2. Attestation Stratigraphy Invariant
1. **Family-Level Attestation (Standard):** In any silo partitioned into families, attestation MUST reside strictly at the **Family level** (`<silo>/<family>/attestation/` with module `vkey/<family>-attestation`). It verifies structural and cryptographic integrity of sibling workspaces using M31 leaf hashes and White3 polynomial certificates.
2. **Silo-Level Attestation (Fallback for Family-Free Silos):** In flat reference silos with no families, attestation resides at `<silo>/attestation/`.
3. **Information Stratigraphy:**
   - **Silo Info (`<silo>-info`):** Catalogs what families are present in the silo.
   - **Family Info (`<family>-info`):** Authoritative domain contracts, interaction models, and member roles.
</attestation_stratigraphy>

---

<family_creation_checklist>
## 3. Canonical 12-Step Family Creation Checklist
When creating or modernizing a family, agents MUST execute all 12 steps in order:
```
[ ] 1.  Create family directory: <silo>/<family>/
[ ] 2.  Create .workspace-dna.webnf with FamilyDNA declaration
[ ] 3.  Create attestation sub-workspace: <family>/attestation/
[ ] 4.  Move/rename member workspaces into family directory
[ ] 5.  Update go.mod module paths (module vkey/<workspace>)
[ ] 6.  Update go.work use and replace directives
[ ] 7.  Update fleet_topology.webnf (with syntax validation)
[ ] 8.  Update fleet_composition.webnf if suite references exist
[ ] 9.  Create/update <silo>-info or <family>-info workspace
[ ] 10. Purge orphan binaries from workspace roots
[ ] 11. Verify via realize.exe -workspace <affected>
[ ] 12. Run nontrivial security SAST scan
```
</family_creation_checklist>
