---
name: com.velocikey.multimodal
displayName: vkey-multimodal
description: >-
  Vision-Language multimodal architecture reanimator powered by Qwen2.5-VL-72B.
  Deconstructs system architecture flowcharts, UI wireframes, and dense LaTeX/PDF tables
  directly into declarative Flutter Dart UI and pure Go SACP actors.
user_invocable: true
arguments:
  - name: image_path
    description: Absolute path to the architecture diagram or UI screenshot
    required: true
  - name: target_workspace
    description: Target workspace to write generated Dart UI or Go actors (e.g. 'web_release', 'fab-aides')
    required: true
---

# vkey℠ Multimodal Architecture Reanimator

`vkey-multimodal` bridges visual system diagrams, wireframes, and LaTeX tables directly into executable **declarative Flutter Dart** and **Go SACP actors** using **Qwen2.5-VL-72B**.

---

## Injected Slash Commands

- `/vkey-multimodal --image_path "<path>" --target_workspace <ws>`: Reanimate visual architecture diagram.
- `/vkey-multimodal-ocr "<path>"`: Extract LaTeX equations and dense tables with pixel precision.

---

## Agent Operational Protocol

1. **Visual Grounding & OCR:**
   * Pass the image to `qwen2.5-vl-72b` to extract node connectivity, data flow directions, and UI widget hierarchies.

2. **Declarative Synthesis:**
   * If UI wireframe $\to$ Emit declarative Dart widgets in `81000-active-source/ui/` (`UI = f(State)`).
   * If Architecture Diagram $\to$ Emit Go SACP transport actors in `81000-active-source/transport/`.

3. **Compilation & Invariant Verification:**
   * Run `realize.exe -workspace <target_workspace>` to certify zero compile errors.
