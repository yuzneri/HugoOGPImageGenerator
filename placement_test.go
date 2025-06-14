package main

import (
	"testing"
)

// TestPlacementMerging tests that partial placement configurations are merged correctly
func TestPlacementMerging(t *testing.T) {
	t.Run("Partial overlay placement from type config", func(t *testing.T) {
		merger := NewConfigMerger()

		// Start with default config
		config := getDefaultConfig()

		// Type config with only X and height specified (Y should remain default)
		typeSettings := &ConfigSettings{
			Overlay: &OverlayConfigSettings{
				Visible: &[]bool{true}[0],
				Image:   &[]string{"cover.jpg"}[0],
				Placement: &PlacementSettings{
					X:      &[]int{100}[0], // Specify X = 100
					Height: &[]int{580}[0], // Specify height = 580
					// Y and Width are nil (should remain as defaults)
				},
			},
		}

		// Apply type settings
		merger.applySettingsToConfig(config, typeSettings)

		// Verify placement
		if config.Overlay.Placement.X != 100 {
			t.Errorf("Expected X=100, got %d", config.Overlay.Placement.X)
		}
		if config.Overlay.Placement.Y != DefaultOverlayY { // Should remain default (50)
			t.Errorf("Expected Y=%d (default), got %d", DefaultOverlayY, config.Overlay.Placement.Y)
		}
		if config.Overlay.Placement.Width != nil {
			t.Errorf("Expected Width=nil (auto-detect), got %v", config.Overlay.Placement.Width)
		}
		if config.Overlay.Placement.Height == nil || *config.Overlay.Placement.Height != 580 {
			t.Errorf("Expected Height=580, got %v", config.Overlay.Placement.Height)
		}
	})

	t.Run("Partial overlay placement from front matter", func(t *testing.T) {
		merger := NewConfigMerger()

		// Start with config that has some overlay settings
		config := getDefaultConfig()
		config.Overlay.Visible = true
		config.Overlay.Image = &[]string{"default.jpg"}[0]
		config.Overlay.Placement.X = 50 // default
		config.Overlay.Placement.Y = 50 // default

		// Front matter with only X and height specified (Y should remain unchanged)
		ogpFM := &OGPFrontMatter{
			Overlay: &ArticleOverlayConfig{
				Placement: &PlacementSettings{
					X:      &[]int{100}[0], // Override X
					Y:      nil,            // Not specified (nil) - should NOT override
					Width:  nil,            // Not specified - should remain nil
					Height: &[]int{580}[0], // Specify height
				},
			},
		}

		// Apply front matter overrides
		result := merger.applyFrontMatterOverrides(config, ogpFM)

		// Verify placement
		if result.Overlay.Placement.X != 100 {
			t.Errorf("Expected X=100, got %d", result.Overlay.Placement.X)
		}
		if result.Overlay.Placement.Y != 50 { // Should remain original value (50)
			t.Errorf("Expected Y=50 (preserved), got %d", result.Overlay.Placement.Y)
		}
		if result.Overlay.Placement.Width != nil {
			t.Errorf("Expected Width=nil (auto-detect), got %v", result.Overlay.Placement.Width)
		}
		if result.Overlay.Placement.Height == nil || *result.Overlay.Placement.Height != 580 {
			t.Errorf("Expected Height=580, got %v", result.Overlay.Placement.Height)
		}
	})

	t.Run("Complex 4-level merge with partial placements", func(t *testing.T) {
		merger := NewConfigMerger()

		// Default config
		defaultConfig := getDefaultConfig()

		// Global config settings (sets Y)
		globalSettings := &ConfigSettings{
			Overlay: &OverlayConfigSettings{
				Placement: &PlacementSettings{
					Y: &[]int{75}[0], // Set Y = 75
				},
			},
		}

		// Type config settings (sets X and height)
		typeSettings := &ConfigSettings{
			Overlay: &OverlayConfigSettings{
				Visible: &[]bool{true}[0],
				Image:   &[]string{"cover.jpg"}[0],
				Placement: &PlacementSettings{
					X:      &[]int{100}[0], // Set X = 100
					Height: &[]int{580}[0], // Set height = 580
				},
			},
		}

		// Front matter (sets X to different value)
		ogpFM := &OGPFrontMatter{
			Overlay: &ArticleOverlayConfig{
				Placement: &PlacementSettings{
					X: &[]int{120}[0], // Override X to 120
					// Y is nil (not specified) - should preserve global Y = 75
				},
			},
		}

		// Apply 4-level merge
		result := merger.MergeConfigsWithSettings(defaultConfig, globalSettings, typeSettings, ogpFM)

		// Verify final placement
		if result.Overlay.Placement.X != 120 {
			t.Errorf("Expected X=120 (from front matter), got %d", result.Overlay.Placement.X)
		}
		if result.Overlay.Placement.Y != 75 {
			t.Errorf("Expected Y=75 (from global), got %d", result.Overlay.Placement.Y)
		}
		if result.Overlay.Placement.Width != nil {
			t.Errorf("Expected Width=nil (auto-detect), got %v", result.Overlay.Placement.Width)
		}
		if result.Overlay.Placement.Height == nil || *result.Overlay.Placement.Height != 580 {
			t.Errorf("Expected Height=580 (from type), got %v", result.Overlay.Placement.Height)
		}
		if !result.Overlay.Visible {
			t.Error("Expected overlay to be visible")
		}
		if result.Overlay.Image == nil || *result.Overlay.Image != "cover.jpg" {
			t.Errorf("Expected overlay image=cover.jpg, got %v", result.Overlay.Image)
		}
	})
}
