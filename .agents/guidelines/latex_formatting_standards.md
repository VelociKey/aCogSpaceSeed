# KaTeX & LaTeX Rendering Standards (Zero KaTeX Parse Errors)

To prevent LaTeX and KaTeX parse errors across all agent outputs, artifacts, and markdown documents:

---

## Hard Invariants for Markdown & LaTeX

1. **CURRENCY FORMATTING (NEVER USE MATH DELIMITERS FOR DOLLARS):**
   - **WRONG**: `$\$0.0380$`, `$\$10.64$`, `$0.0380 / hr`
     - *Why it fails*: KaTeX parser interprets the leading `$` as math mode and `\$` as an unexpected escape character (`Unexpected character: '\'`), or pairs random dollar signs across tables into broken math blocks.
   - **CORRECT**: Write currency in plain text:
     - `$0.0380 / hr` (outside of math delimiters)
     - `USD 0.0380 / hr`
     - `0.0380 USD / hr`

2. **MICROSECONDS & SCIENTIFIC UNITS:**
   - **WRONG**: `\text{ \mu s }` (*\mu is undefined inside \text{}*), `$\mu\text{s}$` with unescaped backslashes in text.
   - **CORRECT**:
     - Use native Unicode: `52.4 µs`, `310 µs`, `100 µs`
     - Or valid math mode: `52.4\ \mu\mathrm{s}`

3. **PLAIN TEXT IN TABLES:**
   - Use clean, standard plain text formatting in markdown tables for numbers, units, and rates (`6.29 ms`, `176,470x`, `$0.0380 / hr`).
   - Reserve LaTeX math blocks `\( ... \)` and `\[ ... \]` strictly for genuine mathematical equations and formulas.
