package main

// OutputConfig represents output format and destination configuration.
type OutputConfig struct {
	Directory string  `yaml:"directory"`          // Output directory for generated images
	Format    string  `yaml:"format"`             // Output image format (png, jpg)
	Filename  *string `yaml:"filename,omitempty"` // Custom filename template (optional, default: "ogp.{format}")
}

// BackgroundConfig represents background image and color configuration.
type BackgroundConfig struct {
	Image *string `yaml:"image,omitempty"` // Path to background image (optional)
	Color string  `yaml:"color"`           // Background color (hex) when image is not specified
}

// LineBreakingConfig represents Japanese line breaking rules configuration.
type LineBreakingConfig struct {
	StartProhibited string `yaml:"start_prohibited"` // Characters that cannot start a line
	EndProhibited   string `yaml:"end_prohibited"`   // Characters that cannot end a line
}

// LineBreakingOverride represents overridable Japanese line breaking rules.
type LineBreakingOverride struct {
	StartProhibited *string `yaml:"start_prohibited,omitempty"` // Characters that cannot start a line
	EndProhibited   *string `yaml:"end_prohibited,omitempty"`   // Characters that cannot end a line
}

// PlacementConfig represents positioning information for overlays.
type PlacementConfig struct {
	X      int `yaml:"x"`
	Y      int `yaml:"y"`
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
}

// TextAreaConfig represents the text rendering area coordinates.
type TextAreaConfig struct {
	X      *int `yaml:"x,omitempty"`
	Y      *int `yaml:"y,omitempty"`
	Width  *int `yaml:"width,omitempty"`
	Height *int `yaml:"height,omitempty"`
}

// OverlayConfigBase represents common overlay configuration.
type OverlayConfigBase struct {
	Image     *string          `yaml:"image,omitempty"`     // Path to image file
	Placement *PlacementConfig `yaml:"placement,omitempty"` // Image positioning
	Fit       *string          `yaml:"fit,omitempty"`       // Fit method ("cover", "contain", "fill", "none")
	Opacity   *float64         `yaml:"opacity,omitempty"`   // Image opacity (0.0-1.0)
}

// MainOverlayConfig represents overlay configuration in main config.
type MainOverlayConfig struct {
	Visible   bool            `yaml:"visible"`         // Whether to render this overlay (default: true)
	Image     *string         `yaml:"image,omitempty"` // Path to image file
	Placement PlacementConfig `yaml:"placement"`       // Image positioning
	Fit       string          `yaml:"fit"`             // Fit method ("cover", "contain", "fill", "none")
	Opacity   float64         `yaml:"opacity"`         // Image opacity (0.0-1.0)
}

// ArticleOverlayConfig represents overlay configuration in front matter.
type ArticleOverlayConfig struct {
	Visible   *bool            `yaml:"visible,omitempty"`   // Whether to render this overlay (default: true)
	Image     *string          `yaml:"image,omitempty"`     // Path to image file
	Placement *PlacementConfig `yaml:"placement,omitempty"` // Image positioning
	Fit       *string          `yaml:"fit,omitempty"`       // Fit method ("cover", "contain", "fill", "none")
	Opacity   *float64         `yaml:"opacity,omitempty"`   // Image opacity (0.0-1.0)
}

// BackgroundOverride represents background configuration overrides in front matter.
type BackgroundOverride struct {
	Image *string `yaml:"image,omitempty"` // Path to background image (relative to article directory)
	Color *string `yaml:"color,omitempty"` // Background color (hex)
}

// OutputOverride represents output configuration overrides in front matter.
type OutputOverride struct {
	Filename *string `yaml:"filename,omitempty"` // Custom filename template (optional)
}
