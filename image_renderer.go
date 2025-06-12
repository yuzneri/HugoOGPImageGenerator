package main

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// ImageRenderer handles text rendering on images with Japanese line breaking support.
// It manages font sizing, text positioning, and layout within defined areas.
// It implements the ImageTextRenderer interface.
type ImageRenderer struct {
	// We'll create text processors on demand based on each text's configuration
}

// Verify that ImageRenderer implements ImageTextRenderer interface
var _ ImageTextRenderer = (*ImageRenderer)(nil)

// NewImageRenderer creates a new ImageRenderer.
func NewImageRenderer() *ImageRenderer {
	return &ImageRenderer{}
}

// RenderOptions contains all parameters needed for text rendering.
type RenderOptions struct {
	Font        *truetype.Font
	Config      *Config
	Title       string
	Description string
	TestMode    bool
}

// RenderTextOnImage renders text onto the provided image using the specified options.
// It handles automatic font sizing, text positioning, and applies Japanese line breaking rules.
func (ir *ImageRenderer) RenderTextOnImage(dst *image.RGBA, options *RenderOptions) error {
	// Render title if visible and provided
	if options.Config.Title.Visible && options.Title != "" {
		err := ir.renderSingleText(dst, options.Font, &options.Config.Title, options.Title, options.TestMode, "title")
		if err != nil {
			return fmt.Errorf("failed to render title: %w", err)
		}
	}

	// Render description if visible and provided
	if options.Config.Description.Visible && options.Description != "" {
		err := ir.renderSingleText(dst, options.Font, &options.Config.Description, options.Description, options.TestMode, "description")
		if err != nil {
			return fmt.Errorf("failed to render description: %w", err)
		}
	}

	return nil
}

// renderSingleText renders a single text element (title or description) onto the image.
func (ir *ImageRenderer) renderSingleText(dst *image.RGBA, font *truetype.Font, textConfig *TextConfig, text string, testMode bool, textType string) error {
	area := textConfig.Area
	alignment := textConfig.BlockPosition
	lineAlignment := textConfig.LineAlignment
	overflow := textConfig.Overflow
	fontSize := textConfig.Size

	if area.X == 0 && area.Y == 0 && area.Width == 0 && area.Height == 0 {
		bounds := dst.Bounds()
		area.X = DefaultAreaPadding
		area.Y = DefaultAreaPadding
		area.Width = bounds.Dx() - (DefaultAreaPadding * 2)
		area.Height = bounds.Dy() - (DefaultAreaPadding * 2)
	}

	if alignment == "" {
		alignment = "middle-center"
	}

	if lineAlignment == "" {
		if strings.Contains(alignment, "left") {
			lineAlignment = "left"
		} else if strings.Contains(alignment, "right") {
			lineAlignment = "right"
		} else {
			lineAlignment = "center"
		}
	}

	if overflow == "" {
		overflow = "shrink"
	}

	maxWidth := area.Width

	// Create text processor for this specific text configuration
	startProhibited, endProhibited := buildProhibitedMaps(textConfig)
	textProcessor := NewTextProcessor(startProhibited, endProhibited, textConfig.LetterSpacing)

	c := freetype.NewContext()
	c.SetFont(font)
	c.SetFontSize(fontSize)
	c.SetClip(dst.Bounds())
	c.SetDst(dst)

	textColor, err := parseHexColor(textConfig.Color)
	if err != nil {
		textColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		DefaultLogger.Warning("Failed to parse color '%s', using white: %v", textConfig.Color, err)
	}
	c.SetSrc(image.NewUniform(textColor))

	face := truetype.NewFace(font, &truetype.Options{Size: fontSize})
	lines := textProcessor.SplitText(text, face, maxWidth)

	if overflow == "shrink" {
		fontSize, face, lines = ir.adjustFontSizeToFit(font, text, textConfig, area, maxWidth, lines, textProcessor)
		c.SetFontSize(fontSize)
	}

	ir.renderTextLines(c, face, lines, textConfig, area, alignment, lineAlignment, testMode, dst, textType)

	return nil
}

