package main

import (
	"fmt"
	"image"
	"image/draw"
	"image/png"
	"os"
	"path/filepath"

	"github.com/golang/freetype/truetype"
)

// ArticleProcessor handles the complete processing pipeline for individual articles.
// It coordinates front matter parsing, configuration merging, and image generation.
type ArticleProcessor struct {
	config            *Config
	contentDir        string
	fontManager       *FontManager
	bgProcessor       *BackgroundProcessor
	imageRenderer     *ImageRenderer
	configDir         string
	configPath        string
	templateProcessor *TemplateProcessor
	configMerger      *ConfigMerger
}

// NewArticleProcessor creates a new ArticleProcessor with the given dependencies.
func NewArticleProcessor(config *Config, contentDir, configDir, configPath string, fontManager *FontManager, bgProcessor *BackgroundProcessor, imageRenderer *ImageRenderer) *ArticleProcessor {
	return &ArticleProcessor{
		config:            config,
		contentDir:        contentDir,
		fontManager:       fontManager,
		bgProcessor:       bgProcessor,
		imageRenderer:     imageRenderer,
		configDir:         configDir,
		configPath:        configPath,
		templateProcessor: NewTemplateProcessor(),
		configMerger:      NewConfigMerger(),
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
	// Parse front matter and build configuration
	fm, finalConfig, err := ap.parseAndConfigureArticle(articlePath)
	if err != nil {
		return err
	}

	// Determine text content for image generation
	title, description, err := ap.determineArticleContent(fm, finalConfig)
	if err != nil {
		return err
	}

	// Generate appropriate output path
	outputPath, err := ap.generateOutputPath(finalConfig, fm, articlePath, options)
	if err != nil {
		return err
	}

	// Handle test mode output
	if options.TestMode {
		ap.handleTestModeOutput(finalConfig, fm, articlePath, title, description, options.OutputDir)
	}

	// Generate the OGP image
	err = ap.generateImage(title, description, outputPath, finalConfig, articlePath, fm.OGP, options.TestMode)
	if err != nil {
		return NewRenderError("OGP image", err)
	}

	ap.logSuccess(articlePath, outputPath, options.TestMode)
	return nil
}

// buildFinalConfigurationWithSettings creates the final configuration by applying the 4-level hierarchy:
// Default Config -> Global ConfigSettings -> Type ConfigSettings -> Front Matter Overrides
func (ap *ArticleProcessor) buildFinalConfigurationWithSettings(fm *FrontMatter, articlePath string) (*Config, error) {
	// Determine content type and load configurations
	_, typeSettings, err := ap.loadContentTypeConfigurationSettings(fm, articlePath)
	if err != nil {
		return nil, err
	}

	// Load global settings from the config file
	globalSettings, err := loadConfigSettings(ap.getConfigPath())
	if err != nil {
		return nil, fmt.Errorf("failed to load global config settings: %w", err)
	}

	// Apply 4-level configuration hierarchy using settings
	finalConfig := ap.configMerger.MergeConfigsWithSettings(
		getDefaultConfig(),
		globalSettings,
		typeSettings,
		fm.OGP,
	)

	return finalConfig, nil
}

// loadContentTypeConfigurationSettings determines content type and loads type-specific configuration as settings
func (ap *ArticleProcessor) loadContentTypeConfigurationSettings(fm *FrontMatter, articlePath string) (string, *ConfigSettings, error) {
	contentType := determineContentType(fm, articlePath, ap.getHugoRootPath())

	typeSettings, err := loadTypeConfigSettings(ap.configDir, contentType)
	if err != nil {
		return "", nil, fmt.Errorf("failed to load type config settings for type '%s' from article '%s': %w",
			contentType, filepath.Base(articlePath), err)
	}

	return contentType, typeSettings, nil
}

// getConfigPath returns the path to the configuration file
func (ap *ArticleProcessor) getConfigPath() string {
	return ap.configPath
}

// getHugoRootPath determines the Hugo root path from the content directory
func (ap *ArticleProcessor) getHugoRootPath() string {
	// The content directory is typically <hugo_root>/content
	// So the Hugo root is the parent of the content directory
	contentDir := ap.contentDir

	// Convert to absolute path if it's relative
	if !filepath.IsAbs(contentDir) {
		absContentDir, err := filepath.Abs(contentDir)
		if err == nil {
			contentDir = absContentDir
		}
	}

	return filepath.Dir(contentDir)
}

// determineText resolves the final text to use for image generation.
// It applies the priority: config template > default value.
func (ap *ArticleProcessor) determineText(fm *FrontMatter, textConfig *TextConfig, defaultValue string) string {
	// Priority 1: Config template (if set)
	if textConfig.Content != nil && *textConfig.Content != "" {
		return ap.processContentTemplate(*textConfig.Content, fm)
	}

	// Priority 2: Default value (article title or description)
	return defaultValue
}

// processContentTemplate processes content template with front matter data.
func (ap *ArticleProcessor) processContentTemplate(contentTemplate string, fm *FrontMatter) string {
	result, err := ap.templateProcessor.ProcessContentTemplate(contentTemplate, fm)
	if err != nil {
		DefaultLogger.Warning("Failed to process content template: %v", err)
		return contentTemplate
	}
	return result
}

// generateTestOutputPath creates a temporary output path for test mode.
func (ap *ArticleProcessor) generateTestOutputPath() string {
	execPath, execErr := os.Executable()
	if execErr != nil {
		execDir, _ := os.Getwd()
		return filepath.Join(execDir, DefaultTestFilename)
	}

	execDir := filepath.Dir(execPath)
	return filepath.Join(execDir, "test.png")
}

// calculateOutputPath calculates the output path without creating directories.
// This is a pure function that can be used for both production and display purposes.
func calculateOutputPath(config *Config, fm *FrontMatter, articlePath, contentDir, outputDir string) (string, error) {
	filename, err := generateOutputFilename(config, fm, articlePath, contentDir)
	if err != nil {
		return "", fmt.Errorf("failed to generate output filename: %w", err)
	}

	if outputDir == "" {
		outputDir = filepath.Join(filepath.Dir(contentDir), config.Output.Directory)
	}

	var outputPath string
	if fm.URL != "" {
		outputPath = filepath.Join(outputDir, fm.URL, filename)
	} else {
		relPath, err := filepath.Rel(contentDir, articlePath)
		if err != nil {
			return "", fmt.Errorf("failed to get relative path: %w", err)
		}
		outputPath = filepath.Join(outputDir, relPath, filename)
	}

	return outputPath, nil
}

// generateProductionOutputPath creates the final output path for production mode.
// It respects custom URLs and creates the necessary directory structure.
func (ap *ArticleProcessor) generateProductionOutputPath(config *Config, fm *FrontMatter, articlePath, outputDir string) (string, error) {
	outputPath, err := calculateOutputPath(config, fm, articlePath, ap.contentDir, outputDir)
	if err != nil {
		return "", err
	}

	err = os.MkdirAll(filepath.Dir(outputPath), DefaultFilePermission)
	if err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	return outputPath, nil
}

// generateImage creates the OGP image by compositing background, overlays, and text.
// It handles both config-level and article-level overlay compositions.
func (ap *ArticleProcessor) generateImage(title, description, outputPath string, config *Config, articlePath string, ogpSettings *OGPFrontMatter, testMode bool) error {
	dst, font, err := ap.setupImageCanvas(config, articlePath)
	if err != nil {
		return err
	}

	err = ap.applyOverlays(dst, config, articlePath, ogpSettings)
	if err != nil {
		return err
	}

	err = ap.renderTextElements(dst, font, config, title, description, testMode)
	if err != nil {
		return err
	}

	return ap.saveImage(dst, outputPath)
}

// setupImageCanvas creates the base image canvas with background.
func (ap *ArticleProcessor) setupImageCanvas(config *Config, articlePath string) (*image.RGBA, *truetype.Font, error) {
	backgroundImage, err := ap.bgProcessor.CreateBackground(config, articlePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create background: %w", err)
	}

	// Load font for title (nil means auto-detect)
	titleFontPath := ""
	if config.Title.Font != nil {
		titleFontPath = *config.Title.Font
	}
	font, err := ap.fontManager.LoadFont(titleFontPath, articlePath)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load font: %w", err)
	}

	// Note: Currently using the same font for both title and description
	// In the future, we could support different fonts for each text element

	bounds := backgroundImage.Bounds()
	dst := image.NewRGBA(bounds)
	draw.Draw(dst, bounds, backgroundImage, image.Point{}, draw.Src)

	return dst, font, nil
}

