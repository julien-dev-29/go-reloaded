# Go Reloaded Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** A Go CLI that reads a text file, applies spec-defined transformations (hex/bin conversion, case changes, punctuation/quote formatting, article fix), and writes the result to an output file.

**Architecture:** Single `main` package. `ProcessText` orchestrates pure transformation stages over a whitespace-split token slice: normalize markers → convert hex/bin → apply case → fix articles → separate punctuation → assemble with quote-aware and punctuation-aware spacing. `main.go` validates args, reads the input, calls `ProcessText`, writes output.

**Tech Stack:** Go 1.26 (installed), standard library only (`os`, `strconv`, `strings`, `regexp`, `fmt`, `path/filepath` for tests, `testing`). `go test ./...` for verification.

## Global Constraints

- Standard Go packages only. No external dependencies.
- Never `panic` for expected failures — return `error` from `run`; print to `os.Stderr` in `main`.
- CLI requires exactly 2 args (`<input.txt> <output.txt>`).
- The 6 spec usage examples MUST pass byte-for-byte (single-space sentence, per spec output).
- Marker format inside parentheses is normalized to `(cmd)` or `(cmd,N)` before tokenizing.
- Invalid operands (non-hex/binary, unknown markers, malformed `(cmd,N)`) are left unchanged, never error.
- Table-driven tests for every rule (DRY), run with `go test ./...`.
- `gofmt`, `go vet`, `go test ./...` and a manual `go run .` sample must all pass at the end.
- Frequent commits; verify `git status` clean end-stage.

---

### Task 0: Scaffold — go.mod, git init, commit spec

**Files:**
- Create: `go.mod`

**Interfaces:**
- Consumes: nothing.
- Produces: a `go-reloaded` module so `go run .`, `go test ./...`, `go vet ./...` work from the project root.

- [ ] **Step 1: Initialize the Go module**

```bash
go mod init go-reloaded
```

Expected: creates `go.mod` containing `module go-reloaded` and `go 1.26.0`.

- [ ] **Step 2: Init git and commit the spec**

```bash
git init
git add docs/superpowers/specs/2026-09-15-go-reloaded-design.md go.mod
git commit -m "chore: add design spec and go module scaffold"
```

Run each command separately; if `git commit` fails with a missing identity error, run:

```bash
git -c user.name="student" -c user.email="student@local" commit -m "chore: add design spec and go module scaffold"
```

Expected: commit succeeds, `git log --oneline` shows 1 commit.

- [ ] **Step 3: Verify scaffold**

```bash
git log --oneline && git status
```

Expected: 1 commit shown; `git status` reports a clean tree.

---

### Task 1: Base conversions (`conversions.go`)

**Files:**
- Create: `conversions.go`
- Test: `conversions_test.go`

**Interfaces:**
- Consumes: nothing.
- Produces:
  - `convertHex(s string) (string, bool)` — decimal string of s parsed as base-16, ok=false if invalid.
  - `convertBin(s string) (string, bool)` — decimal string of s parsed as base-2, ok=false if invalid.

- [ ] **Step 1: Write the failing test**

`conversions_test.go`:

```go
package main

import "testing"

func TestConvertHex(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
		ok   bool
	}{
		{"uppercase", "1E", "30", true},
		{"lowercase", "1e", "30", true},
		{"zero", "0", "0", true},
		{"large", "FF", "255", true},
		{"invalid", "XYZ", "", false},
		{"empty", "", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := convertHex(tc.in)
			if ok != tc.ok {
				t.Fatalf("convertHex(%q) ok=%v, want ok=%v", tc.in, ok, tc.ok)
			}
			if got != tc.want {
				t.Fatalf("convertHex(%q)=%q, want %q", tc.in, got, tc.want)
			}
		})
	}
}

func TestConvertBin(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
		ok   bool
	}{
		{"simple", "10", "2", true},
		{"long", "1010", "10", true},
		{"zero", "0", "0", true},
		{"invalid", "12", "", false},
		{"invalidDigit", "10AB", "", false},
		{"empty", "", "", false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := convertBin(tc.in)
			if ok != tc.ok {
				t.Fatalf("convertBin(%q) ok=%v, want ok=%v", tc.in, ok, tc.ok)
			}
			if got != tc.want {
				t.Fatalf("convertBin(%q)=%q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./... -run TestConvert
```

