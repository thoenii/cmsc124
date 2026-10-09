package parser

import (
	"cmsc124/scanner"
	"fmt"
)

// Parser turns a list of tokens into AST nodes (recursive descent).
type Parser struct {
	tokens  []scanner.Token
	current int // index of the next token to read
}

// NewParser creates a parser for the given tokens.
func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{tokens: tokens}
}

// ---------- Entry points ----------

// ParseProgram parses statements until EOF.
func (p *Parser) ParseProgram() ([]Stmt, error) {
	var stmts []Stmt
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			return nil, err
		}
		stmts = append(stmts, stmt)
	}
	return stmts, nil
}

// Parse parses a single expression.
func (p *Parser) Parse() (Expr, error) {
	return p.expression()
}

// ---------- Statements ----------

// declaration: a fang (variable) declaration, or any other statement.
func (p *Parser) declaration() (Stmt, error) {
	if p.match(scanner.LET) {
		return p.varDeclaration()
	}
	return p.statement()
}

// varDeclaration parses `fang name;` or `fang name = expr;`
func (p *Parser) varDeclaration() (Stmt, error) {
	if !p.check(scanner.IDENTIFIER) {
		return nil, fmt.Errorf("[line %d] Expect variable name", p.peek().Line)
	}
	name := p.advance()

	// The initializer is optional
	var initializer Expr
	if p.match(scanner.EQUAL) {
		var err error
		initializer, err = p.expression()
		if err != nil {
			return nil, err
		}
	}

	if !p.match(scanner.SEMICOLON) {
		return nil, fmt.Errorf("[line %d] Expect ';' after variable declaration", p.peek().Line)
	}
	return VarStmt{Name: name, Initializer: initializer}, nil
}

// statement: a hiss (print) statement, or a plain expression statement.
func (p *Parser) statement() (Stmt, error) {
	if p.match(scanner.PRINT) {
		return p.printStatement()
	}
	return p.expressionStatement()
}

// printStatement parses `hiss(expr);`
func (p *Parser) printStatement() (Stmt, error) {
	if !p.match(scanner.LEFT_PAREN) {
		return nil, fmt.Errorf("[line %d] Expect '(' after hiss", p.peek().Line)
	}
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if !p.match(scanner.RIGHT_PAREN) {
		return nil, fmt.Errorf("[line %d] Expect ')' after value", p.peek().Line)
	}
	if !p.match(scanner.SEMICOLON) {
		return nil, fmt.Errorf("[line %d] Expect ';' after value", p.peek().Line)
	}
	return HissStmt{Expression: expr}, nil
}

// expressionStatement parses `expr;`
func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return nil, err
	}
	if !p.match(scanner.SEMICOLON) {
		return nil, fmt.Errorf("[line %d] Expect ';' after expression", p.peek().Line)
	}
	return ExprStmt{Expression: expr}, nil
}

// ---------- Expressions (lowest to highest precedence) ----------

func (p *Parser) expression() (Expr, error) {
	return p.assignment()
}

// assignment parses `name = value`. Recurses on the right side so x = y = 3 works.
func (p *Parser) assignment() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return nil, err
	}

	if p.match(scanner.EQUAL) {
		equals := p.previous()

		value, err := p.assignment()
		if err != nil {
			return nil, err
		}

		// Only a variable can be assigned to
		if variable, ok := expr.(VariableExpr); ok {
			return AssignExpr{Name: variable.Name, Value: value}, nil
		}
		return nil, fmt.Errorf("[line %d] Invalid assignment target", equals.Line)
	}

	return expr, nil
}

// equality handles == and !=
func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.EQUAL_EQUAL, scanner.BANG_EQUAL) {
		operator := p.previous()

		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		expr = BinaryExpr{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

// comparison handles <, <=, >, and >=
func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.previous()

		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = BinaryExpr{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

// term handles + and -
func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.PLUS, scanner.MINUS) {
		operator := p.previous()

		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = BinaryExpr{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

// factor handles * and /
func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	for p.match(scanner.STAR, scanner.SLASH) {
		operator := p.previous()

		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		expr = BinaryExpr{Left: expr, Operator: operator, Right: right}
	}
	return expr, nil
}

// unary handles a leading minus, e.g. -x or --x
func (p *Parser) unary() (Expr, error) {
	if p.match(scanner.MINUS) {
		operator := p.previous()

		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return UnaryExpr{Operator: operator, Right: right}, nil
	}

	return p.primary()
}

// primary handles literals, variables, and grouped expressions
func (p *Parser) primary() (Expr, error) {
	if p.match(scanner.NUMBER, scanner.STRING, scanner.TRUE, scanner.FALSE, scanner.NONE) {
		return LiteralExpr{Value: p.previous().Literal}, nil
	}

	if p.match(scanner.IDENTIFIER) {
		return VariableExpr{Name: p.previous()}, nil
	}

	if p.match(scanner.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}

		if !p.match(scanner.RIGHT_PAREN) {
			return nil, fmt.Errorf("[line %d] Expect ')' after expression", p.peek().Line)
		}
		return GroupingExpr{Expression: expr}, nil
	}

	return nil, fmt.Errorf("[line %d] Expect expression", p.peek().Line)
}

// ---------- Token helpers ----------

// match consumes the current token if it is one of the given types.
func (p *Parser) match(types ...scanner.TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

// check reports whether the current token has the given type (without consuming it).
func (p *Parser) check(tokenType scanner.TokenType) bool {
	if p.isAtEnd() {
		return tokenType == scanner.EOF
	}
	return p.peek().Type == tokenType
}

// advance consumes the current token and returns it.
func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// peek returns the current token without consuming it.
func (p *Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

// previous returns the most recently consumed token.
func (p *Parser) previous() scanner.Token {
	return p.tokens[p.current-1]
}

// isAtEnd reports whether the current token is EOF.
func (p *Parser) isAtEnd() bool {
	return p.peek().Type == scanner.EOF
}
