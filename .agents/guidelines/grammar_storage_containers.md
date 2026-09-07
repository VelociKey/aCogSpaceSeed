# Guideline: LatentLingua Grammar & Executable Storage Taxonomy

## Core Intent
To maintain precise semantic and lexical distinctions between internal core platforms/protocols and external licensee customizations across the SDLC.

---

## 1. Grammar Storage Containers (LatentLingua)
LatentLingua generated grammars are versioned and stored using two distinct containers based on authority and flexibility:

### A. Die (Core/Product Grammars)
* **Concept:** A rigid, precision-engineered metal stamp that enforces high-performance, standardized shapes.
* **Scope:** Internal / product-specific.
* **Usage:** Houses core platform grammars and structural schemas. These are immutable boundaries used to stamp out our foundational runtime interfaces and code models.

### B. Mold (Licensee Grammars)
* **Concept:** A hollow cavity designed to receive custom materials and cast them into shape.
* **Scope:** External / licensee-specific.
* **Usage:** Houses licensee-facing grammar structures. Licensees "pour" their custom domain vocabularies, business rules, and extensions into these molds, ensuring they conform perfectly to the boundaries of the platform.

---

## 2. Executable & Binary Environments
The environments where compiler-validated binaries are delivered are split similarly to separate platform code from licensee distributions:

### A. Forge (Internal Binaries - sForge)
* **Concept:** The factory where core steel is shaped and tools are forged.
* **Scope:** Internal.
* **Usage:** Serves as the central repository for all internal/external binaries, tooling distributions, and core runtime engines.

### B. Foundry (Licensee Binaries)
* **Concept:** The factory where metal castings are poured and assembled using molds.
* **Scope:** Licensee/Customer.
* **Usage:** Serves as the delivery and runtime environment for licensee-specific compiled applications and service adapters.
