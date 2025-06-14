package main

// ConfigSettings represents configuration structure for reading from YAML files.
// All fields are pointers to distinguish between "not set" and "zero value".
type ConfigSettings struct {
	// Background configuration
	Background *BackgroundSettings `yaml:"background,omitempty"`

	// Output configuration
	Output *OutputSettings `yaml:"output,omitempty"`

	// Text rendering configurations
	Title       *TextSettings `yaml:"title,omitempty"`       // Title text configuration
	Description *TextSettings `yaml:"description,omitempty"` // Description text configuration

	// Default overlay configuration
	Overlay *OverlayConfigSettings `yaml:"overlay,omitempty"`
}

// BackgroundSettings represents background configuration for YAML reading.
type BackgroundSettings struct {
	Image *string `yaml:"image,omitempty"` // Path to background image
	Color *string `yaml:"color,omitempty"` // Background color (hex)
}

// OutputSettings represents output configuration for YAML reading.
type OutputSettings struct {
	Directory *string `yaml:"directory,omitempty"` // Output directory for generated images
	Format    *string `yaml:"format,omitempty"`    // Output image format (png, jpg)
	Filename  *string `yaml:"filename,omitempty"`  // Filename template
}

// TextSettings represents text configuration for YAML reading.
type TextSettings struct {
	// Rendering control
	Visible *bool `yaml:"visible,omitempty"` // Whether to render this text element

	// Content configuration
	Content *string `yaml:"content,omitempty"` // Content template

	// Font configuration
	Font *string  `yaml:"font,omitempty"` // Path to font file
	Size *float64 `yaml:"size,omitempty"` // Font size

	// Text color configuration
	Color *string `yaml:"color,omitempty"` // Hex color code

	// Text rendering area coordinates
	Area *TextAreaSettings `yaml:"area,omitempty"`

	// Text layout configuration
	BlockPosition *string `yaml:"block_position,omitempty"` // Text block position in area
	LineAlignment *string `yaml:"line_alignment,omitempty"` // Individual line alignment within block
	Overflow      *string `yaml:"overflow,omitempty"`       // Overflow handling ("shrink" or "clip")

	// Font sizing configuration
	MinSize *float64 `yaml:"min_size,omitempty"` // Minimum font size for shrink mode

	// Text spacing configuration
	LineHeight    *float64 `yaml:"line_height,omitempty"`    // Line height multiplier
	LetterSpacing *int     `yaml:"letter_spacing,omitempty"` // Letter spacing in pixels

	// Japanese line breaking rules configuration
	LineBreaking *LineBreakingSettings `yaml:"line_breaking,omitempty"`
}

// TextAreaSettings represents text area configuration for YAML reading.
type TextAreaSettings struct {
	X      *int `yaml:"x,omitempty"`
	Y      *int `yaml:"y,omitempty"`
	Width  *int `yaml:"width,omitempty"`
	Height *int `yaml:"height,omitempty"`
}

// LineBreakingSettings represents line breaking configuration for YAML reading.
type LineBreakingSettings struct {
	StartProhibited *string `yaml:"start_prohibited,omitempty"` // Characters that cannot start a line
	EndProhibited   *string `yaml:"end_prohibited,omitempty"`   // Characters that cannot end a line
}

// OverlayConfigSettings represents overlay configuration for YAML reading.
type OverlayConfigSettings struct {
	Visible   *bool              `yaml:"visible,omitempty"`   // Whether to render this overlay
	Image     *string            `yaml:"image,omitempty"`     // Path to image file
	Placement *PlacementSettings `yaml:"placement,omitempty"` // Image positioning
	Fit       *string            `yaml:"fit,omitempty"`       // Fit method
	Opacity   *float64           `yaml:"opacity,omitempty"`   // Image opacity (0.0-1.0)
}

// PlacementSettings represents placement configuration for YAML reading.
type PlacementSettings struct {
	X      *int `yaml:"x,omitempty"`
	Y      *int `yaml:"y,omitempty"`
	Width  *int `yaml:"width,omitempty"`  // nil means auto-detect from image
	Height *int `yaml:"height,omitempty"` // nil means auto-detect from image
}
