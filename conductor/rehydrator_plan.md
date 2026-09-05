# Implementation Plan: Universal Bazel Build Modes for int-rehydrator

* **Status:** ✅ Complete

This plan outlines the modifications required in the internal rehydrator (`int-rehydrator`) toolchain and the workspace Bazel harnesses to support selecting between **Combined In-Browser Mode** and **Split Client/Server Mode**.

---

## 1. Goal
Introduce a new flag `-mode` in the `int-rehydrator` CLI that directs the compiler to output either:
1. **`combined`**: A single deliverable package. For web targets, this produces a unified `main.dart.wasm` file containing the embedded Go simulator engine. For native targets, this compiles a native OS server binary containing embedded web assets.
2. **`split`**: Separate targets. Generates a decoupled static web bundle (to deploy on CDNs like GitHub Pages) and a standalone Go server binary.

---

## 2. Proposed Changes

### A. Compiler Options Update
We modify [options.go](file:///C:/aCogSpaceSeed/00flow/hydration/81000-active-source/options.go) to define the new `-mode` CLI parameter:

```go
// In options.go:
type BuildOptions struct {
    // ... existing options
    Mode string // "combined" or "split"
}

// In ParseFlags():
flag.StringVar(&o.Mode, "mode", "", "Deployment mode: 'combined' (unified single executable/deliverable) or 'split' (decoupled UI and backend server)")

flag.Parse()

// Establish defaults: when Wasm mode is selected, the default mode is "split"; otherwise it is "combined"
if o.Mode == "" {
    if o.Wasm {
        o.Mode = "split"
    } else {
        o.Mode = "combined"
    }
}
```

---

### B. Workspace Bazel Harness Rules Update
We update the Bazel workspace `BUILD.bazel` target specifications to declare conditional compilation logic. The harness will inspect the `REHYDRATOR_MODE` environment variable passed down from `int-rehydrator`:

```python
# In BUILD.bazel:
config_setting(
    name = "mode_combined",
    values = {"define": "rehydrator_mode=combined"},
)

config_setting(
    name = "mode_split",
    values = {"define": "rehydrator_mode=split"},
)
```

---

### C. Build Execution Interception
In [engine.go](file:///C:/aCogSpaceSeed/00flow/hydration/81000-active-source/engine.go), we inject the selected mode into the compilation environments:

```go
// In engine.go (Purify implementation):
modeEnv := fmt.Sprintf("REHYDRATOR_MODE=%s", target.Mode)
cmd.Env = append(cmd.Env, modeEnv)
```

---

## 4. Efficient Multi-Target Compilation Design
To execute builds efficiently without redundant compilation loops, the Bazel harness partitions compilation into three primary targets that share underlying intermediate cache layers.

```mermaid
graph TD
    NormalTarget[//:timewarp_normal <br> Normal Mode] -->|Produces| NormalExe[timewarp_gateway.exe <br> Native Server serves embedded UI]
    SplitTarget[//:timewarp_wasm_split <br> WASM Split Mode - default] -->|Produces| SplitUI[static_ui_release/ <br> Decoupled UI assets for Pages]
    CombinedTarget[//:timewarp_wasm_combined <br> WASM Combined Mode] -->|Produces| CombinedWasm[main.dart.wasm <br> Single in-browser executable]
```

### Build Matrix Execution Rules:
1. **Normal Mode (`-wasm=false`):** 
   Compiles the native Go server executable, compiles the Flutter UI + Go Hub client WASM, and embeds them directly. The compiled native executable is delivered to the workstation's `/96000-internal-executables/` directory.
2. **WASM Split Mode (`-wasm=true -mode=split`):** 
   Generates a decoupled static UI folder containing `main.dart.wasm` and `go_server.wasm`, optimized to be committed directly to GitHub Pages for remote loopback connections.
3. **WASM Combined Mode (`-wasm=true -mode=combined`):** 
   Links the Go simulator engine directly into the Flutter Wasm-GC module, delivering a single unified `.wasm` file for direct browser execution.

### Cache Optimization:
Because Bazel maintains hermetic object cache graphs, the compiler caches intermediate Wasm translation units for the Go Hub and Dart UI. Switching between compiling the `split` UI and the `combined` UI triggers high-speed cache-hit resolutions, completing rebuild operations under ~2 seconds.


---

## 3. Verification Plan
* Run `int-rehydrator.exe -workspace timewarp -mode combined` and verify the compilation produces a single executable in the target directory.
* Run `int-rehydrator.exe -workspace timewarp -mode split` and verify the compilation outputs a decoupled web bundle and server binary.

---

## 5. Deployment Dispatcher & Multi-Target Routing

To support automated delivery across different pipelines after compilation, we introduce a build-all flag and a deployment dispatcher command.

### A. Compiler Flag for All Targets (`-mode=all`)
Setting `-mode=all` instructs the Bazel harness to build all three target layers simultaneously:
```powershell
int-rehydrator.exe -workspace timewarp -mode=all
```
This executes targets `//:timewarp_normal`, `//:timewarp_wasm_split`, and `//:timewarp_wasm_combined` in a single command, generating all OS binaries and browser-local/split WASM bundles in the `96000-internal-executables/` directory.

### B. Deployment CLI Command (`-deploy`)
We introduce a `-deploy` flag to trigger artifact promotion:
```powershell
int-rehydrator.exe -workspace timewarp -deploy=<target>
```

#### Supported Targets:
1. **`gcp` (GCP Artifact Registry & Cloud Storage):**
   Reads `c0100-configuration-registry/110-workspace-configs/deployment.webnf` for credentials, uploads backend Wasm/native binaries to the regional GCS staging bucket, executes Blake3 cryptographic attestation checks, and promotes the raw artifacts to the GCP Artifact Registry.
2. **`pages` (GitHub Pages):**
   Pushes the static UI release assets (split or combined in-browser bundle) to the `gh-pages` branch or updates the root `/docs` folder for Git check-in via:
   ```powershell
   fab-construct.exe preserve -silo 00xper
   ```
3. **`local-wasm` (Local browser-local distribution):**
   Deploys the static Wasm targets locally within the workspace's designated active execution tier paths.
4. **`all` (Universal deployment):**
   Deploys all compiled components to their respective targets (GCP, Pages, and local workstation Forge directories) in a single run.

---

## 6. Workspace Build "Type" Target Routing

To prevent directory structure pollution and reuse the existing semantic directory-tree taxonomy, `int-rehydrator` resolves output file paths by utilizing the existing **workspace build "type"** mapping contract (which maps target compilation categories like `BINARY`, `ACTOR`, and `TOOLCHAIN` to their corresponding Lattice sequence subpaths like `96000-internal-executables`, `99000-internal-actors`, and `97000-internal-toolchains` respectively).

```
Normal Mode Output Path = [Forge Base]  + [Workspace Build "Type" Mapped Directory]
Wasm Mode Output Path   = [Active WS]  + [Workspace Build "Type" Mapped Directory]
```

### Path Resolution Logic:
1. **Retrieve Target "Type" Location:**
   The compiler queries the workspace build "type" mapping to retrieve the compliant directory segment (e.g., `96000-internal-executables`).
2. **Apply Base Target Workspace:**
   * **If `-wasm=false` (Normal compile):** The target base resolves to the global Forge workspace (`C:\aCogSpaceSeed\00flow\forge`).
   * **If `-wasm=true` (Wasm compile):** The target base resolves directly to the active workspace folder (`wsPath` / e.g., `C:\aCogSpaceSeed\00xper\timewarp`).
3. **Write and Store:**
   The Bazel compiler writes the compiled artifacts into the resulting paths, ensuring that all variants reside in identically mapped directories inside both environments.


