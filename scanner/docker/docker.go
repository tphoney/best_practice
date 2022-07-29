package docker

import (
	"context"
	"fmt"
	"strings"

	"github.com/tphoney/best_practice/outputter"
	"github.com/tphoney/best_practice/outputter/buildmaker"
	"github.com/tphoney/best_practice/outputter/dronebuildanalysis"
	"github.com/tphoney/best_practice/scanner"
	"github.com/tphoney/best_practice/scanner/dronescanner"
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
	dockerFilename    = "Dockerfile"
	Name              = scanner.DockerScannerName
	BuildCheck        = "Docker build"
	SecurityScanCheck = "Docker security scan"
	DroneCheck        = "Docker Drone build"
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
	return []string{BuildCheck, SecurityScanCheck}
}

func (sc *scannerConfig) Scan(ctx context.Context, requestedOutputs []string) (returnVal []types.Scanlet, err error) {
	// lets look for any java files.
	dockerFileMatches, err := scanner.FindMatchingFiles(sc.workingDirectory, dockerFilename, true)
	if err != nil || len(dockerFileMatches) == 0 {
		// nothing to see here, lets leave
		return returnVal, nil
	}

	if sc.runAll || slices.Contains(requestedOutputs, BuildCheck) {
		outputResults := sc.buildCheck(dockerFileMatches)
		returnVal = append(returnVal, outputResults...)
	}
	if sc.runAll || slices.Contains(requestedOutputs, SecurityScanCheck) {
		outputResults := sc.securityCheck(dockerFileMatches)
		returnVal = append(returnVal, outputResults...)
	}
	if (sc.runAll || slices.Contains(requestedOutputs, DroneCheck)) && len(dockerFileMatches) > 0 {
		outputResults, err := sc.droneBuildCheck()
		if err == nil {
			returnVal = append(returnVal, outputResults...)
		}
	}

	return returnVal, nil
}

func (sc *scannerConfig) buildCheck(dockerFiles []string) (outputResults []types.Scanlet) {
	// lets check for the build system
	for i := range dockerFiles {
		testResult := types.Scanlet{
			Name:           BuildCheck,
			ScannerFamily:  Name,
			Description:    "add docker build step, we can upload to acr/dockerhub/ecr/gcr/heroku",
			OutputRenderer: buildmaker.Name,
			Spec: buildmaker.OutputFields{
				Build: buildmaker.Build{
					Name:  fmt.Sprintf("docker build %s", dockerFiles[i]),
					Image: "plugins/docker",
					DroneAppend: fmt.Sprintf(`  privileged: true
    settings:
      repo: organization/docker-image-name
      dry_run: true                       # TODO remove this in production
      auto_tag: true
      dockerfile: %s
      username:
        from_secret: docker_username
      password:
        from_secret: docker_password`, dockerFiles[i]),
				},
				CLI:     fmt.Sprintf("docker build  --rm --no-cache -t organization/docker-image-name:latest -f %s .", dockerFiles[i]),
				HelpURL: "https://plugins.drone.io/plugins/docker",
			},
		}
		outputResults = append(outputResults, testResult)
	}
	return outputResults
}

func (sc *scannerConfig) securityCheck(dockerFiles []string) (outputResults []types.Scanlet) {
	// lets check for the build system
	for i := range dockerFiles {
		testResult := types.Scanlet{
			Name:           BuildCheck,
			ScannerFamily:  Name,
			Description:    "run snyk security scan",
			OutputRenderer: buildmaker.Name,
			Spec: buildmaker.OutputFields{
				Build: buildmaker.Build{
					Name:  fmt.Sprintf("docker build %s", dockerFiles[i]),
					Image: "plugins/drone-snyk",
					DroneAppend: fmt.Sprintf(`  privileged: true
    settings:
      dockerfile: %s
      image: organization/docker-image-name
      snyk:
        from_secret: snyk_token`, dockerFiles[i]),
				},
				CLI:     fmt.Sprintf("docker scan drone-plugins/drone-snyk --file= %s", dockerFiles[i]),
				HelpURL: "snyk.io/help/",
			},
		}
		outputResults = append(outputResults, testResult)
	}
	return outputResults
}

func (sc *scannerConfig) droneBuildCheck() (outputResults []types.Scanlet, err error) {
	pipelines, err := dronescanner.ReadDroneFile(sc.workingDirectory, dronescanner.DroneFileLocation)
	if err != nil {
		return outputResults, err
	}
	// iterate over the pipelines
	foundDockerPlugin := false
	foundSnykPlugin := false
	foundDockerScanCommand := false
	foundDockerBuildCommand := false
	for i := range pipelines {
		for j := range pipelines[i].Steps {
			if strings.Contains(pipelines[i].Steps[j].Image, "plugins/docker") {
				foundDockerPlugin = true
			}
			if strings.Contains(pipelines[i].Steps[j].Image, "plugins/drone-snyk") {
				foundSnykPlugin = true
			}
			commands := pipelines[i].Steps[j].Commands
			for k := range commands {
				if strings.Contains(commands[k], "docker build") {
					foundDockerBuildCommand = true
				}
				if strings.Contains(commands[k], "docker scan") {
					foundDockerScanCommand = true
				}
			}
		}
		if !foundDockerPlugin || foundDockerBuildCommand {
			bestPracticeResult := types.Scanlet{
				Name:           BuildCheck,
				ScannerFamily:  Name,
				Description:    fmt.Sprintf("pipeline '%s' should use the drone docker plugin", pipelines[i].Name),
				OutputRenderer: outputter.DroneBuildAnalysis,
				Spec: dronebuildanalysis.OutputFields{
					HelpURL: "https://plugins.drone.io/plugins/docker",
					Command: "docker build  --rm --no-cache -t organization/docker-image-name:latest -f Dockerfile .",
					RawYaml: `
  - name: build docker
    image: plugins/docker
    settings:
      repo: organization/docker-image-name
      dry_run: true                       # TODO remove this in production
      auto_tag: true
      dockerfile: Dockerfile
      username:
        from_secret: docker_username
      password:
        from_secret: docker_password`,
				},
			}
			outputResults = append(outputResults, bestPracticeResult)
		}
		if !foundSnykPlugin || foundDockerScanCommand {
			bestPracticeResult := types.Scanlet{
				Name:           BuildCheck,
				ScannerFamily:  Name,
				Description:    fmt.Sprintf("pipeline '%s' should use the drone snyk plugin", pipelines[i].Name),
				OutputRenderer: outputter.DroneBuildAnalysis,
				Spec: dronebuildanalysis.OutputFields{
					HelpURL: "snyk.io/help/",
					Command: "docker scan drone-plugins/drone-snyk --file=Dockerfile",
					RawYaml: `
  - name: scan image
    image: plugins/drone-snyk
    privileged: true
    settings:
      dockerfile: Dockerfile
      image: organization/docker-image-name
      snyk:
        from_secret: snyk_token`,
				},
			}
			outputResults = append(outputResults, bestPracticeResult)
		}
	}
	return outputResults, err
}
