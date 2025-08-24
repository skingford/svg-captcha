package captcha

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"strings"
)

// MathExpression represents a mathematical expression for the captcha
type MathExpression struct {
	Operand1 int    `json:"operand1"`
	Operand2 int    `json:"operand2"`
	Operator string `json:"operator"`
	Answer   int    `json:"answer"`
	Question string `json:"question"` // Human-readable question like "3 + 5 = ?"
}

// MathExpressionGenerator generates mathematical expressions for captchas
type MathExpressionGenerator struct {
	minValue  int
	maxValue  int
	operators []string
}

// NewMathExpressionGenerator creates a new math expression generator
func NewMathExpressionGenerator(config *Config) *MathExpressionGenerator {
	operators := parseOperators(config.MathOperator)
	return &MathExpressionGenerator{
		minValue:  config.MathMin,
		maxValue:  config.MathMax,
		operators: operators,
	}
}

// parseOperators converts operator string to slice of operators
func parseOperators(operatorStr string) []string {
	var operators []string
	if strings.Contains(operatorStr, "+") {
		operators = append(operators, "+")
	}
	if strings.Contains(operatorStr, "-") {
		operators = append(operators, "-")
	}
	// Default to addition if no valid operators
	if len(operators) == 0 {
		operators = []string{"+"}
	}
	return operators
}

// GenerateExpression creates a random mathematical expression
func (meg *MathExpressionGenerator) GenerateExpression() (*MathExpression, error) {
	// Choose random operator
	operatorIndex, err := secureRandomInt(len(meg.operators))
	if err != nil {
		return nil, NewError(ErrMathGeneration, "failed to generate random operator", 500)
	}
	operator := meg.operators[operatorIndex]

	switch operator {
	case "+":
		return meg.GenerateAddition()
	case "-":
		return meg.GenerateSubtraction()
	default:
		return meg.GenerateAddition() // fallback
	}
}

// GenerateAddition creates an addition expression
func (meg *MathExpressionGenerator) GenerateAddition() (*MathExpression, error) {
	operand1, err := meg.generateOperand()
	if err != nil {
		return nil, err
	}

	operand2, err := meg.generateOperand()
	if err != nil {
		return nil, err
	}

	answer := operand1 + operand2
	question := fmt.Sprintf("%d + %d = ?", operand1, operand2)

	return &MathExpression{
		Operand1: operand1,
		Operand2: operand2,
		Operator: "+",
		Answer:   answer,
		Question: question,
	}, nil
}

// GenerateSubtraction creates a subtraction expression ensuring positive result
func (meg *MathExpressionGenerator) GenerateSubtraction() (*MathExpression, error) {
	operand1, err := meg.generateOperand()
	if err != nil {
		return nil, err
	}

	operand2, err := meg.generateOperand()
	if err != nil {
		return nil, err
	}

	// Ensure result is positive by making operand1 >= operand2
	if operand1 < operand2 {
		operand1, operand2 = operand2, operand1
	}

	answer := operand1 - operand2
	question := fmt.Sprintf("%d - %d = ?", operand1, operand2)

	return &MathExpression{
		Operand1: operand1,
		Operand2: operand2,
		Operator: "-",
		Answer:   answer,
		Question: question,
	}, nil
}

// generateOperand creates a random operand within the configured range
func (meg *MathExpressionGenerator) generateOperand() (int, error) {
	rangeSize := meg.maxValue - meg.minValue + 1
	randomValue, err := secureRandomInt(rangeSize)
	if err != nil {
		return 0, NewError(ErrMathGeneration, "failed to generate random operand", 500)
	}
	return meg.minValue + randomValue, nil
}

// secureRandomInt generates a cryptographically secure random integer in range [0, max)
func secureRandomInt(max int) (int, error) {
	if max <= 0 {
		return 0, fmt.Errorf("max must be positive")
	}

	maxBig := big.NewInt(int64(max))
	n, err := rand.Int(rand.Reader, maxBig)
	if err != nil {
		return 0, err
	}

	return int(n.Int64()), nil
}
