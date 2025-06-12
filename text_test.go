package main

import (
	"testing"

	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font/gofont/goregular"
)

func TestTextProcessor_SplitText_JapaneseLineBreaking(t *testing.T) {
	// Create prohibited character maps for Japanese
	startProhibited := make(map[rune]bool)
	endProhibited := make(map[rune]bool)

	// Add some common Japanese prohibited characters
	for _, r := range []rune("。、！？）】」』") {
		startProhibited[r] = true
	}
	for _, r := range []rune("（【「『") {
		endProhibited[r] = true
	}

	tp := NewTextProcessor(startProhibited, endProhibited, 0)

	// Load a font for testing
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}
	face := truetype.NewFace(font, &truetype.Options{Size: 20})

	tests := []struct {
		name        string
		text        string
		maxWidth    int
		expectLines int
		description string
	}{
		{
			name:        "start prohibited character",
			text:        "これは日本語のテスト。改行されます",
			maxWidth:    200,
			expectLines: 2,
			description: "Should handle start prohibited characters (。)",
		},
		{
			name:        "end prohibited character",
			text:        "これは（日本語）のテストです",
			maxWidth:    200,
			expectLines: 2,
			description: "Should handle end prohibited characters (（)",
		},
		{
			name:        "manual line breaks",
			text:        "第一行\n第二行\n第三行",
			maxWidth:    1000,
			expectLines: 3,
			description: "Should respect manual line breaks",
		},
		{
			name:        "mixed content",
			text:        "English and 日本語 mixed content",
			maxWidth:    150,
			expectLines: 2,
			description: "Should handle mixed English and Japanese",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			lines := tp.SplitText(tt.text, face, tt.maxWidth)

			if len(lines) < tt.expectLines {
				t.Errorf("Expected at least %d lines, got %d. Lines: %v", tt.expectLines, len(lines), lines)
			}

			// Verify no empty lines
			for i, line := range lines {
				if line == "" {
					t.Errorf("Line %d should not be empty", i)
				}
			}
		})
	}
}

func TestTextProcessor_handleStartProhibited(t *testing.T) {
	startProhibited := map[rune]bool{
		'。': true,
		'、': true,
	}
	endProhibited := make(map[rune]bool)

	tp := NewTextProcessor(startProhibited, endProhibited, 0)

	// Test data
	lines := []string{"first line"}
	currentLine := []rune("これはテスト")
	runes := []rune("これはテスト。です")
	i := 7 // Position of '。'
	r := '。'

	newLines, newCurrentLine, newIndex := tp.handleStartProhibited(lines, currentLine, runes, i, r)

	// Should add the prohibited character to current line and finalize it
	if len(newLines) != 2 {
		t.Errorf("Expected 2 lines, got %d", len(newLines))
	}

	if len(newCurrentLine) != 0 {
		t.Errorf("Expected empty current line after handling start prohibited, got %v", newCurrentLine)
	}

	if newIndex <= i {
		t.Errorf("Expected index to advance, got %d", newIndex)
	}
}

func TestTextProcessor_handleEndProhibited(t *testing.T) {
	startProhibited := make(map[rune]bool)
	endProhibited := map[rune]bool{
		'（': true,
		'「': true,
	}

	tp := NewTextProcessor(startProhibited, endProhibited, 0)

	tests := []struct {
		name        string
		lines       []string
		currentLine []rune
		r           rune
		expectLines int
	}{
		{
			name:        "multiple characters in line",
			lines:       []string{"first line"},
			currentLine: []rune("これは（"),
			r:           'テ',
			expectLines: 2,
		},
		{
			name:        "single prohibited character",
			lines:       []string{"first line"},
			currentLine: []rune("（"),
			r:           'テ',
			expectLines: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newLines, newCurrentLine := tp.handleEndProhibited(tt.lines, tt.currentLine, tt.r)

			if len(newLines) != tt.expectLines {
				t.Errorf("Expected %d lines, got %d", tt.expectLines, len(newLines))
			}

			if len(newCurrentLine) == 0 {
				t.Error("Expected non-empty current line after handling end prohibited")
			}
		})
	}
}

func TestTextProcessor_handleWordBoundary(t *testing.T) {
	tp := NewTextProcessor(make(map[rune]bool), make(map[rune]bool), 0)

	tests := []struct {
		name        string
		lines       []string
		currentLine []rune
		r           rune
		expectLines int
	}{
		{
			name:        "word character boundary",
			lines:       []string{},
			currentLine: []rune("hello world test"),
			r:           'a',
			expectLines: 1,
		},
		{
			name:        "non-word character",
			lines:       []string{},
			currentLine: []rune("test、"),
			r:           'あ',
			expectLines: 1,
		},
		{
			name:        "empty current line",
			lines:       []string{},
			currentLine: []rune{},
			r:           'a',
			expectLines: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			newLines, newCurrentLine := tp.handleWordBoundary(tt.lines, tt.currentLine, tt.r)

			if len(newLines) != tt.expectLines {
				t.Errorf("Expected %d lines, got %d", tt.expectLines, len(newLines))
			}

			// Should always have a current line with the new character
			found := false
			for _, char := range newCurrentLine {
				if char == tt.r {
					found = true
					break
				}
			}
			if !found {
				t.Errorf("Expected character %c to be in current line %v", tt.r, newCurrentLine)
			}
		})
	}
}

