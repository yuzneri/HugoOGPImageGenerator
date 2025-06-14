package main

// ConfigMerger handles merging base configuration with front matter overrides.
type ConfigMerger struct{}

// NewConfigMerger creates a new ConfigMerger.
func NewConfigMerger() *ConfigMerger {
	return &ConfigMerger{}
}

// applySettingsToConfig applies ConfigSettings to an existing Config.
// Only explicitly set fields (non-nil) are applied.
func (cm *ConfigMerger) applySettingsToConfig(target *Config, settings *ConfigSettings) {
	if settings == nil {
		return
	}

	// Apply background settings
	if settings.Background != nil {
		cm.applyBackgroundSettings(&target.Background, settings.Background)
	}

	// Apply output settings
	if settings.Output != nil {
		cm.applyOutputSettings(&target.Output, settings.Output)
	}

	// Apply title settings
	if settings.Title != nil {
		cm.applyTextSettings(&target.Title, settings.Title)
	}

	// Apply description settings
	if settings.Description != nil {
		cm.applyTextSettings(&target.Description, settings.Description)
	}

	// Apply overlay settings
	if settings.Overlay != nil {
		cm.applyOverlaySettings(&target.Overlay, settings.Overlay)
	}
}

// applyBackgroundSettings applies BackgroundSettings to BackgroundConfig.
func (cm *ConfigMerger) applyBackgroundSettings(target *BackgroundConfig, settings *BackgroundSettings) {
	if settings.Image != nil {
		target.Image = cm.copyStringPtr(settings.Image)
	}
	if settings.Color != nil {
		target.Color = *settings.Color
	}
}

// applyOutputSettings applies OutputSettings to OutputConfig.
func (cm *ConfigMerger) applyOutputSettings(target *OutputConfig, settings *OutputSettings) {
	if settings.Directory != nil {
		target.Directory = *settings.Directory
	}
	if settings.Format != nil {
		target.Format = *settings.Format
	}
	if settings.Filename != nil {
		target.Filename = *settings.Filename
	}
}

// applyTextSettings applies TextSettings to TextConfig.
func (cm *ConfigMerger) applyTextSettings(target *TextConfig, settings *TextSettings) {
	if settings.Visible != nil {
		target.Visible = *settings.Visible
	}
	if settings.Content != nil {
		target.Content = cm.copyStringPtr(settings.Content)
	}
	if settings.Font != nil {
		target.Font = cm.copyStringPtr(settings.Font)
	}
	if settings.Size != nil {
		target.Size = *settings.Size
	}
	if settings.Color != nil {
		target.Color = *settings.Color
	}
	if settings.Area != nil {
		cm.applyTextAreaSettings(&target.Area, settings.Area)
	}
	if settings.BlockPosition != nil {
		target.BlockPosition = *settings.BlockPosition
	}
	if settings.LineAlignment != nil {
		target.LineAlignment = *settings.LineAlignment
	}
	if settings.Overflow != nil {
		target.Overflow = *settings.Overflow
	}
	if settings.MinSize != nil {
		target.MinSize = *settings.MinSize
	}
	if settings.LineHeight != nil {
		target.LineHeight = *settings.LineHeight
	}
	if settings.LetterSpacing != nil {
		target.LetterSpacing = *settings.LetterSpacing
	}
	if settings.LineBreaking != nil {
		cm.applyLineBreakingSettings(&target.LineBreaking, settings.LineBreaking)
	}
}

// applyTextAreaSettings applies TextAreaSettings to TextArea.
func (cm *ConfigMerger) applyTextAreaSettings(target *TextArea, settings *TextAreaSettings) {
	if settings.X != nil {
		target.X = *settings.X
	}
	if settings.Y != nil {
		target.Y = *settings.Y
	}
	if settings.Width != nil {
		target.Width = *settings.Width
	}
	if settings.Height != nil {
		target.Height = *settings.Height
	}
}

// applyLineBreakingSettings applies LineBreakingSettings to LineBreakingConfig.
func (cm *ConfigMerger) applyLineBreakingSettings(target *LineBreakingConfig, settings *LineBreakingSettings) {
	if settings.StartProhibited != nil {
		target.StartProhibited = *settings.StartProhibited
	}
	if settings.EndProhibited != nil {
		target.EndProhibited = *settings.EndProhibited
	}
}

