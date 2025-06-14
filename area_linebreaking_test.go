package main

import (
	"testing"
)

// TestAreaAndLineBreakingPartialMerging tests that partial area and line_breaking configurations are merged correctly
func TestAreaAndLineBreakingPartialMerging(t *testing.T) {
	t.Run("Partial area settings from type config", func(t *testing.T) {
		merger := NewConfigMerger()

		// Start with default config
		config := getDefaultConfig()

		// Type config with only width and height specified (X, Y should remain default)
		typeSettings := &ConfigSettings{
			Title: &TextSettings{
				Area: &TextAreaSettings{
					Width:  &[]int{800}[0], // Specify width = 800
					Height: &[]int{150}[0], // Specify height = 150
					// X and Y are nil (should remain as defaults)
				},
			},
		}

		// Apply type settings
		merger.applySettingsToConfig(config, typeSettings)

		// Verify area
		if config.Title.Area.X != DefaultTitleAreaX {
			t.Errorf("Expected X=%d (default), got %d", DefaultTitleAreaX, config.Title.Area.X)
		}
		if config.Title.Area.Y != DefaultTitleAreaY {
			t.Errorf("Expected Y=%d (default), got %d", DefaultTitleAreaY, config.Title.Area.Y)
		}
		if config.Title.Area.Width != 800 {
			t.Errorf("Expected Width=800, got %d", config.Title.Area.Width)
		}
		if config.Title.Area.Height != 150 {
			t.Errorf("Expected Height=150, got %d", config.Title.Area.Height)
		}
	})

	t.Run("Partial area settings from front matter", func(t *testing.T) {
		merger := NewConfigMerger()

		// Start with config that has custom area settings
		config := getDefaultConfig()
		config.Title.Area.X = 100
		config.Title.Area.Y = 200
		config.Title.Area.Width = 1000
		config.Title.Area.Height = 300

		// Front matter with only X and width specified (Y, height should remain unchanged)
		ogpFM := &OGPFrontMatter{
			Title: &TextConfigOverride{
				Area: &TextAreaConfig{
					X:     &[]int{150}[0], // Override X
					Width: &[]int{900}[0], // Override width
					// Y and Height are nil (should remain unchanged)
				},
			},
		}

		// Apply front matter overrides
		result := merger.applyFrontMatterOverrides(config, ogpFM)

		// Verify area
		if result.Title.Area.X != 150 {
			t.Errorf("Expected X=150, got %d", result.Title.Area.X)
		}
		if result.Title.Area.Y != 200 { // Should remain original value
			t.Errorf("Expected Y=200 (preserved), got %d", result.Title.Area.Y)
		}
		if result.Title.Area.Width != 900 {
			t.Errorf("Expected Width=900, got %d", result.Title.Area.Width)
		}
		if result.Title.Area.Height != 300 { // Should remain original value
			t.Errorf("Expected Height=300 (preserved), got %d", result.Title.Area.Height)
		}
	})

	t.Run("Partial line breaking settings from type config", func(t *testing.T) {
		merger := NewConfigMerger()

		// Start with default config
		config := getDefaultConfig()

		// Type config with only start_prohibited specified
		typeSettings := &ConfigSettings{
			Description: &TextSettings{
				LineBreaking: &LineBreakingSettings{
					StartProhibited: &[]string{"、。！？"}[0],
					// EndProhibited is nil (should remain as default)
				},
			},
		}

		// Apply type settings
		merger.applySettingsToConfig(config, typeSettings)

		// Verify line breaking
		if config.Description.LineBreaking.StartProhibited != "、。！？" {
			t.Errorf("Expected StartProhibited='、。！？', got %s", config.Description.LineBreaking.StartProhibited)
		}
		if config.Description.LineBreaking.EndProhibited != DefaultEndProhibitedChars {
			t.Errorf("Expected EndProhibited to remain default %s, got %s", DefaultEndProhibitedChars, config.Description.LineBreaking.EndProhibited)
		}
	})

	t.Run("Partial line breaking settings from front matter", func(t *testing.T) {
		merger := NewConfigMerger()

		// Start with config that has custom line breaking
		config := getDefaultConfig()
		config.Title.LineBreaking.StartProhibited = "ABC"
		config.Title.LineBreaking.EndProhibited = "XYZ"

		// Front matter with only end_prohibited specified
		ogpFM := &OGPFrontMatter{
			Title: &TextConfigOverride{
				LineBreaking: &LineBreakingOverride{
					EndProhibited: &[]string{"DEF"}[0],
					// StartProhibited is nil (should remain unchanged)
				},
			},
		}

		// Apply front matter overrides
		result := merger.applyFrontMatterOverrides(config, ogpFM)

		// Verify line breaking
		if result.Title.LineBreaking.StartProhibited != "ABC" { // Should remain original
			t.Errorf("Expected StartProhibited='ABC' (preserved), got %s", result.Title.LineBreaking.StartProhibited)
		}
		if result.Title.LineBreaking.EndProhibited != "DEF" {
			t.Errorf("Expected EndProhibited='DEF', got %s", result.Title.LineBreaking.EndProhibited)
		}
	})

	t.Run("Complex 4-level merge with partial area and line breaking", func(t *testing.T) {
		merger := NewConfigMerger()

		// Default config
		defaultConfig := getDefaultConfig()

		// Global config settings (sets area width)
		globalSettings := &ConfigSettings{
			Title: &TextSettings{
				Area: &TextAreaSettings{
					Width: &[]int{900}[0], // Set width = 900
				},
			},
		}

		// Type config settings (sets area Y and line breaking start)
		typeSettings := &ConfigSettings{
			Title: &TextSettings{
				Area: &TextAreaSettings{
					Y: &[]int{80}[0], // Set Y = 80
				},
				LineBreaking: &LineBreakingSettings{
					StartProhibited: &[]string{"Custom"}[0],
				},
			},
		}

		// Front matter (sets area X and line breaking end)
		ogpFM := &OGPFrontMatter{
			Title: &TextConfigOverride{
				Area: &TextAreaConfig{
					X: &[]int{120}[0], // Override X to 120
				},
				LineBreaking: &LineBreakingOverride{
					EndProhibited: &[]string{"Override"}[0],
				},
			},
		}

		// Apply 4-level merge
		result := merger.MergeConfigsWithSettings(defaultConfig, globalSettings, typeSettings, ogpFM)

		// Verify final area (should combine all levels)
		if result.Title.Area.X != 120 {
			t.Errorf("Expected X=120 (from front matter), got %d", result.Title.Area.X)
		}
		if result.Title.Area.Y != 80 {
			t.Errorf("Expected Y=80 (from type), got %d", result.Title.Area.Y)
		}
		if result.Title.Area.Width != 900 {
			t.Errorf("Expected Width=900 (from global), got %d", result.Title.Area.Width)
		}
		if result.Title.Area.Height != DefaultTitleAreaHeight {
			t.Errorf("Expected Height=%d (default), got %d", DefaultTitleAreaHeight, result.Title.Area.Height)
		}

		// Verify final line breaking (should combine all levels)
		if result.Title.LineBreaking.StartProhibited != "Custom" {
			t.Errorf("Expected StartProhibited='Custom' (from type), got %s", result.Title.LineBreaking.StartProhibited)
		}
		if result.Title.LineBreaking.EndProhibited != "Override" {
			t.Errorf("Expected EndProhibited='Override' (from front matter), got %s", result.Title.LineBreaking.EndProhibited)
		}
	})
}
