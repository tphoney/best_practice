package buildmaker

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/tphoney/best_practice/outputter"
	"github.com/tphoney/best_practice/types"
)

const (
	Name = outputter.BuildMaker
)

var (
	droneFileName = ".drone.yml"
	cieFileName   = ".cie.yml"
)

type (
	OutputFields struct {
		Build
		CLI     string `json:"cli" yaml:"cli"`
		HelpURL string `json:"help_url" yaml:"help_url"`
	}

	Build struct {
		Name        string   `json:"name" yaml:"name"`
		Image       string   `json:"image" yaml:"image"`
		Commands    []string `json:"commands" yaml:"commands"`
		DroneAppend string   `json:"drone_append" yaml:"drone_append"`
		CIEAppend   string   `json:"cie_append" yaml:"cie_append"`
	}

	outputterConfig struct {
		name             string
		description      string
		stdOutput        bool
		workingDirectory string
		outputToFile     bool
		outputDrone      bool
		outputCIE        bool
	}
)

func New(opts ...Option) (types.Outputter, error) {
	oc := new(outputterConfig)
	oc.name = Name
	oc.description = "Creates a full build file"
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
	if len(results) == 0 {
		return nil
	}
	// lets explain what we added to the build
	fmt.Println("Build Generator\n\nAdding the following to the build:")
	for _, result := range results {
		fmt.Printf("- %s, %s\n", result.ScannerFamily, result.Description)
	}
	fmt.Println("")

	droneBuildOutput := `kind: pipeline
type: docker
name: default

platform:
  os: linux
  arch: amd64

steps:`
	cieBuildOutput := "CIECIE\n"

	// add the steps to the build file
	for _, result := range results {
		dbo := result.Spec.(OutputFields)
		// build drone step
		if oc.outputDrone {
			droneBuildOutput += fmt.Sprintf(`
  - name: %s
    image: %s`, dbo.Name, dbo.Image)
			if len(dbo.Commands) > 0 {
				droneBuildOutput += "\n    commands:"
				for _, command := range dbo.Commands {
					droneBuildOutput += fmt.Sprintf("\n      - %s", command)
				}
			}
			if dbo.DroneAppend != "" {
				droneBuildOutput += fmt.Sprintf(`
  %s`, dbo.DroneAppend)
			}
		}
		// build cie step
		if oc.outputCIE {
			cieBuildOutput += fmt.Sprintf("step %s\n", dbo.Name)
		}
	}

	if oc.stdOutput {
		if oc.outputDrone {
			fmt.Printf("Drone build file:\n%s\n", droneBuildOutput)
			fmt.Println(droneBuildOutput)
		}
		if oc.outputCIE {
			fmt.Printf("CIE build file:\n%s\n", cieBuildOutput)
			fmt.Println(cieBuildOutput)
		}
	}
	if oc.outputToFile {
		if oc.outputDrone {
			_, err := os.Stat(filepath.Join(oc.workingDirectory, droneFileName))
			if err == nil {
				// file exists append .new to the file name
				droneFileName += ".new"
			}
			fmt.Printf("Created a new Drone Build file '%s'\n", filepath.Join(oc.workingDirectory, droneFileName))
			writeErr := outputter.WriteToFile(filepath.Join(oc.workingDirectory, droneFileName), droneBuildOutput)
			if writeErr != nil {
				return writeErr
			}
		}
		if oc.outputCIE {
			_, err := os.Stat(filepath.Join(oc.workingDirectory, cieFileName))
			if err == nil {
				// file exists append .new to the file name
				droneFileName += ".new"
			}
			fmt.Printf("Created a new CIE Build file '%s'\n", filepath.Join(oc.workingDirectory, cieFileName))
			writeErr := outputter.WriteToFile(filepath.Join(oc.workingDirectory, cieFileName), cieBuildOutput)
			if writeErr != nil {
				return writeErr
			}
		}
	}
	return nil
}
