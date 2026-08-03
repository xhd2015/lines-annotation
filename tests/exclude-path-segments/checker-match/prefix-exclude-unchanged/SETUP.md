# Scenario

**Feature**: existing prefix Exclude still works with empty segments

```
Options{Exclude:["src"], ExcludePathSegments:[]} + file src/a.go
  -> MatchFile=false (prefix exclude unchanged)
```

## Preconditions

- Segments empty; prefix Exclude is `src`.

## Steps

1. Set Exclude `["src"]`, Segments empty, File `src/a.go`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "prefix-exclude-unchanged")
	req.File = "src/a.go"
	req.Exclude = []string{"src"}
	req.Segments = []string{}
	return nil
}
```
