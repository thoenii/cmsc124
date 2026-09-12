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

// runFile reads a source file, scans it into tokens, and prints them.
// Exits 65 if the scanner reported any error, 0 on a fully clean scan.
func runFile(path string) {
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file %q: %v\n", path, err)
		os.Exit(64)
	}

	sc := scanner.NewScanner(string(source))
	tokens := sc.ScanTokens()

	// print each token in the frozen Token(...) format — same shape as
	// runRepl() below, kept inline here rather than as a shared method
	// since Token has no String() method defined
	for _, tok := range tokens {
		var literal interface{} = "null"
		if tok.Literal != nil {
			literal = tok.Literal
		}
		fmt.Printf("Token(type=%s, lexeme=%s, literal=%v, line=%d)\n", tok.Type, tok.Lexeme, literal, tok.Line)
	}

	// only exit 65 once the whole file has been scanned, per the "keep
	// scanning after an error" contract — never bail out mid-scan
	if sc.HadError() {
		os.Exit(65)
	}
	os.Exit(0)
}

/*
runRepl reads one line at a time from stdin, scans it, and prints its

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
