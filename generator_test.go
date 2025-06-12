package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func TestSanitizeFilename(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "normal filename",
			input:    "normal-filename.png",
			expected: "normal-filename.png",
		},
		{
			name:     "filename with spaces",
			input:    "file name with spaces.png",
			expected: "file name with spaces.png",
		},
		{
			name:     "filename with path traversal",
			input:    "../../../etc/passwd",
			expected: "___etc_passwd",
		},
		{
			name:     "filename with dangerous characters",
			input:    "file<>:\"|?*\\/name.png",
			expected: "file_________name.png",
		},
		{
			name:     "filename with leading/trailing dots and spaces",
			input:    " ..filename.. ",
			expected: "filename",
		},
		{
			name:     "mixed dangerous characters",
			input:    " ../bad<file>name|test?.png ",
			expected: "_bad_file_name_test_.png",
		},
		{
			name:     "only dangerous characters",
			input:    "<>:\"|?*\\/",
			expected: "_________",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only dots and spaces",
			input:    " ... ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := sanitizeFilename(tt.input)
			if result != tt.expected {
				t.Errorf("sanitizeFilename(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestGenerateOutputFilename_Default(t *testing.T) {
	config := getDefaultConfig()
	config.Output.Format = "png"

	// Test with nil filename (should use default)
	config.Output.Filename = nil

	fm := &FrontMatter{
		Title:       "Test Article",
		Description: "Test description",
	}

	filename, err := generateOutputFilename(config, fm, "/content/test", "/content")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "ogp.png"
	if filename != expected {
		t.Errorf("Expected %q, got %q", expected, filename)
	}
}

func TestGenerateOutputFilename_EmptyTemplate(t *testing.T) {
	config := getDefaultConfig()
	config.Output.Format = "jpg"

	// Test with empty filename template (should use default)
	emptyTemplate := ""
	config.Output.Filename = &emptyTemplate

	fm := &FrontMatter{
		Title:       "Test Article",
		Description: "Test description",
	}

	filename, err := generateOutputFilename(config, fm, "/content/test", "/content")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "ogp.jpg"
	if filename != expected {
		t.Errorf("Expected %q, got %q", expected, filename)
	}
}

func TestGenerateOutputFilename_SimpleTemplate(t *testing.T) {
	config := getDefaultConfig()
	config.Output.Format = "png"

	// Test with simple template
	template := "{{.Title}}.{{.Format}}"
	config.Output.Filename = &template

	fm := &FrontMatter{
		Title:       "My Test Article",
		Description: "Test description",
		Date:        time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC),
	}

	filename, err := generateOutputFilename(config, fm, "/content/articles/test", "/content")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "My Test Article.png"
	if filename != expected {
		t.Errorf("Expected %q, got %q", expected, filename)
	}
}

func TestGenerateOutputFilename_ComplexTemplate(t *testing.T) {
	config := getDefaultConfig()
	config.Output.Format = "jpg"

	// Test with complex template including relpath
	template := "{{.RelPath}}-{{.Title}}.{{.Format}}"
	config.Output.Filename = &template

	fm := &FrontMatter{
		Title:       "Test Article",
		Description: "Test description",
		Fields: map[string]interface{}{
			"category": "technology",
		},
	}

	filename, err := generateOutputFilename(config, fm, "/content/articles/test", "/content")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	// RelPath should be "articles/test", but sanitized to "articles_test"
	expected := "articles_test-Test Article.jpg"
	if filename != expected {
		t.Errorf("Expected %q, got %q", expected, filename)
	}
}

func TestGenerateOutputFilename_WithFields(t *testing.T) {
	config := getDefaultConfig()
	config.Output.Format = "png"

	// Test template that uses Fields
	template := "{{.Fields.category}}-{{.Title}}.{{.Format}}"
	config.Output.Filename = &template

	fm := &FrontMatter{
		Title:       "Tech Article",
		Description: "A technology article",
		Fields: map[string]interface{}{
			"category": "technology",
			"author":   "John Doe",
		},
	}

	filename, err := generateOutputFilename(config, fm, "/content/tech/article", "/content")
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	expected := "technology-Tech Article.png"
	if filename != expected {
		t.Errorf("Expected %q, got %q", expected, filename)
	}
}

func TestNewOGPGenerator_ConfigLoadError(t *testing.T) {
	// Test with non-existent config file that contains invalid YAML
	tempDir, err := os.MkdirTemp("", "generator_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create an invalid config file
	configPath := filepath.Join(tempDir, "invalid_config.yaml")
	invalidYAML := "invalid: yaml: content\nthis is not valid yaml ["
	err = os.WriteFile(configPath, []byte(invalidYAML), 0644)
	if err != nil {
		t.Fatalf("Failed to create invalid config: %v", err)
	}

	contentDir := filepath.Join(tempDir, "content")
	projectRoot := tempDir

	_, err = NewOGPGenerator(configPath, contentDir, projectRoot)
	if err == nil {
		t.Error("Expected error when loading invalid config")
	}

	expectedError := "failed to load config"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain %q, got %q", expectedError, err.Error())
	}
}

func TestNewOGPGenerator_Success(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "generator_success_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Test with non-existent config (should use defaults)
	configPath := filepath.Join(tempDir, "nonexistent_config.yaml")
	contentDir := filepath.Join(tempDir, "content")
	projectRoot := tempDir

	generator, err := NewOGPGenerator(configPath, contentDir, projectRoot)
	if err != nil {
		t.Fatalf("Expected no error with default config, got %v", err)
	}

	if generator == nil {
		t.Error("Expected generator to be non-nil")
	}

	if generator.config == nil {
		t.Error("Expected config to be non-nil")
	}

	if generator.fontManager == nil {
		t.Error("Expected fontManager to be non-nil")
	}

	if generator.bgProcessor == nil {
		t.Error("Expected bgProcessor to be non-nil")
	}

	if generator.imageRenderer == nil {
		t.Error("Expected imageRenderer to be non-nil")
	}

	if generator.articleProcessor == nil {
		t.Error("Expected articleProcessor to be non-nil")
	}

	// Verify directory paths
	if generator.contentDir != contentDir {
		t.Errorf("Expected contentDir %q, got %q", contentDir, generator.contentDir)
	}

	if generator.projectRoot != projectRoot {
		t.Errorf("Expected projectRoot %q, got %q", projectRoot, generator.projectRoot)
	}

	expectedConfigDir := filepath.Dir(configPath)
	if generator.configDir != expectedConfigDir {
		t.Errorf("Expected configDir %q, got %q", expectedConfigDir, generator.configDir)
	}
}

func TestTemplateData_Structure(t *testing.T) {
	// Test TemplateData structure directly
	data := TemplateData{
		Title:       "Test Title",
		Description: "Test Description",
		Date:        time.Now(),
		URL:         "/test",
		RelPath:     "articles/test",
		Format:      "png",
		Fields: map[string]interface{}{
			"category": "test",
			"author":   "Test Author",
		},
	}

	if data.Title != "Test Title" {
		t.Errorf("Expected Title 'Test Title', got %q", data.Title)
	}

	if data.Description != "Test Description" {
		t.Errorf("Expected Description 'Test Description', got %q", data.Description)
	}

	if data.URL != "/test" {
		t.Errorf("Expected URL '/test', got %q", data.URL)
	}

	if data.RelPath != "articles/test" {
		t.Errorf("Expected RelPath 'articles/test', got %q", data.RelPath)
	}

	if data.Format != "png" {
		t.Errorf("Expected Format 'png', got %q", data.Format)
	}

	if data.Fields["category"] != "test" {
		t.Errorf("Expected category 'test', got %v", data.Fields["category"])
	}

	if data.Fields["author"] != "Test Author" {
		t.Errorf("Expected author 'Test Author', got %v", data.Fields["author"])
	}
}
