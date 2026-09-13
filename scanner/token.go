package scanner

type TokenType string

const (
	LEFT_PAREN  TokenType = "LEFT_PAREN"
	RIGHT_PAREN TokenType = "RIGHT_PAREN"
	LEFT_BRACE  TokenType = "LEFT_BRACE"
	RIGHT_BRACE TokenType = "RIGHT_BRACE"
	SEMICOLON   TokenType = "SEMICOLON"

	PLUS  TokenType = "PLUS"
	MINUS TokenType = "MINUS"
	STAR  TokenType = "STAR"
	SLASH TokenType = "SLASH"

	LESS    TokenType = "LESS"
	GREATER TokenType = "GREATER"

	EQUAL         TokenType = "EQUAL"
	EQUAL_EQUAL   TokenType = "EQUAL_EQUAL"
	BANG_EQUAL    TokenType = "BANG_EQUAL"
	LESS_EQUAL    TokenType = "LESS_EQUAL"
	GREATER_EQUAL TokenType = "GREATER_EQUAL"

	IDENTIFIER TokenType = "IDENTIFIER"

	LET   TokenType = "LET"
	PRINT TokenType = "PRINT"
	IF    TokenType = "IF"
	ELSE  TokenType = "ELSE"
	FOR   TokenType = "FOR"
	TRUE  TokenType = "TRUE"
	FALSE TokenType = "FALSE"
	NONE  TokenType = "NONE"

	EOF TokenType = "EOF"
)

// Token represents a single lexical token produced by the scanner
type Token struct {
	Type    TokenType   // token category, e.g. "LEFT_PAREN", "NUMBER"
	Lexeme  string      // the raw source text that produced this token
	Literal interface{} // parsed literal value (e.g. number/string), nil if none
	Line    int         // source line number where the token appears
}
