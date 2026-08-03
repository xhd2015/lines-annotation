# Scenario

**Feature**: substring of a component is not an exact segment match

```
file vendored/a.go + segments [vendor]
  -> HasExcludedPathSegment -> false
```

## Preconditions

- Component is `vendored`, not `vendor`.

## Steps

1. Set File `vendored/a.go`, Segments `["vendor"]`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "substring-not-segment")
	req.File = "vendored/a.go"
	req.Segments = []string{"vendor"}
	return nil
}
```
