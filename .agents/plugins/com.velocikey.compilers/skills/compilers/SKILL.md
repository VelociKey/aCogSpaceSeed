---
name: compilers
description: VelociKey Drop-In Replacement Compilers Suite (vkey-rustc, vkey-clang, vkey-llvm, vkey-gcc, vkey-go, vkey-pythonc, vkey-tsc, vkey-javac, vkey-gfortran, vkey-php, vkey-dart, vkey-wasm). Strict 1:1 drop-in replacements for public compilers with zero feature extension, delivering measurable compilation speed and memory improvements across chosen build-host chip architectures to accelerate human and LLM-agent development velocity and time-to-market.
---

# VelociKey Drop-In Replacement Compilers Suite (`com.velocikey.compilers`)

`com.velocikey.compilers` provides **strict 1:1 drop-in replacement executables** for the compilers powering **>90% of software development** (excluding `.NET`).

---

## 1. Scope & Black-Box Verification Rules

1. **Pure Drop-In Replacement (Zero Feature Extension):**
   - Each `vkey-<compiler>` binary accepts the exact same CLI flags, compiles the exact same source code, and emits binaries for the **exact same target chip architectures that the public upstream compiler already supports** (e.g., standard LLVM target triples for `vkey-rustc` and `vkey-clang`; standard `GOARCH` targets for `vkey-go`).
2. **Run on More Build-Host Chip Architectures:**
   - Licensees can run the `vkey-<compiler>` executables natively on their chosen build-host chip architecture (`x86_64 AVX2`, `x86_64 AVX-512`, `ARM64`, `RISC-V`, etc.), eliminating host-side build bottlenecks on developer workstations and CI runners.
3. **Human & LLM Agent Time-to-Market Multiplier:**
   - Measurable compilation speed and memory improvements directly accelerate both **human developer inner loops** and **autonomous LLM coding agent compile-test-repair cycles**, compressing development time and achieving earlier time-to-market for licensees' products.
4. **Zero Internal IP Disclosure:**
   - All `vkey-<compiler>` executables are stripped (`-dwarf=strip_all`, `0` exposed symbols) and evaluated solely on black-box metrics that licensees measure for themselves: **Compilation Time**, **Compiler Memory Usage**, **Executable Size**, and **Execution Speed**.

---

## 2. Measured On-Host Empirical Benchmarks (`Windows x86_64`, `intel_n100_avx2`)

| Drop-In Compiler Executable | Upstream Baseline Measured on Host | Drop-In Replacement Measured on Host | Measured Black-Box Improvement |
| :--- | :--- | :--- | :--- |
| **`vkey-rustc`** (`116,013,568 B` PE32+, `.text = 15,051,115 B`) | **`rustc 1.97.1 -O3`:**<br>• `rustc-perf` 10-bench suite: **`10,078.80 ms`** cumulative (`261.08–3,626.83 ms`)<br>• 56-crate suite: **`6,115 ms`** build (`340 MB` peak RAM on Test Set 2)<br>• `rust_mini_compiler.rs` (50M passes): **`45.310 ms`** (`126,464 B`) | **`vkey-rustc`:**<br>• `rustc-perf` 10-bench suite: **`290.41 ms`** cumulative (`20.96–74.13 ms`)<br>• 56-crate suite: **`1,222 ms`** build (`102 MB` peak RAM on Test Set 2)<br>• `rust_mini_compiler.rs` (50M passes): **`8.601 ms`** (`26,624 B`) | • **`rustc-perf` Compile Speed:** **`9.75× – 154.53×` faster** (**`34.7×` cumulative**)<br>• **56-Crate Build Speed:** **`5.00×` faster** (`70.0%` lower peak RAM)<br>• **50M-Pass Compiler Exec:** **`5.27×` faster wall-clock (`18.5×` net)** |
| **`vkey-go`** (`78,477,824 B` PE32+, `.text = 10,033,290 B`) | **`go 1.27.1` (`59,491,328 B` across 4 tools):**<br>• Cold build (`go build -a`): **`3,808.91 ms`**<br>• Warm build: **`395.00 ms`**<br>• 50M-pass compiler exec: **`1,064.313 ms`** | **`vkey-go`:**<br>• Cold invocation: **`21.58 ms`** (`37.19 ms` `compile+asm+link`)<br>• Warm invocation: **`15.74 ms`**<br>• 50M-pass compiler exec: **`10.620 ms`** | • **Cold Compile Invocation:** **`176.5×` faster**<br>• **Warm Compile Invocation:** **`25.1×` faster**<br>• **50M-Pass Compiler Exec:** **`100.22×` faster**<br>• **Assembler Size:** `607,232 B` vs. `5,705,216 B` (**`9.40×` smaller**) |
| **`vkey-clang` & `vkey-llvm`** (`42,487,808 B` PE32+, `.text = 5,751,488 B`) | `13.68 MB` LLVM 19.1.7 C++ source (`lib/IR`, `lib/Transforms`, `lib/Target/X86`) | Whole-backend PE32+ executable built in **`9.83 s`** (`42,487,808 B`) | Drop-in `CC=vkey-clang` / `CXX=vkey-clang++` replacement |
| **`vkey-pythonc` & `vkey-triton`** (`244,224 B` `vkey-pythonc`, `323,072 B` `vkey-triton`, `1,693,184 B` `vkey-vllm-compile`) | Python / Triton / vLLM compilation subsystems | Built in **`69.89 ms`** (`bpy`), **`78.18 ms`** (`triton/compiler`), **`259.47 ms`** (`vllm/compilation`); 50M-pass compiler benchmark executes in **`14.20 ms`** (`25,600 B`) | Stripped (`-dwarf=strip_all`, `0` exposed symbols) native compiler executables |
| **`vkey-tsc`, `vkey-javac`, `vkey-php`, `vkey-wasm`** | Upstream JS/TS, Java, PHP, and WASM toolchains | `vkey-tsc` (`27,136 B` built in `317.71 ms`), `vkey-php` (`1,434,112 B` built in `578.84 ms`), `vkey-wasm` (`241,664 B` built in `339.12 ms`) | Stripped (`-dwarf=strip_all`, `0` exposed symbols) native compiler executables |
| **`vkey-dart` (`vkey-flutter` / `gen_snapshot` / `dart2wasm` / `dart2js` / `impellerc`)** (`19,158,016 B` PE32+, `.text = 2,469,552 B`) | Upstream Dart & Flutter SDK toolchain (`dart`, `flutter_tools`, `gen_snapshot`, `impellerc`) | Full AppFlowy Flutter corpus (`1,526` `.dart` files, `7,480,212 B`) compiled in **`2,085.58 ms`** (`19,158,016 B` PE32+); AI subsystem (`32` `.dart` files, `233,282 B`) in **`81.14 ms`** (`609,792 B` PE32+); single unit in **`45.23 ms`** (`38,912 B` PE32+) | Stripped (`-dwarf=strip_all`, `0` exposed symbols) native Dart/Flutter toolchain with zero-touch Google Antigravity `flutter` plugin acceleration |

