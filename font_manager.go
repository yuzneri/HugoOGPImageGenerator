package main

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	"github.com/golang/freetype/truetype"
)

// FontManager handles font loading with caching for improved performance.
// It implements the FontLoader interface.
type FontManager struct {
	cache        map[string]*truetype.Font
	pathResolver AssetPathResolver
}

// Verify that FontManager implements FontLoader interface
var _ FontLoader = (*FontManager)(nil)

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
		return nil, NewFileError("read", resolvedPath, err)
	}

	font, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, NewFontError("parse", resolvedPath, err)
	}

	fm.cache[resolvedPath] = font
	return font, nil
}

// resolveFontPath resolves font path using the configured path resolver.
func (fm *FontManager) resolveFontPath(fontPath string, articlePath string) string {
	return fm.pathResolver.ResolveAssetPath(fontPath, articlePath)
}

// getDefaultFont returns a system font with caching.
func (fm *FontManager) getDefaultFont() (*truetype.Font, error) {
	const defaultFontKey = DefaultFontCacheKey

	if font, exists := fm.cache[defaultFontKey]; exists {
		return font, nil
	}

	fontPath := fm.findSystemFont()
	if fontPath == "" {
		return nil, NewFontError("load", "system font", fmt.Errorf("no suitable system font found"))
	}

	fontBytes, err := os.ReadFile(fontPath)
	if err != nil {
		return nil, NewFileError("read", fontPath, err)
	}

	font, err := truetype.Parse(fontBytes)
	if err != nil {
		return nil, NewFontError("parse", fontPath, err)
	}

	fm.cache[defaultFontKey] = font
	return font, nil
}

// findSystemFont tries to find a suitable system font by checking common paths.
func (fm *FontManager) findSystemFont() string {
	// OS-specific font paths in priority order
	var fontPaths []string

	switch runtime.GOOS {
	case "darwin": // macOS
		fontPaths = []string{
			// Japanese fonts first
			"/System/Library/Fonts/ヒラギノ角ゴシック W3.ttc",
			"/System/Library/Fonts/Hiragino Sans GB W3.otf",
			"/Library/Fonts/Osaka.ttf",
			"/System/Library/Fonts/AppleGothic.ttf",
			// Fallback to standard fonts
			"/System/Library/Fonts/Helvetica.ttc",
			"/System/Library/Fonts/Arial.ttf",
		}
	case "linux":
		fontPaths = []string{
			// Japanese fonts first
			"/usr/share/fonts/truetype/fonts-japanese-gothic.ttf",
			"/usr/share/fonts/truetype/noto/NotoSansCJK-Regular.ttc",
			"/usr/share/fonts/opentype/noto/NotoSansCJK-Regular.otf",
			"/usr/share/fonts/truetype/takao-gothic/TakaoGothic.ttf",
			// Fallback to standard fonts
			"/usr/share/fonts/truetype/ubuntu/Ubuntu-R.ttf",
			"/usr/share/fonts/truetype/dejavu/DejaVuSans.ttf",
			"/usr/share/fonts/truetype/liberation/LiberationSans-Regular.ttf",
		}
	case "windows":
		fontPaths = []string{
			// Japanese fonts first
			"C:\\Windows\\Fonts\\meiryo.ttc",
			"C:\\Windows\\Fonts\\YuGothM.ttc",
			"C:\\Windows\\Fonts\\msgothic.ttc",
			"C:\\Windows\\Fonts\\msmincho.ttc",
			// Fallback to standard fonts
			"C:\\Windows\\Fonts\\arial.ttf",
			"C:\\Windows\\Fonts\\calibri.ttf",
			"C:\\Windows\\Fonts\\tahoma.ttf",
		}
	default:
		// Try common paths for unknown OS (Japanese fonts first)
		fontPaths = []string{
			"/usr/share/fonts/truetype/fonts-japanese-gothic.ttf",
			"C:\\Windows\\Fonts\\msgothic.ttc",
			"/System/Library/Fonts/ヒラギノ角ゴシック W3.ttc",
			"/usr/share/fonts/truetype/ubuntu/Ubuntu-R.ttf",
			"/System/Library/Fonts/Helvetica.ttc",
			"C:\\Windows\\Fonts\\arial.ttf",
		}
	}

	// Check each path until we find a valid font
	for _, fontPath := range fontPaths {
		if fm.isValidFont(fontPath) {
			return fontPath
		}
	}

	return ""
}

// isValidFont checks if a font file exists and is readable.
func (fm *FontManager) isValidFont(fontPath string) bool {
	if fontPath == "" {
		return false
	}

	info, err := os.Stat(fontPath)
	if err != nil {
		return false
	}

	// Basic file size check (fonts should be reasonable size)
	return info.Size() > 1024
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
