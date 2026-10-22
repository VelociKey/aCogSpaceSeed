# Sovereign Rule: Declarative Flutter & Go Hub Pipeline

> **Scope:** All user interfaces, dashboards, visual telemetry, and web deployments across the fleet.  
> **Source Locus:** `81000-active-source/ui/` (Dart & Flutter Material Design 3).

---

<hard_invariants>
## 1. Zero Raw HTML/CSS Invariant
1. **NO RAW HTML OR CSS EDITING:** Agents are strictly prohibited from manually editing `.html` or `.css` files directly in `web_release/`, `build/web/`, or deployment assets.
2. **DECLARATIVE DART EXCLUSIVITY:** All UI components, styling, typography, layouts, and responsive behaviors must be authored strictly as declarative Dart ($\text{UI} = f(\text{State})$) in `81000-active-source/ui/`.
3. **AUTOMATED COMPILATION:** Production HTML/CSS/WASM assets are compiled deterministically by the sovereign Flutter pipeline; manual edits are overwritten and forbidden.
</hard_invariants>

---

<gohub_architecture>
## 2. The Go Hub SACP Connection Protocol
1. **NO HTTP REST POLLING:** Client-server interaction between the Flutter UI and the sovereign Go backend must never use traditional HTTP REST polling.
2. **DART JS-INTEROP & SACP TRANSPORT:** Communication flows directly through the **GO Hub** via modern `dart:js_interop` / `package:web` over raw multiplexed `golang.org/x/net/quic` or SACP transport streams.
3. **SUB-MICROSECOND DISPATCH:** Bypasses browser networking overhead, enabling real-time telemetry streaming, interactive canvas updates, and zero-copy binary state synchronization.
</gohub_architecture>