Expected: FAIL — `convertHex` / `convertBin` not defined (the package reference error is the expected failing state). If `go test` reports a "no non-test Go files" style error for the package, that is still the expected FAIL for this step.

- [ ] **Step 3: Write minimal implementation**

`conversions.go`:

```go
package main

import "strconv"

func convertHex(s string) (string, bool) {
	v, err := strconv.ParseInt(s, 16, 64)
	if err != nil {
		return "", false
	}
	return strconv.FormatInt(v, 10), true
}

func convertBin(s string) (string, bool) {
	v, err := strconv.ParseInt(s, 2, 64)
	if err != nil {
		return "", false
	}
	return strconv.FormatInt(v, 10), true
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./... -run TestConvert
```

Expected: ok, all subtests PASS.

- [ ] **Step 5: Commit**

```bash
git add conversions.go conversions_test.go
git commit -m "feat: hex and bin base conversion utilities"
```

---

### Task 2: Marker normalization (`transform.go`)

**Files:**
- Create: `transform.go`
- Test: `transform_test.go`

**Interfaces:**
- Consumes: nothing.
- Produces:
  - `normalizeMarkers(s string) string` — rewrites every `(cmd)` / `(cmd, N)` / `(cmd,N)` occurrence (with any internal whitespace) to canonical `(cmd)` or `(cmd,N)` form. Unknown or malformed parens are left untouched.

- [ ] **Step 1: Write the failing test**

`transform_test.go`:

```go
package main

import "testing"

func TestNormalizeMarkers(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{"bare", "(up)", "(up)"},
		{"spaced inside", "( up )", "(up)"},
		{"number spaced", "(up, 2)", "(up,2)"},
		{"number tight", "(up,2)", "(up,2)"},
		{"space before comma", "(up , 2)", "(up,2)"},
		{"hex", "(hex)", "(hex)"},
		{"bin", "(bin)", "(bin)"},
		{"low", "(low)", "(low)"},
		{"cap", "(cap)", "(cap)"},
		{"unknown left alone", "(foo, 3)", "(foo, 3)"},
		{"not a marker", "hello world", "hello world"},
		{"text untouched", "a (up) b", "a (up) b"},
		{"non-model text", "(x)", "(x)"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := normalizeMarkers(tc.in); got != tc.want {
				t.Fatalf("normalizeMarkers(%q)=%q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./... -run TestNormalizeMarkers
```

Expected: FAIL — `normalizeMarkers` not defined.

- [ ] **Step 3: Write minimal implementation**

Append to `transform.go`:

```go
package main

import (
	"regexp"
)

var markerRegex = regexp.MustCompile(`\(\s*(hex|bin|up|low|cap)\s*(?:,\s*(\d+)\s*)?\)`)

func normalizeMarkers(s string) string {
	return markerRegex.ReplaceAllStringFunc(s, func(m string) string {
		subs := markerRegex.FindStringSubmatch(m)
		if subs[2] != "" {
			return "(" + subs[1] + "," + subs[2] + ")"
		}
		return "(" + subs[1] + ")"
	})
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./... -run TestNormalizeMarkers
```

Expected: ok, all subtests PASS.

- [ ] **Step 5: Commit**

```bash
git add transform.go transform_test.go
git commit -m "feat: normalize transformer markers"
```

---

### Task 3: hex/bin marker conversion (`convertMarkers`)

**Files:**
- Modify: `transform.go`
- Modify: `transform_test.go`

**Interfaces:**
- Consumes: `convertHex`, `convertBin` (Task 1).
- Produces:
  - `convertMarkers(tokens []string) []string` — for each `(hex)`/`(bin)` token, replaces the previous token with its decimal value and drops the marker; if the previous token is not a valid operand, keeps the marker as-is.

