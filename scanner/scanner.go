package scanner

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
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
		s.addToken(SLASH)
	case '<':
		s.addToken(LESS)
	case '>':
		s.addToken(GREATER)
	case '\n':
		s.line++
	case ' ', '\t', '\r':
		// ignore whitespace
	}
}
