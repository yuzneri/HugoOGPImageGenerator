package main

import (
	"testing"
)

// TestExplicitZeroValues tests that explicitly specified zero values are preserved
func TestExplicitZeroValues(t *testing.T) {
	t.Run("Explicit zero in placement X and Y from front matter", func(t *testing.T) {
		merger := NewConfigMerger()

		// Start with config that has non-zero placement
		config := getDefaultConfig()
		config.Overlay.Visible = true
		config.Overlay.Image = &[]string{"test.jpg"}[0]
		config.Overlay.Placement.X = 100
		config.Overlay.Placement.Y = 50

		// Front matter with explicit zeros
		ogpFM := &OGPFrontMatter{
			Overlay: &ArticleOverlayConfig{
				Placement: &PlacementSettings{
					X: &[]int{0}[0], // Explicitly set X = 0
					Y: &[]int{0}[0], // Explicitly set Y = 0
					// Width and Height not specified
				},
			},
		}

		// Apply front matter overrides
		result := merger.applyFrontMatterOverrides(config, ogpFM)

		// Verify that explicit zeros are preserved
		if result.Overlay.Placement.X != 0 {
			t.Errorf("Expected X=0 (explicitly set), got %d", result.Overlay.Placement.X)
		}
		if result.Overlay.Placement.Y != 0 {
			t.Errorf("Expected Y=0 (explicitly set), got %d", result.Overlay.Placement.Y)
		}
	})

	t.Run("Explicit zero in area coordinates from type config", func(t *testing.T) {
		merger := NewConfigMerger()

		// Start with default config
		config := getDefaultConfig()

		// Type config with explicit zeros
		typeSettings := &ConfigSettings{
			Title: &TextSettings{
				Area: &TextAreaSettings{
					X: &[]int{0}[0], // Explicitly set X = 0
					Y: &[]int{0}[0], // Explicitly set Y = 0
					// Width and Height not specified - should remain default
				},
			},
		}

		// Apply type settings
		merger.applySettingsToConfig(config, typeSettings)

		// Verify that explicit zeros are used
		if config.Title.Area.X != 0 {
			t.Errorf("Expected X=0 (explicitly set), got %d", config.Title.Area.X)
		}
		if config.Title.Area.Y != 0 {
			t.Errorf("Expected Y=0 (explicitly set), got %d", config.Title.Area.Y)
		}
		// Width and Height should remain defaults
		if config.Title.Area.Width != DefaultTitleAreaWidth {
			t.Errorf("Expected Width=%d (default), got %d", DefaultTitleAreaWidth, config.Title.Area.Width)
		}
		if config.Title.Area.Height != DefaultTitleAreaHeight {
			t.Errorf("Expected Height=%d (default), got %d", DefaultTitleAreaHeight, config.Title.Area.Height)
		}
	})

	t.Run("Mixed zero and non-zero values", func(t *testing.T) {
		merger := NewConfigMerger()

		// Global settings: set some values
		globalSettings := &ConfigSettings{
			Overlay: &OverlayConfigSettings{
				Visible: &[]bool{true}[0],
				Image:   &[]string{"default.jpg"}[0],
				Placement: &PlacementSettings{
					X: &[]int{100}[0], // X = 100
					Y: &[]int{50}[0],  // Y = 50
				},
			},
		}

		// Type settings: override with zeros and new values
		typeSettings := &ConfigSettings{
			Overlay: &OverlayConfigSettings{
				Placement: &PlacementSettings{
					X:      &[]int{0}[0],   // Explicitly override X to 0
					Height: &[]int{300}[0], // Set height = 300
					// Y not specified - should remain from global (50)
				},
			},
		}

		// Front matter: mix of zero and non-zero
		ogpFM := &OGPFrontMatter{
			Overlay: &ArticleOverlayConfig{
				Placement: &PlacementSettings{
					Y:     &[]int{0}[0],   // Explicitly set Y = 0 (override global)
					Width: &[]int{200}[0], // Set width = 200
					// X not specified - should remain from type (0)
					// Height not specified - should remain from type (300)
				},
			},
		}

		// Apply 4-level merge
		result := merger.MergeConfigsWithSettings(getDefaultConfig(), globalSettings, typeSettings, ogpFM)

		// Verify final placement
		if result.Overlay.Placement.X != 0 {
			t.Errorf("Expected X=0 (from type, explicit zero), got %d", result.Overlay.Placement.X)
		}
		if result.Overlay.Placement.Y != 0 {
			t.Errorf("Expected Y=0 (from front matter, explicit zero), got %d", result.Overlay.Placement.Y)
		}
		if result.Overlay.Placement.Width == nil || *result.Overlay.Placement.Width != 200 {
			t.Errorf("Expected Width=200 (from front matter), got %v", result.Overlay.Placement.Width)
		}
		if result.Overlay.Placement.Height == nil || *result.Overlay.Placement.Height != 300 {
			t.Errorf("Expected Height=300 (from type), got %v", result.Overlay.Placement.Height)
		}
	})

	t.Run("Explicit zero in text size and spacing", func(t *testing.T) {
		merger := NewConfigMerger()

		// Start with default config
		config := getDefaultConfig()

		// Front matter with explicit zero values for text settings
		ogpFM := &OGPFrontMatter{
			Title: &TextConfigOverride{
				LetterSpacing: &[]int{0}[0], // Explicitly set letter spacing = 0
				// Size not specified - should remain default
			},
		}

		// Apply front matter overrides
		result := merger.applyFrontMatterOverrides(config, ogpFM)

		// Verify that explicit zero is preserved
		if result.Title.LetterSpacing != 0 {
			t.Errorf("Expected LetterSpacing=0 (explicitly set), got %d", result.Title.LetterSpacing)
		}
		// Size should remain default
		if result.Title.Size != DefaultTitleFontSize {
			t.Errorf("Expected Size=%f (default), got %f", DefaultTitleFontSize, result.Title.Size)
		}
	})

	t.Run("Distinguish nil vs zero for overlay opacity", func(t *testing.T) {
		merger := NewConfigMerger()

		// Start with default config
		config := getDefaultConfig()
		config.Overlay.Visible = true
		config.Overlay.Image = &[]string{"test.jpg"}[0]
		config.Overlay.Opacity = 0.8

		// Front matter with explicit zero opacity
		ogpFM := &OGPFrontMatter{
			Overlay: &ArticleOverlayConfig{
				Opacity: &[]float64{0.0}[0], // Explicitly set opacity = 0.0 (invisible)
			},
		}

		// Apply front matter overrides
		result := merger.applyFrontMatterOverrides(config, ogpFM)

		// Verify that explicit zero opacity is used (making overlay invisible)
		if result.Overlay.Opacity != 0.0 {
			t.Errorf("Expected Opacity=0.0 (explicitly set), got %f", result.Overlay.Opacity)
		}
	})
}
