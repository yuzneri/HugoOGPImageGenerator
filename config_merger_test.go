package main

import (
	"reflect"
	"testing"
)

func TestNewConfigMerger(t *testing.T) {
	merger := NewConfigMerger()

	if merger == nil {
		t.Error("NewConfigMerger should return a non-nil merger")
	}
}

func TestConfigMerger_MergeConfigs_NilOGPFrontMatter(t *testing.T) {
	merger := NewConfigMerger()

	baseConfig := &Config{}
	baseConfig.Text.Size = 64.0
	baseConfig.Text.Color = "#000000"
	baseConfig.Background.Color = "#FFFFFF"

	result := merger.MergeConfigs(baseConfig, nil)

	if result != baseConfig {
		t.Error("MergeConfigs should return baseConfig when ogpFM is nil")
	}
}

func TestConfigMerger_MergeConfigs_EmptyOGPFrontMatter(t *testing.T) {
	merger := NewConfigMerger()

	baseConfig := &Config{}
	baseConfig.Text.Size = 64.0
	baseConfig.Text.Color = "#000000"
	baseConfig.Background.Color = "#FFFFFF"

	ogpFM := &OGPFrontMatter{}

	result := merger.MergeConfigs(baseConfig, ogpFM)

	// Should be a copy of baseConfig
	if result == baseConfig {
		t.Error("MergeConfigs should return a new config instance")
	}

	if result.Text.Size != baseConfig.Text.Size {
		t.Errorf("Expected text size %f, got %f", baseConfig.Text.Size, result.Text.Size)
	}

	if result.Text.Color != baseConfig.Text.Color {
		t.Errorf("Expected text color %s, got %s", baseConfig.Text.Color, result.Text.Color)
	}
}

func TestConfigMerger_MergeConfigs_TextOverrides(t *testing.T) {
	merger := NewConfigMerger()

	baseConfig := &Config{}
	baseConfig.Text.Size = 64.0
	baseConfig.Text.Color = "#000000"
	baseConfig.Text.Font = "base-font.ttf"
	baseConfig.Text.BlockPosition = "middle-center"
	baseConfig.Text.LineAlignment = "left"
	baseConfig.Text.Overflow = "shrink"
	baseConfig.Text.MinSize = 12.0
	baseConfig.Text.LineHeight = 1.2
	baseConfig.Text.LetterSpacing = 1

	// Create override values
	newSize := 48.0
	newColor := "#FF0000"
	newFont := "override-font.ttf"
	newAlignment := "top-left"
	newLineAlignment := "center"
	newOverflow := "clip"
	newMinSize := 10.0
	newLineHeight := 1.5
	newLetterSpacing := 2

	ogpFM := &OGPFrontMatter{}
	ogpFM.Text = &struct {
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
	}{
		Font:          &newFont,
		Size:          &newSize,
		Color:         &newColor,
		BlockPosition: &newAlignment,
		LineAlignment: &newLineAlignment,
		Overflow:      &newOverflow,
		MinSize:       &newMinSize,
		LineHeight:    &newLineHeight,
		LetterSpacing: &newLetterSpacing,
	}

	result := merger.MergeConfigs(baseConfig, ogpFM)

	if result.Text.Size != newSize {
		t.Errorf("Expected text size %f, got %f", newSize, result.Text.Size)
	}

	if result.Text.Color != newColor {
		t.Errorf("Expected text color %s, got %s", newColor, result.Text.Color)
	}

	if result.Text.Font != newFont {
		t.Errorf("Expected text font %s, got %s", newFont, result.Text.Font)
	}

	if result.Text.BlockPosition != newAlignment {
		t.Errorf("Expected text alignment %s, got %s", newAlignment, result.Text.BlockPosition)
	}

	if result.Text.LineAlignment != newLineAlignment {
		t.Errorf("Expected line alignment %s, got %s", newLineAlignment, result.Text.LineAlignment)
	}

	if result.Text.Overflow != newOverflow {
		t.Errorf("Expected overflow %s, got %s", newOverflow, result.Text.Overflow)
	}

	if result.Text.MinSize != newMinSize {
		t.Errorf("Expected min size %f, got %f", newMinSize, result.Text.MinSize)
	}

	if result.Text.LineHeight != newLineHeight {
		t.Errorf("Expected line height %f, got %f", newLineHeight, result.Text.LineHeight)
	}

	if result.Text.LetterSpacing != newLetterSpacing {
		t.Errorf("Expected letter spacing %d, got %d", newLetterSpacing, result.Text.LetterSpacing)
	}
}

