package outputter

import (
	"context"
	"fmt"

	"github.com/tphoney/best_practice/types"
)

func RunOutput(ctx context.Context, outputters []types.Outputter, scanResults []types.Scanlet) (err error) {
	// iterate over enabled outputs
	for _, outputter := range outputters {
		err = outputter.Output(ctx, scanResults)
		if err != nil {
			fmt.Printf("error running output: %s\n", err)
		}
	}
	// profit
	return nil
}
