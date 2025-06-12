package main

import (
	"image"
	"image/color"
	"testing"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func TestNewImageRenderer(t *testing.T) {
	renderer := NewImageRenderer()

	if renderer == nil {
		t.Error("NewImageRenderer should return a non-nil renderer")
	}
}

func TestImageRenderer_drawTestBorder(t *testing.T) {
	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))

	// Create renderer
	renderer := NewImageRenderer()

	// Define test area
	area := TextArea{
		X:      50,
		Y:      50,
		Width:  300,
		Height: 200,
	}

	// Draw border (test with title type)
	renderer.drawTestBorder(img, area, "title")

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

func TestImageRenderer_drawTestBorder_Description(t *testing.T) {
	// Create a test image
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))

	// Create renderer
	renderer := NewImageRenderer()

	// Define test area
	area := TextArea{
		X:      50,
		Y:      50,
		Width:  300,
		Height: 200,
	}

	// Draw border for description (should be blue)
	renderer.drawTestBorder(img, area, "description")

	// Check if border pixels are blue
	borderColor := color.RGBA{R: 0, G: 0, B: 255, A: 255}

	// Check top border
	if img.RGBAAt(area.X, area.Y) != borderColor {
		t.Error("Top border should be blue for description")
	}

	// Check bottom border
	if img.RGBAAt(area.X, area.Y+area.Height-1) != borderColor {
		t.Error("Bottom border should be blue for description")
	}

	// Check left border
	if img.RGBAAt(area.X, area.Y+10) != borderColor {
		t.Error("Left border should be blue for description")
	}

	// Check right border
	if img.RGBAAt(area.X+area.Width-1, area.Y+10) != borderColor {
		t.Error("Right border should be blue for description")
	}

	// Check interior is not blue (should be transparent/black)
	interior := img.RGBAAt(area.X+10, area.Y+10)
	if interior == borderColor {
		t.Error("Interior should not be blue")
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
	renderer := NewImageRenderer()

	// Create test config
	config := &Config{}
	config.Title.Visible = true
	config.Title.Size = 64.0
	config.Title.Color = "#000000"
	config.Title.MinSize = 12.0
	config.Title.LineHeight = 1.2
	config.Title.LetterSpacing = 1
	config.Title.Area = TextArea{
		X:      50,
		Y:      50,
		Width:  700,
		Height: 200,
	}
	config.Title.BlockPosition = "middle-center"
	config.Title.LineAlignment = "left"
	config.Title.Overflow = "shrink"
	config.Title.LineBreaking.StartProhibited = "。、"
	config.Title.LineBreaking.EndProhibited = "「（"

	config.Description.Visible = true
	config.Description.Size = 32.0
	config.Description.Color = "#666666"
	config.Description.MinSize = 12.0
	config.Description.LineHeight = 1.4
	config.Description.LetterSpacing = 0
	config.Description.Area = TextArea{
		X:      50,
		Y:      280,
		Width:  700,
		Height: 200,
	}
	config.Description.BlockPosition = "top-left"
	config.Description.LineAlignment = "left"
	config.Description.Overflow = "shrink"
	config.Description.LineBreaking.StartProhibited = "。、"
	config.Description.LineBreaking.EndProhibited = "「（"

	options := &RenderOptions{
		Font:        font,
		Config:      config,
		Title:       "テストタイトル",
		Description: "これはテスト用の説明文です。日本語の文章を含みます。",
		TestMode:    false,
	}

	// Render text
	err = renderer.RenderTextOnImage(img, options)
	if err != nil {
		t.Errorf("RenderTextOnImage should not return error: %v", err)
	}

	// Check that the image has some black pixels (text was rendered)
	// Since we have a white background and black text, look for black pixels
	hasBlackPixel := false
	hasGrayPixel := false
	bounds := img.Bounds()
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			pixel := img.RGBAAt(x, y)
			// Look for black pixels (title text)
			if pixel.R == 0 && pixel.G == 0 && pixel.B == 0 && pixel.A == 255 {
				hasBlackPixel = true
			}
			// Look for gray pixels (description text #666666)
			if pixel.R == 102 && pixel.G == 102 && pixel.B == 102 && pixel.A == 255 {
				hasGrayPixel = true
			}
			if hasBlackPixel && hasGrayPixel {
				break
			}
		}
		if hasBlackPixel && hasGrayPixel {
			break
		}
	}

	if !hasBlackPixel {
		t.Error("Expected some black pixels after rendering title text on white background")
	}

	if !hasGrayPixel {
		t.Error("Expected some gray pixels after rendering description text on white background")
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
	renderer := NewImageRenderer()

	// Create test config with minimal settings (testing defaults)
	config := &Config{}
	config.Title.Visible = true
	config.Title.Size = 64.0
	config.Title.Color = "#000000"
	config.Title.LetterSpacing = 1
	config.Title.LineHeight = 1.2
	// Area is zero, should use default
	config.Title.Area = TextArea{}
	config.Title.LineBreaking.StartProhibited = "。、"
	config.Title.LineBreaking.EndProhibited = "「（"

	config.Description.Visible = true
	config.Description.Size = 32.0
	config.Description.Color = "#666666"
	config.Description.Area = TextArea{}
	config.Description.LineBreaking.StartProhibited = "。、"
	config.Description.LineBreaking.EndProhibited = "「（"

	options := &RenderOptions{
		Font:        font,
		Config:      config,
		Title:       "Default Test",
		Description: "",
		TestMode:    true,
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

func TestImageRenderer_RenderTextOnImage_OnlyTitle(t *testing.T) {
	// Load font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}

	// Create test image
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))

	// Create renderer
	renderer := NewImageRenderer()

	// Create test config
	config := getDefaultConfig()

	options := &RenderOptions{
		Font:        font,
		Config:      config,
		Title:       "Only Title Test",
		Description: "", // Empty description
		TestMode:    false,
	}

	// Should not error with empty description
	err = renderer.RenderTextOnImage(img, options)
	if err != nil {
		t.Errorf("RenderTextOnImage should not return error with empty description: %v", err)
	}
}

