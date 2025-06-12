package main

import (
	"image"
	"testing"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func TestImageRenderer_RenderTextOnImage_TitleHidden(t *testing.T) {
	// Load font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}

	// Create test image
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))

	// Create renderer
	renderer := NewImageRenderer()

	// Create test config with title hidden
	config := getDefaultConfig()
	config.Title.Visible = false // Hide title
	config.Description.Visible = true

	options := &RenderOptions{
		Font:        font,
		Config:      config,
		Title:       "This title should not render",
		Description: "This description should render",
		TestMode:    false,
	}

	// Should not error even with title hidden
	err = renderer.RenderTextOnImage(img, options)
	if err != nil {
		t.Errorf("RenderTextOnImage should not return error with title hidden: %v", err)
	}

	// In a real implementation, we could check that title pixels were not drawn
	// For now, we just verify it doesn't crash
}

func TestImageRenderer_RenderTextOnImage_DescriptionHidden(t *testing.T) {
	// Load font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}

	// Create test image
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))

	// Create renderer
	renderer := NewImageRenderer()

	// Create test config with description hidden
	config := getDefaultConfig()
	config.Title.Visible = true
	config.Description.Visible = false // Hide description

	options := &RenderOptions{
		Font:        font,
		Config:      config,
		Title:       "This title should render",
		Description: "This description should not render",
		TestMode:    false,
	}

	// Should not error even with description hidden
	err = renderer.RenderTextOnImage(img, options)
	if err != nil {
		t.Errorf("RenderTextOnImage should not return error with description hidden: %v", err)
	}
}

func TestImageRenderer_RenderTextOnImage_BothHidden(t *testing.T) {
	// Load font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}

	// Create test image
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))

	// Create renderer
	renderer := NewImageRenderer()

	// Create test config with both hidden
	config := getDefaultConfig()
	config.Title.Visible = false       // Hide title
	config.Description.Visible = false // Hide description

	options := &RenderOptions{
		Font:        font,
		Config:      config,
		Title:       "This title should not render",
		Description: "This description should not render",
		TestMode:    false,
	}

	// Should not error even with both hidden
	err = renderer.RenderTextOnImage(img, options)
	if err != nil {
		t.Errorf("RenderTextOnImage should not return error with both text elements hidden: %v", err)
	}
}

func TestConfigMerger_VisibleField(t *testing.T) {
	merger := NewConfigMerger()

	// Test that visible field is properly merged
	baseConfig := &Config{}
	baseConfig.Title.Visible = true
	baseConfig.Description.Visible = true

	// Override to hide description
	descriptionHidden := false
	titleHidden := false

	ogpFM := &OGPFrontMatter{}
	ogpFM.Title = &TextConfigOverride{
		Visible: &titleHidden,
	}
	ogpFM.Description = &TextConfigOverride{
		Visible: &descriptionHidden,
	}

	result := merger.MergeConfigs(baseConfig, ogpFM)

	if result.Title.Visible != titleHidden {
		t.Errorf("Expected title visible %t, got %t", titleHidden, result.Title.Visible)
	}

	if result.Description.Visible != descriptionHidden {
		t.Errorf("Expected description visible %t, got %t", descriptionHidden, result.Description.Visible)
	}
}

func TestConfigMerger_VisibleField_PartialOverride(t *testing.T) {
	merger := NewConfigMerger()

	// Test that only specified visible fields are overridden
	baseConfig := &Config{}
	baseConfig.Title.Visible = true
	baseConfig.Description.Visible = true

	// Override only description
	descriptionHidden := false

	ogpFM := &OGPFrontMatter{}
	ogpFM.Description = &TextConfigOverride{
		Visible: &descriptionHidden,
	}
	// Title override not specified, should remain true

	result := merger.MergeConfigs(baseConfig, ogpFM)

	if result.Title.Visible != true {
		t.Errorf("Expected title visible to remain %t, got %t", true, result.Title.Visible)
	}

	if result.Description.Visible != descriptionHidden {
		t.Errorf("Expected description visible %t, got %t", descriptionHidden, result.Description.Visible)
	}
}

func TestDefaultConfig_VisibleValues(t *testing.T) {
	config := getDefaultConfig()

	// Title should be visible by default
	if !config.Title.Visible {
		t.Error("Title should be visible by default")
	}

	// Description should be hidden by default
	if config.Description.Visible {
		t.Error("Description should be hidden by default")
	}

	// Overlay should be hidden by default
	if config.Overlay.Visible {
		t.Error("Overlay should be hidden by default")
	}
}
