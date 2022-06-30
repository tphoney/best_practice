package golang

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tphoney/best_practice/outputter/bestpractice"
	"github.com/tphoney/best_practice/outputter/dronebuild"
	"github.com/tphoney/best_practice/scanner"
	"github.com/tphoney/best_practice/types"
	"golang.org/x/exp/slices"
)

type scannerConfig struct {
	name             string
	description      string
	workingDirectory string
	checksToRun      []string
	runAll           bool
}

const (
	goModLocation  = "go.mod"
	goLintLocation = ".golangci.yml"
	Name           = "golang"
	ModCheck       = "golang_mod"
	LintCheck      = "golang_lint"
	MainCheck      = "golang_main"
	UnitTestCheck  = "golang_unit_test"
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
	return []string{ModCheck, LintCheck, MainCheck, UnitTestCheck}
}

func (sc *scannerConfig) Scan(ctx context.Context, requestedOutputs []string) (returnVal []types.Scanlet, err error) {
	// lets look for a go.mod file in the directory
	_, err = os.Stat(filepath.Join(sc.workingDirectory, goModLocation))
	if err != nil {
		// nothing to see here, lets leave
		return returnVal, nil
	}
	// check the mod file
	if sc.runAll || slices.Contains(requestedOutputs, ModCheck) {
		match, outputResults := sc.modCheck()
		if match {
			returnVal = append(returnVal, outputResults...)
		}
	}
	// check for go linter
	if sc.runAll || slices.Contains(requestedOutputs, LintCheck) {
		match, lintResult := sc.lintCheck()
		if match {
			returnVal = append(returnVal, lintResult...)
		}
	}
	// find test files
	if sc.runAll || slices.Contains(requestedOutputs, UnitTestCheck) {
		match, testResult := sc.unitTestCheck()
		if match {
			returnVal = append(returnVal, testResult...)
		}
	}
	// find the main.go file
	if sc.runAll || slices.Contains(requestedOutputs, MainCheck) {
		match, mainResult := sc.mainCheck()
		if match {
			returnVal = append(returnVal, mainResult...)
		}
	}
	return returnVal, nil
}

func (sc *scannerConfig) modCheck() (match bool, outputResults []types.Scanlet) {
	// if go mod file does exist
	_, err := os.Stat(filepath.Join(sc.workingDirectory, goModLocation))
	if err == nil {
		droneBuildResult := types.Scanlet{
			Name:           ModCheck,
			ScannerFamily:  Name,
			Description:    "run go mod",
			OutputRenderer: dronebuild.Name,
			Spec: dronebuild.OutputFields{
				RawYaml: `  - name: go mod
    image: golang:1
    commands:
      - go mod tidy
      - diff go.mod go.mod.bak || (echo "go.mod is not up to date" && exit 1)`,
			},
		}
		outputResults = append(outputResults, droneBuildResult)
		bestPracticeResult := types.Scanlet{
			Name:           ModCheck,
			ScannerFamily:  Name,
			Description:    "make sure your go mod file is up to date",
			OutputRenderer: bestpractice.Name,
			Spec: bestpractice.OutputFields{
				Command: "go mod tidy",
				HelpURL: "https://go.dev/ref/mod#go-mod-tidy",
			},
		}
		outputResults = append(outputResults, bestPracticeResult)
		return true, outputResults
	}
	return false, outputResults
}

func (sc *scannerConfig) lintCheck() (match bool, outputResults []types.Scanlet) {
	// if golang lint file does exist
	_, err := os.Stat(filepath.Join(sc.workingDirectory, goLintLocation))
	if err != nil {
		droneBuildResult := types.Scanlet{
			Name:           LintCheck,
			ScannerFamily:  Name,
			Description:    "run go lint as part of the build",
			OutputRenderer: dronebuild.Name,
			Spec: dronebuild.OutputFields{
				RawYaml: `  - name: golangci-lint
    image: golangci/golangci-lint
    commands:
      - golangci-lint run --timeout 500s`,
			},
		}
		outputResults = append(outputResults, droneBuildResult)
		bestPracticeResult := types.Scanlet{
			Name:           LintCheck,
			ScannerFamily:  Name,
			Description:    "go lint teaches you to write better code",
			OutputRenderer: bestpractice.Name,
			Spec: bestpractice.OutputFields{
				Command: "golangci-lint run",
				HelpURL: "https://golangci-lint.run.googlesource.com/golangci-lint",
			},
		}
		outputResults = append(outputResults, bestPracticeResult)
		return true, outputResults
	}
	return false, outputResults
}

func (sc *scannerConfig) mainCheck() (match bool, outputResults []types.Scanlet) {
	matches, err := scanner.FindMatchingFiles(sc.workingDirectory, "main.go")
	if err == nil && len(matches) > 0 {
		// we use the first one found
		mainLocation := strings.TrimPrefix(matches[0], sc.workingDirectory)
		mainLocation = strings.TrimSuffix(mainLocation, "main.go")
		if len(mainLocation) == 1 {
			// dont do anything if main is in the working directory
			mainLocation = ""
		} else {
			mainLocation = "." + mainLocation
		}

		droneBuildResult := types.Scanlet{
			Name:           LintCheck,
			ScannerFamily:  Name,
			Description:    "run go build",
			OutputRenderer: dronebuild.Name,
			Spec: dronebuild.OutputFields{
				RawYaml: fmt.Sprintf(
					`  - name: build
    image: golang:1
    commands:
      - go build %s`, mainLocation),
			},
		}
		outputResults = append(outputResults, droneBuildResult)
		return true, outputResults
	}
	return false, outputResults
}

func (sc *scannerConfig) unitTestCheck() (match bool, outputResults []types.Scanlet) {
	matches, err := scanner.FindMatchingFiles(sc.workingDirectory, "*_test.go")
	if err == nil && len(matches) > 0 {
		droneBuildResult := types.Scanlet{
			Name:           UnitTestCheck,
			ScannerFamily:  Name,
			Description:    "run go unit tests",
			OutputRenderer: dronebuild.Name,
			Spec: dronebuild.OutputFields{
				RawYaml: `  - name: golang unit tests
    image: golang:1
    commands:
      - go test ./...`,
			},
		}
		outputResults = append(outputResults, droneBuildResult)
		return true, outputResults
	}
	return false, outputResults
}
