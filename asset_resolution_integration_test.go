package main

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"
)

// TestUnifiedAssetResolution tests that all asset types use the same resolution priority
func TestUnifiedAssetResolution(t *testing.T) {
	// Create temporary directory structure
	tempDir, err := os.MkdirTemp("", "asset_resolution_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configDir := tempDir
	articleDir := filepath.Join(tempDir, "content", "books", "mybook")

	err = os.MkdirAll(articleDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create article dir: %v", err)
	}

	// Create article-specific cover.jpg
	articleCoverPath := filepath.Join(articleDir, "cover.jpg")
	err = createTestJPEG(articleCoverPath, color.RGBA{255, 0, 0, 255}) // Red image
	if err != nil {
		t.Fatalf("Failed to create article cover: %v", err)
	}

	// Create config-level cover.jpg
	configCoverPath := filepath.Join(configDir, "cover.jpg")
	err = createTestJPEG(configCoverPath, color.RGBA{0, 255, 0, 255}) // Green image
	if err != nil {
		t.Fatalf("Failed to create config cover: %v", err)
	}

	t.Run("Background image resolution", func(t *testing.T) {
		bgProcessor := NewBackgroundProcessor(configDir)

		config := &Config{
			Background: BackgroundConfig{
				Color: "#FFFFFF",
				Image: &[]string{"cover.jpg"}[0],
			},
		}

		img, err := bgProcessor.CreateBackground(config, articleDir)
		if err != nil {
			t.Errorf("Background creation failed: %v", err)
		}
		if img == nil {
			t.Error("Background image should not be nil")
		}
		// Note: We can't easily test which specific file was loaded without modifying the implementation
	})

	t.Run("Font path resolution", func(t *testing.T) {
		// Create article-specific font
		articleFontPath := filepath.Join(articleDir, "test.ttf")
		err = os.WriteFile(articleFontPath, []byte("dummy font content"), 0644)
		if err != nil {
			t.Fatalf("Failed to create article font: %v", err)
		}

		resolver := NewPathResolver(configDir)
		resolvedPath := resolver.ResolveAssetPath("test.ttf", articleDir)

		if resolvedPath != articleFontPath {
			t.Errorf("Expected article font path %s, got %s", articleFontPath, resolvedPath)
		}
	})

	t.Run("Overlay image resolution priority", func(t *testing.T) {
		resolver := NewPathResolver(configDir)

		// Test that article directory is checked first
		resolvedPath := resolver.ResolveAssetPath("cover.jpg", articleDir)

		if resolvedPath != articleCoverPath {
			t.Errorf("Expected article cover path %s, got %s", articleCoverPath, resolvedPath)
		}

		// Test fallback to config directory when article file doesn't exist
		resolvedPath = resolver.ResolveAssetPath("nonexistent.jpg", articleDir)
		expectedConfigPath := filepath.Join(configDir, "nonexistent.jpg")

		if resolvedPath != expectedConfigPath {
			t.Errorf("Expected config fallback path %s, got %s", expectedConfigPath, resolvedPath)
		}
	})
}

// TestAssetResolutionWithTypeConfig tests asset resolution when using type-specific configuration
func TestAssetResolutionWithTypeConfig(t *testing.T) {
	// Create temporary directory structure
	tempDir, err := os.MkdirTemp("", "type_asset_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configDir := tempDir
	articleDir := filepath.Join(tempDir, "content", "books", "book1")

	err = os.MkdirAll(articleDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create article dir: %v", err)
	}

	// Create article-specific assets
	articleCoverPath := filepath.Join(articleDir, "cover.jpg")
	err = createTestJPEG(articleCoverPath, color.RGBA{0, 0, 255, 255}) // Blue image
	if err != nil {
		t.Fatalf("Failed to create article cover: %v", err)
	}

	// Create config-level assets
	configCoverPath := filepath.Join(configDir, "cover.jpg")
	err = createTestJPEG(configCoverPath, color.RGBA{255, 255, 0, 255}) // Yellow image
	if err != nil {
		t.Fatalf("Failed to create config cover: %v", err)
	}

	// Create books.yaml type configuration
	booksConfigPath := filepath.Join(configDir, "books.yaml")
	booksConfigContent := `overlay:
  visible: true
  image: "cover.jpg"
  placement:
    x: 50
    y: 50
    height: 580
  fit: "contain"
  opacity: 1.0`

	err = os.WriteFile(booksConfigPath, []byte(booksConfigContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create books config: %v", err)
	}

	t.Run("Type config overlay resolution prioritizes article directory", func(t *testing.T) {
		// Load type configuration
		typeSettings, err := loadTypeConfigSettings(configDir, "books")
		if err != nil {
			t.Fatalf("Failed to load type config: %v", err)
		}

		if typeSettings == nil || typeSettings.Overlay == nil {
			t.Fatal("Type config should have overlay settings")
		}

		// Verify that the resolveAssetPath function would resolve to article directory first
		resolver := NewPathResolver(configDir)
		resolvedPath := resolver.ResolveAssetPath("cover.jpg", articleDir)

		if resolvedPath != articleCoverPath {
			t.Errorf("Expected article-specific cover %s, got %s", articleCoverPath, resolvedPath)
		}
	})
}

// createTestJPEG creates a test JPEG file with the specified color
func createTestJPEG(filename string, c color.RGBA) error {
	// Create a simple 100x100 image
	img := image.NewRGBA(image.Rect(0, 0, 100, 100))

	// Fill the image with the specified color
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			img.Set(x, y, c)
		}
	}

	// Create the file
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	// Encode as JPEG
	return jpeg.Encode(file, img, nil)
}
