# Scenario

**Feature**: LoadLineChanges with segments drops vendor/node_modules files

```
opts.ExcludePathSegments=[vendor,node_modules]
  -> LoadLineChanges
  -> src/foo.go present; vendor/* and node_modules/* absent
```

## Preconditions

- Git fixture from parent Setup.
- Segments list is the policy CodeLens worker will later pass.

## Steps

1. Set Segments `["vendor","node_modules"]`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "with-segments")
	req.Segments = []string{"vendor", "node_modules"}
	return nil
}
```
