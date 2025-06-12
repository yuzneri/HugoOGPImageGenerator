package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
)

// resolveAssetPath resolves asset paths (fonts, images) relative to config directory.
// It uses the same path resolution logic as font files for consistency.
func resolveAssetPath(assetPath, configDir string) string {
	resolver := NewPathResolver(configDir)
	return resolver.ResolveConfigAssetPath(assetPath)
}

// loadImage loads an image from the filesystem, supporting JPEG and PNG formats.
// It automatically detects the format based on file extension.
func loadImage(imagePath string) (image.Image, error) {
	file, err := os.Open(imagePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open image: %w", err)
	}
	defer file.Close()

	ext := strings.ToLower(filepath.Ext(imagePath))
	switch ext {
	case ".jpg", ".jpeg":
		return jpeg.Decode(file)
	case ".png":
		return png.Decode(file)
	default:
		img, _, err := image.Decode(file)
		return img, err
	}
}

// resizeImage resizes an image according to the specified fit method.
// Supported fit methods: "cover" (fill area, may crop), "contain" (fit within area),
// "fill" (stretch to exact dimensions), or "none" (no resizing).
func resizeImage(src image.Image, targetWidth, targetHeight int, fit string) image.Image {
	srcBounds := src.Bounds()
	srcWidth := srcBounds.Dx()
	srcHeight := srcBounds.Dy()

	if targetWidth <= 0 || targetHeight <= 0 {
		return src
	}

	var dstWidth, dstHeight int

	switch fit {
	case "cover":
		scaleX := float64(targetWidth) / float64(srcWidth)
		scaleY := float64(targetHeight) / float64(srcHeight)
		scale := math.Max(scaleX, scaleY)
		dstWidth = int(float64(srcWidth) * scale)
		dstHeight = int(float64(srcHeight) * scale)
	case "contain":
		scaleX := float64(targetWidth) / float64(srcWidth)
		scaleY := float64(targetHeight) / float64(srcHeight)
		scale := math.Min(scaleX, scaleY)
		dstWidth = int(float64(srcWidth) * scale)
		dstHeight = int(float64(srcHeight) * scale)
	case "fill":
		dstWidth = targetWidth
		dstHeight = targetHeight
	default:
		return src
	}

	if math.Abs(float64(dstWidth-srcWidth)) <= 2 && math.Abs(float64(dstHeight-srcHeight)) <= 2 {
		return src
	}

	return imaging.Resize(src, dstWidth, dstHeight, imaging.Lanczos)
}

// compositeImage composites a source image onto a destination image at the specified position.
// It supports alpha blending with the given opacity (0.0 to 1.0).
func compositeImage(dst *image.RGBA, src image.Image, x, y int, opacity float64) {
	srcBounds := src.Bounds()
	dstBounds := dst.Bounds()

	for sy := srcBounds.Min.Y; sy < srcBounds.Max.Y; sy++ {
		for sx := srcBounds.Min.X; sx < srcBounds.Max.X; sx++ {
			dx := x + sx - srcBounds.Min.X
			dy := y + sy - srcBounds.Min.Y

			if dx < dstBounds.Min.X || dx >= dstBounds.Max.X ||
				dy < dstBounds.Min.Y || dy >= dstBounds.Max.Y {
				continue
			}

			srcColor := color.RGBAModel.Convert(src.At(sx, sy)).(color.RGBA)
			dstColor := color.RGBAModel.Convert(dst.At(dx, dy)).(color.RGBA)

			alpha := float64(srcColor.A) * opacity / 255.0
			invAlpha := 1.0 - alpha

			blendedR := uint8(float64(srcColor.R)*alpha + float64(dstColor.R)*invAlpha)
			blendedG := uint8(float64(srcColor.G)*alpha + float64(dstColor.G)*invAlpha)
			blendedB := uint8(float64(srcColor.B)*alpha + float64(dstColor.B)*invAlpha)
			blendedA := uint8(math.Max(float64(srcColor.A)*opacity, float64(dstColor.A)))

			dst.Set(dx, dy, color.RGBA{R: blendedR, G: blendedG, B: blendedB, A: blendedA})
		}
	}
}

