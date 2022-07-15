package java

import (
	"context"
	"os"
	"path/filepath"
	"strings"

	"github.com/tphoney/best_practice/outputter"
	"github.com/tphoney/best_practice/outputter/buildanalysis"
	"github.com/tphoney/best_practice/outputter/dronebuildmaker"
	"github.com/tphoney/best_practice/outputter/harnessproduct"
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
	androidManifest     = "AndroidManifest.xml"
	antBuildFile        = "build.xml"
	bazelBuildFile      = "BUILD.bazel"
	gradleSettingsFile  = "settings.gradle"
	mavenFolderLocation = ".mvn"

	Name         = scanner.JavaScannerName
	BuildCheck   = "java build"
	TestCheck    = "java test"
	AndroidCheck = "java android"
	DroneCheck   = "java drone build"
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
	return []string{BuildCheck, TestCheck, DroneCheck, AndroidCheck}
}

func (sc *scannerConfig) Scan(ctx context.Context, requestedOutputs []string) (returnVal []types.Scanlet, err error) {
	// lets look for any java files.
	matches, err := scanner.FindMatchingFiles(sc.workingDirectory, "*.java")
	if err != nil || len(matches) == 0 {
		// nothing to see here, lets leave
		return returnVal, nil
	}
	// check for test folders
	testMatches, err := scanner.FindMatchingFolders(sc.workingDirectory, "test")
	if err != nil || len(testMatches) == 0 {
		// add a best practice for adding tests
		bestPracticeResult := types.Scanlet{
			Name:           "add_tests",
			ScannerFamily:  Name,
			Description:    "a java project should have tests, running them depends on the build system",
			OutputRenderer: outputter.DroneBuildMaker,
			Spec: dronebuildmaker.OutputFields{
				Command: "javac -d /absolute/path/for/compiled/classes -cp /absolute/path/to/junit-4.12.jar /absolute/path/to/TestClassName.java",
				HelpURL: "http://users.csc.calpoly.edu/~djanzen/research/TDD08/cdesai/IntroducingJUnit/IntroducingJUnit.html",
			},
		}
		returnVal = append(returnVal, bestPracticeResult)
	}
	if len(testMatches) > 0 {
		// recommend test intelligence
		harnessProductResult := types.Scanlet{
			Name:           "test_intelligence",
			ScannerFamily:  Name,
			Description:    "java tests found",
			OutputRenderer: outputter.HarnessProduct,
			Spec: harnessproduct.OutputFields{
				ProductName: "Test Intelligence",
				URL:         "https://harness.io/blog/continuous-integration/test-intelligence/",
				Explanation: "Test Intelligence reduces the amount of time spent running tests by intelligently running only the necessary tests.",
				Why:         "Detected Java Junit tests",
			},
		}
		returnVal = append(returnVal, harnessProductResult)
	}
	// check for the various build systems
	var buildTypes []string
	if sc.runAll || slices.Contains(requestedOutputs, BuildCheck) {
		_, outputResults := sc.buildCheck(len(testMatches) == 0)
		if len(outputResults) > 0 {
			returnVal = append(returnVal, outputResults...)
		}
	}
	// check for android
	foundAndroid := false
	if sc.runAll || slices.Contains(requestedOutputs, AndroidCheck) {
		androidMatches, err := scanner.FindMatchingFiles(sc.workingDirectory, androidManifest)
		if err == nil || len(androidMatches) > 0 {
			androidScanlet := types.Scanlet{
				Name:           "android",
				ScannerFamily:  Name,
				Description:    "run android specific project tools",
				OutputRenderer: outputter.DroneBuildMaker,
				Spec: dronebuildmaker.OutputFields{
					Command: "sdkmanager --list",
					HelpURL: "https://developer.android.com/studio/command-line/sdkmanager.html",
					RawYaml: `  - name: android sdk
          image: androidsdk/android-31
	        commands:
	          - sdkmanager --list`,
				},
			}
			foundAndroid = true
			returnVal = append(returnVal, androidScanlet)
		}
	}
	if sc.runAll || slices.Contains(requestedOutputs, DroneCheck) {
		outputResults, err := sc.droneCheck(buildTypes, foundAndroid)
		if err == nil {
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
			Name:           BuildCheck,
			ScannerFamily:  Name,
			Description:    "run tests",
			OutputRenderer: dronebuildmaker.Name,
			Spec: dronebuildmaker.OutputFields{
				RawYaml: `  - name: test
    image: google/bazel
    commands:
      - bazel test`,
			},
		}
		outputResults = append(outputResults, testResult)
		droneBuildResult := types.Scanlet{
			Name:           BuildCheck,
			ScannerFamily:  Name,
			Description:    "run bazel build",
			OutputRenderer: dronebuildmaker.Name,
			Spec: dronebuildmaker.OutputFields{
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
			Name:           BuildCheck,
			ScannerFamily:  Name,
			Description:    "run tests",
			OutputRenderer: dronebuildmaker.Name,
			Spec: dronebuildmaker.OutputFields{
				RawYaml: `  - name: test
    image: maven
    commands:
      - mvn test`,
			},
		}
		outputResults = append(outputResults, testResult)
		droneBuildResult := types.Scanlet{
			Name:           BuildCheck,
			ScannerFamily:  Name,
			Description:    "run maven build",
			OutputRenderer: dronebuildmaker.Name,
			Spec: dronebuildmaker.OutputFields{
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
			Name:           BuildCheck,
			ScannerFamily:  Name,
			Description:    "run tests",
			OutputRenderer: dronebuildmaker.Name,
			Spec: dronebuildmaker.OutputFields{
				RawYaml: `  - name: test
    image: gradle/gradle
    commands:
      - ./gradlew test`,
			},
		}
		outputResults = append(outputResults, testResult)
		droneBuildResult := types.Scanlet{
			Name:           BuildCheck,
			ScannerFamily:  Name,
			Description:    "run gradle build",
			OutputRenderer: dronebuildmaker.Name,
			Spec: dronebuildmaker.OutputFields{
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
			Name:           BuildCheck,
			ScannerFamily:  Name,
			Description:    "run tests",
			OutputRenderer: dronebuildmaker.Name,
			Spec: dronebuildmaker.OutputFields{
				RawYaml: `  - name: test
		image: frekele/ant/
		commands:
		  - ant -buildfile build.xml test`,
			},
		}
		outputResults = append(outputResults, testResult)
		droneBuildResult := types.Scanlet{
			Name:           BuildCheck,
			ScannerFamily:  Name,
			Description:    "run ant build",
			OutputRenderer: dronebuildmaker.Name,
			Spec: dronebuildmaker.OutputFields{
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

func (sc *scannerConfig) droneCheck(buildTypes []string, hasAndroid bool) (outputResults []types.Scanlet, err error) {
	pipelines, err := dronescanner.ReadDroneFile(sc.workingDirectory, dronescanner.DroneFileLocation)
	if err != nil {
		return outputResults, err
	}
	foundBazelTest := false
	foundBazelBuild := false
	foundMavenTest := false
	foundMavenBuild := false
	foundGradleTest := false
	foundGradleBuild := false
	foundAndroidCommands := false
	// iterate over the pipelines
	for i := range pipelines {
		for j := range pipelines[i].Steps {
			commands := pipelines[i].Steps[j].Commands
			for k := range commands {
				if strings.Contains(commands[k], "bazel test") {
					foundBazelTest = true
				}
				if strings.Contains(commands[k], "bazel build") {
					foundBazelBuild = true
				}
				if strings.Contains(commands[k], "mvn test") {
					foundMavenTest = true
				}
				if strings.Contains(commands[k], "mvn clean install") {
					foundMavenBuild = true
				}
				if strings.Contains(commands[k], "gradlew test") {
					foundGradleTest = true
				}
				if strings.Contains(commands[k], "gradlew clean build") {
					foundGradleBuild = true
				}
				if strings.Contains(commands[k], "sdkmanager") || strings.Contains(commands[k], "adb") {
					foundAndroidCommands = true
				}
			}
		}
		if foundBazelTest {
			testResult := types.Scanlet{
				Name:           BuildCheck,
				ScannerFamily:  Name,
				Description:    "run bazel tests",
				OutputRenderer: outputter.DroneBuildAnalysis,
				Spec: buildanalysis.OutputFields{
					RawYaml: `  - name: test
		image: google/bazel
		commands:
			- bazel test`,
				},
			}
			outputResults = append(outputResults, testResult)
		}
		if foundBazelBuild {
			buildResult := types.Scanlet{
				Name:           BuildCheck,
				ScannerFamily:  Name,
				Description:    "run bazel build",
				OutputRenderer: outputter.DroneBuildAnalysis,
				Spec: buildanalysis.OutputFields{
					RawYaml: `  - name: build
		image: google/bazel
		commands:
					- bazel build`,
				},
			}
			outputResults = append(outputResults, buildResult)
		}
		if foundMavenTest {
			buildResult := types.Scanlet{
				Name:           BuildCheck,
				ScannerFamily:  Name,
				Description:    "run maven test",
				OutputRenderer: outputter.DroneBuildAnalysis,
				Spec: buildanalysis.OutputFields{
					RawYaml: `  - name: test
			image: maven
			commands:
						- mvn test`,
				},
			}
			outputResults = append(outputResults, buildResult)
		}
		if foundMavenBuild {
			buildResult := types.Scanlet{
				Name:           BuildCheck,
				ScannerFamily:  Name,
				Description:    "run maven build",
				OutputRenderer: outputter.DroneBuildAnalysis,
				Spec: buildanalysis.OutputFields{
					RawYaml: `  - name: build
			image: maven
			commands:
						- mvn clean install`,
				},
			}
			outputResults = append(outputResults, buildResult)
		}
		if foundGradleTest {
			buildResult := types.Scanlet{
				Name:           BuildCheck,
				ScannerFamily:  Name,
				Description:    "run gradle test",
				OutputRenderer: outputter.DroneBuildAnalysis,
				Spec: buildanalysis.OutputFields{
					RawYaml: `  - name: test
			image: gradle/gradle
			commands:
						- ./gradlew test`,
				},
			}
			outputResults = append(outputResults, buildResult)
		}
		if foundGradleBuild {
			buildResult := types.Scanlet{
				Name:           BuildCheck,
				ScannerFamily:  Name,
				Description:    "run gradle build",
				OutputRenderer: outputter.DroneBuildAnalysis,
				Spec: buildanalysis.OutputFields{
					RawYaml: `  - name: build
			image: gradle/gradle
			commands:
						- ./gradlew clean build`,
				},
			}
			outputResults = append(outputResults, buildResult)
		}
		if hasAndroid && !foundAndroidCommands {
			buildResult := types.Scanlet{
				Name:           BuildCheck,
				ScannerFamily:  Name,
				Description:    "run android tests and builds with the android sdk",
				OutputRenderer: outputter.DroneBuildAnalysis,
				Spec: buildanalysis.OutputFields{
					RawYaml: `  - name: build
			image: android/sdk
			commands:
						- sdkmanager --update --no-prompt --all`,
				},
			}
			outputResults = append(outputResults, buildResult)
		}
	}

	return outputResults, err
}
