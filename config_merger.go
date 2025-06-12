package main

// ConfigMerger handles merging base configuration with front matter overrides.
type ConfigMerger struct{}

// NewConfigMerger creates a new ConfigMerger.
func NewConfigMerger() *ConfigMerger {
	return &ConfigMerger{}
}

// MergeConfigs creates a new config by applying front matter overrides to the base config.
// It returns a new config instance without modifying the original.
func (cm *ConfigMerger) MergeConfigs(baseConfig *Config, ogpFM *OGPFrontMatter) *Config {
	if ogpFM == nil {
		return baseConfig
	}

	// Deep copy the base config to avoid modifying the original
	newConfig := *baseConfig

	// Deep copy pointer fields in base config to avoid reference sharing
	if baseConfig.Background.Image != nil {
		imageCopy := *baseConfig.Background.Image
		newConfig.Background.Image = &imageCopy
	}

	// Filename is now a value type, automatically copied

	// Deep copy title content if it exists
	if baseConfig.Title.Content != nil {
		contentCopy := *baseConfig.Title.Content
		newConfig.Title.Content = &contentCopy
	}

	// Deep copy title font if it exists
	if baseConfig.Title.Font != nil {
		fontCopy := *baseConfig.Title.Font
		newConfig.Title.Font = &fontCopy
	}

	// Deep copy description content if it exists
	if baseConfig.Description.Content != nil {
		contentCopy := *baseConfig.Description.Content
		newConfig.Description.Content = &contentCopy
	}

	// Deep copy description font if it exists
	if baseConfig.Description.Font != nil {
		fontCopy := *baseConfig.Description.Font
		newConfig.Description.Font = &fontCopy
	}

	// Deep copy overlay config (value type)
	newConfig.Overlay = baseConfig.Overlay

	// Deep copy overlay pointer fields
	if baseConfig.Overlay.Image != nil {
		imageCopy := *baseConfig.Overlay.Image
		newConfig.Overlay.Image = &imageCopy
	}
	// Fit is now a value type, copy directly
	newConfig.Overlay.Fit = baseConfig.Overlay.Fit
	// Opacity is now a value type, copy directly
	newConfig.Overlay.Opacity = baseConfig.Overlay.Opacity

	// Deep copy placement (it's a value type in MainOverlayConfig)
	newConfig.Overlay.Placement = baseConfig.Overlay.Placement

	// Merge title and description configurations separately
	if ogpFM.Title != nil {
		cm.mergeTextConfigOverride(&newConfig.Title, ogpFM.Title)
	}
	if ogpFM.Description != nil {
		cm.mergeTextConfigOverride(&newConfig.Description, ogpFM.Description)
	}

	cm.mergeBackgroundConfig(&newConfig, ogpFM)
	cm.mergeOutputConfig(&newConfig, ogpFM)
	cm.mergeOverlayConfig(&newConfig, ogpFM)

	return &newConfig
}

func (cm *ConfigMerger) mergeTextConfigOverride(config *TextConfig, override *TextConfigOverride) {
	cm.mergeBoolPtr(&config.Visible, override.Visible)
	cm.mergeStringPtrValue(&config.Content, override.Content)
	cm.mergeStringPtrValue(&config.Font, override.Font)
	cm.mergeFloat64Ptr(&config.Size, override.Size)
	cm.mergeStringPtr(&config.Color, override.Color)
	cm.mergeStringPtr(&config.BlockPosition, override.BlockPosition)
	cm.mergeStringPtr(&config.LineAlignment, override.LineAlignment)
	cm.mergeStringPtr(&config.Overflow, override.Overflow)
	cm.mergeFloat64Ptr(&config.MinSize, override.MinSize)
	cm.mergeFloat64Ptr(&config.LineHeight, override.LineHeight)
	cm.mergeIntPtr(&config.LetterSpacing, override.LetterSpacing)

	cm.mergeTextAreaConfig(&config.Area, override.Area)
	cm.mergeLineBreakingConfig(&config.LineBreaking, override.LineBreaking)
}

