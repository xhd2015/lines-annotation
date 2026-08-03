# Scenario

**Feature**: nil segment list disables segment filtering

```
file vendor/a.go + segments nil
  -> HasExcludedPathSegment -> false
```

## Preconditions

- Segments passed as nil (UseNilSegments).

## Steps

1. Set File `vendor/a.go`, UseNilSegments true.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "nil-segments")
	req.File = "vendor/a.go"
	req.Segments = []string{"vendor"} // ignored when UseNilSegments
	req.UseNilSegments = true
	return nil
}
```
