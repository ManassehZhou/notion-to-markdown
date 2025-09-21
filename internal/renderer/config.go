package renderer

import (
	"fmt"
	"log/slog"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// renderTemplate renders a template string with the given data
func renderTemplate(template string, data map[string]string) string {
	result := template
	for key, value := range data {
		placeholder := "{{." + key + "}}"
		result = strings.ReplaceAll(result, placeholder, value)
	}
	return result
}

// RenderConfig contains configuration for customizing non-standard Markdown rendering
type RenderConfig struct {
	// Math equations template
	MathTemplate string `yaml:"math_template" json:"math_template"`

	// Details/Toggle blocks template
	DetailsTemplate string `yaml:"details_template" json:"details_template"`

	// Video blocks template
	VideoTemplate string `yaml:"video_template" json:"video_template"`

	// PDF blocks template
	PDFTemplate string `yaml:"pdf_template" json:"pdf_template"`

	// Embed blocks template
	EmbedTemplate string `yaml:"embed_template" json:"embed_template"`

	// Callout blocks template
	CalloutTemplate string `yaml:"callout_template" json:"callout_template"`

	// File blocks template (for regular files)
	FileTemplate string `yaml:"file_template" json:"file_template"`
}

// DefaultRenderConfig returns the default configuration for Hugo shortcodes
func DefaultRenderConfig() *RenderConfig {
	return &RenderConfig{
		MathTemplate:    "{{< math >}}\n$$\n{{.Expression}}\n$$\n{{< /math >}}",
		DetailsTemplate: "{{< details summary=\"{{.Summary}}\">}}\n{{.Content}}\n{{< /details >}}",
		VideoTemplate:   "{{< video src=\"{{.URL}}\" >}}",
		PDFTemplate:     "{{< pdf src=\"{{.URL}}\" >}}",
		EmbedTemplate:   "{{< embed url=\"{{.URL}}\" >}}",
		CalloutTemplate: "> {{.Content}}",
		FileTemplate:    "[{{.Text}}]({{.URL}})",
	}
}

// LoadConfigFromYAML loads render configuration from a YAML file
func LoadConfigFromYAML(filepath string) (*RenderConfig, error) {
	// If file doesn't exist, return default config
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		slog.Info("Config file not found, using default configuration", "file", filepath)
		return DefaultRenderConfig(), nil
	}

	data, err := os.ReadFile(filepath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", filepath, err)
	}

	config := DefaultRenderConfig()
	if err := yaml.Unmarshal(data, config); err != nil {
		return nil, fmt.Errorf("failed to parse YAML config: %w", err)
	}

	slog.Info("Loaded configuration", "file", filepath)
	return config, nil
}

// LoadConfigWithFallback tries to load from file, falls back to default if not found
func LoadConfigWithFallback(filepath string) *RenderConfig {
	if config, err := LoadConfigFromYAML(filepath); err == nil {
		return config
	} else {
		slog.Warn("Failed to load config, using default", "error", err)
		return DefaultRenderConfig()
	}
}
