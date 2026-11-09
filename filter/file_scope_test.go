package filter_test

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/xhd2015/lines-annotation/compute/coverage"
	"github.com/xhd2015/lines-annotation/filter"
	"github.com/xhd2015/lines-annotation/model"
	policy "github.com/xhd2015/lines-annotation/model/filter"
	"github.com/xhd2015/lines-annotation/model/lang"
)

func TestFileUniverseRejectsDependencyDirectoriesButKeepsBusinessCode(t *testing.T) {
	for _, path := range []string{"vendor/a.go", "a/vendor/a.go", "test/api_test/runner/vendor/a.go", `test\runner\vendor\a.go`, "a/node_modules/a.go", "a/b/vendor/c/a.go"} {
		t.Run(path, func(t *testing.T) {
			if filter.NewChecker(filter.MakeProjectFilter("", []string{"**"}, nil)).MatchFile(path) {
				t.Fatal("include resurrected excluded dependency", path)
			}
		})
	}
	for _, path := range []string{"src/a.go", "vendor2/a.go", "my_vendor/a.go", "a/node_modules2/a.go"} {
		if !filter.NewChecker(filter.MakeFilter("")).MatchFile(path) {
			t.Fatal("business path excluded", path)
		}
	}
	if filter.NewChecker(filter.MakeFilter("")).MatchFile("src/a_test.go") {
		t.Fatal("Go test suffix policy changed")
	}
	if !filter.NewChecker(filter.MakeFilter(lang.ProfileLanguage_Js)).MatchFile("src/a.ts") {
		t.Fatal("JS policy changed")
	}
	first := policy.DefaultExcludePathSegments()
	first[0] = "changed"
	if policy.DefaultExcludePathSegments()[0] != "vendor" {
		t.Fatal("shared policy is mutable")
	}
}

func annotation(t *testing.T, raw string) *model.ProjectAnnotation {
	t.Helper()
	var result model.ProjectAnnotation
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		t.Fatal(err)
	}
	return &result
}

func TestOldFinishedAnnotationCanBeConsumedWithoutTestDirectoryWorkaround(t *testing.T) {
	project := annotation(t, `{"types":{"change_detail":true,"line_exec_labels":true},"files":{"src/pay.go":{"changeDetail":{"isNew":true},"lines":{"1":{"blockID":"B","changed":true,"execLabels":{"":true}},"2":{"blockID":"B","changed":true}}},"test/api_test/runner/vendor/a.go":{"lines":{"1":{}}}}}`)
	prepared, err := filter.PrepareForConsumption(project, "", nil, nil, true)
	if err != nil {
		t.Fatal(err)
	}
	if len(prepared.Files) != 1 || len(project.Files) != 2 {
		t.Fatal("filter did not isolate the file map")
	}
	summary, err := coverage.ComputeCoverageSummaryChecked(prepared, &coverage.ComputeOptions{DisableFunc: true})
	if err != nil {
		t.Fatal(err)
	}
	if summary[""] == nil || summary[""].Detail.Total.Total != 2 || summary[""].Detail.Total.Covered != 1 || summary[""].Detail.Incrimental.Total != 2 {
		t.Fatalf("incorrect business coverage: %+v", summary)
	}
	view, err := filter.ReserveForLineViewChecked(prepared, &filter.ReserveOptions{ChangedOnly: true})
	if err != nil || len(view.Files) != 1 {
		t.Fatalf("line view failed: %v", err)
	}
}

func TestBusinessMissingDiffFailsWithoutMutationOrInventedUnchangedData(t *testing.T) {
	project := annotation(t, `{"types":{"change_detail":true},"files":{"src/pay.go":{"lines":{"1":{}},"funcs":{"f":{"name":"Pay"}}}}}`)
	_, err := filter.PrepareForConsumption(project, "", nil, nil, true)
	var integrity *model.IntegrityError
	if !errors.As(err, &integrity) || integrity.File != "src/pay.go" {
		t.Fatalf("wanted file error, got %v", err)
	}
	if _, err = coverage.ComputeCoverageSummaryChecked(project, nil); !errors.As(err, &integrity) {
		t.Fatalf("summary did not return typed error: %v", err)
	}
	if _, err = filter.ReserveForLineViewChecked(project, nil); !errors.As(err, &integrity) {
		t.Fatalf("view did not return typed error: %v", err)
	}
	if project.Files["src/pay.go"].ChangeDetail != nil || len(project.Types) != 1 {
		t.Fatal("invalid data was patched or computed")
	}
}

func TestProjectExclusionPrecedesIntegrityAndLegalNoDiffModeSurvives(t *testing.T) {
	project := annotation(t, `{"types":{"change_detail":true},"files":{"test/helper.go":{},"src/pay.go":{"changeDetail":{}}}}`)
	if _, err := filter.PrepareForConsumption(project, "", nil, []string{"test/"}, true); err != nil {
		t.Fatal(err)
	}
	noDiff := annotation(t, `{"files":{"src/pay.go":{"lines":{"1":{}}}}}`)
	if _, err := filter.PrepareForConsumption(noDiff, "", nil, nil, false); err != nil {
		t.Fatal(err)
	}
	if _, err := filter.PrepareForConsumption(noDiff, "", nil, nil, true); err == nil {
		t.Fatal("missing entire diff marker accepted")
	}
}
