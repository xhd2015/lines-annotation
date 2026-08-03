# Scenario

**Feature**: empty (non-nil) segment list disables segment filtering

```
file vendor/a.go + segments []
  -> HasExcludedPathSegment -> false
```

## Preconditions

- Segments is non-nil empty slice.

## Steps

1. Set File `vendor/a.go`, Segments empty slice.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "empty-segments")
	req.File = "vendor/a.go"
	req.Segments = []string{}
	req.UseNilSegments = false
	return nil
}
```
