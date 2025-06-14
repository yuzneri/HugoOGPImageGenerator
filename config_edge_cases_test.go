package main

import (
	"os"
	"testing"
)

// TestEdgeCasesAndErrorHandling tests edge cases and error conditions
func TestEdgeCasesAndErrorHandling(t *testing.T) {
	t.Run("Nonexistent config directory", func(t *testing.T) {
		result, err := loadTypeConfigSettings("/nonexistent/directory", "Book")
		// This should not error - nonexistent directories just return nil
		if err != nil {
			t.Errorf("Unexpected error for nonexistent directory: %v", err)
		}
		if result != nil {
			t.Errorf("Expected nil result for nonexistent directory, got: %+v", result)
		}
	})

	t.Run("Empty type name", func(t *testing.T) {
		tmpDir, err := os.MkdirTemp("", "config_test_")
		if err != nil {
			t.Fatalf("Failed to create temp dir: %v", err)
		}
		defer os.RemoveAll(tmpDir)

		result, err := loadTypeConfigSettings(tmpDir, "")
		if err != nil {
			t.Errorf("Expected no error for empty type name, got: %v", err)
		}
		if result != nil {
			t.Errorf("Expected nil result for empty type name, got: %+v", result)
		}
	})

	t.Run("Malformed YAML file", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "malformed_*.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		malformedYAML := `title:
  size: 48
  invalid_structure:
    - item1
    - item2: {nested: value
  color: "#FF0000"`

		_, err = tmpFile.WriteString(malformedYAML)
		if err != nil {
			t.Fatalf("Failed to write malformed YAML: %v", err)
		}
		tmpFile.Close()

		result, err := loadConfigSettings(tmpFile.Name())
		if err == nil {
			t.Errorf("Expected error for malformed YAML, got none")
		}
		if result != nil {
			t.Errorf("Expected nil result for malformed YAML, got: %+v", result)
		}
	})

	t.Run("Very large values", func(t *testing.T) {
		tmpFile, err := os.CreateTemp("", "large_values_*.yaml")
		if err != nil {
			t.Fatalf("Failed to create temp file: %v", err)
		}
		defer os.Remove(tmpFile.Name())

		largeValuesYAML := `title:
  size: 999999.99
overlay:
  placement:
    x: 2147483647
    y: -2147483648
    width: 1000000
    height: 1000000`

		_, err = tmpFile.WriteString(largeValuesYAML)
		if err != nil {
			t.Fatalf("Failed to write large values YAML: %v", err)
		}
		tmpFile.Close()

		result, err := loadConfigSettings(tmpFile.Name())
		if err != nil {
			t.Errorf("Unexpected error for large values: %v", err)
		}

		expected := &ConfigSettings{
			Title: &TextSettings{
				Size: float64Ptr(999999.99),
			},
			Overlay: &OverlayConfigSettings{
				Placement: &PlacementSettings{
					X:      intPtr(2147483647),
					Y:      intPtr(-2147483648),
					Width:  intPtr(1000000),
					Height: intPtr(1000000),
				},
			},
		}

		if !configSettingsEqual(result, expected) {
			t.Errorf("Large values not handled correctly.\nGot: %+v\nExpected: %+v", result, expected)
		}
	})
}

// TestNilValueHandling tests how nil values are handled throughout the system
func TestNilValueHandling(t *testing.T) {
	t.Run("Nil ConfigMerger", func(t *testing.T) {
		var merger *ConfigMerger = nil

		// This should panic, so we recover and check
		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic when using nil ConfigMerger")
			}
		}()

		merger.applySettingsToConfig(&Config{}, &ConfigSettings{})
	})

	t.Run("Nil target config", func(t *testing.T) {
		merger := NewConfigMerger()

		defer func() {
			if r := recover(); r == nil {
				t.Errorf("Expected panic when target config is nil")
			}
		}()

		merger.applySettingsToConfig(nil, &ConfigSettings{})
	})

	t.Run("Deeply nested nil values", func(t *testing.T) {
		settings := &ConfigSettings{
			Title: &TextSettings{
				Area: nil, // Nil nested structure
			},
			Overlay: &OverlayConfigSettings{
				Placement: nil, // Nil nested structure
			},
		}

		baseConfig := &Config{
			Title: TextConfig{
				Area: TextArea{X: 100, Y: 100},
			},
			Overlay: MainOverlayConfig{
				Placement: PlacementConfig{X: 50, Y: 50},
			},
		}

		merger := NewConfigMerger()
		result := deepCopyConfig(baseConfig)
		merger.applySettingsToConfig(result, settings)

		// Nil nested structures should not affect existing values
		if result.Title.Area.X != 100 || result.Title.Area.Y != 100 {
			t.Errorf("Nil nested structure affected existing values")
		}
		if result.Overlay.Placement.X != 50 || result.Overlay.Placement.Y != 50 {
			t.Errorf("Nil nested structure affected existing placement values")
		}
	})
}

// TestConcurrentAccess tests thread safety (basic test)
func TestConcurrentAccess(t *testing.T) {
	// Create test config file
	tmpFile, err := os.CreateTemp("", "concurrent_test_*.yaml")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	yamlContent := `title:
  size: 48.0
  color: "#FF0000"`

	_, err = tmpFile.WriteString(yamlContent)
	if err != nil {
		t.Fatalf("Failed to write YAML content: %v", err)
	}
	tmpFile.Close()

	// Run multiple goroutines reading the same config
	const numGoroutines = 10
	results := make(chan *ConfigSettings, numGoroutines)
	errors := make(chan error, numGoroutines)

	for i := 0; i < numGoroutines; i++ {
		go func() {
			result, err := loadConfigSettings(tmpFile.Name())
			if err != nil {
				errors <- err
				return
			}
			results <- result
		}()
	}

	// Collect results
	for i := 0; i < numGoroutines; i++ {
		select {
		case err := <-errors:
			t.Errorf("Concurrent access error: %v", err)
		case result := <-results:
			if result.Title == nil || result.Title.Size == nil || *result.Title.Size != 48.0 {
				t.Errorf("Concurrent access produced incorrect result: %+v", result)
			}
		}
	}
}

// TestMemoryLeaks tests for potential memory leaks in pointer handling
func TestMemoryLeaks(t *testing.T) {
	merger := NewConfigMerger()

	// Create a large number of configs and apply settings
	for i := 0; i < 1000; i++ {
		baseConfig := &Config{
			Title: TextConfig{
				Content: stringPtr("initial"),
			},
		}

		settings := &ConfigSettings{
			Title: &TextSettings{
				Content: stringPtr("modified"),
				Size:    float64Ptr(float64(i)),
			},
		}

		merger.applySettingsToConfig(baseConfig, settings)

		// Verify the pointer was properly copied, not shared
		if baseConfig.Title.Content == settings.Title.Content {
			t.Errorf("Pointer sharing detected - potential memory leak")
		}

		if *baseConfig.Title.Content != "modified" {
			t.Errorf("Content not properly applied: got %s, expected modified", *baseConfig.Title.Content)
		}
	}
}

// TestTypeConfigPathSecurity tests path security measures
func TestTypeConfigPathSecurity(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "security_test_")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tmpDir)

	tests := []struct {
		name        string
		typeName    string
		expectError bool
	}{
		{
			name:        "Normal type name",
			typeName:    "Book",
			expectError: false,
		},
		{
			name:        "Path traversal attempt",
			typeName:    "../../../etc/passwd",
			expectError: true,
		},
		{
			name:        "Null byte injection",
			typeName:    "Book\x00.yaml",
			expectError: true,
		},
		{
			name:        "Very long type name",
			typeName:    string(make([]byte, 1000)),
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := loadTypeConfigSettings(tmpDir, tt.typeName)

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for dangerous type name %q, got none", tt.typeName)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for safe type name %q: %v", tt.typeName, err)
				}
				// For safe names with non-existent files, result should be nil
				if result != nil {
					t.Errorf("Expected nil result for non-existent safe type, got: %+v", result)
				}
			}
		})
	}
}

// TestConfigurationBoundaries tests boundary conditions for configuration values
func TestConfigurationBoundaries(t *testing.T) {
	tests := []struct {
		name        string
		yamlContent string
		expectError bool
	}{
		{
			name: "Zero values",
			yamlContent: `title:
  size: 0.0
overlay:
  placement:
    x: 0
    y: 0
    width: 0
    height: 0`,
			expectError: false,
		},
		{
			name: "Negative values",
			yamlContent: `title:
  size: -10.0
overlay:
  placement:
    x: -100
    y: -200`,
			expectError: false,
		},
		{
			name: "Empty strings",
			yamlContent: `title:
  color: ""
  font: ""`,
			expectError: false,
		},
		{
			name: "Very long strings",
			yamlContent: `title:
  color: "` + string(make([]byte, 10000)) + `"`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tmpFile, err := os.CreateTemp("", "boundary_test_*.yaml")
			if err != nil {
				t.Fatalf("Failed to create temp file: %v", err)
			}
			defer os.Remove(tmpFile.Name())

			_, err = tmpFile.WriteString(tt.yamlContent)
			if err != nil {
				t.Fatalf("Failed to write test content: %v", err)
			}
			tmpFile.Close()

			result, err := loadConfigSettings(tmpFile.Name())

			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error for boundary case, got none")
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for boundary case: %v", err)
				}
				if result == nil {
					t.Errorf("Expected valid result for boundary case, got nil")
				}
			}
		})
	}
}
