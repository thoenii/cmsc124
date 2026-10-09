package parser

import (
	"cmsc124/scanner"
	"fmt"
)

// Return error if the parser has encountered any errors during parsing
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