func TestConfigMerger_MergeConfigs_TextAreaOverrides(t *testing.T) {
	merger := NewConfigMerger()

	baseConfig := &Config{}
	baseConfig.Text.Area = TextArea{
		X:      50,
		Y:      50,
		Width:  700,
		Height: 400,
	}

	newX := 100
	newY := 150
	newWidth := 600
	newHeight := 300

	ogpFM := &OGPFrontMatter{}
	ogpFM.Text = &struct {
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
	}{
		Area: &struct {
			X      *int `yaml:"x,omitempty"`
			Y      *int `yaml:"y,omitempty"`
			Width  *int `yaml:"width,omitempty"`
			Height *int `yaml:"height,omitempty"`
		}{
			X:      &newX,
			Y:      &newY,
			Width:  &newWidth,
			Height: &newHeight,
		},
	}

	result := merger.MergeConfigs(baseConfig, ogpFM)

	if result.Text.Area.X != newX {
		t.Errorf("Expected area X %d, got %d", newX, result.Text.Area.X)
	}

	if result.Text.Area.Y != newY {
		t.Errorf("Expected area Y %d, got %d", newY, result.Text.Area.Y)
	}

	if result.Text.Area.Width != newWidth {
		t.Errorf("Expected area width %d, got %d", newWidth, result.Text.Area.Width)
	}

	if result.Text.Area.Height != newHeight {
		t.Errorf("Expected area height %d, got %d", newHeight, result.Text.Area.Height)
	}
}

func TestConfigMerger_MergeConfigs_BackgroundOverrides(t *testing.T) {
	merger := NewConfigMerger()

	baseConfig := &Config{}
	baseConfig.Background.Color = "#FFFFFF"
	baseConfig.Background.Image = nil

	newColor := "#FF0000"
	newImage := "new-background.png"

	ogpFM := &OGPFrontMatter{}
	ogpFM.Background = &struct {
		Image *string `yaml:"image,omitempty"` // Path to background image (relative to article directory)
		Color *string `yaml:"color,omitempty"` // Background color (hex)
	}{
		Color: &newColor,
		Image: &newImage,
	}

	result := merger.MergeConfigs(baseConfig, ogpFM)

	if result.Background.Color != newColor {
		t.Errorf("Expected background color %s, got %s", newColor, result.Background.Color)
	}

	if result.Background.Image == nil || *result.Background.Image != newImage {
		t.Errorf("Expected background image %s, got %v", newImage, result.Background.Image)
	}
}

func TestConfigMerger_MergeConfigs_OutputOverrides(t *testing.T) {
	merger := NewConfigMerger()

	baseConfig := &Config{}
	baseConfig.Output.Filename = nil

	newFilename := "custom-filename"

	ogpFM := &OGPFrontMatter{}
	ogpFM.Output = &struct {
		Filename *string `yaml:"filename,omitempty"` // Custom filename template (optional)
	}{
		Filename: &newFilename,
	}

	result := merger.MergeConfigs(baseConfig, ogpFM)

	if result.Output.Filename == nil || *result.Output.Filename != newFilename {
		t.Errorf("Expected output filename %s, got %v", newFilename, result.Output.Filename)
	}
}

