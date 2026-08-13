package git

import "testing"

func TestGitPathPatterns(t *testing.T) {
	got := GitPathPatterns(&LoadLineChangesOptions{
		ExcludePathSegments: []string{"vendor", "node_modules", ""},
	})
	want := []string{
		"*.go",
		":(exclude)*_test.go",
		":(exclude)vendor",
		":(exclude)vendor/**",
		":(exclude)node_modules",
		":(exclude)node_modules/**",
	}
	if len(got) != len(want) {
		t.Fatalf("len=%d want %d: %#v", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Fatalf("pat[%d]=%q want %q", i, got[i], want[i])
		}
	}
}
