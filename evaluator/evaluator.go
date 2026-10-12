package evaluator

import (
	"cmsc124/parser"
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
