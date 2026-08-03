# Scenario

**Feature**: LoadLineChanges applies ExcludePathSegments in its early file filter

```
# disposable two-commit repo (Dir=abs path, no Chdir)
git fixture -> src/foo.go, vendor/lib.go, pkg/vendor/nested.go, node_modules/pkg/x.go
  -> LoadLineChanges(gitDir, head, base, opts)
  -> Response.Present / Keys
```

## Preconditions

- `git` available on PATH.
- Parallel-safe: each leaf uses `t.TempDir()`; commands use `cmd.Dir` (never process Chdir).
- Commits use fixed author env on the child process only.

## Steps

1. Set Mode to load-line-changes.
2. Build a tiny git repo with base commit + head commit that adds probe files.
3. Store GitDir, Base, Head, ProbeFiles on Request.
4. Leaves set Segments (or empty) and assert Present map.

## Context

- LoadLineChanges already restricts to non-test `.go` files; fixtures are `.go` only.
- Segment filtering must drop keys before/at the same callback as include/exclude.

```go
import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/xhd2015/doctest/session"
)

func Setup(t *testing.T, d *session.Doctest, req *Request) error {
	req.DecisionPath = append(req.DecisionPath, "load-line-changes")
	req.Mode = ModeLoadLineChanges

	dir := t.TempDir()
	base, head, err := initTwoCommitGoRepo(t, dir)
	if err != nil {
		return err
	}
	req.GitDir = dir
	req.Base = base
	req.Head = head
	req.ProbeFiles = []string{
		"src/foo.go",
		"vendor/lib.go",
		"pkg/vendor/nested.go",
		"node_modules/pkg/x.go",
	}
	return nil
}

// initTwoCommitGoRepo creates a repo with an empty base tree and a head tree
// containing the probe .go files. Returns (baseHash, headHash).
func initTwoCommitGoRepo(t *testing.T, dir string) (base, head string, err error) {
	t.Helper()
	run := func(args ...string) (string, error) {
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=doctest",
			"GIT_AUTHOR_EMAIL=doctest@example.com",
			"GIT_COMMITTER_NAME=doctest",
			"GIT_COMMITTER_EMAIL=doctest@example.com",
		)
		var stdout, stderr bytes.Buffer
		cmd.Stdout = &stdout
		cmd.Stderr = &stderr
		if err := cmd.Run(); err != nil {
			return "", fmt.Errorf("git %v: %v\n%s", args, err, stderr.String())
		}
		return strings.TrimSpace(stdout.String()), nil
	}

	if _, err := run("init"); err != nil {
		return "", "", err
	}
	// Detached-friendly default branch name
	_, _ = run("checkout", "-b", "main")

	// Base commit: placeholder so Diff has a parent
	placeholder := filepath.Join(dir, "README")
	if err := os.WriteFile(placeholder, []byte("base\n"), 0o644); err != nil {
		return "", "", err
	}
	if _, err := run("add", "README"); err != nil {
		return "", "", err
	}
	if _, err := run("commit", "-m", "base"); err != nil {
		return "", "", err
	}
	base, err = run("rev-parse", "HEAD")
	if err != nil {
		return "", "", err
	}

	files := map[string]string{
		"src/foo.go":             "package src\n",
		"vendor/lib.go":          "package lib\n",
		"pkg/vendor/nested.go":   "package nested\n",
		"node_modules/pkg/x.go":  "package x\n",
	}
	for rel, body := range files {
		path := filepath.Join(dir, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return "", "", err
		}
		if err := os.WriteFile(path, []byte(body), 0o644); err != nil {
			return "", "", err
		}
		if _, err := run("add", rel); err != nil {
			return "", "", err
		}
	}
	if _, err := run("commit", "-m", "add go files"); err != nil {
		return "", "", err
	}
	head, err = run("rev-parse", "HEAD")
	if err != nil {
		return "", "", err
	}
	return base, head, nil
}
```
