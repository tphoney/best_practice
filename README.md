# Best_practice

A plugin/cli tool/container/library for automating best practice in a code repository. For example:

- Build a drone yaml file based on the technologies in the repo
- Give general programming language specific recommendations
- Recommend other Harness products based on your project

It has the following scanners:

- Drone scanner, analyses your build file
- golang scanner, for best practice
- Java scanner, for best practice
- Javascript scanner, for best practice

And the following output formats:

- Best practice report
- Drone build file creation
- Harness product recommendations.

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

You can use a container locally. This will run it against your current working directory.

```bash
docker pull tphoney/best_practice
docker run -it --rm -v $(pwd):/plugin -e PLUGIN_WORKING_DIRECTORY=/plugin tphoney/best_practice
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
