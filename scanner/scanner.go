package scanner

import (
	"context"
	"fmt"

	"github.com/tphoney/best_practice/types"
)

const (
	JavaScannerName       = "java"
	JavascriptScannerName = "javascript"
	GolangScannerName     = "golang"
	DroneScannerName      = "drone"
)

func RunScanners(ctx context.Context, scannersToRun []types.Scanner, requestedOutputs []string) (scanResults []types.Scanlet, err error) {
	for _, scannerToRun := range scannersToRun {
		results, err := scannerToRun.Scan(ctx, requestedOutputs)
		if err != nil {
			fmt.Printf("error running '%s' scan: %s\n", scannerToRun.Name(), err)
		}
		scanResults = append(scanResults, results...)
	}
	return scanResults, nil
}

func ListScannersNames() []string {
	return []string{JavaScannerName, JavascriptScannerName, GolangScannerName, DroneScannerName}
}
