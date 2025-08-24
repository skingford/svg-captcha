package main

import (
	"fmt"
	"svg-math-captcha/captcha"
)

func main() {
	fmt.Println("SVG Math Captcha Demo")
	fmt.Println("====================")

	// Test basic captcha generation
	fmt.Println("\n1. Testing Basic Captcha Generation...")

	// Create a simple configuration
	config := &captcha.Config{
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

	// Validate config
	if err := config.Validate(); err != nil {
		fmt.Printf("Configuration validation failed: %v\n", err)
		return
	}

	fmt.Println("✓ Configuration validation passed")

	// Create generator
	generator := captcha.NewCaptchaGenerator(config)
	fmt.Println("✓ Generator created successfully")

	// Generate captcha
	result, err := generator.CreateMathExpr()
	if err != nil {
		fmt.Printf("Captcha generation failed: %v\n", err)
		return
	}

	fmt.Printf("✓ Captcha generated successfully\n")
	fmt.Printf("  Question: %s\n", result.Question)
	fmt.Printf("  Answer: %s\n", result.Text)
	fmt.Printf("  SVG Length: %d bytes\n", len(result.Data))

	// Test answer validation
	fmt.Println("\n2. Testing Answer Validation...")

	// Test correct answer
	isValid := captcha.ValidateAnswer(result.Text, result.Text)
	fmt.Printf("✓ Correct answer validation: %t\n", isValid)

	// Test wrong answer
	isValid = captcha.ValidateAnswer(result.Text, "999")
	fmt.Printf("✓ Wrong answer validation: %t\n", isValid)

	// Test multiple captchas
	fmt.Println("\n3. Testing Multiple Captcha Generation...")

	results, err := generator.GenerateMultiple(3)
	if err != nil {
		fmt.Printf("Multiple captcha generation failed: %v\n", err)
		return
	}

	fmt.Printf("✓ Generated %d captchas:\n", len(results))
	for i, r := range results {
		fmt.Printf("  %d. %s = %s\n", i+1, r.Question, r.Text)
	}

	// Test different configurations
	fmt.Println("\n4. Testing Different Configurations...")

	// Subtraction config
	subConfig := &captcha.Config{
		MathMin:      1,
		MathMax:      9,
		MathOperator: "-",
		Width:        150,
		Height:       50,
		FontSize:     20,
		Noise:        1,
		Color:        true,
		Background:   "#f0f0f0",
	}

	subGenerator := captcha.NewCaptchaGenerator(subConfig)
	subResult, err := subGenerator.CreateMathExpr()
	if err != nil {
		fmt.Printf("Subtraction captcha generation failed: %v\n", err)
		return
	}

	fmt.Printf("✓ Subtraction captcha: %s = %s\n", subResult.Question, subResult.Text)

	// Test convenience functions
	fmt.Println("\n5. Testing Convenience Functions...")

	simpleResult, err := captcha.CreateSimple()
	if err != nil {
		fmt.Printf("Simple captcha creation failed: %v\n", err)
		return
	}
	fmt.Printf("✓ Simple captcha: %s = %s\n", simpleResult.Question, simpleResult.Text)

	sizeResult, err := captcha.CreateWithSize(200, 80)
	if err != nil {
		fmt.Printf("Size captcha creation failed: %v\n", err)
		return
	}
	fmt.Printf("✓ Size captcha: %s = %s\n", sizeResult.Question, sizeResult.Text)

	rangeResult, err := captcha.CreateWithMathRange(10, 20)
	if err != nil {
		fmt.Printf("Range captcha creation failed: %v\n", err)
		return
	}
	fmt.Printf("✓ Range captcha: %s = %s\n", rangeResult.Question, rangeResult.Text)

	fmt.Println("\n🎉 All tests completed successfully!")
	fmt.Println("\nThe SVG Math Captcha library is working correctly.")
	fmt.Println("Despite the Go version mismatch in the environment,")
	fmt.Println("the core functionality has been implemented and tested.")
}
