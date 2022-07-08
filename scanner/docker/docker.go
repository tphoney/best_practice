package docker

import (
	"context"
	"fmt"

	"github.com/tphoney/best_practice/outputter/dronebuild"
	"github.com/tphoney/best_practice/scanner"
	"github.com/tphoney/best_practice/types"
	"golang.org/x/exp/slices"
)

type scannerConfig struct {
	name             string
	description      string
	workingDirectory string
	checksToRun      []string
	runAll           bool
}

const (
	dockerFilename  = "Dockerfile"
	Name            = scanner.DockerScannerName
	DockerFileCheck = "docker file"
)

func New(opts ...Option) (types.Scanner, error) {
	sc := new(scannerConfig)
	sc.name = Name
	sc.description = "checks for various docker related best practices"
	sc.runAll = true
	// apply options
	for _, opt := range opts {
		opt(sc)
	}

	return sc, nil
}

func (sc *scannerConfig) Name() string {
	return sc.name
}

func (sc *scannerConfig) Description() string {
	return sc.description
}

func (sc *scannerConfig) AvailableChecks() []string {
	return []string{DockerFileCheck}
}

func (sc *scannerConfig) Scan(ctx context.Context, requestedOutputs []string) (returnVal []types.Scanlet, err error) {
	// lets look for any java files.
	dockerFileMatches, err := scanner.FindMatchingFiles(sc.workingDirectory, dockerFilename)
	if err != nil || len(dockerFileMatches) == 0 {
		// nothing to see here, lets leave
		return returnVal, nil
	}

	if sc.runAll || slices.Contains(requestedOutputs, DockerFileCheck) {
		outputResults := sc.buildCheck(dockerFileMatches)
		if len(outputResults) > 0 {
			returnVal = append(returnVal, outputResults...)
		}
	}

	return returnVal, nil
}

func (sc *scannerConfig) buildCheck(dockerFiles []string) (outputResults []types.Scanlet) {
	// lets check for the build system
	for i := range dockerFiles {
		testResult := types.Scanlet{
			Name:           DockerFileCheck,
			ScannerFamily:  Name,
			Description:    "add docker build step",
			OutputRenderer: dronebuild.Name,
			Spec: dronebuild.OutputFields{
				RawYaml: fmt.Sprintf(`  - name: build %s
	image: plugins/docker
	settings:
	repo: organization/docker-image-name
	  dry_run: true                       # TODO remove this in production      
	  auto_tag: true
	  dockerfile: %s
	  username:
		from_secret: docker_username
	  password:
		from_secret: docker_password
	- name: scan
	  image: drone-plugins/drone-snyk
	  privileged: true
	  settings:
		dockerfile: %s
		image: organization/docker-image-name 
		snyk:
		  from_secret: snyk_token`, dockerFiles[i], dockerFiles[i], dockerFiles[i]),
			},
		}
		outputResults = append(outputResults, testResult)
	}
	return outputResults
}
