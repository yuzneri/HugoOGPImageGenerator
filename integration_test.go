package main

import (
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"
)

// TestArticleProcessor_Integration tests the integration between ArticleProcessor and its dependencies
func TestArticleProcessor_Integration(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "integration_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	contentDir := filepath.Join(tempDir, "content")
	configDir := tempDir
	articleDir := filepath.Join(contentDir, "test-article")

	err = os.MkdirAll(articleDir, DefaultFilePermission)
	if err != nil {
		t.Fatalf("Failed to create article dir: %v", err)
	}

	// Create a test article with front matter
	indexContent := `---
title: "Integration Test Article"
description: "This is a test article for integration testing"
date: 2023-12-25T15:30:45Z
tags: ["test", "integration"]
ogp:
  title:
    content: "Custom OGP Title: {{.Title}}"
    size: 48
  description:
    visible: true
    content: "Published on {{dateFormat \"2006-01-02\" .Date}}"
---

# Integration Test Article

This is the content of the test article.
`
	indexPath := filepath.Join(articleDir, DefaultIndexFilename)
	err = os.WriteFile(indexPath, []byte(indexContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create index file: %v", err)
	}

	// Create a basic config
	config := getDefaultConfig()
	config.Background.Color = DefaultBackgroundColor
	config.Title.Size = DefaultTitleFontSize
	config.Description.Size = DefaultDescriptionFontSize

	// Create processors
	configPath := filepath.Join(configDir, "config.yaml")
	fontManager := NewFontManager(configDir)
	bgProcessor := NewBackgroundProcessor(configDir)
	imageRenderer := NewImageRenderer()
	articleProcessor := NewArticleProcessor(config, contentDir, configDir, configPath, fontManager, bgProcessor, imageRenderer)

	// Test processing in test mode
	options := ProcessOptions{TestMode: true}
	err = articleProcessor.ProcessArticle(articleDir, options)

	if err != nil {
		t.Errorf("ProcessArticle should not return error: %v", err)
	}

	// Test that output file was created
	outputPath := articleProcessor.generateTestOutputPath()
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Expected output file to be created")
	}

	// Clean up output file
	os.Remove(outputPath)
}

