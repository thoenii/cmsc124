package parser

import (
	"cmsc124/scanner"
	"fmt"
)

type Parser struct {
	tokens  []scanner.Token
	current int
	errors  []error // syntax errors encountered during parsing
}

// ParseAll parses every expression in the token stream, one after another,
// until EOF is reached. It returns the list of successfully parsed
// expressions and the first error encountered, if any.
func (p *Parser) ParseAll() ([]Expr, error) {
	var exprs []Expr

	for !p.isAtEnd() {
		expr, err := p.expression()
		if err != nil {
			return nil, err
		}
		exprs = append(exprs, expr)
	}

	return exprs, nil
}

// creates a parser for tokens
func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{tokens: tokens}
}

// parse an expression
func (p *Parser) Parse() (Expr, error) {
	return p.expression()
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return nil, err
	}

	// Parse repeated + and - operations
	for p.match(scanner.PLUS, scanner.MINUS) {
		operator := p.previous()

		right, err := p.factor()
		if err != nil {
			return nil, err
		}

		expr = BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return nil, err
	}

	// Parse repeated * and / operations
	for p.match(scanner.STAR, scanner.SLASH) {
		operator := p.previous()

		right, err := p.unary() // fixed: was p.primary()
		if err != nil {
			return nil, err
		}

		expr = BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(scanner.MINUS) {
		operator := p.previous()

		right, err := p.unary()
		if err != nil {
			return nil, err
		}

		return UnaryExpr{
			Operator: operator,
			Right:    right,
		}, nil
	}

	return p.primary()
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return nil, err
	}

	// Parse repeated <, <=, >, and >= operations
	for p.match(scanner.GREATER, scanner.GREATER_EQUAL, scanner.LESS, scanner.LESS_EQUAL) {
		operator := p.previous()

		right, err := p.term()
		if err != nil {
			return nil, err
		}

		expr = BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return nil, err
	}

	// Parse repeated == and != operations
	for p.match(scanner.EQUAL_EQUAL, scanner.BANG_EQUAL) {
		operator := p.previous()

		right, err := p.comparison()
		if err != nil {
			return nil, err
		}

		expr = BinaryExpr{
			Left:     expr,
			Operator: operator,
			Right:    right,
		}
	}
	return expr, nil
}

// parser for literals and grouped expressions
func (p *Parser) primary() (Expr, error) {
	if p.match(scanner.NUMBER, scanner.STRING, scanner.TRUE, scanner.FALSE, scanner.NONE) {
		return LiteralExpr{Value: p.previous().Literal}, nil
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

// checks current token
func (p *Parser) check(tokenType scanner.TokenType) bool {
	if p.isAtEnd() {
		return tokenType == scanner.EOF
	}

	return p.peek().Type == tokenType
}

// checks the current token to match the token types
func (p *Parser) match(types ...scanner.TokenType) bool {
	for _, tokenType := range types {
		if p.check(tokenType) {
			p.advance()
			return true
		}
	}
	return false
}

// consume and return token
func (p *Parser) advance() scanner.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// peeking
func (p *Parser) peek() scanner.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() scanner.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == scanner.EOF
}


// For Error Handling
func (p *Parser) Errors() []error {
	return p.errors
}

// hadError for syntax errors
func (p *Parser) hadError() bool {
	return len(p.errors) > 0
}

// errorAt points to the current token and returns an error with the given message
func (p *Parser) errorAt(token scanner.Token, message string) error {
	if token.Type == scanner.EOF {
		return fmt.Errorf("[line %d] Error at end: %s", token.Line, message)
	} 
	return fmt.Errorf("[line %d] Error at '%s': %s", token.Line, token.Lexeme, message)
}

// consume and return token, or return an error if the current token is not of the expected type
func (p *Parser) consume(tokenType scanner.TokenType, message string) (scanner.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	return scanner.Token{}, p.errorAt(p.peek(), message)
}

// after an error, synchronize the parser to the next statement to continue parsing
func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == scanner.SEMICOLON {
			return
		}

		switch p.peek().Type {
			case scanner.LET,
			scanner.IF,
			scanner.PRINT,
			scanner.WHILE,
			scanner.FOR,
			scanner.LEFT_BRACE,
			scanner.RIGHT_BRACE:
				return
		}
		p.advance()
	}	
}

// declaration parses a declaration statement and returns it, or returns nil if there was an error
func (p *Parser) safeDeclaration() Stmt {
	stmt, err := p.declaration()
	if err != nil {
		p.errors = append(p.errors, err)
		p.synchronize()
		return nil
	}
	return stmt
}

// blockStatement
func (p *Parser) blockStatement() (Stmt, error) {
	stmts, err := p.block()
	if err != nil {
		return nil, err
	}
	return BlockStmt{Statements: stmts}, nil
}

// block parses statements until matching '{'
func (p *Parser) block() ([]Stmt, error) {
	var stmts []Stmt

	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		if stmt := p.safeDeclaration(); stmt != nil {
			stmts = append(stmts, stmt)
		}
	}

	if _, err := p.consume(scanner.RIGHT_BRACE, "Expect '}' after block."); err != nil {
		return nil, err
	}
	return stmts, nil
}







