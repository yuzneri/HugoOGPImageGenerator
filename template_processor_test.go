package main

import (
	"errors"
	"strings"
	"testing"
	"time"
)

func TestNewTemplateProcessor(t *testing.T) {
	tp := NewTemplateProcessor()
	if tp == nil {
		t.Error("Expected NewTemplateProcessor to return a non-nil processor")
	}

	if tp.funcMap == nil {
		t.Error("Expected funcMap to be initialized")
	}

	// Test that default functions are available
	expectedFuncs := []string{"default", "upper", "lower", "title", "dateFormat", "now", "replace", "split", "trim", "slugify", "add", "sub", "mul", "div", "len", "slice"}
	for _, funcName := range expectedFuncs {
		if tp.funcMap[funcName] == nil {
			t.Errorf("Expected function %s to be available", funcName)
		}
	}
}

func TestTemplateProcessor_ProcessTemplate_NoTemplateMarkers(t *testing.T) {
	tp := NewTemplateProcessor()
	input := "This is plain text without template markers"

	result, err := tp.ProcessTemplate(input, nil)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != input {
		t.Errorf("Expected %q, got %q", input, result)
	}
}

func TestTemplateProcessor_ProcessTemplate_BasicTemplate(t *testing.T) {
	tp := NewTemplateProcessor()

	data := struct {
		Name string
		Age  int
	}{
		Name: "Alice",
		Age:  30,
	}

	template := "Hello {{.Name}}, you are {{.Age}} years old"
	expected := "Hello Alice, you are 30 years old"

	result, err := tp.ProcessTemplate(template, data)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestTemplateProcessor_ProcessTemplate_WithFunctions(t *testing.T) {
	tp := NewTemplateProcessor()

	data := struct {
		Name string
	}{
		Name: "alice",
	}

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "upper function",
			template: "Hello {{upper .Name}}",
			expected: "Hello ALICE",
		},
		{
			name:     "title function",
			template: "Hello {{title .Name}}",
			expected: "Hello Alice",
		},
		{
			name:     "default function with value",
			template: "Hello {{default \"World\" .Name}}",
			expected: "Hello alice",
		},
		{
			name:     "replace function",
			template: "{{replace \"a\" \"@\" .Name}}",
			expected: "@lice",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tp.ProcessTemplate(tt.template, data)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTemplateProcessor_ProcessTemplate_DefaultFunction(t *testing.T) {
	tp := NewTemplateProcessor()

	data := struct {
		Name  string
		Empty string
	}{
		Name:  "Alice",
		Empty: "",
	}

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "default with non-empty value",
			template: "{{default \"World\" .Name}}",
			expected: "Alice",
		},
		{
			name:     "default with empty value",
			template: "{{default \"World\" .Empty}}",
			expected: "World",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tp.ProcessTemplate(tt.template, data)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTemplateProcessor_ProcessTemplate_MathFunctions(t *testing.T) {
	tp := NewTemplateProcessor()

	data := struct {
		A int
		B int
	}{
		A: 10,
		B: 3,
	}

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "add function",
			template: "{{add .A .B}}",
			expected: "13",
		},
		{
			name:     "sub function",
			template: "{{sub .A .B}}",
			expected: "7",
		},
		{
			name:     "mul function",
			template: "{{mul .A .B}}",
			expected: "30",
		},
		{
			name:     "div function",
			template: "{{div .A .B}}",
			expected: "3",
		},
		{
			name:     "div by zero",
			template: "{{div .A 0}}",
			expected: "0",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tp.ProcessTemplate(tt.template, data)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTemplateProcessor_ProcessTemplate_StringFunctions(t *testing.T) {
	tp := NewTemplateProcessor()

	data := struct {
		Text string
	}{
		Text: "  hello world  ",
	}

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "len function",
			template: "{{len .Text}}",
			expected: "15",
		},
		{
			name:     "trim function",
			template: "{{trim .Text}}",
			expected: "hello world",
		},
		{
			name:     "slice function",
			template: "{{slice .Text 2 7}}",
			expected: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tp.ProcessTemplate(tt.template, data)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTemplateProcessor_ProcessTemplate_DateFunctions(t *testing.T) {
	tp := NewTemplateProcessor()

	testTime := time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC)
	data := struct {
		Date time.Time
	}{
		Date: testTime,
	}

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "dateFormat with time.Time",
			template: "{{dateFormat \"2006-01-02\" .Date}}",
			expected: "2023-12-25",
		},
		{
			name:     "dateFormat with RFC3339",
			template: "{{dateFormat \"15:04\" .Date}}",
			expected: "15:30",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tp.ProcessTemplate(tt.template, data)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTemplateProcessor_ProcessTemplate_InvalidTemplate(t *testing.T) {
	tp := NewTemplateProcessor()

	invalidTemplate := "{{.Name"

	_, err := tp.ProcessTemplate(invalidTemplate, nil)
	if err == nil {
		t.Error("Expected error for invalid template")
	}

	var templateErr *AppError
	if !errors.As(err, &templateErr) {
		t.Error("Expected TemplateError type")
	}

	if templateErr.Type != TemplateError {
		t.Errorf("Expected TemplateError type, got %v", templateErr.Type)
	}
}

func TestTemplateProcessor_ProcessFilenameTemplate(t *testing.T) {
	tp := NewTemplateProcessor()

	data := TemplateData{
		Title:   "Test Article",
		Format:  "png",
		RelPath: "articles/test",
	}

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "simple filename template",
			template: "{{.Title}}.{{.Format}}",
			expected: "Test Article.png",
		},
		{
			name:     "filename with relpath sanitized",
			template: "{{.RelPath}}-{{.Title}}.{{.Format}}",
			expected: "articles_test-Test Article.png",
		},
		{
			name:     "empty template uses default",
			template: "",
			expected: "ogp.png",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tp.ProcessFilenameTemplate(tt.template, data)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTemplateProcessor_ProcessFilenameTemplate_SanitizedOutput(t *testing.T) {
	tp := NewTemplateProcessor()

	data := TemplateData{
		Title:  "Test/Article<>",
		Format: "png",
	}

	template := "{{.Title}}.{{.Format}}"
	result, err := tp.ProcessFilenameTemplate(template, data)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if strings.Contains(result, "/") || strings.Contains(result, "<") || strings.Contains(result, ">") {
		t.Errorf("Expected sanitized filename, got %q", result)
	}

	expected := "Test_Article__.png"
	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestTemplateProcessor_ProcessFilenameTemplate_AutoExtension(t *testing.T) {
	tp := NewTemplateProcessor()

	data := TemplateData{
		Title:  "Test",
		Format: "jpg",
	}

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "template without extension",
			template: "{{.Title}}",
			expected: "Test.jpg",
		},
		{
			name:     "template with correct extension",
			template: "{{.Title}}.jpg",
			expected: "Test.jpg",
		},
		{
			name:     "template with wrong extension",
			template: "{{.Title}}.png",
			expected: "Test.png.jpg",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tp.ProcessFilenameTemplate(tt.template, data)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTemplateProcessor_ProcessContentTemplate(t *testing.T) {
	tp := NewTemplateProcessor()

	fm := &FrontMatter{
		Title:       "Test Article",
		Description: "This is a test",
		Date:        time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC),
		URL:         "/test",
		Tags:        []string{"test", "example"},
		Fields: map[string]interface{}{
			"category": "tutorial",
			"author":   "John Doe",
		},
	}

	tests := []struct {
		name     string
		template string
		expected string
	}{
		{
			name:     "simple content template",
			template: "{{.Title}} - {{.Description}}",
			expected: "Test Article - This is a test",
		},
		{
			name:     "template with fields access",
			template: "{{.Fields.category}} by {{.Fields.author}}",
			expected: "tutorial by John Doe",
		},
		{
			name:     "template with date formatting",
			template: "{{dateFormat \"2006-01-02\" .Date}}",
			expected: "2023-12-25",
		},
		{
			name:     "template with default function",
			template: "{{default \"No URL\" .URL}}",
			expected: "/test",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := tp.ProcessContentTemplate(tt.template, fm)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if result != tt.expected {
				t.Errorf("Expected %q, got %q", tt.expected, result)
			}
		})
	}
}

func TestTemplateProcessor_AddFunction(t *testing.T) {
	tp := NewTemplateProcessor()

	// Add a custom function
	tp.AddFunction("custom", func(s string) string {
		return "custom-" + s
	})

	data := struct {
		Name string
	}{
		Name: "test",
	}

	template := "{{custom .Name}}"
	expected := "custom-test"

	result, err := tp.ProcessTemplate(template, data)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}

	if result != expected {
		t.Errorf("Expected %q, got %q", expected, result)
	}
}

func TestTemplateProcessor_parseDate(t *testing.T) {
	tp := NewTemplateProcessor()

	tests := []struct {
		name     string
		input    interface{}
		expected interface{}
	}{
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
		{
			name:     "time.Time input",
			input:    time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC),
			expected: time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "RFC3339 string",
			input:    "2023-12-25T15:30:45Z",
			expected: time.Date(2023, 12, 25, 15, 30, 45, 0, time.UTC),
		},
		{
			name:     "date only string",
			input:    "2023-12-25",
			expected: time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC),
		},
		{
			name:     "unparseable string",
			input:    "invalid-date",
			expected: "invalid-date",
		},
		{
			name:     "integer input",
			input:    123456,
			expected: 123456,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tp.parseDate(tt.input)

			if result != tt.expected {
				// For time.Time comparison, we need special handling
				if expectedTime, ok := tt.expected.(time.Time); ok {
					if resultTime, ok := result.(time.Time); ok {
						if !expectedTime.Equal(resultTime) {
							t.Errorf("Expected %v, got %v", tt.expected, result)
						}
					} else {
						t.Errorf("Expected time.Time, got %T", result)
					}
				} else {
					t.Errorf("Expected %v, got %v", tt.expected, result)
				}
			}
		})
	}
}

func TestTemplateProcessor_buildTemplateData(t *testing.T) {
	tp := NewTemplateProcessor()

	testTime := time.Date(2023, 12, 25, 0, 0, 0, 0, time.UTC)
	fm := &FrontMatter{
		Title:       "Test Article",
		Description: "This is a test",
		Date:        testTime,
		URL:         "/test",
		Tags:        []string{"test", "example"},
		Fields: map[string]interface{}{
			"category": "tutorial",
			"author":   "John Doe",
		},
	}

	data := tp.buildTemplateData(fm)

	if data.Title != "Test Article" {
		t.Errorf("Expected title 'Test Article', got %q", data.Title)
	}

	if data.Description != "This is a test" {
		t.Errorf("Expected description 'This is a test', got %q", data.Description)
	}

	if data.URL != "/test" {
		t.Errorf("Expected URL '/test', got %q", data.URL)
	}

	if dateValue, ok := data.Date.(time.Time); ok {
		if !dateValue.Equal(testTime) {
			t.Errorf("Expected date %v, got %v", testTime, dateValue)
		}
	} else {
		t.Errorf("Expected date to be time.Time, got %T", data.Date)
	}

	// Check fields map
	if data.Fields["category"] != "tutorial" {
		t.Errorf("Expected category 'tutorial', got %v", data.Fields["category"])
	}

	if data.Fields["author"] != "John Doe" {
		t.Errorf("Expected author 'John Doe', got %v", data.Fields["author"])
	}

	// Check standard fields are also in Fields map
	if data.Fields["title"] != "Test Article" {
		t.Errorf("Expected title in fields 'Test Article', got %v", data.Fields["title"])
	}

	if data.Fields["description"] != "This is a test" {
		t.Errorf("Expected description in fields 'This is a test', got %v", data.Fields["description"])
	}
}

func TestTemplateProcessor_SlugifyFunction(t *testing.T) {
	tp := NewTemplateProcessor()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic text",
			input:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "with special characters",
			input:    "Hello, World!",
			expected: "hello-world",
		},
		{
			name:     "with spaces and punctuation",
			input:    "< a, b, & c >",
			expected: "a-b-c",
		},
		{
			name:     "mixed case with numbers",
			input:    "Article 123 Title",
			expected: "article-123-title",
		},
		{
			name:     "with dots (preserved)",
			input:    "main.go",
			expected: "main.go",
		},
		{
			name:     "with hyphens and underscores (preserved)",
			input:    "my-file_name",
			expected: "my-file_name",
		},
		{
			name:     "with consecutive spaces",
			input:    "a  b   c",
			expected: "a-b-c",
		},
		{
			name:     "with leading/trailing spaces",
			input:    "  hello world  ",
			expected: "hello-world",
		},
		{
			name:     "Japanese text (preserved)",
			input:    "こんにちは世界",
			expected: "こんにちは世界",
		},
		{
			name:     "mixed Japanese and English",
			input:    "Hello こんにちは World",
			expected: "hello-こんにちは-world",
		},
		{
			name:     "accented characters (preserved)",
			input:    "Café Résumé",
			expected: "café-résumé",
		},
		{
			name:     "with parentheses and brackets",
			input:    "Test (2023) [Final]",
			expected: "test-2023-final",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "only special characters",
			input:    "!@#$%^&*()",
			expected: "",
		},
		{
			name:     "multiple consecutive hyphens",
			input:    "a---b---c",
			expected: "a-b-c",
		},
		{
			name:     "leading and trailing hyphens",
			input:    "---hello---",
			expected: "hello",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			data := struct {
				Text string
			}{
				Text: tt.input,
			}

			template := "{{.Text | slugify}}"
			result, err := tp.ProcessTemplate(template, data)
			if err != nil {
				t.Errorf("Expected no error, got %v", err)
			}

			if result != tt.expected {
				t.Errorf("slugify(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestTemplateProcessor_SlugifyMethod(t *testing.T) {
	tp := NewTemplateProcessor()

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "basic functionality",
			input:    "Hello World",
			expected: "hello-world",
		},
		{
			name:     "Hugo urlize behavior - spaces",
			input:    "A B C",
			expected: "a-b-c",
		},
		{
			name:     "Hugo urlize behavior - special chars",
			input:    "< a, b, & c >",
			expected: "a-b-c",
		},
		{
			name:     "Hugo urlize behavior - dots preserved",
			input:    "main.go",
			expected: "main.go",
		},
		{
			name:     "Hugo urlize behavior - accented chars preserved",
			input:    "Hugö",
			expected: "hugö",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tp.slugify(tt.input)
			if result != tt.expected {
				t.Errorf("slugify(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}
