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
func Exec(ctx context.Context, args Args) error {
	// setup requested outputs
	if len(args.RequestedOutputs) == 0 {
		args.RequestedOutputs = []string{dronebuild.Name, bestpractice.Name}
	}
	outputters := make([]types.Outputter, 0)
	for _, outputName := range args.RequestedOutputs {
		switch outputName {
		case dronebuild.Name:
			db, _ := dronebuild.New()
			outputters = append(outputters, db)
		case bestpractice.Name:
			bp, _ := bestpractice.New()
			outputters = append(outputters, bp)
		default:
			fmt.Printf("unknown output: %s", outputName)
		}
	}
	if len(outputters) == 0 {
		return fmt.Errorf("no outputters selected")
	}
	// create golang scanner
	golang, err := golang.New()
	if err != nil {
		return err
	}
	scanResults, scanErr := scanner.RunScanners(ctx, []types.Scanner{golang}, args.RequestedOutputs)
	if scanErr != nil {
		fmt.Printf("error running scan failed: %s\n", scanErr)
		return scanErr
	}
	// run output engine
	outputErr := outputter.RunOutput(ctx, outputters, scanResults)
	if outputErr != nil {
		fmt.Printf("error running output failed: %s\n", outputErr)
		return outputErr
	}
	// profit
	return nil
}
