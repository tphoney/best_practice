package dronescanner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/tphoney/best_practice/outputter"
	"github.com/tphoney/best_practice/outputter/buildanalysis"
	"github.com/tphoney/best_practice/scanner"
	"github.com/tphoney/best_practice/types"
	"golang.org/x/exp/slices"
	"gopkg.in/yaml.v3"
)

type scannerConfig struct {
	name             string
	description      string
	workingDirectory string
	checksToRun      []string
	runAll           bool
}

const (
	DroneFileLocation       = ".drone.yml"
	Name                    = scanner.DroneScannerName
	StepsCheck              = "Drone max steps"
	VolumeCachingCheck      = "Drone volume caching"
	MaximumStepsPerPipeline = 6
)

func New(opts ...Option) (types.Scanner, error) {
	sc := new(scannerConfig)
	sc.name = Name
	sc.description = "checks for various drone related best practices"
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
	return []string{StepsCheck, VolumeCachingCheck}
}

func (sc *scannerConfig) Scan(ctx context.Context, requestedOutputs []string) (returnVal []types.Scanlet, err error) {
	// lets look for a go.mod file in the directory
	_, err = os.Stat(filepath.Join(sc.workingDirectory, DroneFileLocation))
	if err != nil {
		// nothing to see here, lets leave
		return returnVal, nil
	}
	pipelines, err := ReadDroneFile(sc.workingDirectory, DroneFileLocation)
	if err != nil {
		return returnVal, err
	}
	// count the number of steps per pipeline
	if sc.runAll || slices.Contains(requestedOutputs, StepsCheck) {
		match, outputResults := droneStepsCheck(pipelines)
		if match {
			returnVal = append(returnVal, outputResults...)
		}
	}
	// if we have multiple go steps / java steps check for shared volumes
	if sc.runAll || slices.Contains(requestedOutputs, VolumeCachingCheck) {
		match, outputResults := droneVolumesCheck(pipelines)
		if match {
			returnVal = append(returnVal, outputResults...)
		}
	}
	return returnVal, nil
}

func droneStepsCheck(pipelines []DronePipeline) (match bool, outputResults []types.Scanlet) {
	// iterate over the pipelines
	for i := range pipelines {
		if len(pipelines[i].Steps) > MaximumStepsPerPipeline {
			bestPracticeResult := types.Scanlet{
				Name:           StepsCheck,
				ScannerFamily:  Name,
				Description:    fmt.Sprintf("pipeline '%s' has more than %d steps, split into multiple pipelines", pipelines[i].Name, MaximumStepsPerPipeline),
				OutputRenderer: outputter.DroneBuildAnalysis,
				Spec: buildanalysis.OutputFields{
					HelpURL: "https://docs.drone.io/yaml/docker/#the-depends_on-attribute",
				},
			}
			outputResults = append(outputResults, bestPracticeResult)
			match = true
		}
	}
	return match, outputResults
}

func droneVolumesCheck(pipelines []DronePipeline) (match bool, outputResults []types.Scanlet) {
	// iterate over the pipelines
	for i := range pipelines {
		numberOfGOSteps := 0
		for j := range pipelines[i].Steps {
			commands := pipelines[i].Steps[j].Commands
			for k := range commands {
				if strings.Contains(commands[k], "go ") {
					numberOfGOSteps++
					// dont count multiple go commands in the same step
					break
				}
			}
		}
		if numberOfGOSteps > 1 {
			bestPracticeResult := types.Scanlet{
				Name:           VolumeCachingCheck,
				ScannerFamily:  Name,
				Description:    fmt.Sprintf("pipeline '%s' has %d golang steps, use a volume", pipelines[i].Name, numberOfGOSteps),
				OutputRenderer: outputter.DroneBuildAnalysis,
				Spec: buildanalysis.OutputFields{
					HelpURL: "https://docs.drone.io/pipeline/docker/syntax/volumes/temporary/",
				},
			}
			outputResults = append(outputResults, bestPracticeResult)
			match = true
		}
	}
	return match, outputResults
}

func ReadDroneFile(workingDir, droneFileLocation string) (pipelines []DronePipeline, err error) {
	file, fileErr := os.Open(filepath.Join(workingDir, droneFileLocation))
	if fileErr != nil {
		return pipelines, fileErr
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	myMap := new(DronePipeline)
	for {
		yamlErr := decoder.Decode(&myMap)
		if myMap == nil {
			continue
		}
		if errors.Is(yamlErr, io.EOF) {
			break
		}
		if yamlErr != nil {
			return pipelines, fmt.Errorf("error reading %s '%s'", filepath.Join(workingDir, droneFileLocation), yamlErr)
		}
		pipelines = append(pipelines, *myMap)
	}
	return pipelines, err
}

// DronePipeline
type DronePipeline struct {
	Kind  string  `yaml:"kind"`
	Type  string  `yaml:"type"`
	Name  string  `yaml:"name"`
	Steps []Steps `yaml:"steps"`
}

// Steps
type Steps struct {
	Name      string   `yaml:"name"`
	Image     string   `yaml:"image"`
	Commands  []string `yaml:"commands"`
	Detach    bool     `yaml:"detach"`
	DependsOn []string `yaml:"depends_on"`
}
