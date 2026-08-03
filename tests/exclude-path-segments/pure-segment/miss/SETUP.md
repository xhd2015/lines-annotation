# Scenario

**Feature**: path does not match any segment rule → kept

```
HasExcludedPathSegment(file, segments) -> false
```

## Preconditions

- Assert expects `HasExcluded == false`.

## Steps

1. Append `miss` to DecisionPath.

## Context

- Empty/nil lists, clean src, substring/filename false positives, empty list entries.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "miss")
	return nil
}
```
