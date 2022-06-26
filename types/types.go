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

	// ScanInput struct {
	// 	FamilyName       string
	// 	RequestedOutputs []string `json:"requested_outputs" yaml:"requested_outputs"`
	// 	RunAll           bool     `json:"run_all" yaml:"run_all"`
	// 	ScanletsToRun    []string `json:"scanlets_to_run" yaml:"scanlets_to_run"`
	// }

	ScanResult struct {
		FamilyName     string    `json:"scanlet_family" yaml:"scanlet_family"`
		Match          bool      `json:"match" yaml:"match"`
		ScansRun       []string  `json:"scans_run" yaml:"scans_run"`
		ScanletResults []Scanlet `json:"outputs" yaml:"outputs"`
	}

	Scanlet struct {
		Name           string      `json:"name" yaml:"name"`
		HumanReasoning string      `json:"human_output" yaml:"human_output"`
		OutputRender   string      `json:"output_render" yaml:"output_render"`
		Spec           interface{} `json:"spec" yaml:"spec"`
	}

	RequestedOutputs struct {
		RequestedOutputs []string `json:"requested_outputs" yaml:"requested_outputs"`
	}

	OutputType struct {
		Name        string `json:"name" yaml:"name"`
		Description string `json:"description" yaml:"description"`
		ToFile      string `json:"to_file" yaml:"to_file"`
		ToStdout    bool   `json:"to_stdout" yaml:"to_stdout"`
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