// applyOverlaySettings applies OverlayConfigSettings to MainOverlayConfig.
func (cm *ConfigMerger) applyOverlaySettings(target *MainOverlayConfig, settings *OverlayConfigSettings) {
	// Apply visible setting if explicitly provided
	if settings.Visible != nil {
		target.Visible = *settings.Visible
	} else if settings.Image != nil && *settings.Image != "" {
		// UX improvement: If an image is provided but visible is not set,
		// automatically make the overlay visible. This makes it more intuitive
		// for users who configure overlay settings but forget to set visible: true
		target.Visible = true
	}

	if settings.Image != nil {
		target.Image = cm.copyStringPtr(settings.Image)
	}
	if settings.Fit != nil {
		target.Fit = *settings.Fit
	}
	if settings.Opacity != nil {
		target.Opacity = *settings.Opacity
	}
	if settings.Placement != nil {
		cm.applyPlacementSettings(&target.Placement, settings.Placement)
	}
}

// applyPlacementSettings applies PlacementSettings to PlacementConfig.
func (cm *ConfigMerger) applyPlacementSettings(target *PlacementConfig, settings *PlacementSettings) {
	if settings.X != nil {
		target.X = *settings.X
	}
	if settings.Y != nil {
		target.Y = *settings.Y
	}
	if settings.Width != nil {
		target.Width = cm.copyIntPtr(settings.Width)
	}
	if settings.Height != nil {
		target.Height = cm.copyIntPtr(settings.Height)
	}
}

// MergeConfigsWithSettings creates a new config by applying the 4-level configuration hierarchy:
// Default Config -> Global ConfigSettings -> Type ConfigSettings -> Front Matter Overrides
// It returns a new config instance without modifying the original configurations.
func (cm *ConfigMerger) MergeConfigsWithSettings(defaultConfig *Config, globalSettings, typeSettings *ConfigSettings, ogpFM *OGPFrontMatter) *Config {
	// Start with default configuration
	result := cm.deepCopyConfig(defaultConfig)

	// Apply global configuration settings if present
	if globalSettings != nil {
		cm.applySettingsToConfig(result, globalSettings)
	}

	// Apply type-specific configuration settings if present
	if typeSettings != nil {
		cm.applySettingsToConfig(result, typeSettings)
	}

	// Apply front matter overrides if present
	if ogpFM != nil {
		result = cm.applyFrontMatterOverrides(result, ogpFM)
	}

	return result
}

// MergeConfigs creates a new config by applying front matter overrides to the base config.
// It returns a new config instance without modifying the original.
func (cm *ConfigMerger) MergeConfigs(baseConfig *Config, ogpFM *OGPFrontMatter) *Config {
	if ogpFM == nil {
		return cm.deepCopyConfig(baseConfig)
	}

	return cm.applyFrontMatterOverrides(cm.deepCopyConfig(baseConfig), ogpFM)
}

// deepCopyConfig creates a deep copy of a configuration
func (cm *ConfigMerger) deepCopyConfig(config *Config) *Config {
	if config == nil {
		return nil
	}

	// Create a new config instance
	newConfig := *config

	// Deep copy all pointer fields
	cm.deepCopyPointerFields(&newConfig, config)

	return &newConfig
}

// deepCopyPointerFields performs deep copying of all pointer fields in a Config
func (cm *ConfigMerger) deepCopyPointerFields(dest, src *Config) {
	// Background pointers
	dest.Background.Image = cm.copyStringPtr(src.Background.Image)

	// Title pointers
	dest.Title.Content = cm.copyStringPtr(src.Title.Content)
	dest.Title.Font = cm.copyStringPtr(src.Title.Font)

	// Description pointers
	dest.Description.Content = cm.copyStringPtr(src.Description.Content)
	dest.Description.Font = cm.copyStringPtr(src.Description.Font)

	// Overlay pointers
	dest.Overlay.Image = cm.copyStringPtr(src.Overlay.Image)

	// Overlay placement pointers
	dest.Overlay.Placement.Width = cm.copyIntPtr(src.Overlay.Placement.Width)
	dest.Overlay.Placement.Height = cm.copyIntPtr(src.Overlay.Placement.Height)
}

// copyStringPtr creates a deep copy of a string pointer
func (cm *ConfigMerger) copyStringPtr(src *string) *string {
	if src == nil {
		return nil
	}
	copy := *src
	return &copy
}

// copyIntPtr creates a deep copy of an int pointer
func (cm *ConfigMerger) copyIntPtr(src *int) *int {
	if src == nil {
		return nil
	}
	copy := *src
	return &copy
}

