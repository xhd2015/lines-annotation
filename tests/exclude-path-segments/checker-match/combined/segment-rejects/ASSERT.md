## Expected

- MatchFile false because path segment `vendor` matches (prefix Exclude does not).

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
		t.Fatal("expected non-nil checker")
	}
	if resp.MatchFile {
		t.Fatalf("MatchFile(%q) = true, want false (segment exclude)", req.File)
	}
}
```