// applyOverlays applies both config-level and article-level overlays to the image.
func (ap *ArticleProcessor) applyOverlays(dst *image.RGBA, config *Config, articlePath string, ogpSettings *OGPFrontMatter) error {
	// Check if overlay should be rendered
	if !config.Overlay.Visible || config.Overlay.Image == nil || *config.Overlay.Image == "" {
		return nil
	}

	// The config already contains merged overlay settings from all sources
	// (defaults -> global -> type -> front matter), so just use the final config
	err := compositeCustomImage(dst, articlePath, &config.Overlay, false, ap.configDir)
	if err != nil {
		DefaultLogger.Warning("Failed to composite overlay: %v", err)
	}

	return nil
}

// renderTextElements renders title and description text onto the image.
func (ap *ArticleProcessor) renderTextElements(dst *image.RGBA, font *truetype.Font, config *Config, title, description string, testMode bool) error {
	renderOptions := &RenderOptions{
		Font:        font,
		Config:      config,
		Title:       title,
		Description: description,
		TestMode:    testMode,
	}

	err := ap.imageRenderer.RenderTextOnImage(dst, renderOptions)
	if err != nil {
		return fmt.Errorf("failed to render text: %w", err)
	}

	return nil
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

// parseAndConfigureArticle reads front matter and builds the final configuration
func (ap *ArticleProcessor) parseAndConfigureArticle(articlePath string) (*FrontMatter, *Config, error) {
	indexPath := filepath.Join(articlePath, DefaultIndexFilename)
	content, err := os.ReadFile(indexPath)
	if err != nil {
		return nil, nil, NewFileError("read", indexPath, err)
	}

	fm, err := parseFrontMatter(content)
	if err != nil {
		return nil, nil, NewConfigError(fmt.Sprintf("failed to parse front matter in %s", indexPath), err)
	}

	// Apply 4-level configuration hierarchy: Default -> Global -> Type -> Front Matter
	finalConfig, err := ap.buildFinalConfigurationWithSettings(fm, articlePath)
	if err != nil {
		return nil, nil, NewConfigError(fmt.Sprintf("failed to build configuration for %s", indexPath), err)
	}

	return fm, finalConfig, nil
}

// determineArticleContent resolves title and description text and validates content exists
func (ap *ArticleProcessor) determineArticleContent(fm *FrontMatter, config *Config) (string, string, error) {
	title := ap.determineText(fm, &config.Title, fm.Title)
	description := ap.determineText(fm, &config.Description, fm.Description)

	if title == "" && description == "" {
		indexPath := filepath.Join("(article)", DefaultIndexFilename)
		return "", "", NewValidationError(fmt.Sprintf("no content found in %s (both title and description are empty)", indexPath))
	}

	return title, description, nil
}

// generateOutputPath creates the appropriate output path based on processing options
func (ap *ArticleProcessor) generateOutputPath(config *Config, fm *FrontMatter, articlePath string, options ProcessOptions) (string, error) {
	if options.TestMode {
		return ap.generateTestOutputPath(), nil
	}

	return ap.generateProductionOutputPath(config, fm, articlePath, options.OutputDir)
}

// handleTestModeOutput prints configuration and path information in test mode
func (ap *ArticleProcessor) handleTestModeOutput(config *Config, fm *FrontMatter, articlePath, title, description, outputDir string) {
	ap.printUsedConfig(config, title, description)
	ap.printOutputPaths(config, fm, articlePath, outputDir)
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

// printUsedConfig prints the configuration used for OGP generation in test mode.
func (ap *ArticleProcessor) printUsedConfig(config *Config, title, description string) {
	fmt.Println("\n=== Configuration Used for OGP Generation ===")

	ap.printImageConfig()
	ap.printOutputConfig(config)
	ap.printBackgroundConfig(config)
	ap.printTitleConfig(config, title)
	ap.printDescriptionConfig(config, description)
	ap.printOverlayConfig(config)

	fmt.Println("\n=== End Configuration ===")
	fmt.Println()
}

// printImageConfig prints image dimension configuration
func (ap *ArticleProcessor) printImageConfig() {
	fmt.Println("\nImage:")
	fmt.Printf("  Size: %dx%d\n", DefaultImageWidth, DefaultImageHeight)
}

// printOutputConfig prints output configuration details
func (ap *ArticleProcessor) printOutputConfig(config *Config) {
	fmt.Println("\nOutput:")
	fmt.Printf("  Format: %s\n", config.Output.Format)
	fmt.Printf("  Directory: %s\n", config.Output.Directory)
	fmt.Printf("  Filename Template: %s\n", config.Output.Filename)
}

// printBackgroundConfig prints background configuration details
func (ap *ArticleProcessor) printBackgroundConfig(config *Config) {
	fmt.Println("\nBackground:")
	if config.Background.Image != nil && *config.Background.Image != "" {
		fmt.Printf("  Image: %s\n", *config.Background.Image)
	} else {
		fmt.Printf("  Color: %s\n", config.Background.Color)
	}
}

// printTitleConfig prints title configuration details
func (ap *ArticleProcessor) printTitleConfig(config *Config, title string) {
	fmt.Println("\nTitle:")
	ap.printTextConfigDetails(&config.Title, title)
}

// printDescriptionConfig prints description configuration details
func (ap *ArticleProcessor) printDescriptionConfig(config *Config, description string) {
	fmt.Println("\nDescription:")
	ap.printTextConfigDetails(&config.Description, description)
}

// printTextConfigDetails prints common text configuration details (shared by title and description)
func (ap *ArticleProcessor) printTextConfigDetails(textConfig *TextConfig, defaultText string) {
	fmt.Printf("  Visible: %t\n", textConfig.Visible)
	if !textConfig.Visible {
		return
	}

	// Print text content
	if textConfig.Content != nil && *textConfig.Content != "" {
		fmt.Printf("  Text: %q\n", *textConfig.Content)
	} else {
		fmt.Printf("  Text: %q\n", defaultText)
	}

	// Print font configuration
	if textConfig.Font != nil && *textConfig.Font != "" {
		fmt.Printf("  Font: %s\n", *textConfig.Font)
	} else {
		fmt.Printf("  Font: (auto-detect)\n")
	}

	// Print text styling configuration
	fmt.Printf("  Size: %.1f\n", textConfig.Size)
	fmt.Printf("  Color: %s\n", textConfig.Color)
	fmt.Printf("  Block Position: %s\n", textConfig.BlockPosition)
	fmt.Printf("  Line Alignment: %s\n", textConfig.LineAlignment)
	fmt.Printf("  Overflow: %s\n", textConfig.Overflow)
	fmt.Printf("  Min Size: %.1f\n", textConfig.MinSize)
	fmt.Printf("  Line Height: %.2f\n", textConfig.LineHeight)
	fmt.Printf("  Letter Spacing: %d\n", textConfig.LetterSpacing)

	// Print area configuration
	fmt.Printf("  Area: X=%d, Y=%d, Width=%d, Height=%d\n",
		textConfig.Area.X, textConfig.Area.Y,
		textConfig.Area.Width, textConfig.Area.Height)

	// Print line breaking configuration
	fmt.Printf("  Line Breaking:\n")
	fmt.Printf("    Start Prohibited: %q\n", textConfig.LineBreaking.StartProhibited)
	fmt.Printf("    End Prohibited: %q\n", textConfig.LineBreaking.EndProhibited)
}

// printOverlayConfig prints overlay configuration details
func (ap *ArticleProcessor) printOverlayConfig(config *Config) {
	fmt.Println("\nOverlay:")
	fmt.Printf("  Visible: %t\n", config.Overlay.Visible)

	if !config.Overlay.Visible {
		return
	}

	if config.Overlay.Image != nil && *config.Overlay.Image != "" {
		fmt.Printf("  Image: %s\n", *config.Overlay.Image)
		ap.printOverlayPlacement(config.Overlay.Placement)
		fmt.Printf("  Fit: %s\n", config.Overlay.Fit)
		fmt.Printf("  Opacity: %.2f\n", config.Overlay.Opacity)
	} else {
		fmt.Printf("  Image: (none)\n")
	}
}

// printOverlayPlacement prints overlay placement configuration details
func (ap *ArticleProcessor) printOverlayPlacement(placement PlacementConfig) {
	fmt.Printf("  Placement:\n")
	fmt.Printf("    X: %d\n", placement.X)
	fmt.Printf("    Y: %d\n", placement.Y)

	if placement.Width != nil {
		fmt.Printf("    Width: %d\n", *placement.Width)
	} else {
		fmt.Printf("    Width: (auto-detect)\n")
	}

	if placement.Height != nil {
		fmt.Printf("    Height: %d\n", *placement.Height)
	} else {
		fmt.Printf("    Height: (auto-detect)\n")
	}
}

// printOutputPaths displays the expected output paths in test mode.
func (ap *ArticleProcessor) printOutputPaths(config *Config, fm *FrontMatter, articlePath, outputDir string) {
	contentDir := ap.contentDir

	outputPath, err := calculateOutputPath(config, fm, articlePath, contentDir, outputDir)
	if err != nil {
		fmt.Printf("Error calculating output path: %v\n", err)
		return
	}

	// Display the output path as relative path from current directory
	cwd, _ := os.Getwd()
	relOutputPath, err := filepath.Rel(cwd, outputPath)
	if err != nil {
		// If cannot make relative, use the absolute path
		relOutputPath = outputPath
	}
	fmt.Printf("Output OGP image Path: %s\n", relOutputPath)
}
