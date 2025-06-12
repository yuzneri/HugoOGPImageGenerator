package main

import (
	"strings"
	"testing"
)

func TestParseFrontMatter_Success(t *testing.T) {
	content := `---
title: "Test Article"
description: "This is a test"
date: 2023-12-25T15:30:45Z
tags: ["test", "example"]
url: "/test"
ogp:
  title:
    content: "Custom OGP Title"
    size: 48
---

# Article Content

This is the main content of the article.
`

	fm, err := parseFrontMatter([]byte(content))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if fm.Title != "Test Article" {
		t.Errorf("Expected title 'Test Article', got %q", fm.Title)
	}

	if fm.Description != "This is a test" {
		t.Errorf("Expected description 'This is a test', got %q", fm.Description)
	}

	if fm.URL != "/test" {
		t.Errorf("Expected URL '/test', got %q", fm.URL)
	}

	if len(fm.Tags) != 2 {
		t.Errorf("Expected 2 tags, got %d", len(fm.Tags))
	}

	if fm.Tags[0] != "test" || fm.Tags[1] != "example" {
		t.Errorf("Expected tags [test, example], got %v", fm.Tags)
	}

	if fm.OGP == nil {
		t.Error("Expected OGP settings to be present")
	}
}

func TestParseFrontMatter_NoFrontMatter(t *testing.T) {
	// Content without front matter delimiters
	content := `# Article Content

This is an article without front matter.
`

	_, err := parseFrontMatter([]byte(content))
	if err == nil {
		t.Error("Expected error for content without front matter")
	}

	expectedError := "no front matter found"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain %q, got %q", expectedError, err.Error())
	}
}

func TestParseFrontMatter_NoEndDelimiter(t *testing.T) {
	// Content with opening delimiter but no closing delimiter
	content := `---
title: "Test Article"
description: "This is a test"

# Article Content

This content has front matter start but no end delimiter.
`

	_, err := parseFrontMatter([]byte(content))
	if err == nil {
		t.Error("Expected error for front matter without end delimiter")
	}

	expectedError := "front matter end delimiter not found"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain %q, got %q", expectedError, err.Error())
	}
}

func TestParseFrontMatter_InvalidYAML(t *testing.T) {
	// Content with invalid YAML in front matter
	content := `---
title: "Test Article"
description: "This is a test
tags: [invalid yaml structure
invalid: yaml: content
---

# Article Content

This content has invalid YAML in front matter.
`

	_, err := parseFrontMatter([]byte(content))
	if err == nil {
		t.Error("Expected error for invalid YAML in front matter")
	}

	expectedError := "failed to unmarshal front matter"
	if !strings.Contains(err.Error(), expectedError) {
		t.Errorf("Expected error to contain %q, got %q", expectedError, err.Error())
	}
}

func TestParseFrontMatter_MinimalFrontMatter(t *testing.T) {
	// Minimal valid front matter
	content := `---
title: "Minimal Article"
---

# Minimal Content
`

	fm, err := parseFrontMatter([]byte(content))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if fm.Title != "Minimal Article" {
		t.Errorf("Expected title 'Minimal Article', got %q", fm.Title)
	}

	// Other fields should have zero values
	if fm.Description != "" {
		t.Errorf("Expected empty description, got %q", fm.Description)
	}

	if len(fm.Tags) != 0 {
		t.Errorf("Expected no tags, got %v", fm.Tags)
	}

	if fm.OGP != nil {
		t.Error("Expected no OGP settings")
	}
}

func TestParseFrontMatter_EmptyFrontMatter(t *testing.T) {
	// Empty front matter section - the parser expects \n---\n not just ---
	content := `---

---

# Article with empty front matter
`

	fm, err := parseFrontMatter([]byte(content))
	if err != nil {
		t.Fatalf("Expected no error for empty front matter, got %v", err)
	}

	// All fields should have zero values
	if fm.Title != "" {
		t.Errorf("Expected empty title, got %q", fm.Title)
	}

	if fm.Description != "" {
		t.Errorf("Expected empty description, got %q", fm.Description)
	}
}

func TestParseFrontMatter_ComplexFrontMatter(t *testing.T) {
	// Complex front matter with nested structures
	content := `---
title: "Complex Article"
description: "A complex test article"
date: 2023-12-25T15:30:45Z
tags: ["golang", "testing", "yaml"]
url: "/articles/complex"
author: "Test Author"
category: "Technology"
draft: false
ogp:
  title:
    content: "Custom OGP: {{.Title}}"
    size: 52
    color: "#2c3e50"
  description:
    visible: true
    content: "Published on {{dateFormat \"2006-01-02\" .Date}}"
  background:
    color: "#ecf0f1"
---

# Complex Article

This article has complex front matter.
`

	fm, err := parseFrontMatter([]byte(content))
	if err != nil {
		t.Fatalf("Expected no error, got %v", err)
	}

	if fm.Title != "Complex Article" {
		t.Errorf("Expected title 'Complex Article', got %q", fm.Title)
	}

	if fm.Description != "A complex test article" {
		t.Errorf("Expected description 'A complex test article', got %q", fm.Description)
	}

	if len(fm.Tags) != 3 {
		t.Errorf("Expected 3 tags, got %d", len(fm.Tags))
	}

	// Check that additional fields are captured
	if fm.Fields["author"] != "Test Author" {
		t.Errorf("Expected author 'Test Author', got %v", fm.Fields["author"])
	}

	if fm.Fields["category"] != "Technology" {
		t.Errorf("Expected category 'Technology', got %v", fm.Fields["category"])
	}

	if fm.Fields["draft"] != false {
		t.Errorf("Expected draft false, got %v", fm.Fields["draft"])
	}

	if fm.OGP == nil {
		t.Fatal("Expected OGP settings to be present")
	}

	if fm.OGP.Title == nil {
		t.Error("Expected OGP title settings")
	}

	if fm.OGP.Description == nil {
		t.Error("Expected OGP description settings")
	}

	if fm.OGP.Background == nil {
		t.Error("Expected OGP background settings")
	}
}

func TestParseFrontMatter_WindowsLineEndings(t *testing.T) {
	// Test with Windows line endings (CRLF) - currently not supported by parseFrontMatter
	// The function expects Unix line endings (LF only)
	content := "---\r\ntitle: \"Windows Article\"\r\ndescription: \"Test with CRLF\"\r\n---\r\n\r\n# Windows Content"

	_, err := parseFrontMatter([]byte(content))
	if err == nil {
		t.Error("Expected error with Windows line endings as current implementation doesn't support CRLF")
	}

	// Test that Unix line endings work
	unixContent := "---\ntitle: \"Unix Article\"\ndescription: \"Test with LF\"\n---\n\n# Unix Content"

	fm, err := parseFrontMatter([]byte(unixContent))
	if err != nil {
		t.Fatalf("Expected no error with Unix line endings, got %v", err)
	}

	if fm.Title != "Unix Article" {
		t.Errorf("Expected title 'Unix Article', got %q", fm.Title)
	}
}
