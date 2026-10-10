package parser

import (
	"cmsc124/scanner"
	"fmt"
	"strconv"
)

// ---------- Interfaces ----------

// Expr is implemented by every expression node in the AST.
type Expr interface {
	exprNode()
}

// Stmt is implemented by every statement node in the AST.
type Stmt interface {
	stmtNode()
}

// ---------- Statements ----------

// BlockStmt is a block of statements, e.g. `{ stmt1; stmt2; }`
type BlockStmt struct {
	Statements []Stmt
}

// IfStmt represents an if statement
type IfStmt struct {
	Condition Expr
	ThenBranch Stmt
	ElseBranch Stmt // nil if there is no else branch
}

// ForStmt
type ForStmt struct {
	Initializer Stmt
	Condition Expr
	Body Stmt
}

// WhileStmt
type WhileStmt struct {
	Condition Expr
	Body Stmt
}

// VarStmt declares a variable, e.g. `fang x = 5;`
type VarStmt struct {
	Name        scanner.Token
	Initializer Expr // nil when there is no "= ..."
}

// ExprStmt is an expression followed by ";"
type ExprStmt struct {
	Expression Expr
}

// HissStmt prints a value, e.g. `hiss x * 2;`
type HissStmt struct {
	Expression Expr
}

func (ForStmt) stmtNode() {}
func (WhileStmt) stmtNode() {}
func (BlockStmt) stmtNode() {}
func (IfStmt) stmtNode()  {}
func (VarStmt) stmtNode()  {}
func (ExprStmt) stmtNode() {}
func (HissStmt) stmtNode() {}

// ---------- Expressions ----------

// LiteralExpr holds a literal value (number, string, boolean, or none).
type LiteralExpr struct {
	Value interface{}
}

// VariableExpr is a variable being used, like x in `x + 1`.
type VariableExpr struct {
	Name scanner.Token
}

// GroupingExpr is an expression inside parentheses.
type GroupingExpr struct {
	Expression Expr
}

// UnaryExpr is an operation with one operand, e.g. `-x`.
type UnaryExpr struct {
	Operator scanner.Token
	Right    Expr
}

// BinaryExpr is an operation with two operands, e.g. `a + b`.
type BinaryExpr struct {
	Left     Expr
	Operator scanner.Token
	Right    Expr
}

// AssignExpr assigns a value to an existing variable, e.g. `x = 10`.
type AssignExpr struct {
	Name  scanner.Token
	Value Expr
}

func (LiteralExpr) exprNode()  {}
func (VariableExpr) exprNode() {}
func (GroupingExpr) exprNode() {}
func (UnaryExpr) exprNode()    {}
func (BinaryExpr) exprNode()   {}
func (AssignExpr) exprNode()   {}

// ---------- Printing ----------

// PrintStmt converts a statement into its string representation.
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

// PrintExpr converts an expression into its string representation.
func PrintExpr(expr Expr) string {
	switch e := expr.(type) {
	case LiteralExpr:
		return printLiteral(e.Value)

	case VariableExpr:
		return string(e.Name.Lexeme)

	case GroupingExpr:
		return "(group " + PrintExpr(e.Expression) + ")"

	case UnaryExpr:
		return "(" + string(e.Operator.Lexeme) + " " + PrintExpr(e.Right) + ")"

	case BinaryExpr:
		return "(" + string(e.Operator.Lexeme) + " " +
			PrintExpr(e.Left) + " " +
			PrintExpr(e.Right) + ")"

	case AssignExpr:
		return "(= " + string(e.Name.Lexeme) + " " + PrintExpr(e.Value) + ")"

	default:
		return ""
	}
}

// printLiteral formats a literal based on its actual value type.
func printLiteral(value interface{}) string {
	switch v := value.(type) {
	case nil:
		return "nil"
	case float64:
		return strconv.FormatFloat(v, 'f', 1, 64) // always one decimal, e.g. 5.0
	case string:
		return v
	default:
		return fmt.Sprint(v)
	}
}
