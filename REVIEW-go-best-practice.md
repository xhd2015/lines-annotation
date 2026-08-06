# Go Best-Practice Review: `lines-annotation`

**Scope:** codebase structure, package layout, API design, error handling, CLI readiness, flags, external commands, and `go:embed` — grounded in [go-best-practice](https://github.com/xhd2015) skill topics (`cli`, `flags-parsing` / less-flags, `cmd-exec`, `kool-create`, `go-embed-assets`) where applicable.

**Project type:** Go **library** (`github.com/xhd2015/lines-annotation`), ~4k LOC, no `cmd/` / `package main`. Consumed by [xgo](https://github.com/xhd2015/xgo) for coverage reporting.

**Review date:** 2026-08-06  
**Module:** `go 1.18` · tags through `v0.0.8`  
**Method:** static inspection + `go build ./...` (succeeds). No fixes applied in this pass.

---

## Executive summary

The library has a **clear domain split** (`model` → `load` / `ast` → `compute` / `merge` / `filter`) that matches how coverage annotations are built and combined. That layout is the main strength.

Against go-best-practice and idiomatic Go library design, the highest-impact gaps are:

1. **Broken / foot-gun public API** (`ProjectAnnotation.Clone` via unset function var).
2. **Widespread panic-based control flow** in library packages (hard for callers to handle).
3. **`model` package doc promise violated** (non-stdlib imports).
4. **No tests**, weak docs, and packaging hygiene issues.
5. **CLI/flags topics largely N/A today**, but `loadcov.LoadAllOptions` is already a natural CLI surface if a thin binary is added later with **less-flags**.

---

## What works well

| Area | Notes |
|------|--------|
| Layering intent | README and package roles: `model` (data), `load` (primitives), `compute` (derived), `merge`, `filter`, `map`, `ast` |
| Annotation type registry | `ProjectAnnotation.Types` + `Has` / `ShouldHave` / `Set` makes pipelines explicit |
| Idempotent compute | `Ensure*` wrappers skip work when types already present |
| High-level facade | `load/loadcov` orchestrates static + profile + git diff + filter (good “product API”) |
| AST abstraction | `ast.LoadInfo` / `File` interfaces keep loaders testable in principle |
| Path filter | `path/filter` with include/exclude patterns is reusable |
| Build | `go build ./...` succeeds on current tree |

---

## Findings (by severity)

### Critical

#### C1. `ProjectAnnotation.Clone` will nil-panic — merge hook never wired

**Where:** `model/project.go`

```go
// clone impl
var MergeAnnotationsInto func(res *ProjectAnnotation, annotations ...*ProjectAnnotation)

func (c *ProjectAnnotation) Clone() *ProjectAnnotation {
	b := &ProjectAnnotation{ ... }
	MergeAnnotationsInto(b, c) // never assigned anywhere in-repo
	return b
}
```

There is **no** assignment to `model.MergeAnnotationsInto` (no `init`, no wiring from `merge`). `load.Clone` correctly calls `merge.MergeAnnotationsInto` instead.

**Why it matters:** Public method on the core type is a latent production crash. The function-var pattern is an anti-pattern for breaking `model` ↔ `merge` import cycles.

**Recommended change:**

- Remove `model.MergeAnnotationsInto` and `(*ProjectAnnotation).Clone`, **or** implement clone purely in `model` without merge, **or** keep clone only as `merge.Clone` / `load.Clone`.
- Prefer explicit package ownership over mutable package-level function vars.

**go-best-practice angle:** library API safety (not a CLI topic, but foundational before any CLI wrapper).

---

#### C2. Library APIs panic on recoverable conditions

**Where (examples):**

| Location | Behavior |
|----------|----------|
| `model.ProjectAnnotation.MustHave` | panics on missing annotation type |
| `compute.*` (`BlockID_ForLine`, `ExecLabels_Block2Line`, …) | call `MustHave` |
| `compute.Changed_ForLineFromChanges` | panics if line changes missing |
| `filter.ReserveForLineView` | panics on missing diff (default option) |
| `load/git.LoadLineChangesWithFilter` | `panic(err)` inside filter callback |
| `merge.MergeAnnotationsInto` | panics if `res == nil` |
| `model.SortIBlocks` | panics on wrong type via `reflect` |

**Why it matters:** Callers (xgo, tools, future CLI) cannot turn many failures into exit codes / JSON errors without `recover`. go-best-practice **CLI UX** expects process-level errors on stderr with exit 1 — that requires library functions to return `error`.

**Recommended change:**

- Keep `MustHave` only for internal asserts if needed; export `ShouldHave` + `error` paths for public compute.
- Change signature of pipeline steps to `func(...) error` (or return `(*ProjectAnnotation, error)`).
- In `load/git`, never `panic` inside a filter: return or skip with error propagation.
- For `ReserveForLineView`, default `MissingDiffOption_Error` should return error, not panic.

---

### High

#### H1. `model` claims stdlib-only but imports `go-coverage`

**Where:** `model/doc.go` vs `model/file.go`

Doc:

> There should have no dependencies other than stdlib…

Actual:

```go
import (
	"github.com/xhd2015/go-coverage/git"
	"github.com/xhd2015/go-coverage/model"
)
```

`FileAnnotation.ChangeDetail` and `LineChanges` are typed from **another module**. `model/coverage` also imports `go-coverage/merge`.

**Why it matters:**

- Core “types only” package pulls transitive deps (and couples JSON/schema to go-coverage).
- Consumers cannot import `model` lightly as documented.

**Recommended change (pick one):**

1. **Own the types** in `lines-annotation/model` (`FileDetail`, `LineChanges`) and convert at the `load/git` boundary; restore true stdlib-only `model`.
2. Or **rewrite the doc** and treat go-coverage as a hard core dependency (document it clearly).

Option 1 aligns better with “minimal substantial model”.

---

#### H2. Zero tests

**Where:** no `*_test.go` anywhere; no CI config in-tree.

**Why it matters:** Merge/diff inheritance and label propagation are easy to regress. No doctest/unit surface for `loadcov.LoadAll`, `merge`, or compute pipelines.

**Recommended change:**

- Table-driven tests for pure functions: `path/filter`, `merge` line/func merge, `Block.ID` / compare, coverage label derivation.
- Fixture-based tests for `loadcov` / git diff with recorded file contents (avoid live git when possible).
- Optional: small `testdata/` trees; if assets grow, apply **go-embed-assets** (placeholder + optional fat fixtures) so `go test` always compiles.

---

#### H3. Bug / dead code in `ast.loadFile` error handling

**Where:** `ast/load_file.go`

```go
ast, syntaxErr := parser.ParseFile(...)
if err != nil {   // checks outer ReadFile err (already nil), not syntaxErr
	return nil, nil
}
```

- Wrong variable: parse failures are stored in `syntaxErr` but never returned here.
- Contrast `parseCode` in `load_code.go`, which correctly returns parse errors.

**Recommended change:** Align policies:

- Either fail load on syntax error (`return nil, syntaxErr`), or deliberately allow partial AST and document that `File.SyntaxError()` must be checked — but do not leave a misleading `if err != nil` that never fires.

---

#### H4. Dual git integration paths; opaque external process usage

**Where:**

- `load/git/load_files.go` → `github.com/xhd2015/gitops/git` (`DiffCommitFiles`, `CatFile`)
- `load/git/line_changes_and_file_change_detail.go` → `github.com/xhd2015/go-coverage/git` (`NewGitDiff`, …)

Two parallel APIs for “git diff + line changes”. Callers can pick inconsistent behavior.

**go-best-practice `cmd-exec`:** Prefer a single, well-wrapped execution layer (`github.com/xhd2015/xgo/support/cmd` when shelling out: Debug, Dir, Output, stderr discard). Today shell-out is hidden inside dependencies — fine if intentional, but **this repo should not grow raw `os/exec`**.

**Recommended change:**

- Consolidate on one git backend (preferably `gitops` or go-coverage, not both).
- Document required git binary / version assumptions.
- If any new process invocation is added here, use **`cmd-exec`** (`cmd.Debug().Dir(...).Output(...)` / `Run`) instead of ad-hoc `os/exec`.

---

#### H5. High-level pipeline is CLI-shaped but has no CLI (and no flag discipline)

**Where:** `load/loadcov.LoadAllOptions`

```go
type LoadAllOptions struct {
	Dir, Args, Profiles, Ref, DiffBase ...
	Include, Exclude []string
	OnlyChangedFiles bool
}
```

This is exactly the shape of a subcommand (`lines-ann load --ref … --diff-base … --profile …`). There is no `cmd/`, no `--help`, no version surface.

**Applicability of skill topics:**

| Topic | Status |
|-------|--------|
| `flags-parsing` / less-flags | Not used; recommended **if** a CLI is added |
| `cli` (color, streaming, dry-run) | N/A until binary exists |
| `kool-create` | N/A for this library; only if scaffolding a demo server/UI |
| `go-embed-assets` | N/A (no UI assets); unused `_ "embed"` imports exist (see M5) |

**Recommended change (when a tool is wanted):**

```text
cmd/lines-annotation/main.go   # or under xgo
  lessflags
    String("--dir", &dir)
    String("--ref", &ref)
    String("--diff-base", &diffBase)
    StringSlice("--profile", &profiles)
    StringSlice("--include", &include)
    StringSlice("--exclude", &exclude)
    Bool("--only-changed", &onlyChanged)
    Help("-h,--help", help)
  → loadcov.LoadAll(...)
  → stream JSON (cli/streaming) or write file
```

- Wire **`--help` at every subcommand** (`flags-parsing/subcommand`).
- Prefer **one pipeline** with side-effect gates if dry-run is ever needed (`cli/dry-run`), not a second copy of LoadAll.
- Map library `error` → stderr + exit 1 (requires C2).

Until then: document `loadcov` as the supported programmatic entry and keep options structs stable.

---

### Medium

#### M1. Package path / name collisions and awkward import paths

| Path | Package name | Issue |
|------|--------------|--------|
| `./map` | `annotation_map` | Dir is Go keyword-ish; import path `.../map` forces aliases; name ≠ path |
| `./filter`, `./model/filter`, `./path/filter` | all `filter` | Call sites need aliases (`ann_filter`, `path_filter`, `filter_model`) |
| `./compute/fun` | `fun` | Cryptic; prefer `fn`/`funcinfo`/`funcexcl` with a full word |
| `./model/coverage` and `./compute/coverage` | both `coverage` | Alias required |

**Recommended change:**

- Rename `map` → `remap` or `annmap` (package + path aligned).
- Rename top-level `filter` → `annfilter` / `projectfilter`, or nest under `filter/project`.
- Expand `compute/fun` to a descriptive name.

---

#### M2. Non-idiomatic exported names (`Underscore_Style`)

**Where:** `compute` package — `BlockID_ForLine`, `ExecLabels_Block2Line`, `CoverageLabels_Line2Func`, `Changed_ForLineFromChanges`, etc.

Go convention is `BlockIDForLine` / `ExecLabelsBlockToLine` (or nested packages). Underscores in exported identifiers are rare outside generated code / CGO.

**Recommended change:** Gradual rename with deprecated aliases, or accept as house style and document it (pipeline notation: `Source_ToTarget`).

---

#### M3. Typos and naming debt in public / serialized surface

| Item | Location | Note |
|------|----------|------|
| `LodOpts` | `load/loadcov` | Likely `LoadOpts` |
| `requries` | `load_all.go` error string | Typo |
| `Incrimental` | `model/coverage/summary.go` | Documented intentional JSON typo — add `json` alias or versioned field |
| `inerited` | `model/func.go` comment | Docs only |
| README: `primitvie`, `Annoation` | README.md | Docs only |
| `compulteLineCoverage` | `compute/coverage` | unexported but noisy |
| `Slashlize` | via xgo fileutil | External; wrap as `ToSlash` locally if desired |

**Recommended change:** Fix non-serialized typos (`LodOpts`, error strings, README). For `Incrimental`, add correct JSON name with dual-read if wire format can change.

---

#### M4. `interface{}` extension bags without contracts

**Where:** `ProjectData`, `FileData`, `LineData`, `BlockData` — untyped `interface{}` plus parallel `*DataType string` fields.

**Why it matters:** No shared registry of type strings; merge/clone behavior is shallow / unclear; hard to test.

**Recommended change:** Small registry or typed unions for known kinds; document extension protocol for xgo.

---

#### M5. Blank `embed` imports with no `//go:embed`

**Where:** `load/loadcov/load_all.go`, `load_profile.go` — `_ "embed"` only.

**go-best-practice `go-embed-assets`:** Not applicable for SPA/assets here, but empty embed setup is noise. Remove unused imports. If embedding **testdata** later:

- Commit `placeholder.txt` under embed roots if generated trees are gitignored.
- Use `//go:embed testdata/**` so bare `go test` / module consumers always compile.

---

#### M6. Module / dependency hygiene

- **`go 1.18`** while toolchain is much newer — consider bumping (e.g. 1.21+) for security and stdlib features; keep min version intentional for consumers.
- **`xgo` pinned to a pseudo-version** — prefer tagged releases for reproducibility.
- Transitive **goja / regexp2** via xgo — heavy for a coverage annotation library; ensure only needed xgo subpackages are imported (`support/fileutil`, `support/goinfo`, `support/coverage`, `support/pattern`) to avoid pulling runtime JS stack into unexpected builds (verify with `go mod graph` / build tags if possible).

---

#### M7. Duplicate / shallow merge helpers

**Where:** `merge.MergeStringSet` and `MergeMap` are identical; `MergeLine` does not deep-clone nested structs (shared pointers after merge); `MergeStringSet` ignores `false` values in source sets inconsistently vs maps.

**Recommended change:** One helper; document merge semantics (OR for bools, first-wins for IDs, pointer replace for Remark/Code); deep-copy when merge must not alias.

---

#### M8. `SortIBlocks` uses `reflect` + panic

**Where:** `model/block.go`

Prefer generics (if go version allows) or concrete `[]*BlockAnnotation` / `[]*coverage.BlockStats` sort helpers. Keeps `model` free of reflection and panics.

---

### Low

#### L1. README incomplete for library consumers

- No import examples, no pipeline diagram (static → profile → diff → compute → summary).
- No mention of `loadcov.LoadAll` as the main entry.
- Spelling issues (see M3).
- “Why Lines Annotation” is useful; **Usage** section should show 10–20 lines of real API usage.

#### L2. Sparse `doc.go` coverage

Only `model`, `load`, `loadcov`, `merge` have package docs. Add one-liners for `compute`, `filter`, `ast`, `path`, `map`.

#### L3. Commented-out prerequisites

e.g. `filter/reserve_for_line_view.go` and `compute/changed_line2func.go` leave `MustHave` / TODO checks disabled — either enforce or document optional deps.

#### L4. `StaticChecksum_Col8bits` “temporary util for bugfix”

Public API clutter; move to internal or delete once fixed upstream.

#### L5. No `internal/` boundary

Everything is importable. Consider `internal/parsecov` for coverage line regex parsing (`load_profile.go`) so only stable facades stay public.

#### L6. `kool-create` / scaffolding

Not required for an existing library. If a **coverage review demo** (Go API + React UI) is planned, scaffold with **`kool create go-react`** / `server`, and ship UI via **`go-embed-assets`** (placeholders + release hydrate). Out of scope for current pure-library layout.

---

## Package layout assessment

```text
lines-annotation/
  model/           # core types (+ accidental go-coverage coupling)
    coverage/      # binary/block profile types
    filter/        # filter option DTOs
    lang/          # go/js profile language
  ast/             # LoadInfo, parse files/dirs
    ast2ann/       # AST → block profile
  load/            # primitive loaders (empty lines, funcs, blocks)
    git/           # git line changes (two backends)
    loadcov/       # high-level coverage load facade
  compute/         # derived fields (labels, changed, uncoverable)
    coverage/      # summary aggregation
    fun/           # first-line excluded
  merge/           # annotation merge
  filter/          # project-level filtering / line-view reserve
  map/             # path prefix remap (package annotation_map)
  path/            # path trim + path/filter
```

**Dependency direction (intended):**

```text
model ← merge, compute, filter, map, load, ast
ast ← load, loadcov
load ← loadcov
compute ← loadcov, filter, map
```

**Friction:**

- `model.Clone` wants to call `merge` without importing it → broken function var (C1).
- `filter` → `compute` for `LineUncoverable` couples filter to compute (acceptable but worth noting).
- No `internal/` for implementation details.

**Verdict:** Layout is **conceptually sound** for a multi-stage annotation pipeline. Rename collisions (M1) and model purity (H1) are the main structural fixes.

---

## CLI / flags / embed / cmd matrix (skill topics)

| Topic | Applies now? | Recommendation |
|-------|----------------|----------------|
| **flags-parsing / less-flags** | No binary | When adding `cmd/`, parse with less-flags; Help on every subcommand; StringSlice for profiles/include/exclude |
| **cli/color** | No | Only if TTY reports; honor `NO_COLOR` |
| **cli/streaming** | No | Stream NDJSON per file for large projects; avoid buffering full `ProjectAnnotation` print unless `--output file` |
| **cli/dry-run** | No | Single LoadAll path; gate writes (if any) with `if !dryRun` |
| **cli/skill-cli** | No | N/A unless packaging this as a skill |
| **cmd-exec** | Indirect | Keep git behind libraries; new exec → `xgo/support/cmd` |
| **go-embed-assets** | No assets | Remove dead `_ "embed"`; use pattern for test fixtures if needed |
| **kool-create** | No | Only for optional demo app |

---

## Prioritized recommendation roadmap

### Phase 1 — Correctness / API safety (do first)

1. Fix or remove `ProjectAnnotation.Clone` / `MergeAnnotationsInto` (C1).
2. Replace panic paths with `error` returns on public APIs (C2); start with `load/git` and `ReserveForLineView`.
3. Fix `ast.loadFile` error handling (H3).
4. Align `model` docs with deps or decouple go-coverage types (H1).

### Phase 2 — Quality bar

5. Add unit tests for merge, path filter, compute pure helpers (H2).
6. Consolidate git backends (H4).
7. Rename `LodOpts`, fix error/README typos (M3).
8. Remove unused `embed` imports (M5).

### Phase 3 — Ergonomics / optional CLI

9. Rename packages for import clarity (`map` → `remap`, etc.) (M1).
10. Optional `cmd/lines-annotation` with **less-flags** wrapping `loadcov.LoadAll` (H5).
11. Bump Go version + pin xgo to a tag (M6).
12. Expand README with a real end-to-end example (L1).

---

## Suggested non-goals (for this library)

- Do **not** force a full skill-CLI or TUI onto a pure library.
- Do **not** introduce `kool-create` scaffolding into this repo unless building a demo product.
- Do **not** fat-embed a coverage UI here unless product scope expands; xgo can own UI.

---

## Appendix A — Inventory

| Metric | Value |
|--------|--------|
| Go packages (`go list ./...`) | 17 |
| Approx. LOC | ~3981 |
| `package main` / `cmd/` | none |
| `*_test.go` | none |
| Direct modules | gitops, go-coverage, xgo |
| Tags observed | v0.0.1 … v0.0.8 |

## Appendix B — Key entrypoints for reviewers / consumers

| Entrypoint | Role |
|------------|------|
| `loadcov.LoadAll` | Full project annotation pipeline |
| `loadcov.LoadStatic` | AST-derived static structure |
| `loadcov.LoadCoverageProfileFiles` | Go cover profiles → annotation |
| `loadcov.LoadGitDiff` / `load/git.LoadFiles` | Diff + line changes |
| `compute.*` / `compute/coverage.ComputeCoverageSummary` | Derived labels + summary |
| `merge.MergeAnnotations` | Combine annotation layers |
| `filter.ReserveForLineView` | Slim view for UI |

---

## Appendix C — Traceability to go-best-practice topics

| Finding | Topics |
|---------|--------|
| C2, H5 | `cli` (error UX), future `flags-parsing` |
| H4 | `cmd-exec` |
| H5 | `flags-parsing`, `flags-parsing/subcommand`, `cli/streaming`, `cli/dry-run` |
| M5, H2 fixtures | `go-embed-assets` |
| L6 | `kool-create` |
| C1, H1, M1, M2 | General Go package layout / library API design (skill index context) |

---

*End of review. Implementation intentionally deferred per request.*
