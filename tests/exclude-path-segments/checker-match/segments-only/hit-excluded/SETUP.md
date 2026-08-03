# Scenario

**Feature**: Checker rejects file with excluded path segment

```
Options{ExcludePathSegments:[vendor]} + file vendor/lib/a.go
  -> CheckerNil=false, MatchFile=false
```

## Preconditions

- Segments non-empty; file under `vendor/`.

## Steps

1. Set File `vendor/lib/a.go`, Segments `["vendor"]`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "hit-excluded")
	req.File = "vendor/lib/a.go"
	req.Segments = []string{"vendor"}
	return nil
}
```
