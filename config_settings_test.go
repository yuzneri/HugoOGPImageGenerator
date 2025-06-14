package main

import (
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

// TestLoadConfigSettings tests loading ConfigSettings from YAML files
func TestLoadConfigSettings(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		expected    *ConfigSettings
		expectError bool
	}{
		{
			name: "Complete configuration",
			yamlContent: `title:
  visible: true
  size: 48.0
  color: "#FF0000"
  area:
    x: 100
    y: 200
overlay:
  visible: false
  placement:
    x: 50
    width: 400`,
			expected: &ConfigSettings{
				Title: &TextSettings{
					Visible: boolPtrSettings(true),
					Size:    float64PtrSettings(48.0),
					Color:   stringPtrSettings("#FF0000"),
					Area: &TextAreaSettings{
						X: intPtrSettings(100),
						Y: intPtrSettings(200),
					},
				},
				Overlay: &OverlayConfigSettings{
					Visible: boolPtrSettings(false),
					Placement: &PlacementSettings{
						X:     intPtrSettings(50),
						Width: intPtrSettings(400),
					},
				},
			},
		},
		{
			name: "Partial configuration",
			yamlContent: `title:
  size: 64.0
overlay:
  placement:
    height: 580`,
			expected: &ConfigSettings{
				Title: &TextSettings{
					Size: float64PtrSettings(64.0),
				},
				Overlay: &OverlayConfigSettings{
					Placement: &PlacementSettings{
						Height: intPtrSettings(580),
					},
				},
			},
		},
		{
			name:        "Empty configuration",
			yamlContent: ``,
			expected:    &ConfigSettings{},
		},
		{
			name: "Invalid YAML",
			yamlContent: `title:
  size: "invalid_number"`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create temporary file
			tmpFile, err := os.CreateTemp("", "config_test_*.yaml")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			// Write test content
			_, err = tmpFile.WriteString(tt.yamlContent)
			if err != nil {
				t.Fatalf("Failed to write test content: %v", err)
			}
			tmpFile.Close()

			// Test loadConfigSettings
			result, err := loadConfigSettings(tmpFile.Name())

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			if !configSettingsEqual(result, tt.expected) {
				t.Errorf("Result doesn't match expected.\nGot: %+v\nExpected: %+v", result, tt.expected)
			}
		})
	}
}

// TestLoadConfigSettingsNonExistent tests behavior when config file doesn't exist
func TestLoadConfigSettingsNonExistent(t *testing.T) {
	result, err := loadConfigSettings("nonexistent_file.yaml")

	if err != nil {
		t.Errorf("Expected no error for non-existent file, got: %v", err)
	}

	if result != nil {
		t.Errorf("Expected nil result for non-existent file, got: %+v", result)
	}
}

// TestLoadTypeConfigSettings tests loading type-specific ConfigSettings
func TestLoadTypeConfigSettings(t *testing.T) {
	// Create temporary directory
	tmpDir, err := os.MkdirTemp("", "config_test_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	// Create book.yaml
	bookConfigPath := filepath.Join(tmpDir, "book.yaml")
	bookContent := `title:
  area:
    x: 462
    y: 100
overlay:
  placement:
    x: 100
    height: 580`

	err = os.WriteFile(bookConfigPath, []byte(bookContent), 0644)
	if err != nil {
		t.Fatalf("Failed to write book config: %v", err)
	}

	// Test loading existing type config
	result, err := loadTypeConfigSettings(tmpDir, "Book")
	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	expected := &ConfigSettings{
		Title: &TextSettings{
			Area: &TextAreaSettings{
				X: intPtrSettings(462),
				Y: intPtrSettings(100),
			},
		},
		Overlay: &OverlayConfigSettings{
			Placement: &PlacementSettings{
				X:      intPtrSettings(100),
				Height: intPtrSettings(580),
			},
		},
	}

	if !configSettingsEqual(result, expected) {
		t.Errorf("Result doesn't match expected.\nGot: %+v\nExpected: %+v", result, expected)
	}

	// Test loading non-existent type config
	result, err = loadTypeConfigSettings(tmpDir, "NonExistent")
	if err != nil {
		t.Errorf("Expected no error for non-existent type, got: %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil result for non-existent type, got: %+v", result)
	}

	// Test empty type name
	result, err = loadTypeConfigSettings(tmpDir, "")
	if err != nil {
		t.Errorf("Expected no error for empty type name, got: %v", err)
	}
	if result != nil {
		t.Errorf("Expected nil result for empty type name, got: %+v", result)
	}
}

// Helper functions for creating pointers (avoid redeclaration)
func stringPtrSettings(s string) *string    { return &s }
func intPtrSettings(i int) *int             { return &i }
func float64PtrSettings(f float64) *float64 { return &f }
func boolPtrSettings(b bool) *bool          { return &b }

// configSettingsEqual compares two ConfigSettings for deep equality
func configSettingsEqual(a, b *ConfigSettings) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	return reflect.DeepEqual(a, b)
}