// TestTemplateProcessor_Integration tests template processing with real front matter
func TestTemplateProcessor_Integration(t *testing.T) {
	tp := NewTemplateProcessor()

	// Test front matter data
	fm := &FrontMatter{
		Title:       "My Test Article",
		Description: "A comprehensive test of template processing",
		Date:        time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC),
		URL:         "/articles/test",
		Tags:        []string{"golang", "testing", "templates"},
		Fields: map[string]interface{}{
			"author":   "Test Author",
			"category": "Technology",
			"rating":   4.5,
		},
	}

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "Basic field substitution",
			template: "{{.Title}} by {{.Fields.author}}",
			expected: "My Test Article by Test Author",
		},
		{
			name:     "Date formatting",
			template: "Published: {{dateFormat \"January 2, 2006\" .Date}}",
			expected: "Published: December 25, 2023",
		},
		{
			name:     "String manipulation",
			template: "{{upper .Fields.category}} - {{lower .Title}}",
			expected: "TECHNOLOGY - my test article",
		},
		{
			name:     "Default function",
			template: "Status: {{default \"Draft\" .Fields.status}}",
			expected: "Status: Draft",
		},
		{
			name:     "Complex template",
			template: "Article: {{.Title}} ({{.Fields.category}}) - {{dateFormat \"2006-01-02\" .Date}}",
			expected: "Article: My Test Article (Technology) - 2023-12-25",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tp.ProcessContentTemplate(tt.template, fm)
			if err != nil {
				t.Errorf("ProcessContentTemplate failed: %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

// TestErrorHandling_Integration tests error propagation through the system
func TestErrorHandling_Integration(t *testing.T) {
	// Create processors with non-existent directories to trigger errors
	fontManager := NewFontManager("/nonexistent")
	bgProcessor := NewBackgroundProcessor("/nonexistent")
	imageRenderer := NewImageRenderer()

	config := getDefaultConfig()
	config.Background.Image = StringPtr("nonexistent.png")

	articleProcessor := NewArticleProcessor(config, "/nonexistent", "/nonexistent", "/nonexistent/config.yaml", fontManager, bgProcessor, imageRenderer)

	// Test that errors propagate correctly
	err := articleProcessor.ProcessArticle("/nonexistent/article", ProcessOptions{TestMode: true})

	if err == nil {
		t.Error("Expected error when processing non-existent article")
	}

	// Should be a FileError for missing index.md
	if !IsFileError(err) {
		t.Errorf("Expected FileError, got %T", err)
	}
}

// TestConfigMerger_Integration tests config merging with real scenarios
func TestConfigMerger_Integration(t *testing.T) {
	merger := NewConfigMerger()
	baseConfig := getDefaultConfig()

	// Test complex override scenario
	titleContent := "Custom Title: {{.Title}}"
	descContent := "By {{.Fields.author}} on {{dateFormat \"2006-01-02\" .Date}}"
	newFontSize := 56.0
	visible := true

	ogpFM := &OGPFrontMatter{
		Title: &TextConfigOverride{
			Content: &titleContent,
			Size:    &newFontSize,
		},
		Description: &TextConfigOverride{
			Visible: &visible,
			Content: &descContent,
		},
		Background: &BackgroundOverride{
			Color: StringPtr("#f0f0f0"),
		},
	}

	result := merger.MergeConfigs(baseConfig, ogpFM)

	// Verify title overrides
	if result.Title.Content == nil || *result.Title.Content != titleContent {
		t.Errorf("Expected title content %q, got %v", titleContent, result.Title.Content)
	}

	if result.Title.Size != newFontSize {
		t.Errorf("Expected title size %f, got %f", newFontSize, result.Title.Size)
	}

	// Verify description overrides
	if !result.Description.Visible {
		t.Error("Expected description to be visible")
	}

	if result.Description.Content == nil || *result.Description.Content != descContent {
		t.Errorf("Expected description content %q, got %v", descContent, result.Description.Content)
	}

	// Verify background override
	if result.Background.Color != "#f0f0f0" {
		t.Errorf("Expected background color #f0f0f0, got %s", result.Background.Color)
	}

	// Verify that non-overridden values remain unchanged
	if result.Title.Color != DefaultTitleColor {
		t.Errorf("Expected title color to remain %s, got %s", DefaultTitleColor, result.Title.Color)
	}
}

// TestEndToEnd_Integration tests the complete pipeline from config to image generation
func TestEndToEnd_Integration(t *testing.T) {
	// Create a temporary directory structure
	tempDir, err := os.MkdirTemp("", "e2e_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	contentDir := filepath.Join(tempDir, "content")
	articleDir := filepath.Join(contentDir, "e2e-article")

	err = os.MkdirAll(articleDir, DefaultFilePermission)
	if err != nil {
		t.Fatalf("Failed to create article dir: %v", err)
	}

	// Create article with comprehensive front matter
	indexContent := `---
title: "End-to-End Test"
description: "Complete pipeline testing"
date: 2023-12-25T10:00:00Z
tags: ["e2e", "test"]
ogp:
  title:
    content: "{{upper .Title}} - {{.Fields.category}}"
    size: 52
    color: "#2c3e50"
  description:
    visible: true
    content: "{{.Description}} | {{dateFormat \"Jan 2006\" .Date}}"
    size: 28
  background:
    color: "#ecf0f1"
category: "Testing"
---

# End-to-End Test

This tests the complete pipeline.
`

	indexPath := filepath.Join(articleDir, DefaultIndexFilename)
	err = os.WriteFile(indexPath, []byte(indexContent), 0644)
	if err != nil {
		t.Fatalf("Failed to create index file: %v", err)
	}

	// Create OGP generator
	generator, err := NewOGPGenerator("", contentDir, tempDir)
	if err != nil {
		t.Fatalf("Failed to create OGP generator: %v", err)
	}

	// Test complete generation
	err = generator.GenerateTest(articleDir)
	if err != nil {
		t.Errorf("GenerateTest should not return error: %v", err)
	}

	// Verify output file exists
	outputPath := generator.articleProcessor.generateTestOutputPath()
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		t.Error("Expected output file to be created")
	}

	// Clean up
	os.Remove(outputPath)
}

// Helper function to create string pointer
func StringPtr(s string) *string {
	return &s
}

// TestTemplateProcessor_ErrorHandling tests error handling in template processing
func TestTemplateProcessor_ErrorHandling(t *testing.T) {
	tp := NewTemplateProcessor()

	// Test invalid template syntax
	invalidTemplate := "{{.Title"
	fm := &FrontMatter{Title: "Test"}

	_, err := tp.ProcessContentTemplate(invalidTemplate, fm)
	if err == nil {
		t.Error("Expected error for invalid template")
	}

	// Check error type
	var appErr *AppError
	if !errors.As(err, &appErr) {
		t.Error("Expected AppError type")
	} else {
		if appErr.Type != TemplateError {
			t.Errorf("Expected TemplateError, got %v", appErr.Type)
		}
	}
}

// TestFontManager_Integration tests font loading with fallback scenarios
func TestFontManager_Integration(t *testing.T) {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "font_integration")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	fm := NewFontManager(tempDir)

	// Test loading default font (empty path)
	font, err := fm.LoadFont("", "")
	if err != nil {
		t.Errorf("Loading default font should not return error: %v", err)
	}

	if font == nil {
		t.Error("Default font should not be nil")
	}

	// Test with fallback
	defaultFont, _ := fm.LoadFont("", "")
	fallbackFont := fm.LoadFontWithFallback("nonexistent.ttf", "", defaultFont)

	if fallbackFont != defaultFont {
		t.Error("Should return fallback font when original font fails to load")
	}
}

// TestGenerateAll_HugoPageKinds tests that every kind of Hugo page gets an image at its permalink path.
func TestGenerateAll_HugoPageKinds(t *testing.T) {
	projectRoot := t.TempDir()
	contentDir := filepath.Join(projectRoot, "content")

	pages := map[string]string{
		"_index.md":                 "Home",
		"about.md":                  "About",
		"events/_index.md":          "Events",
		"events/2026-10-meetup.md":  "Meetup",
		"events/custom-url.md":      "Custom URL\nurl: /special/",
		"posts/bundle/index.md":     "Bundle",
		"posts/bundle/resource.md":  "Resource",
		"posts/bundle/sub/note.md":  "Nested resource",
		"posts/nested/deep/page.md": "Deep",
	}
	for path, title := range pages {
		full := filepath.Join(contentDir, path)
		if err := os.MkdirAll(filepath.Dir(full), DefaultFilePermission); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte("---\ntitle: "+title+"\n---\n"), 0644); err != nil {
			t.Fatal(err)
		}
	}

	// Pages without YAML front matter are skipped rather than aborting the whole run.
	skipped := map[string]string{
		"no-front-matter.md": "# Plain page\n",
		"toml/_index.md":     "+++\ntitle = \"TOML\"\n+++\n",
	}
	for path, body := range skipped {
		full := filepath.Join(contentDir, path)
		if err := os.MkdirAll(filepath.Dir(full), DefaultFilePermission); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}

	configPath := filepath.Join(projectRoot, "config.yaml")
	if err := os.WriteFile(configPath, []byte("output:\n  directory: public\n"), 0644); err != nil {
		t.Fatal(err)
	}

	generator, err := NewOGPGenerator(configPath, contentDir, projectRoot)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}
	if err := generator.GenerateAll(); err != nil {
		t.Fatalf("GenerateAll should not return error: %v", err)
	}

	var got []string
	err = filepath.WalkDir(filepath.Join(projectRoot, "public"), func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(projectRoot, path)
		got = append(got, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		"public/about/ogp.png",
		"public/events/2026-10-meetup/ogp.png",
		"public/events/ogp.png",
		"public/ogp.png",
		"public/posts/bundle/ogp.png",
		"public/posts/nested/deep/page/ogp.png",
		"public/special/ogp.png",
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("generated files:\n got %v\nwant %v", got, want)
	}
}

// TestGenerateSingle_File tests that a single-file page can be specified directly.
func TestGenerateSingle_File(t *testing.T) {
	projectRoot := t.TempDir()
	contentDir := filepath.Join(projectRoot, "content")
	page := filepath.Join(contentDir, "events", "meetup.md")
	if err := os.MkdirAll(filepath.Dir(page), DefaultFilePermission); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(page, []byte("---\ntitle: Meetup\n---\n"), 0644); err != nil {
		t.Fatal(err)
	}

	configPath := filepath.Join(projectRoot, "config.yaml")
	generator, err := NewOGPGenerator(configPath, contentDir, projectRoot)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}
	if err := generator.GenerateSingle("events/meetup.md"); err != nil {
		t.Fatalf("GenerateSingle should not return error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(projectRoot, "public", "events", "meetup", "ogp.png")); err != nil {
		t.Errorf("Expected output file: %v", err)
	}
}

// TestGenerateAll_Exclusion tests excluding pages by kind (global and type config) and by front matter.
func TestGenerateAll_Exclusion(t *testing.T) {
	projectRoot := t.TempDir()
	contentDir := filepath.Join(projectRoot, "content")

	files := map[string]string{
		"content/_index.md":        "---\ntitle: Home\n---\n",
		"content/about.md":         "---\ntitle: About\n---\n",
		"content/hidden.md":        "---\ntitle: Hidden\nogp:\n  enabled: false\n---\n",
		"content/draft.md":         "---\ntitle: Draft\ndraft: true\n---\n",
		"content/posts/_index.md":  "---\ntitle: Posts\n---\n",
		"content/posts/post.md":    "---\ntitle: Post\n---\n",
		"content/events/_index.md": "---\ntitle: Events\n---\n",
		"content/events/meetup.md": "---\ntitle: Meetup\n---\n",
		"ogp/config.yaml":          "kinds: [home, page]\n",
		"ogp/events.yaml":          "kinds: [section, page]\n",
		"ogp/posts.yaml":           "title:\n  size: 40\n",
	}
	for path, body := range files {
		full := filepath.Join(projectRoot, path)
		if err := os.MkdirAll(filepath.Dir(full), DefaultFilePermission); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(full, []byte(body), 0644); err != nil {
			t.Fatal(err)
		}
	}

	generator, err := NewOGPGenerator(filepath.Join(projectRoot, "ogp", "config.yaml"), contentDir, projectRoot)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}
	if err := generator.GenerateAll(); err != nil {
		t.Fatalf("GenerateAll should not return error: %v", err)
	}

	var got []string
	err = filepath.WalkDir(filepath.Join(projectRoot, "public"), func(path string, d os.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		rel, _ := filepath.Rel(projectRoot, path)
		got = append(got, filepath.ToSlash(rel))
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}

	want := []string{
		"public/about/ogp.png",
		"public/events/meetup/ogp.png",
		"public/events/ogp.png", // events.yaml re-enables sections
		"public/ogp.png",
		"public/posts/post/ogp.png", // posts/_index.md excluded by the global kinds
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("generated files:\n got %v\nwant %v", got, want)
	}
}

// TestGenerateAll_UnknownKind tests that a typo in kinds is reported instead of silently generating nothing.
func TestGenerateAll_UnknownKind(t *testing.T) {
	projectRoot := t.TempDir()
	contentDir := filepath.Join(projectRoot, "content")
	if err := os.MkdirAll(contentDir, DefaultFilePermission); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(contentDir, "about.md"), []byte("---\ntitle: About\n---\n"), 0644); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(projectRoot, "config.yaml")
	if err := os.WriteFile(configPath, []byte("kinds: [pages]\n"), 0644); err != nil {
		t.Fatal(err)
	}

	generator, err := NewOGPGenerator(configPath, contentDir, projectRoot)
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}
	err = generator.GenerateAll()
	if err == nil || !strings.Contains(err.Error(), `unknown page kind "pages"`) {
		t.Errorf("expected unknown page kind error, got %v", err)
	}
}