func TestImageRenderer_RenderTextOnImage_OnlyDescription(t *testing.T) {
	// Load font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}

	// Create test image
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))

	// Create renderer
	renderer := NewImageRenderer()

	// Create test config
	config := getDefaultConfig()

	options := &RenderOptions{
		Font:        font,
		Config:      config,
		Title:       "", // Empty title
		Description: "Only Description Test",
		TestMode:    false,
	}

	// Should not error with empty title
	err = renderer.RenderTextOnImage(img, options)
	if err != nil {
		t.Errorf("RenderTextOnImage should not return error with empty title: %v", err)
	}
}

func TestImageRenderer_adjustFontSizeToFit(t *testing.T) {
	// Load font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}

	// Create renderer
	renderer := NewImageRenderer()

	// Create text config with small area to force font size reduction
	textConfig := &TextConfig{
		Size:          100.0, // Large font size
		MinSize:       12.0,
		LineHeight:    1.2,
		LetterSpacing: 1,
		LineBreaking: struct {
			StartProhibited string `yaml:"start_prohibited"`
			EndProhibited   string `yaml:"end_prohibited"`
		}{
			StartProhibited: "。、",
			EndProhibited:   "「（",
		},
	}

	area := TextArea{
		X:      0,
		Y:      0,
		Width:  100, // Small area
		Height: 50,  // Small area
	}

	title := "Very Long Title That Should Be Shrunk"
	initialLines := []string{title}

	// Create text processor for testing
	startProhibited, endProhibited := buildProhibitedMaps(textConfig)
	textProcessor := NewTextProcessor(startProhibited, endProhibited, textConfig.LetterSpacing)

	// Test font size adjustment
	adjustedSize, adjustedFace, adjustedLines := renderer.adjustFontSizeToFit(
		font, title, textConfig, area, area.Width, initialLines, textProcessor)

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
	renderer := NewImageRenderer()

	// Test with zero/negative minimum size (should default to 12.0)
	textConfig := &TextConfig{
		Size:          50.0,
		MinSize:       0.0, // Should default to 12.0
		LineHeight:    1.2,
		LetterSpacing: 1,
		LineBreaking: struct {
			StartProhibited string `yaml:"start_prohibited"`
			EndProhibited   string `yaml:"end_prohibited"`
		}{
			StartProhibited: "。、",
			EndProhibited:   "「（",
		},
	}

	area := TextArea{
		X:      0,
		Y:      0,
		Width:  50,
		Height: 20,
	}

	title := "Test"
	initialLines := []string{title}

	// Create text processor for testing
	startProhibited, endProhibited := buildProhibitedMaps(textConfig)
	textProcessor := NewTextProcessor(startProhibited, endProhibited, textConfig.LetterSpacing)

	adjustedSize, _, _ := renderer.adjustFontSizeToFit(
		font, title, textConfig, area, area.Width, initialLines, textProcessor)

	// Should not go below default minimum of 12.0
	if adjustedSize < 12.0 {
		t.Errorf("Font size should not go below default minimum 12.0, got %f", adjustedSize)
	}
}

func TestImageRenderer_renderSingleText_ColorParsing(t *testing.T) {
	// Load font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}

	// Create test image
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))

	// Create renderer
	renderer := NewImageRenderer()

	// Test with invalid color (should use white)
	textConfig := &TextConfig{
		Size:  32.0,
		Color: "invalid-color",
		Area: TextArea{
			X:      10,
			Y:      10,
			Width:  180,
			Height: 180,
		},
		LineBreaking: struct {
			StartProhibited string `yaml:"start_prohibited"`
			EndProhibited   string `yaml:"end_prohibited"`
		}{
			StartProhibited: "",
			EndProhibited:   "",
		},
	}

	err = renderer.renderSingleText(img, font, textConfig, "Test", false, "title")
	if err != nil {
		t.Errorf("renderSingleText should not return error even with invalid color: %v", err)
	}
}

func TestImageRenderer_TestMode(t *testing.T) {
	// Load font
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}

	// Create test image
	img := image.NewRGBA(image.Rect(0, 0, 800, 600))

	// Create renderer
	renderer := NewImageRenderer()

	// Create test config
	config := getDefaultConfig()

	// Set specific areas for both texts and enable description for this test
	config.Title.Area = TextArea{X: 50, Y: 50, Width: 300, Height: 100}
	config.Description.Visible = true // Enable description for this test
	config.Description.Area = TextArea{X: 50, Y: 200, Width: 400, Height: 150}

	options := &RenderOptions{
		Font:        font,
		Config:      config,
		Title:       "Test Title",
		Description: "Test Description",
		TestMode:    true, // Enable test mode
	}

	err = renderer.RenderTextOnImage(img, options)
	if err != nil {
		t.Errorf("RenderTextOnImage should not return error: %v", err)
	}

	// Check if borders are drawn with correct colors
	titleBorderColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}       // Red for title
	descriptionBorderColor := color.RGBA{R: 0, G: 0, B: 255, A: 255} // Blue for description

	// Check title area border (should be red)
	if img.RGBAAt(config.Title.Area.X, config.Title.Area.Y) != titleBorderColor {
		t.Error("Expected red border for title area in test mode")
	}

	// Check description area border (should be blue)
	if img.RGBAAt(config.Description.Area.X, config.Description.Area.Y) != descriptionBorderColor {
		t.Error("Expected blue border for description area in test mode")
	}
}
