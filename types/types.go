package types

import "context"

type (
	Scanner interface {
		Scan(ctx context.Context, RequestedOutputs []string) (scanResults []Scanlet, err error)
		Name() string
	}

	Outputter interface {
		Output(ctx context.Context, scanResults []Scanlet) error
		Name() string
	}

	ScanResult struct {
		FamilyName     string    `json:"scanlet_family" yaml:"scanlet_family"`
		Match          bool      `json:"match" yaml:"match"`
		ScansRun       []string  `json:"scans_run" yaml:"scans_run"`
		ScanletResults []Scanlet `json:"outputs" yaml:"outputs"`
	}

	Scanlet struct {
		Name           string      `json:"name" yaml:"name"`
		Description    string      `json:"description" yaml:"description"`
		OutputRenderer string      `json:"output_renderer" yaml:"output_renderer"`
		Spec           interface{} `json:"spec" yaml:"spec"`
	}

	OutputterConfiguration struct {
		Name        string      `json:"name" yaml:"name"`
		Description string      `json:"description" yaml:"description"`
		Spec        interface{} `json:"spec" yaml:"spec"`
	}

	DroneBuildOutput struct {
		RawYaml string `json:"raw_yaml" yaml:"raw_yaml"`
	}

	BestPracticeOutput struct {
		Command string `json:"command" yaml:"command"`
		Url     string `json:"url" yaml:"url"`
	}

	ProductRecommendationOutput struct {
		Stdout       bool   `json:"stdout" yaml:"stdout"`
		FileLocation string `json:"file_location" yaml:"file_location"`
	}
)
