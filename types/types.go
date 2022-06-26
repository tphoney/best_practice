package types

const (
	// output plugin names
	OutputNameDroneBuild        = "dronebuild"
	OutputBestPractice          = "bestpractice"
	OutputProductRecommendation = "productrecommendation"
)

type (
	ScanType struct {
		FamilyName string
		ScanTypes  []string
	}

	ScanInput struct {
		RequestedOutputs []string `json:"requested_outputs" yaml:"requested_outputs"`
		RunAll           bool     `json:"run_all" yaml:"run_all"`
		ScansToRun       []string `json:"scans_to_run" yaml:"scans_to_run"`
	}

	ScanResult struct {
		FamilyName     string    `json:"scanlet_family" yaml:"scanlet_family"`
		Match          bool      `json:"match" yaml:"match"`
		ScansRun       []string  `json:"scans_run" yaml:"scans_run"`
		ScanletResults []Scanlet `json:"outputs" yaml:"outputs"`
	}

	Scanlet struct {
		Name           string      `json:"scanlet_family" yaml:"scanlet_family"`
		HumanReasoning string      `json:"human_output" yaml:"human_output"`
		OutputRender   string      `json:"output_render" yaml:"output_render"`
		Spec           interface{} `json:"spec" yaml:"spec"`
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