- [ ] **Step 1: Write the failing test**

Append to `transform_test.go`:

```go
func TestConvertMarkers(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want []string
	}{
		{
			"hex converts previous word",
			[]string{"1E", "(hex)", "files", "were", "added"},
			[]string{"30", "files", "were", "added"},
		},
		{
			"bin converts previous word",
			[]string{"It", "has", "been", "10", "(bin)", "years"},
			[]string{"It", "has", "been", "2", "years"},
		},
		{
			"invalid hex marker kept",
			[]string{"XYZ", "(hex)"},
			[]string{"XYZ", "(hex)"},
		},
		{
			"marker at start kept",
			[]string{"(hex)", "value"},
			[]string{"(hex)", "value"},
		},
		{
			"multiple markers",
			[]string{"42", "(hex)", "and", "10", "(bin)"},
			[]string{"66", "and", "2"},
		},
		{
			"no markers unchanged",
			[]string{"hello", "world"},
			[]string{"hello", "world"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := convertMarkers(tc.in)
			if len(got) != len(tc.want) {
				t.Fatalf("convertMarkers(%v) length=%d, want %d: %v", tc.in, len(got), len(tc.want), got)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("convertMarkers(%v)=%v, want %v", tc.in, got, tc.want)
				}
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./... -run TestConvertMarkers
```

Expected: FAIL — `convertMarkers` not defined.

- [ ] **Step 3: Write minimal implementation**

Append to `transform.go`:

```go
func convertMarkers(tokens []string) []string {
	out := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		if len(out) == 0 {
			out = append(out, tok)
			continue
		}
		switch tok {
		case "(hex)":
			if val, ok := convertHex(out[len(out)-1]); ok {
				out[len(out)-1] = val
			} else {
				out = append(out, tok)
			}
		case "(bin)":
			if val, ok := convertBin(out[len(out)-1]); ok {
				out[len(out)-1] = val
			} else {
				out = append(out, tok)
			}
		default:
			out = append(out, tok)
		}
	}
	return out
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./... -run TestConvertMarkers
```

Expected: ok, all cases PASS.

- [ ] **Step 5: Commit**

```bash
git add transform.go transform_test.go
git commit -m "feat: convert hex and bin markers to decimal"
```

---

### Task 4: Case transformations (`applyCase`)

**Files:**
- Modify: `transform.go`
- Modify: `transform_test.go`

**Interfaces:**
- Consumes: nothing new.
- Produces:
  - `applyCase(tokens []string) []string` — applies `(up)`/`(low)`/`(cap)` to the preceding word, or `(up,N)`/`(low,N)`/`(cap,N)` to the N preceding words (fewer if fewer exist), then removes the marker.
  - `parseCaseMarker(tok string) (cmd string, n int, ok bool)` — returns command name, word count (default 1), false for non-markers or malformed.
  - `capitalize(s string) string` — uppercases first byte, preserves the rest.

- [ ] **Step 1: Write the failing test**

Append to `transform_test.go`:

