package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// Article identifies a Hugo content file and the page it produces.
type Article struct {
	File    string // Content file (index.md, _index.md, or a single-file page such as about.md)
	Dir     string // Directory used for asset lookup and content type detection
	PageDir string // Directory under content/ that corresponds to the page URL
}

// resolveArticle normalizes a page bundle directory or a content file path into an Article.
// For a directory, index.md takes precedence over _index.md.
func resolveArticle(path string) (*Article, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, NewFileError("stat", path, err)
	}

	if info.IsDir() {
		for _, name := range []string{DefaultIndexFilename, BranchIndexFilename} {
			file := filepath.Join(path, name)
			if _, err := os.Stat(file); err == nil {
				return &Article{File: file, Dir: path, PageDir: path}, nil
			}
		}
		return nil, NewValidationError(fmt.Sprintf("no %s or %s found in %s", DefaultIndexFilename, BranchIndexFilename, path))
	}

	if filepath.Ext(path) != MarkdownExtension {
		return nil, NewValidationError(fmt.Sprintf("not a markdown content file: %s", path))
	}

	dir := filepath.Dir(path)
	name := filepath.Base(path)
	if name == DefaultIndexFilename || name == BranchIndexFilename {
		return &Article{File: path, Dir: dir, PageDir: dir}, nil
	}

	// Hugo publishes content/a/b.md at /a/b/, so the page directory is the file name without its extension.
	return &Article{File: path, Dir: dir, PageDir: filepath.Join(dir, strings.TrimSuffix(name, MarkdownExtension))}, nil
}

// findArticles returns every content file under contentDir that Hugo renders as a page.
// Inside a leaf bundle (a directory with index.md) only index.md is a page;
// the other markdown files there are page resources, so the bundle is not descended into.
func findArticles(contentDir string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.IsDir() {
			index := filepath.Join(path, DefaultIndexFilename)
			if _, err := os.Stat(index); err == nil {
				files = append(files, index)
				return filepath.SkipDir
			}
			return nil
		}

		if filepath.Ext(path) == MarkdownExtension {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return files, nil
}

// pageKind returns the Hugo page kind of a content file: home, section or page.
func pageKind(file, contentDir string) string {
	if filepath.Base(file) != BranchIndexFilename {
		return PageKindPage
	}
	if filepath.Clean(filepath.Dir(file)) == filepath.Clean(contentDir) {
		return PageKindHome
	}
	return PageKindSection
}
