package annotation_map

import (
	"strings"

	"github.com/xhd2015/lines-annotation/model"
)

func TrimOrAddPrefixDir(project *model.ProjectAnnotation, dir string) {
	dir = strings.TrimSpace(dir)
	remove := false
	if strings.HasPrefix(dir, "-") {
		remove = true
		dir = dir[1:]
	}
	remapDir(project, remove, dir)
}

func TrimPrefixDir(project *model.ProjectAnnotation, dir string) {
	remapDir(project, true, dir)
}
func AddPrefixDir(project *model.ProjectAnnotation, dir string) {
	remapDir(project, false, dir)
}

func remapDir(project *model.ProjectAnnotation, remove bool, dir string) {
	if dir == "" {
		return
	}

	// trim ./
	dir = strings.TrimPrefix(dir, ".")
	dir = strings.TrimPrefix(dir, "/")
	dir = strings.TrimSuffix(dir, "/")
	if dir == "" {
		return
	}
	files := project.Files
	newFiles := make(model.FileAnnotationMapping, len(files))
	project.Files = newFiles
	for file, val := range files {
		var newFile string
		if !remove {
			// add prefix
			newFile = dir + "/" + string(file)
		} else {
			// remove prefix
			newFile = trimPathPrefixForRemap(string(file), dir)
			if newFile == "" {
				continue
			}
		}
		newFiles[model.RelativeFile(newFile)] = val
	}
}

func trimPathPrefixForRemap(file string, prefixNoSlash string) string {
	if !strings.HasPrefix(file, prefixNoSlash) {
		return ""
	}
	return strings.TrimPrefix(file[len(prefixNoSlash):], "/")
}
