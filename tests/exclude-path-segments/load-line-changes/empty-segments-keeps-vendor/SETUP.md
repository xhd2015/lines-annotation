# Scenario

**Feature**: empty segments leave LoadLineChanges backward compatible

```
opts.ExcludePathSegments=[] 
  -> LoadLineChanges
  -> vendor/lib.go still present (no segment filter)
```

## Preconditions

- Git fixture from parent; empty segment list.

## Steps

1. Set Segments to empty slice.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "empty-segments-keeps-vendor")
	req.Segments = []string{}
	return nil
}
```
