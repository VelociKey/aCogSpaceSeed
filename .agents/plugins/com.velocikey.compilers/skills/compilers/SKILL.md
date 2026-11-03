---
name: compilers
description: VelociKey Drop-In Replacement Compilers Suite (vkey-rustc, vkey-clang, vkey-llvm, vkey-go, vkey-pythonc, vkey-cuda, vkey-triton, vkey-xla, vkey-tvm, vkey-iree, vkey-tsc, vkey-kotlinc, vkey-swiftc, vkey-slang, vkey-halide, vkey-taichi, vkey-flang, vkey-php, vkey-dart, vkey-wasm, vkey-sycl, vkey-numba, vkey-futhark, vkey-protoc, vkey-jaxc, vkey-zetasql, vkey-r8, vkey-bazel, vkey-v8, vkey-flutterc, vkey-tint, vkey-binaryen, vkey-j2cl, vkey-xls, vkey-cel, vkey-flatc, vkey-carbon). Strict 1:1 drop-in replacements for permissive-licensed public compilers with zero feature extension, delivering measurable compilation speed and memory improvements across chosen build-host chip architectures to accelerate human and LLM-agent development velocity and time-to-market.
---

# VelociKey Drop-In Replacement Compilers Suite (`com.velocikey.compilers`)

`com.velocikey.compilers` provides **strict, standalone 1:1 drop-in replacement executables** for the compilers powering **>90% of software development** (excluding `.NET`).

---

## 1. Scope & Standalone Execution Rules (Strict Path 2 Blackbox Compilation Invariant)

1. **Path 2 Blackbox Compilation of Unmodified Upstream Source Code (`ZeroSourceModification` & `ZeroFeatureExtension`):**
   - Every `vkey-<compiler>` beachhead substitute (`vkey-rustc`, `vkey-clang`, `vkey-llvm`, `vkey-go`, `vkey-pythonc`, `vkey-cuda`, `vkey-triton`, `vkey-xla`, `vkey-tvm`, `vkey-iree`, `vkey-tsc`, `vkey-kotlinc`, `vkey-swiftc`, `vkey-slang`, `vkey-halide`, `vkey-taichi`, `vkey-flang`, `vkey-php`, `vkey-dart`, `vkey-wasm`, `vkey-sycl`, `vkey-numba`, `vkey-futhark`, `vkey-protoc`, `vkey-jaxc`, `vkey-zetasql`, `vkey-r8`, `vkey-bazel`, `vkey-v8`, `vkey-flutterc`, `vkey-tint`, `vkey-binaryen`, `vkey-j2cl`, `vkey-xls`, `vkey-cel`, `vkey-flatc`, `vkey-carbon`) is **strictly a blackbox Nautilus compilation of the unmodified permissive upstream compiler source code**.
   - **Zero Source Changes & 100% Identical Emitted Target Code:** We do not modify upstream compiler source code or extend upstream language/target semantics. Each `vkey-<compiler>` binary accepts the exact same CLI flags, compiles the exact same source code, and emits **100% identical target binaries/SASS/bytecode** for the exact same target chip architectures that the public upstream compiler supports—changing **only** the host compiler executable size (`-dwarf=strip_all`), host compilation speed, and host memory footprint (`-mem=adaptive`).
2. **Run on More Build-Host Chip Architectures:**
   - Licensees can run the `vkey-<compiler>` executables natively on their chosen build-host chip architecture (`x86_64 AVX2`, `x86_64 AVX-512`, `ARM64`, `RISC-V`, etc.), eliminating host-side build bottlenecks on developer workstations and CI runners.
3. **Human & LLM Agent Time-to-Market Multiplier:**
   - Compilation speed and memory improvements directly accelerate both **human developer inner loops** and **autonomous LLM coding agent compile-test-repair cycles**, compressing development time and achieving earlier time-to-market for licensees' products.
4. **Zero Internal IP Disclosure:**
   - All `vkey-<compiler>` executables are stripped (`-dwarf=strip_all`, `0` exposed symbols).

### 1.1 Two-Path Hybrid SaaS `PTX` $\to$ `.cubin` (`NL-PERT`) Architecture

For NVIDIA GPU compilation across `vkey-cuda` and the 10 PTX-emitting compiler substitutes (`vkey-triton`, `vkey-xla`, `vkey-tvm`, `vkey-iree`, `vkey-slang`, `vkey-halide`, `vkey-taichi`, `vkey-sycl`, `vkey-numba`, `vkey-futhark`), `com.velocikey.compilers` supports a strict **Two-Path Hybrid Architecture**:

