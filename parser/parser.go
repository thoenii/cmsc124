package parser

import (
	"cmsc124/scanner"
	"errors"
	"fmt"
)

// Parser turns a list of tokens into AST nodes (recursive descent).
type Parser struct {
	tokens  []scanner.Token
	current int // index of the next token to read
	errors []error // errors encountered during parsing
}

// NewParser creates a parser for the given tokens.
func NewParser(tokens []scanner.Token) *Parser {
	return &Parser{tokens: tokens}
}

// ---------- Entry points ----------

// ParseProgram parses statements until EOF.
func (p *Parser) ParseProgram() ([]Stmt, error) {
	p.errors = nil // reset errors for each parse
	var stmts []Stmt
	for !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			p.errors = append(p.errors, err)
			p.synchronize() // skip to the next statement
			continue
		}
		stmts = append(stmts, stmt)
	}
	if len(p.errors) > 0 {
		return nil, errors.Join(p.errors...)
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
	if p.match(scanner.FOR) {
		return p.forStatement()
	}
	if p.match(scanner.WHILE) {
		return p.whileStatement()
	}
	if p.match(scanner.IF) {
		return p.ifStatement()
	}
	if p.match(scanner.LEFT_BRACE) {
		stmts, err := p.block()
		if err != nil {
			return nil, err
		}
		return BlockStmt{Statements: stmts}, nil
	}
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

// ---------- Control flow ----------

// block parses declarations up until the closing brace
func (p *Parser) block() ([]Stmt, error) {
	open := p.previous() 
	var stmts []Stmt

	for !p.check(scanner.RIGHT_BRACE) && !p.isAtEnd() {
		stmt, err := p.declaration()
		if err != nil {
			p.errors = append(p.errors, err)
			p.synchronize() // skip to the next statement
			continue
		}
		stmts = append(stmts, stmt)
	}

	if _, err := p.consume(scanner.RIGHT_BRACE, fmt.Sprintf("Expect '}' to close block opened at line %d", open.Line)); err != nil {
		return nil, err
	}
	return stmts, nil
}

// ifStatement parses `if (condition) thenBranch else elseBranch`
func (p *Parser) ifStatement() (Stmt, error) {
	if _, err := p.consume(scanner.LEFT_PAREN, "Expect '(' after 'if'"); err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(scanner.RIGHT_PAREN, "Expect ')' after if condition"); err != nil {
		return nil, err
	}

	thenBranch, err := p.statement()
	if err != nil {
		return nil, err
	}

	//else
	var elseBranch Stmt
	if p.match(scanner.ELSE) {
		elseBranch, err = p.statement()
		if err != nil {
			return nil, err
		}
	}
	return IfStmt{Condition: condition, ThenBranch: thenBranch, ElseBranch: elseBranch}, nil
}

//forStatement parses `for (initializer; condition; increment) body`
func (p *Parser) forStatement() (Stmt, error) {
	if _, err := p.consume(scanner.LEFT_PAREN, "Expect '(' after 'for'"); err != nil {
		return nil, err
	}

	//initialiazer
	var initializer Stmt
	var err error
	if p.match(scanner.SEMICOLON) {
		initializer, err = p.varDeclaration()
		if err != nil {
			return nil, err
		}
	}

	var condition Expr
	if !p.check(scanner.SEMICOLON) {
		condition, err = p.expression()
		if err != nil {
			return nil, err
		}
	}
	if _, err := p.consume(scanner.SEMICOLON, "Expect ';' after loop condition"); err != nil {
		return nil, err
	}
	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	return ForStmt{Initializer: initializer, Condition: condition, Body: body}, nil
}

// whileStatement parses while condition statements
func (p *Parser) whileStatement() (Stmt, error) {
	if _, err := p.consume(scanner.LEFT_PAREN, "Expect '(' after 'while'"); err != nil {
		return nil, err
	}
	condition, err := p.expression()
	if err != nil {
		return nil, err
	}
	if _, err := p.consume(scanner.RIGHT_PAREN, "Expect ')' after while condition"); err != nil {
		return nil, err
	}

	body, err := p.statement()
	if err != nil {
		return nil, err
	}
	return WhileStmt{Condition: condition, Body: body}, nil
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

// ---------- Error handling ----------
type ParseError struct {
	Line int
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("[line %d] %s", e.Line, e.Message)
}

// errorAt reports an error at the current token.
func (p *Parser) errorAt(token scanner.Token, message string) error {
	return &ParseError{Line: token.Line, Message: message}
}

// consume checks that the current token is of the expected type, consumes it, and returns it. If not, it reports an error.
func (p *Parser) consume(tokenType scanner.TokenType, message string) (scanner.Token, error) {
	if p.check(tokenType) {
		return p.advance(), nil
	}
	return scanner.Token{}, p.errorAt(p.peek(), message)
} 

func (p *Parser) synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == scanner.SEMICOLON {
			return
		}
	switch p.peek().Type {
		case scanner.LET,
		scanner.PRINT,
		scanner.LEFT_BRACE,
		scanner.RIGHT_BRACE:
			return
		}
		p.advance()

	}
}

