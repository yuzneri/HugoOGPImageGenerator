package main

import (
	"bytes"
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"
	"strings"
	tmpl "text/template"
	"time"
)

// ArticleProcessor handles the complete processing pipeline for individual articles.
// It coordinates front matter parsing, configuration merging, and image generation.
type ArticleProcessor struct {
	config        *Config
	contentDir    string
	fontManager   *FontManager
	bgProcessor   *BackgroundProcessor
	imageRenderer *ImageRenderer
	configDir     string
}

// NewArticleProcessor creates a new ArticleProcessor with the given dependencies.
func NewArticleProcessor(config *Config, contentDir, configDir string, fontManager *FontManager, bgProcessor *BackgroundProcessor, imageRenderer *ImageRenderer) *ArticleProcessor {
	return &ArticleProcessor{
		config:        config,
		contentDir:    contentDir,
		fontManager:   fontManager,
		bgProcessor:   bgProcessor,
		imageRenderer: imageRenderer,
		configDir:     configDir,
	}
}

// ProcessOptions controls how articles are processed.
type ProcessOptions struct {
	TestMode  bool   // Generate test output to temporary location
	OutputDir string // Override default output directory
}

// ProcessArticle processes a single article and generates its OGP image.
// It reads the front matter, applies configuration overrides, and orchestrates the rendering pipeline.
func (ap *ArticleProcessor) ProcessArticle(articlePath string, options ProcessOptions) error {
	indexPath := filepath.Join(articlePath, "index.md")
	content, err := os.ReadFile(indexPath)
	if err != nil {
		return fmt.Errorf("failed to read %s: %w", indexPath, err)
	}

	fm, err := parseFrontMatter(content)
	if err != nil {
		return fmt.Errorf("failed to parse front matter in %s: %w", indexPath, err)
	}

	overrideConfig := applyFrontMatterOverrides(ap.config, fm.OGP)
	title := ap.determineTitle(fm, overrideConfig)

	if title == "" {
		return fmt.Errorf("no title found in %s (checked: front matter title, config content template, and OGP content)", indexPath)
	}

	var outputPath string
	if options.TestMode {
		outputPath = ap.generateTestOutputPath()
	} else {
		var err error
		outputPath, err = ap.generateProductionOutputPath(overrideConfig, fm, articlePath, options.OutputDir)
		if err != nil {
			return err
		}
	}

	err = ap.generateImage(title, outputPath, overrideConfig, articlePath, fm.OGP, options.TestMode)
	if err != nil {
		return fmt.Errorf("failed to generate OGP image for %s: %w", title, err)
	}

	ap.logSuccess(articlePath, outputPath, options.TestMode)
	return nil
}

// determineTitle resolves the final title to use for image generation.
// It applies the priority: front matter OGP content > config template > article title.
func (ap *ArticleProcessor) determineTitle(fm *FrontMatter, config *Config) string {
	// Priority 1: Front matter OGP content (highest)
	if fm.OGP != nil && fm.OGP.Text != nil && fm.OGP.Text.Content != nil {
		return ap.processContentTemplate(*fm.OGP.Text.Content, fm)
	}

	// Priority 2: Config template (medium)
	if config.Text.Content != nil {
		return ap.processContentTemplate(*config.Text.Content, fm)
	}

	// Priority 3: Article title (lowest)
	return fm.Title
}

// processContentTemplate processes content template with front matter data.
func (ap *ArticleProcessor) processContentTemplate(contentTemplate string, fm *FrontMatter) string {
	// If no template markers, return as-is
	if !strings.Contains(contentTemplate, "{{") {
		return contentTemplate
	}

	// Prepare template data with proper date handling
	data := TemplateData{
		Title:       fm.Title,
		Description: fm.Description,
		Date:        ap.parseDate(fm.Date),
		URL:         fm.URL,
		Fields:      make(map[string]interface{}),
	}

	// Add all front matter fields
	if fm.Fields != nil {
		data.Fields = fm.Fields
	}

	// Add standard fields to Fields map for template access
	data.Fields["title"] = fm.Title
	data.Fields["description"] = fm.Description
	data.Fields["date"] = ap.parseDate(fm.Date)
	data.Fields["url"] = fm.URL
	data.Fields["tags"] = fm.Tags

	// Parse and execute template with Hugo-like functions
	funcMap := tmpl.FuncMap{
		// Hugo-compatible functions
		"default": func(defaultValue interface{}, value interface{}) interface{} {
			if value == nil || value == "" {
				return defaultValue
			}
			return value
		},
		"upper": strings.ToUpper,
		"lower": strings.ToLower,
		"title": strings.Title,
		// Date functions like Hugo
		"dateFormat": func(layout string, date interface{}) string {
			if date == nil {
				return ""
			}
			if t, ok := date.(time.Time); ok {
				return t.Format(layout)
			}
			return fmt.Sprintf("%v", date)
		},
		"now": time.Now,
		// String functions
		"replace": func(old, new, s string) string {
			return strings.ReplaceAll(s, old, new)
		},
		"split": strings.Split,
		"trim":  strings.TrimSpace,
	}

	t, err := tmpl.New("content").Funcs(funcMap).Parse(contentTemplate)
	if err != nil {
		DefaultLogger.Warning("Failed to parse content template: %v", err)
		return contentTemplate
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		DefaultLogger.Warning("Failed to execute content template: %v", err)
		return contentTemplate
	}

	return buf.String()
}

