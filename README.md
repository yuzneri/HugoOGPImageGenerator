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

text:
  # content: null      # No default content template (uses article title)
  font: ""             # Auto-detect system fonts
  size: 64             # Font size in pixels
  color: "#000000"     # Black text
  area:
    x: 100             # Text area X position
    y: 100             # Text area Y position  
    width: 1000        # Text area width
    height: 430        # Text area height
  block_position: "middle-center" # Text block position in area
  line_alignment: "left"          # Individual line alignment
  overflow: "shrink"              # Shrink text to fit area
  min_size: 12.0                  # Minimum font size when shrinking
  line_height: 1.2                # Line height multiplier
  letter_spacing: 1               # Letter spacing in pixels
  line_breaking:
    start_prohibited: ".)}]>!?、。，．！？)）］｝〉》」』ー～ぁぃぅぇぉっゃゅょゎァィゥェォッャュョヮヵヶ々"  # Japanese line-start prohibited chars
    end_prohibited: "({[<（［｛〈《「『"                                                            # Japanese line-end prohibited chars

# overlay: null        # No overlay by default
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

#### Text Configuration
```yaml
text:
  content: "{{.Title | default \"Default\"}}"  # Content template (optional)
  font: "fonts/custom.ttf"                     # Font file path (empty = auto-detect)
  size: 72                                     # Font size in pixels
  color: "#FF0000"                             # Text color (hex)
  area:                                        # Text rendering area
    x: 50                                      # X position
    y: 50                                      # Y position
    width: 1200                                # Area width
    height: 500                                # Area height
  block_position: "top-left"                  # Text block position
  line_alignment: "center"                     # Line alignment within block
  overflow: "clip"                             # Overflow handling
  min_size: 8.0                                # Minimum font size for shrink mode
  line_height: 1.5                             # Line height multiplier
  letter_spacing: 2                            # Letter spacing in pixels
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

### Global Configuration (`config.yaml`)

Create a `config.yaml` file to override defaults:

```yaml
background:
  color: "#F5F5F5"
  image: "assets/bg.jpg"

output:
  directory: "static/images"
  format: "jpg"
  filename: "ogp-{{.Title | slugify}}.{{.Format}}"

text:
  content: "{{.Title | upper}} - {{.Fields.author | default \"Blog\"}}"
  font: "fonts/NotoSansCJK-Bold.ttc"
  size: 48
  color: "#2C3E50"
  area:
    x: 80
    y: 120
    width: 1040
    height: 400
  block_position: "middle-left"
  line_alignment: "left"
  overflow: "shrink"
  min_size: 20.0
  line_height: 1.4
  letter_spacing: 0

overlay:
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
date: 2024-01-01
ogp:
  text:
    content: "Custom: {{.Title | upper}} - {{dateFormat \"2006年1月2日\" .Date}}"
    size: 72
    color: "#FF0000"
    area:
      x: 50
      y: 50
  background:
    color: "#F0F0F0"
    image: "custom-bg.jpg"
  overlay:
    image: "article-overlay.png"
  output:
    filename: "custom-{{.Title}}.{{.Format}}"
---
```

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

### Available Template Data
- `.Title`: Article title
- `.Description`: Article description
- `.Date`: Article date (parsed as time.Time)
- `.URL`: Custom URL from front matter
- `.Fields`: All front matter fields (access with `.Fields.fieldname`)

**Examples:**
- `.Fields.title` - Access the title field directly from front matter
- `.Fields.author` - Custom author field
- `.Fields.category` - Article category
- `.Fields.tags` - Article tags (array)