// TestGenerateSingle_AbsolutePathWithRelativeRoot tests an absolute article path combined with a relative project root.
func TestGenerateSingle_AbsolutePathWithRelativeRoot(t *testing.T) {
	projectRoot := t.TempDir()
	contentDir := filepath.Join(projectRoot, "content")
	if err := os.MkdirAll(contentDir, DefaultFilePermission); err != nil {
		t.Fatal(err)
	}
	home := filepath.Join(contentDir, "_index.md")
	if err := os.WriteFile(home, []byte("---\ntitle: Home\n---\n"), 0644); err != nil {
		t.Fatal(err)
	}
	configPath := filepath.Join(projectRoot, "config.yaml")
	if err := os.WriteFile(configPath, []byte("kinds: [home]\n"), 0644); err != nil {
		t.Fatal(err)
	}

	cwd, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}
	if err := os.Chdir(projectRoot); err != nil {
		t.Fatal(err)
	}
	defer os.Chdir(cwd)

	generator, err := NewOGPGenerator(configPath, "content", ".")
	if err != nil {
		t.Fatalf("Failed to create generator: %v", err)
	}
	if err := generator.GenerateSingle(home); err != nil {
		t.Fatalf("GenerateSingle should not return error: %v", err)
	}

	if _, err := os.Stat(filepath.Join(projectRoot, "public", "ogp.png")); err != nil {
		t.Errorf("Expected home page image: %v", err)
	}
}
