# Scenario

**Feature**: optional ExcludePathSegments on lines-annotation filter/load APIs

```
# harness narrows Request along the tree, then Run hits public APIs
Setup chain -> Request{Mode, File|GitDir, Segments, ...}
  -> HasExcludedPathSegment | NewChecker.MatchFile | LoadLineChanges
  -> Response for Assert
```

## Preconditions

- Working tree is module `github.com/xhd2015/lines-annotation` (this worktree).
- Parallel-safe: no `t.Setenv` / `t.Chdir` / process cwd reliance.
- Absolute paths only via `t.TempDir()` and `d.DOCTEST_*` fields.
- Classic TDD RED until production implements the public surface called by `Run`.

## Steps

1. Initialize empty `DecisionPath` when nil.
2. Branch/leaf Setup appends path and sets Mode + case fields.
3. Root does not create git fixtures (only `load-line-changes` grouping does).

## Context

- Product policy (`vendor`/`node_modules`) belongs in CodeLens worker (P2), not here.
- Pure segment + Checker are the bulk; LoadLineChanges is a thin integration leaf.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	if req.DecisionPath == nil {
		req.DecisionPath = []string{}
	}
	return nil
}
```
