# Scenario

**Feature**: filename containing the word is not a path segment

```
file src/vendor.go + segments [vendor]
  -> components [src, vendor.go]
  -> HasExcludedPathSegment -> false
```

## Preconditions

- Final component is `vendor.go`, not `vendor`.

## Steps

1. Set File `src/vendor.go`, Segments `["vendor"]`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "filename-not-segment")
	req.File = "src/vendor.go"
	req.Segments = []string{"vendor"}
	return nil
}
```
