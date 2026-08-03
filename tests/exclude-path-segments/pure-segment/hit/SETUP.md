# Scenario

**Feature**: path has a listed segment → excluded

```
HasExcludedPathSegment(file, segments) -> true
```

## Preconditions

- Assert expects `HasExcluded == true`.

## Steps

1. Append `hit` to DecisionPath.

## Context

- Covers leading, mid-path, node_modules, Windows separators, multi-list any-match.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "hit")
	return nil
}
```
