package dronebuild

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
