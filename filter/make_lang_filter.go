package filter

import (
	"github.com/xhd2015/lines-annotation/model/filter"
	"github.com/xhd2015/lines-annotation/model/lang"
)

func MakeFilter(profileLang lang.ProfileLanguage) *filter.Options {
	opts := &filter.Options{ExcludePathSegments: filter.DefaultExcludePathSegments()}
	if profileLang.IsGo() {
		opts.Suffix = []string{".go"}
		opts.ExcludeSuffix = []string{"_test.go"}
	} else if profileLang == lang.ProfileLanguage_Js {
		opts.Suffix = []string{".js", ".ts"}
	}
	return opts
}

// MakeProjectFilter intersects system/language rules with project rules.
func MakeProjectFilter(profileLang lang.ProfileLanguage, include, exclude []string) *filter.Options {
	opts := MakeFilter(profileLang)
	opts.Include = append([]string(nil), include...)
	opts.Exclude = append([]string(nil), exclude...)
	return opts
}
