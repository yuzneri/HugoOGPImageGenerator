package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func TestNewFontManager(t *testing.T) {
	configDir := "/test/config"

	fm := NewFontManager(configDir)

	if fm == nil {
		t.Error("NewFontManager should return a non-nil font manager")
	}

	if fm.cache == nil {
		t.Error("FontManager should have an initialized cache")
	}

	if fm.pathResolver == nil {
		t.Error("FontManager should have a path resolver")
	}
}

func TestFontManager_LoadFont_Success(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "font_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test font file
	fontPath := filepath.Join(tempDir, "test_font.ttf")
	err = os.WriteFile(fontPath, goregular.TTF, 0644)
	if err != nil {
		t.Fatalf("Failed to create test font file: %v", err)
	}

	fm := NewFontManager(tempDir)

	// Test loading font
	font, err := fm.LoadFont("test_font.ttf", tempDir)

	if err != nil {
		t.Errorf("LoadFont should not return error: %v", err)
	}

	if font == nil {
		t.Error("LoadFont should return a non-nil font")
	}

	// Verify font is cached
	if _, exists := fm.cache[fontPath]; !exists {
		t.Error("Font should be cached after loading")
	}
}

func TestFontManager_LoadFont_Cache(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "font_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test font file
	fontPath := filepath.Join(tempDir, "test_font.ttf")
	err = os.WriteFile(fontPath, goregular.TTF, 0644)
	if err != nil {
		t.Fatalf("Failed to create test font file: %v", err)
	}

	fm := NewFontManager(tempDir)

	// Load font first time
	font1, err := fm.LoadFont("test_font.ttf", tempDir)
	if err != nil {
		t.Fatalf("First LoadFont should not return error: %v", err)
	}

	// Load font second time (should use cache)
	font2, err := fm.LoadFont("test_font.ttf", tempDir)
	if err != nil {
		t.Errorf("Second LoadFont should not return error: %v", err)
	}

	// Should return the same font instance from cache
	if font1 != font2 {
		t.Error("LoadFont should return cached font instance on second call")
	}
}

func TestFontManager_LoadFont_FileNotFound(t *testing.T) {
	fm := NewFontManager("/nonexistent")

	font, err := fm.LoadFont("nonexistent.ttf", "/nonexistent")

	if err == nil {
		t.Error("LoadFont should return error for nonexistent file")
	}

	if font != nil {
		t.Error("LoadFont should return nil font on error")
	}
}

func TestFontManager_LoadFont_InvalidFont(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "font_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create an invalid font file (just text)
	fontPath := filepath.Join(tempDir, "invalid_font.ttf")
	err = os.WriteFile(fontPath, []byte("This is not a font file"), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid font file: %v", err)
	}

	fm := NewFontManager(tempDir)

	font, err := fm.LoadFont("invalid_font.ttf", tempDir)

	if err == nil {
		t.Error("LoadFont should return error for invalid font file")
	}

	if font != nil {
		t.Error("LoadFont should return nil font for invalid file")
	}
}

func TestFontManager_LoadFontWithFallback_Success(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "font_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test font file
	fontPath := filepath.Join(tempDir, "test_font.ttf")
	err = os.WriteFile(fontPath, goregular.TTF, 0644)
	if err != nil {
		t.Fatalf("Failed to create test font file: %v", err)
	}

	// Create default font
	defaultFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("Failed to parse default font: %v", err)
	}

	fm := NewFontManager(tempDir)

	// Test loading font with fallback
	font := fm.LoadFontWithFallback("test_font.ttf", tempDir, defaultFont)

	if font == nil {
		t.Error("LoadFontWithFallback should return a non-nil font")
	}

	// The returned font should be the loaded font, not the fallback
	if font == defaultFont {
		t.Error("LoadFontWithFallback should return loaded font, not fallback when successful")
	}
}

func TestFontManager_LoadFontWithFallback_UseFallback(t *testing.T) {
	// Create default font
	defaultFont, err := truetype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("Failed to parse default font: %v", err)
	}

	fm := NewFontManager("/nonexistent")

	// Test loading nonexistent font with fallback
	font := fm.LoadFontWithFallback("nonexistent.ttf", "/nonexistent", defaultFont)

	if font == nil {
		t.Error("LoadFontWithFallback should return a non-nil font")
	}

	// Should return the fallback font
	if font != defaultFont {
		t.Error("LoadFontWithFallback should return fallback font when loading fails")
	}
}

func TestFontManager_resolveFontPath(t *testing.T) {
	configDir := "/test/config"
	fm := NewFontManager(configDir)

	// This test verifies that resolveFontPath calls pathResolver.ResolveAssetPath
	// Since we can't easily mock the pathResolver, we'll just verify it doesn't panic
	result := fm.resolveFontPath("font.ttf", "/test/article")

	// Should return some path (the exact path depends on PathResolver implementation)
	if result == "" {
		t.Error("resolveFontPath should return a non-empty path")
	}
}

func TestFontManager_LoadFont_MultipleFonts(t *testing.T) {
	// Create a temporary directory for testing
	tempDir, err := os.MkdirTemp("", "font_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create two test font files
	font1Path := filepath.Join(tempDir, "font1.ttf")
	font2Path := filepath.Join(tempDir, "font2.ttf")

	err = os.WriteFile(font1Path, goregular.TTF, 0644)
	if err != nil {
		t.Fatalf("Failed to create first font file: %v", err)
	}

	err = os.WriteFile(font2Path, goregular.TTF, 0644)
	if err != nil {
		t.Fatalf("Failed to create second font file: %v", err)
	}

	fm := NewFontManager(tempDir)

	// Load both fonts
	font1, err := fm.LoadFont("font1.ttf", tempDir)
	if err != nil {
		t.Errorf("Failed to load first font: %v", err)
	}

	font2, err := fm.LoadFont("font2.ttf", tempDir)
	if err != nil {
		t.Errorf("Failed to load second font: %v", err)
	}

	// Both should be non-nil
	if font1 == nil || font2 == nil {
		t.Error("Both fonts should be loaded successfully")
	}

	// Both should be cached
	if len(fm.cache) != 2 {
		t.Errorf("Expected 2 fonts in cache, got %d", len(fm.cache))
	}

	// Verify cache contains both fonts
	if _, exists := fm.cache[font1Path]; !exists {
		t.Error("First font should be in cache")
	}

	if _, exists := fm.cache[font2Path]; !exists {
		t.Error("Second font should be in cache")
	}
}

func TestFontManager_LoadFont_EmptyPath_UsesDefault(t *testing.T) {
	fm := NewFontManager("/test/config")

	// Test loading font with empty path
	font, err := fm.LoadFont("", "/test/article")

	if err != nil {
		t.Errorf("LoadFont with empty path should not return error: %v", err)
	}

	if font == nil {
		t.Error("LoadFont with empty path should return a non-nil font")
	}

	// Test loading font with whitespace-only path
	font2, err := fm.LoadFont("   ", "/test/article")

	if err != nil {
		t.Errorf("LoadFont with whitespace path should not return error: %v", err)
	}

	if font2 == nil {
		t.Error("LoadFont with whitespace path should return a non-nil font")
	}

	// Both should return the same cached default font
	if font != font2 {
		t.Error("LoadFont should return the same cached default font for empty paths")
	}
}
