# Scenario

**Feature**: pure HasExcludedPathSegment without git or Checker

```
# segment helper only — no Include/Exclude/suffix, no LoadLineChanges
Request{Mode=pure-segment, File, Segments}
  -> model/filter.HasExcludedPathSegment
  -> Response.HasExcluded
```

## Preconditions

- Mode is pure-segment for all descendants.

## Steps

1. Set Mode to pure-segment.
2. Record decision path segment `pure-segment`.

## Context

- Fast L2 leaves; expected RED until helper exists.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "pure-segment")
	req.Mode = ModePureSegment
	return nil
}
```
