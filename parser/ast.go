package parser

import (
	"cmsc124/scanner"
	"fmt"
	"strconv"
)

// interface for all types of AST expressions
type Expr interface {
	exprNode()
}

type Stmt interface {
	stmtNode()
}

// fang x = 5;
type VarStmt struct {
	Name        scanner.Token
	Initializer Expr // nil when there is no "= ..."
}

func (VarStmt) stmtNode() {}

// an expression followed by ";"
type ExprStmt struct {
	Expression Expr
}

func (ExprStmt) stmtNode() {}

// A variable being used, like x in "x + 1"
type VariableExpr struct {
	Name scanner.Token
}

func (VariableExpr) exprNode() {}

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

// For operations with only one operand
type UnaryExpr struct {
	Operator scanner.Token
	Right    Expr
}

// marks UnaryExpr as an AST expression
func (UnaryExpr) exprNode() {}

// Expression inside the parenthesis
type GroupingExpr struct {
	Expression Expr
}

func (GroupingExpr) exprNode() {}

// hiss x * 2;
type HissStmt struct {
	Expression Expr
}

func (HissStmt) stmtNode() {}

// x = 10
type AssignExpr struct {
	Name  scanner.Token
	Value Expr
}

func (AssignExpr) exprNode() {}

func PrintStmt(stmt Stmt) string {
	switch s := stmt.(type) {
	case VarStmt:
		if s.Initializer == nil {
			return "(var " + string(s.Name.Lexeme) + ")"
		}
		return "(var " + string(s.Name.Lexeme) + " " + PrintExpr(s.Initializer) + ")"
	case ExprStmt:
		return PrintExpr(s.Expression)
	case HissStmt:
		return "(print " + PrintExpr(s.Expression) + ")"
	default:
		return ""
	}
}

// PrintExpr converts an AST expression into its string representation.
func PrintExpr(expr Expr) string {
	switch e := expr.(type) {

	case VariableExpr:
		return string(e.Name.Lexeme)
	case LiteralExpr:
		// Prints literal based on actual value type
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
		return "(" + string(e.Operator.Lexeme) + " " +
			PrintExpr(e.Right) + ")"

	case GroupingExpr:
		return "(group " + PrintExpr(e.Expression) + ")"

	case AssignExpr:
		return "(= " + string(e.Name.Lexeme) + " " + PrintExpr(e.Value) + ")"

	default:
		return ""
	}
}
