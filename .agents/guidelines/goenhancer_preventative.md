# Guideline: GoEnhancer Preventative Coding Standards (Shift-Left)

## Core Intent
To write Go code that complies with the platform's AST-verification standards from the start, preventing violations and ensuring that the `goenhancer` runner discovers zero issues.

---

## 1. Concurrency Alignment
* **No Mutex + Map Structs:** Never define structs that mix standard lock fields and map fields:
  ```go
  // ❌ Prohibited (flags ConcurrencyEnhancer)
  type Registry struct {
      mu sync.Mutex
      m  map[string]Value
  }
  ```
* **Use Adaptive Map:** Initialize concurrent maps using the `adaptivelock` package:
  ```go
  // ✅ Correct
  type Registry struct {
      m adaptivelock.Map[string, Value]
  }
  ```

---

## 2. Explicit Variable Discards
* **No Blank Discards:** Do not use the blank identifier `_ = someFunc()` to discard returned values.
* **Use Discard Package:** Import `sov.fleet/logiclibrary/81000-active-source/pkg/200-enhancers/discard` and call:
  ```go
  discard.Discard(someFunc()) // Explicit, auditable intent
  ```

---

## 3. Logging & Formatting
* **Structured Logging Only:** Never use `log.Print/Fatal` or `fmt.Println` for diagnostics.
* **Slog Mandate:** Always use Go's structured `log/slog` package. Include `"workspace"` and `"operation"` correlation attributes.
* **Format Optimization:**
  * Use `strconv.Itoa(n)` instead of `fmt.Sprintf("%d", n)`.
  * Use `hex.EncodeToString(bytes)` instead of `fmt.Sprintf("%x", bytes)`.

---

## 4. Error Handling & Chain Propagation
* **Wrapping Errors:** Always wrap underlying errors using `%w`:
  ```go
  // ❌ Prohibited (flags ErrArchitect)
  return fmt.Errorf("read failed: %v", err)

  // ✅ Correct
  return fmt.Errorf("read failed: %w", err)
  ```
* **Error Returns:** Every exported function that can fail must return `error` as its last return value (never return booleans for error status).

---

## 5. Library Safety & Process Isolation
* **No os.Exit() in Libraries:** Code in library packages must never call `os.Exit()`. Return error values. `os.Exit()` is strictly restricted to `func main()`.
* **No panics:** Library code must not panic. Return `(T, error)` tuples and propagate context.
* **Modern Types:** Always write `any` instead of empty `interface{}`.
