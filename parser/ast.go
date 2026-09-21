package parser

import "cmsc124/scanner"

type Expr interface {
	exprNode()
}

type LiteralExpr struct {
	Value interface{}
}

func (LiteralExpr) exprNode() {}

type BinaryExpr struct {
	Left 	  Expr
	Operator  scanner.Token
	Right 	  Expr
}

func (BinaryExpr) exprNode() {}

type UnaryExpr struct {
	Operator  scanner.Token
	Right 	  Expr
}

func (UnaryExpr) exprNode() {}

type GroupingExpr struct {
	Expression Expr
}

func (GroupingExpr) exprNode() {}