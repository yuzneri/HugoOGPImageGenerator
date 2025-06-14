package main

import (
	"testing"
)

// TestApplySettingsToConfig tests the core merge logic
func TestApplySettingsToConfig(t *testing.T) {
	tests := []struct {
		name       string
		baseConfig *Config
		settings   *ConfigSettings
		expected   *Config
	}{
		{
			name: "Partial title settings application",
			baseConfig: &Config{
				Title: TextConfig{
					Visible: true,
					Size:    32.0,
					Color:   "#000000",
					Area: TextArea{
						X: 0, Y: 0, Width: 800, Height: 200,
					},
				},
			},
			settings: &ConfigSettings{
				Title: &TextSettings{
					Size:  float64Ptr(48.0),
					Color: stringPtr("#FF0000"),
				},
			},
			expected: &Config{
				Title: TextConfig{
					Visible: true,      // 保持
					Size:    48.0,      // 上書き
					Color:   "#FF0000", // 上書き
					Area: TextArea{ // 保持
						X: 0, Y: 0, Width: 800, Height: 200,
					},
				},
			},
		},
		{
			name: "Partial placement settings application",
			baseConfig: &Config{
				Overlay: MainOverlayConfig{
					Visible: true,
					Placement: PlacementConfig{
						X: 50, Y: 50,
						Width:  nil,
						Height: nil,
					},
				},
			},
			settings: &ConfigSettings{
				Overlay: &OverlayConfigSettings{
					Placement: &PlacementSettings{
						X:      intPtr(100),
						Height: intPtr(580),
					},
				},
			},
			expected: &Config{
				Overlay: MainOverlayConfig{
					Visible: true, // 保持
					Placement: PlacementConfig{
						X:      100,         // 上書き
						Y:      50,          // 保持
						Width:  nil,         // 保持
						Height: intPtr(580), // 上書き
					},
				},
			},
		},
		{
			name: "Multiple sections application",
			baseConfig: &Config{
				Title: TextConfig{
					Visible: false,
					Size:    24.0,
				},
				Description: TextConfig{
					Visible: true,
					Color:   "#333333",
				},
				Overlay: MainOverlayConfig{
					Visible: false,
				},
			},
			settings: &ConfigSettings{
				Title: &TextSettings{
					Visible: boolPtr(true),
				},
				Description: &TextSettings{
					Size: float64Ptr(16.0),
				},
				Overlay: &OverlayConfigSettings{
					Visible: boolPtr(true),
				},
			},
			expected: &Config{
				Title: TextConfig{
					Visible: true, // 上書き
					Size:    24.0, // 保持
				},
				Description: TextConfig{
					Visible: true,      // 保持
					Color:   "#333333", // 保持
					Size:    16.0,      // 上書き
				},
				Overlay: MainOverlayConfig{
					Visible: true, // 上書き
				},
			},
		},
		{
			name: "Nil settings should not change anything",
			baseConfig: &Config{
				Title: TextConfig{
					Visible: true,
					Size:    32.0,
				},
			},
			settings: nil,
			expected: &Config{
				Title: TextConfig{
					Visible: true,
					Size:    32.0,
				},
			},
		},
		{
			name: "Empty settings should not change anything",
			baseConfig: &Config{
				Title: TextConfig{
					Visible: true,
					Size:    32.0,
				},
			},
			settings: &ConfigSettings{},
			expected: &Config{
				Title: TextConfig{
					Visible: true,
					Size:    32.0,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merger := NewConfigMerger()

			// Make a copy of baseConfig to avoid modifying the original
			target := deepCopyConfig(tt.baseConfig)

			// Apply settings
			merger.applySettingsToConfig(target, tt.settings)

			// Compare results
			if !configsEqual(target, tt.expected) {
				t.Errorf("Result doesn't match expected.\nGot: %+v\nExpected: %+v", target, tt.expected)
			}
		})
	}
}

// TestMergeConfigsWithSettings tests the full 4-level hierarchy
func TestMergeConfigsWithSettings(t *testing.T) {
	// Default config
	defaultConfig := &Config{
		Title: TextConfig{
			Visible: true,
			Size:    32.0,
			Color:   "#000000",
		},
		Overlay: MainOverlayConfig{
			Visible: true,
			Placement: PlacementConfig{
				X: 50, Y: 50,
			},
		},
	}

	// Global settings
	globalSettings := &ConfigSettings{
		Title: &TextSettings{
			Size:  float64Ptr(48.0),
			Color: stringPtr("#333333"),
		},
	}

	// Type settings
	typeSettings := &ConfigSettings{
		Title: &TextSettings{
			Color: stringPtr("#FF0000"), // This should override global
		},
		Overlay: &OverlayConfigSettings{
			Placement: &PlacementSettings{
				X: intPtr(100),
			},
		},
	}

	// Front matter (nil in this test)
	ogpFM := (*OGPFrontMatter)(nil)

	// Expected result after applying hierarchy
	expected := &Config{
		Title: TextConfig{
			Visible: true,      // from default
			Size:    48.0,      // from global
			Color:   "#FF0000", // from type (overrides global)
		},
		Overlay: MainOverlayConfig{
			Visible: true, // from default
			Placement: PlacementConfig{
				X: 100, // from type
				Y: 50,  // from default
			},
		},
	}

	merger := NewConfigMerger()
	result := merger.MergeConfigsWithSettings(defaultConfig, globalSettings, typeSettings, ogpFM)

	if !configsEqual(result, expected) {
		t.Errorf("4-level hierarchy result doesn't match expected.\nGot: %+v\nExpected: %+v", result, expected)
	}
}

// TestPartialPlacementMerge tests the specific case mentioned in the issue
func TestPartialPlacementMerge(t *testing.T) {
	// Base config with default placement
	baseConfig := &Config{
		Overlay: MainOverlayConfig{
			Visible: true,
			Placement: PlacementConfig{
				X: 50, Y: 50,
				Width:  nil,
				Height: nil,
			},
		},
	}

	// Settings with only X and Height specified
	settings := &ConfigSettings{
		Overlay: &OverlayConfigSettings{
			Placement: &PlacementSettings{
				X:      intPtr(100),
				Height: intPtr(580),
				// Y and Width are nil - should not override
			},
		},
	}

	// Expected result: only X and Height should change
	expected := &Config{
		Overlay: MainOverlayConfig{
			Visible: true,
			Placement: PlacementConfig{
				X:      100,         // changed
				Y:      50,          // preserved
				Width:  nil,         // preserved
				Height: intPtr(580), // changed
			},
		},
	}

	merger := NewConfigMerger()
	result := deepCopyConfig(baseConfig)
	merger.applySettingsToConfig(result, settings)

	if !configsEqual(result, expected) {
		t.Errorf("Partial placement merge failed.\nGot: %+v\nExpected: %+v", result, expected)
	}
}

// TestPointerFieldHandling tests proper handling of pointer fields
func TestPointerFieldHandling(t *testing.T) {
	tests := []struct {
		name       string
		baseConfig *Config
		settings   *ConfigSettings
		expected   *Config
	}{
		{
			name: "Setting Content pointer field",
			baseConfig: &Config{
				Title: TextConfig{Content: nil},
			},
			settings: &ConfigSettings{
				Title: &TextSettings{
					Content: stringPtr("{{.Title}}"),
				},
			},
			expected: &Config{
				Title: TextConfig{Content: stringPtr("{{.Title}}")},
			},
		},
		{
			name: "Clearing Content pointer field",
			baseConfig: &Config{
				Title: TextConfig{Content: stringPtr("existing")},
			},
			settings: &ConfigSettings{
				Title: &TextSettings{
					Content: stringPtr(""),
				},
			},
			expected: &Config{
				Title: TextConfig{Content: stringPtr("")},
			},
		},
		{
			name: "Not setting Content pointer field should preserve existing",
			baseConfig: &Config{
				Title: TextConfig{Content: stringPtr("existing")},
			},
			settings: &ConfigSettings{
				Title: &TextSettings{
					Size: float64Ptr(48.0),
					// Content is nil - should not change existing
				},
			},
			expected: &Config{
				Title: TextConfig{
					Content: stringPtr("existing"), // preserved
					Size:    48.0,                  // changed
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merger := NewConfigMerger()
			result := deepCopyConfig(tt.baseConfig)
			merger.applySettingsToConfig(result, tt.settings)

			if !configsEqual(result, tt.expected) {
				t.Errorf("Pointer field handling failed.\nGot: %+v\nExpected: %+v", result, tt.expected)
			}
		})
	}
}

// Helper function to deep copy Config for testing
func deepCopyConfig(original *Config) *Config {
	if original == nil {
		return nil
	}

	copied := *original

	// Deep copy pointer fields
	if original.Background.Image != nil {
		copied.Background.Image = stringPtr(*original.Background.Image)
	}
	if original.Title.Content != nil {
		copied.Title.Content = stringPtr(*original.Title.Content)
	}
	if original.Title.Font != nil {
		copied.Title.Font = stringPtr(*original.Title.Font)
	}
	if original.Description.Content != nil {
		copied.Description.Content = stringPtr(*original.Description.Content)
	}
	if original.Description.Font != nil {
		copied.Description.Font = stringPtr(*original.Description.Font)
	}
	if original.Overlay.Image != nil {
		copied.Overlay.Image = stringPtr(*original.Overlay.Image)
	}
	if original.Overlay.Placement.Width != nil {
		copied.Overlay.Placement.Width = intPtr(*original.Overlay.Placement.Width)
	}
	if original.Overlay.Placement.Height != nil {
		copied.Overlay.Placement.Height = intPtr(*original.Overlay.Placement.Height)
	}

	return &copied
}

// Helper function to compare configs (simplified for testing)
func configsEqual(a, b *Config) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	// Compare basic fields
	if a.Title.Visible != b.Title.Visible ||
		a.Title.Size != b.Title.Size ||
		a.Title.Color != b.Title.Color {
		return false
	}

	// Compare pointer fields
	if !stringPtrEqual(a.Title.Content, b.Title.Content) ||
		!stringPtrEqual(a.Title.Font, b.Title.Font) {
		return false
	}

	// Compare placement
	if a.Overlay.Placement.X != b.Overlay.Placement.X || a.Overlay.Placement.Y != b.Overlay.Placement.Y {
		return false
	}

	if !intPtrEqual(a.Overlay.Placement.Width, b.Overlay.Placement.Width) ||
		!intPtrEqual(a.Overlay.Placement.Height, b.Overlay.Placement.Height) {
		return false
	}

	return true
}

func stringPtrEqual(a, b *string) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

func intPtrEqual(a, b *int) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// Helper function for creating int pointers
func intPtr(i int) *int { return &i }
