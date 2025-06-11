package main

import (
	"image"
	"image/color"
	"testing"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func TestNewImageRenderer(t *testing.T) {
	textProcessor := NewTextProcessor(
		map[rune]bool{'。': true, '、': true},
		map[rune]bool{'「': true, '（': true},
		1,
	)

	renderer := NewImageRenderer(textProcessor)

	if renderer == nil {
		t.Error("NewImageRenderer should return a non-nil renderer")
	}

	if renderer.textProcessor != textProcessor {
		t.Error("NewImageRenderer should set the text processor correctly")
	}
}

func TestImageRenderer_drawTestBorder(t *testing.T) {
	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))

	// Create renderer
	textProcessor := NewTextProcessor(
		map[rune]bool{'。': true, '、': true},
		map[rune]bool{'「': true, '（': true},
		1,
	)
	renderer := NewImageRenderer(textProcessor)

	// Define test area
	area := TextArea{
		X:      50,
		Y:      50,
		Width:  300,
		Height: 200,
	}

	// Draw border
	renderer.drawTestBorder(img, area)

	// Check if border pixels are red
	borderColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	// Check top border
	if img.RGBAAt(area.X, area.Y) != borderColor {
		t.Error("Top border should be red")
	}

	// Check bottom border
	if img.RGBAAt(area.X, area.Y+area.Height-1) != borderColor {
		t.Error("Bottom border should be red")
	}

	// Check left border
	if img.RGBAAt(area.X, area.Y+10) != borderColor {
		t.Error("Left border should be red")
	}

	// Check right border
	if img.RGBAAt(area.X+area.Width-1, area.Y+10) != borderColor {
		t.Error("Right border should be red")
	}

	// Check interior is not red (should be transparent/black)
	interior := img.RGBAAt(area.X+10, area.Y+10)
	if interior == borderColor {
		t.Error("Interior should not be red")
	}
}

func TestImageRenderer_RenderTextOnImage(t *testing.T) {
	// Load font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}

	// Create test image with white background
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))
	// Fill with white background
	white := color.RGBA{R: 255, G: 255, B: 255, A: 255}
	for y := 0; y < 600; y++ {
		for x := 0; x < 800; x++ {
			img.Set(x, y, white)
		}
	}

	// Create renderer
	textProcessor := NewTextProcessor(
		map[rune]bool{'。': true, '、': true},
		map[rune]bool{'「': true, '（': true},
		1,
	)
	renderer := NewImageRenderer(textProcessor)

	// Create test config
	config := &Config{}
	config.Text.Size = 64.0
	config.Text.Color = "#000000"
	config.Text.MinSize = 12.0
	config.Text.LineHeight = 1.2
	config.Text.LetterSpacing = 1
	config.Text.Area = TextArea{
		X:      50,
		Y:      50,
		Width:  700,
		Height: 400,
	}
	config.Text.BlockPosition = "middle-center"
	config.Text.LineAlignment = "left"
	config.Text.Overflow = "shrink"

	options := &RenderOptions{
		Font:     font,
		Config:   config,
		Title:    "テストタイトル",
		TestMode: false,
	}

	// Render text
	err = renderer.RenderTextOnImage(img, options)
	if err != nil {
		t.Errorf("RenderTextOnImage should not return error: %v", err)
	}

	// Check that the image has some black pixels (text was rendered)
	// Since we have a white background and black text, look for black pixels
	hasBlackPixel := false
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.RGBAAt(x, y)
			// Look for black pixels (text)
			if pixel.R == 0 && pixel.G == 0 && pixel.B == 0 && pixel.A == 255 {
				hasBlackPixel = true
				break
			}
		}
		if hasBlackPixel {
			break
		}
	}

	if !hasBlackPixel {
		t.Error("Expected some black pixels after rendering text on white background")
	}
}

func TestImageRenderer_RenderTextOnImage_WithDefaults(t *testing.T) {
	// Load font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}

	// Create test image
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))

	// Create renderer
	textProcessor := NewTextProcessor(
		map[rune]bool{'。': true, '、': true},
		map[rune]bool{'「': true, '（': true},
		1,
	)
	renderer := NewImageRenderer(textProcessor)

	// Create test config with minimal settings (testing defaults)
	config := &Config{}
	config.Text.Size = 64.0
	config.Text.Color = "#000000"
	config.Text.LetterSpacing = 1
	config.Text.LineHeight = 1.2
	// Area is zero, should use default
	config.Text.Area = TextArea{}

	options := &RenderOptions{
		Font:     font,
		Config:   config,
		Title:    "Default Test",
		TestMode: true,
	}

	// Render text
	err = renderer.RenderTextOnImage(img, options)
	if err != nil {
		t.Errorf("RenderTextOnImage with defaults should not return error: %v", err)
	}

	// In test mode, there should be a red border drawn
	borderColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	// Check if border was drawn (default area should be 50,50,700,500)
	if img.RGBAAt(50, 50) != borderColor {
		t.Error("Expected red border pixel at default area position")
	}
}

func TestImageRenderer_adjustFontSizeToFit(t *testing.T) {
	// Load font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}

	// Create renderer
	textProcessor := NewTextProcessor(
		map[rune]bool{'。': true, '、': true},
		map[rune]bool{'「': true, '（': true},
		1,
	)
	renderer := NewImageRenderer(textProcessor)

	// Create config with small area to force font size reduction
	config := &Config{}
	config.Text.Size = 100.0 // Large font size
	config.Text.MinSize = 12.0
	config.Text.LineHeight = 1.2
	config.Text.LetterSpacing = 1

	area := TextArea{
		X:      0,
		Y:      0,
		Width:  100, // Small area
		Height: 50,  // Small area
	}

	title := "Very Long Title That Should Be Shrunk"
	initialLines := []string{title}

	// Test font size adjustment
	adjustedSize, adjustedFace, adjustedLines := renderer.adjustFontSizeToFit(
		font, title, config, area, area.Width, initialLines)

	// Font size should be reduced
	if adjustedSize >= 100.0 {
		t.Errorf("Font size should be reduced from 100.0, got %f", adjustedSize)
	}

	// Should not go below minimum
	if adjustedSize < 12.0 {
		t.Errorf("Font size should not go below minimum 12.0, got %f", adjustedSize)
	}

	// Face should not be nil
	if adjustedFace == nil {
		t.Error("Adjusted face should not be nil")
	}

	// Lines should not be empty
	if len(adjustedLines) == 0 {
		t.Error("Adjusted lines should not be empty")
	}
}

func TestImageRenderer_adjustFontSizeToFit_WithMinSize(t *testing.T) {
	// Load font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}

	// Create renderer
	textProcessor := NewTextProcessor(
		map[rune]bool{'。': true, '、': true},
		map[rune]bool{'「': true, '（': true},
		1,
	)
	renderer := NewImageRenderer(textProcessor)

	// Test with zero/negative minimum size (should default to 12.0)
	config := &Config{}
	config.Text.Size = 50.0
	config.Text.MinSize = 0.0 // Should default to 12.0
	config.Text.LineHeight = 1.2
	config.Text.LetterSpacing = 1

	area := TextArea{
		X:      0,
		Y:      0,
		Width:  50,
		Height: 20,
	}

	title := "Test"
	initialLines := []string{title}

	adjustedSize, _, _ := renderer.adjustFontSizeToFit(
		font, title, config, area, area.Width, initialLines)

	// Should not go below default minimum of 12.0
	if adjustedSize < 12.0 {
		t.Errorf("Font size should not go below default minimum 12.0, got %f", adjustedSize)
	}
}
