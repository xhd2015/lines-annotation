package load

import (
	"github.com/xhd2015/lines-annotation/model"
)

// see ast.LoadFuncs
func FuncInfoMappingToAnnotation(mapping map[model.RelativeFile][]*model.FuncAnnotation) *model.ProjectAnnotation {
	if len(mapping) == 0 {
		return nil
	}
	files := make(model.FileAnnotationMapping, len(mapping))
	for file, funcInfos := range mapping {
		funcs := make(model.FuncAnnotationMapping, len(funcInfos))
		for _, funcInfo := range funcInfos {
			if funcInfo.Block == nil {
				continue
			}
			blockID := funcInfo.Block.ID()
			funcs[blockID] = funcInfo
		}
		files[file] = &model.FileAnnotation{
			Funcs: funcs,
		}
	}
	return &model.ProjectAnnotation{
		Files: files,
		Types: map[model.AnnotationType]bool{
			model.AnnotationType_FileFuncs: true,
		},
	}
}
