## Expected

- `src/foo.go` is present in result.Files.
- `vendor/lib.go`, `pkg/vendor/nested.go`, and `node_modules/pkg/x.go` are absent.

## Errors

- Run must not return an error (git/load failures are test failures).

```go
import (
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Assert(t *testing.T, d *session.Doctest, req *Request, resp *Response, err error) {
	if err != nil {
		t.Fatalf("LoadLineChanges error: %v", err)
	}
	if resp == nil {
		t.Fatal("nil Response")
	}
	if !resp.Present["src/foo.go"] {
		t.Fatalf("expected src/foo.go present; keys=%v present=%v", resp.Keys, resp.Present)
	}
	for _, bad := range []string{"vendor/lib.go", "pkg/vendor/nested.go", "node_modules/pkg/x.go"} {
		if resp.Present[bad] {
			t.Fatalf("expected %s absent under ExcludePathSegments; keys=%v", bad, resp.Keys)
		}
	}
}
```
