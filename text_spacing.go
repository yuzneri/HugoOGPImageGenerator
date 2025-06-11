package main

import (
	"github.com/golang/freetype"
	"golang.org/x/image/font"
	"golang.org/x/image/math/fixed"
)

// measureStringWithSpacing calculates text width including letter spacing.
// It measures each character individually and adds spacing between them.
func measureStringWithSpacing(face font.Face, text string, letterSpacingPx int) int {
	if text == "" {
		return 0
	}

	runes := []rune(text)
	if len(runes) == 0 {
		return 0
	}

	// Measure character by character
	totalWidth := 0
	for i, r := range runes {
		charWidth := font.MeasureString(face, string(r))
		totalWidth += int(charWidth >> 6)

		// Add letter spacing between characters (not after the last character)
		if i < len(runes)-1 {
			totalWidth += letterSpacingPx
		}
	}

	return totalWidth
}

// drawStringWithSpacing draws text with custom letter spacing.
// It renders each character individually with the specified spacing.
func drawStringWithSpacing(c *freetype.Context, face font.Face, text string, x, y int, letterSpacingPx int) error {
	if text == "" {
		return nil
	}

	runes := []rune(text)
	if len(runes) == 0 {
		return nil
	}

	currentX := x
	for _, r := range runes {
		pt := fixed.Point26_6{X: fixed.I(currentX), Y: fixed.I(y)}
		_, err := c.DrawString(string(r), pt)
		if err != nil {
			return err
		}

		// Move to next character position
		charWidth := font.MeasureString(face, string(r))
		currentX += int(charWidth>>6) + letterSpacingPx
	}

	return nil
}
