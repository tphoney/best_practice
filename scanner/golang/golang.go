package golang

import (
	"context"
	"os"

	"github.com/tphoney/best_practice/types"
	"golang.org/x/exp/slices"
)

const (
	goModLocation  = "./go.mod"
	goLintLocation = "./.golangci.yml"
)

var (
	scanType = types.ScanType{
		FamilyName: "golang",
		ScanTypes:  []string{"go_mod", "go_lint"},
	}
)

func ReturnScanType() types.ScanType {
	return scanType
}

func Scan(ctx context.Context, scanInput types.ScanInput) (returnVal []types.Scanlet, err error) {
	// lets look for a go.mod file in the directory
	goModFile, err := os.Open(goModLocation)
	if err != nil {
		// nothing to see here
		return returnVal, err
	}
	goModFile.Close()
	// check the mod file
	if scanInput.RunAll || slices.Contains(scanInput.RequestedOutputs, "go_mod") {
		match, outputResults := modCheck()
		if match {
			returnVal = append(returnVal, outputResults...)
		}
	}
	// check for go linter
	if scanInput.RunAll || slices.Contains(scanInput.RequestedOutputs, "go_lint") {
		match, lintResult := lintCheck()
		if match {
			returnVal = append(returnVal, lintResult...)
		}
	}
	// find the main.go file
	// find test files
	return returnVal, nil
}

func modCheck() (match bool, outputResults []types.Scanlet) {
	// if go mod file does exist
	_, err := os.Open(goModLocation)
	if err == nil {
		droneBuildResult := types.Scanlet{
			Name:           "go_mod",
			HumanReasoning: "run go mod",
			Spec: types.DroneBuildOutput{
				RawYaml: `- name: go mod
	image: golang:1
	commands:
	  - go mod tidy
	  - diff go.mod go.mod.bak || (echo "go.mod is not up to date" && exit 1)`,
			},
		}
		outputResults = append(outputResults, droneBuildResult)
		bestPracticeResult := types.Scanlet{
			Name:           "go_mod",
			HumanReasoning: "make sure your go mod file is up to date",
			OutputRender:   types.OutputBestPractice,
			Spec: types.BestPracticeOutput{
				Command: "go mod tidy",
				Url:     "https://golang.org/cmd/go/#hdr-Modules",
			},
		}
		outputResults = append(outputResults, bestPracticeResult)
		return true, outputResults
	}
	return false, outputResults
}

func lintCheck() (match bool, outputResults []types.Scanlet) {
	// if golang lint file does exist
	_, err := os.Open(goLintLocation)
	if err != nil {
		droneBuildResult := types.Scanlet{
			Name:           "go_lint",
			HumanReasoning: "run go lint as part of the build",
			OutputRender:   types.OutputNameDroneBuild,
			Spec: types.DroneBuildOutput{
				RawYaml: `- name: golangci-lint
image: golangci/golangci-lint
commands:
  - golangci-lint run --timeout 500s`,
			},
		}
		outputResults = append(outputResults, droneBuildResult)
		bestPracticeResult := types.Scanlet{
			Name:           "go_lint",
			HumanReasoning: "go lint teaches you to write better code",
			OutputRender:   types.OutputBestPractice,
			Spec: types.BestPracticeOutput{
				Command: "golangci-lint run",
				Url:     "https://golangci-lint.run.googlesource.com/golangci-lint",
			},
		}
		outputResults = append(outputResults, bestPracticeResult)
		return true, outputResults
	}
	return false, outputResults
}
