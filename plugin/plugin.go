// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"
	"os"

	"github.com/tphoney/best_practice/outputter"
	"github.com/tphoney/best_practice/outputter/buildanalysis"
	"github.com/tphoney/best_practice/outputter/dronebuildmaker"
	"github.com/tphoney/best_practice/outputter/harnessproduct"
	"github.com/tphoney/best_practice/scanner"
	"github.com/tphoney/best_practice/scanner/docker"
	"github.com/tphoney/best_practice/scanner/dronescanner"
	"github.com/tphoney/best_practice/scanner/golang"
	"github.com/tphoney/best_practice/scanner/java"
	"github.com/tphoney/best_practice/scanner/javascript"
	"github.com/tphoney/best_practice/scanner/ruby"
	"github.com/tphoney/best_practice/types"
)

// Args provides plugin execution arguments.
type Args struct {
	Pipeline

	// Level defines the plugin log level.
	Level string `envconfig:"PLUGIN_LOG_LEVEL"`

	RequestedScanners []string `envconfig:"PLUGIN_REQUESTED_SCANNERS"`
	RequestedOutputs  []string `envconfig:"PLUGIN_REQUESTED_OUTPUTS"`
	WorkingDirectory  string   `envconfig:"PLUGIN_WORKING_DIRECTORY"`
}

// Exec executes the plugin.
func Exec(ctx context.Context, args *Args) error { // nolint:gocyclo
	fmt.Println("==========================")
	// setup the base directory
	if args.WorkingDirectory == "" {
		args.WorkingDirectory = os.Getenv("DRONE_WORKSPACE")
		if args.WorkingDirectory == "" {
			args.WorkingDirectory, _ = os.Getwd()
		}
	}
	fmt.Println("working directory:", args.WorkingDirectory)
	// setup requested scanners
	if len(args.RequestedScanners) == 0 {
		args.RequestedScanners = scanner.ListScannersNames()
	}
	scanners := make([]types.Scanner, 0)
	for _, scannerName := range args.RequestedScanners {
		switch scannerName {
		case scanner.DockerScannerName:
			d, err := docker.New(docker.WithWorkingDirectory(args.WorkingDirectory))
			if err != nil {
				return err
			}
			scanners = append(scanners, d)
		case scanner.DroneScannerName:
			d, err := dronescanner.New(dronescanner.WithWorkingDirectory(args.WorkingDirectory))
			if err != nil {
				return err
			}
			scanners = append(scanners, d)
		case scanner.GolangScannerName:
			g, err := golang.New(golang.WithWorkingDirectory(args.WorkingDirectory))
			if err != nil {
				return err
			}
			scanners = append(scanners, g)
		case scanner.JavascriptScannerName:
			j, err := javascript.New(javascript.WithWorkingDirectory(args.WorkingDirectory))
			if err != nil {
				return err
			}
			scanners = append(scanners, j)
		case scanner.JavaScannerName:
			j, err := java.New(java.WithWorkingDirectory(args.WorkingDirectory))
			if err != nil {
				return err
			}
			scanners = append(scanners, j)
		case scanner.RubyScannerName:
			r, err := ruby.New(ruby.WithWorkingDirectory(args.WorkingDirectory))
			if err != nil {
				return err
			}
			scanners = append(scanners, r)
		default:
			fmt.Printf("unknown scanner: %s\n", scannerName)
		}
	}
	if len(scanners) == 0 {
		return fmt.Errorf("no scanners requested")
	}
	// setup requested outputs
	if len(args.RequestedOutputs) == 0 {
		args.RequestedOutputs = outputter.ListOutputterNames()
	}
	outputters := make([]types.Outputter, 0)
	for _, outputName := range args.RequestedOutputs {
		switch outputName {
		case outputter.DroneBuildMaker:
			db, _ := dronebuildmaker.New(dronebuildmaker.WithWorkingDirectory(args.WorkingDirectory), dronebuildmaker.WithStdOutput(false), dronebuildmaker.WithOutputToFile(true))
			outputters = append(outputters, db)
		case outputter.DroneBuildAnalysis:
			bp, _ := buildanalysis.New(buildanalysis.WithStdOutput(true))
			outputters = append(outputters, bp)
		case outputter.HarnessProduct:
			hp, _ := harnessproduct.New()
			outputters = append(outputters, hp)
		default:
			fmt.Printf("unknown output: %s", outputName)
		}
	}
	if len(outputters) == 0 {
		return fmt.Errorf("no outputters selected")
	}

	fmt.Println("scanners used:")
	for i := range scanners {
		fmt.Printf("%s\n", scanners[i].Name())
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
