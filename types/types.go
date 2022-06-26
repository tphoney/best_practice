package types

import "context"

type (
	Scanner interface {
		Name() string
		Description() string
		AvailableChecks() []string
		Scan(ctx context.Context, RequestedOutputs []string) (scanResults []Scanlet, err error)
	}

	Outputter interface {
		Name() string
		Description() string
		Output(ctx context.Context, scanResults []Scanlet) error
	}

	Scanlet struct {
		Name           string      `json:"name" yaml:"name"`
		ScannerFamily  string      `json:"scanner_family" yaml:"scanner_family"`
		Description    string      `json:"description" yaml:"description"`
		OutputRenderer string      `json:"output_renderer" yaml:"output_renderer"`
		Spec           interface{} `json:"spec" yaml:"spec"`
	}
)
