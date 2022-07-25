package harnessproduct

import (
	"context"
	"fmt"

	"github.com/tphoney/best_practice/outputter"
	"github.com/tphoney/best_practice/types"
)

const (
	Name = outputter.HarnessProduct
)

type (
	OutputFields struct {
		ProductName string `json:"product_name"`
		URL         string `json:"url"`
		Explanation string `json:"explanation"`
		Why         string `json:"why"`
	}

	outputterConfig struct {
		name             string
		description      string
		workingDirectory string
	}
)

func New(opts ...Option) (types.Outputter, error) {
	oc := new(outputterConfig)
	oc.name = Name
	oc.description = "Shows harness product recommendations"
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
	if len(results) == 0 {
		return nil
	}
	productOutput := "Product Recommendations\n"
	// add the steps to the build file
	for _, result := range results {
		dbo := result.Spec.(OutputFields)
		productOutput += fmt.Sprintf(
			`- %s
  URL: %s
  Explanation: %s
  Why: %s
`, dbo.ProductName, dbo.URL, dbo.Explanation, dbo.Why)
	}
	fmt.Println(productOutput)
	return nil
}
