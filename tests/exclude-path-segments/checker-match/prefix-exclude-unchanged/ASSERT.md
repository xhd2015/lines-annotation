## Expected

- Prefix Exclude `src` still rejects `src/a.go` when segment list is empty.
- Feature must not break Include/Exclude prefix behavior.

## Errors

- Run must not return an error.

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("Run error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil Response")
	}
	if resp.CheckerNil {
		t.Fatal("NewChecker nil with Exclude prefix set; want non-nil")
	}
	if resp.MatchFile {
		t.Fatalf("MatchFile(%q) = true, want false (prefix Exclude)", req.File)
	}
}
```
