package scanner

import (
	"context"
	"fmt"

	"github.com/tphoney/best_practice/scanner/golang"
	"github.com/tphoney/best_practice/types"
)

func RunScan(ctx context.Context, scansToRun []string) (scanResults []types.Scanlet, err error) {
	for _, scanToRun := range scansToRun {
		switch scanToRun {
		case "golang":
			golangInput := types.ScanInput{
				RequestedOutputs: []string{"dronebuild", "bestpractice", "productrecommendation"},
				RunAll:           true,
			}
			golangResults, err := golang.Scan(ctx, golangInput)
			if err != nil {
				fmt.Printf("error running golang scan: %s\n", err)
			}
			scanResults = append(scanResults, golangResults...)
		default:
			fmt.Printf("scanlet %s not found", scanToRun)
		}
	}
	return scanResults, nil
}
