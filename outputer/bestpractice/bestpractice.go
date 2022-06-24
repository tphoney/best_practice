package bestpractice

import (
	"context"
	"fmt"

	"github.com/tphoney/best_practice/types"
)

func Output(ctx context.Context, scanResults []types.Scanlet) error {
	var bestPracticeResults []types.Scanlet
	// iterate over enabled outputs
	for _, output := range scanResults {
		if output.OutputRender == types.OutputBestPractice {
			bestPracticeResults = append(bestPracticeResults, output)
		}
	}
	// lets render the best practice results to stdout
	if len(bestPracticeResults) <= 0 {
		return nil
	}

	fmt.Println("Best Practice Results:")
	for _, result := range bestPracticeResults {
		fmt.Printf(`
		%s:
		%+v`, result.HumanReasoning, result.Spec)
	}

	return nil
}
