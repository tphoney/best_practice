# Best_practice

A plugin/cli tool/container/library for automating best practice in a code repository.
It has the following scanners:

- golang, scanner for golang best practice.

And the following output formats:

- Drone build file creation
- best practice report

Example output:

<img width="447" alt="image" src="https://user-images.githubusercontent.com/10402706/175973905-0eaa76f9-5d9e-4f4e-8305-03c1021169b0.png">

## Usage

There are 4 ways to use this project:

### Using it as a cli tool

Download the Binaries from the release section. Then, you can use it as a cli tool.

```bash
./best-practice 
```

### Using the container locally

```bash
docker pull tphoney/best_practice
docker run -it --rm -v $(pwd):/best_practice tphoney/best_practice
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

Have a look at the calls in `plugin\plugin.go

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
