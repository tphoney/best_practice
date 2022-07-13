package buildanalysis

import (
	"context"
	"fmt"

	"github.com/tphoney/best_practice/outputter"
	"github.com/tphoney/best_practice/types"
)

const (
	Name = outputter.DroneBuildAnalysis
)

type (
	outputterConfig struct {
		name         string
		description  string
		stdOutput    bool
		outputToFile string
	}

	OutputFields struct {
		RawYaml string `json:"raw_yaml" yaml:"raw_yaml"`
		Command string `json:"command" yaml:"command"`
		HelpURL string `json:"url" yaml:"url"`
	}
)

func New(opts ...Option) (types.Outputter, error) {
	oc := new(outputterConfig)
	oc.name = Name
	oc.description = "Suggests practical changes based on your project layout and build file"
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
	var bestPracticeResults []types.Scanlet
	// iterate over enabled outputs
	for _, output := range scanResults {
		if output.OutputRenderer == outputter.DroneBuildAnalysis {
			bestPracticeResults = append(bestPracticeResults, output)
		}
	}
	// lets render the best practice results to stdout
	if len(bestPracticeResults) == 0 {
		return nil
	}
	fmt.Println("Best Practice Results:")
	for _, result := range bestPracticeResults {
		bp := result.Spec.(OutputFields)
		fmt.Printf(`- %s: %s
  command to run: "%s"
  url: %s
  Drone YAML:
%s
`, result.Name, result.Description, bp.Command, bp.HelpURL, bp.RawYaml)
	}
	fmt.Println("")
	return nil
}
