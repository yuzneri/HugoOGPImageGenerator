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
		contentDir := filepath.Join(dir, ContentDirectory)
		staticDir := filepath.Join(dir, StaticDirectory)

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

	return "", NewValidationError("Hugo project root not found (no content and static directories found)")
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

		if d.Name() == DefaultIndexFilename {
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
				if (fm.OGP.Title != nil && fm.OGP.Title.Content != nil) ||
					(fm.OGP.Description != nil && fm.OGP.Description.Content != nil) {
					settings = append(settings, "custom content")
				}
				if fm.OGP.Overlay != nil {
					settings = append(settings, "overlay composition")
				}
				if fm.OGP.Title != nil || fm.OGP.Description != nil {
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
	fmt.Println("  ogp-generator <project-root>                          # Generate all OGP images")
	fmt.Println("  ogp-generator --single <project-root> <article-path>  # Generate single article OGP")
	fmt.Println("  ogp-generator --test <article-directory-path>         # Test single article OGP (output to current dir)")
	fmt.Println("  ogp-generator --list <project-root>                   # List all available articles")
	fmt.Println("  ogp-generator --version                               # Show version information")
	fmt.Println("")
	fmt.Println("Global Options:")
	fmt.Println("  --config <config-file>   # Specify custom config file (can be used with any mode)")
	fmt.Println("                           # Default: config.yaml in executable directory")
	fmt.Println("")
	fmt.Println("Examples:")
	fmt.Println("  # Generate all with default config")
	fmt.Println("  ogp-generator /path/to/project")
	fmt.Println("")
	fmt.Println("  # Generate all with custom config")
	fmt.Println("  ogp-generator /path/to/project --config custom.yaml")
	fmt.Println("")
	fmt.Println("  # Test single article with custom config")
	fmt.Println("  ogp-generator --test \"/path/to/article\" --config custom.yaml")
	fmt.Println("")
	fmt.Println("  # Generate single article with custom config")
	fmt.Println("  ogp-generator --single /path/to/project \"article/path\" --config custom.yaml")
	fmt.Println("")
	fmt.Println("  # List articles (config file not required)")
	fmt.Println("  ogp-generator --list /path/to/project")
	fmt.Println("")
	fmt.Println("Notes:")
	fmt.Println("  - For --test mode, specify the full path to the article directory containing index.md")
	fmt.Println("  - Relative paths are resolved from current working directory")
	fmt.Println("  - The --config flag can be placed anywhere in the command line")
	fmt.Println("  - If --config flag is not specified, uses config.yaml from executable directory")
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

// parseConfigFlag extracts the --config flag value from arguments and returns the remaining args.
func parseConfigFlag(args []string) (configPath string, remainingArgs []string) {
	configPath = ""
	remainingArgs = make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		if args[i] == "--config" && i+1 < len(args) {
			configPath = args[i+1]
			i++ // Skip the config value
		} else {
			remainingArgs = append(remainingArgs, args[i])
		}
	}

	return configPath, remainingArgs
}

// parseArgs parses command-line arguments and returns a CLIArgs structure.
func parseArgs(args []string) (*CLIArgs, error) {
	if len(args) < 2 {
		return nil, NewValidationError("insufficient arguments")
	}

	// Extract --config flag from all arguments
	configFromFlag, filteredArgs := parseConfigFlag(args)

	cli := &CLIArgs{}

	// Set default config path
	execDir, _ := filepath.Abs(filepath.Dir(args[0]))
	cli.ConfigPath = filepath.Join(execDir, DefaultConfigFilename)

	// Override with --config flag if provided
	if configFromFlag != "" {
		cli.ConfigPath = configFromFlag
	}

	if filteredArgs[1] == "--version" {
		cli.Mode = "--version"
		return cli, nil
	} else if filteredArgs[1] == "--single" {
		if len(filteredArgs) < 4 {
			return nil, NewValidationError("--single mode requires project-root and article-path")
		}

		cli.Mode = filteredArgs[1]
		cli.ProjectRoot = filteredArgs[2]
		cli.ArticlePath = filteredArgs[3]
	} else if filteredArgs[1] == "--test" {
		if len(filteredArgs) < 3 {
			return nil, NewValidationError("--test mode requires article-path")
		}

		cli.Mode = filteredArgs[1]
		// 記事パスを絶対パスに解決
		resolver := NewPathResolver("")
		articlePath, _ := resolver.ResolveFromCwd(filteredArgs[2])
		cli.ArticlePath = articlePath
		cli.ProjectRoot = "" // testモードではプロジェクトルート不要
	} else if filteredArgs[1] == "--list" {
		if len(filteredArgs) < 3 {
			return nil, NewValidationError("--list mode requires project-root")
		}

		cli.Mode = filteredArgs[1]
		cli.ProjectRoot = filteredArgs[2]
	} else {
		cli.Mode = "--all"
		cli.ProjectRoot = filteredArgs[1]
	}

	return cli, nil
}