// adjustFontSizeToFit automatically reduces font size until text fits within the specified area.
// It iteratively shrinks the font while respecting the minimum size constraint.
func (ir *ImageRenderer) adjustFontSizeToFit(font *truetype.Font, title string, textConfig *TextConfig, area TextArea, maxWidth int, lines []string, textProcessor *TextProcessor) (float64, font.Face, []string) {
	fontSize := textConfig.Size
	maxHeight := area.Height
	minFontSize := textConfig.MinSize
	if minFontSize <= 0 {
		minFontSize = DefaultMinFontSize
	}

	for fontSize > minFontSize {
		lineHeight := int(fontSize * textConfig.LineHeight)
		totalHeight := len(lines) * lineHeight

		var maxTextWidth int
		face := truetype.NewFace(font, &truetype.Options{Size: fontSize})
		for _, line := range lines {
			textWidthPx := measureStringWithSpacing(face, line, textConfig.LetterSpacing)
			if textWidthPx > maxTextWidth {
				maxTextWidth = textWidthPx
			}
		}

		if totalHeight <= maxHeight && maxTextWidth <= area.Width {
			break
		}

		fontSize = fontSize * FontSizeShrinkFactor
		face = truetype.NewFace(font, &truetype.Options{Size: fontSize})
		lines = textProcessor.SplitText(title, face, maxWidth)
	}

	face := truetype.NewFace(font, &truetype.Options{Size: fontSize})
	return fontSize, face, lines
}

// renderTextLines draws multiple lines of text with proper positioning and alignment.
// It handles both block-level alignment (within the text area) and line-level alignment.
func (ir *ImageRenderer) renderTextLines(c *freetype.Context, face font.Face, lines []string, textConfig *TextConfig, area TextArea, alignment, lineAlignment string, testMode bool, dst *image.RGBA, textType string) {
	fontSize := textConfig.Size
	lineHeight := int(fontSize * textConfig.LineHeight)
	totalHeight := len(lines) * lineHeight

	var maxTextWidth int
	for _, line := range lines {
		textWidthPx := measureStringWithSpacing(face, line, textConfig.LetterSpacing)
		if textWidthPx > maxTextWidth {
			maxTextWidth = textWidthPx
		}
	}

	blockX, blockY := calculateTextPosition(area, alignment, maxTextWidth, totalHeight)

	for i, line := range lines {
		textWidthPx := measureStringWithSpacing(face, line, textConfig.LetterSpacing)

		var lineX int
		switch lineAlignment {
		case "left":
			lineX = blockX
		case "right":
			lineX = blockX + maxTextWidth - textWidthPx
		default:
			lineX = blockX + (maxTextWidth-textWidthPx)/2
		}

		y := blockY + lineHeight + i*lineHeight

		err := drawStringWithSpacing(c, face, line, lineX, y, textConfig.LetterSpacing)
		if err != nil {
			DefaultLogger.Warning("Failed to draw text line %d: %v", i, err)
		}
	}

	if testMode {
		ir.drawTestBorder(dst, area, textType)
	}
}

// drawTestBorder draws a colored border around the text area for debugging purposes.
// Title areas use red borders, description areas use blue borders.
func (ir *ImageRenderer) drawTestBorder(dst *image.RGBA, area TextArea, textType string) {
	var borderColor color.RGBA

	// Determine border color based on text type
	if textType == "title" {
		borderColor = color.RGBA{R: 255, G: 0, B: 0, A: 255} // Red for title
	} else {
		borderColor = color.RGBA{R: 0, G: 0, B: 255, A: 255} // Blue for description
	}

	for x := area.X; x < area.X+area.Width; x++ {
		for i := 0; i < TestBorderThickness; i++ {
			dst.Set(x, area.Y+i, borderColor)
			dst.Set(x, area.Y+area.Height-1-i, borderColor)
		}
	}

	for y := area.Y; y < area.Y+area.Height; y++ {
		for i := 0; i < TestBorderThickness; i++ {
			dst.Set(area.X+i, y, borderColor)
			dst.Set(area.X+area.Width-1-i, y, borderColor)
		}
	}
}
