# Sovereign Agent Guideline: 12 Grand Category Functors Quick-Reference

> [!IMPORTANT]
> **PURPOSE (THE CANDLES & MATCHES DRAWER):**
> Whenever an AI Agent needs to apply category theory, bridge discrete algorithms to continuous mathematics, or formulate Lean 4 / WebNF proof reductions, do not guess or search arbitrarily. Use this quick-reference.

## The 12 Category Functor Decision Matrix

| Symbol | Name | Primary Target Domains | Lean 4 Type / Structure |
| :--- | :--- | :--- | :--- |
| **𝓕_Yoneda** | Presheaf Embedding | Splay Trees, Online Lookahead, Metric Geodesics | `YonedaPresheafFunctor (C : Type*) [Category C]` |
| **𝓕_Tropical** | Tropical Valuation | Polytopes, Associahedra, Monge Monotonicity | `TropicalValuationFunctor (R Γ : Type*)` |
| **𝓕_Dirac** | Dirac Spectral Trace | Variance Constants, Selberg Poles, GKW Operators | `DiracSpectralTraceFunctor (H : Type*)` |
| **𝓕_SAT** | Ideal Refutation | Graph Chromaticity (Hadwiger-Nelson), 3-SAT | `SATRefutationFunctor (R : Type*) [CommRing R]` |
| **𝓕_Motivic** | Motivic Homotopy | Polynomial Factorization, Galois Cohomology | `MotivicHomotopyFunctor (Var Mot : Type*)` |
| **𝓕_Langlands** | Automorphic Transfer | Arithmetic Progressions (Zaremba), Lie Harmonics | `LanglandsAutomorphicFunctor (G A : Type*)` |
| **𝓕_HodgeRiemann** | Combinatorial Chow | Matroid Minors (Rota), Graph Minors (Hadwiger) | `CombinatorialHodgeRiemannFunctor (M C : Type*)` |
| **𝓕_Noncommutative**| Cyclic C*-Algebra | Vertex-Transitive Hamiltonicity, Cayley Diameter | `NoncommutativeSpectralFunctor (D C : Type*)` |
| **𝓕_Fukaya** | Lagrangian Floer | Rotation Distance, Sequential Game Holonomy | `FukayaLagrangianIntersectionFunctor (M : Type*)` |
| **𝓕_Tannakian** | Group Scheme Duality | Fast Matrix Multiplication ω=2 GIT Polytopes | `TannakianGaloisDualityFunctor (T G : Type*)` |
| **𝓕_Perfectoid** | Perfectoid Tilting | Collatz 3x+1 Solenoid, 2-Adic Haar Measures | `PerfectoidTiltingFunctor (K0 Kp : Type*)` |
| **𝓕_Koszul** | Koszul Syzygy | Multiplier OBDD Lower Bounds, Resolution Width | `KoszulHomologicalSyzygyFunctor (A A_dual : Type*)` |

## Canonical Sovereign File Locations

* **Central Monograph:** `000ALL/00000-knowledge-foundations/120-system-architecture/grand_category_functors_central_registry.md`
* **Formal Lean 4 Kernel:** `00xper/knuth-open-conjectures/81000-active-source/lean4/MathesisCategoryCalculus.lean`
* **Native Go Registry:** `00xper/knuth-open-conjectures/81000-active-source/knuth_open_conjectures.go` (`GetGrand12CategoryFunctors()`)
