package dronebuild

import (
	"context"
	"fmt"

	"github.com/tphoney/best_practice/outputter"
	"github.com/tphoney/best_practice/types"
)

const (
	Name = "drone build"
)

type (
	DroneBuildOutput struct {
		RawYaml string `json:"raw_yaml" yaml:"raw_yaml"`
	}

	outputterConfig struct {
		name         string
		description  string
		stdOutput    bool
		outputToFile string
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
	if len(results) <= 0 {
		return nil
	}
	// lets explain what we added to the drone build
	fmt.Println("")
	fmt.Printf("Created a new Drone Build file '%s' with:\n", oc.outputToFile)
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
		dbo := result.Spec.(DroneBuildOutput)
		buildOutput += fmt.Sprintln(dbo.RawYaml)
	}

	if oc.stdOutput {
		fmt.Println(buildOutput)
	}
	if oc.outputToFile != "" {
		err := outputter.WriteToFile(oc.outputToFile, buildOutput)
		if err != nil {
			return err
		}
	}
	return nil
}
