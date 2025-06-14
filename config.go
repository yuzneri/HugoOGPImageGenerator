package main

import (
	"fmt"
	"image/color"
	"os"
	"strconv"
	"strings"

	"gopkg.in/yaml.v3"
)

// TextConfig represents configuration for a single text element (title or description).
// It contains all settings for fonts, text rendering, and line breaking for that text.
// TextConfig represents a complete text configuration with all fields having values (runtime use)
type TextConfig struct {
	// Rendering control
	Visible bool `yaml:"visible"` // Whether to render this text element
	// Content configuration
	Content *string `yaml:"content"` // Content template (nil means use template)
	// Font configuration
	Font *string `yaml:"font"` // Path to font file (nil means auto-detect)
	Size float64 `yaml:"size"` // Font size
	// Text color configuration
	Color string `yaml:"color"` // Hex color code
	// Text rendering area coordinates
	Area          TextArea `yaml:"area"`
	BlockPosition string   `yaml:"block_position"` // Text block position in area
	LineAlignment string   `yaml:"line_alignment"` // Individual line alignment within block
	Overflow      string   `yaml:"overflow"`       // Overflow handling ("shrink" or "clip")
	MinSize       float64  `yaml:"min_size"`       // Minimum font size for shrink mode
	LineHeight    float64  `yaml:"line_height"`    // Line height multiplier
	LetterSpacing int      `yaml:"letter_spacing"` // Letter spacing in pixels
	// Japanese line breaking rules configuration
	LineBreaking LineBreakingConfig `yaml:"line_breaking"`
}

// Legacy type aliases for backwards compatibility
type PlacementInfo = PlacementConfig
type ConfigOverlay = MainOverlayConfig
type FrontMatterOverlay = ArticleOverlayConfig

// GetImage implements OverlaySettings interface
func (c MainOverlayConfig) GetImage() *string {
	return c.Image
}

// GetPlacement implements OverlaySettings interface
func (c MainOverlayConfig) GetPlacement() *PlacementConfig {
	return &c.Placement
}

// GetFit implements OverlaySettings interface
func (c MainOverlayConfig) GetFit() *string {
	return &c.Fit
}

// GetOpacity implements OverlaySettings interface
func (c MainOverlayConfig) GetOpacity() *float64 {
	return &c.Opacity
}

// GetImage implements OverlaySettings interface
func (f *ArticleOverlayConfig) GetImage() *string {
	return f.Image
}

// GetPlacement implements OverlaySettings interface
// NOTE: ArticleOverlayConfig should not be used for rendering, only for configuration merging
func (f *ArticleOverlayConfig) GetPlacement() *PlacementConfig {
	// This method should not be called in normal operation
	// ArticleOverlayConfig is for configuration merging, not for rendering
	if f.Placement == nil {
		return nil
	}
	// Convert PlacementSettings to PlacementConfig for interface compatibility
	// IMPORTANT: Preserve all explicitly specified values, including zero
	return &PlacementConfig{
		X:      valueOrDefault(f.Placement.X, 0),
		Y:      valueOrDefault(f.Placement.Y, 0),
		Width:  f.Placement.Width,
		Height: f.Placement.Height,
	}
}

// valueOrDefault returns the value pointed to by ptr, or defaultValue if ptr is nil
// IMPORTANT: If ptr is non-nil, the pointed value is used even if it's zero
func valueOrDefault(ptr *int, defaultValue int) int {
	if ptr == nil {
		return defaultValue
	}
	return *ptr // 明示的に指定された値（0を含む）をそのまま使用
}

// GetFit implements OverlaySettings interface
func (f *ArticleOverlayConfig) GetFit() *string {
	return f.Fit
}

// GetOpacity implements OverlaySettings interface
func (f *ArticleOverlayConfig) GetOpacity() *float64 {
	return f.Opacity
}

// Config represents the main configuration structure for OGP image generation.
// It contains all settings for fonts, text rendering, image processing, and line breaking.
type Config struct {
	// Background configuration
	Background BackgroundConfig `yaml:"background"`

	// Output configuration
	Output OutputConfig `yaml:"output"`

	// Text rendering configurations
	Title       TextConfig `yaml:"title"`       // Title text configuration
	Description TextConfig `yaml:"description"` // Description text configuration

	// Default overlay configuration
	Overlay MainOverlayConfig `yaml:"overlay"`
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
		return color.RGBA{}, NewValidationError(fmt.Sprintf("invalid hex color format: %s (expected 6 or 8 characters)", hex))
	}
}

// TextConfigOverride represents overrides for a text configuration in front matter.
// All fields are optional pointers to allow partial overrides.
type TextConfigOverride struct {
	Visible       *bool                 `yaml:"visible,omitempty"`
	Content       *string               `yaml:"content,omitempty"`
	Font          *string               `yaml:"font,omitempty"`
	Size          *float64              `yaml:"size,omitempty"`
	Color         *string               `yaml:"color,omitempty"` // Hex color code
	Area          *TextAreaConfig       `yaml:"area,omitempty"`
	BlockPosition *string               `yaml:"block_position,omitempty"`
	LineAlignment *string               `yaml:"line_alignment,omitempty"`
	Overflow      *string               `yaml:"overflow,omitempty"`
	MinSize       *float64              `yaml:"min_size,omitempty"`
	LineHeight    *float64              `yaml:"line_height,omitempty"`
	LetterSpacing *int                  `yaml:"letter_spacing,omitempty"`
	LineBreaking  *LineBreakingOverride `yaml:"line_breaking,omitempty"`
}

