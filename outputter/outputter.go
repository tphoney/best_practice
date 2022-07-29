package outputter

import (
	"context"
	"fmt"

	"github.com/tphoney/best_practice/types"
)

const (
	BuildMaker         = "build maker"
	HarnessProduct     = "harness product"
	DroneBuildAnalysis = "drone build analysis"
)

func RunOutput(ctx context.Context, outputters []types.Outputter, scanResults []types.Scanlet) (err error) {
	// iterate over enabled outputs
	for _, outputter := range outputters {
		fmt.Println("++++++++++++++++++++++++++")
		err = outputter.Output(ctx, scanResults)
		if err != nil {
			fmt.Printf("error running output: %s\n", err)
		}
	}
	// profit
	return nil
}

func ListOutputterNames() []string {
	return []string{BuildMaker, HarnessProduct, DroneBuildAnalysis}
}
