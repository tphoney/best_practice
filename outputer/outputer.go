package outputer

import (
	"context"
	"fmt"

	"github.com/tphoney/best_practice/outputer/bestpractice"
	"github.com/tphoney/best_practice/outputer/dronebuild"
	"github.com/tphoney/best_practice/types"
)

func RunOutput(ctx context.Context, scanResults []types.Scanlet) (err error) {
	// iterate over enabled outputs
	// best practice
	bestPracticeErr := bestpractice.Output(ctx, scanResults)
	if bestPracticeErr != nil {
		fmt.Printf("error running best practice output: %s\n", bestPracticeErr)
		return bestPracticeErr
	}
	dronebuild := dronebuild.Output(ctx, scanResults)
	if dronebuild != nil {
		fmt.Printf("error running drone build output: %s\n", dronebuild)
		return dronebuild
	}
	// profit
	return nil
}
