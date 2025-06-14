package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// OGPGenerator is the main orchestrator for OGP image generation.
// It manages configuration, font loading, background processing, and text rendering.
type OGPGenerator struct {
	config           *Config
	contentDir       string
	projectRoot      string
	configDir        string
	fontManager      *FontManager
	bgProcessor      *BackgroundProcessor
	imageRenderer    *ImageRenderer
	articleProcessor *ArticleProcessor
}

// NewOGPGenerator creates a new OGPGenerator instance with all necessary components.
// It loads the configuration, initializes processors, and sets up the text rendering pipeline.
func NewOGPGenerator(configPath, contentDir, projectRoot string) (*OGPGenerator, error) {
	config, err := loadConfig(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load config: %w", err)
	}

	configDir := filepath.Dir(configPath)

	fontManager := NewFontManager(configDir)
	bgProcessor := NewBackgroundProcessor(configDir)
	imageRenderer := NewImageRenderer()

	articleProcessor := NewArticleProcessor(config, contentDir, configDir, configPath, fontManager, bgProcessor, imageRenderer)

	return &OGPGenerator{
		config:           config,
		contentDir:       contentDir,
		projectRoot:      projectRoot,
		configDir:        configDir,
		fontManager:      fontManager,
		bgProcessor:      bgProcessor,
		imageRenderer:    imageRenderer,
		articleProcessor: articleProcessor,
	}, nil
}

// TemplateData represents the data available to filename templates.
// It provides access to article metadata for dynamic filename generation.
type TemplateData struct {
	Title       string                 // Article title
	Description string                 // Article description
	Date        interface{}            // Article date (from front matter)
	URL         string                 // Custom URL (if set in front matter)
	RelPath     string                 // Relative path from content directory
	Format      string                 // Output format (png, jpg)
	Fields      map[string]interface{} // All front matter fields
}

// sanitizeFilename removes potentially dangerous characters from filename.
// It prevents path traversal attacks and ensures filesystem compatibility.
func sanitizeFilename(filename string) string {
	// Remove path traversal attempts
	filename = strings.ReplaceAll(filename, "..", "")

	// Replace potentially problematic characters with safe alternatives
	re := regexp.MustCompile(`[<>:"|?*\\/]`)
	filename = re.ReplaceAllString(filename, "_")

	// Remove leading/trailing whitespace and dots
	filename = strings.Trim(filename, " .")

	return filename
}

// generateOutputFilename creates the output filename using template or default logic.
// It supports Go template syntax with access to article metadata and automatically
// appends file extensions based on the output format.
func generateOutputFilename(config *Config, fm *FrontMatter, articlePath, contentDir string) (string, error) {
	// Use template if configured
	if config.Output.Filename != "" {
		templateStr := config.Output.Filename

		// Prepare template data
		relPath, err := filepath.Rel(contentDir, articlePath)
		if err != nil {
			relPath = ""
		}

		data := TemplateData{
			Title:       fm.Title,
			Description: fm.Description,
			Date:        fm.Date,
			URL:         fm.URL,
			RelPath:     relPath,
			Format:      config.Output.Format,
			Fields:      make(map[string]interface{}),
		}

		// Add all front matter fields to Fields map for flexible access
		if fm.Fields != nil {
			data.Fields = fm.Fields
		}

		// Add standard fields to Fields map for template access
		data.Fields["title"] = fm.Title
		data.Fields["description"] = fm.Description
		data.Fields["date"] = fm.Date
		data.Fields["url"] = fm.URL
		data.Fields["tags"] = fm.Tags

		// Use template processor to generate filename
		templateProcessor := NewTemplateProcessor()
		filename, err := templateProcessor.ProcessFilenameTemplate(templateStr, data)
		if err != nil {
			return "", err
		}

		return filename, nil
	}

	// Default behavior: ogp.{format}
	return "ogp." + config.Output.Format, nil
}

// GenerateSingle generates an OGP image for a single article.
// The articlePath can be relative (from contentDir) or absolute.
func (g *OGPGenerator) GenerateSingle(articlePath string) error {
	var fullArticlePath string
	if filepath.IsAbs(articlePath) {
		fullArticlePath = articlePath
	} else {
		fullArticlePath = filepath.Join(g.contentDir, articlePath)
	}

	if _, err := os.Stat(fullArticlePath); os.IsNotExist(err) {
		return NewFileError("stat", fullArticlePath, err)
	}

	indexPath := filepath.Join(fullArticlePath, DefaultIndexFilename)
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return NewFileError("stat", indexPath, err)
	}

	fmt.Printf("Generating OGP image for: %s\n", articlePath)
	return g.articleProcessor.ProcessArticle(fullArticlePath, ProcessOptions{TestMode: false})
}

// GenerateTest generates a test OGP image to a temporary location.
// This is useful for previewing images during development without overwriting production files.
func (g *OGPGenerator) GenerateTest(articlePath string) error {
	// testモードでは記事パスを直接使用（contentDir不要）
	fullArticlePath := articlePath

	if _, err := os.Stat(fullArticlePath); os.IsNotExist(err) {
		return NewFileError("stat", fullArticlePath, err)
	}

	indexPath := filepath.Join(fullArticlePath, DefaultIndexFilename)
	if _, err := os.Stat(indexPath); os.IsNotExist(err) {
		return NewFileError("stat", indexPath, err)
	}

	fmt.Printf("Testing OGP image for: %s\n", articlePath)
	return g.articleProcessor.ProcessArticle(fullArticlePath, ProcessOptions{TestMode: true})
}

// GenerateAll generates OGP images for all articles in the content directory.
// It walks through all directories containing index.md files and processes them.
func (g *OGPGenerator) GenerateAll() error {
	return filepath.WalkDir(g.contentDir, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Name() == DefaultIndexFilename {
			articleDir := filepath.Dir(path)
			return g.articleProcessor.ProcessArticle(articleDir, ProcessOptions{TestMode: false})
		}

		return nil
	})
}
