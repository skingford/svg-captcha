package captcha

import "fmt"

// ColorManager handles color selection for captcha elements
type ColorManager struct {
	enableColor bool
	background  string
	textColors  []string
	noiseColors []string
}

// NewColorManager creates a new color manager
func NewColorManager(config *Config) *ColorManager {
	textColors := []string{
		"#1a1a1a", "#2c3e50", "#34495e", "#7f8c8d",
		"#c0392b", "#e74c3c", "#d35400", "#e67e22",
		"#16a085", "#27ae60", "#2980b9", "#8e44ad",
	}

	noiseColors := []string{
		"#bdc3c7", "#95a5a6", "#ecf0f1", "#d5dbdb",
		"#f8c471", "#f7dc6f", "#aed6f1", "#a9dfbf",
	}

	if !config.Color {
		// Use grayscale colors only
		textColors = []string{
			"#1a1a1a", "#2c2c2c", "#3f3f3f", "#525252",
			"#666666", "#7a7a7a", "#8d8d8d", "#a0a0a0",
		}
		noiseColors = []string{
			"#bdc3c7", "#d5d5d5", "#e8e8e8", "#f2f2f2",
		}
	}

	return &ColorManager{
		enableColor: config.Color,
		background:  config.Background,
		textColors:  textColors,
		noiseColors: noiseColors,
	}
}

// GetRandomTextColor returns a random color suitable for text
func (cm *ColorManager) GetRandomTextColor() string {
	if len(cm.textColors) == 0 {
		return "#000000" // fallback
	}

	index, err := secureRandomInt(len(cm.textColors))
	if err != nil {
		return cm.textColors[0] // fallback to first color
	}

	return cm.textColors[index]
}

// GetRandomNoiseColor returns a random color suitable for noise elements
func (cm *ColorManager) GetRandomNoiseColor() string {
	if len(cm.noiseColors) == 0 {
		return "#cccccc" // fallback
	}

	index, err := secureRandomInt(len(cm.noiseColors))
	if err != nil {
		return cm.noiseColors[0] // fallback to first color
	}

	return cm.noiseColors[index]
}

// GetBackgroundColor returns the configured background color
func (cm *ColorManager) GetBackgroundColor() string {
	return cm.background
}

// GetRandomColorWithOpacity returns a random color with specified opacity
func (cm *ColorManager) GetRandomColorWithOpacity(opacity float64) string {
	color := cm.GetRandomNoiseColor()
	// Convert hex to rgba if needed (simplified implementation)
	if opacity < 1.0 {
		return fmt.Sprintf("%s;opacity:%.2f", color, opacity)
	}
	return color
}
