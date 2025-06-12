package main

import (
	"image/color"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseHexColor(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected color.RGBA
		hasError bool
	}{
		{
			name:     "6-character hex without #",
			input:    "FF0000",
			expected: color.RGBA{R: 255, G: 0, B: 0, A: 255},
			hasError: false,
		},
		{
			name:     "6-character hex with #",
			input:    "#00FF00",
			expected: color.RGBA{R: 0, G: 255, B: 0, A: 255},
			hasError: false,
		},
		{
			name:     "8-character hex with alpha",
			input:    "#0000FF80",
			expected: color.RGBA{R: 0, G: 0, B: 255, A: 128},
			hasError: false,
		},
		{
			name:     "8-character hex without #",
			input:    "FF00FF40",
			expected: color.RGBA{R: 255, G: 0, B: 255, A: 64},
			hasError: false,
		},
		{
			name:     "lowercase hex",
			input:    "#abcdef",
			expected: color.RGBA{R: 171, G: 205, B: 239, A: 255},
			hasError: false,
		},
		{
			name:     "mixed case hex",
			input:    "#AbCdEf",
			expected: color.RGBA{R: 171, G: 205, B: 239, A: 255},
			hasError: false,
		},
		{
			name:     "white color",
			input:    "#FFFFFF",
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			hasError: false,
		},
		{
			name:     "black color",
			input:    "#000000",
			expected: color.RGBA{R: 0, G: 0, B: 0, A: 255},
			hasError: false,
		},
		{
			name:     "transparent black",
			input:    "#00000000",
			expected: color.RGBA{R: 0, G: 0, B: 0, A: 0},
			hasError: false,
		},
		{
			name:     "fully opaque white",
			input:    "#FFFFFFFF",
			expected: color.RGBA{R: 255, G: 255, B: 255, A: 255},
			hasError: false,
		},
		{
			name:     "invalid hex character",
			input:    "#GGGGGG",
			expected: color.RGBA{},
			hasError: true,
		},
		{
			name:     "too short",
			input:    "#FFF",
			expected: color.RGBA{},
			hasError: true,
		},
		{
			name:     "too long",
			input:    "#FFFFFFFFF",
			expected: color.RGBA{},
			hasError: true,
		},
		{
			name:     "empty string",
			input:    "",
			expected: color.RGBA{},
			hasError: true,
		},
		{
			name:     "only hash",
			input:    "#",
			expected: color.RGBA{},
			hasError: true,
		},
		{
			name:     "7 characters",
			input:    "#FFFFFFF",
			expected: color.RGBA{},
			hasError: true,
		},
		{
			name:     "non-hex characters",
			input:    "#ZZZZZZ",
			expected: color.RGBA{},
			hasError: true,
		},
		{
			name:     "special characters",
			input:    "#@!$%^&",
			expected: color.RGBA{},
			hasError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseHexColor(tt.input)

			if tt.hasError {
				if err == nil {
					t.Errorf("Expected error for input %q, but got none", tt.input)
				}
				return
			}

			if err != nil {
				t.Errorf("Expected no error for input %q, but got %v", tt.input, err)
				return
			}

			if result.R != tt.expected.R {
				t.Errorf("Expected R=%d, got R=%d", tt.expected.R, result.R)
			}
			if result.G != tt.expected.G {
				t.Errorf("Expected G=%d, got G=%d", tt.expected.G, result.G)
			}
			if result.B != tt.expected.B {
				t.Errorf("Expected B=%d, got B=%d", tt.expected.B, result.B)
			}
			if result.A != tt.expected.A {
				t.Errorf("Expected A=%d, got A=%d", tt.expected.A, result.A)
			}
		})
	}
}

func TestGetDefaultConfig(t *testing.T) {
	config := getDefaultConfig()

	if config == nil {
		t.Fatal("getDefaultConfig should return a non-nil config")
	}

	// Test background defaults
	if config.Background.Color != DefaultBackgroundColor {
		t.Errorf("Expected background color %q, got %q", DefaultBackgroundColor, config.Background.Color)
	}

	// Test output defaults
	if config.Output.Directory != DefaultOutputDirectory {
		t.Errorf("Expected output directory %q, got %q", DefaultOutputDirectory, config.Output.Directory)
	}

	if config.Output.Format != FormatPNG {
		t.Errorf("Expected output format %q, got %q", FormatPNG, config.Output.Format)
	}

	// Test title defaults
	if !config.Title.Visible {
		t.Error("Expected title to be visible by default")
	}

	if config.Title.Size != DefaultTitleFontSize {
		t.Errorf("Expected title size %f, got %f", DefaultTitleFontSize, config.Title.Size)
	}

	if config.Title.Color != DefaultTitleColor {
		t.Errorf("Expected title color %q, got %q", DefaultTitleColor, config.Title.Color)
	}

	// Test description defaults
	if config.Description.Visible {
		t.Error("Expected description to be hidden by default")
	}

	if config.Description.Size != DefaultDescriptionFontSize {
		t.Errorf("Expected description size %f, got %f", DefaultDescriptionFontSize, config.Description.Size)
	}

	// Test overlay defaults
	if config.Overlay.Visible {
		t.Error("Expected overlay to be hidden by default")
	}
}

