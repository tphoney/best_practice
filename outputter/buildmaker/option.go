package buildmaker

type Option func(*outputterConfig)

func WithStdOutput(stdOutput bool) Option {
	return func(p *outputterConfig) {
		p.stdOutput = stdOutput
	}
}

func WithOutputToFile(i bool) Option {
	return func(p *outputterConfig) {
		p.outputToFile = i
	}
}

func WithWorkingDirectory(i string) Option {
	return func(p *outputterConfig) {
		p.workingDirectory = i
	}
}

func WithDroneOutput(i bool) Option {
	return func(p *outputterConfig) {
		p.outputDrone = i
	}
}

func WithCIEOutput(i bool) Option {
	return func(p *outputterConfig) {
		p.outputCIE = i
	}
}
