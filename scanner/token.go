package scanner

type TokenType string

const (
	LEFT_PAREN TokenType = "LEFT_PAREN"
	RIGHT_PAREN TokenType = "RIGHT_PAREN"
	LEFT_BRACE TokenType = "LEFT_BRACE"
	RIGHT_BRACE TokenType = "RIGHT_BRACE"
	SEMICOLON TokenType = "SEMICOLON"

	PLUS TokenType = "PLUS"
	MINUS TokenType = "MINUS"
	STAR TokenType = "STAR"
	SLASH TokenType = "SLASH"

	EQUAL TokenType = "EQUAL"
	EQUAL_EQUAL TokenType = "EQUAL_EQUAL"

	BANG TokenType = "BANG"
	BANG_EQUAL TokenType = "BANG_EQUAL"

	LESS TokenType = "LESS"
	LESS_EQUAL TokenType = "LESS_EQUAL"

	GREATER TokenType = "GREATER"
	GREATER_EQUAL TokenType = "GREATER_EQUAL"

	IDENTIFIER TokenType = "IDENTIFIER"
	NUMBER TokenType = "NUMBER"
	STRING TokenType = "STRING"
	
	LET TokenType = "LET"
	PRINT TokenType = "PRINT"
	IF TokenType = "IF"
	ELSE TokenType = "ELSE"
	WHILE TokenType = "WHILE"
	TRUE TokenType = "TRUE"
	FALSE TokenType = "FALSE"
	NONE TokenType = "NONE"

	EOF TokenType = "EOF"
)

type Token struct {
	Type TokenType
	Lexeme string
	Literal interface{}
	Line int
}

func NewToken(tokenType TokenType, lexeme string, literal interface{}, line int) *Token {
	return &Token{
		Type: tokenType,
		Lexeme: lexeme,
		Literal: literal,
		Line: line,
	}
}
