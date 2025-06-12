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
	templateProcessor *TemplateProcessor
}

// NewArticleProcessor creates a new ArticleProcessor with the given dependencies.
func NewArticleProcessor(config *Config, contentDir, configDir string, fontManager *FontManager, bgProcessor *BackgroundProcessor, imageRenderer *ImageRenderer) *ArticleProcessor {
	return &ArticleProcessor{
		config:            config,
		contentDir:        contentDir,
		fontManager:       fontManager,
		bgProcessor:       bgProcessor,
		imageRenderer:     imageRenderer,
		configDir:         configDir,
		templateProcessor: NewTemplateProcessor(),
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
	indexPath := filepath.Join(articlePath, DefaultIndexFilename)
	content, err := os.ReadFile(indexPath)
	if err != nil {
		return NewFileError("read", indexPath, err)
	}

	fm, err := parseFrontMatter(content)
	if err != nil {
		return NewConfigError(fmt.Sprintf("failed to parse front matter in %s", indexPath), err)
	}

	overrideConfig := applyFrontMatterOverrides(ap.config, fm.OGP)
	title := ap.determineText(fm, &overrideConfig.Title, fm.Title)
	description := ap.determineText(fm, &overrideConfig.Description, fm.Description)

	if title == "" && description == "" {
		return NewValidationError(fmt.Sprintf("no content found in %s (both title and description are empty)", indexPath))
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

	// Output config details in test mode
	if options.TestMode {
		ap.printUsedConfig(overrideConfig, title, description)
		ap.printOutputPaths(overrideConfig, fm, articlePath, options.OutputDir)
	}

	err = ap.generateImage(title, description, outputPath, overrideConfig, articlePath, fm.OGP, options.TestMode)
	if err != nil {
		return NewRenderError("OGP image", err)
	}

	ap.logSuccess(articlePath, outputPath, options.TestMode)
	return nil
}

// determineText resolves the final text to use for image generation.
// It applies the priority: config template > default value.
func (ap *ArticleProcessor) determineText(fm *FrontMatter, textConfig *TextConfig, defaultValue string) string {
	// Priority 1: Config template (if set)
	if textConfig.Content != nil {
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

	// Load font for title (prefer title font, fall back to global font)
	var titleFont string
	if config.Title.Font != nil {
		titleFont = *config.Title.Font
	}
	font, err := ap.fontManager.LoadFont(titleFont, articlePath)
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
	// Determine if article has explicit overlay settings
	hasArticleOverlay := ogpSettings != nil && ogpSettings.Overlay != nil

	if hasArticleOverlay && ogpSettings.Overlay.Image != nil {
		// Article has explicit overlay image - use article-level processing
		if ogpSettings.Overlay.Visible == nil || *ogpSettings.Overlay.Visible {
			err := compositeCustomImage(dst, articlePath, ogpSettings.Overlay, false, ap.configDir)
			if err != nil {
				DefaultLogger.Warning("Failed to composite article overlay: %v", err)
			}
		}
	} else if config.Overlay.Visible && config.Overlay.Image != nil {
		// No article overlay, use config overlay (which may have been influenced by article settings)
		// Determine if this overlay comes from original config or was modified by article settings
		if ap.config.Overlay.Image != nil &&
			config.Overlay.Image != nil && *config.Overlay.Image == *ap.config.Overlay.Image {
			// This is the original config overlay
			err := compositeCustomImage(dst, "", &config.Overlay, true, ap.configDir)
			if err != nil {
				DefaultLogger.Warning("Failed to composite config overlay: %v", err)
			}
		} else {
			// This overlay was modified by article settings - treat as article overlay
			err := compositeCustomImage(dst, articlePath, &config.Overlay, false, ap.configDir)
			if err != nil {
				DefaultLogger.Warning("Failed to composite merged overlay: %v", err)
			}
		}
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

	// Image configuration
	fmt.Println("\nImage:")
	fmt.Printf("  Size: %dx%d\n", DefaultImageWidth, DefaultImageHeight)

	// Output configuration
	fmt.Println("\nOutput:")
	fmt.Printf("  Format: %s\n", config.Output.Format)
	fmt.Printf("  Directory: %s\n", config.Output.Directory)
	fmt.Printf("  Filename Template: %s\n", config.Output.Filename)

	// Background configuration
	fmt.Println("\nBackground:")
	if config.Background.Image != nil {
		fmt.Printf("  Image: %s\n", *config.Background.Image)
	} else {
		fmt.Printf("  Color: %s\n", config.Background.Color)
	}

	// Title configuration
	fmt.Println("\nTitle:")
	fmt.Printf("  Visible: %t\n", config.Title.Visible)
	if config.Title.Visible {

		if config.Title.Content != nil {
			fmt.Printf("  Text: %q\n", *config.Title.Content)
		} else {
			fmt.Printf("  Text: %q\n", title)
		}

		if config.Title.Font != nil {
			fmt.Printf("  Font: %s\n", *config.Title.Font)
		} else {
			fmt.Printf("  Font: (auto-detect)\n")
		}
		fmt.Printf("  Size: %.1f\n", config.Title.Size)
		fmt.Printf("  Color: %s\n", config.Title.Color)
		fmt.Printf("  Block Position: %s\n", config.Title.BlockPosition)
		fmt.Printf("  Line Alignment: %s\n", config.Title.LineAlignment)
		fmt.Printf("  Overflow: %s\n", config.Title.Overflow)
		fmt.Printf("  Min Size: %.1f\n", config.Title.MinSize)
		fmt.Printf("  Line Height: %.2f\n", config.Title.LineHeight)
		fmt.Printf("  Letter Spacing: %d\n", config.Title.LetterSpacing)
		fmt.Printf("  Area: X=%d, Y=%d, Width=%d, Height=%d\n",
			config.Title.Area.X, config.Title.Area.Y,
			config.Title.Area.Width, config.Title.Area.Height)
		fmt.Printf("  Line Breaking:\n")
		fmt.Printf("    Start Prohibited: %q\n", config.Title.LineBreaking.StartProhibited)
		fmt.Printf("    End Prohibited: %q\n", config.Title.LineBreaking.EndProhibited)
	}

	// Description configuration
	fmt.Println("\nDescription:")
	fmt.Printf("  Visible: %t\n", config.Description.Visible)
	if config.Description.Visible {

		if config.Description.Content != nil {
			fmt.Printf("  Text: %q\n", *config.Description.Content)
		} else {
			fmt.Printf("  Text: %q\n", description)
		}
		if config.Description.Font != nil {
			fmt.Printf("  Font: %s\n", *config.Description.Font)
		} else {
			fmt.Printf("  Font: (auto-detect)\n")
		}
		fmt.Printf("  Size: %.1f\n", config.Description.Size)
		fmt.Printf("  Color: %s\n", config.Description.Color)
		fmt.Printf("  Block Position: %s\n", config.Description.BlockPosition)
		fmt.Printf("  Line Alignment: %s\n", config.Description.LineAlignment)
		fmt.Printf("  Overflow: %s\n", config.Description.Overflow)
		fmt.Printf("  Min Size: %.1f\n", config.Description.MinSize)
		fmt.Printf("  Line Height: %.2f\n", config.Description.LineHeight)
		fmt.Printf("  Letter Spacing: %d\n", config.Description.LetterSpacing)
		fmt.Printf("  Area: X=%d, Y=%d, Width=%d, Height=%d\n",
			config.Description.Area.X, config.Description.Area.Y,
			config.Description.Area.Width, config.Description.Area.Height)
		fmt.Printf("  Line Breaking:\n")
		fmt.Printf("    Start Prohibited: %q\n", config.Description.LineBreaking.StartProhibited)
		fmt.Printf("    End Prohibited: %q\n", config.Description.LineBreaking.EndProhibited)
	}

	// Overlay configuration
	fmt.Println("\nOverlay:")
	fmt.Printf("  Visible: %t\n", config.Overlay.Visible)
	if config.Overlay.Visible {
		if config.Overlay.Image != nil {
			fmt.Printf("  Image: %s\n", *config.Overlay.Image)
			fmt.Printf("  Placement:\n")
			fmt.Printf("    X: %d\n", config.Overlay.Placement.X)
			fmt.Printf("    Y: %d\n", config.Overlay.Placement.Y)

			if config.Overlay.Placement.Width != nil {
				fmt.Printf("    Width: %d\n", *config.Overlay.Placement.Width)
			} else {
				fmt.Printf("    Width: (auto-detect)\n")
			}

			if config.Overlay.Placement.Height != nil {
				fmt.Printf("    Height: %d\n", *config.Overlay.Placement.Height)
			} else {
				fmt.Printf("    Height: (auto-detect)\n")
			}

			fmt.Printf("  Fit: %s\n", config.Overlay.Fit)
			fmt.Printf("  Opacity: %.2f\n", config.Overlay.Opacity)
		} else {
			fmt.Printf("  Image: (none)\n")
		}
	}

	fmt.Println("\n=== End Configuration ===")
	fmt.Println()
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
