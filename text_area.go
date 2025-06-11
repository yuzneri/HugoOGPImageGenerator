package main

// TextArea defines a rectangular area for text rendering.
type TextArea struct {
	X      int `yaml:"x"`
	Y      int `yaml:"y"`
	Width  int `yaml:"width"`
	Height int `yaml:"height"`
}

// NewTextArea creates a new TextArea with the given dimensions.
func NewTextArea(x, y, width, height int) TextArea {
	return TextArea{
		X:      x,
		Y:      y,
		Width:  width,
		Height: height,
	}
}

// Contains checks if the given dimensions fit within this text area.
func (ta TextArea) Contains(width, height int) bool {
	return width <= ta.Width && height <= ta.Height
}

// SetDefaults returns a TextArea with default values if all dimensions are zero.
// The default area has 50px margins on all sides.
func (ta TextArea) SetDefaults(imageWidth, imageHeight int) TextArea {
	if ta.X == 0 && ta.Y == 0 && ta.Width == 0 && ta.Height == 0 {
		return TextArea{
			X:      50,
			Y:      50,
			Width:  imageWidth - 100,
			Height: imageHeight - 100,
		}
	}
	return ta
}
