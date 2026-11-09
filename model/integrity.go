package model

import "fmt"

// IntegrityError identifies corrupt/incomplete annotation data without turning
// a missing diff into an unchanged file (which would undercount new code).
type IntegrityError struct {
	File   RelativeFile
	Reason string
}

func (e *IntegrityError) Error() string {
	return fmt.Sprintf("annotation integrity: %s: %s", e.Reason, e.File)
}

// ValidateForConsumption runs after file filtering, before any derived labels,
// statistics or line-view projection. requireDiff is false for legal modes
// without a baseline; a project declaring change_detail is always checked.
func ValidateForConsumption(project *ProjectAnnotation, requireDiff bool) error {
	if project == nil {
		return &IntegrityError{Reason: "missing project"}
	}
	if requireDiff && !project.Has(AnnotationType_ChangeDetail) {
		return &IntegrityError{Reason: "missing change_detail annotation"}
	}
	for file, data := range project.Files {
		if data == nil {
			return &IntegrityError{File: file, Reason: "nil file annotation"}
		}
		if (requireDiff || project.Has(AnnotationType_ChangeDetail)) && data.ChangeDetail == nil {
			return &IntegrityError{File: file, Reason: "missing diff"}
		}
		for _, line := range data.Lines {
			if line == nil {
				return &IntegrityError{File: file, Reason: "nil line annotation"}
			}
		}
		for _, fn := range data.Funcs {
			if fn == nil {
				return &IntegrityError{File: file, Reason: "nil function annotation"}
			}
		}
		for _, block := range data.Blocks {
			if block == nil {
				return &IntegrityError{File: file, Reason: "nil block annotation"}
			}
		}
	}
	return nil
}
