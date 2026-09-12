package main

import (
	"bufio"
	"cmsc124/scanner"
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:] // grab CLI args, skip program name

	switch {
	case len(args) == 0:
		// no arguments: start the interactive REPL
		runRepl()

	case len(args) == 2 && args[0] == "--tokenize":
		// tokenize flag: scan one file and print its token stream
		runFile(args[1])

	case len(args) == 1:
		// bare "./run <path>": preserve Lab 0's original contract unchanged,
		// since every lab folder ever committed has to keep passing in CI
		fmt.Println("Hello from CMSC 124!")
		os.Exit(0)

	default:
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
		var literal interface{} = "null"
		if tok.Literal != nil {
			literal = tok.Literal
		}
		fmt.Printf("Token(type=%s, lexeme=%s, literal=%v, line=%d)\n", tok.Type, tok.Lexeme, literal, tok.Line)
	}
}

/* runRepl reads one line at a time from stdin, scans it, and prints its
   tokens. A bad line must not kill the session, so errors from a single
   line are reported and the loop continues to the next prompt.
*/

func runRepl() {
	scannerInput := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")

		if !scannerInput.Scan() {
			// EOF (e.g. Ctrl+D) or read error: end the session cleanly
			break
		}

		line := scannerInput.Text()

		s := scanner.NewScanner(line)
		tokens := s.ScanTokens()

		// print each token in a readable format, same as runFile
		for _, tok := range tokens {
			var literal interface{} = "null"
			if tok.Literal != nil {
				literal = tok.Literal
			}
			fmt.Printf("Token(type=%s, lexeme=%s, literal=%v, line=%d)\n", tok.Type, tok.Lexeme, literal, tok.Line)
		}
	}
}
