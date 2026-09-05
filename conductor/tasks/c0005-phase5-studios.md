# Phase 5: Symmetrical Studio & Portal Frontends

* **Status:** Planned
* **Target Promotion Silo:** `00floo`
* **Primary Workspaces:** [o-afflume-realization](file:///c:/aCogSpaceSeed/00floo/o-afflume-realization), [o-afflume](file:///c:/aCogSpaceSeed/00floo/o-afflume) (backend service layer)
* **SDK Reference:** `s-a2a` (CODEC / MV4 bridge), `s-logiclibrary` (Gatekeeper mTLS)

## Objectives
Build the Flutter Web UIs — **Value Vector Studio**, **Ledger Studio**, and **Ordering Portal** — as Material View 4 surfaces backed by Go-native state projections via the CODEC bridge. All frontends must comply with M3 brand tokens, i18n/l10n ARB bundles, and WCAG 2.2 AA accessibility.

**Architecture Model:** Go backends own all business logic and state. Dart/Flutter is a thin MV4 surface that renders Go-projected state. No business logic in Dart.

> [!IMPORTANT]
> Per **ADR-019**, all Flutter web targets must be compiled using standard JS compilation (`flutter build web`) only. Do **not** use `--wasm` or `--dart2wasm` flags; `dart:html` usage in the legacy library surface blocks Wasm-GC compilation.

## Checklist

### 5.1 Design System & Brand Tokens
- [ ] Define the M3 seed colors (3 seeds → 47 role-based color tokens) in `lib/theme/color_schemes.dart`
- [ ] Generate the full M3 typography scale using `Outfit` (display) and `Inter` (body) from Google Fonts
- [ ] Create a `AppTheme` class that applies the color scheme and typography scale as a `ThemeData` for both light and dark modes
- [ ] Define reusable M3 component tokens: `borderRadius`, `elevation`, `stateLayerOpacity`, `motionDuration`
- [ ] Create a `DesignTokens` constants file to avoid magic numbers throughout the UI

### 5.2 i18n / l10n / a11y Scaffolding
- [ ] Initialize `.arb` ARB bundle files for `en` (English), `fr` (French), `de` (German), `ja` (Japanese) as the initial locales
- [ ] Configure `flutter_localizations` in `pubspec.yaml` and generate `AppLocalizations` delegates
- [ ] Apply semantic labels (`Semantics` widget) to all interactive controls for WCAG 2.2 AA screen reader support
- [ ] Ensure all text contrast ratios meet WCAG 2.2 AA (≥ 4.5:1 for normal text, ≥ 3:1 for large text) — validate with `flutter_accessibility_scanner`
- [ ] Add focus traversal ordering (`FocusTraversalGroup`) to all modal dialogs and data entry forms

### 5.3 Go → Flutter CODEC Bridge
- [ ] Define `afflume_surface.wag` grammar for the Go-projected state types that Flutter will render (ValueVector summaries, invoice line items, session status)
- [ ] Generate Go structs from `afflume_surface.wag` via `s-latentlingua` codec generation
- [ ] Implement the Go-side `SurfaceProjector` that serializes current state into CODEC frames and broadcasts over the local QUIC stream
- [ ] Implement the Dart-side `CodcReceiver` that deserializes CODEC frames into `ChangeNotifier`-backed state objects for MV4 rendering

### 5.4 Value Vector Studio
The primary tool for the platform operator to view real-time and historical VPU consumption.

- [ ] **Dashboard Screen:** Per-workspace VPU consumption gauges (current period vs. included quota), trend sparkline (last 30 days), active sessions count
- [ ] **Offering Configuration Screen:** Read/edit `current_offering_rules.webnf` entries — tier name, included VPUs, overage price; requires `billing_admin` RBAC role
- [ ] **VIP Metric Weights Screen:** View/edit `ActiveVIPRules` mapping (metric name, feature SKU, VPU weight per unit)
- [ ] **Session Drill-Down:** Select a workspace and view its raw spool file entries grouped by metric type and day

### 5.5 Ledger Studio
The auditor-facing tool for inspecting and manually adjusting the double-entry ledger.

- [ ] **Chart of Accounts Browser:** Hierarchical tree view of `masterledgercoa.webnf` — families → accounts, with purpose annotations and account type badges (Asset, Liability, Equity, Revenue, Expense)
- [ ] **Transaction Journal Viewer:** Paginated list of journal events with filter by date range, event type, and workspace ID; each entry expandable to show full debit/credit posting breakdown
- [ ] **Ledger Chain Integrity Panel:** Triggers `VerifyChainIntegrity()` and displays a pass/fail badge with hash mismatch details if any broken links are found
- [ ] **Manual Adjustment Entry:** Form for authorized accountants (`auditor` RBAC role) to post manual correcting journal entries; writes to the append-only chain with `correction_of` reference field
- [ ] **Period Close Report:** Generates a printable trial balance for a selected date range, exportable as PDF

### 5.6 Ordering Portal
The customer-facing checkout wizard for initiating subscriptions.

- [ ] **Tier Selection Step:** Card grid displaying all 5 subscription tiers (Solo Founder → Enterprise) with included VPUs, price, and feature highlights; selecting a tier emits a `tier_selected` receptor event
- [ ] **Customer Details Step:** Form collecting company name, billing contact, billing address (country-aware jurisdiction field for tax calculation)
- [ ] **Invoice Preview Step:** Renders the `BillableInvoice` breakdown (tier fee + overage estimate + tax) before Stripe payment; uses CODEC-projected data from Phase 3 resolver
- [ ] **Payment Step:** Embeds Stripe.js hosted payment link (URL from `CheckoutSession.PaymentURL`); listens for payment confirmation webhook → displays success state
- [ ] **Confirmation Step:** Displays invoice number, ledger reference, and provides a download link for the signed PDF invoice

### 5.7 Navigation & Routing
- [ ] Use `go_router` for declarative URL-based routing with deep-link support
- [ ] Implement route guards: `ValueVectorStudio` and `LedgerStudio` require `billing_admin` or `auditor` role; `OrderingPortal` is public
- [ ] Add a unified left-rail `NavigationDrawer` with M3 `NavigationDestination` items for each studio/portal

### 5.8 Build & Deployment
- [ ] Configure `flutter build web` (JS-only, no WASM) in the `71000-build-harness` build target for `o-afflume-realization`
- [ ] Serve compiled Flutter web assets via the `o-afflume` backend's embedded file server (`pkg/web` package)
- [ ] Add a `--dev-mode` flag to `o-afflume` binary that uses `flutter run -d chrome` hot-reload instead of embedded assets

## Verification
- `flutter analyze` reports zero errors or warnings for all studio targets
- Accessibility scanner reports WCAG 2.2 AA compliance for all 5 screens in each studio
- All 4 locales render without missing ARB key warnings
- End-to-end routing works: deep-link to `/ledger/chain-integrity` loads the panel directly
- `flutter build web` produces a valid deployable bundle under `o-afflume-realization/96000-internal-executables/web/`
