// Package main implements an OGP (Open Graph Protocol) image generator
// for static site generators like Hugo.
package main

import (
	"log"
	"os"
	"path/filepath"
)

// version is set during build time via ldflags
var version = "dev"

// main is the entry point of the OGP generator application.
// It handles command line arguments and orchestrates the image generation process.
func main() {
	cli, err := parseArgs(os.Args)
	if err != nil {
		printUsage()
		os.Exit(1)
	}

	if cli.Mode == "--version" {
		printVersion()
		return
	}

	var contentDir string
	if cli.Mode == "--test" {
		// testモードではcontentディレクトリは不要（記事パスを直接指定）
		contentDir = ""
	} else {
		contentDir = filepath.Join(cli.ProjectRoot, "content")
	}

	if cli.Mode == "--list" {
		err := listArticles(contentDir)
		if err != nil {
			log.Fatalf("Failed to list articles: %v", err)
		}
		return
	}

	generator, err := NewOGPGenerator(cli.ConfigPath, contentDir, cli.ProjectRoot)
	if err != nil {
		log.Fatalf("Failed to initialize OGP generator: %v", err)
	}

	switch cli.Mode {
	case "--single":
		err = generator.GenerateSingle(cli.ArticlePath)
		if err != nil {
			log.Fatalf("Failed to generate OGP image: %v", err)
		}
		log.Println("Single OGP image generation completed!")

	case "--test":
		err = generator.GenerateTest(cli.ArticlePath)
		if err != nil {
			log.Fatalf("Failed to generate test OGP image: %v", err)
		}

	default:
		err = generator.GenerateAll()
		if err != nil {
			log.Fatalf("Failed to generate OGP images: %v", err)
		}
		log.Println("OGP image generation completed!")
	}
}