```go
func TestApplyCase(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want []string
	}{
		{
			"up single",
			[]string{"Ready,", "set,", "go", "(up)", "!"},
			[]string{"Ready,", "set,", "GO", "!"},
		},
		{
			"low single",
			[]string{"I", "should", "stop", "SHOUTING", "(low)"},
			[]string{"I", "should", "stop", "shouting"},
		},
		{
			"cap single",
			[]string{"Welcome", "to", "the", "Brooklyn", "bridge", "(cap)"},
			[]string{"Welcome", "to", "the", "Brooklyn", "Bridge"},
		},
		{
			"up with count",
			[]string{"This", "is", "so", "exciting", "(up,2)"},
			[]string{"This", "is", "SO", "EXCITING"},
		},
		{
			"cap with count",
			[]string{"it", "was", "the", "age", "of", "foolishness", "(cap,6)"},
			[]string{"It", "Was", "The", "Age", "Of", "Foolishness"},
		},
		{
			"low with count",
			[]string{"IT", "WAS", "THE", "(low,3)", "winter", "of", "despair."},
			[]string{"it", "was", "the", "winter", "of", "despair."},
		},
		{
			"count exceeds available",
			[]string{"ONE", "TWO", "(low,5)"},
			[]string{"one", "two"},
		},
		{
			"unknown marker kept",
			[]string{"hello", "(foo)"},
			[]string{"hello", "(foo)"},
		},
		{
			"marker at start",
			[]string{"(cap)", "hello"},
			[]string{"(cap)", "hello"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := applyCase(tc.in)
			if len(got) != len(tc.want) {
				t.Fatalf("applyCase(%v) length=%d, want %d: %v", tc.in, len(got), len(tc.want), got)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("applyCase(%v)=%v, want %v", tc.in, got, tc.want)
				}
			}
		})
	}
}

func TestParseCaseMarker(t *testing.T) {
	cases := []struct {
		name      string
		in        string
		wantCmd   string
		wantN     int
		wantParse bool
	}{
		{"up", "(up)", "up", 1, true},
		{"low", "(low)", "low", 1, true},
		{"cap", "(cap)", "cap", 1, true},
		{"up count", "(up,3)", "up", 3, true},
		{"zero count invalid", "(up,0)", "", 0, false},
		{"alpha count invalid", "(up,x)", "", 0, false},
		{"unknown", "(foo)", "", 0, false},
		{"not a marker", "hello", "", 0, false},
		{"hex not case", "(hex)", "", 0, false},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cmd, n, ok := parseCaseMarker(tc.in)
			if ok != tc.wantParse || cmd != tc.wantCmd || n != tc.wantN {
				t.Fatalf("parseCaseMarker(%q)=(%q,%d,%v), want (%q,%d,%v)",
					tc.in, cmd, n, ok, tc.wantCmd, tc.wantN, tc.wantParse)
			}
		})
	}
}

func TestCapitalize(t *testing.T) {
	cases := []struct{ in, want string }{
		{"bridge", "Bridge"},
		{"Bridge", "Bridge"},
		{"", ""},
		{"123", "123"},
	}
	for _, tc := range cases {
		if got := capitalize(tc.in); got != tc.want {
			t.Fatalf("capitalize(%q)=%q, want %q", tc.in, got, tc.want)
		}
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./... -run "TestApplyCase|TestParseCaseMarker|TestCapitalize"
```

Expected: FAIL — functions not defined.

- [ ] **Step 3: Write minimal implementation**

Append to `transform.go` (add `strconv` and `strings` to its imports):

```go
func applyCase(tokens []string) []string {
	out := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		cmd, n, ok := parseCaseMarker(tok)
		if !ok {
			out = append(out, tok)
			continue
		}
		start := len(out) - n
		if start < 0 {
			start = 0
		}
		for i := start; i < len(out); i++ {
			switch cmd {
			case "up":
				out[i] = strings.ToUpper(out[i])
			case "low":
				out[i] = strings.ToLower(out[i])
			case "cap":
				out[i] = capitalize(out[i])
			}
		}
	}
	return out
}

func parseCaseMarker(tok string) (string, int, bool) {
	if len(tok) < 4 || tok[0] != '(' || tok[len(tok)-1] != ')' {
		return "", 0, false
	}
	inner := tok[1 : len(tok)-1]
	if inner == "up" || inner == "low" || inner == "cap" {
		return inner, 1, true
	}
	comma := strings.IndexByte(inner, ',')
	if comma < 0 {
		return "", 0, false
	}
	cmd := inner[:comma]
	switch cmd {
	case "up", "low", "cap":
	default:
		return "", 0, false
	}
	n, err := strconv.Atoi(inner[comma+1:])
	if err != nil || n <= 0 {
		return "", 0, false
	}
	return cmd, n, true
}

func capitalize(s string) string {
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./... -run "TestApplyCase|TestParseCaseMarker|TestCapitalize"
```

