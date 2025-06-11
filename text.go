package main

import (
	"strings"

	"golang.org/x/image/font"
)

// TextProcessor handles Japanese text processing with line breaking rules (禁則処理).
// It maintains maps of characters that have special line breaking constraints.
type TextProcessor struct {
	startProhibited map[rune]bool // Characters that cannot start a line (行頭禁則文字)
	endProhibited   map[rune]bool // Characters that cannot end a line (行末禁則文字)
	letterSpacing   int           // Letter spacing in pixels
}

// NewTextProcessor creates a new TextProcessor with the given prohibited character maps.
func NewTextProcessor(startProhibited, endProhibited map[rune]bool, letterSpacing int) *TextProcessor {
	return &TextProcessor{
		startProhibited: startProhibited,
		endProhibited:   endProhibited,
		letterSpacing:   letterSpacing,
	}
}

// SplitText breaks text into multiple lines respecting Japanese line breaking rules and maximum width.
// It handles both manual line breaks (\n) and automatic wrapping based on text width measurement.
func (t *TextProcessor) SplitText(text string, face font.Face, maxWidth int) []string {
	if text == "" {
		return []string{}
	}

	// Handle manual line breaks first
	if strings.Contains(text, "\n") {
		lines := strings.Split(text, "\n")
		var result []string

		// Process each manually broken line for automatic wrapping
		for _, line := range lines {
			subLines := t.splitTextSingle(line, face, maxWidth)
			result = append(result, subLines...)
		}

		return result
	}

	// No manual line breaks, process as single line
	return t.splitTextSingle(text, face, maxWidth)
}

// splitTextSingle processes a single line of text, applying automatic line breaking
// with Japanese prohibitions and English word boundary awareness.
func (t *TextProcessor) splitTextSingle(text string, face font.Face, maxWidth int) []string {
	if text == "" {
		return []string{}
	}

	// Convert to runes for proper Unicode handling
	runes := []rune(text)
	if len(runes) == 0 {
		return []string{text}
	}

	var lines []string
	var currentLine []rune

	i := 0
	for i < len(runes) {
		r := runes[i]
		testLine := append(currentLine, r)
		testLineStr := string(testLine)

		// Measure text width using font metrics with letter spacing
		textWidthPx := measureStringWithSpacing(face, testLineStr, t.letterSpacing)

		if textWidthPx <= maxWidth {
			// Character fits within the line width
			currentLine = testLine
			i++
		} else {
			// Character would exceed line width, apply breaking rules
			if t.startProhibited[r] {
				// 行頭禁則文字: Force character to stay on current line
				currentLine = testLine
				i++

				// Include any consecutive prohibited characters
				for i < len(runes) && t.startProhibited[runes[i]] {
					currentLine = append(currentLine, runes[i])
					i++
				}

				// Finalize the current line
				if len(currentLine) > 0 {
					lines = append(lines, string(currentLine))
					currentLine = []rune{}
				}
			} else if len(currentLine) > 0 && t.endProhibited[currentLine[len(currentLine)-1]] {
				// 行末禁則文字: Move prohibited character to next line
				if len(currentLine) > 1 {
					// Move the prohibited character to the next line
					endChar := currentLine[len(currentLine)-1]
					lines = append(lines, string(currentLine[:len(currentLine)-1]))
					currentLine = []rune{endChar, r}
				} else {
					// Only the prohibited character on this line, break anyway
					lines = append(lines, string(currentLine))
					currentLine = []rune{r}
				}
				i++
			} else {
				// Handle English word boundaries
				if len(currentLine) > 0 && isWordChar(r) {
					// Find word boundary to avoid breaking words mid-way
					boundaryPos := findWordBoundary(currentLine, len(currentLine)-1)

					if boundaryPos >= 0 && boundaryPos < len(currentLine)-1 {
						// Break at word boundary
						lines = append(lines, string(currentLine[:boundaryPos+1]))
						currentLine = append([]rune{}, currentLine[boundaryPos+1:]...)
						currentLine = append(currentLine, r)
						i++
					} else {
						// No good word boundary found, break normally
						if len(currentLine) > 0 {
							lines = append(lines, string(currentLine))
						}
						currentLine = []rune{r}
						i++
					}
				} else {
					// Normal character break
					if len(currentLine) > 0 {
						lines = append(lines, string(currentLine))
					}
					currentLine = []rune{r}
					i++
				}
			}
		}
	}

	// Add any remaining characters as the last line
	if len(currentLine) > 0 {
		lines = append(lines, string(currentLine))
	}

	return lines
}

// isWordChar determines if a character is part of an English word.
// Used to avoid breaking words in the middle when possible.
func isWordChar(r rune) bool {
	return (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
		(r >= '0' && r <= '9') || r == '_' || r == '-'
}

// findWordBoundary searches backwards from maxPos to find a suitable word break position.
// Returns the index where a line break would be appropriate to avoid splitting English words.
func findWordBoundary(runes []rune, maxPos int) int {
	if maxPos >= len(runes) {
		maxPos = len(runes) - 1
	}

	// If the character at maxPos is not a word character, break there
	if maxPos < 0 || !isWordChar(runes[maxPos]) {
		return maxPos
	}

	// Search backwards for the beginning of the current word
	for i := maxPos; i >= 0; i-- {
		if i == 0 {
			// Reached the beginning, word is too long, break at maxPos anyway
			return maxPos
		}
		if !isWordChar(runes[i-1]) {
			// Found word boundary
			return i - 1
		}
	}

	return maxPos
}

// calculateTextPosition determines the starting position for text rendering
// within the specified area, taking into account the desired alignment.
func calculateTextPosition(area TextArea, alignment string, textWidth, textHeight int) (x, y int) {
	// Calculate horizontal alignment
	var alignX int
	switch {
	case strings.Contains(alignment, "left"):
		alignX = area.X
	case strings.Contains(alignment, "right"):
		alignX = area.X + area.Width - textWidth
	default: // center
		alignX = area.X + (area.Width-textWidth)/2
	}

	// Calculate vertical alignment
	var alignY int
	switch {
	case strings.Contains(alignment, "top"):
		alignY = area.Y
	case strings.Contains(alignment, "bottom"):
		alignY = area.Y + area.Height - textHeight
	default: // middle
		alignY = area.Y + (area.Height-textHeight)/2
	}

	return alignX, alignY
}
