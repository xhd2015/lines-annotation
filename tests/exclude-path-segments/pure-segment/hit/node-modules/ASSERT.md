## Expected

- `node_modules` as a path component excludes the file.

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
	if !resp.HasExcluded {
		t.Fatalf("HasExcludedPathSegment(%q, %v) = false, want true", req.File, req.Segments)
	}
}
```
