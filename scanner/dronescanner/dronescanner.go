package dronescanner

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/tphoney/best_practice/outputter/bestpractice"
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
	droneFileLocation       = ".drone.yml"
	Name                    = scanner.DroneScannerName
	DroneCheck              = "drone"
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
	return []string{DroneCheck}
}

func (sc *scannerConfig) Scan(ctx context.Context, requestedOutputs []string) (returnVal []types.Scanlet, err error) {
	// lets look for a go.mod file in the directory
	_, err = os.Stat(filepath.Join(sc.workingDirectory, droneFileLocation))
	if err != nil {
		// nothing to see here, lets leave
		return returnVal, nil
	}
	droneContents, err := readDroneFile(sc.workingDirectory, droneFileLocation)
	if err != nil {
		return returnVal, err
	}
	// count the number of steps per pipeline
	// check for build, test, lint, and deploy
	// check image versions
	// if we use the docker plugin, make sure we use snyk
	// if we have multiple go steps / java steps check for shared volumes
	if sc.runAll || slices.Contains(requestedOutputs, DroneCheck) {
		match, outputResults := droneStepsCheck(droneContents)
		if match {
			returnVal = append(returnVal, outputResults...)
		}
	}
	return returnVal, nil
}

func droneStepsCheck(droneContents []DronePipeline) (match bool, outputResults []types.Scanlet) {
	// iterate over the pipelines
	for i := range droneContents {
		if len(droneContents[i].Steps) > MaximumStepsPerPipeline {
			bestPracticeResult := types.Scanlet{
				Name:           DroneCheck,
				ScannerFamily:  Name,
				Description:    fmt.Sprintf("pipeline '%s' has more than %d steps, split into multiple pipelines", droneContents[i].Name, MaximumStepsPerPipeline),
				OutputRenderer: bestpractice.Name,
				Spec: bestpractice.OutputFields{
					Command: "",
					HelpURL: "https://docs.drone.io/yaml/digitalocean/#the-depends_on-attribute",
				},
			}
			outputResults = append(outputResults, bestPracticeResult)
			match = true
		}
	}
	return match, outputResults
}

func readDroneFile(workingDir, droneFileLocation string) (bla []DronePipeline, err error) {
	file, fileErr := os.Open(filepath.Join(workingDir, droneFileLocation))
	if fileErr != nil {
		return bla, fileErr
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	for {
		myMap := new(DronePipeline)
		yamlErr := decoder.Decode(&myMap)
		if myMap == nil {
			continue
		}
		if errors.Is(yamlErr, io.EOF) {
			break
		}
		// we may want to consider an anoymous struct here
		if yamlErr == nil {
			bla = append(bla, *myMap)
		}
	}
	return bla, err
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
