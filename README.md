# Hugo OGP Image Generator

A Go application that automatically generates Open Graph Protocol (OGP) images for Hugo static sites.
It processes Hugo content with front matter and creates custom social media preview images with advanced text rendering, background processing, and overlay composition.

## Installation

```bash
go build -o ogp
```

## Usage

### Generate OGP for all articles
```bash
./ogp /path/to/hugo/project
```

### Generate for a specific article
```bash
./ogp --single /path/to/hugo/project "article/path"
```

### Test mode
```bash
./ogp --test "/path/to/article/directory"
```

### List all articles with OGP settings
```bash
./ogp --list /path/to/hugo/project
```

## Configuration

The application uses YAML configuration files with comprehensive default values.
All configuration options can be overridden per-article via Hugo front matter.

### Default Values

When no `config.yaml` file is provided, or when specific options are omitted, the following defaults are used:

```yaml
background:
  color: "#FFFFFF"  # White background
  # image: null     # No background image by default

output:
  directory: "public"  # Output to Hugo's public directory
  format: "png"        # PNG format for images
  # filename: null     # Default: "ogp.{format}"

# Title text configuration
title:
  visible: true        # Whether to render title (default: true)
  # content: null      # No default content template (uses article title)
  font: ""             # Auto-detect system fonts
  size: 64             # Font size in pixels
  color: "#000000"     # Black text
  area:
    x: 100             # Text area X position
    y: 50              # Text area Y position
    width: 1000        # Text area width
    height: 250        # Text area height
  block_position: "middle-center" # Text block position in area
  line_alignment: "center"        # Individual line alignment
  overflow: "shrink"              # Shrink text to fit area
  min_size: 24.0                  # Minimum font size when shrinking
  line_height: 1.2                # Line height multiplier
  letter_spacing: 1               # Letter spacing in pixels
  line_breaking:
    start_prohibited: ".)}]>!?、。，．！？)）］｝〉》」』ー～ぁぃぅぇぉっゃゅょゎァィゥェォッャュョヮヵヶ々"  # Japanese line-start prohibited chars
    end_prohibited: "({[<（［｛〈《「『"                                                            # Japanese line-end prohibited chars

# Description text configuration
description:
  visible: false       # Whether to render description (default: false)
  # content: null      # No default content template (uses article description)
  font: ""             # Auto-detect system fonts
  size: 32             # Font size in pixels
  color: "#666666"     # Gray text
  area:
    x: 100             # Text area X position
    y: 350             # Text area Y position
    width: 1000        # Text area width
    height: 200        # Text area height
  block_position: "top-left"      # Text block position in area
  line_alignment: "left"          # Individual line alignment
  overflow: "clip"                # Overflow handling
  min_size: 16.0                  # Minimum font size when shrinking
  line_height: 1.2                # Line height multiplier
  letter_spacing: 0               # Letter spacing in pixels
  line_breaking:
    start_prohibited: ".)}]>!?、。，．！？)）］｝〉》」』ー～ぁぃぅぇぉっゃゅょゎァィゥェォッャュョヮヵヶ々"  # Japanese line-start prohibited chars
    end_prohibited: "({[<（［｛〈《「『"                                                            # Japanese line-end prohibited chars

# Overlay configuration
overlay:
  visible: false       # Whether to render overlay (default: false)
  # image: null        # No overlay image by default
  # placement:         # Placement configuration when overlay is used
  #   x: 0
  #   y: 0
  #   width: 100
  #   height: 100
  # fit: "contain"     # Default fit method when overlay is configured
  # opacity: 1.0       # Default opacity when overlay is configured
```

### Configuration Options

#### Background Settings
```yaml
background:
  color: "#FFFFFF"           # Hex color (6 or 8 chars: #RRGGBB or #RRGGBBAA)
  image: "path/to/image.jpg" # Background image path (optional)
```

#### Output Settings
```yaml
output:
  directory: "public"                    # Output directory
  format: "png"                          # Image format (png, jpg)
  filename: "custom-{{.Title}}.{{.Format}}" # Custom filename template (optional)
```

#### Title Text Configuration
```yaml
title:
  visible: true                                # Whether to render title (default: true)
  content: "{{.Title}}"  # Content template (optional)
  font: "fonts/custom.ttf"                     # Font file path (empty = auto-detect)
  size: 72                                     # Font size in pixels
  color: "#000000"                             # Text color (hex)
  area:                                        # Text rendering area
    x: 50                                      # X position
    y: 50                                      # Y position
    width: 1200                                # Area width
    height: 300                                # Area height
  block_position: "middle-center"              # Text block position
  line_alignment: "center"                     # Line alignment within block
  overflow: "shrink"                           # Overflow handling
  min_size: 24.0                               # Minimum font size for shrink mode
  line_height: 1.2                             # Line height multiplier
  letter_spacing: 1                            # Letter spacing in pixels
  line_breaking:                               # Japanese line breaking rules
    start_prohibited: "、。！？」』）"          # Chars that cannot start a line
    end_prohibited: "「『（"                   # Chars that cannot end a line
```

#### Description Text Configuration
```yaml
description:
  visible: true                                # Whether to render description (default: true)
  content: "{{.Description}}"                  # Content template (optional)
  font: "fonts/description.ttf"               # Font file path (can be different from title)
  size: 32                                     # Font size in pixels
  color: "#666666"                             # Text color (hex)
  area:                                        # Text rendering area
    x: 50                                      # X position
    y: 380                                     # Y position (below title)
    width: 1200                                # Area width
    height: 200                                # Area height
  block_position: "top-left"                  # Text block position
  line_alignment: "left"                       # Line alignment within block
  overflow: "shrink"                           # Overflow handling
  min_size: 16.0                               # Minimum font size for shrink mode
  line_height: 1.4                             # Line height multiplier
  letter_spacing: 0                            # Letter spacing in pixels
  line_breaking:                               # Japanese line breaking rules
    start_prohibited: "、。！？」』）"          # Chars that cannot start a line
    end_prohibited: "「『（"                   # Chars that cannot end a line
```

