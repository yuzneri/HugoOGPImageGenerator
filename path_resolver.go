package main

import (
	"os"
	"path/filepath"
)

// PathResolver handles path resolution for assets relative to config and article directories.
type PathResolver struct {
	configDir string
}

// NewPathResolver creates a new PathResolver with the given config directory.
func NewPathResolver(configDir string) *PathResolver {
	return &PathResolver{
		configDir: configDir,
	}
}

// ResolveAssetPath resolves asset paths with fallback from article to config directory.
// It checks the article directory first, then falls back to the config directory.
func (pr *PathResolver) ResolveAssetPath(assetPath string, articlePath string) string {
	if filepath.IsAbs(assetPath) {
		return assetPath
	}

	if articlePath != "" {
		articleAssetPath := filepath.Join(articlePath, assetPath)
		if _, err := os.Stat(articleAssetPath); err == nil {
			return articleAssetPath
		}
	}

	return filepath.Join(pr.configDir, assetPath)
}

// ResolveConfigAssetPath resolves asset paths relative to the config directory only.
func (pr *PathResolver) ResolveConfigAssetPath(assetPath string) string {
	if filepath.IsAbs(assetPath) {
		return assetPath
	}
	return filepath.Join(pr.configDir, assetPath)
}

// ResolvePath resolves a relative path against a base path.
func (pr *PathResolver) ResolvePath(basePath, relativePath string) string {
	if filepath.IsAbs(relativePath) {
		return relativePath
	}
	return filepath.Join(basePath, relativePath)
}

// ResolveFromCwd resolves a path relative to the current working directory.
func (pr *PathResolver) ResolveFromCwd(relativePath string) (string, error) {
	if filepath.IsAbs(relativePath) {
		return relativePath, nil
	}

	cwd, err := os.Getwd()
	if err != nil {
		return "", err
	}

	return filepath.Join(cwd, relativePath), nil
}