Expected: ok, all PASS.

- [ ] **Step 5: Commit**

```bash
git add transform.go transform_test.go
git commit -m "feat: apply case transformation markers"
```

---

### Task 5: Article fix (`fixArticles`)

**Files:**
- Modify: `transform.go`
- Modify: `transform_test.go`

**Interfaces:**
- Consumes: nothing.
- Produces:
  - `fixArticles(tokens []string) []string` — rewrites standalone `a`→`an` and `A`→`An` when the next token starts with a vowel or `h` (case-insensitive).
  - `startsWithVowelOrH(s string) bool`.

- [ ] **Step 1: Write the failing test**

Append to `transform_test.go`:

```go
func TestFixArticles(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want []string
	}{
		{
			"vowel lower",
			[]string{"bearing", "a", "untold", "story"},
			[]string{"bearing", "an", "untold", "story"},
		},
		{
			"vowel upper",
			[]string{"There", "it", "was.", "A", "amazing", "rock!"},
			[]string{"There", "it", "was.", "An", "amazing", "rock!"},
		},
		{
			"h letter",
			[]string{"a", "hour"},
			[]string{"an", "hour"},
		},
		{
			"consonant unchanged",
			[]string{"a", "cat"},
			[]string{"a", "cat"},
		},
		{
			"last word unchanged",
			[]string{"I", "am", "a"},
			[]string{"I", "am", "a"},
		},
		{
			"non-lower a untouched",
			[]string{"magic", "A"},
			[]string{"magic", "A"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := fixArticles(tc.in)
			if len(got) != len(tc.want) {
				t.Fatalf("fixArticles(%v) length=%d, want %d", tc.in, len(got), len(tc.want))
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("fixArticles(%v)=%v, want %v", tc.in, got, tc.want)
				}
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./... -run TestFixArticles
```

Expected: FAIL — `fixArticles` not defined.

- [ ] **Step 3: Write minimal implementation**

Append to `transform.go`:

```go
func fixArticles(tokens []string) []string {
	out := append([]string{}, tokens...)
	for i := 0; i < len(out)-1; i++ {
		if (out[i] == "a" || out[i] == "A") && startsWithVowelOrH(out[i+1]) {
			if out[i] == "A" {
				out[i] = "An"
			} else {
				out[i] = "an"
			}
		}
	}
	return out
}

func startsWithVowelOrH(s string) bool {
	if s == "" {
		return false
	}
	switch s[0] {
	case 'a', 'e', 'i', 'o', 'u', 'h', 'A', 'E', 'I', 'O', 'U', 'H':
		return true
	}
	return false
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./... -run TestFixArticles
```

Expected: ok, all PASS.

- [ ] **Step 5: Commit**

```bash
git add transform.go transform_test.go
git commit -m "feat: fix a/an articles before vowel or h"
```

---

### Task 6: Punctuation and quotes (`punctuation.go`)

**Files:**
- Create: `punctuation.go`
- Test: `punctuation_test.go`

**Interfaces:**
- Consumes: nothing (token slice built by `ProcessText`).
- Produces:
  - `separatePunct(tokens []string) []string` — peels any leading punctuation run (from `.,!?:;`) off each token into its own unit token.
  - `assemble(tokens []string) string` — joins tokens with single spaces except: no space before a punctuation unit, no space between an opening `'` and the next token, no space between a token and a closing `'`, always closing/opening quotes paired in order.

- [ ] **Step 1: Write the failing test**

`punctuation_test.go`:

