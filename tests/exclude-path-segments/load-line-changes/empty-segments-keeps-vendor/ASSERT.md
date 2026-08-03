## Expected

- With empty ExcludePathSegments, `vendor/lib.go` remains in result.Files.
- `src/foo.go` remains present as well.

## Errors

- Run must not return an error.

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
		t.Fatalf("expected src/foo.go present; keys=%v", resp.Keys)
	}
	if !resp.Present["vendor/lib.go"] {
		t.Fatalf("empty segments must keep vendor/lib.go; keys=%v present=%v", resp.Keys, resp.Present)
	}
}
```
