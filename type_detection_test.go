package main

import (
	"path/filepath"
	"testing"
)

func TestDetermineContentType(t *testing.T) {
	tests := []struct {
		name         string
		frontMatter  *FrontMatter
		articlePath  string
		hugoRootPath string
		expected     string
	}{
		{
			name: "front matter type takes precedence",
			frontMatter: &FrontMatter{
				Type: "custom",
			},
			articlePath:  "/hugo/content/blog/article1",
			hugoRootPath: "/hugo",
			expected:     "custom",
		},
		{
			name:         "directory-based type detection",
			frontMatter:  &FrontMatter{},
			articlePath:  "/hugo/content/blog/article1",
			hugoRootPath: "/hugo",
			expected:     "blog",
		},
		{
			name:         "nested directory - first level determines type",
			frontMatter:  &FrontMatter{},
			articlePath:  "/hugo/content/posts/2023/january/article1",
			hugoRootPath: "/hugo",
			expected:     "posts",
		},
		{
			name:         "content root direct file",
			frontMatter:  &FrontMatter{},
			articlePath:  "/hugo/content/about",
			hugoRootPath: "/hugo",
			expected:     "about",
		},
		{
			name: "empty type field in front matter falls back to directory",
			frontMatter: &FrontMatter{
				Type: "",
			},
			articlePath:  "/hugo/content/news/breaking",
			hugoRootPath: "/hugo",
			expected:     "news",
		},
		{
			name:         "relative path from content root",
			frontMatter:  &FrontMatter{},
			articlePath:  "/hugo/content/events/meetup",
			hugoRootPath: "/hugo",
			expected:     "events",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := determineContentType(tt.frontMatter, tt.articlePath, tt.hugoRootPath)
			if result != tt.expected {
				t.Errorf("determineContentType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestExtractDirectoryType(t *testing.T) {
	tests := []struct {
		name         string
		articlePath  string
		hugoRootPath string
		expected     string
	}{
		{
			name:         "simple blog path",
			articlePath:  "/hugo/content/blog/article1",
			hugoRootPath: "/hugo",
			expected:     "blog",
		},
		{
			name:         "nested path",
			articlePath:  "/hugo/content/posts/2023/article1",
			hugoRootPath: "/hugo",
			expected:     "posts",
		},
		{
			name:         "single level under content",
			articlePath:  "/hugo/content/about",
			hugoRootPath: "/hugo",
			expected:     "about",
		},
		{
			name:         "content root equals article path",
			articlePath:  "/hugo/content",
			hugoRootPath: "/hugo",
			expected:     "page",
		},
		{
			name:         "deep nesting",
			articlePath:  "/hugo/content/docs/guides/advanced/config",
			hugoRootPath: "/hugo",
			expected:     "docs",
		},
		{
			name:         "windows path",
			articlePath:  "C:\\hugo\\content\\blog\\article1",
			hugoRootPath: "C:\\hugo",
			expected:     "page", // Should return page due to absolute path handling
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := extractDirectoryType(tt.articlePath, tt.hugoRootPath)
			if result != tt.expected {
				t.Errorf("extractDirectoryType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestFrontMatterWithType(t *testing.T) {
	tests := []struct {
		name     string
		yaml     string
		expected string
	}{
		{
			name: "type field present",
			yaml: `---
title: "Test Article"
type: "tutorial"
---
Content here`,
			expected: "tutorial",
		},
		{
			name: "no type field",
			yaml: `---
title: "Test Article"
description: "A test"
---
Content here`,
			expected: "",
		},
		{
			name: "empty type field",
			yaml: `---
title: "Test Article"
type: ""
---
Content here`,
			expected: "",
		},
		{
			name: "type with other fields",
			yaml: `---
title: "Test Article"
date: 2023-01-01
type: "documentation"
tags: ["test"]
---
Content here`,
			expected: "documentation",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fm, err := parseFrontMatter([]byte(tt.yaml))
			if err != nil {
				t.Fatalf("Failed to parse front matter: %v", err)
			}

			if fm.Type != tt.expected {
				t.Errorf("Type = %v, want %v", fm.Type, tt.expected)
			}
		})
	}
}

// Helper function tests for path manipulation
func TestPathHelpers(t *testing.T) {
	tests := []struct {
		name     string
		path     string
		expected []string
	}{
		{
			name:     "unix path",
			path:     "/hugo/content/blog/article1",
			expected: []string{"", "hugo", "content", "blog", "article1"},
		},
		{
			name:     "windows path",
			path:     "C:\\hugo\\content\\blog\\article1",
			expected: []string{"C:", "hugo", "content", "blog", "article1"},
		},
		{
			name:     "relative path",
			path:     "content/blog/article1",
			expected: []string{"content", "blog", "article1"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test the path splitting logic that will be used in implementation
			parts := filepath.SplitList(tt.path)
			if len(parts) == 1 {
				// SplitList doesn't work as expected for paths, use manual split
				parts = splitPath(tt.path)
			}

			// This test is mainly to understand filepath behavior
			// The actual implementation may differ
			t.Logf("Path parts for %s: %v", tt.path, parts)
		})
	}
}

// Helper function that will be needed for implementation
func splitPath(path string) []string {
	return filepath.SplitList(path)
}
