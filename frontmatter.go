package main

import (
	"bytes"
	"fmt"

	"gopkg.in/yaml.v3"
)

// FrontMatter represents the YAML front matter structure commonly used in static site generators.
// It contains basic article metadata and optional OGP-specific settings.
type FrontMatter struct {
	Title       string                 `yaml:"title"`         // Article title
	Date        interface{}            `yaml:"date"`          // Publication date (can be string or time.Time)
	Description string                 `yaml:"description"`   // Article description
	Tags        []string               `yaml:"tags"`          // Article tags
	Type        string                 `yaml:"type"`          // Hugo content type
	URL         string                 `yaml:"url"`           // Custom URL (overrides default)
	OGP         *OGPFrontMatter        `yaml:"ogp,omitempty"` // OGP-specific settings
	Fields      map[string]interface{} `yaml:",inline"`       // Additional fields for template access
}

// parseFrontMatter extracts and parses YAML front matter from article content.
// It expects front matter to be delimited by "---" lines at the beginning of the content.
func parseFrontMatter(content []byte) (*FrontMatter, error) {
	if !bytes.HasPrefix(content, []byte("---\n")) {
		return nil, fmt.Errorf("no front matter found")
	}

	// Skip the opening "---\n"
	content = content[4:]

	// Find the closing "---" delimiter
	endIndex := bytes.Index(content, []byte("\n---\n"))
	if endIndex == -1 {
		return nil, fmt.Errorf("front matter end delimiter not found")
	}

	// Extract just the front matter content
	frontMatterContent := content[:endIndex]

	// Parse the YAML content
	var fm FrontMatter
	err := yaml.Unmarshal(frontMatterContent, &fm)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal front matter: %w", err)
	}

	return &fm, nil
}
