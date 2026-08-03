# Scenario

**Feature**: empty strings in the segment list are ignored

```
file vendor/a.go + segments ["", ""]
  -> ignore empty entries → effective list empty
  -> HasExcludedPathSegment -> false
```

## Preconditions

- Segment list has only empty strings (no real names).

## Steps

1. Set File `vendor/a.go`, Segments `["", ""]`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "empty-entries-ignored")
	req.File = "vendor/a.go"
	req.Segments = []string{"", ""}
	return nil
}
```
