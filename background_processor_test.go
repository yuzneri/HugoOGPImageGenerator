package main

import (
	"errors"
	"image"
	"image/color"
	"image/png"
	"os"
	"path/filepath"
	"testing"
)

func TestNewBackgroundProcessor(t *testing.T) {
	configDir := "/test/config"

	processor := NewBackgroundProcessor(configDir)

	if processor == nil {
		t.Error("NewBackgroundProcessor should return a non-nil processor")
	}

	if processor.pathResolver == nil {
		t.Error("BackgroundProcessor should have a path resolver")
	}
}

func TestBackgroundProcessor_createColorBackground(t *testing.T) {
	processor := NewBackgroundProcessor("/test")

	tests := []struct {
		name          string
		colorHex      string
		expectedErr   bool
		expectedColor color.RGBA
	}{
		{
			name:          "Valid black color",
			colorHex:      "#000000",
			expectedErr:   false,
			expectedColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		},
		{
			name:          "Valid white color",
			colorHex:      "#FFFFFF",
			expectedErr:   false,
			expectedColor: color.RGBA{R: 255, G: 255, B: 255, A: 255},
		},
		{
			name:          "Valid red color",
			colorHex:      "#FF0000",
			expectedErr:   false,
			expectedColor: color.RGBA{R: 255, G: 0, B: 0, A: 255},
		},
		{
			name:          "Invalid color - defaults to black",
			colorHex:      "invalid",
			expectedErr:   false,
			expectedColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		},
		{
			name:          "Empty color - defaults to black",
			colorHex:      "",
			expectedErr:   false,
			expectedColor: color.RGBA{R: 0, G: 0, B: 0, A: 255},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			img, err := processor.createColorBackground(tt.colorHex)

			if (err != nil) != tt.expectedErr {
				t.Errorf("createColorBackground() error = %v, expectedErr = %v", err, tt.expectedErr)
				return
			}

			if img == nil {
				t.Error("createColorBackground() should return a non-nil image")
				return
			}

			// Check image dimensions
			bounds := img.Bounds()
			if bounds.Dx() != 1200 || bounds.Dy() != 630 {
				t.Errorf("Expected image size 1200x630, got %dx%d", bounds.Dx(), bounds.Dy())
			}

			// Check that the image is the expected color
			if rgbaImg, ok := img.(*image.RGBA); ok {
				// Sample a few pixels to verify color
				pixel1 := rgbaImg.RGBAAt(100, 100)
				pixel2 := rgbaImg.RGBAAt(600, 300)

				if pixel1 != tt.expectedColor {
					t.Errorf("Expected pixel color %+v, got %+v", tt.expectedColor, pixel1)
				}

				if pixel2 != tt.expectedColor {
					t.Errorf("Expected pixel color %+v, got %+v", tt.expectedColor, pixel2)
				}
			} else {
				t.Error("Expected *image.RGBA type")
			}
		})
	}
}

func TestBackgroundProcessor_CreateBackground_WithColor(t *testing.T) {
	processor := NewBackgroundProcessor("/test")

	config := &Config{}
	config.Background.Color = "#FF0000" // Red
	config.Background.Image = nil       // No image

	img, err := processor.CreateBackground(config, "/test/article")

	if err != nil {
		t.Errorf("CreateBackground() should not return error: %v", err)
	}

	if img == nil {
		t.Error("CreateBackground() should return a non-nil image")
	}

	// Check that it's red
	if rgbaImg, ok := img.(*image.RGBA); ok {
		pixel := rgbaImg.RGBAAt(100, 100)
		expected := color.RGBA{R: 255, G: 0, B: 0, A: 255}
		if pixel != expected {
			t.Errorf("Expected red pixel %+v, got %+v", expected, pixel)
		}
	}
}

