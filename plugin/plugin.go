// Copyright 2020 the Drone Authors. All rights reserved.
// Use of this source code is governed by the Blue Oak Model License
// that can be found in the LICENSE file.

package plugin

import (
	"context"
	"fmt"

	"github.com/tphoney/best_practice/outputer"
	"github.com/tphoney/best_practice/scanner"
)

// Args provides plugin execution arguments.
type Args struct {
	Pipeline

	// Level defines the plugin log level.
	Level string `envconfig:"PLUGIN_LOG_LEVEL"`

	// TODO replace or remove
	Param1 string `envconfig:"PLUGIN_PARAM1"`
	Param2 string `envconfig:"PLUGIN_PARAM2"`
}

// Exec executes the plugin.
func Exec(ctx context.Context, args Args) error {
	// create scan engine
	scanResults, scanErr := scanner.RunScan(ctx, []string{"golang"})
	if scanErr != nil {
		fmt.Printf("error running scan failed: %s\n", scanErr)
		return scanErr
	}
	// run output engine
	outputErr := outputer.RunOutput(ctx, scanResults)
	if outputErr != nil {
		fmt.Printf("error running output failed: %s\n", outputErr)
		return outputErr
	}
	// iterate over outputs
	// profit
	return nil
}
