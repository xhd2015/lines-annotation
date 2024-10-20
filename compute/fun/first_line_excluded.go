package fun

import "github.com/xhd2015/lines-annotation/model"

// depends on Remark.Excluded
func FirstLineExcludedByRemark(project *model.ProjectAnnotation) {
	if project.Has(model.AnnotationType_FirstLineExcluded) {
		return
	}
	for _, fileData := range project.Files {
		for _, fn := range fileData.Funcs {
			lineData := fileData.Lines[model.LineNum(fn.Block.StartLine)]
			if lineData == nil || lineData.Remark == nil {
				continue
			}
			if lineData.Remark.Excluded {
				fn.FirstLineExcluded = true
			}
		}
	}
	project.Set(model.AnnotationType_FirstLineExcluded)
}
