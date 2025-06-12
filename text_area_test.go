package main

import (
	"testing"
)

func TestNewTextArea(t *testing.T) {
	tests := []struct {
		name     string
		x, y     int
		w, h     int
		expected TextArea
	}{
		{
			name: "positive values",
			x:    10, y: 20, w: 100, h: 200,
			expected: TextArea{X: 10, Y: 20, Width: 100, Height: 200},
		},
		{
			name: "zero values",
			x:    0, y: 0, w: 0, h: 0,
			expected: TextArea{X: 0, Y: 0, Width: 0, Height: 0},
		},
		{
			name: "mixed values",
			x:    50, y: 75, w: 500, h: 300,
			expected: TextArea{X: 50, Y: 75, Width: 500, Height: 300},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := NewTextArea(tt.x, tt.y, tt.w, tt.h)

			if result.X != tt.expected.X {
				t.Errorf("Expected X=%d, got X=%d", tt.expected.X, result.X)
			}
			if result.Y != tt.expected.Y {
				t.Errorf("Expected Y=%d, got Y=%d", tt.expected.Y, result.Y)
			}
			if result.Width != tt.expected.Width {
				t.Errorf("Expected Width=%d, got Width=%d", tt.expected.Width, result.Width)
			}
			if result.Height != tt.expected.Height {
				t.Errorf("Expected Height=%d, got Height=%d", tt.expected.Height, result.Height)
			}
		})
	}
}

func TestTextArea_Contains(t *testing.T) {
	area := TextArea{X: 0, Y: 0, Width: 200, Height: 100}

	tests := []struct {
		name     string
		width    int
		height   int
		expected bool
	}{
		{
			name:  "fits exactly",
			width: 200, height: 100,
			expected: true,
		},
		{
			name:  "fits with room",
			width: 150, height: 80,
			expected: true,
		},
		{
			name:  "width too large",
			width: 250, height: 80,
			expected: false,
		},
		{
			name:  "height too large",
			width: 150, height: 120,
			expected: false,
		},
		{
			name:  "both too large",
			width: 250, height: 120,
			expected: false,
		},
		{
			name:  "zero dimensions",
			width: 0, height: 0,
			expected: true,
		},
		{
			name:  "edge case - one dimension zero",
			width: 200, height: 0,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := area.Contains(tt.width, tt.height)
			if result != tt.expected {
				t.Errorf("Contains(%d, %d) = %v, expected %v", tt.width, tt.height, result, tt.expected)
			}
		})
	}
}

func TestTextArea_SetDefaults(t *testing.T) {
	imageWidth := 1200
	imageHeight := 630

	tests := []struct {
		name     string
		area     TextArea
		expected TextArea
	}{
		{
			name:     "all zero - should use defaults",
			area:     TextArea{X: 0, Y: 0, Width: 0, Height: 0},
			expected: TextArea{X: 50, Y: 50, Width: 1100, Height: 530},
		},
		{
			name:     "partially zero - should keep original",
			area:     TextArea{X: 0, Y: 0, Width: 500, Height: 0},
			expected: TextArea{X: 0, Y: 0, Width: 500, Height: 0},
		},
		{
			name:     "non-zero values - should keep original",
			area:     TextArea{X: 100, Y: 100, Width: 400, Height: 300},
			expected: TextArea{X: 100, Y: 100, Width: 400, Height: 300},
		},
		{
			name:     "only X non-zero - should keep original",
			area:     TextArea{X: 25, Y: 0, Width: 0, Height: 0},
			expected: TextArea{X: 25, Y: 0, Width: 0, Height: 0},
		},
		{
			name:     "only Y non-zero - should keep original",
			area:     TextArea{X: 0, Y: 25, Width: 0, Height: 0},
			expected: TextArea{X: 0, Y: 25, Width: 0, Height: 0},
		},
		{
			name:     "only Width non-zero - should keep original",
			area:     TextArea{X: 0, Y: 0, Width: 500, Height: 0},
			expected: TextArea{X: 0, Y: 0, Width: 500, Height: 0},
		},
		{
			name:     "only Height non-zero - should keep original",
			area:     TextArea{X: 0, Y: 0, Width: 0, Height: 400},
			expected: TextArea{X: 0, Y: 0, Width: 0, Height: 400},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.area.SetDefaults(imageWidth, imageHeight)

			if result.X != tt.expected.X {
				t.Errorf("Expected X=%d, got X=%d", tt.expected.X, result.X)
			}
			if result.Y != tt.expected.Y {
				t.Errorf("Expected Y=%d, got Y=%d", tt.expected.Y, result.Y)
			}
			if result.Width != tt.expected.Width {
				t.Errorf("Expected Width=%d, got Width=%d", tt.expected.Width, result.Width)
			}
			if result.Height != tt.expected.Height {
				t.Errorf("Expected Height=%d, got Height=%d", tt.expected.Height, result.Height)
			}
		})
	}
}

func TestTextArea_SetDefaults_DifferentImageSizes(t *testing.T) {
	tests := []struct {
		name        string
		imageWidth  int
		imageHeight int
		expected    TextArea
	}{
		{
			name:       "standard size",
			imageWidth: 1200, imageHeight: 630,
			expected: TextArea{X: 50, Y: 50, Width: 1100, Height: 530},
		},
		{
			name:       "small image",
			imageWidth: 400, imageHeight: 300,
			expected: TextArea{X: 50, Y: 50, Width: 300, Height: 200},
		},
		{
			name:       "very small image",
			imageWidth: 100, imageHeight: 100,
			expected: TextArea{X: 50, Y: 50, Width: 0, Height: 0},
		},
		{
			name:       "large image",
			imageWidth: 2000, imageHeight: 1000,
			expected: TextArea{X: 50, Y: 50, Width: 1900, Height: 900},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			area := TextArea{X: 0, Y: 0, Width: 0, Height: 0}
			result := area.SetDefaults(tt.imageWidth, tt.imageHeight)

			if result.X != tt.expected.X {
				t.Errorf("Expected X=%d, got X=%d", tt.expected.X, result.X)
			}
			if result.Y != tt.expected.Y {
				t.Errorf("Expected Y=%d, got Y=%d", tt.expected.Y, result.Y)
			}
			if result.Width != tt.expected.Width {
				t.Errorf("Expected Width=%d, got Width=%d", tt.expected.Width, result.Width)
			}
			if result.Height != tt.expected.Height {
				t.Errorf("Expected Height=%d, got Height=%d", tt.expected.Height, result.Height)
			}
		})
	}
}
