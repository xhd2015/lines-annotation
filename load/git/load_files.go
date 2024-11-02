package git

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/xhd2015/gitops/git"
	git2 "github.com/xhd2015/go-coverage/git"
	"github.com/xhd2015/go-coverage/line"
	cov_model "github.com/xhd2015/go-coverage/model"
	"github.com/xhd2015/lines-annotation/model"
)

func LoadFiles(dir string, relFiles []string, diffBase string) (*model.ProjectAnnotation, error) {
	files := make(model.FileAnnotationMapping, len(relFiles))
	for _, relFile := range relFiles {
		newContent, err := os.ReadFile(filepath.Join(dir, relFile))
		if err != nil {
			return nil, err
		}
		newContentStr := string(newContent)
		ok, oldContent, err := git.CatFile(dir, diffBase, relFile)
		if err != nil {
			return nil, err
		}
		var isNew bool
		var contentChanged bool
		var lineChanges *cov_model.LineChanges
		if !ok {
			isNew = true
		} else if newContentStr != oldContent {
			contentChanged = true
			var err error
			lineChanges, err = line.Diff(newContentStr, oldContent)
			if err != nil {
				return nil, fmt.Errorf("diff %s: %w", relFile, err)
			}
		}

		files[model.RelativeFile(relFile)] = &model.FileAnnotation{
			ChangeDetail: &git2.FileDetail{
				IsNew:          isNew,
				ContentChanged: contentChanged,
			},
			LineChanges: lineChanges,
		}
	}

	return &model.ProjectAnnotation{
		Files: files,
		Types: map[model.AnnotationType]bool{
			model.AnnotationType_LineChanges:  true,
			model.AnnotationType_ChangeDetail: true,
		},
	}, nil
}
