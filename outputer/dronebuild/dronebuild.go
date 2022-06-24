package dronebuild

import (
	"context"
	"fmt"

	"github.com/tphoney/best_practice/types"
)

func Output(ctx context.Context, scanResults []types.Scanlet) error {
	var results []types.Scanlet
	// iterate over enabled outputs
	for _, output := range scanResults {
		if output.OutputRender == types.OutputBestPractice {
			results = append(results, output)
		}
	}
	// lets render the best practice results to stdout
	if len(results) <= 0 {
		return nil
	}
	// lets explain what we added to the drone build
	fmt.Println("Added to Drone Build:")
	for _, result := range results {
		fmt.Printf(`- %s
`, result.HumanReasoning)
	}
	fmt.Println("")
	fmt.Println("")

	fmt.Printf(`kind: pipeline
type: docker
name: default

platform:
	os: linux
	arch: amd64

steps:
`)
	for _, result := range results {
		fmt.Printf(`%+v`, result.Spec)
	}

	return nil
}
