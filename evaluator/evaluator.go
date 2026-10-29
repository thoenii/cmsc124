package evaluator

import (
	"cmsc124/parser"
	"cmsc124/scanner"
	"fmt"
)

// Error Handling

// RuntimeError carries the operator token so the report can include a line number.
type RuntimeError struct {
	Token   scanner.Token
	Message string
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("[line %d] Runtime error: %s", e.Token.Line, e.Message)
}

func checkNumberOperand(op scanner.Token, v Value) (float64, error) {
	if n, ok := v.(float64); ok {
		return n, nil
	}
	return 0, RuntimeError{op, "Operand must be a number."}
}

func checkNumberOperands(op scanner.Token, l, r Value) (float64, float64, error) {
	a, okA := l.(float64)
	b, okB := r.(float64)
	if okA && okB {
		return a, b, nil
	}
	return 0, 0, RuntimeError{op, "Operands must be numbers."}
}

// Expressions

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
    	return nil, fmt.Errorf( "unsupported expression type: %T", expr)
	}
}

func unaryOp(op scanner.Token, right Value) (Value, error) {
	switch op.Type {
	case scanner.MINUS:
		n, err := checkNumberOperand(op, right)
		if err != nil {
			return nil, err
		}
		return 0 - n, nil // 0 - n so -0 never prints as "-0"
	case scanner.BANG:
		return !isTruthy(right), nil
	}
	return nil, RuntimeError{op, "Unknown unary operator."}
}

func binaryOp(op scanner.Token, left, right Value) (Value, error) {
	switch op.Type {
	case scanner.EQUAL_EQUAL:
		return isEqual(left, right), nil

	case scanner.BANG_EQUAL:
		return !isEqual(left, right), nil

	case scanner.PLUS:
		if a, ok := left.(string); ok {
			if b, ok := right.(string); ok {
				return a + b, nil
			}
		}
		if a, ok := left.(float64); ok {
			if b, ok := right.(float64); ok {
				return a + b, nil
			}
		}
		return nil, RuntimeError{op, "Operands must be two numbers or two strings."}

	case scanner.MINUS, scanner.STAR, scanner.SLASH,
		scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL:

		// one check for every numbers-only operator
		a, b, err := checkNumberOperands(op, left, right)
		if err != nil {
			return nil, err
		}
		switch op.Type {
		case scanner.MINUS:
			return a - b, nil
		case scanner.STAR:
			return a * b, nil
		case scanner.SLASH:
			if b == 0 {
				return nil, RuntimeError{op, "Division by zero."}
			}
			return a / b, nil
		case scanner.GREATER:
			return a > b, nil
		case scanner.GREATER_EQUAL:
			return a >= b, nil
		case scanner.LESS:
			return a < b, nil
		case scanner.LESS_EQUAL:
			return a <= b, nil
		}
	}
	return nil, RuntimeError{op, "Unknown binary operator."}
}

// Statement executions

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
	return fmt.Errorf("unsupported statement type: %T", stmt)
}
