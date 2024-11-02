package loadcov

import (
	"fmt"

	"github.com/xhd2015/lines-annotation/load/git"
	"github.com/xhd2015/lines-annotation/model"
)

func LoadGitDiff(dir string, relFiles []string, diffBase string) (*model.ProjectAnnotation, error) {
	if diffBase == "" {
		return nil, fmt.Errorf("requires diffBase")
	}

	return git.LoadFiles(dir, relFiles, diffBase)
}
