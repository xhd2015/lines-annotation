package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	annfilter "github.com/xhd2015/lines-annotation/filter"
	"github.com/xhd2015/lines-annotation/model"
)

func TestGitAndAnnotationFileUniversesAgree(t *testing.T) {
	dir := t.TempDir()
	run := func(args ...string) string {
		t.Helper()
		cmd := exec.Command("git", args...)
		cmd.Dir = dir
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
		return strings.TrimSpace(string(out))
	}
	run("init", "-q")
	paths := []string{"src/a.go", "src/a_test.go", "vendor/a.go", "a/vendor/a.go", "test/api_test/runner/vendor/a.go", "a/node_modules/a.go", "vendor2/a.go", "my_vendor/a.go", "src/a.txt"}
	for _, file := range paths {
		target := filepath.Join(dir, file)
		if err := os.MkdirAll(filepath.Dir(target), 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(target, []byte("package fixture\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	run("add", ".")
	run("-c", "user.name=fixture", "-c", "user.email=fixture@example.invalid", "commit", "-qm", "fixture")
	head := run("rev-parse", "HEAD")
	for _, projectExclude := range [][]string{nil, {"test/"}, {"src/"}} {
		opts := &LoadLineChangesOptions{Exclude: projectExclude}
		annotation, err := LoadLineChanges(dir, head, head, opts)
		if err != nil {
			t.Fatal(err)
		}
		check := annfilter.NewChecker(annfilter.MakeProjectFilter("", nil, projectExclude))
		for _, file := range paths {
			_, present := annotation.Files[model.RelativeFile(file)]
			if present != check.MatchFile(file) {
				t.Fatalf("Git/annotation policy differs at %s excludes=%v: git=%v", file, projectExclude, present)
			}
		}
	}
}
