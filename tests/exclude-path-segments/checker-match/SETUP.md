# Scenario

**Feature**: NewChecker / MatchFile honors ExcludePathSegments

```
Request{Mode=checker-match, Options fields, File}
  -> filter.NewChecker(opts)
  -> MatchFile(file)  # true = keep
```

## Preconditions

- Composite Options: segments alone, prefix Exclude alone, or both.

## Steps

1. Set Mode to checker-match.
2. Append `checker-match` to DecisionPath.

## Context

- MatchFile false means rejected/excluded.
- Non-empty ExcludePathSegments alone must yield a non-nil checker.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "checker-match")
	req.Mode = ModeCheckerMatch
	return nil
}
```