func TestBackgroundProcessor_CreateBackground_WithImage(t *testing.T) {
	// Create a temporary directory and test image
	tempDir, err := os.MkdirTemp("", "bg_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test image file
	testImagePath := filepath.Join(tempDir, "test_bg.png")
	testImg := image.NewRGBA(image.Rect(0, 0, 100, 100))
	// Fill with blue
	for y := 0; y < 100; y++ {
		for x := 0; x < 100; x++ {
			testImg.Set(x, y, color.RGBA{R: 0, G: 0, B: 255, A: 255})
		}
	}

	file, err := os.Create(testImagePath)
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}
	err = png.Encode(file, testImg)
	file.Close()
	if err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}

	processor := NewBackgroundProcessor(tempDir)

	imagePath := "test_bg.png"
	config := &Config{}
	config.Background.Color = "#000000"
	config.Background.Image = &imagePath

	img, err := processor.CreateBackground(config, tempDir)

	if err != nil {
		t.Errorf("CreateBackground() should not return error: %v", err)
	}

	if img == nil {
		t.Error("CreateBackground() should return a non-nil image")
	}

	// Check image dimensions
	bounds := img.Bounds()
	if bounds.Dx() != 100 || bounds.Dy() != 100 {
		t.Errorf("Expected loaded image size 100x100, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestBackgroundProcessor_loadBackgroundImage_Success(t *testing.T) {
	// Create a temporary directory and test image
	tempDir, err := os.MkdirTemp("", "bg_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test image file
	testImagePath := filepath.Join(tempDir, "test_bg.png")
	testImg := image.NewRGBA(image.Rect(0, 0, 200, 150))
	// Fill with green
	for y := 0; y < 150; y++ {
		for x := 0; x < 200; x++ {
			testImg.Set(x, y, color.RGBA{R: 0, G: 255, B: 0, A: 255})
		}
	}

	file, err := os.Create(testImagePath)
	if err != nil {
		t.Fatalf("Failed to create test image file: %v", err)
	}
	err = png.Encode(file, testImg)
	file.Close()
	if err != nil {
		t.Fatalf("Failed to encode test image: %v", err)
	}

	processor := NewBackgroundProcessor(tempDir)

	img, err := processor.loadBackgroundImage("test_bg.png", tempDir)

	if err != nil {
		t.Errorf("loadBackgroundImage() should not return error: %v", err)
	}

	if img == nil {
		t.Error("loadBackgroundImage() should return a non-nil image")
	}

	// Check image dimensions
	bounds := img.Bounds()
	if bounds.Dx() != 200 || bounds.Dy() != 150 {
		t.Errorf("Expected loaded image size 200x150, got %dx%d", bounds.Dx(), bounds.Dy())
	}
}

func TestBackgroundProcessor_loadBackgroundImage_FileNotFound(t *testing.T) {
	processor := NewBackgroundProcessor("/nonexistent")

	img, err := processor.loadBackgroundImage("nonexistent.png", "/nonexistent")

	if err == nil {
		t.Error("loadBackgroundImage() should return error for nonexistent file")
	}

	if img != nil {
		t.Error("loadBackgroundImage() should return nil image on error")
	}

	// Check that the error is the correct type
	if !IsFileError(err) {
		t.Errorf("Expected FileError, got %T", err)
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		if appErr.Type != FileError {
			t.Errorf("Expected FileError type, got %v", appErr.Type)
		}
		if appErr.Context["operation"] != "open" {
			t.Errorf("Expected operation context 'open', got %v", appErr.Context["operation"])
		}
	} else {
		t.Error("Expected error to be AppError")
	}
}

func TestBackgroundProcessor_CreateBackground_EmptyImagePath(t *testing.T) {
	processor := NewBackgroundProcessor("/test")

	emptyImagePath := ""
	config := &Config{}
	config.Background.Color = "#00FF00"       // Green
	config.Background.Image = &emptyImagePath // Empty string, should use color instead

	img, err := processor.CreateBackground(config, "/test/article")

	if err != nil {
		t.Errorf("CreateBackground() should not return error for empty image path: %v", err)
	}

	if img == nil {
		t.Error("CreateBackground() should return a non-nil image")
	}

	// Should fall back to color background (green)
	if rgbaImg, ok := img.(*image.RGBA); ok {
		pixel := rgbaImg.RGBAAt(100, 100)
		expected := color.RGBA{R: 0, G: 255, B: 0, A: 255}
		if pixel != expected {
			t.Errorf("Expected green pixel %+v, got %+v", expected, pixel)
		}
	}
}

func TestBackgroundProcessor_loadBackgroundImage_InvalidImage(t *testing.T) {
	// Create a temporary directory and invalid image file
	tempDir, err := os.MkdirTemp("", "bg_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create an invalid image file (just text)
	invalidImagePath := filepath.Join(tempDir, "invalid.png")
	err = os.WriteFile(invalidImagePath, []byte("This is not an image file"), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid image file: %v", err)
	}

	processor := NewBackgroundProcessor(tempDir)

	img, err := processor.loadBackgroundImage("invalid.png", tempDir)

	if err == nil {
		t.Error("loadBackgroundImage() should return error for invalid image file")
	}

	if img != nil {
		t.Error("loadBackgroundImage() should return nil image for invalid file")
	}

	// Check that the error is an ImageError
	var appErr *AppError
	if errors.As(err, &appErr) {
		if appErr.Type != ImageError {
			t.Errorf("Expected ImageError type, got %v", appErr.Type)
		}
		if appErr.Context["operation"] != "decode" {
			t.Errorf("Expected operation context 'decode', got %v", appErr.Context["operation"])
		}
	} else {
		t.Error("Expected error to be AppError")
	}
}
