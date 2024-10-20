package filter

import (
	"github.com/xhd2015/lines-annotation/model"
	"github.com/xhd2015/lines-annotation/model/filter"
	path_filter "github.com/xhd2015/lines-annotation/path/filter"
)

func FilterFiles(project *model.ProjectAnnotation, opts *filter.Options) {
	if opts == nil {
		return
	}
	pathFilter := path_filter.NewFileFilter(opts.Include, opts.Exclude)
	FilterFilesWithCheck(project, func(file model.RelativeFile) bool {
		if !opts.MatchSuffix(string(file)) {
			return false
		}
		if !pathFilter.MatchFile(string(file)) {
			return false
		}
		return true
	})
}

func FilterFilesWithCheck(project *model.ProjectAnnotation, check func(file model.RelativeFile) bool) {
	if check == nil {
		return
	}
	files := project.Files
	newFiles := make(model.FileAnnotationMapping, len(files))
	for file, val := range files {
		if !check(file) {
			continue
		}
		newFiles[file] = val
	}
	project.Files = newFiles
}
