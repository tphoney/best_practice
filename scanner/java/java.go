package java

import (
	"context"
	"os"
	"path/filepath"

	"github.com/tphoney/best_practice/outputter/bestpractice"
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
	mavenFolderLocation = ".mvn"
	antBuildFile        = "build.xml"
	gradleSettingsFile  = "settings.gradle"
	bazelBuildFile      = "BUILD.bazel"
	Name                = scanner.JavaScannerName
	BuildSystemCheck    = "build"
	UnitTestCheck       = "junit"
)

func New(opts ...Option) (types.Scanner, error) {
	sc := new(scannerConfig)
	sc.name = Name
	sc.description = "checks for various java related best practices"
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
	return []string{BuildSystemCheck, UnitTestCheck}
}

func (sc *scannerConfig) Scan(ctx context.Context, requestedOutputs []string) (returnVal []types.Scanlet, err error) {
	// lets look for any java files.
	matches, err := scanner.FindMatchingFiles(sc.workingDirectory, "*.java")
	if err != nil || len(matches) == 0 {
		// nothing to see here, lets leave
		return returnVal, nil
	}
	// check for test folders
	testMatches, err := scanner.FindMatchingFiles(sc.workingDirectory, "**/test/*.java")
	if err != nil || len(testMatches) == 0 {
		// add a best practice for adding tests
		bestPracticeResult := types.Scanlet{
			Name:           "add_tests",
			ScannerFamily:  Name,
			Description:    "a java project should have tests, running them depends on the build system",
			OutputRenderer: bestpractice.Name,
			Spec: bestpractice.OutputFields{
				Command: "javac -d /absolute/path/for/compiled/classes -cp /absolute/path/to/junit-4.12.jar /absolute/path/to/TestClassName.java",
				HelpURL: "some help",
			},
		}
		returnVal = append(returnVal, bestPracticeResult)
	}
	// check for the various build systems
	if sc.runAll || slices.Contains(requestedOutputs, BuildSystemCheck) {
		_, outputResults := sc.buildCheck(len(testMatches) == 0)
		if len(outputResults) > 0 {
			returnVal = append(returnVal, outputResults...)
		}
	}

	return returnVal, nil
}

func (sc *scannerConfig) buildCheck(hasTests bool) (buildType []string, outputResults []types.Scanlet) {
	// lets check for the build system
	_, err := os.Stat(filepath.Join(sc.workingDirectory, bazelBuildFile))
	if err == nil {
		testResult := types.Scanlet{
			Name:           BuildSystemCheck,
			ScannerFamily:  Name,
			Description:    "run tests",
			OutputRenderer: dronebuild.Name,
			Spec: dronebuild.OutputFields{
				RawYaml: `  - name: test
    image: google/bazel
    commands:
      - bazel test`,
			},
		}
		outputResults = append(outputResults, testResult)
		droneBuildResult := types.Scanlet{
			Name:           BuildSystemCheck,
			ScannerFamily:  Name,
			Description:    "run bazel build",
			OutputRenderer: dronebuild.Name,
			Spec: dronebuild.OutputFields{
				RawYaml: `  - name: bazel build
    image: google/bazel
    commands:
      - bazel build :all`,
			},
		}
		outputResults = append(outputResults, droneBuildResult)
	}
	// it may be a maven project
	_, err = os.Stat(filepath.Join(sc.workingDirectory, mavenFolderLocation))
	if err == nil {
		testResult := types.Scanlet{
			Name:           BuildSystemCheck,
			ScannerFamily:  Name,
			Description:    "run tests",
			OutputRenderer: dronebuild.Name,
			Spec: dronebuild.OutputFields{
				RawYaml: `  - name: test
    image: maven
    commands:
      - mvn test`,
			},
		}
		outputResults = append(outputResults, testResult)
		droneBuildResult := types.Scanlet{
			Name:           BuildSystemCheck,
			ScannerFamily:  Name,
			Description:    "run maven build",
			OutputRenderer: dronebuild.Name,
			Spec: dronebuild.OutputFields{
				RawYaml: `  - name: maven build
    image: maven
    commands:
      - mvn clean install`,
			},
		}
		outputResults = append(outputResults, droneBuildResult)
		buildType = append(buildType, "maven")
	}
	// it may be a gradle project
	_, err = os.Stat(filepath.Join(sc.workingDirectory, gradleSettingsFile))
	if err == nil {
		testResult := types.Scanlet{
			Name:           BuildSystemCheck,
			ScannerFamily:  Name,
			Description:    "run tests",
			OutputRenderer: dronebuild.Name,
			Spec: dronebuild.OutputFields{
				RawYaml: `  - name: test
    image: gradle/gradle
    commands:
      - ./gradlew test`,
			},
		}
		outputResults = append(outputResults, testResult)
		droneBuildResult := types.Scanlet{
			Name:           BuildSystemCheck,
			ScannerFamily:  Name,
			Description:    "run gradle build",
			OutputRenderer: dronebuild.Name,
			Spec: dronebuild.OutputFields{
				RawYaml: `  - name: gradle build
	image: gradle/gradle
	commands:
	- ./gradlew clean build`,
			},
		}
		outputResults = append(outputResults, droneBuildResult)
		buildType = append(buildType, "gradle")
	}
	// it may be an ant project
	_, err = os.Stat(filepath.Join(sc.workingDirectory, antBuildFile))
	if err == nil {
		testResult := types.Scanlet{
			Name:           BuildSystemCheck,
			ScannerFamily:  Name,
			Description:    "run tests",
			OutputRenderer: dronebuild.Name,
			Spec: dronebuild.OutputFields{
				RawYaml: `  - name: test
		image: frekele/ant/
		commands:
		  - ant -buildfile build.xml test`,
			},
		}
		outputResults = append(outputResults, testResult)
		droneBuildResult := types.Scanlet{
			Name:           BuildSystemCheck,
			ScannerFamily:  Name,
			Description:    "run ant build",
			OutputRenderer: dronebuild.Name,
			Spec: dronebuild.OutputFields{
				RawYaml: `  - name: ant build
		image: frekele/ant
		commands:
		- ant -buildfile build.xml`,
			},
		}
		outputResults = append(outputResults, droneBuildResult)
		buildType = append(buildType, "gradle")
	}
	return buildType, outputResults
}
