package path

import (
	"fmt"
	"path/filepath"
	"strings"
)

// TrimPrefix (/a/b/c, /a/b) -> c
func TrimPrefix(s string, prefix string) (string, error) {
	return doTrimPathPrefix(s, prefix, false)
}

func TrimPrefixOrEmpty(s string, prefix string) string {
	res, _ := doTrimPathPrefix(s, prefix, true)
	return res
}

func doTrimPathPrefix(s string, prefix string, allowUnprefixed bool) (string, error) {
	if !strings.HasPrefix(s, prefix) {
		if !allowUnprefixed {
			return "", fmt.Errorf("string %s not prefixed with %s", s, prefix)
		}
		return "", nil
	}
	s = s[len(prefix):]
	return strings.TrimPrefix(s, string(filepath.Separator)), nil
}
