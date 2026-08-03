# Scenario

**Feature**: clean application path is not segment-excluded

```
file src/a.go + segments [vendor, node_modules]
  -> HasExcludedPathSegment -> false
```

## Preconditions

- Path has no component named `vendor` or `node_modules`.

## Steps

1. Set File `src/a.go`, Segments `["vendor","node_modules"]`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "clean-src")
	req.File = "src/a.go"
	req.Segments = []string{"vendor", "node_modules"}
	return nil
}
```