func TestConfigMerger_MergeConfigs_OverlayOverrides(t *testing.T) {
	merger := NewConfigMerger()

	baseConfig := &Config{}
	baseConfig.Overlay = nil

	newOverlay := &struct {
		Image     *string `yaml:"image,omitempty"` // Path to image file (relative to article directory)
		Placement *struct {
			X      *int `yaml:"x,omitempty"`      // X position
			Y      *int `yaml:"y,omitempty"`      // Y position
			Width  *int `yaml:"width,omitempty"`  // Image width
			Height *int `yaml:"height,omitempty"` // Image height
		} `yaml:"placement,omitempty"`
		Fit     *string  `yaml:"fit,omitempty"`     // Fit method ("cover", "contain", "fill", "none")
		Opacity *float64 `yaml:"opacity,omitempty"` // Image opacity (0.0-1.0)
	}{}
	imageValue := "overlay.png"
	newOverlay.Image = &imageValue
	placementValue := struct {
		X      *int `yaml:"x,omitempty"`      // X position
		Y      *int `yaml:"y,omitempty"`      // Y position
		Width  *int `yaml:"width,omitempty"`  // Image width
		Height *int `yaml:"height,omitempty"` // Image height
	}{}
	x := 100
	y := 200
	width := 300
	height := 400
	placementValue.X = &x
	placementValue.Y = &y
	placementValue.Width = &width
	placementValue.Height = &height
	newOverlay.Placement = &placementValue
	fitValue := "contain"
	newOverlay.Fit = &fitValue
	opacityValue := 0.8
	newOverlay.Opacity = &opacityValue

	ogpFM := &OGPFrontMatter{}
	ogpFM.Overlay = newOverlay

	result := merger.MergeConfigs(baseConfig, ogpFM)

	if !reflect.DeepEqual(result.Overlay, newOverlay) {
		t.Errorf("Expected overlay config %+v, got %+v", newOverlay, result.Overlay)
	}
}

func TestConfigMerger_mergeStringPtr(t *testing.T) {
	merger := NewConfigMerger()

	target := "original"
	source := "updated"

	merger.mergeStringPtr(&target, &source)

	if target != source {
		t.Errorf("Expected target to be %s, got %s", source, target)
	}

	// Test with nil source
	target = "original"
	merger.mergeStringPtr(&target, nil)

	if target != "original" {
		t.Errorf("Expected target to remain %s when source is nil, got %s", "original", target)
	}
}

func TestConfigMerger_mergeFloat64Ptr(t *testing.T) {
	merger := NewConfigMerger()

	target := 1.0
	source := 2.0

	merger.mergeFloat64Ptr(&target, &source)

	if target != source {
		t.Errorf("Expected target to be %f, got %f", source, target)
	}

	// Test with nil source
	target = 1.0
	merger.mergeFloat64Ptr(&target, nil)

	if target != 1.0 {
		t.Errorf("Expected target to remain %f when source is nil, got %f", 1.0, target)
	}
}

func TestConfigMerger_mergeIntPtr(t *testing.T) {
	merger := NewConfigMerger()

	target := 1
	source := 2

	merger.mergeIntPtr(&target, &source)

	if target != source {
		t.Errorf("Expected target to be %d, got %d", source, target)
	}

	// Test with nil source
	target = 1
	merger.mergeIntPtr(&target, nil)

	if target != 1 {
		t.Errorf("Expected target to remain %d when source is nil, got %d", 1, target)
	}
}

func TestConfigMerger_MergeConfigs_PartialOverrides(t *testing.T) {
	merger := NewConfigMerger()

	baseConfig := &Config{}
	baseConfig.Text.Size = 64.0
	baseConfig.Text.Color = "#000000"
	baseConfig.Text.Font = "base-font.ttf"
	baseConfig.Background.Color = "#FFFFFF"

	// Only override text size
	newSize := 48.0

	ogpFM := &OGPFrontMatter{}
	ogpFM.Text = &struct {
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
	}{
		Size: &newSize,
	}

	result := merger.MergeConfigs(baseConfig, ogpFM)

	// Size should be overridden
	if result.Text.Size != newSize {
		t.Errorf("Expected text size %f, got %f", newSize, result.Text.Size)
	}

	// Color and Font should remain from base config
	if result.Text.Color != baseConfig.Text.Color {
		t.Errorf("Expected text color %s to remain unchanged, got %s", baseConfig.Text.Color, result.Text.Color)
	}

	if result.Text.Font != baseConfig.Text.Font {
		t.Errorf("Expected text font %s to remain unchanged, got %s", baseConfig.Text.Font, result.Text.Font)
	}

	// Background should remain unchanged
	if result.Background.Color != baseConfig.Background.Color {
		t.Errorf("Expected background color %s to remain unchanged, got %s", baseConfig.Background.Color, result.Background.Color)
	}
}
