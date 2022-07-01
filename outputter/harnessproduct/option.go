package harnessproduct

type Option func(*outputterConfig)

func WithWorkingDirectory(i string) Option {
	return func(p *outputterConfig) {
		p.workingDirectory = i
	}
}
