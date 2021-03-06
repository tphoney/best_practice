package javascript

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
	packageLocation = "package.json"
	Name            = scanner.JavascriptScannerName
	BuildCheck      = "Javascript build"
	TestCheck       = "Javascript test"
	LintCheck       = "Javascript lint"
	DroneCheck      = "Javascript Drone build"
	nodeVersion     = "18"
)

func New(opts ...Option) (types.Scanner, error) {
	sc := new(scannerConfig)
	sc.name = Name
	sc.description = "checks for various javascript related best practices"
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
	return []string{BuildCheck, TestCheck, LintCheck, DroneCheck}
}

func (sc *scannerConfig) Scan(ctx context.Context, requestedChecks []string) (returnVal []types.Scanlet, err error) {
	// lets look for a package file in the directory
	_, err = os.Stat(filepath.Join(sc.workingDirectory, packageLocation))
	if err != nil {
		// nothing to see here, lets leave
		return returnVal, nil
	}
	var scriptMap map[string]interface{}
	packageStruct, err := scanner.ReadJSONFile(filepath.Join(sc.workingDirectory, packageLocation))
	if err == nil {
		// look for declared scripts
		if packageStruct["scripts"] != nil {
			scriptMap = packageStruct["scripts"].(map[string]interface{})
		}
	} else {
		return returnVal, err
	}
	if sc.runAll || slices.Contains(requestedChecks, TestCheck) {
		match, outputResults := sc.testCheck(scriptMap)
		if match {
			returnVal = append(returnVal, outputResults...)
		}
	}
	if sc.runAll || slices.Contains(requestedChecks, LintCheck) {
		match, outputResults := sc.lintCheck(scriptMap)
		if match {
			returnVal = append(returnVal, outputResults...)
		}
	}
	if sc.runAll || slices.Contains(requestedChecks, BuildCheck) {
		match, outputResults := sc.buildCheck(scriptMap)
		if match {
			returnVal = append(returnVal, outputResults...)
		}
	}
	if sc.runAll || slices.Contains(requestedChecks, DroneCheck) {
		outputResults, err := sc.droneCheck()
		if err == nil {
			returnVal = append(returnVal, outputResults...)
		}
	}
	return returnVal, nil
}

func (sc *scannerConfig) buildCheck(scriptMap map[string]interface{}) (match bool, outputResults []types.Scanlet) {
	if scriptMap["build"] != "" {
		buildResult := types.Scanlet{
			Name:           BuildCheck,
			ScannerFamily:  Name,
			Description:    "run npm build",
			OutputRenderer: buildmaker.Name,
			Spec: buildmaker.OutputFields{
				Build: buildmaker.Build{
					Name:     "run npm build",
					Image:    fmt.Sprintf("node:%s-alpine", nodeVersion),
					Commands: []string{"npm run build"},
				},
				CLI:     "npm run build",
				HelpURL: "https://docs.npmjs.com/misc/build",
			},
		}
		outputResults = append(outputResults, buildResult)
		return true, outputResults
	}
	return false, outputResults
}

func (sc *scannerConfig) lintCheck(scriptMap map[string]interface{}) (match bool, outputResults []types.Scanlet) {
	if scriptMap["lint"] != "" {
		lintResult := types.Scanlet{
			Name:           LintCheck,
			ScannerFamily:  Name,
			Description:    "run npm lint",
			OutputRenderer: buildmaker.Name,
			Spec: buildmaker.OutputFields{
				Build: buildmaker.Build{
					Name:     "run npm lint",
					Image:    fmt.Sprintf("node:%s-alpine", nodeVersion),
					Commands: []string{"npm run lint"},
				},
				CLI:     "npm run lint",
				HelpURL: "https://docs.npmjs.com/misc/lint",
			},
		}
		outputResults = append(outputResults, lintResult)
		return true, outputResults
	}
	return false, outputResults
}

func (sc *scannerConfig) testCheck(scriptMap map[string]interface{}) (match bool, outputResults []types.Scanlet) {
	if scriptMap["test"] != "" {
		testResult := types.Scanlet{
			Name:           TestCheck,
			ScannerFamily:  Name,
			Description:    "run npm test",
			OutputRenderer: buildmaker.Name,
			Spec: buildmaker.OutputFields{
				Build: buildmaker.Build{
					Name:     "run npm test",
					Image:    fmt.Sprintf("node:%s-alpine", nodeVersion),
					Commands: []string{"npm run test"},
				},
				CLI:     "npm run test",
				HelpURL: "https://docs.npmjs.com/misc/test",
			},
		}
		outputResults = append(outputResults, testResult)
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
	foundNPMBuild := false
	foundNPMLint := false
	foundNPMTest := false
	for i := range pipelines {
		for j := range pipelines[i].Steps {
			commands := pipelines[i].Steps[j].Commands
			for k := range commands {
				if strings.Contains(commands[k], "npm run build") {
					foundNPMBuild = true
				}
				if strings.Contains(commands[k], "npm run lint") {
					foundNPMLint = true
				}
				if strings.Contains(commands[k], "npm run test") {
					foundNPMTest = true
				}
			}
		}
		if foundNPMBuild {
			bestPracticeResult := types.Scanlet{
				Name:           DroneCheck,
				ScannerFamily:  Name,
				Description:    fmt.Sprintf("pipeline '%s' should run npm build", pipelines[i].Name),
				OutputRenderer: outputter.DroneBuildAnalysis,
				Spec: dronebuildanalysis.OutputFields{
					HelpURL: "https://docs.npmjs.com/misc/build",
					RawYaml: fmt.Sprintf(`
    - name: run npm build
      image: node:%s-alpine
      commands:
        - npm run build`, nodeVersion),
				},
			}
			outputResults = append(outputResults, bestPracticeResult)
		}
		if foundNPMLint {
			bestPracticeResult := types.Scanlet{
				Name:           DroneCheck,
				ScannerFamily:  Name,
				Description:    fmt.Sprintf("pipeline '%s' should run npm lint", pipelines[i].Name),
				OutputRenderer: outputter.DroneBuildAnalysis,
				Spec: dronebuildanalysis.OutputFields{
					HelpURL: "https://docs.npmjs.com/misc/lint",
					RawYaml: fmt.Sprintf(`
  - name: run npm build
    image: node:%s-alpine
    commands:
    - npm run lint`, nodeVersion),
				},
			}
			outputResults = append(outputResults, bestPracticeResult)
		}
		if foundNPMTest {
			bestPracticeResult := types.Scanlet{
				Name:           DroneCheck,
				ScannerFamily:  Name,
				Description:    fmt.Sprintf("pipeline '%s' should run npm test", pipelines[i].Name),
				OutputRenderer: outputter.DroneBuildAnalysis,
				Spec: dronebuildanalysis.OutputFields{
					HelpURL: "https://docs.npmjs.com/misc/test",
					RawYaml: fmt.Sprintf(`
  - name: run npm build
    image: node:%s-alpine
    commands:
      - npm run test`, nodeVersion),
				},
			}
			outputResults = append(outputResults, bestPracticeResult)
		}
	}
	return outputResults, err
}
