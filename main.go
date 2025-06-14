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
		// testモードでは記事パスから contentDir を推測
		// 記事パスから遡って content ディレクトリを探す
		dir := filepath.Dir(cli.ArticlePath)
		contentDir = ""
		for {
			if filepath.Base(dir) == ContentDirectory {
				contentDir = dir
				break
			}
			parent := filepath.Dir(dir)
			if parent == dir {
				// ルートに到達
				break
			}
			dir = parent
		}

		// content ディレクトリが見つからない場合は、現在のディレクトリを使用
		if contentDir == "" {
			cwd, _ := os.Getwd()
			contentDir = filepath.Join(cwd, "test_content")
		}
	} else {
		contentDir = filepath.Join(cli.ProjectRoot, ContentDirectory)
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
