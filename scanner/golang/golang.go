package golang

import (
	"context"
	"os"

	"github.com/tphoney/best_practice/outputter/bestpractice"
	"github.com/tphoney/best_practice/outputter/dronebuild"
	"github.com/tphoney/best_practice/types"
	"golang.org/x/exp/slices"
)

type scannerConfig struct {
	name        string
	description string
	checksToRun []string
	runAll      bool
}

const (
	goModLocation  = "./go.mod"
	goLintLocation = "./.golangci.yml"
	Name           = "golang"
	ModCheckName   = "go_mod"
	LintCheckName  = "go_lint"
)

func New(opts ...Option) (types.Scanner, error) {
	sc := new(scannerConfig)
	sc.name = Name
	sc.description = "checks for various go related best practices"
	sc.runAll = true
	// apply options
	for _, opt := range opts {
		opt(sc)
	}

	return sc, nil
}

func (sc *scannerConfig) Name() string {
	return sc.name
}

func (sc *scannerConfig) Description() string {
	return sc.description
}

func (sc *scannerConfig) AvailableChecks() []string {
	return []string{ModCheckName, LintCheckName}
}

func (sc *scannerConfig) Scan(ctx context.Context, requestedOutputs []string) (returnVal []types.Scanlet, err error) {
	// lets look for a go.mod file in the directory
	goModFile, err := os.Open(goModLocation)
	if err != nil {
		// nothing to see here
		return returnVal, err
	}
	goModFile.Close()
	// check the mod file
	if sc.runAll || slices.Contains(requestedOutputs, ModCheckName) {
		match, outputResults := modCheck()
		if match {
			returnVal = append(returnVal, outputResults...)
		}
	}
	// check for go linter
	if sc.runAll || slices.Contains(requestedOutputs, LintCheckName) {
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
			Name:           ModCheckName,
			ScannerFamily:  Name,
			Description:    "run go mod",
			OutputRenderer: dronebuild.Name,
			Spec: dronebuild.DroneBuildOutput{
				RawYaml: `  - name: go mod
  image: golang:1
  commands:
    - go mod tidy
	- diff go.mod go.mod.bak || (echo "go.mod is not up to date" && exit 1)`,
			},
		}
		outputResults = append(outputResults, droneBuildResult)
		bestPracticeResult := types.Scanlet{
			Name:           ModCheckName,
			ScannerFamily:  Name,
			Description:    "make sure your go mod file is up to date",
			OutputRenderer: bestpractice.Name,
			Spec: bestpractice.BestPracticeOutput{
				Command: "go mod tidy",
				Url:     "https://go.dev/ref/mod#go-mod-tidy",
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
			Name:           LintCheckName,
			ScannerFamily:  Name,
			Description:    "run go lint as part of the build",
			OutputRenderer: dronebuild.Name,
			Spec: dronebuild.DroneBuildOutput{
				RawYaml: `  - name: golangci-lint
  image: golangci/golangci-lint
  commands:
    - golangci-lint run --timeout 500s`,
			},
		}
		outputResults = append(outputResults, droneBuildResult)
		bestPracticeResult := types.Scanlet{
			Name:           LintCheckName,
			ScannerFamily:  Name,
			Description:    "go lint teaches you to write better code",
			OutputRenderer: bestpractice.Name,
			Spec: bestpractice.BestPracticeOutput{
				Command: "golangci-lint run",
				Url:     "https://golangci-lint.run.googlesource.com/golangci-lint",
			},
		}
		outputResults = append(outputResults, bestPracticeResult)
		return true, outputResults
	}
	return false, outputResults
}
