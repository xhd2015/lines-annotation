# Scenario

**Feature**: segment `node_modules` excludes dependency tree paths

```
file node_modules/x/y.go + segments [node_modules]
  -> HasExcludedPathSegment -> true
```

## Preconditions

- Path under `node_modules`.

## Steps

1. Set File `node_modules/x/y.go`, Segments `["node_modules"]`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "node-modules")
	req.File = "node_modules/x/y.go"
	req.Segments = []string{"node_modules"}
	return nil
}
```
