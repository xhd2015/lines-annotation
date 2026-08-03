# Scenario

**Feature**: under combined options, prefix Exclude alone is enough to reject

```
Exclude:[internal] + segments:[vendor] + file internal/pkg/a.go
  -> MatchFile=false (prefix; no vendor segment required)
```

## Preconditions

- File matches prefix Exclude, not segment vendor.

## Steps

1. Set File `internal/pkg/a.go`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "prefix-rejects")
	req.File = "internal/pkg/a.go"
	return nil
}
```