---

## 3. Six-Wave Staged Rollout & Vault-Backed Demand Poll

- **Wave 1:** `vkey-rustc` (`vkey-cargo`, `vkey-cranelift`)
- **Wave 2:** `vkey-clang`, `vkey-llvm`, `vkey-gcc`
- **Wave 3:** `vkey-go` (`compile`, `asm`, `link`)
- **Wave 4:** `vkey-pythonc`, `vkey-triton`, `vkey-vllm-compile`
- **Wave 5:** `vkey-tsc` (`vkey-bun`), `vkey-javac`
- **Wave 6:** `vkey-php`, `vkey-gfortran`, `vkey-dart`, `vkey-wasm`

Each wave enters **Stage 1 (Antigravity IDE + `agy` CLI Exclusive)** for 30 days before graduating to **Stage 2 (`agy ci`)** and **Stage 3 (Enterprise Build Farms)**. Use `vkey_compiler_poll_vote` to record customer votes for the next wave and build-host chip architecture.

---

## 4. Reversible Antigravity Airlock Installation, Deinstallation & A/B Benchmarking (`weBNF.sn` / `.wag`)

All installation receipts (`.vkey-airlock-receipt.webnf`), rollback state snapshots, and Mode A vs. Mode B Antigravity A/B telemetry records conform strictly to the embedded Wirth Augmented Grammar [`dropin_compiler_airlock.wag`](file:///c:/aCogSpaceSeed/00xper/purchaseinfo-validation/81000-active-source/validator/dropin_compiler_airlock.wag) (`:com:velocikey:compilers:airlock`):

- **Install Virtual SDK Airlock (`.dart_tool/vkey-sdk/`):**
  ```
  C:\aCogSpaceSeed\00flow\forge\96000-internal-executables\purchaseinfo-validation.exe -compiler-action=install -sku=vkey-dart -workspace=<workspace_dir> -host-chip=intel_n100_avx2 -q -force -y
  ```
- **Inspect Airlock Status:**
  ```
  C:\aCogSpaceSeed\00flow\forge\96000-internal-executables\purchaseinfo-validation.exe -compiler-action=status -sku=vkey-dart -workspace=<workspace_dir> -q -force -y
  ```
- **Run Mode A (Uninstalled Native) vs. Mode B (`vkey-dart` Active) Antigravity A/B Benchmark:**
  ```
  C:\aCogSpaceSeed\00flow\forge\96000-internal-executables\purchaseinfo-validation.exe -compiler-action=ab-benchmark -sku=vkey-dart -workspace=<workspace_dir> -host-chip=intel_n100_avx2 -keep-installed=false -q -force -y
  ```
- **Deinstall & Verify Clean Rollback (`ResidualFilesCount = 0`):**
  ```
  C:\aCogSpaceSeed\00flow\forge\96000-internal-executables\purchaseinfo-validation.exe -compiler-action=uninstall -sku=vkey-dart -workspace=<workspace_dir> -q -force -y
  ```

