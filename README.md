# Hugo OGP Image Generator

A Go application that automatically generates Open Graph Protocol (OGP) images for Hugo static sites.
It processes Hugo content with front matter and creates custom social media preview images with advanced text rendering, background processing, and overlay composition.

## Installation

Download the latest release binary for your platform (Linux/macOS, amd64/arm64) as `./ogp`:
```bash
curl -fsSL https://raw.githubusercontent.com/yuzneri/HugoOGPImageGenerator/main/install.sh | sh
```

Set `INSTALL_DIR` to change the destination and `VERSION` (e.g. `v19`) to pin a release:
```bash
curl -fsSL https://raw.githubusercontent.com/yuzneri/HugoOGPImageGenerator/main/install.sh | INSTALL_DIR=./bin VERSION=v19 sh
```

For Windows, download the zip from the [Releases](https://github.com/yuzneri/HugoOGPImageGenerator/releases/latest) page.

## Usage

### Generate OGP for all articles
```bash
./ogp /path/to/hugo/project
```

### Generate for a specific article
```bash
./ogp --single /path/to/hugo/project "posts/my-bundle"
./ogp --single /path/to/hugo/project "posts/my-post.md"
```

### Test mode
```bash
./ogp --test "/path/to/hugo/project/content/posts/my-bundle"
./ogp --test "/path/to/hugo/project/content/posts/my-post.md"
```

An article can be given as a page bundle directory (`index.md`, or `_index.md` for a section) or as a content file.

### List all articles with OGP settings
```bash
./ogp --list /path/to/hugo/project
```

## Configuration

The application uses a 4-level configuration hierarchy with comprehensive default values:

1. **Default Configuration** - Built-in sensible defaults
2. **Global Configuration** - `config.yaml` file settings
3. **Type-Specific Configuration** - Content type-based settings
4. **Front Matter Overrides** - Per-article customizations

Each level can override any setting from the previous levels, allowing for fine-grained control while maintaining consistent defaults.

### Default Values

When no `config.yaml` file is provided, or when specific options are omitted, the following defaults are used:

```yaml
kinds: [home, section, page]

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

#### Page Selection
```yaml
kinds: [home, section, page]   # Page kinds to generate images for
```

| Kind | Content file |
|---|---|
| `home` | `content/_index.md` |
| `section` | Other `_index.md` (sections, taxonomies and terms) |
| `page` | Page bundles (`index.md`) and single-file pages |

`kinds` can be set in the global config and in type-specific configs; a type config replaces the global list.
An empty list (`kinds: []`) generates nothing for that type, and an unknown kind is an error.

To skip a single page, set `enabled: false` in its front matter (front matter only):
```yaml
ogp:
  enabled: false
```

Pages with `draft: true` are always skipped, because Hugo does not publish drafts by default and an image in `public/` would expose an unpublished page.

Skipped pages are reported as `Skipped: <path> (<reason>)`, also when given to `--single` or `--test`.

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

### Type-Specific Configuration

Create `{type}.yaml` files to apply consistent styling by content type. Content types are determined by:

1. **Front matter `type` field** (priority) 
2. **Directory path** (`content/posts/` → type "posts")
3. **Default** ("page")

**Example: posts.yaml**
```yaml
title:
  size: 64
  color: "#1A202C"

description:
  visible: true
  content: "{{.Description}} - {{.Fields.author | default \"Blog\"}}"
  size: 28

overlay:
  visible: false
```

Type configurations support partial overrides - only specify changed values from defaults.

### Front Matter Overrides

Override any configuration setting per article:

```yaml
---
title: "Article Title"
description: "Article description text"
date: 2024-01-01
author: "John Doe"
ogp:
  enabled: true              # Set to false to skip this page
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
    image: "custom-bg.jpg"
  overlay:
    visible: true
    image: "article-overlay.png"
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

All assets (images, fonts) are resolved using a simple priority system:

1. **Article directory first** - `{article-directory}/{asset-path}`
2. **Config directory second** - `{config-directory}/{asset-path}`

**Examples:**
- `cover.jpg` → `content/books/book1/cover.jpg` then `config/cover.jpg`
- `images/logo.png` → `content/posts/article/images/logo.png` then `config/images/logo.png`
- Absolute paths are used directly

This applies to all asset references regardless of where they're defined (global config, type config, or front matter). 
Japanese fonts are auto-detected when no font path is specified.

### Target Pages
All `.md` pages under `content/` are processed:

| Content file | Output (default) |
|---|---|
| `content/_index.md` | `public/ogp.png` |
| `content/posts/_index.md` | `public/posts/ogp.png` |
| `content/posts/my-bundle/index.md` | `public/posts/my-bundle/ogp.png` |
| `content/posts/my-post.md` | `public/posts/my-post/ogp.png` |

- `.md` files inside a leaf bundle (other than `index.md`) are page resources and are skipped.
- Pages without YAML front matter (no front matter, or TOML/JSON front matter) are skipped with a warning.
- Drafts (`draft: true`) and pages excluded by `kinds` or `enabled: false` are skipped (see [Page Selection](#page-selection)).
- For a single-file page, assets are resolved from the directory that contains the file.
- A section page (`_index.md`) gets the type of its top-level section, like any other page, so `content/posts/_index.md` also uses `posts.yaml`. Override it with `ogp:` in the section's front matter, or exclude sections with `kinds` (see [Page Selection](#page-selection)).
- `slug`, lowercased paths (`disablePathToLower`) and other content formats or extensions (`.markdown`, `.html`, ...) are not taken into account.

### Output Paths
- **Default**: `{output.directory}/{article-path}/ogp.{format}`
- **With custom URL**: `{output.directory}/{custom-url}/ogp.{format}`
- **With filename template**: `{output.directory}/{article-path}/{generated-filename}`

## Requirements

- Go 1.18 or later
- Hugo static site with content directory structure
- Optional: Japanese fonts for Japanese text support (auto-detected when available)
