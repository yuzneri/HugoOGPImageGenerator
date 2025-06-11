package main

import (
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config represents the main configuration structure for OGP image generation.
// It contains all settings for fonts, text rendering, image processing, and line breaking.
type Config struct {
	// Background configuration
	Background struct {
		Image *string `yaml:"image,omitempty"` // Path to background image (optional)
		Color string  `yaml:"color"`           // Background color (hex) when image is not specified
	} `yaml:"background"`

	// Output configuration
	Output struct {
		Directory string  `yaml:"directory"`          // Output directory for generated images
		Format    string  `yaml:"format"`             // Output image format (png, jpg)
		Filename  *string `yaml:"filename,omitempty"` // Custom filename template (optional, default: "ogp.{format}")
	} `yaml:"output"`

	// Text rendering configuration
	Text struct {
		// Content configuration
		Content *string `yaml:"content,omitempty"` // Default content template (optional)
		// Font configuration
		Font string  `yaml:"font"` // Path to font file
		Size float64 `yaml:"size"` // Default font size
		// Text color configuration (supports hex color codes)
		Color string `yaml:"color"` // Hex color code (e.g., "#FF00FF", "#ff00ff80")
		// Text rendering area coordinates
		Area          TextArea `yaml:"area"`
		BlockPosition string   `yaml:"block_position"` // Text block position in area (e.g., "middle-center")
		LineAlignment string   `yaml:"line_alignment"` // Individual line alignment within block ("left", "center", "right")
		Overflow      string   `yaml:"overflow"`       // Overflow handling ("shrink" or "clip")
		MinSize       float64  `yaml:"min_size"`       // Minimum font size for shrink mode
		LineHeight    float64  `yaml:"line_height"`    // Line height multiplier
		LetterSpacing int      `yaml:"letter_spacing"` // Letter spacing in pixels
		// Japanese line breaking rules configuration
		LineBreaking struct {
			StartProhibited string `yaml:"start_prohibited"` // Characters that cannot start a line
			EndProhibited   string `yaml:"end_prohibited"`   // Characters that cannot end a line
		} `yaml:"line_breaking"`
	} `yaml:"text"`

	// Default overlay configuration (optional)
	Overlay *struct {
		Image     *string `yaml:"image,omitempty"` // Path to image file (relative to project root)
		Placement *struct {
			X      *int `yaml:"x,omitempty"`      // X position
			Y      *int `yaml:"y,omitempty"`      // Y position
			Width  *int `yaml:"width,omitempty"`  // Image width
			Height *int `yaml:"height,omitempty"` // Image height
		} `yaml:"placement,omitempty"`
		Fit     *string  `yaml:"fit,omitempty"`     // Fit method ("cover", "contain", "fill", "none")
		Opacity *float64 `yaml:"opacity,omitempty"` // Image opacity (0.0-1.0)
	} `yaml:"overlay,omitempty"`
}

// parseHexColor parses hex color codes like "#FF00FF" or "#ff00ff80"
func parseHexColor(hex string) (color.RGBA, error) {
	// Remove # prefix if present
	hex = strings.TrimPrefix(hex, "#")

	// Support both 6-char (#RRGGBB) and 8-char (#RRGGBBAA) formats
	switch len(hex) {
	case 6:
		// Parse RGB, default alpha to 255
		val, err := strconv.ParseUint(hex, 16, 32)
		if err != nil {
			return color.RGBA{}, err
		}
		return color.RGBA{
			R: uint8(val >> 16),
			G: uint8(val >> 8),
			B: uint8(val),
			A: 255,
		}, nil
	case 8:
		// Parse RGBA
		val, err := strconv.ParseUint(hex, 16, 32)
		if err != nil {
			return color.RGBA{}, err
		}
		return color.RGBA{
			R: uint8(val >> 24),
			G: uint8(val >> 16),
			B: uint8(val >> 8),
			A: uint8(val),
		}, nil
	default:
		return color.RGBA{}, fmt.Errorf("invalid hex color format: %s (expected 6 or 8 characters)", hex)
	}
}

// OGPFrontMatter represents OGP-specific settings in article front matter.
// All fields are optional and override the corresponding config values.
type OGPFrontMatter struct {
	Text *struct {
		Content *string  `yaml:"content,omitempty"`
		Font    *string  `yaml:"font,omitempty"`
		Size    *float64 `yaml:"size,omitempty"`
		Color   *string  `yaml:"color,omitempty"` // Hex color code
		Area    *struct {
			X      *int `yaml:"x,omitempty"`
			Y      *int `yaml:"y,omitempty"`
			Width  *int `yaml:"width,omitempty"`
			Height *int `yaml:"height,omitempty"`
		} `yaml:"area,omitempty"`
		BlockPosition *string  `yaml:"block_position,omitempty"`
		LineAlignment *string  `yaml:"line_alignment,omitempty"`
		Overflow      *string  `yaml:"overflow,omitempty"`
		MinSize       *float64 `yaml:"min_size,omitempty"`
		LineHeight    *float64 `yaml:"line_height,omitempty"`
		LetterSpacing *int     `yaml:"letter_spacing,omitempty"`
		LineBreaking  *struct {
			StartProhibited *string `yaml:"start_prohibited,omitempty"` // Characters that cannot start a line
			EndProhibited   *string `yaml:"end_prohibited,omitempty"`   // Characters that cannot end a line
		} `yaml:"line_breaking,omitempty"`
	} `yaml:"text,omitempty"`

	// Background settings
	Background *struct {
		Image *string `yaml:"image,omitempty"` // Path to background image (relative to article directory)
		Color *string `yaml:"color,omitempty"` // Background color (hex)
	} `yaml:"background,omitempty"`

	// Overlay image composition settings
	Overlay *struct {
		Image     *string `yaml:"image,omitempty"` // Path to image file (relative to article directory)
		Placement *struct {
			X      *int `yaml:"x,omitempty"`      // X position
			Y      *int `yaml:"y,omitempty"`      // Y position
			Width  *int `yaml:"width,omitempty"`  // Image width
			Height *int `yaml:"height,omitempty"` // Image height
		} `yaml:"placement,omitempty"`
		Fit     *string  `yaml:"fit,omitempty"`     // Fit method ("cover", "contain", "fill", "none")
		Opacity *float64 `yaml:"opacity,omitempty"` // Image opacity (0.0-1.0)
	} `yaml:"overlay,omitempty"`

	// Output settings
	Output *struct {
		Filename *string `yaml:"filename,omitempty"` // Custom filename template (optional)
	} `yaml:"output,omitempty"`
}

// getDefaultConfig returns a config with sensible defaults for OGP generation.
func getDefaultConfig() *Config {
	config := &Config{}

	// Background defaults
	config.Background.Color = "#FFFFFF"

	// Output defaults
	config.Output.Directory = "public"
	config.Output.Format = "png"

	// Text defaults
	config.Text.Font = "" // システムのフォントを自動選択
	config.Text.Size = 64
	config.Text.Color = "#000000"
	config.Text.Area.X = 100
	config.Text.Area.Y = 100
	config.Text.Area.Width = 1000
	config.Text.Area.Height = 430
	config.Text.BlockPosition = "middle-center"
	config.Text.LineAlignment = "left"
	config.Text.Overflow = "shrink"
	config.Text.MinSize = 12.0
	config.Text.LineHeight = 1.2
	config.Text.LetterSpacing = 1

	// Japanese line breaking defaults
	config.Text.LineBreaking.StartProhibited = ".)}]>!?、。，．！？)）］｝〉》」』ー～ぁぃぅぇぉっゃゅょゎァィゥェォッャュョヮヵヶ々"
	config.Text.LineBreaking.EndProhibited = "({[<（［｛〈《「『"

	return config
}

// loadConfig reads and parses a YAML configuration file.
// If the file doesn't exist or fields are missing, defaults are applied.
func loadConfig(configPath string) (*Config, error) {
	// Start with default configuration
	config := getDefaultConfig()

	// Try to read config file
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		// If config file doesn't exist, return defaults
		return config, nil
	}

	// Parse the config file and merge with defaults
	err = yaml.Unmarshal(configBytes, config)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}

// buildProhibitedMaps creates lookup maps for Japanese line breaking rules.
// Returns two maps: one for characters that cannot start a line,
// and one for characters that cannot end a line.
func buildProhibitedMaps(config *Config) (map[rune]bool, map[rune]bool) {
	startProhibited := make(map[rune]bool)
	endProhibited := make(map[rune]bool)

	// Build map for characters that cannot start a line (行頭禁則文字)
	for _, r := range []rune(config.Text.LineBreaking.StartProhibited) {
		startProhibited[r] = true
	}

	// Build map for characters that cannot end a line (行末禁則文字)
	for _, r := range []rune(config.Text.LineBreaking.EndProhibited) {
		endProhibited[r] = true
	}

	return startProhibited, endProhibited
}

// applyFrontMatterOverrides creates a new config with front matter overrides applied.
// This allows per-article customization of OGP generation settings.
func applyFrontMatterOverrides(config *Config, ogpFM *OGPFrontMatter) *Config {
	merger := NewConfigMerger()
	return merger.MergeConfigs(config, ogpFM)
}
