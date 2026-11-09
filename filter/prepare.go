package filter

import (
	"github.com/xhd2015/lines-annotation/model"
	"github.com/xhd2015/lines-annotation/model/lang"
)

// PrepareForConsumption filters a request-owned, decoded annotation before
// validating it. Only file/type maps are copied: no second large decoded tree
// is allocated. File contents remain request-owned and may be computed in place.
func PrepareForConsumption(project *model.ProjectAnnotation, language lang.ProfileLanguage, include, exclude []string, requireDiff bool) (*model.ProjectAnnotation, error) {
	if project == nil {
		return nil, model.ValidateForConsumption(nil, requireDiff)
	}
	copy := *project
	copy.Types = make(map[model.AnnotationType]bool, len(project.Types))
	for key, value := range project.Types {
		copy.Types[key] = value
	}
	FilterFiles(&copy, MakeProjectFilter(language, include, exclude))
	if err := model.ValidateForConsumption(&copy, requireDiff); err != nil {
		return nil, err
	}
	return &copy, nil
}
