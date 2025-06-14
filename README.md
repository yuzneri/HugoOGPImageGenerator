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

The application uses a 4-level configuration hierarchy with comprehensive default values:

1. **Default Configuration** - Built-in sensible defaults
2. **Global Configuration** - `config.yaml` file settings
3. **Type-Specific Configuration** - Content type-based settings
4. **Front Matter Overrides** - Per-article customizations

Each level can override any setting from the previous levels, allowing for fine-grained control while maintaining consistent defaults.

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

### Type-Specific Configuration

Create configuration files for different content types to apply consistent styling across all content of that type. The content type is automatically detected from the Hugo directory structure.

#### Content Type Detection

The content type is determined using Hugo's standard type detection with the following priority:

1. **Hugo's native `type` field in front matter** (highest priority)
2. **Directory path within Hugo's `content/` directory** (fallback)
3. **Default to `"page"`** (if neither above applies)

**Front matter type field example:**
```yaml
---
title: "Special Article"
description: "Custom content type"
type: "featured"        # ← Explicit type override
---
```

**Directory-based type detection:**
```
content/
├── about.md           → type: "page"  
├── posts/
│   └── article1.md    → type: "posts"
├── books/
│   └── book1/
│       └── index.md   → type: "books"
└── tutorials/
    └── tutorial1.md   → type: "tutorials"
```

The front matter `type` field takes precedence over directory structure, allowing you to override the default directory-based type for specific articles.

#### Type Configuration Files

Create type-specific configuration files in your config directory using the pattern `{type}.yaml`:

**books.yaml** - Configuration for all book content:
```yaml
title:
  visible: true
  size: 72
  color: "#2C3E50"
  area:
    x: 460
    y: 100
    width: 710
    height: 370
  block_position: "top-center"
  line_alignment: "center"

description:
  visible: false

overlay:
  visible: true
  image: "cover.jpg"           # Look for cover.jpg with priority: 1) article directory 2) config directory
  placement:
    x: 50
    y: 50
    height: 580                # Fixed height, auto-width
  fit: "contain"
  opacity: 1.0

background:
  color: "#F8F9FA"
```

**posts.yaml** - Configuration for blog posts:
```yaml
title:
  visible: true
  size: 64
  color: "#1A202C"
  area:
    x: 100
    y: 80
    width: 1000
    height: 200

description:
  visible: true
  content: "{{.Description}} - {{.Fields.author | default \"Blog\"}}"
  size: 28
  color: "#718096"
  area:
    x: 100
    y: 300
    width: 1000
    height: 150

overlay:
  visible: false

background:
  color: "#FFFFFF"
```

**tutorials.yaml** - Configuration for tutorial content:
```yaml
title:
  visible: true
  content: "Tutorial: {{.Title}}"
  size: 56
  color: "#2B6CB0"
  
description:
  visible: true
  content: "Step-by-step guide"
  size: 24
  color: "#4A5568"

overlay:
  visible: true
  image: "tutorial-badge.png"   # Common badge for all tutorials
  placement:
    x: 950
    y: 50
    width: 200
    height: 100
  fit: "contain"
```

#### Partial Configuration Support

Type configurations support partial settings - you only need to specify the values you want to change from the defaults:

**minimal-books.yaml**:
```yaml
# Only override specific values
title:
  color: "#8B5CF6"            # Purple title for books
  area:
    x: 400                    # Move title to the right
    
overlay:
  visible: true
  placement:
    height: 500              # Only set height, X/Y/width inherited from defaults
```

This partial configuration will:
- Change title color to purple
- Move title area X position to 400 (Y, width, height remain as defaults)
- Enable overlay visibility  
- Set overlay height to 500 (X, Y, width remain as defaults)

#### Type Configuration Examples by Use Case

**E-commerce Products**:
```yaml
# product.yaml
title:
  content: "{{.Title | upper}}"
  color: "#E53E3E"
  
description:
  content: "Price: {{.Fields.price | default \"Contact for pricing\"}}"
  
overlay:
  image: "product-badge.png"
```

**Documentation Pages**:
```yaml
# docs.yaml  
title:
  content: "📖 {{.Title}}"
  color: "#2D3748"
  
description:
  content: "{{.Fields.category | default \"Documentation\"}} | {{.Description}}"
  
background:
  color: "#F7FAFC"
```

**Event Announcements**:
```yaml
# events.yaml
title:
  content: "🎉 {{.Title}}"
  size: 68
  color: "#9F7AEA"
  
description:
  content: "{{dateFormat \"January 2, 2006\" .Date}} | {{.Description}}"
  
background:
  image: "event-bg.jpg"
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

The system uses a 2-tier priority system for finding asset files (images, fonts):

1. **Article directory** (highest priority) - `{article-directory}/{asset-path}`
2. **Config directory** (fallback) - `{config-directory}/{asset-path}`

**Priority examples:**
- For `cover.jpg`: First check `content/books/book1/cover.jpg`, then `config/cover.jpg`
- For `images/logo.png`: First check `content/posts/article/images/logo.png`, then `config/images/logo.png`
- **Absolute paths**: Used directly without any resolution

**Unified behavior for all contexts:**
- **All asset types** (background images, fonts, overlay images): Use the 2-tier priority system
- **Type/Global config files**: Assets search article directory first, then config directory
- **Front matter overrides**: Assets use the same 2-tier system (article → config)
- **Font auto-detection**: When font path is omitted, the system automatically detects Japanese fonts

### Output Paths
- **Default**: `{output.directory}/{article-path}/ogp.{format}`
- **With custom URL**: `{output.directory}/{custom-url}/ogp.{format}`
- **With filename template**: `{output.directory}/{article-path}/{generated-filename}`

## Requirements

- Go 1.18 or later
- Hugo static site with content directory structure
- Optional: Japanese fonts for Japanese text support (auto-detected when available)
