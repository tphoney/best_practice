# Best_practice

A plugin/cli tool/container/library for automating best practice in a code repository. For example:

- Build a drone yaml file based on the programming language technologies in the repository
- Analyse your drone build file and give recommendations based on languages used in the repository
- Recommend other Harness products based on your repository

A scanner will check a language for build, lint, testing capabilities and language specific features. EG Android builds in Java. We have the following language specific scanners:

- Docker scanner
- Drone scanner
- Golang scanner
- Java scanner
- Javascript scanner
- Ruby scanner

And the following output formats:

- Best practice for existing Drone builds
- Build file creation, either Drone or CIE (*.new file if you have an existing file)
- Harness product recommendations

Example output:

<img width="447" alt="image" src="https://user-images.githubusercontent.com/10402706/175973905-0eaa76f9-5d9e-4f4e-8305-03c1021169b0.png">

## Usage

There are 4 ways to use this project:

### Quick start (docker container)

To run the best_practice tool against the current working directory.

```bash
docker pull tphoney/best_practice
docker run -it --rm -v $(pwd):/plugin -e PLUGIN_WORKING_DIRECTORY=/plugin tphoney/best_practice
```

It will create a drone build file, give best practice (if a drone file exists) and harness product recommendations.

To execute the newly created drone build file.

```bash
# install drone-cli if necessary
brew install drone-cli
# execute the drone build
drone exec .drone.yml
# your build should run !
```

**NB if there is an existing `.drone.yml` file, it will create one called `.drone.yml.new`**

### Using it as a cli tool

Download the Binaries from the release section. Then, you can use it as a cli tool.

```bash
./best-practice
```

Execute the newly created drone build file

```bash
# install drone-cli if necessary
brew install drone-cli
# execute the drone build
drone exec .drone.yml
```

### Using it in your drone build

Below is an example `.drone.yml` that uses this plugin.

```yaml
kind: pipeline
name: default

steps:
- name: run tphoney/best_practice plugin
  image: tphoney/best_practice
  pull: if-not-exists
```

### Using it as a library

Select your scanners and pass it through to the output formatters:

```go
# set the working directory to the root of your project
workingDirectory, err := os.Getwd()
# set your scanners, this uses all of the scanners by default
requestedScanners = scanner.ListScannersNames()
# set your output formatters, this uses all of the output formatters by default
requestesOutputFormatters = output.ListOutputFormattersNames()
# run the scanners
scanResults, scanErr := scanner.RunScanners(ctx, requestedScanners, requestesOutputFormatters)
# run the output formatters
outputErr := outputter.RunOutput(ctx, outputters, scanResults)
```

## Developer notes

### Building

Build the plugin binaries:

```text
scripts/build.sh
```

Build the plugin image:

```text
docker build -t tphoney/best_practice -f docker/Dockerfile .
```
