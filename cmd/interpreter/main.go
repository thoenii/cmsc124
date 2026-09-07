package main

import (
	"cmsc124/scanner"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:] // grab CLI args, skip program name
	if len(args) == 2 && args[0] == "--tokenize" {
		runFile(args[1])
	} else {
		fmt.Println("Usage: ./run --tokenize <file>")
		os.Exit(64) // exit code 64: bad usage
	}
}

// runFile reads a source file, scans it into tokens, and prints them
func runFile(path string) {
	data, err := os.ReadFile(path)

	if err != nil {
		fmt.Fprintln(os.Stderr, "Could not read file:", err)
		os.Exit(65) // exit code 65: bad input file
	}

	s := scanner.NewScanner(string(data))
	tokens := s.ScanTokens()

	// print each token in a readable format
	for _, tok := range tokens {
		literal := "null"
		if tok.Literal != nil {
			literal = fmt.Sprintf("%v", tok.Literal) // format literal value if present
		}
		fmt.Printf("Token(type=%s, lexeme=%s, literal=%v, line=%d)\n", tok.Type, tok.Lexeme, literal, tok.Line)
	}
}
