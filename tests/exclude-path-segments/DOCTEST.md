# ExcludePathSegments — L2 doctests

Classic TDD for optional path-**segment** exclusion on
`github.com/xhd2015/lines-annotation`. Callers pass names such as `vendor` /
`node_modules`; any exact path component match drops the file. The library stays
policy-free (no hardcoded defaults).

Primary surface: pure segment helper + `filter.NewChecker` / `MatchFile`.
Secondary: `load/git.LoadLineChanges` with a disposable two-commit git fixture.

## Version

0.0.2

# DSN (Domain Specific Notion)

Optional **ExcludePathSegments** on filter Options lets callers drop files by
exact path-component name early in load/filter.

**Participants**

- **Options** — `model/filter.Options` carries `ExcludePathSegments []string`
  plus existing Include/Exclude/suffix fields.
- **HasExcludedPathSegment** — pure helper: normalize `\`→`/`, split components,
  ignore empty segments; true when any component equals a listed name.
- **Checker** — `filter.NewChecker` / `MatchFile` applies suffix, prefix
  Include/Exclude, and segment exclusion (reject if any rule rejects).
- **LoadLineChanges** — `load/git.LoadLineChanges` applies the same segment rule
  in its early file filter (with existing `.go` / `_test.go` rules).
- **Test harness** — in-process L2: pure API calls, or a tiny temp git repo for load.

**Behaviors**

- Empty/nil segment list → no segment-based exclude (backward compatible).
- Exact segment only: `pkg/vendor/a.go` hits `vendor`; `vendored/a.go` and
  `src/vendor.go` do not.
- Prefix `Exclude` still works; combined with segments, either match rejects.
- Windows-style `\` separators normalize before component split.

## Decision Tree

```text
exclude-path-segments
|-- pure-segment                         # HasExcludedPathSegment (no git)
|   |-- hit                              # should exclude
|   |   |-- leading-vendor
|   |   |-- mid-path-vendor
|   |   |-- node-modules
|   |   |-- windows-separators
|   |   `-- multi-listed
|   `-- miss                             # should keep
|       |-- clean-src
|       |-- empty-segments
|       |-- nil-segments
|       |-- substring-not-segment
|       |-- filename-not-segment
|       `-- empty-entries-ignored
|-- checker-match                        # NewChecker + MatchFile
|   |-- segments-only
|   |   |-- hit-excluded
|   |   `-- miss-kept
|   |-- prefix-exclude-unchanged
|   `-- combined                         # prefix OR segment rejects
|       |-- prefix-rejects
|       `-- segment-rejects
`-- load-line-changes                    # tiny git fixture
    |-- with-segments
    `-- empty-segments-keeps-vendor
```

### Parameter significance (high → low)

1. **Surface** — pure segment helper vs Checker.MatchFile vs LoadLineChanges
2. **Outcome** — excluded (hit) vs kept (miss)
3. **Path / options shape** — position of segment, separators, empty/nil list,
   interaction with prefix Exclude

## Test Index

| Leaf | Scenario |
|------|----------|
| `pure-segment/hit/leading-vendor` | `vendor/a.go` + `[vendor]` → excluded |
| `pure-segment/hit/mid-path-vendor` | `pkg/vendor/a.go` + `[vendor]` → excluded |
| `pure-segment/hit/node-modules` | `node_modules/x/y.go` + `[node_modules]` → excluded |
| `pure-segment/hit/windows-separators` | `pkg\vendor\a.go` + `[vendor]` → excluded |
| `pure-segment/hit/multi-listed` | `node_modules/x.go` + `[vendor,node_modules]` → excluded |
| `pure-segment/miss/clean-src` | `src/a.go` + `[vendor,node_modules]` → kept |
| `pure-segment/miss/empty-segments` | `vendor/a.go` + `[]` → kept |
| `pure-segment/miss/nil-segments` | `vendor/a.go` + nil → kept |
| `pure-segment/miss/substring-not-segment` | `vendored/a.go` + `[vendor]` → kept |
| `pure-segment/miss/filename-not-segment` | `src/vendor.go` + `[vendor]` → kept |
| `pure-segment/miss/empty-entries-ignored` | `vendor/a.go` + `["",""]` → kept (empties ignored) |
| `checker-match/segments-only/hit-excluded` | Checker non-nil; MatchFile false for vendor path |
| `checker-match/segments-only/miss-kept` | Checker non-nil; MatchFile true for src path |
| `checker-match/prefix-exclude-unchanged` | `Exclude:["src"]`, empty segments; prefix still rejects |
| `checker-match/combined/prefix-rejects` | both set; file under prefix exclude rejected |
| `checker-match/combined/segment-rejects` | both set; file with vendor segment rejected |
| `load-line-changes/with-segments` | vendor/node_modules absent; `src/foo.go` present |
| `load-line-changes/empty-segments-keeps-vendor` | empty segments → vendor file still loaded |

## How to Run

From module root `lines-annotation-mad-max`:

```sh
doctest vet ./tests/exclude-path-segments
doctest test ./tests/exclude-path-segments
doctest test -v ./tests/exclude-path-segments
```

Classic TDD: leaves are expected **RED** until `ExcludePathSegments`,
`HasExcludedPathSegment`, Checker wiring, and LoadLineChanges filter exist.

```go
import (
	"fmt"
	"testing"

	filefilter "github.com/xhd2015/lines-annotation/filter"
	loadgit "github.com/xhd2015/lines-annotation/load/git"
	modelfilter "github.com/xhd2015/lines-annotation/model/filter"

	"github.com/xhd2015/doctest/session"
)

