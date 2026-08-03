# Scenario

**Feature**: mid-path segment `vendor` excludes nested vendored package

```
file pkg/vendor/a.go + segments [vendor]
  -> HasExcludedPathSegment -> true
```

## Preconditions

- `vendor` is neither first nor last component.

## Steps

1. Set File `pkg/vendor/a.go`, Segments `["vendor"]`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "mid-path-vendor")
	req.File = "pkg/vendor/a.go"
	req.Segments = []string{"vendor"}
	return nil
}
```
