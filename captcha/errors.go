package captcha

import "fmt"

// Error type constants
const (
	ErrInvalidConfig  = "INVALID_CONFIG"
	ErrMathGeneration = "MATH_GENERATION_FAILED"
	ErrSVGGeneration  = "SVG_GENERATION_FAILED"
	ErrFontLoadFailed = "FONT_LOAD_FAILED"
	ErrRenderFailed   = "RENDER_FAILED"
)

// CaptchaError represents an error that occurred during captcha generation
type CaptchaError struct {
	Type    string `json:"type"`
	Message string `json:"message"`
	Code    int    `json:"code"`
}

// Error implements the error interface
func (e *CaptchaError) Error() string {
	return fmt.Sprintf("[%s] %s (code: %d)", e.Type, e.Message, e.Code)
}

// NewError creates a new CaptchaError
func NewError(errorType, message string, code int) *CaptchaError {
	return &CaptchaError{
		Type:    errorType,
		Message: message,
		Code:    code,
	}
}
