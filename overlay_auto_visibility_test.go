package main

import (
	"testing"
)

// TestOverlayAutoVisibility tests the auto-visibility feature for overlays.
// When a type config provides an overlay image but doesn't explicitly set visible,
// the overlay should automatically become visible.
func TestOverlayAutoVisibility(t *testing.T) {
	tests := []struct {
		name          string
		typeConfig    *ConfigSettings
		frontMatter   *OGPFrontMatter
		expectVisible bool
		description   string
	}{
		{
			name: "Type config with image but no visible field should auto-enable",
			typeConfig: &ConfigSettings{
				Overlay: &OverlayConfigSettings{
					Image: stringPtr("book-overlay.png"),
					Placement: &PlacementSettings{
						X:      intPtr(25),
						Y:      intPtr(25),
						Height: intPtr(580),
					},
				},
			},
			frontMatter:   &OGPFrontMatter{},
			expectVisible: true,
			description:   "Auto-visibility when image provided in type config",
		},
		{
			name: "Type config with explicit visible: true should remain visible",
			typeConfig: &ConfigSettings{
				Overlay: &OverlayConfigSettings{
					Visible: boolPtr(true),
					Image:   stringPtr("book-overlay.png"),
				},
			},
			frontMatter:   &OGPFrontMatter{},
			expectVisible: true,
			description:   "Explicit visible: true should be preserved",
		},
		{
			name: "Type config with explicit visible: false should remain hidden",
			typeConfig: &ConfigSettings{
				Overlay: &OverlayConfigSettings{
					Visible: boolPtr(false),
					Image:   stringPtr("book-overlay.png"),
				},
			},
			frontMatter:   &OGPFrontMatter{},
			expectVisible: false,
			description:   "Explicit visible: false should override auto-visibility",
		},
		{
			name: "Type config with empty image should not auto-enable",
			typeConfig: &ConfigSettings{
				Overlay: &OverlayConfigSettings{
					Image: stringPtr(""),
					Placement: &PlacementSettings{
						X: intPtr(25),
						Y: intPtr(25),
					},
				},
			},
			frontMatter:   &OGPFrontMatter{},
			expectVisible: false,
			description:   "Empty image should not trigger auto-visibility",
		},
		{
			name: "Front matter can override type config auto-visibility",
			typeConfig: &ConfigSettings{
				Overlay: &OverlayConfigSettings{
					Image: stringPtr("book-overlay.png"),
				},
			},
			frontMatter: &OGPFrontMatter{
				Overlay: &ArticleOverlayConfig{
					Visible: boolPtr(false),
				},
			},
			expectVisible: false,
			description:   "Front matter visible: false should override type config auto-visibility",
		},
		{
			name: "Type config placement without image should not auto-enable",
			typeConfig: &ConfigSettings{
				Overlay: &OverlayConfigSettings{
					Placement: &PlacementSettings{
						X:      intPtr(25),
						Y:      intPtr(25),
						Height: intPtr(580),
					},
				},
			},
			frontMatter:   &OGPFrontMatter{},
			expectVisible: false,
			description:   "Placement settings without image should not trigger auto-visibility",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			merger := NewConfigMerger()
			defaultConfig := getDefaultConfig()

			finalConfig := merger.MergeConfigsWithSettings(
				defaultConfig,
				nil,            // no global settings
				tt.typeConfig,  // type config
				tt.frontMatter, // front matter
			)

			if finalConfig.Overlay.Visible != tt.expectVisible {
				t.Errorf("Expected overlay visible to be %v, got %v. %s",
					tt.expectVisible, finalConfig.Overlay.Visible, tt.description)
			}

			// Verify that other overlay settings are preserved
			if tt.typeConfig.Overlay != nil && tt.typeConfig.Overlay.Image != nil {
				expectedImage := *tt.typeConfig.Overlay.Image
				if finalConfig.Overlay.Image == nil {
					if expectedImage != "" {
						t.Errorf("Expected overlay image to be preserved, got nil")
					}
				} else if *finalConfig.Overlay.Image != expectedImage {
					t.Errorf("Expected overlay image %q, got %q", expectedImage, *finalConfig.Overlay.Image)
				}
			}
		})
	}
}
