package scanner

import (
	"fmt"
	"os"
	"unicode"
)

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
}

var keywords = map[string]TokenType{
	"let":   LET,
	"print": PRINT,
	"if":    IF,
	"else":  ELSE,
	"for":   FOR,
	"true":  TRUE,
	"false": FALSE,
	"none":  NONE,
}

// NewScanner creates a Scanner for the given source string, starting at line 1
func NewScanner(source string) *Scanner {
	return &Scanner{source: source, line: 1}
}

// ScanTokens scans the entire source and returns the resulting list of tokens
func (s *Scanner) ScanTokens() []Token {
	for !s.isAtEnd() {
		s.start = s.current
		s.scanToken()
	}
	s.tokens = append(s.tokens, Token{Type: EOF, Lexeme: "", Literal: nil, Line: s.line})
	return s.tokens
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(rune(s.peek())) {
		s.advance()
	}
	text := s.source[s.start:s.current]

	if tokenType, ok := keywords[text]; ok {
		s.addToken(tokenType)
	} else {
		s.addToken(IDENTIFIER)
	}
}

// isAtEnd reports whether the scanner has consumed all source characters
func (s *Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

// advance consumes and returns the current character, moving the cursor forward
func (s *Scanner) advance() byte {
	c := s.source[s.current]
	s.current++
	return c
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() {
		return false
	} else if s.source[s.current] != expected {
		return false
	} else {
		s.current++
		return true
	}
}

func (s *Scanner) peek() byte {
	if s.isAtEnd() {
		return 0
	} else {
		return s.source[s.current]
	}
}

func (s *Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return 0
	} else {
		return s.source[s.current+1]
	}
}

func (s *Scanner) blockComment() {
	for !s.isAtEnd() {
		if s.peek() == '*' && s.peekNext() == '/' {
			s.advance()
			s.advance()
			return
		} else if s.peek() == '\n' {
			s.line++
			s.advance()
		} else {
			s.advance()
		}
	}
	// TODO: replace with s.reportError(s.line, "Unterminated block comment")
	// once reportError() is implemented, so this also sets the had-error
	// flag and gets picked up by main.go's exit-65 check.
	fmt.Fprintf(os.Stderr, "[line %d] Error: Unterminated block comment\n", s.line)
}

func isAlpha(c rune) bool {
	return c == '_' || unicode.IsLetter(c)
}

func isAlphaNumeric(c rune) bool {
	return c == '_' || unicode.IsLetter(c) || unicode.IsDigit(c)
}

// addToken creates a token from source[start:current] and appends it to tokens
func (s *Scanner) addToken(tokenType TokenType) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{Type: tokenType, Lexeme: text, Literal: nil, Line: s.line})
}

// scanToken consumes one character and produces the matching token, if any
func (s *Scanner) scanToken() {
	c := s.advance()
	switch c {
	case '(':
		s.addToken(LEFT_PAREN)
	case ')':
		s.addToken(RIGHT_PAREN)
	case '{':
		s.addToken(LEFT_BRACE)
	case '}':
		s.addToken(RIGHT_BRACE)
	case ';':
		s.addToken(SEMICOLON)
	case '+':
		s.addToken(PLUS)
	case '-':
		s.addToken(MINUS)
	case '*':
		s.addToken(STAR)
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			s.blockComment()
		} else {
			s.addToken(SLASH)
		}
	case '=':
		if s.match('=') {
			s.addToken(EQUAL_EQUAL)
		} else {
			s.addToken(EQUAL)
		}
	case '<':
		if s.match('=') {
			s.addToken(LESS_EQUAL)
		} else {
			s.addToken(LESS)
		}
	case '>':
		if s.match('=') {
			s.addToken(GREATER_EQUAL)
		} else {
			s.addToken(GREATER)
		}
	case '\n':
		s.line++
	case ' ', '\t', '\r':
		// ignore whitespace
	default:
		if isAlpha(rune(c)) {
			s.identifier()
		} else {
			// TODO: replace with s.reportError(s.line, fmt.Sprintf("Unexpected character '%c'", c))
			// once reportError() is implemented, so this also sets the had-error
			// flag and gets picked up by main.go's exit-65 check.
			fmt.Fprintf(os.Stderr, "[line %d] Error: Unexpected character '%c'\n", s.line, c)
		}
	}
}
