package bestpractice

type Option func(*outputterConfig)

func WithStdOutput(stdOutput bool) Option {
	return func(p *outputterConfig) {
		p.stdOutput = stdOutput
	}
}

func WithOutputToFile(i string) Option {
	return func(p *outputterConfig) {
		p.outputToFile = i
	}
}
