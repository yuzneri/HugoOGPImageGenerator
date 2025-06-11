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

	newConfig := *baseConfig

	cm.mergeTextConfig(&newConfig, ogpFM)
	cm.mergeBackgroundConfig(&newConfig, ogpFM)
	cm.mergeOutputConfig(&newConfig, ogpFM)
	cm.mergeOverlayConfig(&newConfig, ogpFM)

	return &newConfig
}

func (cm *ConfigMerger) mergeTextConfig(config *Config, ogpFM *OGPFrontMatter) {
	if ogpFM.Text == nil {
		return
	}

	text := ogpFM.Text
	if text.Content != nil {
		config.Text.Content = text.Content
	}
	cm.mergeStringPtr(&config.Text.Font, text.Font)
	cm.mergeFloat64Ptr(&config.Text.Size, text.Size)
	cm.mergeStringPtr(&config.Text.Color, text.Color)
	cm.mergeStringPtr(&config.Text.BlockPosition, text.BlockPosition)
	cm.mergeStringPtr(&config.Text.LineAlignment, text.LineAlignment)
	cm.mergeStringPtr(&config.Text.Overflow, text.Overflow)
	cm.mergeFloat64Ptr(&config.Text.MinSize, text.MinSize)
	cm.mergeFloat64Ptr(&config.Text.LineHeight, text.LineHeight)
	cm.mergeIntPtr(&config.Text.LetterSpacing, text.LetterSpacing)

	cm.mergeTextAreaConfig(&config.Text.Area, text.Area)
	cm.mergeLineBreakingConfig(&config.Text.LineBreaking, text.LineBreaking)
}

func (cm *ConfigMerger) mergeTextAreaConfig(area *TextArea, overrideArea *struct {
	X      *int `yaml:"x,omitempty"`
	Y      *int `yaml:"y,omitempty"`
	Width  *int `yaml:"width,omitempty"`
	Height *int `yaml:"height,omitempty"`
}) {
	if overrideArea == nil {
		return
	}

	cm.mergeIntPtr(&area.X, overrideArea.X)
	cm.mergeIntPtr(&area.Y, overrideArea.Y)
	cm.mergeIntPtr(&area.Width, overrideArea.Width)
	cm.mergeIntPtr(&area.Height, overrideArea.Height)
}

func (cm *ConfigMerger) mergeLineBreakingConfig(lineBreaking *struct {
	StartProhibited string `yaml:"start_prohibited"`
	EndProhibited   string `yaml:"end_prohibited"`
}, overrideLineBreaking *struct {
	StartProhibited *string `yaml:"start_prohibited,omitempty"`
	EndProhibited   *string `yaml:"end_prohibited,omitempty"`
}) {
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
	if bg.Image != nil {
		config.Background.Image = bg.Image
	}
	cm.mergeStringPtr(&config.Background.Color, bg.Color)
}

func (cm *ConfigMerger) mergeOutputConfig(config *Config, ogpFM *OGPFrontMatter) {
	if ogpFM.Output == nil {
		return
	}

	output := ogpFM.Output
	if output.Filename != nil {
		config.Output.Filename = output.Filename
	}
}

func (cm *ConfigMerger) mergeOverlayConfig(config *Config, ogpFM *OGPFrontMatter) {
	if ogpFM.Overlay != nil {
		config.Overlay = ogpFM.Overlay
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
