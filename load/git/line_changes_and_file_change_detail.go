package git

import (
	"strings"

	"github.com/xhd2015/go-coverage/git"
	"github.com/xhd2015/go-coverage/line"

	"github.com/xhd2015/lines-annotation/model"
	"github.com/xhd2015/lines-annotation/model/filter"
	path_filter "github.com/xhd2015/lines-annotation/path/filter"
)

type LoadLineChangesOptions = filter.Options

// GitPathPatterns builds git pathspec from LoadLineChanges options so vendor /
// node_modules / non-.go files are excluded in ls-files and diff, not only in Go.
func GitPathPatterns(opts *LoadLineChangesOptions) []string {
	patterns := []string{"*.go", ":(exclude)*_test.go"}
	opts = filter.WithDefaultPathSegments(opts)
	for _, seg := range opts.ExcludePathSegments {
		if seg == "" {
			continue
		}
		patterns = append(patterns, ":(exclude)"+seg, ":(exclude)"+seg+"/**")
	}
	return patterns
}

func LoadLineChangesWithFilter(gitDir string, commitHash string, compareCommitHash string, check func(file string) bool) (*model.ProjectAnnotation, error) {
	segments := filter.DefaultExcludePathSegments()
	return loadLineChangesWithFilter(gitDir, commitHash, compareCommitHash, GitPathPatterns(nil), func(file string) bool {
		return !filter.HasExcludedPathSegment(file, segments) && (check == nil || check(file))
	})
}

func loadLineChangesWithFilter(gitDir string, commitHash string, compareCommitHash string, pathPatterns []string, filter func(file string) bool) (*model.ProjectAnnotation, error) {
	gitDiff := git.NewGitDiffWithPatterns(gitDir, compareCommitHash, commitHash, pathPatterns)

	// map the content
	changedFilesRanger, err := line.MakeGitDiffRanger(gitDiff, func(file string) bool {
		_, ok, err := gitDiff.GetOldFile(file)
		if err != nil {
			panic(err)
		}
		if !ok {
			// if it does not have a corresponding old file,
			// then don't modify anything
			return false
		}
		if filter != nil && !filter(file) {
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}

	changes, err := line.CollectChanges(changedFilesRanger)
	if err != nil {
		return nil, err
	}

	// for all files:
	//   if the file is new, do nothing
	//   if there is a diff, iterate each change to collect line changes
	allFiles, err := gitDiff.AllFilesDetailsV2()
	if err != nil {
		return nil, err
	}

	files := make(model.FileAnnotationMapping)
	for file, fd := range allFiles {
		if filter != nil && !filter(file) {
			continue
		}
		fileAnn := &model.FileAnnotation{
			ChangeDetail: fd,
		}

		if !fd.IsNew {
			if fd.ContentChanged {
				fileAnn.LineChanges = changes[file]
			}
		}
		files[model.RelativeFile(file)] = fileAnn
	}

	return &model.ProjectAnnotation{
		Files: files,
		Types: map[model.AnnotationType]bool{
			model.AnnotationType_LineChanges:  true,
			model.AnnotationType_ChangeDetail: true,
		},
	}, nil
}

// for
func LoadLineChanges(gitDir string, commitHash string, compareCommitHash string, opts *LoadLineChangesOptions) (*model.ProjectAnnotation, error) {
	opts = filter.WithDefaultPathSegments(opts)
	var fileFilter *path_filter.FileFilter
	if opts != nil && (len(opts.Include) > 0 || len(opts.Exclude) > 0) {
		fileFilter = path_filter.NewFileFilter(opts.Include, opts.Exclude)
	}
	return loadLineChangesWithFilter(gitDir, commitHash, compareCommitHash, GitPathPatterns(opts), func(file string) bool {
		if !strings.HasSuffix(file, ".go") {
			return false
		}
		if strings.HasSuffix(file, "_test.go") {
			return false
		}
		if fileFilter != nil && !fileFilter.MatchFile(file) {
			return false
		}
		if opts != nil && filter.HasExcludedPathSegment(file, opts.ExcludePathSegments) {
			return false
		}
		return true
	})
}
