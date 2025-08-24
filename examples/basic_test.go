package main

import (
	"fmt"
	"log"
	"os"
	"testing"

	"svg-math-captcha/captcha"
)

func TestSvgMathCaptcha(t *testing.T) {
	fmt.Println("🔢 SVG Math Captcha - Basic Usage Example")
	fmt.Println("========================================")

	// Example 1: Basic usage with default configuration
	fmt.Println("\n1. Basic Usage:")
	result1, err := captcha.CreateSimple()
	if err != nil {
		log.Fatalf("Failed to generate captcha: %v", err)
	}

	fmt.Printf("   Question: %s\n", result1.Question)
	fmt.Printf("   Answer: %s\n", result1.Text)
	fmt.Printf("   SVG Length: %d bytes\n", len(result1.Data))

	// Save to file
	err = saveToFile("basic_captcha.svg", result1.Data)
	if err != nil {
		log.Printf("Failed to save basic captcha: %v", err)
	} else {
		fmt.Println("   ✓ Saved as basic_captcha.svg")
	}

	// Example 2: Custom configuration
	fmt.Println("\n2. Custom Configuration:")
	config := &captcha.Config{
		MathMin:      5,
		MathMax:      15,
		MathOperator: "+-", // Both addition and subtraction
		Width:        250,
		Height:       80,
		FontSize:     28,
		Noise:        3,
		Color:        true,
		Background:   "#cc9966",
	}

	generator := captcha.NewCaptchaGenerator(config)
	result2, err := generator.CreateMathExpr()
	if err != nil {
		log.Fatalf("Failed to generate custom captcha: %v", err)
	}

	fmt.Printf("   Question: %s\n", result2.Question)
	fmt.Printf("   Answer: %s\n", result2.Text)
	fmt.Printf("   SVG Length: %d bytes\n", len(result2.Data))

	err = saveToFile("custom_captcha.svg", result2.Data)
	if err != nil {
		log.Printf("Failed to save custom captcha: %v", err)
	} else {
		fmt.Println("   ✓ Saved as custom_captcha.svg")
	}

	// Example 3: Generate multiple captchas
	fmt.Println("\n3. Multiple Captchas:")
	results, err := generator.GenerateMultiple(5)
	if err != nil {
		log.Fatalf("Failed to generate multiple captchas: %v", err)
	}

	for i, result := range results {
		fmt.Printf("   Captcha %d: %s = %s\n", i+1, result.Question, result.Text)

		filename := fmt.Sprintf("captcha_%d.svg", i+1)
		err = saveToFile(filename, result.Data)
		if err != nil {
			log.Printf("Failed to save %s: %v", filename, err)
		}
	}
	fmt.Printf("   ✓ Generated %d captchas\n", len(results))

	// Example 4: Different size captchas
	fmt.Println("\n4. Different Sizes:")
	sizes := []struct{ width, height int }{
		{100, 40},
		{200, 60},
		{300, 100},
	}

	for _, size := range sizes {
		result, err := captcha.CreateWithSize(size.width, size.height)
		if err != nil {
			log.Printf("Failed to generate captcha with size %dx%d: %v", size.width, size.height, err)
			continue
		}

		filename := fmt.Sprintf("captcha_size_%dx%d.svg", size.width, size.height)
		err = saveToFile(filename, result.Data)
		if err != nil {
			log.Printf("Failed to save %s: %v", filename, err)
		} else {
			fmt.Printf("   ✓ %dx%d: %s = %s\n", size.width, size.height, result.Question, result.Text)
		}
	}

	// Example 5: Different math ranges
	fmt.Println("\n5. Different Math Ranges:")
	ranges := []struct{ min, max int }{
		{1, 5},    // Easy
		{10, 50},  // Medium
		{50, 100}, // Hard
	}

	for _, r := range ranges {
		result, err := captcha.CreateWithMathRange(r.min, r.max)
		if err != nil {
			log.Printf("Failed to generate captcha with range %d-%d: %v", r.min, r.max, err)
			continue
		}

		difficulty := "Easy"
		if r.max > 20 {
			difficulty = "Medium"
		}
		if r.max > 50 {
			difficulty = "Hard"
		}

		fmt.Printf("   %s (%d-%d): %s = %s\n", difficulty, r.min, r.max, result.Question, result.Text)
	}

	// Example 6: Validation
	fmt.Println("\n6. Answer Validation:")
	testResult, _ := captcha.CreateSimple()

	// Test correct answer
	isValid := captcha.ValidateAnswer(testResult.Text, testResult.Text)
	fmt.Printf("   Correct answer validation: %t ✓\n", isValid)

	// Test incorrect answer
	isValid = captcha.ValidateAnswer(testResult.Text, "999")
	fmt.Printf("   Incorrect answer validation: %t ✓\n", isValid)

	fmt.Println("\n🎉 Example completed successfully!")
	fmt.Println("📁 Check the current directory for generated SVG files")
}

func saveToFile(filename, content string) error {
	return os.WriteFile(filename, []byte(content), 0644)
}
