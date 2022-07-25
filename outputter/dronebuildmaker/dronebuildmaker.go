package dronebuildmaker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tphoney/best_practice/outputter"
	"github.com/tphoney/best_practice/types"
)

const (
	Name = outputter.DroneBuildMaker
)

var (
	FileName = ".drone.yml"
)

type (
	OutputFields struct {
		RawYaml string `json:"raw_yaml" yaml:"raw_yaml"`
		Command string `json:"command" yaml:"command"`
		HelpURL string `json:"help_url" yaml:"help_url"`
	}

	outputterConfig struct {
		name             string
		description      string
		stdOutput        bool
		workingDirectory string
		outputToFile     bool
	}
)

func New(opts ...Option) (types.Outputter, error) {
	oc := new(outputterConfig)
	oc.name = Name
	oc.description = "Creates a full Drone build file"
	// apply options
	for _, opt := range opts {
		opt(oc)
	}

	return oc, nil
}

func (oc outputterConfig) Name() string {
	return oc.name
}

func (oc outputterConfig) Description() string {
	return oc.description
}

func (oc outputterConfig) Output(ctx context.Context, scanResults []types.Scanlet) error {
	var results []types.Scanlet
	// iterate over enabled outputs
	for _, output := range scanResults {
		if output.OutputRenderer == Name {
			results = append(results, output)
		}
	}
	// lets render the best practice results to stdout
	if len(results) == 0 {
		return nil
	}
	// lets explain what we added to the drone build
	fmt.Println("Drone Build Generator\n\nAdded the following to the drone build:")
	for _, result := range results {
		fmt.Printf("- %s, %s\n", result.ScannerFamily, result.Description)
	}
	fmt.Println("")

	buildOutput := `kind: pipeline
type: docker
name: default

platform:
  os: linux
  arch: amd64

steps:`
	// add the steps to the build file
	for _, result := range results {
		dbo := result.Spec.(OutputFields)
		buildOutput += dbo.RawYaml
	}

	if oc.stdOutput {
		fmt.Printf("Drone build file:\n%s\n", buildOutput)
		fmt.Println(buildOutput)
	}
	if oc.outputToFile {
		_, err := os.Stat(filepath.Join(oc.workingDirectory, FileName))
		if err == nil {
			// file exists append .new to the file name
			FileName += ".new"
		}
		fmt.Printf("Created a new Drone Build file '%s'\n", filepath.Join(oc.workingDirectory, FileName))
		writeErr := outputter.WriteToFile(filepath.Join(oc.workingDirectory, FileName), buildOutput)
		if writeErr != nil {
			return writeErr
		}
	}
	return nil
}