```go
package main

import "testing"

func TestSeparatePunct(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want []string
	}{
		{
			"leading punctuation",
			[]string{"there", ",and", "then"},
			[]string{"there", ",", "and", "then"},
		},
		{
			"punct group",
			[]string{"are", "...", "kinda"},
			[]string{"are", "...", "kinda"},
		},
		{
			"question prefix",
			[]string{"boring", ",what", "do"},
			[]string{"boring", ",", "what", "do"},
		},
		{
			"exclamation suffix untouched",
			[]string{"BAMM", "!!"},
			[]string{"BAMM", "!!"},
		},
		{
			"quote untouched",
			[]string{"said:", "'", "I", "am"},
			[]string{"said:", "'", "I", "am"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := separatePunct(tc.in)
			if len(got) != len(tc.want) {
				t.Fatalf("separatePunct(%v) length=%d, want %d: %v", tc.in, len(got), len(tc.want), got)
			}
			for i := range got {
				if got[i] != tc.want[i] {
					t.Fatalf("separatePunct(%v)=%v, want %v", tc.in, got, tc.want)
				}
			}
		})
	}
}

func TestAssemble(t *testing.T) {
	cases := []struct {
		name string
		in   []string
		want string
	}{
		{
			"punctuation attaches left",
			[]string{"I", "was", "sitting", "over", "there,", ",", "and", "then", "BAMM", "!!"},
			"I was sitting over there, and then BAMM!!",
		},
		{
			"punct group attaches left",
			[]string{"I", "was", "thinking", "...", "You", "were", "right"},
			"I was thinking... You were right",
		},
		{
			"question attaches left",
			[]string{"Punctuation", "tests", "are", "...", "kinda", "boring", ",", "what", "do", "you", "think", "?"},
			"Punctuation tests are... kinda boring, what do you think?",
		},
		{
			"single word quotes",
			[]string{"me:", "'", "awesome", "'"},
			"me: 'awesome'",
		},
		{
			"multi word quotes",
			[]string{"said:", "'", "I", "am", "the", "most", "well-known", "homosexual", "in", "the", "world", "'"},
			"said: 'I am the most well-known homosexual in the world'",
		},
		{
			"no special tokens",
			[]string{"hello", "world"},
			"hello world",
		},
		{
			"single quote left alone",
			[]string{"it", "'", "is"},
			"it ' is",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := assemble(tc.in); got != tc.want {
				t.Fatalf("assemble(%v)=%q, want %q", tc.in, got, tc.want)
			}
		})
	}
}
```

- [ ] **Step 2: Run test to verify it fails**

```bash
go test ./... -run TestAssemble -count=1
go test ./... -run TestSeparatePunct -count=1
```

Expected: FAIL — `assemble` / `separatePunct` not defined.

- [ ] **Step 3: Write minimal implementation**

`punctuation.go`:

```go
package main

import "strings"

const punctChars = ".,!?:;"

func isPunctChar(c byte) bool {
	return strings.ContainsRune(punctChars, rune(c))
}

func separatePunct(tokens []string) []string {
	out := make([]string, 0, len(tokens))
	for _, tok := range tokens {
		i := 0
		for i < len(tok) && isPunctChar(tok[i]) {
			i++
		}
		if i > 0 {
			out = append(out, tok[:i])
		}
		if i < len(tok) {
			out = append(out, tok[i:])
		}
	}
	return out
}

func assemble(tokens []string) string {
	isOpen := make([]bool, len(tokens))
	isClose := make([]bool, len(tokens))
	open := false
	for i, tok := range tokens {
		if tok == "'" {
			if open {
				isClose[i] = true
			} else {
				isOpen[i] = true
			}
			open = !open
		}
	}
	var sb strings.Builder
	for i, tok := range tokens {
		if i > 0 {
			sep := " "
			if isOpen[i-1] || isClose[i] || isPunctChar(tok[0]) {
				sep = ""
			}
			sb.WriteString(sep)
		}
		sb.WriteString(tok)
	}
	return sb.String()
}
```

- [ ] **Step 4: Run test to verify it passes**

```bash
go test ./... -count=1
```

Expected: ok, all PASS (Tasks 1–6).

- [ ] **Step 5: Commit**

```bash
git add punctuation.go punctuation_test.go
git commit -m "feat: format punctuation and quote pairs"
```

---

### Task 7: Orchestrator + CLI (`ProcessText`, `run`, `main`)

