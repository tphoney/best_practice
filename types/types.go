package types

const (
	// output plugin names
	OutputNameDroneBuild        = "dronebuild"
	OutputBestPractice          = "bestpractice"
	OutputProductRecommendation = "productrecommendation"
)

type (
	Scan struct {
		FamilyName         string    `json:"scanlet_family" yaml:"scanlet_family"`
		Match              bool      `json:"match" yaml:"match"`
		OutputCapabilities []string  `json:"output_capabilities" yaml:"output_capabilities"`
		Outputs            []Scanlet `json:"outputs" yaml:"outputs"`
	}
	Scanlet struct {
		Name           string      `json:"scanlet_family" yaml:"scanlet_family"`
		HumanReasoning string      `json:"human_output" yaml:"human_output"`
		Enabled        bool        `json:"enabled" yaml:"enabled"`
		OutputRender   string      `json:"output_render" yaml:"output_render"`
		Spec           interface{} `json:"spec" yaml:"spec"`
	}

	DroneBuildOutput struct {
		RawYaml string `json:"raw_yaml" yaml:"raw_yaml"`
	}

	BestPracticeOutput struct {
		Command string `json:"command" yaml:"command"`
	}

	ProductRecommendationOutput struct {
		Stdout       bool   `json:"stdout" yaml:"stdout"`
		FileLocation string `json:"file_location" yaml:"file_location"`
	}
)
