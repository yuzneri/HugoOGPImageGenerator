package main

import (
	"os"
	"path/filepath"
	"testing"
)

// TestOverlayIntegration tests the complete overlay functionality with real file operations
func TestOverlayIntegration(t *testing.T) {
	// Create temporary directory structure
	tempDir, err := os.MkdirTemp("", "overlay_integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configDir := tempDir
	contentDir := filepath.Join(tempDir, "content")
	articleDir := filepath.Join(contentDir, "books", "test-book")

	err = os.MkdirAll(articleDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create article dir: %v", err)
	}

	// Create article-specific cover image (simple PNG-like content)
	articleCoverPath := filepath.Join(articleDir, "cover.png")
	createTestImage(t, articleCoverPath)

	// Create config-level default cover
	configCoverPath := filepath.Join(configDir, "default-cover.png")
	createTestImage(t, configCoverPath)

	// Create books.yaml type configuration
	booksConfigPath := filepath.Join(configDir, "books.yaml")
	booksConfigContent := `overlay:
  visible: true
  image: "cover.png"
  placement:
    x: 100
    height: 580
    # y not specified - should use default (50)
    # width not specified - should auto-detect
  fit: "contain"
  opacity: 1.0`

	err = os.WriteFile(booksConfigPath, []byte(booksConfigContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create books config: %v", err)
	}

	// Create article with front matter
	indexContent := `---
title: "Test Book"
description: "A test book for overlay integration"
ogp:
  title:
    area:
      width: 800
      # height not specified - should use default
    line_breaking:
      start_prohibited: "Custom"
      # end_prohibited not specified - should use default
---

# Test Book Content
`
	indexPath := filepath.Join(articleDir, "index.md")
	err = os.WriteFile(indexPath, []byte(indexContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create index file: %v", err)
	}

	t.Run("End-to-end overlay with partial settings", func(t *testing.T) {
		// Create OGP generator
		generator, err := NewOGPGenerator("", contentDir, configDir)
		if err != nil {
			t.Fatalf("Failed to create OGP generator: %v", err)
		}

		// Generate test image
		err = generator.GenerateTest(articleDir)
		if err != nil {
			t.Errorf("GenerateTest should not return error: %v", err)
		}

		// Verify output file exists
		outputPath := generator.articleProcessor.generateTestOutputPath()
		if _, err := os.Stat(outputPath); os.IsNotExist(err) {
			t.Error("Expected output file to be created")
		}

		// Clean up
		os.Remove(outputPath)
	})

	t.Run("Type config overlay asset resolution", func(t *testing.T) {
		// Load type configuration
		typeSettings, err := loadTypeConfigSettings(configDir, "books")
		if err != nil {
			t.Fatalf("Failed to load type config: %v", err)
		}

		if typeSettings == nil || typeSettings.Overlay == nil {
			t.Fatal("Type config should have overlay settings")
		}

		// Verify overlay image path
		if typeSettings.Overlay.Image == nil || *typeSettings.Overlay.Image != "cover.png" {
			t.Errorf("Expected overlay image 'cover.png', got %v", typeSettings.Overlay.Image)
		}

		// Verify placement settings
		if typeSettings.Overlay.Placement == nil {
			t.Fatal("Type config should have placement settings")
		}

		placement := typeSettings.Overlay.Placement
		if placement.X == nil || *placement.X != 100 {
			t.Errorf("Expected X=100, got %v", placement.X)
		}
		if placement.Y != nil {
			t.Errorf("Expected Y=nil (unspecified), got %v", placement.Y)
		}
		if placement.Width != nil {
			t.Errorf("Expected Width=nil (auto-detect), got %v", placement.Width)
		}
		if placement.Height == nil || *placement.Height != 580 {
			t.Errorf("Expected Height=580, got %v", placement.Height)
		}
	})

	t.Run("Asset resolution priority verification", func(t *testing.T) {
		resolver := NewPathResolver(configDir)

		// Test cover.png resolution - should find article-specific version
		resolvedPath := resolver.ResolveAssetPath("cover.png", articleDir)
		if resolvedPath != articleCoverPath {
			t.Errorf("Expected article-specific cover %s, got %s", articleCoverPath, resolvedPath)
		}

		// Test non-existent file - should fallback to config directory
		resolvedPath = resolver.ResolveAssetPath("nonexistent.png", articleDir)
		expectedConfigPath := filepath.Join(configDir, "nonexistent.png")
		if resolvedPath != expectedConfigPath {
			t.Errorf("Expected config fallback %s, got %s", expectedConfigPath, resolvedPath)
		}
	})

	t.Run("Fallback to config directory when article asset missing", func(t *testing.T) {
		// Remove article cover temporarily
		os.Rename(articleCoverPath, articleCoverPath+".backup")
		defer os.Rename(articleCoverPath+".backup", articleCoverPath)

		resolver := NewPathResolver(configDir)
		resolvedPath := resolver.ResolveAssetPath("cover.png", articleDir)

		// Should fallback to config directory (even though file doesn't exist there)
		expectedConfigPath := filepath.Join(configDir, "cover.png")
		if resolvedPath != expectedConfigPath {
			t.Errorf("Expected config fallback %s, got %s", expectedConfigPath, resolvedPath)
		}
	})
}

// TestPartialConfigurationMerging tests that partial configurations are merged correctly across all levels
func TestPartialConfigurationMerging(t *testing.T) {
	t.Run("Complex partial configuration scenario", func(t *testing.T) {
		merger := NewConfigMerger()

		// Start with default config
		defaultConfig := getDefaultConfig()

		// Global config (sets some overlay and title area settings)
		globalSettings := &ConfigSettings{
			Overlay: &OverlayConfigSettings{
				Placement: &PlacementSettings{
					Y: &[]int{75}[0], // Set Y = 75
				},
			},
			Title: &TextSettings{
				Area: &TextAreaSettings{
					Width: &[]int{900}[0], // Set width = 900
				},
			},
		}

		// Type config (sets overlay image and more placement settings)
		typeSettings := &ConfigSettings{
			Overlay: &OverlayConfigSettings{
				Visible: &[]bool{true}[0],
				Image:   &[]string{"cover.jpg"}[0],
				Placement: &PlacementSettings{
					X:      &[]int{100}[0], // Set X = 100
					Height: &[]int{580}[0], // Set height = 580
				},
				Opacity: &[]float64{0.9}[0],
			},
			Title: &TextSettings{
				Area: &TextAreaSettings{
					Y: &[]int{80}[0], // Set Y = 80
				},
				LineBreaking: &LineBreakingSettings{
					StartProhibited: &[]string{"Custom"}[0],
				},
			},
		}

		// Front matter (overrides some settings)
		ogpFM := &OGPFrontMatter{
			Overlay: &ArticleOverlayConfig{
				Placement: &PlacementSettings{
					X: &[]int{120}[0], // Override X to 120
					// Y is nil (not specified) - should preserve type/global Y
				},
				Fit: &[]string{"cover"}[0],
			},
			Title: &TextConfigOverride{
				Area: &TextAreaConfig{
					X: &[]int{150}[0], // Override X to 150
					// Other fields nil - should preserve previous values
				},
				LineBreaking: &LineBreakingOverride{
					EndProhibited: &[]string{"Override"}[0],
					// StartProhibited nil - should preserve type value
				},
			},
		}

		// Apply 4-level merge
		result := merger.MergeConfigsWithSettings(defaultConfig, globalSettings, typeSettings, ogpFM)

		// Verify overlay settings
		if !result.Overlay.Visible {
			t.Error("Expected overlay to be visible")
		}
		if result.Overlay.Image == nil || *result.Overlay.Image != "cover.jpg" {
			t.Errorf("Expected overlay image 'cover.jpg', got %v", result.Overlay.Image)
		}
		if result.Overlay.Fit != "cover" {
			t.Errorf("Expected overlay fit 'cover', got %s", result.Overlay.Fit)
		}
		if result.Overlay.Opacity != 0.9 {
			t.Errorf("Expected overlay opacity 0.9, got %f", result.Overlay.Opacity)
		}

		// Verify overlay placement (complex merge from multiple levels)
		if result.Overlay.Placement.X != 120 {
			t.Errorf("Expected overlay X=120 (from front matter), got %d", result.Overlay.Placement.X)
		}
		if result.Overlay.Placement.Y != 75 {
			t.Errorf("Expected overlay Y=75 (from global), got %d", result.Overlay.Placement.Y)
		}
		if result.Overlay.Placement.Width != nil {
			t.Errorf("Expected overlay Width=nil (auto-detect), got %v", result.Overlay.Placement.Width)
		}
		if result.Overlay.Placement.Height == nil || *result.Overlay.Placement.Height != 580 {
			t.Errorf("Expected overlay Height=580 (from type), got %v", result.Overlay.Placement.Height)
		}

		// Verify title area (complex merge from multiple levels)
		if result.Title.Area.X != 150 {
			t.Errorf("Expected title area X=150 (from front matter), got %d", result.Title.Area.X)
		}
		if result.Title.Area.Y != 80 {
			t.Errorf("Expected title area Y=80 (from type), got %d", result.Title.Area.Y)
		}
		if result.Title.Area.Width != 900 {
			t.Errorf("Expected title area Width=900 (from global), got %d", result.Title.Area.Width)
		}
		if result.Title.Area.Height != DefaultTitleAreaHeight {
			t.Errorf("Expected title area Height=%d (default), got %d", DefaultTitleAreaHeight, result.Title.Area.Height)
		}

		// Verify line breaking (complex merge from multiple levels)
		if result.Title.LineBreaking.StartProhibited != "Custom" {
			t.Errorf("Expected StartProhibited='Custom' (from type), got %s", result.Title.LineBreaking.StartProhibited)
		}
		if result.Title.LineBreaking.EndProhibited != "Override" {
			t.Errorf("Expected EndProhibited='Override' (from front matter), got %s", result.Title.LineBreaking.EndProhibited)
		}
	})
}

// TestEdgeCasesPartialConfiguration tests edge cases for partial configuration merging
func TestEdgeCasesPartialConfiguration(t *testing.T) {
	t.Run("Zero values vs nil values", func(t *testing.T) {
		merger := NewConfigMerger()

		// Base config with non-zero values
		config := getDefaultConfig()
		config.Overlay.Placement.X = 50
		config.Overlay.Placement.Y = 50

		// Settings with explicit zero values (should override)
		settings := &ConfigSettings{
			Overlay: &OverlayConfigSettings{
				Placement: &PlacementSettings{
					X: &[]int{0}[0], // Explicit zero - should override
					Y: nil,          // Nil - should NOT override
				},
			},
		}

		merger.applySettingsToConfig(config, settings)

		if config.Overlay.Placement.X != 0 {
			t.Errorf("Expected X=0 (explicit zero), got %d", config.Overlay.Placement.X)
		}
		if config.Overlay.Placement.Y != 50 {
			t.Errorf("Expected Y=50 (preserved), got %d", config.Overlay.Placement.Y)
		}
	})

	t.Run("Empty string vs nil string", func(t *testing.T) {
		merger := NewConfigMerger()

		// Base config with non-empty values
		config := getDefaultConfig()
		config.Title.LineBreaking.StartProhibited = "ABC"
		config.Title.LineBreaking.EndProhibited = "XYZ"

		// Settings with explicit empty string (should override)
		settings := &ConfigSettings{
			Title: &TextSettings{
				LineBreaking: &LineBreakingSettings{
					StartProhibited: &[]string{""}[0], // Explicit empty - should override
					EndProhibited:   nil,              // Nil - should NOT override
				},
			},
		}

		merger.applySettingsToConfig(config, settings)

		if config.Title.LineBreaking.StartProhibited != "" {
			t.Errorf("Expected StartProhibited='' (explicit empty), got %s", config.Title.LineBreaking.StartProhibited)
		}
		if config.Title.LineBreaking.EndProhibited != "XYZ" {
			t.Errorf("Expected EndProhibited='XYZ' (preserved), got %s", config.Title.LineBreaking.EndProhibited)
		}
	})

	t.Run("Nested nil structures", func(t *testing.T) {
		merger := NewConfigMerger()

		// Base config
		config := getDefaultConfig()

		// Settings with nil nested structures
		settings := &ConfigSettings{
			Title: &TextSettings{
				Area:         nil, // Nil nested structure
				LineBreaking: nil, // Nil nested structure
			},
		}

		// Should not crash or affect existing values
		merger.applySettingsToConfig(config, settings)

		// Values should remain as defaults
		if config.Title.Area.X != DefaultTitleAreaX {
			t.Errorf("Expected area X to remain default %d, got %d", DefaultTitleAreaX, config.Title.Area.X)
		}
		if config.Title.LineBreaking.StartProhibited != DefaultStartProhibitedChars {
			t.Errorf("Expected line breaking to remain default")
		}
	})
}

// createTestImage creates a minimal test image file
func createTestImage(t *testing.T, path string) {
	// Create a minimal PNG-like content (just for testing file resolution)
	// This is not a real PNG, but enough for path resolution tests
	content := []byte("FAKE_PNG_CONTENT_FOR_TESTING")
	err := os.WriteFile(path, content, 0644)
	if err != nil {
		t.Fatalf("Failed to create test image %s: %v", path, err)
	}
}
