# Scenario

**Feature**: prefix Exclude and path segments both apply (OR reject)

```
Options{Exclude + ExcludePathSegments}
  -> reject if prefix hits OR segment hits
```

## Preconditions

- Both prefix Exclude and non-empty segments configured.

## Steps

1. Append `combined`.
2. Default Exclude `["internal"]`, Segments `["vendor"]` (leaves may refine File).

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "combined")
	req.Exclude = []string{"internal"}
	req.Segments = []string{"vendor"}
	return nil
}
```
