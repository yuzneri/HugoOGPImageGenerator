package main

import (
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
	baseConfig.Title.Size = 64.0
	baseConfig.Title.Color = "#000000"
	baseConfig.Description.Size = 32.0
	baseConfig.Description.Color = "#666666"
	baseConfig.Background.Color = "#FFFFFF"

	result := merger.MergeConfigs(baseConfig, nil)

	// Should be a copy of baseConfig, not the same pointer
	if result == baseConfig {
		t.Error("MergeConfigs should return a new config instance")
	}

	// Verify the values are the same
	if result.Title.Size != baseConfig.Title.Size {
		t.Errorf("Expected title size %f, got %f", baseConfig.Title.Size, result.Title.Size)
	}
	if result.Title.Color != baseConfig.Title.Color {
		t.Errorf("Expected title color %s, got %s", baseConfig.Title.Color, result.Title.Color)
	}
}

func TestConfigMerger_MergeConfigs_EmptyOGPFrontMatter(t *testing.T) {
	merger := NewConfigMerger()

	baseConfig := &Config{}
	baseConfig.Title.Size = 64.0
	baseConfig.Title.Color = "#000000"
	baseConfig.Description.Size = 32.0
	baseConfig.Description.Color = "#666666"
	baseConfig.Background.Color = "#FFFFFF"

	ogpFM := &OGPFrontMatter{}

	result := merger.MergeConfigs(baseConfig, ogpFM)

	// Should be a copy of baseConfig
	if result == baseConfig {
		t.Error("MergeConfigs should return a new config instance")
	}

	if result.Title.Size != baseConfig.Title.Size {
		t.Errorf("Expected title size %f, got %f", baseConfig.Title.Size, result.Title.Size)
	}

	if result.Title.Color != baseConfig.Title.Color {
		t.Errorf("Expected title color %s, got %s", baseConfig.Title.Color, result.Title.Color)
	}

	if result.Description.Size != baseConfig.Description.Size {
		t.Errorf("Expected description size %f, got %f", baseConfig.Description.Size, result.Description.Size)
	}

	if result.Description.Color != baseConfig.Description.Color {
		t.Errorf("Expected description color %s, got %s", baseConfig.Description.Color, result.Description.Color)
	}
}

