package annotation_map

import (
	"github.com/xhd2015/lines-annotation/compute"
	"github.com/xhd2015/lines-annotation/model"
)

// ShiftLineRemark maps `baseAnnotation` into new annotation using changes mapping
func ShiftLineRemark(baseAnnotation *model.ProjectAnnotation, changes *model.ProjectAnnotation) (*model.ProjectAnnotation, error) {
	if err := changes.ShouldHave("ShiftLineRemark", model.AnnotationType_LineChanges, model.AnnotationType_ChangeDetail); err != nil {
		return nil, err
	}
	// block to line
	newProject := &model.ProjectAnnotation{}
	if baseAnnotation.Has(model.AnnotationType_LineRemark) {
		compute.LineChangesMergeLineRemark(newProject, changes, baseAnnotation)
	}

	return newProject, nil
}
