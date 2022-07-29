package golang

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/tphoney/best_practice/outputter"

	"github.com/tphoney/best_practice/outputter/buildmaker"
	"github.com/tphoney/best_practice/outputter/dronebuildanalysis"
	"github.com/tphoney/best_practice/scanner"
	"github.com/tphoney/best_practice/scanner/dronescanner"
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
	Name           = scanner.GolangScannerName
	ModCheck       = "Golang mod"
	LintCheck      = "Golang lint"
	MainCheck      = "Golang main"
	testCheck      = "Golang test"
	DroneCheck     = "Golang Drone build"
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
	return []string{ModCheck, LintCheck, MainCheck, testCheck, DroneCheck}
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
	if sc.runAll || slices.Contains(requestedOutputs, testCheck) {
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
	if sc.runAll || slices.Contains(requestedOutputs, DroneCheck) {
		droneResult, err := sc.droneCheck()
		if err == nil {
			returnVal = append(returnVal, droneResult...)
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
			OutputRenderer: buildmaker.Name,
			Spec: buildmaker.OutputFields{
				Build: buildmaker.Build{
					Name:     "go mod",
					Image:    "golang:1",
					Commands: []string{"go mod tidy", `diff go.mod go.mod.bak || (echo "go.mod is not up to date" && exit 1)`},
				},
				CLI:     "go mod tidy",
				HelpURL: "https://go.dev/ref/mod#go-mod-tidy",
			},
		}
		outputResults = append(outputResults, droneBuildResult)
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
			OutputRenderer: buildmaker.Name,
			Spec: buildmaker.OutputFields{
				Build: buildmaker.Build{
					Name:     "go lint",
					Image:    "golangci/golangci-lint",
					Commands: []string{"golangci-lint run --timeout 500s"},
				},
				CLI:     "golangci-lint run",
				HelpURL: "https://golangci-lint.run.googlesource.com/golangci-lint",
			},
		}
		outputResults = append(outputResults, droneBuildResult)
		return true, outputResults
	}
	return false, outputResults
}

func (sc *scannerConfig) mainCheck() (match bool, outputResults []types.Scanlet) {
	matches, err := scanner.FindMatchingFiles(sc.workingDirectory, "main.go", true)
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
			OutputRenderer: buildmaker.Name,
			Spec: buildmaker.OutputFields{
				Build: buildmaker.Build{
					Name:     "go build",
					Image:    "golang:1",
					Commands: []string{fmt.Sprintf("go build %s", mainLocation)},
				},
				CLI:     fmt.Sprintf("go build %s", mainLocation),
				HelpURL: "https://pkg.go.dev/cmd/go#hdr-Build_constraints",
			},
		}
		outputResults = append(outputResults, droneBuildResult)
		return true, outputResults
	}
	return false, outputResults
}

func (sc *scannerConfig) unitTestCheck() (match bool, outputResults []types.Scanlet) {
	matches, err := scanner.FindMatchingFiles(sc.workingDirectory, "*_test.go", true)
	if err == nil && len(matches) > 0 {
		droneBuildResult := types.Scanlet{
			Name:           testCheck,
			ScannerFamily:  Name,
			Description:    "run go unit tests",
			OutputRenderer: buildmaker.Name,
			Spec: buildmaker.OutputFields{
				Build: buildmaker.Build{
					Name:     "go unit tests",
					Image:    "golang:1",
					Commands: []string{"go test ./..."},
				},
				CLI:     "go test ./...",
				HelpURL: "https://golang.org/cmd/go/#hdr-Testing_tools",
			},
		}
		outputResults = append(outputResults, droneBuildResult)
		return true, outputResults
	}
	return false, outputResults
}

func (sc *scannerConfig) droneCheck() (outputResults []types.Scanlet, err error) {
	pipelines, err := dronescanner.ReadDroneFile(sc.workingDirectory, dronescanner.DroneFileLocation)
	if err != nil {
		return outputResults, err
	}
	// iterate over the pipelines
	for i := range pipelines {
		foundGoLint := false
		foundGoUnit := false
		foundGoBuild := false
		foundGoMod := false
		for j := range pipelines[i].Steps {
			commands := pipelines[i].Steps[j].Commands
			for k := range commands {
				if strings.Contains(commands[k], "go build") {
					foundGoBuild = true
				}
				if strings.Contains(commands[k], "golangci-lint run") {
					foundGoLint = true
				}
				if strings.Contains(commands[k], "go mod tidy") {
					foundGoMod = true
				}
				if strings.Contains(commands[k], "go test") {
					foundGoUnit = true
				}
			}
		}
		if !foundGoMod && foundGoBuild {
			bestPracticeResult := types.Scanlet{
				Name:           DroneCheck,
				ScannerFamily:  Name,
				Description:    fmt.Sprintf("pipeline '%s' should check mod file is up to date", pipelines[i].Name),
				OutputRenderer: outputter.DroneBuildAnalysis,
				Spec: dronebuildanalysis.OutputFields{
					HelpURL: "https://go.dev/ref/mod#go-mod-tidy",
					Command: "go mod tidy",
					RawYaml: `
  - name: go mod tidy
    image: golang:1
    commands:
      - go mod tidy
      - diff go.mod go.mod.bak || (echo "go.mod is not up to date" && exit 1)`,
				},
			}
			outputResults = append(outputResults, bestPracticeResult)
		}
		if !foundGoLint && foundGoBuild {
			bestPracticeResult := types.Scanlet{
				Name:           DroneCheck,
				ScannerFamily:  Name,
				Description:    fmt.Sprintf("pipeline '%s' should check go lint", pipelines[i].Name),
				OutputRenderer: outputter.DroneBuildAnalysis,
				Spec: dronebuildanalysis.OutputFields{
					HelpURL: "https://golangci-lint.run.googlesource.com/golangci-lint",
					Command: "golangci-lint run",
					RawYaml: `
  - name: golangci-lint
    image: golangci/golangci-lint
    commands:
      - golangci-lint run --timeout 500s`,
				},
			}
			outputResults = append(outputResults, bestPracticeResult)
		}
		if !foundGoUnit && foundGoBuild {
			bestPracticeResult := types.Scanlet{
				Name:           DroneCheck,
				ScannerFamily:  Name,
				Description:    fmt.Sprintf("pipeline '%s' should check go unit tests", pipelines[i].Name),
				OutputRenderer: outputter.DroneBuildAnalysis,
				Spec: dronebuildanalysis.OutputFields{
					HelpURL: "https://golang.org/cmd/go/#hdr-Testing_tools",
					Command: "go test ./...",
					RawYaml: `
  - name: go unit tests
    image: golang:1
    commands:
      - go test ./...`,
				},
			}
			outputResults = append(outputResults, bestPracticeResult)
		}
	}
	return outputResults, err
}
