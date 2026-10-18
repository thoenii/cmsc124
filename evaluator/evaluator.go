package evaluator

import (
	"cmsc124/parser"
	"cmsc124/scanner"
	"fmt"
)

func Evaluate(expr parser.Expr) (Value, error) {
	switch e := expr.(type) {
	case parser.LiteralExpr:
		return e.Value, nil
	case parser.GroupingExpr:
		return Evaluate(e.Expression)
	default:
		return nil, fmt.Errorf("not implemented yet")
	}
}

func unaryOp(op scanner.Token, right Value) (Value, error) {
	if op.Type == scanner.MINUS {
		n, _ := right.(float64)
		return 0 - n, nil // 0 - n = 0 // no -0 value
	}
	return nil, fmt.Errorf("not implemented yet")
}

func binaryOp(op scanner.Token, left, right Value) (Value, error) {
	switch op.Type {
	case scanner.MINUS:
		a, _ := left.(float64)
		b, _ := right.(float64)
		return a - b, nil
	case scanner.STAR:
		a, _ := left.(float64)
		b, _ := right.(float64)
		return a * b, nil
	case scanner.SLASH:
		a, _ := left.(float64)
		b, _ := right.(float64)
		return a / b, nil
	case scanner.PLUS:
		if a, ok := left.(string); ok {
			b, _ := right.(string)
			return a + b, nil
		}
		a, _ := left.(float64)
		b, _ := right.(float64)
		return a + b, nil
	}
	return nil, fmt.Errorf("not implemented yet")
}
