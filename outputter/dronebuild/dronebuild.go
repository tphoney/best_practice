package dronebuild

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/tphoney/best_practice/outputter"
	"github.com/tphoney/best_practice/types"
)

const (
	Name     = "drone build"
	FileName = ".drone.yml.new"
)

type (
	OutputFields struct {
		RawYaml string `json:"raw_yaml" yaml:"raw_yaml"`
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
	oc.description = "Creates a Drone build file"
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
	fmt.Println("")
	fmt.Println("Added the following steps to the Drone build file:")
	for _, result := range results {
		fmt.Printf("- %s\n", result.Description)
	}
	fmt.Println("")

	buildOutput := `kind: pipeline
type: docker
name: default

platform:
  os: linux
  arch: amd64

steps:
`
	// add the steps to the build file
	for _, result := range results {
		dbo := result.Spec.(OutputFields)
		buildOutput += fmt.Sprintln(dbo.RawYaml)
	}

	if oc.stdOutput {
		fmt.Printf("Drone build file:\n%s\n", buildOutput)
		fmt.Println(buildOutput)
	}
	if oc.outputToFile {
		fmt.Printf("Created a new Drone Build file '%s'\n", filepath.Join(oc.workingDirectory, FileName))
		err := outputter.WriteToFile(filepath.Join(oc.workingDirectory, FileName), buildOutput)
		if err != nil {
			return err
		}
	}
	fmt.Println("")
	return nil
}
