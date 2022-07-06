package bestpractice

import (
	"context"
	"fmt"

	"github.com/tphoney/best_practice/outputter"
	"github.com/tphoney/best_practice/types"
)

const (
	Name = outputter.BestPractice
)

type (
	OutputFields struct {
		Command string `json:"command" yaml:"command"`
		HelpURL string `json:"url" yaml:"url"`
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
	oc.description = "Suggests practical changes based on your project layout"
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
		if output.OutputRenderer == Name {
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
`, result.Name, result.Description, bp.Command, bp.HelpURL)
	}
	fmt.Println("")
	return nil
}