func (cm *ConfigMerger) mergeTextAreaConfig(area *TextArea, overrideArea *TextAreaConfig) {
	if overrideArea == nil {
		return
	}

	cm.mergeIntPtr(&area.X, overrideArea.X)
	cm.mergeIntPtr(&area.Y, overrideArea.Y)
	cm.mergeIntPtr(&area.Width, overrideArea.Width)
	cm.mergeIntPtr(&area.Height, overrideArea.Height)
}

func (cm *ConfigMerger) mergeLineBreakingConfig(lineBreaking *LineBreakingConfig, overrideLineBreaking *LineBreakingOverride) {
	if overrideLineBreaking == nil {
		return
	}

	cm.mergeStringPtr(&lineBreaking.StartProhibited, overrideLineBreaking.StartProhibited)
	cm.mergeStringPtr(&lineBreaking.EndProhibited, overrideLineBreaking.EndProhibited)
}

func (cm *ConfigMerger) mergeBackgroundConfig(config *Config, ogpFM *OGPFrontMatter) {
	if ogpFM.Background == nil {
		return
	}

	bg := ogpFM.Background
	cm.mergeStringPtrValue(&config.Background.Image, bg.Image)
	cm.mergeStringPtr(&config.Background.Color, bg.Color)
}

func (cm *ConfigMerger) mergeOutputConfig(config *Config, ogpFM *OGPFrontMatter) {
	if ogpFM.Output == nil {
		return
	}

	output := ogpFM.Output
	if output.Filename != nil {
		config.Output.Filename = *output.Filename
	}
}

func (cm *ConfigMerger) mergeOverlayConfig(config *Config, ogpFM *OGPFrontMatter) {
	if ogpFM.Overlay == nil {
		return
	}

	// Merge overlay configuration (config.Overlay is always a value, never nil)
	overlay := ogpFM.Overlay
	cm.mergeBoolPtr(&config.Overlay.Visible, overlay.Visible)
	cm.mergeStringPtrValue(&config.Overlay.Image, overlay.Image)
	cm.mergeStringPtr(&config.Overlay.Fit, overlay.Fit)
	cm.mergeFloat64Ptr(&config.Overlay.Opacity, overlay.Opacity)

	// Merge placement configuration
	if overlay.Placement != nil {
		// Directly merge values from article placement to main config placement
		if overlay.Placement.X != 0 {
			config.Overlay.Placement.X = overlay.Placement.X
		}
		if overlay.Placement.Y != 0 {
			config.Overlay.Placement.Y = overlay.Placement.Y
		}
		if overlay.Placement.Width != nil {
			config.Overlay.Placement.Width = overlay.Placement.Width
		}
		if overlay.Placement.Height != nil {
			config.Overlay.Placement.Height = overlay.Placement.Height
		}
	}
}

func (cm *ConfigMerger) mergeStringPtr(target *string, source *string) {
	if source != nil {
		*target = *source
	}
}

func (cm *ConfigMerger) mergeFloat64Ptr(target *float64, source *float64) {
	if source != nil {
		*target = *source
	}
}

func (cm *ConfigMerger) mergeIntPtr(target *int, source *int) {
	if source != nil {
		*target = *source
	}
}

func (cm *ConfigMerger) mergeBoolPtr(target *bool, source *bool) {
	if source != nil {
		*target = *source
	}
}

// mergeStringPtrValue copies the value of source pointer to target pointer (avoiding reference sharing)
func (cm *ConfigMerger) mergeStringPtrValue(target **string, source *string) {
	if source != nil {
		value := *source
		*target = &value
	}
}

// mergeFloat64PtrValue copies the value of source pointer to target pointer (avoiding reference sharing)
func (cm *ConfigMerger) mergeFloat64PtrValue(target **float64, source *float64) {
	if source != nil {
		value := *source
		*target = &value
	}
}

// mergeIntPtrValue copies the value of source pointer to target pointer (avoiding reference sharing)
func (cm *ConfigMerger) mergeIntPtrValue(target **int, source *int) {
	if source != nil {
		value := *source
		*target = &value
	}
}
