package main

// Image dimensions constants
const (
	// DefaultImageWidth is the standard OGP image width
	DefaultImageWidth = 1200

	// DefaultImageHeight is the standard OGP image height
	DefaultImageHeight = 630
)

// File and directory constants
const (
	// DefaultFilePermission for creating directories
	DefaultFilePermission = 0755

	// DefaultConfigFilename is the default config file name
	DefaultConfigFilename = "config.yaml"

	// DefaultIndexFilename is the Hugo content index file
	DefaultIndexFilename = "index.md"

	// DefaultTestFilename for test mode output
	DefaultTestFilename = "test.png"
)

// Default text configuration constants
const (
	// DefaultTitleFontSize for title text
	DefaultTitleFontSize = 64.0

	// DefaultDescriptionFontSize for description text
	DefaultDescriptionFontSize = 32.0

	// DefaultMinFontSize minimum font size for shrinking
	DefaultMinFontSize = 12.0

	// DefaultTitleMinSize minimum title font size
	DefaultTitleMinSize = 24.0

	// DefaultDescriptionMinSize minimum description font size
	DefaultDescriptionMinSize = 16.0

	// DefaultLineHeight multiplier for line spacing
	DefaultLineHeight = 1.2

	// DefaultTitleLetterSpacing for title text
	DefaultTitleLetterSpacing = 1

	// DefaultDescriptionLetterSpacing for description text
	DefaultDescriptionLetterSpacing = 0

	// FontSizeShrinkFactor for iterative font size reduction
	FontSizeShrinkFactor = 0.9
)

// Default area dimensions
const (
	// DefaultTitleAreaX starting X position for title area
	DefaultTitleAreaX = 100

	// DefaultTitleAreaY starting Y position for title area
	DefaultTitleAreaY = 50

	// DefaultTitleAreaWidth width of title text area
	DefaultTitleAreaWidth = 1000

	// DefaultTitleAreaHeight height of title text area
	DefaultTitleAreaHeight = 250

	// DefaultDescriptionAreaX starting X position for description area
	DefaultDescriptionAreaX = 100

	// DefaultDescriptionAreaY starting Y position for description area
	DefaultDescriptionAreaY = 350

	// DefaultDescriptionAreaWidth width of description text area
	DefaultDescriptionAreaWidth = 1000

	// DefaultDescriptionAreaHeight height of description text area
	DefaultDescriptionAreaHeight = 200

	// DefaultAreaPadding general padding for text areas
	DefaultAreaPadding = 50
)

// Default color constants
const (
	// DefaultBackgroundColor white background
	DefaultBackgroundColor = "#FFFFFF"

	// DefaultTitleColor black title text
	DefaultTitleColor = "#000000"

	// DefaultDescriptionColor gray description text
	DefaultDescriptionColor = "#666666"

	// DefaultWhiteColor for fallback
	DefaultWhiteColor = "#FFFFFF"

	// DefaultBlackColor for fallback
	DefaultBlackColor = "#000000"
)

// Position and alignment constants
const (
	// DefaultTitleBlockPosition for title text positioning
	DefaultTitleBlockPosition = "middle-center"

	// DefaultTitleLineAlignment for title text alignment
	DefaultTitleLineAlignment = "center"

	// DefaultDescriptionBlockPosition for description text positioning
	DefaultDescriptionBlockPosition = "top-left"

	// DefaultDescriptionLineAlignment for description text alignment
	DefaultDescriptionLineAlignment = "left"
)

// Overflow handling constants
const (
	// OverflowShrink shrinks text to fit
	OverflowShrink = "shrink"

	// OverflowClip clips text that doesn't fit
	OverflowClip = "clip"
)

// Output format constants
const (
	// FormatPNG PNG image format
	FormatPNG = "png"

	// FormatJPG JPG image format
	FormatJPG = "jpg"
)

// Directory constants
const (
	// DefaultOutputDirectory for generated images
	DefaultOutputDirectory = "public"

	// ContentDirectory name in Hugo projects
	ContentDirectory = "content"

	// StaticDirectory name in Hugo projects
	StaticDirectory = "static"
)

// Default overlay configuration constants
const (
	// DefaultOverlayVisible whether overlay is shown by default
	DefaultOverlayVisible = false

	// DefaultOverlayFit default image fit method
	DefaultOverlayFit = "contain"

	// DefaultOverlayOpacity default image opacity
	DefaultOverlayOpacity = 1.0

	// DefaultOverlayX default X position for overlay
	DefaultOverlayX = 50

	// DefaultOverlayY default Y position for overlay
	DefaultOverlayY = 50
)

// Japanese line breaking character sets
const (
	// DefaultStartProhibitedChars characters that cannot start a line
	DefaultStartProhibitedChars = ".)}]>!?、。，．！？)）］｝〉》」』ー～ぁぃぅぇぉっゃゅょゎァィゥェォッャュョヮヵヶ々"

	// DefaultEndProhibitedChars characters that cannot end a line
	DefaultEndProhibitedChars = "({[<（［｛〈《「『"
)

// Font cache constants
const (
	// DefaultFontCacheKey for embedded default font
	DefaultFontCacheKey = "__default_embedded_font__"
)

// Test border dimensions
const (
	// TestBorderThickness for debugging borders
	TestBorderThickness = 2
)
