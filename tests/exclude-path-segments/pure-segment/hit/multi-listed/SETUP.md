# Scenario

**Feature**: any listed segment is enough (OR across segment names)

```
file node_modules/x.go + segments [vendor, node_modules]
  -> HasExcludedPathSegment -> true
```

## Preconditions

- Multiple segment names; file matches only one of them.

## Steps

1. Set File `node_modules/x.go`, Segments `["vendor","node_modules"]`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "multi-listed")
	req.File = "node_modules/x.go"
	req.Segments = []string{"vendor", "node_modules"}
	return nil
}
```
