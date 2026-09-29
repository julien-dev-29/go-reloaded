# Go Reloaded — Text Editing / Auto-Correction Tool

## Overview

A Go CLI tool that reads a text file, applies a fixed set of textual transformations, and writes the result to an output file. Written from scratch for the Zone01 "go-reloaded" project. It will be corrected by auditors, so behavior must match the spec examples exactly and include unit tests.

## Usage

```
go run . <input.txt> <output.txt>
```

- Exactly two argument names required. Wrong arg count returns an error.
- Missing/unreadable input file returns an error.
- Output file is created (or overwritten) with the transformed text.

## Architecture

Single `main` package with pure, separately-testable functions:

| File | Responsibility |
|------|----------------|
| `main.go` | CLI entry: validate args, read input, call `ProcessText`, write output |
| `transform.go` | `ProcessText(text string) string` orchestrator + case/hex/bin/article rules |
| `punctuation.go` | Punctuation attachment and quote-pair formatting |
| `conversions.go` | Hex-to-decimal and bin-to-decimal base conversion utilities |
| `*_test.go` files | Table-driven unit tests |

No external dependencies. Standard library only.

## Data Flow

1. Read entire input file as a string.
2. `ProcessText(text)` applies transformation stages in order.
3. Write result string to output file.

## Processing Order

Inside `ProcessText`:

1. **hex/bin conversion** — replace `(hex)` / `(bin)` markers together with their preceding word.
2. **case transformations** — `(up)` / `(low)` / `(cap)` and `(up, N)` / `(low, N)` / `(cap, N)`.
3. **punctuation formatting** — attach `., ,, !, ?, :, ;` (including groups like `...` / `!?`) to the preceding word; format `'...'` quote pairs.
4. **article fix** — `a`/`A` → `an`/`An` before a word starting with a vowel or `h`.

## Transformation Rules

### Hex conversion

- `(hex)` converts the immediately preceding word (a hexadecimal number) to its decimal form.
- Input may be upper- or lowercase hex digits; output is a plain decimal string.
- Example: `1E (hex) files were added` → `30 files were added`
- If the preceding word is not a valid hex number, the text is left unchanged.

### Binary conversion

- `(bin)` converts the immediately preceding word (a binary number) to its decimal form.
- Example: `It has been 10 (bin) years` → `It has been 2 years`
- If the preceding word is not a valid binary number, the text is left unchanged.

### Case transformations

- `(up)` — uppercase the preceding word.
- `(low)` — lowercase the preceding word.
- `(cap)` — capitalize the preceding word (first letter uppercase, rest preserved as-is).
- Bare forms apply to exactly the one preceding word.
- `(up, N)`, `(low, N)`, `(cap, N)` — apply to the N words preceding the marker.
  - N is a positive integer.
  - If N exceeds the number of words before the marker, apply to all available words.
  - Marker and its parenthesized argument are removed from output.
- Examples:
  - `Ready, set, go (up) !` → `Ready, set, GO!`
  - `I should stop SHOUTING (low)` → `I should stop shouting`
  - `Welcome to the Brooklyn bridge (cap)` → `Welcome to the Brooklyn Bridge`
  - `This is so exciting (up, 2)` → `This is SO EXCITING`

### Punctuation

- `., ,, !, ?, :, ;` attach to the preceding word with **no space before** and **one space after**.
- Consecutive punctuation groups (e.g. `...`, `!?`, `?!`, `!!`) are treated as a single unit and keep the same attachment/spacing rule.
- Examples:
  - `I was sitting over there ,and then BAMM !!` → `I was sitting over there, and then BAMM!!`
  - `I was thinking ... You were right` → `I was thinking... You were right`
- Punctuation should also attach across newlines in the input (text is treated as a single stream).

### Quote pairs

- An even set of `'` characters form pairs (left quote / right quote).
- Left quote attaches directly to the following content, right quote attaches directly to the preceding content — no spaces.
- Spaces immediately after a left quote and immediately before a right quote are removed.
- Works for single words and multi-word content.
- Examples:
  - `' awesome '` → `'awesome'`
  - `' I am the most well-known homosexual in the world '` → `'I am the most well-known homosexual in the world'`
- Single (unpaired) quote: leave as-is (no special handling).

### Article fix

- A standalone word `a` becomes `an` when the next word begins with a vowel (`a, e, i, o, u`) or `h`.
- Case is preserved: `A` → `An`.
- The check is on the first character of the next word.
- Example: `A amazing rock` → `An amazing rock`

## Error Handling

- Standard Go errors; never panic for expected failures.
- CLI: exactly 2 args required, missing file → clear error message.
- Invalid hex/bin operands: left unchanged rather than erroring.
- Unrecognized markers or malformed `(<cmd>, N)`: left unchanged.

## Testing

- `TestProcessText` — table-driven tests for every rule using the spec examples.
- `TestConvertHex` / `TestConvertBin` — valid and invalid inputs.
- `TestPunctuation` — punctuation + quote grouping edge cases.
- `TestMain` — CLI validation (arg counts, missing files) via a `run` helper function.
- Integration test — write `sample.txt`, run the program, assert `result.txt` matches the spec's expected output.
- Run with `go test ./...`.