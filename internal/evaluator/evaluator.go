package evaluator

import (
	"fmt"

	"github.com/expr-lang/expr"
)

// Evaluate checks if the condition is true given the value
func Evaluate(condition string, value any) (bool, error) {
	// 1. Compile the expression
	program, err := expr.Compile(condition, expr.Env(map[string]any{"value": value}))
	if err != nil {
		return false, fmt.Errorf("invalid condition '%s': %w", condition, err)
	}

	// 2. Run the expression
	output, err := expr.Run(program, map[string]any{"value": value})
	if err != nil {
		return false, fmt.Errorf("execution failed: %w", err)
	}

	// 3. Assert boolean result
	result, ok := output.(bool)
	if !ok {
		return false, fmt.Errorf("condition must return boolean, got %T", output)
	}

	return result, nil
}
