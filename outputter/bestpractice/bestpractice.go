package bestpractice

import (
	"context"
	"fmt"

	"github.com/tphoney/best_practice/types"
)

const (
	Name = "best practice"
)

type outputterConfig struct {
	name         string
	stdOutput    bool
	outputToFile string
}

func New() (types.Outputter, error) {
	c := new(outputterConfig)
	c.name = Name
	c.stdOutput = true
	c.outputToFile = ""

	return c, nil
}

func (oc outputterConfig) Name() string {
	return oc.name
}

func (oc outputterConfig) Output(ctx context.Context, scanResults []types.Scanlet) error {
	var bestPracticeResults []types.Scanlet
	// iterate over enabled outputs
	for _, output := range scanResults {
		if output.OutputRender == Name {
			bestPracticeResults = append(bestPracticeResults, output)
		}
	}
	// lets render the best practice results to stdout
	if len(bestPracticeResults) <= 0 {
		return nil
	}
	fmt.Println("Best Practice Results:")
	for _, result := range bestPracticeResults {
		bp := result.Spec.(types.BestPracticeOutput)
		fmt.Printf(`
%s: %s
command: %s
url: %s
`, result.Name, result.HumanReasoning, bp.Command, bp.Url)
	}
	fmt.Println("")
	return nil
}