// parseDate converts various date formats to time.Time for template usage.
func (ap *ArticleProcessor) parseDate(dateValue interface{}) interface{} {
	if dateValue == nil {
		return nil
	}

	// If already a time.Time, return as-is
	if t, ok := dateValue.(time.Time); ok {
		return t
	}

	// Try to parse string dates
	if dateStr, ok := dateValue.(string); ok {
		// Common date formats
		formats := []string{
			time.RFC3339,
			"2006-01-02T15:04:05Z07:00",
			"2006-01-02T15:04:05",
			"2006-01-02 15:04:05",
			"2006-01-02",
		}

		for _, format := range formats {
			if t, err := time.Parse(format, dateStr); err == nil {
				return t
			}
		}
	}

	// Return original value if parsing fails
	return dateValue
}

// generateTestOutputPath creates a temporary output path for test mode.
func (ap *ArticleProcessor) generateTestOutputPath() string {
	execPath, execErr := os.Executable()
	if execErr != nil {
		execDir, _ := os.Getwd()
		return filepath.Join(execDir, "test.png")
	}

	execDir := filepath.Dir(execPath)
	return filepath.Join(execDir, "test.png")
}

// generateProductionOutputPath creates the final output path for production mode.
// It respects custom URLs and creates the necessary directory structure.
func (ap *ArticleProcessor) generateProductionOutputPath(config *Config, fm *FrontMatter, articlePath, outputDir string) (string, error) {
	filename, err := generateOutputFilename(config, fm, articlePath, ap.contentDir)
	if err != nil {
		return "", fmt.Errorf("failed to generate output filename: %w", err)
	}

	var outputPath string
	if outputDir == "" {
		outputDir = filepath.Join(filepath.Dir(ap.contentDir), config.Output.Directory)
	}

	if fm.URL != "" {
		outputPath = filepath.Join(outputDir, fm.URL, filename)
	} else {
		relPath, err := filepath.Rel(ap.contentDir, articlePath)
		if err != nil {
			return "", fmt.Errorf("failed to get relative path: %w", err)
		}
		outputPath = filepath.Join(outputDir, relPath, filename)
	}

	err = os.MkdirAll(filepath.Dir(outputPath), 0755)
	if err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	return outputPath, nil
}

// generateImage creates the OGP image by compositing background, overlays, and text.
// It handles both config-level and article-level overlay compositions.
func (ap *ArticleProcessor) generateImage(title, outputPath string, config *Config, articlePath string, ogpSettings *OGPFrontMatter, testMode bool) error {
	backgroundImage, err := ap.bgProcessor.CreateBackground(config, articlePath)
	if err != nil {
		return fmt.Errorf("failed to create background: %w", err)
	}

	font, err := ap.fontManager.LoadFont(config.Text.Font, articlePath)
	if err != nil {
		return fmt.Errorf("failed to load font: %w", err)
	}

	bounds := backgroundImage.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, backgroundImage, image.Point{}, draw.Src)

	if ap.config.Overlay != nil && ap.config.Overlay.Image != nil {
		err := compositeCustomImage(dst, "", ap.config.Overlay, true, ap.configDir)
		if err != nil {
			DefaultLogger.Warning("Failed to composite config overlay: %v", err)
		}
	}

	if ogpSettings != nil && ogpSettings.Overlay != nil && ogpSettings.Overlay.Image != nil {
		err := compositeCustomImage(dst, articlePath, ogpSettings.Overlay, false, "")
		if err != nil {
			DefaultLogger.Warning("Failed to composite front matter overlay: %v", err)
		}
	}

	renderOptions := &RenderOptions{
		Font:     font,
		Config:   config,
		Title:    title,
		TestMode: testMode,
	}

	err = ap.imageRenderer.RenderTextOnImage(dst, renderOptions)
	if err != nil {
		return fmt.Errorf("failed to render text: %w", err)
	}

	return ap.saveImage(dst, outputPath)
}

// saveImage writes the generated image to the specified path as PNG.
func (ap *ArticleProcessor) saveImage(img *image.RGBA, outputPath string) error {
	outputFile, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create output file: %w", err)
	}
	defer outputFile.Close()

	err = png.Encode(outputFile, img)
	if err != nil {
		return fmt.Errorf("failed to encode PNG: %w", err)
	}

	return nil
}

// logSuccess outputs information about the successfully generated image.
func (ap *ArticleProcessor) logSuccess(articlePath, outputPath string, testMode bool) {
	relPath, _ := filepath.Rel(ap.contentDir, articlePath)

	if testMode {
		fmt.Printf("Test OGP image generated: %s\n", relPath)
		fmt.Printf("Output file: %s\n", outputPath)

		if stat, err := os.Stat(outputPath); err == nil {
			fmt.Printf("File size: %d bytes\n", stat.Size())
		}
	} else {
		fmt.Printf("Generated OGP image: %s -> %s\n", relPath, outputPath)
	}
}
