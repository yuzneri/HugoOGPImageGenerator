# Hugo OGP Image Generator

A Go application that automatically generates Open Graph Protocol (OGP) images for Hugo static sites.
It processes Hugo content with front matter and creates custom social media preview images with advanced text rendering, background processing, and overlay composition.

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
  color: "#FFFFFF"
  # image: null

output:
  directory: "public"
  format: "png"
  filename: "ogp"

# Title text configuration
title:
  visible: true
  # content: null
  # font: null
  size: 64
  color: "#000000"
  area:
    x: 100
    y: 50
    width: 1000
    height: 250
  block_position: "middle-center"
  line_alignment: "center"
  overflow: "shrink"
  min_size: 24.0
  line_height: 1.2
  letter_spacing: 1
  line_breaking:
    start_prohibited: ".)}]>!?、。，．！？)）］｝〉》」』ー～ぁぃぅぇぉっゃゅょゎァィゥェォッャュョヮヵヶ々"
    end_prohibited: "({[<（［｛〈《「『"

# Description text configuration
description:
  visible: false
  # content: null
  # font: null
  size: 32
  color: "#666666"
  area:
    x: 100
    y: 350
    width: 1000
    height: 200
  block_position: "top-left"
  line_alignment: "left"
  overflow: "clip"
  min_size: 16.0
  line_height: 1.2
  letter_spacing: 0
  line_breaking:
    start_prohibited: ".)}]>!?、。，．！？)）］｝〉》」』ー～ぁぃぅぇぉっゃゅょゎァィゥェォッャュョヮヵヶ々"
    end_prohibited: "({[<（［｛〈《「『"

# Overlay configuration
overlay:
  visible: false
  # image: null
  placement:
    x: 50
    y: 50
  # width: null
  # height: null
  fit: "contain"
  opacity: 1.0
```

### Configuration Options

#### Background Settings
```yaml
background:
  color: "#FFFFFF"           # Background color (hex format)
  image: "path/to/image.jpg" # Background image path (optional)
```

#### Output Settings
```yaml
output:
  directory: "public"                       # Output directory
  format: "png"                             # Image format (png, jpg)
  filename: "custom-{{.Title}}.{{.Format}}" # Filename template (optional)
```

#### Title Text Configuration
```yaml
title:
  visible: true                                # Show title text
  content: "{{.Title}}"                        # Content template (optional, uses article title)
  font: "fonts/custom.ttf"                     # Font file path (optional, auto-detects if omitted)
  size: 72                                     # Font size in pixels
  color: "#000000"                             # Text color (hex format)
  area:                                        # Text rendering area
    x: 50                                      # X position
    y: 50                                      # Y position
    width: 1200                                # Area width
    height: 300                                # Area height
  block_position: "middle-center"              # Text block position
  line_alignment: "center"                     # Line alignment within block
  overflow: "shrink"                           # Overflow handling method
  min_size: 24.0                               # Minimum font size for shrink mode
  line_height: 1.2                             # Line height multiplier
  letter_spacing: 1                            # Letter spacing in pixels
  line_breaking:                               # Japanese line breaking rules
    start_prohibited: "、。！？」』）"           # Characters that cannot start a line
    end_prohibited: "「『（"                    # Characters that cannot end a line
```

#### Description Text Configuration
```yaml
description:
  visible: true                                # Show description text
  content: "{{.Description}}"                  # Content template (optional, uses article description)
  font: "fonts/description.ttf"                # Font file path (optional, can differ from title)
  size: 32                                     # Font size in pixels
  color: "#666666"                             # Text color (hex format)
  area:                                        # Text rendering area
    x: 50                                      # X position
    y: 380                                     # Y position
    width: 1200                                # Area width
    height: 200                                # Area height
  block_position: "top-left"                   # Text block position
  line_alignment: "left"                       # Line alignment within block
  overflow: "shrink"                           # Overflow handling method
  min_size: 16.0                               # Minimum font size for shrink mode
  line_height: 1.4                             # Line height multiplier
  letter_spacing: 0                            # Letter spacing in pixels
  line_breaking:                               # Japanese line breaking rules
    start_prohibited: "、。！？」』）"           # Characters that cannot start a line
    end_prohibited: "「『（"                    # Characters that cannot end a line
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
  visible: true             # Show overlay image
  image: "overlay.png"      # Overlay image path
  placement:                # Image positioning
    x: 100                  # X position
    y: 100                  # Y position
    width: 300              # Image width (optional, auto-detects if omitted)
    height: 200             # Image height (optional, auto-detects if omitted)
  fit: "cover"              # Image fit method
  opacity: 0.8              # Image opacity (0.0-1.0)
```

**Placement Width/Height Behavior:**

The final image size depends on both the `placement` dimensions and the `fit` option:

**Step 1 - Dimension Calculation:**
- **Both width/height specified**: Uses specified dimensions as target
- **Width only**: Auto-calculates height maintaining aspect ratio
- **Height only**: Auto-calculates width maintaining aspect ratio  
- **Neither specified**: Uses original image dimensions

**Step 2 - Fit Processing:**
- **`cover`**: Scale to cover target area completely (may crop excess)
- **`contain`**: Scale to fit within target area (preserves aspect ratio)
- **`fill`**: Stretch to exact target dimensions (may distort)
- **`none`**: No resizing, uses original dimensions

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
- **Font auto-detection**: When font path is omitted, the system automatically detects Japanese fonts

### Output Paths
- **Default**: `{output.directory}/{article-path}/ogp.{format}`
- **With custom URL**: `{output.directory}/{custom-url}/ogp.{format}`
- **With filename template**: `{output.directory}/{article-path}/{generated-filename}`

## Requirements

- Go 1.18 or later
- Hugo static site with content directory structure
- Optional: Japanese fonts for Japanese text support (auto-detected when available)
