package parser

import (
	"cmsc124/scanner"
	"fmt"
	"strconv"
)

type Expr interface {
	exprNode()
}

// For literal values (e.g. number, string, boolean, or none)
type LiteralExpr struct {
	Value interface{}
}

// marks LiteralExpr as an AST expression
func (LiteralExpr) exprNode() {}

// For operatons with two operands (left and right)
type BinaryExpr struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

// marks BinaryExpr as an AST expression
func (BinaryExpr) exprNode() {}

type UnaryExpr struct {
	Operator scanner.Token
	Right    Expr
}

func (UnaryExpr) exprNode() {}

type GroupingExpr struct {
	Expression Expr
}

func (GroupingExpr) exprNode() {}

// PrintExpr converts an AST expression into its string representation.
func PrintExpr(expr Expr) string {
	switch e := expr.(type) {

	case LiteralExpr:
		switch v := e.Value.(type) {
		case nil:
			return "nil"
		case float64:
			return strconv.FormatFloat(v, 'f', 1, 64)
		case string:
			return v
		default:
			return fmt.Sprint(v)
		}

	case BinaryExpr:
		return "(" + string(e.Operator.Lexeme) + " " +
			PrintExpr(e.Left) + " " +
			PrintExpr(e.Right) + ")"

	case UnaryExpr:
		return "(" + string(e.Operator.Type) + " " +
			PrintExpr(e.Right) + ")"

	case GroupingExpr:
		return "(group " + PrintExpr(e.Expression) + ")"

	default:
		return ""
	}
}
