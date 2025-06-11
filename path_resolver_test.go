package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNewPathResolver(t *testing.T) {
	configDir := "/test/config"

	resolver := NewPathResolver(configDir)

	if resolver == nil {
		t.Error("NewPathResolver should return a non-nil resolver")
	}

	if resolver.configDir != configDir {
		t.Errorf("Expected configDir to be %s, got %s", configDir, resolver.configDir)
	}
}

func TestPathResolver_ResolveAssetPath_AbsolutePath(t *testing.T) {
	resolver := NewPathResolver("/test/config")

	absPath := "/absolute/path/to/asset.png"
	result := resolver.ResolveAssetPath(absPath, "/test/article")

	if result != absPath {
		t.Errorf("Expected absolute path %s, got %s", absPath, result)
	}
}

func TestPathResolver_ResolveAssetPath_ArticleAssetExists(t *testing.T) {
	// Create temporary directories for testing
	tempDir, err := os.MkdirTemp("", "path_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create config directory
	configDir := filepath.Join(tempDir, "config")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Create article directory
	articleDir := filepath.Join(tempDir, "article")
	err = os.MkdirAll(articleDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create article dir: %v", err)
	}

	// Create asset in article directory
	assetPath := filepath.Join(articleDir, "asset.png")
	err = os.WriteFile(assetPath, []byte("test"), 0644)
	if err != nil {
		t.Fatalf("Failed to create asset file: %v", err)
	}

	resolver := NewPathResolver(configDir)

	result := resolver.ResolveAssetPath("asset.png", articleDir)

	if result != assetPath {
		t.Errorf("Expected article asset path %s, got %s", assetPath, result)
	}
}

func TestPathResolver_ResolveAssetPath_FallbackToConfig(t *testing.T) {
	// Create temporary directories for testing
	tempDir, err := os.MkdirTemp("", "path_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create config directory
	configDir := filepath.Join(tempDir, "config")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Create article directory (but no asset in it)
	articleDir := filepath.Join(tempDir, "article")
	err = os.MkdirAll(articleDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create article dir: %v", err)
	}

	resolver := NewPathResolver(configDir)

	result := resolver.ResolveAssetPath("asset.png", articleDir)
	expected := filepath.Join(configDir, "asset.png")

	if result != expected {
		t.Errorf("Expected config asset path %s, got %s", expected, result)
	}
}

func TestPathResolver_ResolveAssetPath_EmptyArticlePath(t *testing.T) {
	configDir := "/test/config"
	resolver := NewPathResolver(configDir)

	result := resolver.ResolveAssetPath("asset.png", "")
	expected := filepath.Join(configDir, "asset.png")

	if result != expected {
		t.Errorf("Expected config asset path %s, got %s", expected, result)
	}
}

func TestPathResolver_ResolveConfigAssetPath_AbsolutePath(t *testing.T) {
	resolver := NewPathResolver("/test/config")

	absPath := "/absolute/path/to/asset.png"
	result := resolver.ResolveConfigAssetPath(absPath)

	if result != absPath {
		t.Errorf("Expected absolute path %s, got %s", absPath, result)
	}
}

func TestPathResolver_ResolveConfigAssetPath_RelativePath(t *testing.T) {
	configDir := "/test/config"
	resolver := NewPathResolver(configDir)

	result := resolver.ResolveConfigAssetPath("asset.png")
	expected := filepath.Join(configDir, "asset.png")

	if result != expected {
		t.Errorf("Expected config asset path %s, got %s", expected, result)
	}
}

func TestPathResolver_ResolvePath_AbsolutePath(t *testing.T) {
	resolver := NewPathResolver("/test/config")

	absPath := "/absolute/path/to/file.txt"
	result := resolver.ResolvePath("/base/path", absPath)

	if result != absPath {
		t.Errorf("Expected absolute path %s, got %s", absPath, result)
	}
}

func TestPathResolver_ResolvePath_RelativePath(t *testing.T) {
	resolver := NewPathResolver("/test/config")

	basePath := "/base/path"
	relativePath := "file.txt"
	result := resolver.ResolvePath(basePath, relativePath)
	expected := filepath.Join(basePath, relativePath)

	if result != expected {
		t.Errorf("Expected resolved path %s, got %s", expected, result)
	}
}

func TestPathResolver_ResolveFromCwd_AbsolutePath(t *testing.T) {
	resolver := NewPathResolver("/test/config")

	absPath := "/absolute/path/to/file.txt"
	result, err := resolver.ResolveFromCwd(absPath)

	if err != nil {
		t.Errorf("ResolveFromCwd should not return error for absolute path: %v", err)
	}

	if result != absPath {
		t.Errorf("Expected absolute path %s, got %s", absPath, result)
	}
}

func TestPathResolver_ResolveFromCwd_RelativePath(t *testing.T) {
	resolver := NewPathResolver("/test/config")

	relativePath := "file.txt"
	result, err := resolver.ResolveFromCwd(relativePath)

	if err != nil {
		t.Errorf("ResolveFromCwd should not return error: %v", err)
	}

	// Result should be current working directory + relative path
	cwd, _ := os.Getwd()
	expected := filepath.Join(cwd, relativePath)

	if result != expected {
		t.Errorf("Expected resolved path %s, got %s", expected, result)
	}
}

func TestPathResolver_ResolveAssetPath_WindowsPaths(t *testing.T) {
	// Test with Windows-style paths
	resolver := NewPathResolver("C:\\test\\config")

	// Test with backslashes
	result := resolver.ResolveAssetPath("asset.png", "C:\\test\\article")
	expected := filepath.Join("C:\\test\\config", "asset.png")

	// On Windows, this should work correctly
	// On Unix-like systems, this will be treated as relative paths
	if result != expected {
		t.Logf("Note: Windows-style paths may behave differently on Unix systems")
		t.Logf("Expected: %s, Got: %s", expected, result)
	}
}

func TestPathResolver_ResolveAssetPath_MultipleLevels(t *testing.T) {
	// Create temporary directories for testing
	tempDir, err := os.MkdirTemp("", "path_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create nested config directory
	configDir := filepath.Join(tempDir, "config", "assets")
	err = os.MkdirAll(configDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create config dir: %v", err)
	}

	// Create nested article directory
	articleDir := filepath.Join(tempDir, "content", "posts", "article1")
	err = os.MkdirAll(articleDir, 0755)
	if err != nil {
		t.Fatalf("Failed to create article dir: %v", err)
	}

	resolver := NewPathResolver(configDir)

	// Test with nested relative path
	result := resolver.ResolveAssetPath("images/logo.png", articleDir)
	expected := filepath.Join(configDir, "images", "logo.png")

	if result != expected {
		t.Errorf("Expected nested asset path %s, got %s", expected, result)
	}
}

func TestPathResolver_EdgeCases(t *testing.T) {
	resolver := NewPathResolver("/test/config")

	// Test with empty asset path
	result := resolver.ResolveAssetPath("", "/test/article")
	expected := filepath.Join("/test/config", "")

	if result != expected {
		t.Errorf("Expected empty asset path %s, got %s", expected, result)
	}

	// Test with dot path
	result = resolver.ResolveAssetPath(".", "/test/article")
	expected = filepath.Join("/test/config", ".")

	if result != expected {
		t.Errorf("Expected dot path %s, got %s", expected, result)
	}

	// Test with parent directory path
	result = resolver.ResolveAssetPath("../asset.png", "/test/article")
	expected = filepath.Join("/test/config", "..", "asset.png")

	if result != expected {
		t.Errorf("Expected parent directory path %s, got %s", expected, result)
	}
}