func TestConfigMerger_MergeConfigs_TitleOverrides(t *testing.T) {
	merger := NewConfigMerger()

	baseConfig := &Config{}
	baseConfig.Title.Size = 64.0
	baseConfig.Title.Color = "#000000"
	baseTitleFont := "base-font.ttf"
	baseConfig.Title.Font = &baseTitleFont
	baseConfig.Title.BlockPosition = "middle-center"
	baseConfig.Title.LineAlignment = "left"
	baseConfig.Title.Overflow = "shrink"
	baseConfig.Title.MinSize = 12.0
	baseConfig.Title.LineHeight = 1.2
	baseConfig.Title.LetterSpacing = 1

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
	ogpFM.Title = &TextConfigOverride{
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

	if result.Title.Size != newSize {
		t.Errorf("Expected title size %f, got %f", newSize, result.Title.Size)
	}

	if result.Title.Color != newColor {
		t.Errorf("Expected title color %s, got %s", newColor, result.Title.Color)
	}

	if result.Title.Font == nil || *result.Title.Font != newFont {
		t.Errorf("Expected title font %s, got %v", newFont, result.Title.Font)
	}

	if result.Title.BlockPosition != newAlignment {
		t.Errorf("Expected title alignment %s, got %s", newAlignment, result.Title.BlockPosition)
	}

	if result.Title.LineAlignment != newLineAlignment {
		t.Errorf("Expected line alignment %s, got %s", newLineAlignment, result.Title.LineAlignment)
	}

	if result.Title.Overflow != newOverflow {
		t.Errorf("Expected overflow %s, got %s", newOverflow, result.Title.Overflow)
	}

	if result.Title.MinSize != newMinSize {
		t.Errorf("Expected min size %f, got %f", newMinSize, result.Title.MinSize)
	}

	if result.Title.LineHeight != newLineHeight {
		t.Errorf("Expected line height %f, got %f", newLineHeight, result.Title.LineHeight)
	}

	if result.Title.LetterSpacing != newLetterSpacing {
		t.Errorf("Expected letter spacing %d, got %d", newLetterSpacing, result.Title.LetterSpacing)
	}
}

func TestConfigMerger_MergeConfigs_DescriptionOverrides(t *testing.T) {
	merger := NewConfigMerger()

	baseConfig := &Config{}
	baseConfig.Description.Visible = true
	baseConfig.Description.Size = 32.0
	baseConfig.Description.Color = "#666666"
	baseDescriptionFont := "base-font.ttf"
	baseConfig.Description.Font = &baseDescriptionFont

	// Create override values
	newVisible := false
	newSize := 24.0
	newColor := "#333333"
	newFont := "description-font.ttf"

	ogpFM := &OGPFrontMatter{}
	ogpFM.Description = &TextConfigOverride{
		Visible: &newVisible,
		Font:    &newFont,
		Size:    &newSize,
		Color:   &newColor,
	}

	result := merger.MergeConfigs(baseConfig, ogpFM)

	if result.Description.Visible != newVisible {
		t.Errorf("Expected description visible %t, got %t", newVisible, result.Description.Visible)
	}

	if result.Description.Size != newSize {
		t.Errorf("Expected description size %f, got %f", newSize, result.Description.Size)
	}

	if result.Description.Color != newColor {
		t.Errorf("Expected description color %s, got %s", newColor, result.Description.Color)
	}

	if result.Description.Font == nil || *result.Description.Font != newFont {
		t.Errorf("Expected description font %s, got %v", newFont, result.Description.Font)
	}
}

func TestConfigMerger_MergeConfigs_TitleAreaOverrides(t *testing.T) {
	merger := NewConfigMerger()

	baseConfig := &Config{}
	baseConfig.Title.Area = TextArea{
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
	ogpFM.Title = &TextConfigOverride{
		Area: &TextAreaConfig{
			X:      &newX,
			Y:      &newY,
			Width:  &newWidth,
			Height: &newHeight,
		},
	}

	result := merger.MergeConfigs(baseConfig, ogpFM)

	if result.Title.Area.X != newX {
		t.Errorf("Expected area X %d, got %d", newX, result.Title.Area.X)
	}

	if result.Title.Area.Y != newY {
		t.Errorf("Expected area Y %d, got %d", newY, result.Title.Area.Y)
	}

	if result.Title.Area.Width != newWidth {
		t.Errorf("Expected area width %d, got %d", newWidth, result.Title.Area.Width)
	}

	if result.Title.Area.Height != newHeight {
		t.Errorf("Expected area height %d, got %d", newHeight, result.Title.Area.Height)
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
	ogpFM.Background = &BackgroundOverride{
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

func TestConfigMerger_MergeConfigs_LineBreakingOverrides(t *testing.T) {
	merger := NewConfigMerger()

	baseConfig := &Config{}
	baseConfig.Title.LineBreaking.StartProhibited = "。、"
	baseConfig.Title.LineBreaking.EndProhibited = "「（"

	newStartProhibited := "。、！？"
	newEndProhibited := "「（【"

	ogpFM := &OGPFrontMatter{}
	ogpFM.Title = &TextConfigOverride{
		LineBreaking: &LineBreakingOverride{
			StartProhibited: &newStartProhibited,
			EndProhibited:   &newEndProhibited,
		},
	}

	result := merger.MergeConfigs(baseConfig, ogpFM)

	if result.Title.LineBreaking.StartProhibited != newStartProhibited {
		t.Errorf("Expected start prohibited %s, got %s", newStartProhibited, result.Title.LineBreaking.StartProhibited)
	}

	if result.Title.LineBreaking.EndProhibited != newEndProhibited {
		t.Errorf("Expected end prohibited %s, got %s", newEndProhibited, result.Title.LineBreaking.EndProhibited)
	}
}

func TestConfigMerger_MergeConfigs_ComplexOverrides(t *testing.T) {
	merger := NewConfigMerger()

	baseConfig := getDefaultConfig()

	// Complex override scenario
	titleContent := "{{.Title}} - Custom"
	descContent := "{{.Description}} - Summary"
	titleSize := 72.0
	descSize := 28.0
	bgImage := "custom-bg.jpg"
	overlayImage := "logo.png"

	ogpFM := &OGPFrontMatter{}

	// Title overrides
	ogpFM.Title = &TextConfigOverride{
		Content: &titleContent,
		Size:    &titleSize,
	}

	// Description overrides
	ogpFM.Description = &TextConfigOverride{
		Content: &descContent,
		Size:    &descSize,
	}

	// Background override
	ogpFM.Background = &BackgroundOverride{
		Image: &bgImage,
	}

	// Overlay configuration
	ogpFM.Overlay = &ArticleOverlayConfig{
		Image: &overlayImage,
	}

	result := merger.MergeConfigs(baseConfig, ogpFM)

	// Verify title overrides
	if result.Title.Content == nil || *result.Title.Content != titleContent {
		t.Errorf("Expected title content %s, got %v", titleContent, result.Title.Content)
	}
	if result.Title.Size != titleSize {
		t.Errorf("Expected title size %f, got %f", titleSize, result.Title.Size)
	}

	// Verify description overrides
	if result.Description.Content == nil || *result.Description.Content != descContent {
		t.Errorf("Expected description content %s, got %v", descContent, result.Description.Content)
	}
	if result.Description.Size != descSize {
		t.Errorf("Expected description size %f, got %f", descSize, result.Description.Size)
	}

	// Verify background override
	if result.Background.Image == nil || *result.Background.Image != bgImage {
		t.Errorf("Expected background image %s, got %v", bgImage, result.Background.Image)
	}

	// Verify overlay is set
	if result.Overlay.Image == nil || *result.Overlay.Image != overlayImage {
		t.Errorf("Expected overlay image %s, got %v", overlayImage, result.Overlay.Image)
	}
}

func TestConfigMerger_MergeHelpers(t *testing.T) {
	merger := NewConfigMerger()

	// Test mergeStringPtr
	t.Run("mergeStringPtr", func(t *testing.T) {
		base := "base"
		override := "override"
		merger.mergeStringPtr(&base, &override)
		if base != override {
			t.Errorf("Expected %s, got %s", override, base)
		}

		base = "base"
		merger.mergeStringPtr(&base, nil)
		if base != "base" {
			t.Errorf("Expected base to remain unchanged")
		}
	})

	// Test mergeFloat64Ptr
	t.Run("mergeFloat64Ptr", func(t *testing.T) {
		base := 10.0
		override := 20.0
		merger.mergeFloat64Ptr(&base, &override)
		if base != override {
			t.Errorf("Expected %f, got %f", override, base)
		}

		base = 10.0
		merger.mergeFloat64Ptr(&base, nil)
		if base != 10.0 {
			t.Errorf("Expected base to remain unchanged")
		}
	})

	// Test mergeIntPtr
	t.Run("mergeIntPtr", func(t *testing.T) {
		base := 10
		override := 20
		merger.mergeIntPtr(&base, &override)
		if base != override {
			t.Errorf("Expected %d, got %d", override, base)
		}

		base = 10
		merger.mergeIntPtr(&base, nil)
		if base != 10 {
			t.Errorf("Expected base to remain unchanged")
		}
	})
}

func TestConfigMerger_OverlayConfiguration(t *testing.T) {
	merger := NewConfigMerger()

	// Test that empty overlay in base config is properly overridden
	baseConfig := &Config{}
	// Overlay is initialized as zero value by default

	overlayImage := "overlay.png"
	x := 10
	y := 20
	opacity := 0.8
	fit := "contain"

	ogpFM := &OGPFrontMatter{}
	ogpFM.Overlay = &ArticleOverlayConfig{
		Image:   &overlayImage,
		Opacity: &opacity,
		Fit:     &fit,
		Placement: &PlacementSettings{
			X: &x,
			Y: &y,
		},
	}

	result := merger.MergeConfigs(baseConfig, ogpFM)

	// Overlay is always present as a value type, check if it has meaningful data
	if result.Overlay.Image == nil {
		t.Error("Expected overlay to have configuration set")
		return
	}

	if result.Overlay.Image == nil || *result.Overlay.Image != overlayImage {
		t.Errorf("Expected overlay image %s, got %v", overlayImage, result.Overlay.Image)
	}

	if result.Overlay.Opacity != opacity {
		t.Errorf("Expected overlay opacity %f, got %f", opacity, result.Overlay.Opacity)
	}

	if result.Overlay.Fit != fit {
		t.Errorf("Expected overlay fit %s, got %s", fit, result.Overlay.Fit)
	}

	if result.Overlay.Placement.X != x {
		t.Errorf("Expected overlay placement X %d, got %d", x, result.Overlay.Placement.X)
	}
}

func TestConfigMerger_OutputConfiguration(t *testing.T) {
	merger := NewConfigMerger()

	baseConfig := &Config{}
	baseConfig.Output.Directory = "public"
	baseConfig.Output.Format = "png"

	// Only filename can be overridden in front matter
	filename := "custom-{.Title}.png"

	ogpFM := &OGPFrontMatter{}
	ogpFM.Output = &OutputOverride{
		Filename: &filename,
	}

	result := merger.MergeConfigs(baseConfig, ogpFM)

	// Directory and format should remain unchanged
	if result.Output.Directory != baseConfig.Output.Directory {
		t.Errorf("Output directory should not change: expected %s, got %s",
			baseConfig.Output.Directory, result.Output.Directory)
	}

	if result.Output.Format != baseConfig.Output.Format {
		t.Errorf("Output format should not change: expected %s, got %s",
			baseConfig.Output.Format, result.Output.Format)
	}

	// Filename should be updated
	if result.Output.Filename != filename {
		t.Errorf("Expected output filename %s, got %s", filename, result.Output.Filename)
	}
}

// Test that deep copy works correctly
func TestConfigMerger_DeepCopy(t *testing.T) {
	merger := NewConfigMerger()

	baseConfig := &Config{}
	baseConfig.Title.Size = 64.0
	baseConfig.Title.Color = "#000000"
	baseConfig.Title.Area = TextArea{X: 10, Y: 20, Width: 100, Height: 200}

	ogpFM := &OGPFrontMatter{}

	result := merger.MergeConfigs(baseConfig, ogpFM)

	// Modify the result
	result.Title.Size = 128.0
	result.Title.Area.X = 50

	// Original should be unchanged
	if baseConfig.Title.Size != 64.0 {
		t.Error("Base config was modified")
	}

	if baseConfig.Title.Area.X != 10 {
		t.Error("Base config area was modified")
	}
}

// Test pointer reference isolation for all pointer fields
func TestConfigMerger_PointerReferenceIsolation(t *testing.T) {
	merger := NewConfigMerger()

	// Set up base config with pointer fields
	originalTitleContent := "Original Title"
	originalBgImage := "original-bg.jpg"
	originalFilename := "original.png"
	originalOverlayImage := "original-overlay.png"
	originalFit := "contain"
	originalOpacity := 0.5

	baseConfig := &Config{}
	baseConfig.Title.Content = &originalTitleContent
	baseConfig.Background.Image = &originalBgImage
	baseConfig.Output.Filename = originalFilename
	baseConfig.Overlay = MainOverlayConfig{
		Visible: true,
		Image:   &originalOverlayImage,
		Fit:     originalFit,
		Opacity: originalOpacity,
	}

	t.Run("title content isolation", func(t *testing.T) {
		ogpFM := &OGPFrontMatter{}
		result := merger.MergeConfigs(baseConfig, ogpFM)

		// Modify merged config
		if result.Title.Content != nil {
			*result.Title.Content = "Modified Title"
		}

		// Original should remain unchanged
		if *baseConfig.Title.Content != originalTitleContent {
			t.Errorf("Base config title content was modified: expected %s, got %s",
				originalTitleContent, *baseConfig.Title.Content)
		}
	})

	t.Run("background image isolation", func(t *testing.T) {
		ogpFM := &OGPFrontMatter{}
		result := merger.MergeConfigs(baseConfig, ogpFM)

		// Modify merged config
		if result.Background.Image != nil {
			*result.Background.Image = "modified-bg.jpg"
		}

		// Original should remain unchanged
		if *baseConfig.Background.Image != originalBgImage {
			t.Errorf("Base config background image was modified: expected %s, got %s",
				originalBgImage, *baseConfig.Background.Image)
		}
	})

	t.Run("output filename isolation", func(t *testing.T) {
		ogpFM := &OGPFrontMatter{}
		result := merger.MergeConfigs(baseConfig, ogpFM)

		// Modify merged config
		result.Output.Filename = "modified.png"

		// Original should remain unchanged (value type, so automatically isolated)
		if baseConfig.Output.Filename != originalFilename {
			t.Errorf("Base config output filename was modified: expected %s, got %s",
				originalFilename, baseConfig.Output.Filename)
		}
	})

	t.Run("overlay pointer isolation", func(t *testing.T) {
		ogpFM := &OGPFrontMatter{}
		result := merger.MergeConfigs(baseConfig, ogpFM)

		// Modify merged config overlay fields
		if result.Overlay.Image != nil {
			*result.Overlay.Image = "modified-overlay.png"
		}
		// Fit is now a value type, can be assigned directly
		result.Overlay.Fit = "cover"
		// Opacity is now a value type, can be assigned directly
		result.Overlay.Opacity = 0.8

		// Original should remain unchanged
		if *baseConfig.Overlay.Image != originalOverlayImage {
			t.Errorf("Base config overlay image was modified: expected %s, got %s",
				originalOverlayImage, *baseConfig.Overlay.Image)
		}
		if baseConfig.Overlay.Fit != originalFit {
			t.Errorf("Base config overlay fit was modified: expected %s, got %s",
				originalFit, baseConfig.Overlay.Fit)
		}
		if baseConfig.Overlay.Opacity != originalOpacity {
			t.Errorf("Base config overlay opacity was modified: expected %f, got %f",
				originalOpacity, baseConfig.Overlay.Opacity)
		}
	})
}

// Test that merging with overrides doesn't affect base config pointers
func TestConfigMerger_OverridePointerIsolation(t *testing.T) {
	merger := NewConfigMerger()

	// Base config
	originalContent := "Original"
	baseConfig := &Config{}
	baseConfig.Title.Content = &originalContent

	// Override values
	overrideContent := "Override"
	ogpFM := &OGPFrontMatter{
		Title: &TextConfigOverride{
			Content: &overrideContent,
		},
	}

	result := merger.MergeConfigs(baseConfig, ogpFM)

	// Result should have override value
	if result.Title.Content == nil || *result.Title.Content != overrideContent {
		t.Errorf("Expected result content %s, got %v", overrideContent, result.Title.Content)
	}

	// Base config should remain unchanged
	if *baseConfig.Title.Content != originalContent {
		t.Errorf("Base config was modified: expected %s, got %s",
			originalContent, *baseConfig.Title.Content)
	}

	// Modify override source - result should not be affected
	overrideContent = "Modified Source"
	if result.Title.Content != nil && *result.Title.Content == "Modified Source" {
		t.Error("Result was affected by modifying override source")
	}
}

// Test deep nested pointer isolation (Overlay.Placement)
func TestConfigMerger_DeepNestedPointerIsolation(t *testing.T) {
	merger := NewConfigMerger()

	// Base config with nested pointers
	originalX := 10
	originalY := 20
	originalWidth := 100
	originalHeight := 200

	baseConfig := &Config{}
	baseConfig.Overlay = MainOverlayConfig{
		Visible: true,
		Placement: PlacementConfig{
			X:      originalX,
			Y:      originalY,
			Width:  &originalWidth,
			Height: &originalHeight,
		},
	}

	ogpFM := &OGPFrontMatter{}
	result := merger.MergeConfigs(baseConfig, ogpFM)

	// Modify result's placement values
	result.Overlay.Placement.X = 50
	result.Overlay.Placement.Y = 60
	newWidth := 300
	newHeight := 400
	result.Overlay.Placement.Width = &newWidth
	result.Overlay.Placement.Height = &newHeight

	// Original should remain unchanged
	if baseConfig.Overlay.Placement.X != originalX {
		t.Errorf("Base config placement X was modified: expected %d, got %d",
			originalX, baseConfig.Overlay.Placement.X)
	}
	if baseConfig.Overlay.Placement.Y != originalY {
		t.Errorf("Base config placement Y was modified: expected %d, got %d",
			originalY, baseConfig.Overlay.Placement.Y)
	}
	if baseConfig.Overlay.Placement.Width == nil || *baseConfig.Overlay.Placement.Width != originalWidth {
		t.Errorf("Base config placement Width was modified: expected %d, got %v",
			originalWidth, baseConfig.Overlay.Placement.Width)
	}
	if baseConfig.Overlay.Placement.Height == nil || *baseConfig.Overlay.Placement.Height != originalHeight {
		t.Errorf("Base config placement Height was modified: expected %d, got %v",
			originalHeight, baseConfig.Overlay.Placement.Height)
	}
}

// Helper function for creating string pointers in tests
func stringPtr(s string) *string {
	return &s
}

// Helper function for creating float64 pointers in tests
func float64Ptr(f float64) *float64 {
	return &f
}

// Test multi-article processing isolation
func TestConfigMerger_MultiArticleIsolation(t *testing.T) {
	merger := NewConfigMerger()

	// Base config (global)
	baseConfig := &Config{}
	baseConfig.Title.Size = 64.0
	baseConfig.Title.Color = "#000000"
	baseConfig.Background.Color = "#FFFFFF"
	baseConfig.Overlay = MainOverlayConfig{
		Visible: false,
	}

	// Article 1: Has overlay configuration
	article1Content := "Article 1 Title"
	article1OverlayImage := "article1-overlay.png"
	ogpFM1 := &OGPFrontMatter{
		Title: &TextConfigOverride{
			Content: &article1Content,
		},
		Overlay: &ArticleOverlayConfig{
			Visible: boolPtr(true),
			Image:   &article1OverlayImage,
		},
	}

	// Article 2: No overlay configuration, different title
	article2Content := "Article 2 Title"
	ogpFM2 := &OGPFrontMatter{
		Title: &TextConfigOverride{
			Content: &article2Content,
		},
	}

	// Article 3: Different overlay configuration
	article3Content := "Article 3 Title"
	article3OverlayImage := "article3-overlay.png"
	ogpFM3 := &OGPFrontMatter{
		Title: &TextConfigOverride{
			Content: &article3Content,
		},
		Overlay: &ArticleOverlayConfig{
			Visible: boolPtr(true),
			Image:   &article3OverlayImage,
		},
	}

	// Process articles in sequence (simulating real usage)
	result1 := merger.MergeConfigs(baseConfig, ogpFM1)
	result2 := merger.MergeConfigs(baseConfig, ogpFM2)
	result3 := merger.MergeConfigs(baseConfig, ogpFM3)

	// Verify Article 1 config
	if result1.Title.Content == nil || *result1.Title.Content != article1Content {
		t.Errorf("Article 1 title content incorrect: expected %s, got %v",
			article1Content, result1.Title.Content)
	}
	if !result1.Overlay.Visible {
		t.Error("Article 1 overlay should be visible")
	}
	if result1.Overlay.Image == nil || *result1.Overlay.Image != article1OverlayImage {
		t.Errorf("Article 1 overlay image incorrect: expected %s, got %v",
			article1OverlayImage, result1.Overlay.Image)
	}

	// Verify Article 2 config (should not have overlay from Article 1)
	if result2.Title.Content == nil || *result2.Title.Content != article2Content {
		t.Errorf("Article 2 title content incorrect: expected %s, got %v",
			article2Content, result2.Title.Content)
	}
	if result2.Overlay.Visible {
		t.Error("Article 2 overlay should not be visible (no overlay in front matter)")
	}
	if result2.Overlay.Image != nil && *result2.Overlay.Image == article1OverlayImage {
		t.Error("Article 2 incorrectly inherited overlay image from Article 1")
	}

	// Verify Article 3 config
	if result3.Title.Content == nil || *result3.Title.Content != article3Content {
		t.Errorf("Article 3 title content incorrect: expected %s, got %v",
			article3Content, result3.Title.Content)
	}
	if !result3.Overlay.Visible {
		t.Error("Article 3 overlay should be visible")
	}
	if result3.Overlay.Image == nil || *result3.Overlay.Image != article3OverlayImage {
		t.Errorf("Article 3 overlay image incorrect: expected %s, got %v",
			article3OverlayImage, result3.Overlay.Image)
	}

	// Verify base config remains unchanged
	if baseConfig.Title.Size != 64.0 {
		t.Error("Base config was modified during processing")
	}
	if baseConfig.Overlay.Visible {
		t.Error("Base config overlay visibility was modified")
	}
}

// Test new helper functions for value copying
func TestConfigMerger_NewHelperFunctions(t *testing.T) {
	merger := NewConfigMerger()

	t.Run("mergeStringPtrValue", func(t *testing.T) {
		source := "source value"
		var target *string

		merger.mergeStringPtrValue(&target, &source)

		// Target should have the value
		if target == nil || *target != source {
			t.Errorf("Expected target to have value %s, got %v", source, target)
		}

		// Modify source - target should not be affected
		source = "modified source"
		if target == nil || *target == "modified source" {
			t.Error("Target was affected by modifying source")
		}

		// Test with nil source
		merger.mergeStringPtrValue(&target, nil)
		if target == nil || *target != "source value" {
			t.Error("Target should not change when source is nil")
		}
	})

	t.Run("mergeFloat64PtrValue", func(t *testing.T) {
		source := 42.5
		var target *float64

		merger.mergeFloat64PtrValue(&target, &source)

		// Target should have the value
		if target == nil || *target != source {
			t.Errorf("Expected target to have value %f, got %v", source, target)
		}

		// Modify source - target should not be affected
		source = 99.9
		if target == nil || *target == 99.9 {
			t.Error("Target was affected by modifying source")
		}

		// Test with nil source
		merger.mergeFloat64PtrValue(&target, nil)
		if target == nil || *target != 42.5 {
			t.Error("Target should not change when source is nil")
		}
	})

	t.Run("mergeIntPtrValue", func(t *testing.T) {
		source := 123
		var target *int

		merger.mergeIntPtrValue(&target, &source)

		// Target should have the value
		if target == nil || *target != source {
			t.Errorf("Expected target to have value %d, got %v", source, target)
		}

		// Modify source - target should not be affected
		source = 456
		if target == nil || *target == 456 {
			t.Error("Target was affected by modifying source")
		}

		// Test with nil source
		merger.mergeIntPtrValue(&target, nil)
		if target == nil || *target != 123 {
			t.Error("Target should not change when source is nil")
		}
	})
}

// Test complex scenario with nested pointer overrides
func TestConfigMerger_ComplexNestedOverrides(t *testing.T) {
	merger := NewConfigMerger()

	// Base config with nested structure
	baseImage := "base-overlay.png"
	baseFit := "contain"
	baseOpacity := 0.5
	baseX := 10
	baseY := 20

	baseConfig := &Config{}
	baseConfig.Overlay = MainOverlayConfig{
		Visible: true,
		Image:   &baseImage,
		Fit:     baseFit,
		Opacity: baseOpacity,
		Placement: PlacementConfig{
			X: baseX,
			Y: baseY,
		},
	}

	// Override with nested changes
	overrideImage := "override-overlay.png"
	overrideOpacity := 0.8
	overrideX := 50
	overrideWidth := 200

	ogpFM := &OGPFrontMatter{
		Overlay: &ArticleOverlayConfig{
			Image:   &overrideImage,
			Opacity: &overrideOpacity,
			Placement: &PlacementSettings{
				X:     &overrideX,
				Width: &overrideWidth,
				// Y and Height not specified - should keep base values
			},
		},
	}

	result := merger.MergeConfigs(baseConfig, ogpFM)

	// Check overridden values
	if result.Overlay.Image == nil || *result.Overlay.Image != overrideImage {
		t.Errorf("Expected overlay image %s, got %v", overrideImage, result.Overlay.Image)
	}
	if result.Overlay.Opacity != overrideOpacity {
		t.Errorf("Expected overlay opacity %f, got %f", overrideOpacity, result.Overlay.Opacity)
	}
	if result.Overlay.Placement.X != overrideX {
		t.Errorf("Expected placement X %d, got %d", overrideX, result.Overlay.Placement.X)
	}
	if result.Overlay.Placement.Width == nil || *result.Overlay.Placement.Width != overrideWidth {
		t.Errorf("Expected placement width %d, got %v", overrideWidth, result.Overlay.Placement.Width)
	}

	// Check preserved values (should maintain base values for non-overridden fields)
	if result.Overlay.Fit != baseFit {
		t.Errorf("Expected fit to remain %s, got %s", baseFit, result.Overlay.Fit)
	}
	if result.Overlay.Placement.Y != baseY {
		t.Errorf("Expected placement Y to remain %d, got %d", baseY, result.Overlay.Placement.Y)
	}

	// Check base config isolation
	if *baseConfig.Overlay.Image != baseImage {
		t.Errorf("Base config image was modified: expected %s, got %s",
			baseImage, *baseConfig.Overlay.Image)
	}
	if baseConfig.Overlay.Placement.X != baseX {
		t.Errorf("Base config placement X was modified: expected %d, got %d",
			baseX, baseConfig.Overlay.Placement.X)
	}
}

// Helper function for creating bool pointers
func boolPtr(b bool) *bool {
	return &b
}
