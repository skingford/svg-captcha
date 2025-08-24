package captcha

import (
	"log"
	"strconv"
	"sync"
)

// CaptchaResult represents the result of captcha generation
type CaptchaResult struct {
	Data     string `json:"data"`     // SVG XML content
	Text     string `json:"text"`     // Answer to the math expression
	Question string `json:"question"` // Human-readable question
}

// CaptchaGenerator is the main engine for generating captchas
type CaptchaGenerator struct {
	config      *Config
	mathGen     *MathExpressionGenerator
	svgRenderer *SVGRenderer
	noiseGen    *NoiseGenerator
	mutex       sync.RWMutex
}

// NewCaptchaGenerator creates a new captcha generator with the given configuration
func NewCaptchaGenerator(config *Config) *CaptchaGenerator {
	if config == nil {
		config = DefaultConfig()
	}

	// Validate configuration
	if err := config.Validate(); err != nil {
		log.Printf("Warning: Invalid configuration, using defaults: %v", err)
		config = DefaultConfig()
	}

	return &CaptchaGenerator{
		config:      config,
		mathGen:     NewMathExpressionGenerator(config),
		svgRenderer: NewSVGRenderer(config),
		noiseGen:    NewNoiseGenerator(),
	}
}

// CreateMathExpr generates a math expression captcha with default settings
func (cg *CaptchaGenerator) CreateMathExpr() (*CaptchaResult, error) {
	return cg.CreateMathExprWithOptions(cg.config)
}

// CreateMathExprWithOptions generates a math expression captcha with custom configuration
func (cg *CaptchaGenerator) CreateMathExprWithOptions(opts *Config) (*CaptchaResult, error) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("Captcha generation panic recovered: %v", r)
		}
	}()

	if opts == nil {
		opts = cg.config
	}

	// Validate options
	if err := opts.Validate(); err != nil {
		return nil, err
	}

	// Generate math expression
	expr, err := cg.mathGen.GenerateExpression()
	if err != nil {
		return nil, err
	}

	// Create temporary renderer with new options if different
	renderer := cg.svgRenderer
	if opts != cg.config {
		renderer = NewSVGRenderer(opts)
	}

	// Render SVG
	svgData, err := renderer.RenderMathExpression(expr, opts)
	if err != nil {
		return nil, err
	}

	return &CaptchaResult{
		Data:     svgData,
		Text:     strconv.Itoa(expr.Answer),
		Question: expr.Question,
	}, nil
}

// UpdateConfig updates the generator's configuration
func (cg *CaptchaGenerator) UpdateConfig(config *Config) error {
	if config == nil {
		return NewError(ErrInvalidConfig, "config cannot be nil", 400)
	}

	if err := config.Validate(); err != nil {
		return err
	}

	cg.mutex.Lock()
	defer cg.mutex.Unlock()

	cg.config = config
	cg.mathGen = NewMathExpressionGenerator(config)
	cg.svgRenderer = NewSVGRenderer(config)

	return nil
}

// GetConfig returns a copy of the current configuration
func (cg *CaptchaGenerator) GetConfig() *Config {
	cg.mutex.RLock()
	defer cg.mutex.RUnlock()

	// Return a copy to prevent external modification
	configCopy := *cg.config
	return &configCopy
}

// GenerateMultiple generates multiple captchas at once
func (cg *CaptchaGenerator) GenerateMultiple(count int) ([]*CaptchaResult, error) {
	if count <= 0 {
		return nil, NewError(ErrInvalidConfig, "count must be positive", 400)
	}

	if count > 100 {
		return nil, NewError(ErrInvalidConfig, "count cannot exceed 100", 400)
	}

	results := make([]*CaptchaResult, 0, count)

	for i := 0; i < count; i++ {
		result, err := cg.CreateMathExpr()
		if err != nil {
			return nil, err
		}
		results = append(results, result)
	}

	return results, nil
}

// ValidateAnswer checks if the provided answer matches the expected result
func ValidateAnswer(expected, provided string) bool {
	return expected == provided
}

// CreateSimple is a convenience function to quickly create a captcha with default settings
func CreateSimple() (*CaptchaResult, error) {
	generator := NewCaptchaGenerator(DefaultConfig())
	return generator.CreateMathExpr()
}

// CreateWithSize is a convenience function to create a captcha with specific dimensions
func CreateWithSize(width, height int) (*CaptchaResult, error) {
	config := DefaultConfig()
	config.Width = width
	config.Height = height

	generator := NewCaptchaGenerator(config)
	return generator.CreateMathExpr()
}

// CreateWithMathRange is a convenience function to create a captcha with specific number range
func CreateWithMathRange(min, max int) (*CaptchaResult, error) {
	config := DefaultConfig()
	config.MathMin = min
	config.MathMax = max

	generator := NewCaptchaGenerator(config)
	return generator.CreateMathExpr()
}
