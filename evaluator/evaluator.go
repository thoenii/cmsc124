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
	case parser.UnaryExpr:
		right, err := Evaluate(e.Right)
		if err != nil {
			return nil, err
		}
		return unaryOp(e.Operator, right)
	case parser.BinaryExpr:
		left, err := Evaluate(e.Left)
		if err != nil {
			return nil, err
		}
		right, err := Evaluate(e.Right)
		if err != nil {
			return nil, err
		}
		return binaryOp(e.Operator, left, right)
	default:
		return nil, fmt.Errorf("not implemented yet")
	}
}

func unaryOp(op scanner.Token, right Value) (Value, error) {
	switch op.Type {
	case scanner.MINUS:
		n, _ := right.(float64)
		return 0 - n, nil // 0 - n so -0 never prints as "-0"
	case scanner.BANG:
		return !isTruthy(right), nil
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

	case scanner.EQUAL_EQUAL:
		return isEqual(left, right), nil

	case scanner.BANG_EQUAL:
		return !isEqual(left, right), nil
	}
	return nil, fmt.Errorf("not implemented yet")
}

// statement executions
func Exec(stmt parser.Stmt) error {
	switch s := stmt.(type) {
	case parser.ExprStmt:
		_, err := Evaluate(s.Expression)
		return err

	case parser.HissStmt:
		val, err := Evaluate(s.Expression)
		if err != nil {
			return err
		}
		fmt.Println(Stringify(val))
		return nil

	case parser.BlockStmt:
		for _, st := range s.Statements {
			if err := Exec(st); err != nil {
				return err
			}
		}
		return nil

	case parser.IfStmt:
		cond, err := Evaluate(s.Condition)
		if err != nil {
			return err
		}
		if isTruthy(cond) {
			return Exec(s.ThenBranch)
		}
		if s.ElseBranch != nil {
			return Exec(s.ElseBranch)
		}
		return nil

	case parser.WhileStmt:
		for {
			cond, err := Evaluate(s.Condition)
			if err != nil {
				return err
			}
			if !isTruthy(cond) {
				return nil
			}
			if err := Exec(s.Body); err != nil {
				return err
			}
		}
	}
	return fmt.Errorf("not implemented yet")
}

// Helper functions

// Any types allowed. Different types are never equal, nil == nil is true,
// strings compare by value. Go's == on two `any` values already does this.
func isEqual(a, b Value) bool {
	return a == b
}

// Only false and nil are falsey. 0 and "" are truthy.
func isTruthy(v Value) bool {
	if v == nil {
		return false
	}
	if b, ok := v.(bool); ok {
		return b
	}
	return true
}
