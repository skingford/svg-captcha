package captcha

import (
	"strconv"
	"strings"
	"testing"
)

func TestDefaultConfig(t *testing.T) {
	config := DefaultConfig()

	if config.MathMin != 1 {
		t.Errorf("Expected MathMin to be 1, got %d", config.MathMin)
	}

	if config.MathMax != 9 {
		t.Errorf("Expected MathMax to be 9, got %d", config.MathMax)
	}

	if config.MathOperator != "+" {
		t.Errorf("Expected MathOperator to be '+', got %s", config.MathOperator)
	}

	if config.Width != 150 {
		t.Errorf("Expected Width to be 150, got %d", config.Width)
	}

	if config.Height != 50 {
		t.Errorf("Expected Height to be 50, got %d", config.Height)
	}
}

func TestConfigValidation(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
	}{
		{
			name:    "valid config",
			config:  DefaultConfig(),
			wantErr: false,
		},
		{
			name: "negative MathMin",
			config: &Config{
				MathMin:  -1,
				MathMax:  9,
				Width:    150,
				Height:   50,
				FontSize: 20,
				Noise:    1,
			},
			wantErr: true,
		},
		{
			name: "MathMax <= MathMin",
			config: &Config{
				MathMin:  5,
				MathMax:  5,
				Width:    150,
				Height:   50,
				FontSize: 20,
				Noise:    1,
			},
			wantErr: true,
		},
		{
			name: "zero width",
			config: &Config{
				MathMin:  1,
				MathMax:  9,
				Width:    0,
				Height:   50,
				FontSize: 20,
				Noise:    1,
			},
			wantErr: true,
		},
		{
			name: "negative noise",
			config: &Config{
				MathMin:  1,
				MathMax:  9,
				Width:    150,
				Height:   50,
				FontSize: 20,
				Noise:    -1,
			},
			wantErr: true,
		},
		{
			name: "noise too high",
			config: &Config{
				MathMin:  1,
				MathMax:  9,
				Width:    150,
				Height:   50,
				FontSize: 20,
				Noise:    15,
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Config.Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestMathExpressionGenerator(t *testing.T) {
	config := DefaultConfig()
	generator := NewMathExpressionGenerator(config)

	// Test addition generation
	expr, err := generator.GenerateAddition()
	if err != nil {
		t.Fatalf("Failed to generate addition: %v", err)
	}

	if expr.Operator != "+" {
		t.Errorf("Expected operator '+', got %s", expr.Operator)
	}

	if expr.Answer != expr.Operand1+expr.Operand2 {
		t.Errorf("Wrong addition result: %d + %d = %d, expected %d",
			expr.Operand1, expr.Operand2, expr.Answer, expr.Operand1+expr.Operand2)
	}

	// Check operand ranges
	if expr.Operand1 < config.MathMin || expr.Operand1 > config.MathMax {
		t.Errorf("Operand1 %d out of range [%d, %d]", expr.Operand1, config.MathMin, config.MathMax)
	}

	if expr.Operand2 < config.MathMin || expr.Operand2 > config.MathMax {
		t.Errorf("Operand2 %d out of range [%d, %d]", expr.Operand2, config.MathMin, config.MathMax)
	}
}

func TestMathExpressionGeneratorSubtraction(t *testing.T) {
	config := DefaultConfig()
	generator := NewMathExpressionGenerator(config)

	// Test subtraction generation
	expr, err := generator.GenerateSubtraction()
	if err != nil {
		t.Fatalf("Failed to generate subtraction: %v", err)
	}

	if expr.Operator != "-" {
		t.Errorf("Expected operator '-', got %s", expr.Operator)
	}

	if expr.Answer != expr.Operand1-expr.Operand2 {
		t.Errorf("Wrong subtraction result: %d - %d = %d, expected %d",
			expr.Operand1, expr.Operand2, expr.Answer, expr.Operand1-expr.Operand2)
	}

	// Subtraction should always result in non-negative answer
	if expr.Answer < 0 {
		t.Errorf("Subtraction result is negative: %d", expr.Answer)
	}
}

func TestCaptchaGenerator(t *testing.T) {
	generator := NewCaptchaGenerator(DefaultConfig())

	result, err := generator.CreateMathExpr()
	if err != nil {
		t.Fatalf("Failed to generate captcha: %v", err)
	}

	// Check result fields
	if result.Data == "" {
		t.Error("Expected SVG data, got empty string")
	}

	if result.Text == "" {
		t.Error("Expected answer text, got empty string")
	}

	if result.Question == "" {
		t.Error("Expected question text, got empty string")
	}

	// Verify SVG format
	if !strings.HasPrefix(result.Data, "<?xml") {
		t.Error("Expected SVG to start with XML declaration")
	}

	if !strings.Contains(result.Data, "<svg") {
		t.Error("Expected SVG to contain <svg tag")
	}

	if !strings.HasSuffix(strings.TrimSpace(result.Data), "</svg>") {
		t.Error("Expected SVG to end with </svg> tag")
	}

	// Verify answer is numeric
	_, err = strconv.Atoi(result.Text)
	if err != nil {
		t.Errorf("Answer is not a valid integer: %s", result.Text)
	}
}

func TestCaptchaGeneratorCustomConfig(t *testing.T) {
	config := &Config{
		MathMin:      5,
		MathMax:      15,
		MathOperator: "+-",
		Width:        200,
		Height:       80,
		FontSize:     24,
		Noise:        3,
		Color:        false,
		Background:   "#ffffff",
	}

	generator := NewCaptchaGenerator(config)
	result, err := generator.CreateMathExprWithOptions(config)
	if err != nil {
		t.Fatalf("Failed to generate captcha with custom config: %v", err)
	}

	// Verify SVG contains expected dimensions
	if !strings.Contains(result.Data, `width="200"`) {
		t.Error("SVG does not contain expected width")
	}

	if !strings.Contains(result.Data, `height="80"`) {
		t.Error("SVG does not contain expected height")
	}

	// Verify background color
	if !strings.Contains(result.Data, config.Background) {
		t.Error("SVG does not contain expected background color")
	}
}

func TestGenerateMultiple(t *testing.T) {
	generator := NewCaptchaGenerator(DefaultConfig())

	// Test generating multiple captchas
	results, err := generator.GenerateMultiple(5)
	if err != nil {
		t.Fatalf("Failed to generate multiple captchas: %v", err)
	}

	if len(results) != 5 {
		t.Errorf("Expected 5 captchas, got %d", len(results))
	}

	// Verify each result is valid
	for i, result := range results {
		if result.Data == "" {
			t.Errorf("Captcha %d has empty data", i)
		}
		if result.Text == "" {
			t.Errorf("Captcha %d has empty text", i)
		}
		if result.Question == "" {
			t.Errorf("Captcha %d has empty question", i)
		}
	}

	// Test edge cases
	_, err = generator.GenerateMultiple(0)
	if err == nil {
		t.Error("Expected error for zero count")
	}

	_, err = generator.GenerateMultiple(-1)
	if err == nil {
		t.Error("Expected error for negative count")
	}

	_, err = generator.GenerateMultiple(101)
	if err == nil {
		t.Error("Expected error for count > 100")
	}
}

func TestValidateAnswer(t *testing.T) {
	tests := []struct {
		expected string
		provided string
		want     bool
	}{
		{"5", "5", true},
		{"10", "10", true},
		{"5", "6", false},
		{"", "", true},
		{"5", "", false},
		{"", "5", false},
		{"5", " 5 ", false}, // No trimming
	}

	for _, tt := range tests {
		got := ValidateAnswer(tt.expected, tt.provided)
		if got != tt.want {
			t.Errorf("ValidateAnswer(%q, %q) = %v, want %v", tt.expected, tt.provided, got, tt.want)
		}
	}
}

func TestConvenienceFunctions(t *testing.T) {
	// Test CreateSimple
	result, err := CreateSimple()
	if err != nil {
		t.Fatalf("CreateSimple failed: %v", err)
	}
	if result.Data == "" {
		t.Error("CreateSimple returned empty data")
	}

	// Test CreateWithSize
	result, err = CreateWithSize(300, 100)
	if err != nil {
		t.Fatalf("CreateWithSize failed: %v", err)
	}
	if !strings.Contains(result.Data, `width="300"`) {
		t.Error("CreateWithSize did not set correct width")
	}

	// Test CreateWithMathRange
	result, err = CreateWithMathRange(10, 20)
	if err != nil {
		t.Fatalf("CreateWithMathRange failed: %v", err)
	}
	// Parse answer to verify it's in range
	answer, err := strconv.Atoi(result.Text)
	if err != nil {
		t.Fatalf("CreateWithMathRange returned non-numeric answer: %s", result.Text)
	}
	if answer < 10 || answer > 40 { // Max possible is 20+20=40
		t.Errorf("CreateWithMathRange answer %d seems out of expected range", answer)
	}
}

func TestColorManager(t *testing.T) {
	config := DefaultConfig()
	colorMgr := NewColorManager(config)

	// Test color generation
	textColor := colorMgr.GetRandomTextColor()
	if textColor == "" {
		t.Error("GetRandomTextColor returned empty string")
	}
	if !strings.HasPrefix(textColor, "#") {
		t.Errorf("Expected color to start with #, got %s", textColor)
	}

	noiseColor := colorMgr.GetRandomNoiseColor()
	if noiseColor == "" {
		t.Error("GetRandomNoiseColor returned empty string")
	}

	backgroundColor := colorMgr.GetBackgroundColor()
	if backgroundColor != config.Background {
		t.Errorf("Expected background color %s, got %s", config.Background, backgroundColor)
	}
}

func TestNoiseGenerator(t *testing.T) {
	noiseGen := NewNoiseGenerator()
	colorMgr := NewColorManager(DefaultConfig())

	// Test line generation
	lines := noiseGen.GenerateLines(5, 200, 100, colorMgr)
	if len(lines) > 5 { // Might be less due to random generation failures
		t.Errorf("Expected at most 5 lines, got %d", len(lines))
	}

	// Test circle generation
	circles := noiseGen.GenerateDots(3, 200, 100, colorMgr)
	if len(circles) > 3 {
		t.Errorf("Expected at most 3 circles, got %d", len(circles))
	}
}

func TestSVGRenderer(t *testing.T) {
	config := DefaultConfig()
	renderer := NewSVGRenderer(config)

	expr := &MathExpression{
		Operand1: 3,
		Operand2: 5,
		Operator: "+",
		Answer:   8,
		Question: "3 + 5 = ?",
	}

	svgData, err := renderer.RenderMathExpression(expr, config)
	if err != nil {
		t.Fatalf("Failed to render SVG: %v", err)
	}

	if svgData == "" {
		t.Error("SVG data is empty")
	}

	// Verify SVG structure
	if !strings.Contains(svgData, "<svg") {
		t.Error("SVG does not contain svg tag")
	}

	if !strings.Contains(svgData, "xmlns=\"http://www.w3.org/2000/svg\"") {
		t.Error("SVG does not contain proper namespace")
	}

	if !strings.Contains(svgData, "<rect") {
		t.Error("SVG does not contain background rectangle")
	}
}

func TestCaptchaError(t *testing.T) {
	err := NewError(ErrInvalidConfig, "test message", 400)

	if err.Type != ErrInvalidConfig {
		t.Errorf("Expected error type %s, got %s", ErrInvalidConfig, err.Type)
	}

	if err.Message != "test message" {
		t.Errorf("Expected error message 'test message', got %s", err.Message)
	}

	if err.Code != 400 {
		t.Errorf("Expected error code 400, got %d", err.Code)
	}

	errorString := err.Error()
	if !strings.Contains(errorString, "INVALID_CONFIG") {
		t.Error("Error string does not contain error type")
	}

	if !strings.Contains(errorString, "test message") {
		t.Error("Error string does not contain error message")
	}
}

// Benchmark tests
func BenchmarkCaptchaGeneration(b *testing.B) {
	generator := NewCaptchaGenerator(DefaultConfig())

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := generator.CreateMathExpr()
		if err != nil {
			b.Fatalf("Failed to generate captcha: %v", err)
		}
	}
}

func BenchmarkMathExpressionGeneration(b *testing.B) {
	config := DefaultConfig()
	generator := NewMathExpressionGenerator(config)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := generator.GenerateExpression()
		if err != nil {
			b.Fatalf("Failed to generate math expression: %v", err)
		}
	}
}

func BenchmarkSVGRendering(b *testing.B) {
	config := DefaultConfig()
	renderer := NewSVGRenderer(config)
	expr := &MathExpression{
		Operand1: 3,
		Operand2: 5,
		Operator: "+",
		Answer:   8,
		Question: "3 + 5 = ?",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := renderer.RenderMathExpression(expr, config)
		if err != nil {
			b.Fatalf("Failed to render SVG: %v", err)
		}
	}
}