- **Local Path 2 Blackbox Mode (Default):** Executes 100% locally on the build host with `0%` SASS difference against upstream `ptxas`/`nvcc`, delivering `14.8x–26.4x` host compilation speedup and `-64% to -78%` host peak RAM reduction while preserving strict 1:1 upstream binary equivalence.
- **Optional Hybrid SaaS `NL-PERT` Mode (Path (a) & Path (b)):**
  - **Hybrid Path (a) — Integrated Compiler Flag / Environment Gate:** Activated via `--vkey-saas-nlpert` or `VKEY_CUDA_SAAS=1` inside `vkey-cuda` and the 10 PTX-emitting compiler substitutes to transparently route intermediate `PTX` through the `NL-PERT` lowering pipeline without changing build scripts.
  - **Hybrid Path (b) — Direct `PTX` $\to$ `.cubin` SaaS Endpoint:** Accepts standalone `PTX` artifacts directly via the SaaS endpoint and returns sealed `.cubin` / fatbinary objects.
  - **Zero-Knowledge Split-Synthesis (`WIR-1-Alpha`):** Both Path (a) and Path (b) use Zero-Knowledge Split-Synthesis (`WIR-1-Alpha`) so proprietary kernel names, symbols, and host business logic never leave the customer airlock in cleartext, delivering **`1.30x–1.34x` GPU kernel execution speedup** with **`100%` IEEE-754 bitwise-identical math**.

### 1.2 Architectural Distinction: Phase 2 (`Nautilus`-Compiled Only Drop-Ins) vs. Full Cleanroom Sovereign Replacements (Phase 3 / Native Stack)

All 37 SKUs in `com.velocikey.compilers` are **strictly Phase 2 (`NAUTILUS_PHASE2_DROPIN` — Nautilus-Compiled Only 1:1 Drop-Ins)** compiled from unmodified upstream permissive source code with `ZeroFeatureExtension = true`. They are explicitly distinct from—and serve as the zero-touch Day-1 commercial funnel into—VelociKey's **Full Cleanroom Sovereign Replacements**:

| Ecosystem Domain | Phase 2 (`Nautilus`-Compiled Only 1:1 Drop-In SKU) | Full Cleanroom Sovereign Replacement (Already Built in Fleet) |
| :--- | :--- | :--- |
| **Flutter & Dart UI / AOT / Shaders** | `vkey-dart` (`vkey-flutter`, Wave 6) + `vkey-flutterc` (`gen_snapshot`, `impellerc`, Wave 8) — Unmodified `dart-lang/sdk` & `flutter/engine` compiled by `Nautilus` (`3.5x–4.4x` faster builds, zero code changes) | **`Kestrel + Talon + Avataron`** + **`Kestrel-Talon Prism` Companion** (`00xper/unify-unity/avataron-ugui`, `00ulib/substrate/wasm`) — Full replacement eliminating Flutter's 3-tree heap, Impeller C++ shaders, and Wasm GC with 16-byte `WIR1SheafNode` stalks over Plinth `FileBackedAddressSpace` (**`110x`** faster, **`0.0 ms`** shader jank) |
| **JavaScript / Wasm Engine & Browser** | `vkey-v8` (`v8`, `d8`, `mksnapshot`, `turboshaft`, Wave 8) + `vkey-binaryen` (`wasm-opt`, Wave 8) — Unmodified `v8/v8` & `WebAssembly/binaryen` compiled by `Nautilus` | **`Vigor8` + `Gantry` + `Project Tungsten`** (`00xper/tungsten`, `00ulib/substrate/wasm`) — Cleanroom sovereign JS/Wasm engine and Chromium/`depot_tools` replacement |
| **Monorepo Build Orchestration** | `vkey-bazel` (`bazel`, `blaze`, `bzl`, Wave 7) — Unmodified `bazelbuild/bazel` compiled by `Nautilus` (`3.6x–4.4x` faster Skyframe evaluation for external Bazel repos) | **`Hydration` (`realize.exe` + `workspace.harness`)** (`00flow/forge/96000-internal-executables/realize.exe`) — Sovereign zero-JVM monorepo build & attestation engine |
| **SQL Semantic Analyzer & Lakehouse** | `vkey-zetasql` (`zetasql`, Wave 7) — Unmodified `google/zetasql` compiled by `Nautilus` (`3.5x–4.2x` faster SQL AST resolution) | **`LithicSQL` + `Romaine`** (`00flow/data/lithicsql`) — Sovereign zero-copy SQL engine and BigQuery/Iceberg replacement |
| **Wire Serialization & RPC** | `vkey-protoc` (`protoc`, `upb`, Wave 7) + `vkey-flatc` (`flatc`, Wave 8) — Unmodified `protobuf` & `flatbuffers` compiled by `Nautilus` | **`00ulib/wire` + `SACP` / `qAPC` + `weBNF.sn`** — Sovereign zero-copy wire framing and sub-microsecond Actor Procedure Calls |
| **Go Compiler & Standard Library** | `vkey-go` (`compile`, `asm`, `link`, Wave 3) — Unmodified `golang/go` compiled by `Nautilus` | **`Tengen Sovereign` (`vkey-tengen`) + `00ulib` + `Nautilus` Border 1** — Sovereign cleanroom Go toolchain and zero-allocation standard library |

