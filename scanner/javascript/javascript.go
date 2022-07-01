package javascript

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

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
	packageLocation = "package.json"
	Name            = scanner.JavascriptScannerName
	BuildCheck      = "build_check"
	TestCheck       = "test_check"
	LintCheck       = "lint_check"
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
	return []string{BuildCheck}
}

func (sc *scannerConfig) Scan(ctx context.Context, requestedChecks []string) (returnVal []types.Scanlet, err error) {
	// lets look for a package file in the directory
	_, err = os.Stat(filepath.Join(sc.workingDirectory, packageLocation))
	if err != nil {
		// nothing to see here, lets leave
		return returnVal, nil
	}
	var scriptMap map[string]interface{}
	var dependencyMap map[string]interface{}
	var reactVersion string
	packageStruct, err := scanner.ReadJSONFile(filepath.Join(sc.workingDirectory, packageLocation))
	if err == nil {
		// look for declared scripts
		if packageStruct["scripts"] != nil {
			scriptMap = packageStruct["scripts"].(map[string]interface{})
		}
		if packageStruct["dependencies"] != nil {
			dependencyMap = packageStruct["dependencies"].(map[string]interface{})
			rawReactVersion := dependencyMap["react"].(string)
			v, versionErr := scanner.ReturnVersionObject(rawReactVersion)
			if versionErr != nil {
				fmt.Printf("error parsing react version: %s\n", versionErr.Error())
			}
			reactVersion = fmt.Sprint(v.Major())
		}
	} else {
		return returnVal, err
	}
	if sc.runAll || slices.Contains(requestedChecks, TestCheck) {
		match, outputResults := sc.testCheck(scriptMap, reactVersion)
		if match {
			returnVal = append(returnVal, outputResults...)
		}
	}
	if sc.runAll || slices.Contains(requestedChecks, LintCheck) {
		match, outputResults := sc.lintCheck(scriptMap, reactVersion)
		if match {
			returnVal = append(returnVal, outputResults...)
		}
	}
	if sc.runAll || slices.Contains(requestedChecks, BuildCheck) {
		match, outputResults := sc.buildCheck(scriptMap, reactVersion)
		if match {
			returnVal = append(returnVal, outputResults...)
		}
	}
	return returnVal, nil
}

func (sc *scannerConfig) buildCheck(scriptMap map[string]interface{}, reactVersion string) (match bool, outputResults []types.Scanlet) {
	if scriptMap["build"] != "" {
		droneBuildResult := types.Scanlet{
			Name:           BuildCheck,
			ScannerFamily:  Name,
			Description:    "run npm build",
			OutputRenderer: dronebuild.Name,
			Spec: dronebuild.OutputFields{
				RawYaml: fmt.Sprintf(`  - name: run npm build
    image: node:%s-alpine
    commands:
      - npm run build`, reactVersion)},
		}
		outputResults = append(outputResults, droneBuildResult)
		return true, outputResults
	}
	return false, outputResults
}

func (sc *scannerConfig) lintCheck(scriptMap map[string]interface{}, reactVersion string) (match bool, outputResults []types.Scanlet) {
	if scriptMap["lint"] != "" {
		droneBuildResult := types.Scanlet{
			Name:           BuildCheck,
			ScannerFamily:  Name,
			Description:    "run npm lint",
			OutputRenderer: dronebuild.Name,
			Spec: dronebuild.OutputFields{
				RawYaml: fmt.Sprintf(`  - name: run npm lint
    image: node:%s-alpine
    commands:
      - npm run lint`, reactVersion)},
		}
		outputResults = append(outputResults, droneBuildResult)
		return true, outputResults
	}
	return false, outputResults
}

func (sc *scannerConfig) testCheck(scriptMap map[string]interface{}, reactVersion string) (match bool, outputResults []types.Scanlet) {
	if scriptMap["test"] != "" {
		droneBuildResult := types.Scanlet{
			Name:           BuildCheck,
			ScannerFamily:  Name,
			Description:    "run npm test",
			OutputRenderer: dronebuild.Name,
			Spec: dronebuild.OutputFields{
				RawYaml: fmt.Sprintf(`  - name: run npm test
    image: node:%s-alpine
    commands:
      - npm run test`, reactVersion)},
		}
		outputResults = append(outputResults, droneBuildResult)
		return true, outputResults
	}

	return false, outputResults
}
