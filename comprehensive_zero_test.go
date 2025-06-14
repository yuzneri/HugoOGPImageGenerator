package main

import (
	"testing"
)

// TestComprehensiveZeroValueHandling tests all parameter types to ensure explicit zero values are preserved
func TestComprehensiveZeroValueHandling(t *testing.T) {
	t.Run("All numeric zero values in overlay", func(t *testing.T) {
		merger := NewConfigMerger()

		// Start with default config
		defaultConfig := getDefaultConfig()

		// Type config with explicit zeros
		typeSettings := &ConfigSettings{
			Overlay: &OverlayConfigSettings{
				Visible: &[]bool{true}[0],
				Image:   &[]string{"test.jpg"}[0],
				Opacity: &[]float64{0.0}[0], // Explicit zero opacity (invisible)
				Placement: &PlacementSettings{
					X:      &[]int{0}[0], // Explicit zero X
					Y:      &[]int{0}[0], // Explicit zero Y
					Width:  &[]int{0}[0], // Explicit zero width
					Height: &[]int{0}[0], // Explicit zero height
				},
			},
		}

		// Apply settings
		result := merger.MergeConfigsWithSettings(defaultConfig, nil, typeSettings, nil)

		// Verify all explicit zeros are preserved
		if !result.Overlay.Visible {
			t.Error("Expected overlay to be visible")
		}
		if result.Overlay.Opacity != 0.0 {
			t.Errorf("Expected Opacity=0.0 (explicit zero), got %f", result.Overlay.Opacity)
		}
		if result.Overlay.Placement.X != 0 {
			t.Errorf("Expected X=0 (explicit zero), got %d", result.Overlay.Placement.X)
		}
		if result.Overlay.Placement.Y != 0 {
			t.Errorf("Expected Y=0 (explicit zero), got %d", result.Overlay.Placement.Y)
		}
		if result.Overlay.Placement.Width == nil || *result.Overlay.Placement.Width != 0 {
			t.Errorf("Expected Width=0 (explicit zero), got %v", result.Overlay.Placement.Width)
		}
		if result.Overlay.Placement.Height == nil || *result.Overlay.Placement.Height != 0 {
			t.Errorf("Expected Height=0 (explicit zero), got %v", result.Overlay.Placement.Height)
		}
	})

	t.Run("Text configuration with zero values", func(t *testing.T) {
		merger := NewConfigMerger()

		// Start with default config
		defaultConfig := getDefaultConfig()

		// Type config with explicit zeros for text
		typeSettings := &ConfigSettings{
			Title: &TextSettings{
				Size:          &[]float64{0.0}[0], // Explicit zero size
				LetterSpacing: &[]int{0}[0],       // Explicit zero letter spacing
				LineHeight:    &[]float64{0.0}[0], // Explicit zero line height
				MinSize:       &[]float64{0.0}[0], // Explicit zero min size
				Area: &TextAreaSettings{
					X:      &[]int{0}[0], // Explicit zero X
					Y:      &[]int{0}[0], // Explicit zero Y
					Width:  &[]int{0}[0], // Explicit zero width
					Height: &[]int{0}[0], // Explicit zero height
				},
			},
		}

		// Apply settings
		result := merger.MergeConfigsWithSettings(defaultConfig, nil, typeSettings, nil)

		// Verify all explicit zeros are preserved
		if result.Title.Size != 0.0 {
			t.Errorf("Expected Size=0.0 (explicit zero), got %f", result.Title.Size)
		}
		if result.Title.LetterSpacing != 0 {
			t.Errorf("Expected LetterSpacing=0 (explicit zero), got %d", result.Title.LetterSpacing)
		}
		if result.Title.LineHeight != 0.0 {
			t.Errorf("Expected LineHeight=0.0 (explicit zero), got %f", result.Title.LineHeight)
		}
		if result.Title.MinSize != 0.0 {
			t.Errorf("Expected MinSize=0.0 (explicit zero), got %f", result.Title.MinSize)
		}
		if result.Title.Area.X != 0 {
			t.Errorf("Expected Area.X=0 (explicit zero), got %d", result.Title.Area.X)
		}
		if result.Title.Area.Y != 0 {
			t.Errorf("Expected Area.Y=0 (explicit zero), got %d", result.Title.Area.Y)
		}
		if result.Title.Area.Width != 0 {
			t.Errorf("Expected Area.Width=0 (explicit zero), got %d", result.Title.Area.Width)
		}
		if result.Title.Area.Height != 0 {
			t.Errorf("Expected Area.Height=0 (explicit zero), got %d", result.Title.Area.Height)
		}
	})

	t.Run("String values - empty vs nil", func(t *testing.T) {
		merger := NewConfigMerger()

		// Start with default config
		defaultConfig := getDefaultConfig()

		// Type config with explicit empty strings
		typeSettings := &ConfigSettings{
			Title: &TextSettings{
				Content:       &[]string{""}[0], // Explicit empty content
				Font:          &[]string{""}[0], // Explicit empty font path
				Color:         &[]string{""}[0], // Explicit empty color
				BlockPosition: &[]string{""}[0], // Explicit empty position
				LineAlignment: &[]string{""}[0], // Explicit empty alignment
				Overflow:      &[]string{""}[0], // Explicit empty overflow
				LineBreaking: &LineBreakingSettings{
					StartProhibited: &[]string{""}[0], // Explicit empty start prohibited
					EndProhibited:   &[]string{""}[0], // Explicit empty end prohibited
				},
			},
		}

		// Apply settings
		result := merger.MergeConfigsWithSettings(defaultConfig, nil, typeSettings, nil)

		// Verify explicit empty strings are preserved
		if result.Title.Content == nil || *result.Title.Content != "" {
			t.Errorf("Expected Content='' (explicit empty), got %v", result.Title.Content)
		}
		if result.Title.Font == nil || *result.Title.Font != "" {
			t.Errorf("Expected Font='' (explicit empty), got %v", result.Title.Font)
		}
		if result.Title.Color != "" {
			t.Errorf("Expected Color='' (explicit empty), got %s", result.Title.Color)
		}
		if result.Title.BlockPosition != "" {
			t.Errorf("Expected BlockPosition='' (explicit empty), got %s", result.Title.BlockPosition)
		}
		if result.Title.LineAlignment != "" {
			t.Errorf("Expected LineAlignment='' (explicit empty), got %s", result.Title.LineAlignment)
		}
		if result.Title.Overflow != "" {
			t.Errorf("Expected Overflow='' (explicit empty), got %s", result.Title.Overflow)
		}
		if result.Title.LineBreaking.StartProhibited != "" {
			t.Errorf("Expected StartProhibited='' (explicit empty), got %s", result.Title.LineBreaking.StartProhibited)
		}
		if result.Title.LineBreaking.EndProhibited != "" {
			t.Errorf("Expected EndProhibited='' (explicit empty), got %s", result.Title.LineBreaking.EndProhibited)
		}
	})

	t.Run("Boolean values - explicit false", func(t *testing.T) {
		merger := NewConfigMerger()

		// Start with config where overlay is visible
		defaultConfig := getDefaultConfig()
		defaultConfig.Overlay.Visible = true

		// Type config with explicit false
		typeSettings := &ConfigSettings{
			Title: &TextSettings{
				Visible: &[]bool{false}[0], // Explicit false
			},
			Description: &TextSettings{
				Visible: &[]bool{false}[0], // Explicit false
			},
			Overlay: &OverlayConfigSettings{
				Visible: &[]bool{false}[0], // Explicit false (hide overlay)
			},
		}

		// Apply settings
		result := merger.MergeConfigsWithSettings(defaultConfig, nil, typeSettings, nil)

		// Verify explicit false values are preserved
		if result.Title.Visible {
			t.Error("Expected Title.Visible=false (explicit false), got true")
		}
		if result.Description.Visible {
			t.Error("Expected Description.Visible=false (explicit false), got true")
		}
		if result.Overlay.Visible {
			t.Error("Expected Overlay.Visible=false (explicit false), got true")
		}
	})

	t.Run("Front matter overrides with zero values", func(t *testing.T) {
		merger := NewConfigMerger()

		// Start with config that has non-zero values
		config := getDefaultConfig()
		config.Overlay.Visible = true
		config.Overlay.Image = &[]string{"test.jpg"}[0]
		config.Overlay.Opacity = 0.8
		config.Overlay.Placement.X = 100
		config.Overlay.Placement.Y = 200

		// Front matter with explicit zeros
		ogpFM := &OGPFrontMatter{
			Title: &TextConfigOverride{
				Visible:       &[]bool{false}[0],  // Explicit false
				Size:          &[]float64{0.0}[0], // Explicit zero size
				LetterSpacing: &[]int{0}[0],       // Explicit zero spacing
			},
			Overlay: &ArticleOverlayConfig{
				Opacity: &[]float64{0.0}[0], // Explicit zero opacity
				Placement: &PlacementSettings{
					X: &[]int{0}[0], // Explicit zero X
					Y: &[]int{0}[0], // Explicit zero Y
				},
			},
		}

		// Apply front matter overrides
		result := merger.applyFrontMatterOverrides(config, ogpFM)

		// Verify all explicit zeros and false values are preserved
		if result.Title.Visible {
			t.Error("Expected Title.Visible=false (explicit false), got true")
		}
		if result.Title.Size != 0.0 {
			t.Errorf("Expected Title.Size=0.0 (explicit zero), got %f", result.Title.Size)
		}
		if result.Title.LetterSpacing != 0 {
			t.Errorf("Expected Title.LetterSpacing=0 (explicit zero), got %d", result.Title.LetterSpacing)
		}
		if result.Overlay.Opacity != 0.0 {
			t.Errorf("Expected Overlay.Opacity=0.0 (explicit zero), got %f", result.Overlay.Opacity)
		}
		if result.Overlay.Placement.X != 0 {
			t.Errorf("Expected Overlay.Placement.X=0 (explicit zero), got %d", result.Overlay.Placement.X)
		}
		if result.Overlay.Placement.Y != 0 {
			t.Errorf("Expected Overlay.Placement.Y=0 (explicit zero), got %d", result.Overlay.Placement.Y)
		}
	})
}
