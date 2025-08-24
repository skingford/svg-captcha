package captcha

import (
	"os"
	"strconv"
)

// Config defines the configuration options for SVG math captcha generation
type Config struct {
	// Math expression settings
	MathMin      int    `json:"mathMin"`      // Minimum operand value (default: 1)
	MathMax      int    `json:"mathMax"`      // Maximum operand value (default: 9)
	MathOperator string `json:"mathOperator"` // Operators to use: "+", "-", "+-" (default: "+")

	// Visual settings
	Width      int    `json:"width"`      // SVG width in pixels (default: 150)
	Height     int    `json:"height"`     // SVG height in pixels (default: 50)
	FontSize   int    `json:"fontSize"`   // Font size (default: 20)
	Noise      int    `json:"noise"`      // Noise level 0-10 (default: 1)
	Color      bool   `json:"color"`      // Use random colors (default: true)
	Background string `json:"background"` // Background color (default: "#f0f0f0")

	// Text settings
	IgnoreChars string `json:"ignoreChars"` // Characters to avoid (default: "0o1i")
}

// DefaultConfig returns a configuration with sensible default values
func DefaultConfig() *Config {
	return &Config{
		MathMin:      1,
		MathMax:      9,
		MathOperator: "+",
		Width:        150,
		Height:       50,
		FontSize:     20,
		Noise:        1,
		Color:        true,
		Background:   "#f0f0f0",
		IgnoreChars:  "0o1i",
	}
}

// LoadConfigFromEnv loads configuration from environment variables
// Falls back to default values if environment variables are not set
func LoadConfigFromEnv() *Config {
	config := DefaultConfig()

	if val := os.Getenv("CAPTCHA_MATH_MIN"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			config.MathMin = parsed
		}
	}

	if val := os.Getenv("CAPTCHA_MATH_MAX"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			config.MathMax = parsed
		}
	}

	if val := os.Getenv("CAPTCHA_OPERATOR"); val != "" {
		config.MathOperator = val
	}

	if val := os.Getenv("CAPTCHA_WIDTH"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			config.Width = parsed
		}
	}

	if val := os.Getenv("CAPTCHA_HEIGHT"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			config.Height = parsed
		}
	}

	if val := os.Getenv("CAPTCHA_FONT_SIZE"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			config.FontSize = parsed
		}
	}

	if val := os.Getenv("CAPTCHA_NOISE"); val != "" {
		if parsed, err := strconv.Atoi(val); err == nil {
			config.Noise = parsed
		}
	}

	if val := os.Getenv("CAPTCHA_COLOR"); val != "" {
		if parsed, err := strconv.ParseBool(val); err == nil {
			config.Color = parsed
		}
	}

	if val := os.Getenv("CAPTCHA_BACKGROUND"); val != "" {
		config.Background = val
	}

	if val := os.Getenv("CAPTCHA_IGNORE_CHARS"); val != "" {
		config.IgnoreChars = val
	}

	return config
}

// Validate checks if the configuration values are valid
func (c *Config) Validate() error {
	if c.MathMin < 0 {
		return &CaptchaError{Type: ErrInvalidConfig, Message: "MathMin must be >= 0", Code: 400}
	}
	if c.MathMax <= c.MathMin {
		return &CaptchaError{Type: ErrInvalidConfig, Message: "MathMax must be > MathMin", Code: 400}
	}
	if c.Width <= 0 || c.Height <= 0 {
		return &CaptchaError{Type: ErrInvalidConfig, Message: "Width and Height must be > 0", Code: 400}
	}
	if c.FontSize <= 0 {
		return &CaptchaError{Type: ErrInvalidConfig, Message: "FontSize must be > 0", Code: 400}
	}
	if c.Noise < 0 || c.Noise > 10 {
		return &CaptchaError{Type: ErrInvalidConfig, Message: "Noise must be between 0 and 10", Code: 400}
	}
	return nil
}
