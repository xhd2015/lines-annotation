# Scenario

**Feature**: leading path segment `vendor` excludes the file

```
file vendor/a.go + segments [vendor]
  -> HasExcludedPathSegment -> true
```

## Preconditions

- File path starts with segment `vendor`.

## Steps

1. Set File `vendor/a.go`, Segments `["vendor"]`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "leading-vendor")
	req.File = "vendor/a.go"
	req.Segments = []string{"vendor"}
	return nil
}
```
