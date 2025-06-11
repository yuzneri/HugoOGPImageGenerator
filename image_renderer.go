package main

import (
	"image"
	"image/color"
	"strings"

	"github.com/golang/freetype"
	"github.com/golang/freetype/truetype"
	"golang.org/x/image/font"
)

// ImageRenderer handles text rendering on images with Japanese line breaking support.
// It manages font sizing, text positioning, and layout within defined areas.
type ImageRenderer struct {
	textProcessor *TextProcessor
}

// NewImageRenderer creates a new ImageRenderer with the given text processor.
func NewImageRenderer(textProcessor *TextProcessor) *ImageRenderer {
	return &ImageRenderer{
		textProcessor: textProcessor,
	}
}

// RenderOptions contains all parameters needed for text rendering.
type RenderOptions struct {
	Font     *truetype.Font
	Config   *Config
	Title    string
	TestMode bool
}

// RenderTextOnImage renders text onto the provided image using the specified options.
// It handles automatic font sizing, text positioning, and applies Japanese line breaking rules.
func (ir *ImageRenderer) RenderTextOnImage(dst *image.RGBA, options *RenderOptions) error {
	config := options.Config
	area := config.Text.Area
	alignment := config.Text.BlockPosition
	lineAlignment := config.Text.LineAlignment
	overflow := config.Text.Overflow
	fontSize := config.Text.Size

	if area.X == 0 && area.Y == 0 && area.Width == 0 && area.Height == 0 {
		bounds := dst.Bounds()
		area.X = 50
		area.Y = 50
		area.Width = bounds.Dx() - 100
		area.Height = bounds.Dy() - 100
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

	c := freetype.NewContext()
	c.SetFont(options.Font)
	c.SetFontSize(fontSize)
	c.SetClip(dst.Bounds())
	c.SetDst(dst)

	textColor, err := parseHexColor(config.Text.Color)
	if err != nil {
		textColor = color.RGBA{R: 255, G: 255, B: 255, A: 255}
		DefaultLogger.Warning("Failed to parse color '%s', using white: %v", config.Text.Color, err)
	}
	c.SetSrc(image.NewUniform(textColor))

	face := truetype.NewFace(options.Font, &truetype.Options{Size: fontSize})
	lines := ir.textProcessor.SplitText(options.Title, face, maxWidth)

	if overflow == "shrink" {
		fontSize, face, lines = ir.adjustFontSizeToFit(options.Font, options.Title, config, area, maxWidth, lines)
		c.SetFontSize(fontSize)
	}

	ir.renderTextLines(c, face, lines, config, area, alignment, lineAlignment, options.TestMode, dst)

	return nil
}

// adjustFontSizeToFit automatically reduces font size until text fits within the specified area.
// It iteratively shrinks the font while respecting the minimum size constraint.
func (ir *ImageRenderer) adjustFontSizeToFit(font *truetype.Font, title string, config *Config, area TextArea, maxWidth int, lines []string) (float64, font.Face, []string) {
	fontSize := config.Text.Size
	maxHeight := area.Height
	minFontSize := config.Text.MinSize
	if minFontSize <= 0 {
		minFontSize = 12.0
	}

	for fontSize > minFontSize {
		lineHeight := int(fontSize * config.Text.LineHeight)
		totalHeight := len(lines) * lineHeight

		var maxTextWidth int
		face := truetype.NewFace(font, &truetype.Options{Size: fontSize})
		for _, line := range lines {
			textWidthPx := measureStringWithSpacing(face, line, config.Text.LetterSpacing)
			if textWidthPx > maxTextWidth {
				maxTextWidth = textWidthPx
			}
		}

		if totalHeight <= maxHeight && maxTextWidth <= area.Width {
			break
		}

		fontSize = fontSize * 0.9
		face = truetype.NewFace(font, &truetype.Options{Size: fontSize})
		lines = ir.textProcessor.SplitText(title, face, maxWidth)
	}

	face := truetype.NewFace(font, &truetype.Options{Size: fontSize})
	return fontSize, face, lines
}

// renderTextLines draws multiple lines of text with proper positioning and alignment.
// It handles both block-level alignment (within the text area) and line-level alignment.
func (ir *ImageRenderer) renderTextLines(c *freetype.Context, face font.Face, lines []string, config *Config, area TextArea, alignment, lineAlignment string, testMode bool, dst *image.RGBA) {
	fontSize := config.Text.Size
	lineHeight := int(fontSize * config.Text.LineHeight)
	totalHeight := len(lines) * lineHeight

	var maxTextWidth int
	for _, line := range lines {
		textWidthPx := measureStringWithSpacing(face, line, config.Text.LetterSpacing)
		if textWidthPx > maxTextWidth {
			maxTextWidth = textWidthPx
		}
	}

	blockX, blockY := calculateTextPosition(area, alignment, maxTextWidth, totalHeight)

	for i, line := range lines {
		textWidthPx := measureStringWithSpacing(face, line, config.Text.LetterSpacing)

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

		err := drawStringWithSpacing(c, face, line, lineX, y, config.Text.LetterSpacing)
		if err != nil {
			DefaultLogger.Warning("Failed to draw text line %d: %v", i, err)
		}
	}

	if testMode {
		ir.drawTestBorder(dst, area)
	}
}

// drawTestBorder draws a red border around the text area for debugging purposes.
func (ir *ImageRenderer) drawTestBorder(dst *image.RGBA, area TextArea) {
	borderColor := color.RGBA{R: 255, G: 0, B: 0, A: 255}

	for x := area.X; x < area.X+area.Width; x++ {
		dst.Set(x, area.Y, borderColor)
		dst.Set(x, area.Y+1, borderColor)
		dst.Set(x, area.Y+area.Height-1, borderColor)
		dst.Set(x, area.Y+area.Height-2, borderColor)
	}

	for y := area.Y; y < area.Y+area.Height; y++ {
		dst.Set(area.X, y, borderColor)
		dst.Set(area.X+1, y, borderColor)
		dst.Set(area.X+area.Width-1, y, borderColor)
		dst.Set(area.X+area.Width-2, y, borderColor)
	}
}
