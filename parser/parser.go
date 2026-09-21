package parser

import (
	"fmt"
	"cmsc124/scanner"
)

type Parser struct {
	tokens []scanner.Token
	current int
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
	return p.term()
}

func (p *Parser) term() (Expr, error){
	expr, err := p.primary()
	if err != nil {
		return nil, err
	}

	// Parse repeated + and - operations
	for p.match(scanner.PLUS, scanner.MINUS) {
		operator := p.previous()

		right, err := p.primary()
		if err != nil {return nil, err}

		expr = BinaryExpr{
			Left: 	  expr, 
			Operator: operator, 
			Right: 	  right,
		}
	}
	return expr, nil
}

// parser for literals and grouped expressions
func (p *Parser) primary() (Expr, error) {
	if p.match(scanner.NUMBER, scanner.STRING, scanner.TRUE, scanner.FALSE, scanner.NONE){
		return LiteralExpr{Value: p.previous().Literal}, nil
	}

	if p.match(scanner.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {return nil, err}
	

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
		if p.check(tokenType) {p.advance() 
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

