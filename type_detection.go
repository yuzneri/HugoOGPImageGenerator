package main

import (
	"fmt"
	"path/filepath"
	"strings"
)

// Hugo content type constants
const (
	DefaultContentType = "page"
	ContentDirName     = "content"
)

// determineContentType determines the Hugo content type for an article.
// Priority order: front matter type field > directory-based detection > default "page"
func determineContentType(frontMatter *FrontMatter, articlePath, hugoRootPath string) string {
	// 1. Check front matter type field first (highest priority)
	if frontMatter != nil && frontMatter.Type != "" {
		return frontMatter.Type
	}

	// 2. Convert paths to absolute if they are relative
	absArticlePath := articlePath
	if !filepath.IsAbs(articlePath) {
		if abs, err := filepath.Abs(articlePath); err == nil {
			absArticlePath = abs
		}
	}

	absHugoRootPath := hugoRootPath
	if !filepath.IsAbs(hugoRootPath) {
		if abs, err := filepath.Abs(hugoRootPath); err == nil {
			absHugoRootPath = abs
		}
	}

	// 3. Fall back to directory-based type detection
	return extractDirectoryType(absArticlePath, absHugoRootPath)
}

// extractDirectoryType extracts the content type from the directory structure.
// According to Hugo's convention, the first directory under content/ determines the type.
// If the article is directly under content/, the type defaults to DefaultContentType.
func extractDirectoryType(articlePath, hugoRootPath string) string {
	relPath, err := getRelativePathFromContent(articlePath, hugoRootPath)
	if err != nil {
		return DefaultContentType
	}

	contentType := extractTypeFromPath(relPath)
	return contentType
}

// getRelativePathFromContent calculates the relative path from the content directory.
func getRelativePathFromContent(articlePath, hugoRootPath string) (string, error) {
	cleanArticlePath := filepath.Clean(articlePath)
	cleanHugoRoot := filepath.Clean(hugoRootPath)

	// Validate that the article path is absolute
	if !filepath.IsAbs(cleanArticlePath) {
		return "", NewValidationError(fmt.Sprintf("article path must be absolute: %s", cleanArticlePath))
	}

	contentDir := filepath.Join(cleanHugoRoot, ContentDirName)
	relPath, err := filepath.Rel(contentDir, cleanArticlePath)
	if err != nil {
		return "", NewFileError("calculate relative path", contentDir, err)
	}

	return relPath, nil
}

// extractTypeFromPath extracts the content type from a relative path.
func extractTypeFromPath(relPath string) string {
	// Handle special cases
	if relPath == "." || relPath == "" {
		return DefaultContentType
	}

	components := parsePathComponents(relPath)
	if len(components) == 0 {
		return DefaultContentType
	}

	firstComponent := components[0]
	if isFileComponent(firstComponent) {
		return DefaultContentType
	}

	return firstComponent
}

// parsePathComponents splits a path into valid components, filtering out empty and invalid parts.
func parsePathComponents(relPath string) []string {
	// Normalize path separators
	normalizedPath := filepath.ToSlash(relPath)

	// Split into components
	pathComponents := strings.Split(normalizedPath, "/")

	// Filter valid components
	var validComponents []string
	for _, component := range pathComponents {
		if isValidPathComponent(component) {
			validComponents = append(validComponents, component)
		}
	}

	return validComponents
}

// isValidPathComponent checks if a path component is valid (not empty, not current/parent dir).
func isValidPathComponent(component string) bool {
	return component != "" && component != "." && component != ".."
}

// isFileComponent checks if a component represents a file (contains an extension).
func isFileComponent(component string) bool {
	return strings.Contains(component, ".")
}
