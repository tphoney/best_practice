// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"

	"github.com/tphoney/best_practice/outputter"
	"github.com/tphoney/best_practice/outputter/bestpractice"
	"github.com/tphoney/best_practice/outputter/dronebuild"
	"github.com/tphoney/best_practice/scanner"
	"github.com/tphoney/best_practice/scanner/golang"
	"github.com/tphoney/best_practice/types"
)

// Args provides plugin execution arguments.
type Args struct {
	Pipeline

	// Level defines the plugin log level.
	Level string `envconfig:"PLUGIN_LOG_LEVEL"`

	RequestedScanners []string `envconfig:"PLUGIN_REQUESTED_SCANNERS"`
	RequestedOutputs  []string `envconfig:"PLUGIN_REQUESTED_OUTPUTS"`
}

// Exec executes the plugin.
func Exec(ctx context.Context, args *Args) error {
	fmt.Println("==========================")
	// setup requested scanners
	if len(args.RequestedScanners) == 0 {
		args.RequestedScanners = []string{golang.Name}
	}
	scanners := make([]types.Scanner, 0)
	for _, scannerName := range args.RequestedScanners {
		switch scannerName {
		case golang.Name:
			// create golang scanner
			g, err := golang.New()
			if err != nil {
				return err
			}
			scanners = append(scanners, g)
		default:
			fmt.Printf("unknown scanner: %s\n", scannerName)
		}
	}
	if len(scanners) == 0 {
		return fmt.Errorf("no scanners requested")
	}
	// setup requested outputs
	if len(args.RequestedOutputs) == 0 {
		args.RequestedOutputs = []string{dronebuild.Name, bestpractice.Name}
	}
	outputters := make([]types.Outputter, 0)
	for _, outputName := range args.RequestedOutputs {
		switch outputName {
		case dronebuild.Name:
			db, _ := dronebuild.New(dronebuild.WithOutputToFile(".drone.yml.new"), dronebuild.WithStdOutput(false))
			outputters = append(outputters, db)
		case bestpractice.Name:
			bp, _ := bestpractice.New(bestpractice.WithStdOutput(true))
			outputters = append(outputters, bp)
		default:
			fmt.Printf("unknown output: %s", outputName)
		}
	}
	if len(outputters) == 0 {
		return fmt.Errorf("no outputters selected")
	}

	fmt.Println("scanners used:")
	for i := range scanners {
		fmt.Printf("%s - %s\n", scanners[i].Name(), scanners[i].Description())
	}
	scanResults, scanErr := scanner.RunScanners(ctx, scanners, args.RequestedOutputs)
	if scanErr != nil {
		fmt.Printf("error running scan failed: %s\n", scanErr)
		return scanErr
	}
	fmt.Println("outputs used:")
	for i := range outputters {
		fmt.Printf("%s - %s\n", outputters[i].Name(), outputters[i].Description())
	}
	fmt.Println("==========================")
	// run output engine
	outputErr := outputter.RunOutput(ctx, outputters, scanResults)
	if outputErr != nil {
		fmt.Printf("error running output failed: %s\n", outputErr)
		return outputErr
	}
	// profit
	return nil
}