// OverlaySettings defines the interface for overlay configuration
type OverlaySettings interface {
	GetImage() *string
	GetPlacement() *PlacementConfig
	GetFit() *string
	GetOpacity() *float64
}

// compositeCustomImage composites an overlay image with full configuration support.
// It handles path resolution, resizing, cropping (for cover fit), and alpha blending.
// The isConfigOverlay parameter determines whether to use config-relative or article-relative paths.
func compositeCustomImage(dst *image.RGBA, basePath string, overlaySettings OverlaySettings, isConfigOverlay bool, configDir string) error {
	var imagePath string
	imagePtr := overlaySettings.GetImage()
	if imagePtr == nil {
		return fmt.Errorf("overlay image is nil")
	}
	if isConfigOverlay {
		// For config overlays, use the same path resolution as fonts/background images
		imagePath = resolveAssetPath(*imagePtr, configDir)
	} else {
		// For front matter overlays, try article directory first, then config directory
		if filepath.IsAbs(*imagePtr) {
			imagePath = *imagePtr
		} else {
			articleImagePath := filepath.Join(basePath, *imagePtr)
			if _, err := os.Stat(articleImagePath); err == nil {
				imagePath = articleImagePath
			} else {
				imagePath = resolveAssetPath(*imagePtr, configDir)
			}
		}
	}

	img, err := loadImage(imagePath)
	if err != nil {
		return fmt.Errorf("failed to load image %s: %w", imagePath, err)
	}

	x, y := 0, 0
	originalWidth := img.Bounds().Dx()
	originalHeight := img.Bounds().Dy()
	width, height := originalWidth, originalHeight
	fit := "contain"
	opacity := 1.0

	var widthSpecified, heightSpecified bool
	placement := overlaySettings.GetPlacement()
	if placement != nil {
		x = placement.X
		y = placement.Y
		if placement.Width != 0 {
			width = placement.Width
			widthSpecified = true
		}
		if placement.Height != 0 {
			height = placement.Height
			heightSpecified = true
		}
	}

	if widthSpecified && !heightSpecified {
		aspectRatio := float64(originalHeight) / float64(originalWidth)
		height = int(float64(width) * aspectRatio)
	} else if !widthSpecified && heightSpecified {
		aspectRatio := float64(originalWidth) / float64(originalHeight)
		width = int(float64(height) * aspectRatio)
	}

	if overlaySettings.GetFit() != nil {
		fit = *overlaySettings.GetFit()
	}

	if overlaySettings.GetOpacity() != nil {
		opacity = *overlaySettings.GetOpacity()
		if opacity < 0 {
			opacity = 0
		}
		if opacity > 1 {
			opacity = 1
		}
	}

	resizedImg := resizeImage(img, width, height, fit)

	if fit == "cover" {
		resizedBounds := resizedImg.Bounds()
		resizedWidth := resizedBounds.Dx()
		resizedHeight := resizedBounds.Dy()

		if resizedWidth > width || resizedHeight > height {
			cropX := (resizedWidth - width) / 2
			cropY := (resizedHeight - height) / 2

			croppedImg := image.NewRGBA(image.Rect(0, 0, width, height))

			for cy := 0; cy < height; cy++ {
				for cx := 0; cx < width; cx++ {
					srcX := cropX + cx
					srcY := cropY + cy

					if srcX < resizedBounds.Max.X && srcY < resizedBounds.Max.Y {
						at := resizedImg.At(resizedBounds.Min.X+srcX, resizedBounds.Min.Y+srcY)
						croppedImg.Set(cx, cy, at)
					}
				}
			}

			resizedImg = croppedImg
		}
	}

	compositeImage(dst, resizedImg, x, y, opacity)

	return nil
}
