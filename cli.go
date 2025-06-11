package main

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

// findProjectRoot searches for a Hugo project root by looking for content and static directories.
func findProjectRoot(startDir string) (string, error) {
	dir := startDir
	for {
		// contentとstaticフォルダの存在をチェック
		contentDir := filepath.Join(dir, "content")
		staticDir := filepath.Join(dir, "static")

		contentExists := false
		staticExists := false

		if _, err := os.Stat(contentDir); err == nil {
			contentExists = true
		}
		if _, err := os.Stat(staticDir); err == nil {
			staticExists = true
		}

		if contentExists && staticExists {
			return dir, nil
		}

		// 親ディレクトリに移動
		parent := filepath.Dir(dir)
		if parent == dir {
			// ルートディレクトリに到達
			break
		}
		dir = parent
	}

	return "", fmt.Errorf("Hugo project root not found (no content and static directories found)")
}

// listArticles displays all available articles with their titles and OGP settings.
func listArticles(contentDir string) error {
	fmt.Println("Available articles:")
	fmt.Println("==================")

	var count int
	err := filepath.WalkDir(contentDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if d.Name() == "index.md" {
			articleDir := filepath.Dir(path)
			relPath, err := filepath.Rel(contentDir, articleDir)
			if err != nil {
				return err
			}

			content, err := os.ReadFile(path)
			if err != nil {
				fmt.Printf("  %s (failed to read title)\n", relPath)
				count++
				return nil
			}

			fm, err := parseFrontMatter(content)
			if err != nil {
				fmt.Printf("  %s (failed to parse front matter)\n", relPath)
				count++
				return nil
			}

			title := fm.Title
			if title == "" {
				title = "(no title)"
			}

			ogpSettings := ""
			if fm.OGP != nil {
				var settings []string
				if fm.OGP.Text != nil && fm.OGP.Text.Content != nil {
					settings = append(settings, "custom content")
				}
				if fm.OGP.Overlay != nil {
					settings = append(settings, "overlay composition")
				}
				if fm.OGP.Text != nil {
					settings = append(settings, "custom text")
				}
				if len(settings) > 0 {
					ogpSettings = fmt.Sprintf(" [%s]", strings.Join(settings, ", "))
				}
			}

			fmt.Printf("  %s\n    Title: %s%s\n", relPath, title, ogpSettings)
			count++
		}

		return nil
	})

	if err != nil {
		return err
	}

	fmt.Printf("\nTotal: %d articles found\n", count)
	fmt.Println("\nUsage:")
	fmt.Println("  ogp-generator --test \"<article-path>\"")
	fmt.Println("  ogp-generator --single <project-root> \"<article-path>\"")

	return nil
}

// printUsage displays command-line usage information.
func printUsage() {
	fmt.Println("Usage:")
	fmt.Println("  ogp-generator <project-root> [config-file]           # Generate all OGP images")
	fmt.Println("  ogp-generator --single <project-root> <article-path> # Generate single article OGP")
	fmt.Println("  ogp-generator --test <article-directory-path>        # Test single article OGP (output to current dir)")
	fmt.Println("  ogp-generator --list <project-root>                  # List all available articles")
	fmt.Println("  ogp-generator --version                              # Show version information")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  ogp-generator --test \"../../content/Book/JPEGの裏側\"")
	fmt.Println("  ogp-generator --test \"/absolute/path/to/article/directory\"")
	fmt.Println("")
	fmt.Println("Notes:")
	fmt.Println("  - For --test mode, specify the full path to the article directory containing index.md")
	fmt.Println("  - Relative paths are resolved from current working directory")
}

// printVersion displays the application version.
func printVersion() {
	fmt.Printf("OGP Generator %s\n", version)
}

// CLIArgs represents parsed command-line arguments.
type CLIArgs struct {
	Mode        string
	ProjectRoot string
	ConfigPath  string
	ArticlePath string
}

// parseArgs parses command-line arguments and returns a CLIArgs structure.
func parseArgs(args []string) (*CLIArgs, error) {
	if len(args) < 2 {
		return nil, fmt.Errorf("insufficient arguments")
	}

	cli := &CLIArgs{}

	if args[1] == "--version" {
		cli.Mode = "--version"
		return cli, nil
	} else if args[1] == "--single" {
		if len(args) < 4 {
			return nil, fmt.Errorf("--single mode requires project-root and article-path")
		}

		cli.Mode = args[1]
		cli.ProjectRoot = args[2]
		cli.ArticlePath = args[3]

		execDir, _ := filepath.Abs(filepath.Dir(args[0]))
		cli.ConfigPath = filepath.Join(execDir, "config.yaml")
		if len(args) >= 5 {
			cli.ConfigPath = args[4]
		}
	} else if args[1] == "--test" {
		if len(args) < 3 {
			return nil, fmt.Errorf("--test mode requires article-path")
		}

		cli.Mode = args[1]
		// 記事パスを絶対パスに解決
		resolver := NewPathResolver("")
		articlePath, _ := resolver.ResolveFromCwd(args[2])
		cli.ArticlePath = articlePath
		cli.ProjectRoot = "" // testモードではプロジェクトルート不要

		execDir, _ := filepath.Abs(filepath.Dir(args[0]))
		cli.ConfigPath = filepath.Join(execDir, "config.yaml")
		if len(args) >= 4 {
			cli.ConfigPath = args[3]
		}
	} else if args[1] == "--list" {
		if len(args) < 3 {
			return nil, fmt.Errorf("--list mode requires project-root")
		}

		cli.Mode = args[1]
		cli.ProjectRoot = args[2]
	} else {
		cli.Mode = "--all"
		cli.ProjectRoot = args[1]

		execDir, _ := filepath.Abs(filepath.Dir(args[0]))
		cli.ConfigPath = filepath.Join(execDir, "config.yaml")
		if len(args) >= 3 {
			cli.ConfigPath = args[2]
		}
	}

	return cli, nil
}
