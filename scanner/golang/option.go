package golang

import "golang.org/x/exp/slices"

type Option func(*scannerConfig)

func WithChecksToRun(i []string) Option {
	return func(p *scannerConfig) {
		if len(i) > 0 {
			validChecks := []string{}
			// only add valid checks
			for _, check := range i {
				if slices.Contains(p.AvailableChecks(), check) {
					validChecks = append(validChecks, check)
				}
			}
			p.runAll = false
			p.checksToRun = validChecks
		} else {
			p.runAll = true
		}
	}
}