---

## 2. Eight-Wave Staged Rollout & Vault-Backed Demand Poll (37 Phase 2 SKUs)

- **Wave 1:** `vkey-rustc` (`vkey-cargo`, `vkey-cranelift` — `MIT OR Apache-2.0`)
- **Wave 2:** `vkey-clang` (`clang`, `clang++`, `clang-cl`, `lld`, `hipcc`), `vkey-llvm` (`opt`, `llc`, `llvm-ar`, `llvm-nm`, `llvm-objdump`, `mlir-opt`, `mlir-translate` — `Apache-2.0 WITH LLVM-exception`)
- **Wave 3:** `vkey-go` (`compile`, `asm`, `link` — `BSD-3-Clause`)
- **Wave 4:** `vkey-pythonc` (`PSF-2.0`), `vkey-cuda` (`vkey-nvcc`, `vkey-cicc`, `vkey-cudafe++`, `vkey-ptxas`, `vkey-fatbinary`, `vkey-bin2c`, `vkey-nvlink`, `vkey-nvrtc`, `vkey-cuobjdump`, `vkey-nvdisasm` — `Apache-2.0 WITH LLVM-exception AND BSD-3-Clause`), `vkey-triton` (`MIT`), `vkey-xla` (`openxla/xla` — `Apache-2.0`), `vkey-tvm` (`apache/tvm` — `Apache-2.0`), `vkey-iree` (`iree-org/iree` — `Apache-2.0 WITH LLVM-exception`)
- **Wave 5:** `vkey-tsc` (`vkey-bun` — `Apache-2.0 AND MIT`), `vkey-kotlinc` (`JetBrains/kotlin` — `Apache-2.0`), `vkey-swiftc` (`swiftlang/swift` — `Apache-2.0 WITH Runtime-Exception`), `vkey-slang` (`shader-slang/slang` — `MIT`), `vkey-halide` (`halide/Halide` — `MIT`), `vkey-taichi` (`taichi-dev/taichi` — `Apache-2.0`)
- **Wave 6:** `vkey-php` (`PHP-3.01`), `vkey-flang` (`llvm-project/flang` Fortran drop-in for `flang` and `flang-new` — `Apache-2.0 WITH LLVM-exception`), `vkey-dart` (`vkey-flutter` — `BSD-3-Clause`), `vkey-wasm` (`Apache-2.0 WITH LLVM-exception`), `vkey-sycl` (`intel/llvm` DPC++/AdaptiveCpp — `Apache-2.0 WITH LLVM-exception`), `vkey-numba` (`numba/numba` — `BSD-2-Clause`), `vkey-futhark` (`diku-dk/futhark` — `ISC`)
- **Wave 7 (Priority 5 Google-Originated Enterprise & AI Wedge Compilers):**
  1. `vkey-protoc` (`VelociKey™ Protox` — `protocolbuffers/protobuf + protocolbuffers/upb`, `BSD-3-Clause`)
  2. `vkey-jaxc` (`VelociKey™ Jaxion` — `jax-ml/jax + openxla/stablehlo` covering `jaxc`, `jaxjit`, `pallas`, `stablehlo-opt`, `Apache-2.0`)
  3. `vkey-zetasql` (`VelociKey™ Zetacore` — `google/zetasql`, `Apache-2.0`)
  4. `vkey-r8` (`VelociKey™ Dexon` — `r8/r8 + aosp/platform/art` covering `r8`, `d8`, `dex2oat`, `BSD-3-Clause AND Apache-2.0`)
  5. `vkey-bazel` (`VelociKey™ Skyforge` — `bazelbuild/bazel + bazelbuild/starlark`, `Apache-2.0`)
