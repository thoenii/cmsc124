package scanner

// Token represents a single lexical token produced by the scanner
type Token struct {
	Type    string      // token category, e.g. "LEFT_PAREN", "NUMBER"
	Lexeme  string      // the raw source text that produced this token
	Literal interface{} // parsed literal value (e.g. number/string), nil if none
	Line    int         // source line number where the token appears
}
