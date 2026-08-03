# Scenario

**Feature**: Checker keeps file without matching segment

```
Options{ExcludePathSegments:[vendor,node_modules]} + file src/pkg/a.go
  -> CheckerNil=false, MatchFile=true
```

## Preconditions

- Segments non-empty; clean src path.

## Steps

1. Set File `src/pkg/a.go`, Segments `["vendor","node_modules"]`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "miss-kept")
	req.File = "src/pkg/a.go"
	req.Segments = []string{"vendor", "node_modules"}
	return nil
}
```
