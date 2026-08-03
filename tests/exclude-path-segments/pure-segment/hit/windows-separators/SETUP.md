# Scenario

**Feature**: backslash separators normalize before segment split

```
file pkg\vendor\a.go + segments [vendor]
  -> normalize \ -> /
  -> HasExcludedPathSegment -> true
```

## Preconditions

- Input uses Windows-style `\` separators (string form; OS-independent).

## Steps

1. Set File to `pkg\vendor\a.go` (raw backslashes), Segments `["vendor"]`.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "windows-separators")
	req.File = `pkg\vendor\a.go`
	req.Segments = []string{"vendor"}
	return nil
}
```