**Text Block Position Options:**
- `top-left`, `top-center`, `top-right`
- `middle-left`, `middle-center`, `middle-right`  
- `bottom-left`, `bottom-center`, `bottom-right`

**Line Alignment Options:** `left`, `center`, `right`

**Overflow Options:**
- `shrink`: Reduce font size to fit text in area
- `clip`: Truncate text that doesn't fit

#### Overlay Configuration
```yaml
overlay:
  visible: true             # Whether to render overlay (default: false)
  image: "overlay.png"      # Overlay image path
  placement:                # Image positioning
    x: 100                  # X position
    y: 100                  # Y position
    width: 300              # Image width
    height: 200             # Image height
  fit: "cover"              # Image fit method
  opacity: 0.8              # Image opacity (0.0-1.0)
```

**Overlay Fit Options:**
- `cover`: Scale image to cover area (may crop)
- `contain`: Scale image to fit within area (preserves aspect ratio)
- `fill`: Stretch image to fill exact dimensions
- `none`: Use original image size

### Global Configuration

Create a `config.yaml` file in the same directory as the executable, or specify a custom path with the `--config` flag:

```yaml
background:
  color: "#F5F5F5"
  image: "assets/bg.jpg"

output:
  directory: "static/images"
  format: "jpg"
  filename: "ogp-{{.Title | slugify}}.{{.Format}}"

title:
  visible: true
  content: "{{.Title | upper}}"
  font: "fonts/NotoSansCJK-Bold.ttc"
  size: 64
  color: "#2C3E50"
  area:
    x: 80
    y: 80
    width: 1040
    height: 240
  block_position: "middle-center"
  line_alignment: "center"
  overflow: "shrink"
  min_size: 24.0
  line_height: 1.2
  letter_spacing: 1

description:
  visible: true
  content: "{{.Description}} - {{.Fields.author | default \"Blog\"}}"
  font: "fonts/NotoSansCJK-Regular.ttc"
  size: 28
  color: "#7F8C8D"
  area:
    x: 80
    y: 340
    width: 1040
    height: 180
  block_position: "top-left"
  line_alignment: "left"
  overflow: "shrink"
  min_size: 16.0
  line_height: 1.4
  letter_spacing: 0

overlay:
  visible: true
  image: "assets/logo.png"
  placement:
    x: 950
    y: 50
    width: 200
    height: 100
  fit: "contain"
  opacity: 0.9
```

### Front Matter Overrides

Override any configuration setting per article:

```yaml
---
title: "Article Title"
description: "Article description text"
date: 2024-01-01
author: "John Doe"
ogp:
  title:
    visible: true
    content: "Custom: {{.Title | upper}}"
    size: 80
    color: "#FF0000"
    area:
      x: 50
      y: 50
      width: 1100
      height: 280
  description:
    visible: true
    content: "{{.Description}} - by {{.Fields.author}}"
    size: 24
    color: "#0066CC"
    area:
      x: 50
      y: 350
      width: 1100
      height: 150
  background:
    color: "#F0F0F0"
    image: "custom-bg.jpg"  # Relative to article directory
  overlay:
    visible: true
    image: "article-overlay.png"  # Relative to article directory
    placement:
      x: 900
      y: 50
      width: 250
      height: 125
    fit: "contain"
    opacity: 0.9
  output:
    filename: "custom-{{.Title}}.{{.Format}}"
---
```

**Note:** Image paths in front matter are resolved relative to the article directory first.
If not found, they fall back to paths relative to the config directory.

## Template Functions

Content templates support Hugo-compatible functions:

### Text Functions
- `default`: `{{.Title | default "Default Value"}}`
- `upper`: `{{.Title | upper}}`
- `lower`: `{{.Title | lower}}`
- `title`: `{{.Title | title}}`
- `replace`: `{{.Title | replace "old" "new"}}`
- `trim`: `{{.Description | trim}}`

### Date Functions
- `dateFormat`: `{{dateFormat "2006-01-02" .Date}}`
- `now`: `{{dateFormat "2006" now}}`

### Filename Template Functions
- `slugify`: `{{.Title | slugify}}` - Convert text to URL-friendly format

### Available Template Data
- `.Title`: Article title
- `.Description`: Article description
- `.Date`: Article date (parsed as time.Time)
- `.URL`: Custom URL from front matter
- `.RelPath`: Relative path from content directory
- `.Format`: Output format (png, jpg)
- `.Fields`: All front matter fields (access with `.Fields.fieldname`)

**Examples:**
- `.Fields.title` - Access the title field directly from front matter
- `.Fields.author` - Custom author field
- `.Fields.category` - Article category
- `.Fields.tags` - Article tags (array)

## Path Resolution

### Asset Paths
- **Config file**: Paths are relative to the config file directory
- **Front matter**: Paths are relative to the article directory first, then fall back to config directory
- **Font auto-detection**: When font path is empty, the system automatically detects Japanese fonts

### Output Paths
- **Default**: `{output.directory}/{article-path}/ogp.{format}`
- **With custom URL**: `{output.directory}/{custom-url}/ogp.{format}`
- **With filename template**: `{output.directory}/{article-path}/{generated-filename}`

## Requirements

- Go 1.18 or later
- Hugo static site with content directory structure
- Optional: Japanese fonts for Japanese text support (auto-detected when available)
