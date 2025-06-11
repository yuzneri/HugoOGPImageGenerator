package main

import (
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"os"
)

// BackgroundProcessor handles background image loading and color background creation.
type BackgroundProcessor struct {
	pathResolver *PathResolver
}

// NewBackgroundProcessor creates a new BackgroundProcessor.
func NewBackgroundProcessor(configDir string) *BackgroundProcessor {
	return &BackgroundProcessor{
		pathResolver: NewPathResolver(configDir),
	}
}

// CreateBackground creates a background image from either an image file or solid color.
func (bp *BackgroundProcessor) CreateBackground(config *Config, articlePath string) (image.Image, error) {
	if config.Background.Image != nil && *config.Background.Image != "" {
		return bp.loadBackgroundImage(*config.Background.Image, articlePath)
	}

	return bp.createColorBackground(config.Background.Color)
}

// loadBackgroundImage loads a background image from the filesystem.
func (bp *BackgroundProcessor) loadBackgroundImage(imagePath string, articlePath string) (image.Image, error) {
	bgPath := bp.pathResolver.ResolveAssetPath(imagePath, articlePath)

	bgFile, err := os.Open(bgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open background image %s: %w", bgPath, err)
	}
	defer bgFile.Close()

	backgroundImage, _, err := image.Decode(bgFile)
	if err != nil {
		return nil, fmt.Errorf("failed to decode background image: %w", err)
	}

	return backgroundImage, nil
}

// createColorBackground creates a solid color background image (1200x630 pixels).
func (bp *BackgroundProcessor) createColorBackground(colorHex string) (image.Image, error) {
	bgColor, err := parseHexColor(colorHex)
	if err != nil {
		bgColor = color.RGBA{R: 0, G: 0, B: 0, A: 255}
	}

	backgroundImage := image.NewRGBA(image.Rect(0, 0, 1200, 630))
	draw.Draw(backgroundImage, backgroundImage.Bounds(), &image.Uniform{bgColor}, image.Point{}, draw.Src)

	return backgroundImage, nil
}
