package main

import (
	"image"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// FontLoader interface for font management operations.
type FontLoader interface {
	LoadFont(fontPath string, articlePath string) (*truetype.Font, error)
	LoadFontWithFallback(fontPath string, articlePath string, defaultFont *truetype.Font) *truetype.Font
}

// BackgroundCreator interface for background processing operations.
type BackgroundCreator interface {
	CreateBackground(config *Config, articlePath string) (image.Image, error)
}

// ImageTextRenderer interface for text rendering operations.
type ImageTextRenderer interface {
	RenderTextOnImage(dst *image.RGBA, options *RenderOptions) error
}

// AssetPathResolver interface for path resolution operations.
type AssetPathResolver interface {
	ResolveAssetPath(assetPath string, articlePath string) string
	ResolveFromCwd(path string) (string, error)
}

// TextLineProcessor interface for text processing and line breaking.
type TextLineProcessor interface {
	SplitText(text string, face font.Face, maxWidth int) []string
}

// ConfigurationMerger interface for configuration merging operations.
type ConfigurationMerger interface {
	MergeConfigs(baseConfig *Config, ogpFM *OGPFrontMatter) *Config
}

// AppLogger interface for logging operations.
type AppLogger interface {
	Info(format string, args ...interface{})
	Warning(format string, args ...interface{})
	Error(format string, args ...interface{})
}

// ArticleImageProcessor interface for processing articles.
type ArticleImageProcessor interface {
	ProcessArticle(articlePath string, options ProcessOptions) error
}

// ImageGenerator interface for the main OGP generation orchestrator.
type ImageGenerator interface {
	GenerateSingle(articlePath string) error
	GenerateTest(articlePath string) error
	GenerateAll() error
}

// ImageOverlaySettings interface for overlay configuration.
type ImageOverlaySettings interface {
	GetImage() *string
	GetPlacement() *PlacementConfig
	GetFit() *string
	GetOpacity() *float64
}

// FileSystem interface for file system operations (useful for testing).
type FileSystem interface {
	ReadFile(filename string) ([]byte, error)
	WriteFile(filename string, data []byte, perm int) error
	MkdirAll(path string, perm int) error
	Stat(name string) (FileInfo, error)
}

// FileInfo interface for file information (useful for testing).
type FileInfo interface {
	Name() string
	Size() int64
	IsDir() bool
}
