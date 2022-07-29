package ruby

import (
	"context"
	"fmt"
	"os"
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
	Name        = scanner.RubyScannerName
	BuildCheck  = "Ruby build"
	TestCheck   = "Ruby test"
	LintCheck   = "Ruby lint"
	DroneCheck  = "Ruby Drone build"
	rubyVersion = "latest"
)

func New(opts ...Option) (types.Scanner, error) {
	sc := new(scannerConfig)
	sc.name = Name
	sc.description = "checks for various ruby related best practices"
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
	matches, err := scanner.FindMatchingFiles(sc.workingDirectory, "*.rb", true)
	if err != nil || len(matches) == 0 {
		// nothing to see here, lets leave
		return returnVal, nil
	}
	if sc.runAll || slices.Contains(requestedChecks, TestCheck) {
		_, testpathErr := os.Stat(fmt.Sprintf("%s/spec", sc.workingDirectory))
		if testpathErr == nil {
			droneBuildResult := types.Scanlet{
				Name:           TestCheck,
				ScannerFamily:  Name,
				Description:    "run rspec",
				OutputRenderer: buildmaker.Name,
				Spec: buildmaker.OutputFields{
					Build: buildmaker.Build{
						Name:     "run rspec",
						Image:    fmt.Sprintf("ruby:%s", rubyVersion),
						Commands: []string{"bundle install", "bundle exec rspec spec"},
					},
					CLI:     "bundle exec rspec spec",
					HelpURL: "https://docs.npmjs.com/misc/test",
				},
			}
			returnVal = append(returnVal, droneBuildResult)
		}
	}
	if sc.runAll || slices.Contains(requestedChecks, LintCheck) {
		match, outputResults := sc.lintCheck(rubyVersion)
		if match {
			returnVal = append(returnVal, outputResults...)
		}
	}
	if sc.runAll || slices.Contains(requestedChecks, BuildCheck) {
		match, outputResults := sc.buildCheck(rubyVersion)
		if match {
			returnVal = append(returnVal, outputResults...)
		}
	}
	if sc.runAll || slices.Contains(requestedChecks, DroneCheck) {
		outputResults, err := sc.droneCheck(rubyVersion)
		if err == nil {
			returnVal = append(returnVal, outputResults...)
		}
	}
	return returnVal, nil
}

func (sc *scannerConfig) buildCheck(rubyVersion string) (match bool, outputResults []types.Scanlet) {
	// do we have a rakefile?
	rakefileExist, _ := scanner.FindMatchingFiles(sc.workingDirectory, "Rakefile", true)
	if len(rakefileExist) > 0 {
		droneBuildResult := types.Scanlet{
			Name:           BuildCheck,
			ScannerFamily:  Name,
			Description:    "build using rake",
			OutputRenderer: buildmaker.Name,
			Spec: buildmaker.OutputFields{
				Build: buildmaker.Build{
					Name:     "build with rake",
					Image:    fmt.Sprintf("ruby:%s", rubyVersion),
					Commands: []string{"bundle install", "bundle exec rake"},
				},
				CLI:     "bundle exec rake build",
				HelpURL: "https://bundler.io/man/bundle-exec.1.html",
			},
		}
		outputResults = append(outputResults, droneBuildResult)
		return true, outputResults
	}
	return false, outputResults
}

func (sc *scannerConfig) lintCheck(rubyVersion string) (match bool, outputResults []types.Scanlet) {
	rubocopExist, _ := scanner.FindMatchingFiles(sc.workingDirectory, ".rubocop.yml", false)
	if len(rubocopExist) > 0 {
		lintResult := types.Scanlet{
			Name:           LintCheck,
			ScannerFamily:  Name,
			Description:    "run rubocop",
			OutputRenderer: buildmaker.Name,
			Spec: buildmaker.OutputFields{
				Build: buildmaker.Build{
					Name:     "run rubocop",
					Image:    fmt.Sprintf("ruby:%s", rubyVersion),
					Commands: []string{"bundle install", "bundle exec rubocop"},
				},
				CLI:     "rubocop",
				HelpURL: "https://docs.rubygems.org/rubocop",
			},
		}
		outputResults = append(outputResults, lintResult)
		return true, outputResults
	}
	return false, outputResults
}

func (sc *scannerConfig) droneCheck(nodeVersion string) (outputResults []types.Scanlet, err error) {
	pipelines, err := dronescanner.ReadDroneFile(sc.workingDirectory, dronescanner.DroneFileLocation)
	if err != nil {
		return outputResults, err
	}
	// iterate over the pipelines
	foundRubyBuild := false
	foundRubyLint := false
	foundRubyTest := false
	for i := range pipelines {
		for j := range pipelines[i].Steps {
			commands := pipelines[i].Steps[j].Commands
			for k := range commands {
				if strings.Contains(commands[k], "bundle exec rake") {
					foundRubyBuild = true
				}
				if strings.Contains(commands[k], "rubocop") {
					foundRubyLint = true
				}
				if strings.Contains(commands[k], "bundle exec rspec") {
					foundRubyTest = true
				}
			}
		}
		if foundRubyBuild {
			bestPracticeResult := types.Scanlet{
				Name:           DroneCheck,
				ScannerFamily:  Name,
				Description:    "pipeline '%s' should run ruby build",
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
		if foundRubyLint {
			bestPracticeResult := types.Scanlet{
				Name:           DroneCheck,
				ScannerFamily:  Name,
				Description:    "pipeline '%s' should run rubocop",
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
		if foundRubyTest {
			bestPracticeResult := types.Scanlet{
				Name:           DroneCheck,
				ScannerFamily:  Name,
				Description:    "pipeline '%s' should run npm test",
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
