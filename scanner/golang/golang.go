package golang

import (
	"context"
	"fmt"
	"os"

	"github.com/tphoney/best_practice/types"
)

const (
	goModLocation  = "./go.mod"
	goLintLocation = "./.golangci.yml"
)

var (
// Output information

)

func Scan(ctx context.Context) (returnVal []types.Scanlet, err error) {
	// lets look for a go.mod file in the directory
	goModFile, err := os.Open(goModLocation)
	if err != nil {
		// nothing to see here
		return returnVal, err
	}
	goModFile.Close()
	fmt.Printf("found go.mod file: %s\n", goModFile.Name())
	// find the main.go file
	// find test files
	// check for go linter
	match, lintResult := lintCheck()
	if match {
		returnVal = append(returnVal, lintResult...)
	}
	// build up the output
	droneOutput := types.Scanlet{
		Name:           "go build file",
		HumanReasoning: "run go build",
		Spec: types.DroneBuildOutput{
			RawYaml: `
- name: check go mod file
    image: golang:1
    commands:
      - cp go.mod go.mod.bak
      - go mod tidy
      - diff go.mod go.mod.bak || (echo "go.mod is not up to date" && exit 1)
- name: build
    image: golang:1
    commands:
      - go build
`,
		},
	}
	returnVal = append(returnVal, droneOutput)
	return returnVal, nil
}

func lintCheck() (match bool, outputResults []types.Scanlet) {
	// if golang lint file does exist
	_, err := os.Open(goLintLocation)
	if err != nil {
		droneBuildResult := types.Scanlet{
			Name:           "go lint",
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
			Name:           "go lint",
			HumanReasoning: "go lint teaches you to write better code",
			OutputRender:   types.OutputBestPractice,
			Spec: types.BestPracticeOutput{
				Command: "golangci-lint run",
			},
		}
		outputResults = append(outputResults, bestPracticeResult)
		return true, outputResults
	}
	return false, outputResults
}
