---
name: com.velocikey.translate
displayName: vkey-translate
description: >-
  Universal and Indigenous translation engine powered by Google DeepMind's Gemma-4-31B
  and the 256k SentencePiece tokenizer. Translates Flutter localization .arb files,
  technical documentation, and UI strings into 100+ global languages and 20+ Indigenous
  American languages (Navajo, Cherokee, Lakota, Quechua, etc.) with zero placeholder corruption.
user_invocable: true
arguments:
  - name: workspace
    description: Target workspace name or path containing localization files (e.g. 'web_release', 'cyberneticstudio')
    required: true
  - name: languages
    description: Comma-separated list of ISO-639 language codes or names (e.g. 'nv, chr, es, ja, de')
    required: true
---

# vkey℠ Universal & Indigenous Translator

`vkey-translate` provides high-fidelity, culturally nuanced translation for software UI, documentation, and Flutter applications using **Gemma-4-31B**'s 256k Universal SentencePiece vocabulary.

---

## Injected Slash Commands

- `/vkey-translate --workspace <name> --languages "nv, chr, es, ja"`: Translate all `.arb` bundles.
- `/vkey-translate-string --text "<text>" --to "<lang>"`: Translate a standalone text segment.

---

## Agent Operational Protocol

1. **Extract Localization Keys:**
   * Scan the target workspace for template localization files (e.g. `lib/l10n/app_en.arb` or `81000-active-source/ui/l10n/`).
   * Preserve all `{placeholder}` and `$variable` tokens verbatim.

2. **Gemma-4-31B Polyglot Synthesis:**
   * Dispatch translation requests to `gemma-4-31b` with strict grammatical and cultural preservation rails.

3. **Deterministic Dart Verification:**
   * Write generated `app_<lang>.arb` files.
   * Run `dart run intl_translation:generate_from_arb` or `flutter gen-l10n` to confirm 100% syntactic validity.