// OGPFrontMatter represents OGP-specific settings in article front matter.
// All fields are optional and override the corresponding config values.
type OGPFrontMatter struct {
	// Text configurations
	Title       *TextConfigOverride `yaml:"title,omitempty"`       // Title text overrides
	Description *TextConfigOverride `yaml:"description,omitempty"` // Description text overrides

	// Background settings
	Background *BackgroundOverride `yaml:"background,omitempty"`

	// Overlay image composition settings
	Overlay *ArticleOverlayConfig `yaml:"overlay,omitempty"`

	// Output settings
	Output *OutputOverride `yaml:"output,omitempty"`
}

// getDefaultConfig returns a config with sensible defaults for OGP generation.
func getDefaultConfig() *Config {
	config := &Config{}

	setDefaultBackground(config)
	setDefaultOutput(config)
	setDefaultTitle(config)
	setDefaultDescription(config)
	setDefaultOverlay(config)

	return config
}

// setDefaultBackground configures default background settings
func setDefaultBackground(config *Config) {
	config.Background.Color = DefaultBackgroundColor
}

// setDefaultOutput configures default output settings
func setDefaultOutput(config *Config) {
	config.Output.Directory = DefaultOutputDirectory
	config.Output.Format = FormatPNG
	config.Output.Filename = "ogp"
}

// setDefaultTitle configures default title settings
func setDefaultTitle(config *Config) {
	config.Title.Visible = true
	config.Title.Font = nil
	config.Title.Size = DefaultTitleFontSize
	config.Title.Color = DefaultTitleColor
	config.Title.Area.X = DefaultTitleAreaX
	config.Title.Area.Y = DefaultTitleAreaY
	config.Title.Area.Width = DefaultTitleAreaWidth
	config.Title.Area.Height = DefaultTitleAreaHeight
	config.Title.BlockPosition = DefaultTitleBlockPosition
	config.Title.LineAlignment = DefaultTitleLineAlignment
	config.Title.Overflow = OverflowShrink
	config.Title.MinSize = DefaultTitleMinSize
	config.Title.LineHeight = DefaultLineHeight
	config.Title.LetterSpacing = DefaultTitleLetterSpacing
	config.Title.LineBreaking.StartProhibited = DefaultStartProhibitedChars
	config.Title.LineBreaking.EndProhibited = DefaultEndProhibitedChars
}

// setDefaultDescription configures default description settings
func setDefaultDescription(config *Config) {
	config.Description.Visible = false
	config.Description.Font = nil
	config.Description.Size = DefaultDescriptionFontSize
	config.Description.Color = DefaultDescriptionColor
	config.Description.Area.X = DefaultDescriptionAreaX
	config.Description.Area.Y = DefaultDescriptionAreaY
	config.Description.Area.Width = DefaultDescriptionAreaWidth
	config.Description.Area.Height = DefaultDescriptionAreaHeight
	config.Description.BlockPosition = DefaultDescriptionBlockPosition
	config.Description.LineAlignment = DefaultDescriptionLineAlignment
	config.Description.Overflow = OverflowClip
	config.Description.MinSize = DefaultDescriptionMinSize
	config.Description.LineHeight = DefaultLineHeight
	config.Description.LetterSpacing = DefaultDescriptionLetterSpacing
	config.Description.LineBreaking.StartProhibited = DefaultStartProhibitedChars
	config.Description.LineBreaking.EndProhibited = DefaultEndProhibitedChars
}

// setDefaultOverlay configures default overlay settings
func setDefaultOverlay(config *Config) {
	config.Overlay.Visible = DefaultOverlayVisible
	config.Overlay.Placement.X = DefaultOverlayX
	config.Overlay.Placement.Y = DefaultOverlayY
	config.Overlay.Fit = DefaultOverlayFit
	config.Overlay.Opacity = DefaultOverlayOpacity
}

// loadConfig reads and parses a YAML configuration file.
// If the file doesn't exist or fields are missing, defaults are applied.

// loadConfigSettings loads ConfigSettings from a file.
// Returns nil if the file doesn't exist (not an error).
func loadConfigSettings(configPath string) (*ConfigSettings, error) {
	// Try to read config file
	configBytes, err := os.ReadFile(configPath)
	if err != nil {
		// If config file doesn't exist, return nil (not an error)
		return nil, nil
	}

	// Parse the config file
	var settings ConfigSettings
	err = yaml.Unmarshal(configBytes, &settings)
	if err != nil {
		return nil, NewConfigError("failed to unmarshal config", err)
	}

	return &settings, nil
}

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
		return nil, NewConfigError("failed to unmarshal config", err)
	}

	return config, nil
}

// buildProhibitedMaps creates lookup maps for Japanese line breaking rules.
// Returns two maps: one for characters that cannot start a line,
// and one for characters that cannot end a line.
func buildProhibitedMaps(textConfig *TextConfig) (map[rune]bool, map[rune]bool) {
	startProhibited := make(map[rune]bool)
	endProhibited := make(map[rune]bool)

	// Build map for characters that cannot start a line (行頭禁則文字)
	for _, r := range []rune(textConfig.LineBreaking.StartProhibited) {
		startProhibited[r] = true
	}

	// Build map for characters that cannot end a line (行末禁則文字)
	for _, r := range []rune(textConfig.LineBreaking.EndProhibited) {
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