func TestLoadConfig_NonExistentFile(t *testing.T) {
	// Test loading non-existent config file (should return defaults)
	config, err := loadConfig("/nonexistent/path/config.yaml")
	if err != nil {
		t.Errorf("Expected no error for non-existent config file, got %v", err)
	}

	if config == nil {
		t.Error("Expected config to be non-nil even for non-existent file")
	}

	// Should return default values
	if config.Background.Color != DefaultBackgroundColor {
		t.Errorf("Expected default background color %q, got %q", DefaultBackgroundColor, config.Background.Color)
	}
}

func TestLoadConfig_ValidFile(t *testing.T) {
	// Create a temporary config file
	tempDir, err := os.MkdirTemp("", "config_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "config.yaml")
	configContent := `
background:
  color: "#FF0000"
output:
  directory: "./custom_output"
  format: "jpg"
title:
  visible: true
  size: 72.0
  color: "#FFFFFF"
description:
  visible: true
  size: 24.0
  color: "#CCCCCC"
`

	err = os.WriteFile(configPath, []byte(configContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write config file: %v", err)
	}

	config, err := loadConfig(configPath)
	if err != nil {
		t.Fatalf("Expected no error loading valid config, got %v", err)
	}

	if config.Background.Color != "#FF0000" {
		t.Errorf("Expected background color #FF0000, got %q", config.Background.Color)
	}

	if config.Output.Directory != "./custom_output" {
		t.Errorf("Expected output directory ./custom_output, got %q", config.Output.Directory)
	}

	if config.Output.Format != "jpg" {
		t.Errorf("Expected output format jpg, got %q", config.Output.Format)
	}

	if config.Title.Size != 72.0 {
		t.Errorf("Expected title size 72.0, got %f", config.Title.Size)
	}

	if config.Title.Color != "#FFFFFF" {
		t.Errorf("Expected title color #FFFFFF, got %q", config.Title.Color)
	}

	if !config.Description.Visible {
		t.Error("Expected description to be visible")
	}

	if config.Description.Size != 24.0 {
		t.Errorf("Expected description size 24.0, got %f", config.Description.Size)
	}
}

func TestLoadConfig_InvalidYAML(t *testing.T) {
	// Create a temporary invalid config file
	tempDir, err := os.MkdirTemp("", "config_invalid_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	configPath := filepath.Join(tempDir, "invalid_config.yaml")
	invalidContent := `
background:
  color: "#FF0000"
invalid yaml structure [
missing: quotation marks"
`

	err = os.WriteFile(configPath, []byte(invalidContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write invalid config file: %v", err)
	}

	_, err = loadConfig(configPath)
	if err == nil {
		t.Error("Expected error for invalid YAML config")
	}

	expectedError := "failed to unmarshal config"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain %q, got %q", expectedError, err.Error())
	}
}

func TestBuildProhibitedMaps(t *testing.T) {
	textConfig := &TextConfig{
		LineBreaking: LineBreakingConfig{
			StartProhibited: "。、！？",
			EndProhibited:   "（【「",
		},
	}

	startProhibited, endProhibited := buildProhibitedMaps(textConfig)

	// Test start prohibited characters
	if !startProhibited['。'] {
		t.Error("Expected '。' to be in start prohibited map")
	}
	if !startProhibited['、'] {
		t.Error("Expected '、' to be in start prohibited map")
	}
	if !startProhibited['！'] {
		t.Error("Expected '！' to be in start prohibited map")
	}
	if !startProhibited['？'] {
		t.Error("Expected '？' to be in start prohibited map")
	}

	// Test end prohibited characters
	if !endProhibited['（'] {
		t.Error("Expected '（' to be in end prohibited map")
	}
	if !endProhibited['【'] {
		t.Error("Expected '【' to be in end prohibited map")
	}
	if !endProhibited['「'] {
		t.Error("Expected '「' to be in end prohibited map")
	}

	// Test characters not in the maps
	if startProhibited['a'] {
		t.Error("Expected 'a' to not be in start prohibited map")
	}
	if endProhibited['a'] {
		t.Error("Expected 'a' to not be in end prohibited map")
	}

	// Test empty config
	emptyConfig := &TextConfig{
		LineBreaking: LineBreakingConfig{
			StartProhibited: "",
			EndProhibited:   "",
		},
	}

	emptyStart, emptyEnd := buildProhibitedMaps(emptyConfig)
	if len(emptyStart) != 0 {
		t.Errorf("Expected empty start prohibited map, got %d entries", len(emptyStart))
	}
	if len(emptyEnd) != 0 {
		t.Errorf("Expected empty end prohibited map, got %d entries", len(emptyEnd))
	}
}

func TestApplyFrontMatterOverrides(t *testing.T) {
	baseConfig := getDefaultConfig()

	// Test with nil OGP front matter
	result := applyFrontMatterOverrides(baseConfig, nil)
	if result == nil {
		t.Error("Expected non-nil result")
	}

	// Should be equivalent to base config
	if result.Title.Size != baseConfig.Title.Size {
		t.Errorf("Expected title size to remain %f, got %f", baseConfig.Title.Size, result.Title.Size)
	}

	// Test with simple override
	newSize := 48.0
	ogpFM := &OGPFrontMatter{
		Title: &TextConfigOverride{
			Size: &newSize,
		},
	}

	result = applyFrontMatterOverrides(baseConfig, ogpFM)
	if result.Title.Size != newSize {
		t.Errorf("Expected title size %f, got %f", newSize, result.Title.Size)
	}

	// Other values should remain unchanged
	if result.Title.Color != baseConfig.Title.Color {
		t.Errorf("Expected title color to remain %q, got %q", baseConfig.Title.Color, result.Title.Color)
	}
}
