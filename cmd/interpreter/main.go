package main

import (
	"bufio"
	"cmsc124/evaluator"
	"cmsc124/parser"
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

	case len(args) == 2 && args[0] == "--parse":
		// for parsing
		runParseFile(args[1])

	case len(args) == 2 && args[0] == "--eval":
		runEvalFile(args[1])

	case len(args) == 1:
		runProgramFile(args[0])

	default:
		fmt.Fprintln(os.Stderr, "Usage: ./run [--tokenize | --parse | --eval] <file>")
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

	// print each token in the frozen Token(...) format, kept inline here
	// rather than as a shared method since Token has no String() method defined
	for _, tok := range tokens {
		var literal interface{} = "null"
		if tok.Literal != nil {
			literal = tok.Literal
		}
		fmt.Printf("Token(type=%s, lexeme=%s, literal=%v, line=%d)\n", tok.Type, tok.Lexeme, literal, tok.Line)
	}

	// only exit 65 once the whole file has been scanned, per the "keep
	// scanning after an error" contract - never bail out mid-scan
	if sc.HadError() {
		os.Exit(65)
	}
	os.Exit(0)
}

// runParseFile scans a file, parses its tokens, and prints the AST.
func runParseFile(path string) {
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file %q: %v\n", path, err)
		os.Exit(64)
	}

	sc := scanner.NewScanner(string(source))
	tokens := sc.ScanTokens()

	// Stop if scanning encountered an error.
	if sc.HadError() {
		os.Exit(65)
	}

	p := parser.NewParser(tokens)
	prog, err := p.ParseProgram()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(65)
	}

	// An empty program must print nothing, not a blank line.
	if out := parser.PrintProgram(prog); out != "" {
		fmt.Println(out)
	}
	os.Exit(0)
}

func execStmt(interp *evaluator.Interpreter, stmt parser.Stmt) error {
	// bare expression: evaluate and print the value (REPL only)
	if s, ok := stmt.(parser.ExprStmt); ok {
		val, err := interp.Evaluate(s.Expression)
		if err != nil {
			return err
		}
		fmt.Println(evaluator.Stringify(val))
		return nil
	}

	// everything else: Exec handles it (hiss prints inside Exec)
	return interp.Exec(stmt)
}

/*
runRepl reads one line at a time from stdin, scans it, parses it, and prints
the value of each statement. A bad line must not kill the session, so errors
from a single line are reported and the loop continues to the next prompt.
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
		if s.HadError() {
			continue // scanner already reported it; back to the prompt
		}

		p := parser.NewParser(tokens)
		prog, err := p.ParseProgram()
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			continue
		}

		interpreter := evaluator.NewInterpreter()

		for _, stmt := range prog.Statements {
			if err := execStmt(interpreter, stmt); err != nil {
				fmt.Fprintln(os.Stderr, err)
				break // stop this line, but keep the session alive
			}
		}
	}
}

// runEvalFile scans, parses, then evaluates each expression statement and
// prints its value. Static errors exit 65; runtime errors exit 70.
func runEvalFile(path string) {
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file %q: %v\n", path, err)
		os.Exit(64)
	}

	sc := scanner.NewScanner(string(source))
	tokens := sc.ScanTokens()
	if sc.HadError() {
		os.Exit(65)
	}

	p := parser.NewParser(tokens)
	prog, err := p.ParseProgram()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(65)
	}

	interpreter := evaluator.NewInterpreter()

	for _, stmt := range prog.Statements {
		var err error
		switch s := stmt.(type) {
		case parser.ExprStmt:
			// Lab 3 behavior: a bare expression prints its value
			var val evaluator.Value
			val, err = interpreter.Evaluate(s.Expression)
			if err == nil {
				fmt.Println(evaluator.Stringify(val))
			}
		default:
			// hiss, blocks, if, slither
			err = interpreter.Exec(stmt)
		}
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(70)
		}
	}
	os.Exit(0)
}

func runProgramFile(path string) {
	source, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file %q: %v\n", path, err)
		os.Exit(64)
	}

	sc := scanner.NewScanner(string(source))
	tokens := sc.ScanTokens()
	if sc.HadError() {
		os.Exit(65)
	}

	p := parser.NewParser(tokens)
	prog, err := p.ParseProgram()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(65)
	}

	interpreter := evaluator.NewInterpreter()

	for _, stmt := range prog.Statements {
		if err := interpreter.Exec(stmt); err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(70)
		}
	}
	os.Exit(0)
}
