package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

// FontManager handles font loading with caching for improved performance.
type FontManager struct {
	cache        map[string]*truetype.Font
	pathResolver *PathResolver
}

// NewFontManager creates a new FontManager with an empty cache.
func NewFontManager(configDir string) *FontManager {
	return &FontManager{
		cache:        make(map[string]*truetype.Font),
		pathResolver: NewPathResolver(configDir),
	}
}

// LoadFont loads a font from the filesystem with caching.
// It resolves the font path relative to config or article directories.
// If fontPath is empty, it uses the embedded Go regular font as default.
func (fm *FontManager) LoadFont(fontPath string, articlePath string) (*truetype.Font, error) {
	// Handle empty font path by using embedded default font
	if strings.TrimSpace(fontPath) == "" {
		return fm.getDefaultFont()
	}

	resolvedPath := fm.resolveFontPath(fontPath, articlePath)

	if font, exists := fm.cache[resolvedPath]; exists {
		return font, nil
	}

	fontBytes, err := os.ReadFile(resolvedPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read font file %s: %w", resolvedPath, err)
	}

	font, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse font: %w", err)
	}

	fm.cache[resolvedPath] = font
	return font, nil
}

// resolveFontPath resolves font path using the configured path resolver.
func (fm *FontManager) resolveFontPath(fontPath string, articlePath string) string {
	return fm.pathResolver.ResolveAssetPath(fontPath, articlePath)
}

// getDefaultFont returns the embedded default font with caching.
func (fm *FontManager) getDefaultFont() (*truetype.Font, error) {
	const defaultFontKey = "__default_embedded_font__"

	if font, exists := fm.cache[defaultFontKey]; exists {
		return font, nil
	}

	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		return nil, fmt.Errorf("failed to parse embedded default font: %w", err)
	}

	fm.cache[defaultFontKey] = font
	return font, nil
}

// LoadFontWithFallback loads a font with fallback to a default font on error.
// It logs warnings when font loading fails but continues execution.
func (fm *FontManager) LoadFontWithFallback(fontPath string, articlePath string, defaultFont *truetype.Font) *truetype.Font {
	font, err := fm.LoadFont(fontPath, articlePath)
	if err != nil {
		DefaultLogger.Warning("Failed to load font %s: %v, using default font", fontPath, err)
		return defaultFont
	}
	return font
}