**Files:**
- Modify: `transform.go`
- Create: `main.go`
- Modify: `transform_test.go`
- Create: `main_test.go`

**Interfaces:**
- Consumes: `normalizeMarkers`, `convertMarkers`, `applyCase`, `fixArticles`, `separatePunct`, `assemble`.
- Produces:
  - `ProcessText(text string) string` — full pipeline.
  - `run(args []string) error` — CLI logic, returns errors instead of exiting.
  - `main()` — wires `run(os.Args[1:])`, prints error to stderr, exits 1.

- [ ] **Step 1: Write the failing tests**

Append to `transform_test.go` (acceptance tests for the spec's transformation examples):

```go
func TestProcessTextSpecExamples(t *testing.T) {
	cases := []struct {
		name string
		in   string
		want string
	}{
		{
			"full sample",
			"it (cap) was the best of times, it was the worst of times (up) , it was the age of wisdom, it was the age of foolishness (cap, 6) , it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, IT WAS THE (low, 3) winter of despair.",
			"It was the best of times, it was the worst of TIMES, it was the age of wisdom, It Was The Age Of Foolishness, it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, it was the winter of despair.",
		},
		{
			"hex and bin",
			"Simply add 42 (hex) and 10 (bin) and you will see the result is 68.",
			"Simply add 66 and 2 and you will see the result is 68.",
		},
		{
			"article a to an",
			"There is no greater agony than bearing a untold story inside you.",
			"There is no greater agony than bearing an untold story inside you.",
		},
		{
			"punctuation",
			"Punctuation tests are ... kinda boring ,what do you think ?",
			"Punctuation tests are... kinda boring, what do you think?",
		},
		{
			"quotes single",
			"I am exactly how they describe me: ' awesome '",
			"I am exactly how they describe me: 'awesome'",
		},
		{
			"quotes multi",
			"As Elton John said: ' I am the most well-known homosexual in the world '",
			"As Elton John said: 'I am the most well-known homosexual in the world'",
		},
		{
			"up punctuation example",
			"Ready, set, go (up) !",
			"Ready, set, GO!",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if got := ProcessText(tc.in); got != tc.want {
				t.Fatalf("ProcessText(%q)\n got: %q\nwant: %q", tc.in, got, tc.want)
			}
		})
	}
}
```

`main_test.go`:

```go
package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestRun(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		dir := t.TempDir()
		in := filepath.Join(dir, "sample.txt")
		out := filepath.Join(dir, "result.txt")
		content := "it (cap) was the best of times"
		if err := os.WriteFile(in, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
		if err := run([]string{in, out}); err != nil {
			t.Fatalf("run() unexpected error: %v", err)
		}
		got, err := os.ReadFile(out)
		if err != nil {
			t.Fatalf("output file not created: %v", err)
		}
		if string(got) != "It was the best of times" {
			t.Fatalf("output=%q, want %q", string(got), "It was the best of times")
		}
	})

	t.Run("wrong arg count", func(t *testing.T) {
		if err := run(nil); err == nil {
			t.Fatal("expected error for zero args")
		}
		if err := run([]string{"a", "b", "c"}); err == nil {
			t.Fatal("expected error for three args")
		}
	})

	t.Run("missing input file", func(t *testing.T) {
		err := run([]string{filepath.Join(t.TempDir(), "nope.txt"), filepath.Join(t.TempDir(), "out.txt")})
		if err == nil || !strings.Contains(err.Error(), "nope.txt") {
			t.Fatalf("expected error mentioning missing file, got %v", err)
		}
	})
}
```

- [ ] **Step 2: Run tests to verify they fail**

```bash
go test ./... -run TestProcessTextSpecExamples -count=1
go test ./... -run TestRun -count=1
```

Expected: FAIL — `ProcessText` / `run` not defined.

- [ ] **Step 3: Write minimal implementation**

Append to `transform.go` (add `strings` import if not already present):

```go
func ProcessText(text string) string {
	tokens := strings.Fields(normalizeMarkers(text))
	tokens = convertMarkers(tokens)
	tokens = applyCase(tokens)
	tokens = fixArticles(tokens)
	tokens = separatePunct(tokens)
	return assemble(tokens)
}
```

`main.go`:

```go
package main

import (
	"fmt"
	"os"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "Error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) != 2 {
		return fmt.Errorf("usage: go run . <input.txt> <output.txt>")
	}
	data, err := os.ReadFile(args[0])
	if err != nil {
		return err
	}
	return os.WriteFile(args[1], []byte(ProcessText(string(data))), 0o644)
}
```

- [ ] **Step 4: Run tests to verify they pass**

```bash
go test ./... -count=1
```

Expected: ok — every test across all tasks PASSES.

- [ ] **Step 5: Commit**

```bash
git add transform.go transform_test.go main.go main_test.go
git commit -m "feat: process text pipeline and CLI entrypoint"
```

---

### Task 8: Final quality pass

**Files:**
- (Verify) all files above.

**Interfaces:**
- Consumes: everything built in Tasks 1–7.

- [ ] **Step 1: Check formatting and vet**

```bash
gofmt -l .
go vet ./...
```

Expected: `gofmt -l .` prints no files; `go vet ./...` prints nothing and exits 0.

- [ ] **Step 2: Run the full test suite**

```bash
go test ./... -count=1
```

Expected: `ok` and 0 failures.

- [ ] **Step 3: Manual end-to-end run with the spec's sample**

Create `sample.txt` in the project root with exactly:

```text
it (cap) was the best of times, it was the worst of times (up) , it was the age of wisdom, it was the age of foolishness (cap, 6) , it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, IT WAS THE (low, 3) winter of despair.
```

Then:

```bash
go run . sample.txt result.txt
```

Expected: `result.txt` reads:

```text
It was the best of times, it was the worst of TIMES, it was the age of wisdom, It Was The Age Of Foolishness, it was the epoch of belief, it was the epoch of incredulity, it was the season of Light, it was the season of darkness, it was the spring of hope, it was the winter of despair.
```

- [ ] **Step 4: Run the three other spec examples manually**

```bash
printf 'Simply add 42 (hex) and 10 (bin) and you will see the result is 68.\n' > sample.txt
go run . sample.txt result.txt
cat result.txt
printf 'There is no greater agony than bearing a untold story inside you.\n' > sample.txt
go run . sample.txt result.txt
cat result.txt
printf 'Punctuation tests are ... kinda boring ,what do you think ?\n' > sample.txt
go run . sample.txt result.txt
cat result.txt
```

Expected outputs (in order):
`Simply add 66 and 2 and you will see the result is 68.`
`There is no greater agony than bearing an untold story inside you.`
`Punctuation tests are... kinda boring, what do you think?`

- [ ] **Step 5: Remove scratch files and commit**

```bash
rm -f sample.txt result.txt
git status
```

Expected: clean tree (only source files + docs tracked).

```bash
git add -A
git commit -m "chore: final verification pass"
```

Expected: commit succeeds; `git log --oneline` shows 10 commits.

---

## Self-Review

**Spec coverage:**
- hex/bin conversion → Tasks 1 + 3 ✓
- case transformations incl. `(cmd, N)` → Task 4 ✓
- punctuation incl. groups → Task 6 ✓
- quote pairs (single + multi-word) → Task 6 ✓
- article fix → Task 5 ✓
- CLI args + file IO + errors → Task 7 ✓
- tests/unit testing → every task + acceptance tests Task 7 ✓
- spec usage examples → Task 7 + Task 8 ✓

**Placeholder scan:** none — every step has exact code and commands.

**Type consistency:** `convertHex`/`convertBin` (Task 1) match Task 3 usage; `normalizeMarkers`/`convertMarkers`/`applyCase`/`fixArticles`/`separatePunct`/`assemble` all match their definitions; `run(args []string) error` consistent across Task 7 code and tests.