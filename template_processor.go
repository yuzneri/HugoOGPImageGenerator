package main

import (
	"bytes"
	"regexp"
	"strings"
	tmpl "text/template"
	"time"
	"unicode"
)

// TemplateProcessor handles template processing with Hugo-like functions.
type TemplateProcessor struct {
	funcMap tmpl.FuncMap
}

// NewTemplateProcessor creates a new template processor with default functions.
func NewTemplateProcessor() *TemplateProcessor {
	tp := &TemplateProcessor{}
	tp.initializeFunctions()
	return tp
}

// ProcessTemplate processes a template string with the given data.
func (tp *TemplateProcessor) ProcessTemplate(templateStr string, data interface{}) (string, error) {
	// If no template markers, return as-is
	if !strings.Contains(templateStr, "{{") {
		return templateStr, nil
	}

	t, err := tmpl.New("content").Funcs(tp.funcMap).Parse(templateStr)
	if err != nil {
		return "", NewTemplateError("parse", err)
	}

	var buf bytes.Buffer
	err = t.Execute(&buf, data)
	if err != nil {
		return "", NewTemplateError("execute", err)
	}

	return buf.String(), nil
}

// ProcessFilenameTemplate processes a filename template with the given data.
func (tp *TemplateProcessor) ProcessFilenameTemplate(templateStr string, data TemplateData) (string, error) {
	result, err := tp.ProcessTemplate(templateStr, data)
	if err != nil {
		return "", err
	}

	// Sanitize the result
	filename := sanitizeFilename(result)

	// Ensure we have a valid filename
	if filename == "" {
		return "ogp." + data.Format, nil
	}

	// Auto-append extension if not present
	expectedExt := "." + data.Format
	if !strings.HasSuffix(strings.ToLower(filename), strings.ToLower(expectedExt)) {
		filename = filename + expectedExt
	}

	return filename, nil
}

// ProcessContentTemplate processes a content template with front matter data.
func (tp *TemplateProcessor) ProcessContentTemplate(templateStr string, fm *FrontMatter) (string, error) {
	data := tp.buildTemplateData(fm)
	return tp.ProcessTemplate(templateStr, data)
}

// AddFunction adds a custom function to the template processor.
func (tp *TemplateProcessor) AddFunction(name string, fn interface{}) {
	tp.funcMap[name] = fn
}

// initializeFunctions sets up the default Hugo-like template functions.
func (tp *TemplateProcessor) initializeFunctions() {
	tp.funcMap = tmpl.FuncMap{
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
			return ""
		},
		"now": time.Now,

		// String functions
		"replace": func(old, new, s string) string {
			return strings.ReplaceAll(s, old, new)
		},
		"split": strings.Split,
		"trim":  strings.TrimSpace,

		// URL/Slug functions (Hugo-compatible)
		"slugify": func(s string) string {
			return tp.slugify(s)
		},

		// Math functions
		"add": func(a, b int) int {
			return a + b
		},
		"sub": func(a, b int) int {
			return a - b
		},
		"mul": func(a, b int) int {
			return a * b
		},
		"div": func(a, b int) int {
			if b == 0 {
				return 0
			}
			return a / b
		},

		// Utility functions
		"len": func(s string) int {
			return len(s)
		},
		"slice": func(s string, start, end int) string {
			if start < 0 || end > len(s) || start > end {
				return s
			}
			return s[start:end]
		},
	}
}

// buildTemplateData creates template data from front matter.
func (tp *TemplateProcessor) buildTemplateData(fm *FrontMatter) TemplateData {
	data := TemplateData{
		Title:       fm.Title,
		Description: fm.Description,
		Date:        tp.parseDate(fm.Date),
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
	data.Fields["date"] = tp.parseDate(fm.Date)
	data.Fields["url"] = fm.URL
	data.Fields["tags"] = fm.Tags

	return data
}

// parseDate converts various date formats to time.Time for template usage.
func (tp *TemplateProcessor) parseDate(dateValue interface{}) interface{} {
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

// slugify converts a string to a URL-friendly slug following Hugo's urlize behavior.
// Based on Hugo's urlize function specification:
// - Converts to lowercase
// - Replaces non-alphanumeric characters with hyphens (except dots, hyphens, underscores)
// - Preserves non-ASCII characters (like Japanese, accented characters)
// - Collapses multiple consecutive hyphens into single hyphen
// - Trims leading and trailing hyphens
func (tp *TemplateProcessor) slugify(s string) string {
	if s == "" {
		return ""
	}

	// Convert to lowercase
	result := strings.ToLower(s)

	// Replace problematic characters with hyphens, but preserve:
	// - Alphanumeric characters (including non-ASCII like Japanese)
	// - Dots (.)
	// - Hyphens (-)
	// - Underscores (_)
	var processed strings.Builder
	for _, r := range result {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '.' || r == '-' || r == '_' {
			processed.WriteRune(r)
		} else if unicode.IsSpace(r) || unicode.IsPunct(r) || unicode.IsSymbol(r) {
			processed.WriteRune('-')
		} else {
			processed.WriteRune(r) // Preserve other characters (like CJK)
		}
	}

	result = processed.String()

	// Collapse multiple consecutive hyphens into single hyphen
	re := regexp.MustCompile(`-+`)
	result = re.ReplaceAllString(result, "-")

	// Trim leading and trailing hyphens
	result = strings.Trim(result, "-")

	return result
}
