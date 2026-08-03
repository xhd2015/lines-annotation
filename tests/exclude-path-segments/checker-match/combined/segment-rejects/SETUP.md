# Scenario

**Feature**: under combined options, segment exclude alone is enough to reject

```
Exclude:[internal] + segments:[vendor] + file pkg/vendor/a.go
  -> MatchFile=false (segment; not under internal/)
```

## Preconditions

- File matches segment `vendor`, not prefix `internal`.

## Steps

1. Set File `pkg/vendor/a.go`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "segment-rejects")
	req.File = "pkg/vendor/a.go"
	return nil
}
```
