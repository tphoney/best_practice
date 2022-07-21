package scanner

import (
	"context"
	"fmt"

	"github.com/tphoney/best_practice/types"
)

const (
	DockerScannerName     = "Docker"
	DroneScannerName      = "Drone"
	GolangScannerName     = "Golang"
	JavaScannerName       = "Java"
	JavascriptScannerName = "Javascript"
	RubyScannerName       = "Ruby"
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
	// run language scanners first, then scanners that may depend on them
	return []string{GolangScannerName, JavaScannerName, JavascriptScannerName, RubyScannerName, DockerScannerName, DroneScannerName}
}