const (
	ModePureSegment     = "pure-segment"
	ModeCheckerMatch    = "checker-match"
	ModeLoadLineChanges = "load-line-changes"
)

type Request struct {
	Mode         string
	DecisionPath []string

	// pure-segment + checker-match
	File           string
	Segments       []string
	UseNilSegments bool // force nil segment slice (distinct from empty non-nil)

	// checker-match Options extras
	Include []string
	Exclude []string // prefix exclude patterns

	// load-line-changes
	GitDir     string
	Head       string
	Base       string
	ProbeFiles []string // relative paths to report presence for
}

type Response struct {
	// pure-segment
	HasExcluded bool

	// checker-match
	CheckerNil bool
	MatchFile  bool // true => keep file

	// load-line-changes
	Present map[string]bool
	Keys    []string
}

func Run(t *testing.T, d *session.Doctest, req *Request) (*Response, error) {
	t.Helper()
	segs := req.Segments
	if req.UseNilSegments {
		segs = nil
	}

	switch req.Mode {
	case ModePureSegment:
		// Public pure helper (to be added under model/filter).
		return &Response{
			HasExcluded: modelfilter.HasExcludedPathSegment(req.File, segs),
		}, nil

	case ModeCheckerMatch:
		opts := &modelfilter.Options{
			Include:             req.Include,
			Exclude:             req.Exclude,
			ExcludePathSegments: segs,
		}
		c := filefilter.NewChecker(opts)
		if c == nil {
			// Historical: nil checker means "match all". After the feature,
			// a non-empty ExcludePathSegments alone must produce a non-nil checker.
			return &Response{CheckerNil: true, MatchFile: true}, nil
		}
		return &Response{
			CheckerNil: false,
			MatchFile:  c.MatchFile(req.File),
		}, nil

	case ModeLoadLineChanges:
		if req.GitDir == "" || req.Head == "" || req.Base == "" {
			return nil, fmt.Errorf("load-line-changes requires GitDir, Head, Base from Setup")
		}
		opts := &modelfilter.Options{
			ExcludePathSegments: segs,
		}
		ann, err := loadgit.LoadLineChanges(req.GitDir, req.Head, req.Base, opts)
		if err != nil {
			return nil, err
		}
		present := make(map[string]bool, len(req.ProbeFiles))
		keys := make([]string, 0, len(ann.Files))
		for k := range ann.Files {
			keys = append(keys, string(k))
			present[string(k)] = true
		}
		// Ensure probe keys exist in map (false when absent)
		for _, p := range req.ProbeFiles {
			if _, ok := present[p]; !ok {
				present[p] = false
			}
		}
		return &Response{Present: present, Keys: keys}, nil

	default:
		return nil, fmt.Errorf("unknown Mode %q", req.Mode)
	}
}
```
