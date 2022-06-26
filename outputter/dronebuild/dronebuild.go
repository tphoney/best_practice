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

type outputterConfig struct {
	name         string
	description  string
	stdOutput    bool
	outputToFile string
}

func New() (types.Outputter, error) {
	c := new(outputterConfig)
	c.name = Name
	c.description = "Creates a Drone build file"
	c.stdOutput = true
	c.outputToFile = ""

	return c, nil
}

func (oc outputterConfig) Name() string {
	return oc.name
}
func (oc outputterConfig) Output(ctx context.Context, scanResults []types.Scanlet) error {
	var results []types.Scanlet
	// iterate over enabled outputs
	for _, output := range scanResults {
		if output.OutputRender == Name {
			results = append(results, output)
		}
	}
	// lets render the best practice results to stdout
	if len(results) <= 0 {
		return nil
	}
	// lets explain what we added to the drone build
	fmt.Println("")
	fmt.Println("Added to Drone Build:")
	for _, result := range results {
		fmt.Printf(`- %s
`, result.HumanReasoning)
	}

	buildOutput := `kind: pipeline
type: docker
name: default

platform:
	os: linux
	arch: amd64

steps:
`
	for _, result := range results {
		dbo := result.Spec.(types.DroneBuildOutput)
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
