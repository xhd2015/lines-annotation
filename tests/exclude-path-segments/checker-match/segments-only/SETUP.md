# Scenario

**Feature**: Options with only ExcludePathSegments (no Include/Exclude prefix)

```
Options{ExcludePathSegments} -> NewChecker non-nil -> MatchFile
```

## Preconditions

- Include and Exclude prefix lists stay empty.

## Steps

1. Append `segments-only`.
2. Clear Include/Exclude.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "segments-only")
	req.Include = nil
	req.Exclude = nil
	return nil
}
```