func TestIsWordChar(t *testing.T) {
	tests := []struct {
		char     rune
		expected bool
	}{
		{'a', true},
		{'Z', true},
		{'0', true},
		{'9', true},
		{'_', true},
		{'-', true},
		{' ', false},
		{'あ', false},
		{'!', false},
		{'@', false},
	}

	for _, tt := range tests {
		t.Run(string(tt.char), func(t *testing.T) {
			result := isWordChar(tt.char)
			if result != tt.expected {
				t.Errorf("Expected isWordChar(%c) = %v, got %v", tt.char, tt.expected, result)
			}
		})
	}
}

func TestFindWordBoundary(t *testing.T) {
	tests := []struct {
		name     string
		runes    []rune
		maxPos   int
		expected int
	}{
		{
			name:     "find boundary in middle",
			runes:    []rune("hello world test"),
			maxPos:   10, // 't' in "test"
			expected: 5,  // space after "world"
		},
		{
			name:     "boundary at start",
			runes:    []rune("word"),
			maxPos:   2,
			expected: 2, // fallback to maxPos
		},
		{
			name:     "non-word character at maxPos",
			runes:    []rune("hello, world"),
			maxPos:   5, // comma
			expected: 5,
		},
		{
			name:     "maxPos out of bounds",
			runes:    []rune("test"),
			maxPos:   10,
			expected: 3, // adjusted to bounds
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := findWordBoundary(tt.runes, tt.maxPos)
			if result != tt.expected {
				t.Errorf("Expected %d, got %d", tt.expected, result)
			}
		})
	}
}

func TestCalculateTextPosition(t *testing.T) {
	area := TextArea{X: 100, Y: 50, Width: 400, Height: 200}
	textWidth := 200
	textHeight := 100

	tests := []struct {
		name      string
		alignment string
		expectedX int
		expectedY int
	}{
		{
			name:      "top-left",
			alignment: "top-left",
			expectedX: 100,
			expectedY: 50,
		},
		{
			name:      "top-center",
			alignment: "top-center",
			expectedX: 200, // 100 + (400-200)/2
			expectedY: 50,
		},
		{
			name:      "top-right",
			alignment: "top-right",
			expectedX: 300, // 100 + 400 - 200
			expectedY: 50,
		},
		{
			name:      "middle-center",
			alignment: "middle-center",
			expectedX: 200, // 100 + (400-200)/2
			expectedY: 100, // 50 + (200-100)/2
		},
		{
			name:      "bottom-right",
			alignment: "bottom-right",
			expectedX: 300, // 100 + 400 - 200
			expectedY: 150, // 50 + 200 - 100
		},
		{
			name:      "default (center)",
			alignment: "",
			expectedX: 200,
			expectedY: 100,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			x, y := calculateTextPosition(area, tt.alignment, textWidth, textHeight)

			if x != tt.expectedX {
				t.Errorf("Expected X=%d, got X=%d", tt.expectedX, x)
			}

			if y != tt.expectedY {
				t.Errorf("Expected Y=%d, got Y=%d", tt.expectedY, y)
			}
		})
	}
}

func TestTextProcessor_fitsInWidth(t *testing.T) {
	tp := NewTextProcessor(make(map[rune]bool), make(map[rune]bool), 0)

	// Load a font for testing
	font, err := truetype.Parse(goregular.TTF)
	if err != nil {
		t.Fatalf("Failed to parse font: %v", err)
	}
	face := truetype.NewFace(font, &truetype.Options{Size: 12})

	tests := []struct {
		name     string
		testLine []rune
		maxWidth int
		expected bool
	}{
		{
			name:     "short text fits",
			testLine: []rune("Hi"),
			maxWidth: 1000,
			expected: true,
		},
		{
			name:     "long text doesn't fit",
			testLine: []rune("This is a very long line that should not fit"),
			maxWidth: 50,
			expected: false,
		},
		{
			name:     "empty text fits",
			testLine: []rune(""),
			maxWidth: 100,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tp.fitsInWidth(tt.testLine, face, tt.maxWidth)
			if result != tt.expected {
				t.Errorf("Expected %v, got %v", tt.expected, result)
			}
		})
	}
}