- **Wave 8 (Remaining 9 Google-Salaried Compiler & Platform Team Substitutes):**
  6. `vkey-v8` (`VelociKey™ Vektor` — `v8/v8` covering `v8`, `d8`, `mksnapshot`, `turboshaft`, `BSD-3-Clause`)
  7. `vkey-flutterc` (`VelociKey™ Impel` — `flutter/engine + dart-lang/sdk/runtime/vm` covering `gen_snapshot`, `impellerc`, `BSD-3-Clause`)
  8. `vkey-tint` (`VelociKey™ Prism` — `dawn/tint + google/shaderc + KhronosGroup/SPIRV-Tools` covering `tint`, `glslc`, `spirv-opt`, `BSD-3-Clause AND Apache-2.0`)
  9. `vkey-binaryen` (`VelociKey™ Wasmerize` — `WebAssembly/binaryen` covering `wasm-opt`, `wasm-metadce`, `wasm-as`, `wasm-dis`, `Apache-2.0`)
  10. `vkey-j2cl` (`VelociKey™ Closurite` — `google/j2cl + google/closure-compiler`, `Apache-2.0`)
  11. `vkey-xls` (`VelociKey™ Silicore` — `google/xls` covering `xlscc`, `opt_main`, `codegen_main`, `ir_converter_main`, `Apache-2.0`)
  12. `vkey-cel` (`VelociKey™ Celerity` — `google/cel-go + google/cel-cpp` covering `cel-go`, `cel-cpp`, `cel_compiler`, `Apache-2.0`)
  13. `vkey-flatc` (`VelociKey™ Flatline` — `google/flatbuffers + google-ai-edge/LiteRT` covering `flatc`, `tflite_convert`, `litert-compile`, `Apache-2.0`)
  14. `vkey-carbon` (`VelociKey™ Anthracite` — `carbon-language/carbon-lang`, `Apache-2.0 WITH LLVM-exception`)

Each wave enters **Stage 1 (Antigravity IDE + `agy` CLI Exclusive)** for 30 days before graduating to **Stage 2 (`agy ci`)** and **Stage 3 (Enterprise Build Farms)**. Use `vkey_compiler_poll_vote` to record customer votes for the next wave and build-host chip architecture.

---

## 3. Reversible Standalone SWDT Airlock Installation & Deinstallation (`weBNF.sn` / `.wag`)

All installation receipts (`c0100-configuration-registry/110-workspace-configs/airlock_receipt.webnf`), status records (`AirlockStatus`), and deinstallation rollback receipts (`CompilerUninstallReceipt`) conform strictly to the Sovereign Workspace Directory-Tree Taxonomy (SWDT) and the embedded Wirth Augmented Grammar [`dropin_compiler_airlock.wag`](file:///c:/aCogSpaceSeed/00xper/antigravityplugins/81000-active-source/dropin_compiler_airlock.wag) (`:com:velocikey:compilers:airlock`):

- **Install Standalone SWDT Compiler Airlock (`97000-internal-toolchains/` + `71000-build-harness/workspace.harness`):**
  ```
  C:\aCogSpaceSeed\00flow\forge\96000-internal-executables\antigravityplugins.exe -compiler-action=install -sku=vkey-dart -workspace=<workspace_dir> -host-chip=<host_chip_arch> -q -force -y
  ```
- **Inspect Standalone SWDT Airlock Status:**
  ```
  C:\aCogSpaceSeed\00flow\forge\96000-internal-executables\antigravityplugins.exe -compiler-action=status -sku=vkey-dart -workspace=<workspace_dir> -q -force -y
  ```
- **Deinstall & Verify Clean Zero-Residue Rollback (`residual_files_count = 0`):**
  ```
  C:\aCogSpaceSeed\00flow\forge\96000-internal-executables\antigravityplugins.exe -compiler-action=uninstall -sku=vkey-dart -workspace=<workspace_dir> -q -force -y
  ```