// applyFrontMatterOverrides applies front matter overrides to a config
func (cm *ConfigMerger) applyFrontMatterOverrides(config *Config, ogpFM *OGPFrontMatter) *Config {
	result := cm.deepCopyConfig(config)

	// Merge title and description configurations separately
	if ogpFM.Title != nil {
		cm.mergeTextConfigOverride(&result.Title, ogpFM.Title)
	}
	if ogpFM.Description != nil {
		cm.mergeTextConfigOverride(&result.Description, ogpFM.Description)
	}

	cm.mergeBackgroundConfig(result, ogpFM)
	cm.mergeOutputConfig(result, ogpFM)
	cm.mergeOverlayConfig(result, ogpFM)

	return result
}

// mergeTextConfigOverride applies front matter overrides to text config
// Simple rule: if override field is not nil, use the override value
func (cm *ConfigMerger) mergeTextConfigOverride(config *TextConfig, override *TextConfigOverride) {
	if override.Visible != nil {
		config.Visible = *override.Visible
	}
	if override.Content != nil {
		config.Content = override.Content
	}
	if override.Font != nil {
		config.Font = override.Font
	}
	if override.Size != nil {
		config.Size = *override.Size
	}
	if override.Color != nil {
		config.Color = *override.Color
	}
	if override.BlockPosition != nil {
		config.BlockPosition = *override.BlockPosition
	}
	if override.LineAlignment != nil {
		config.LineAlignment = *override.LineAlignment
	}
	if override.Overflow != nil {
		config.Overflow = *override.Overflow
	}
	if override.MinSize != nil {
		config.MinSize = *override.MinSize
	}
	if override.LineHeight != nil {
		config.LineHeight = *override.LineHeight
	}
	if override.LetterSpacing != nil {
		config.LetterSpacing = *override.LetterSpacing
	}
	if override.Area != nil {
		cm.mergeTextAreaConfig(&config.Area, override.Area)
	}
	if override.LineBreaking != nil {
		cm.mergeLineBreakingConfig(&config.LineBreaking, override.LineBreaking)
	}
}

// mergeTextAreaConfig applies area overrides to text area
// Simple rule: if override field is not nil, use the override value
func (cm *ConfigMerger) mergeTextAreaConfig(area *TextArea, overrideArea *TextAreaConfig) {
	if overrideArea == nil {
		return
	}
	if overrideArea.X != nil {
		area.X = *overrideArea.X
	}
	if overrideArea.Y != nil {
		area.Y = *overrideArea.Y
	}
	if overrideArea.Width != nil {
		area.Width = *overrideArea.Width
	}
	if overrideArea.Height != nil {
		area.Height = *overrideArea.Height
	}
}

// mergeLineBreakingConfig applies line breaking overrides
func (cm *ConfigMerger) mergeLineBreakingConfig(lineBreaking *LineBreakingConfig, overrideLineBreaking *LineBreakingOverride) {
	if overrideLineBreaking == nil {
		return
	}
	if overrideLineBreaking.StartProhibited != nil {
		lineBreaking.StartProhibited = *overrideLineBreaking.StartProhibited
	}
	if overrideLineBreaking.EndProhibited != nil {
		lineBreaking.EndProhibited = *overrideLineBreaking.EndProhibited
	}
}

// mergeBackgroundConfig applies background overrides
func (cm *ConfigMerger) mergeBackgroundConfig(config *Config, ogpFM *OGPFrontMatter) {
	if ogpFM.Background == nil {
		return
	}
	if ogpFM.Background.Image != nil {
		config.Background.Image = ogpFM.Background.Image
	}
	if ogpFM.Background.Color != nil {
		config.Background.Color = *ogpFM.Background.Color
	}
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

// mergeOverlayConfig applies overlay overrides
func (cm *ConfigMerger) mergeOverlayConfig(config *Config, ogpFM *OGPFrontMatter) {
	if ogpFM.Overlay == nil {
		return
	}

	overlay := ogpFM.Overlay
	if overlay.Visible != nil {
		config.Overlay.Visible = *overlay.Visible
	}
	if overlay.Image != nil {
		config.Overlay.Image = overlay.Image
	}
	if overlay.Fit != nil {
		config.Overlay.Fit = *overlay.Fit
	}
	if overlay.Opacity != nil {
		config.Overlay.Opacity = *overlay.Opacity
	}
	if overlay.Placement != nil {
		cm.applyPlacementSettings(&config.Overlay.Placement, overlay.Placement)
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
