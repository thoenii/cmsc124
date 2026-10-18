# Viper

## Creators

- Marie Toney Fay S. Gelvezon (@thoenii)
- Aleighia Keith L. Reyes (@akreyes9)

## Overview

Viper is a small, dynamically typed programming language designed as an educational interpreter project for CMSC 124. It is intended for students and beginner programmers who want a simple and familiar syntax for writing basic programs involving variables, expressions, conditions, loops, and output.

## Host language and build

- Host language: Go 1.26.7
- Version metadata: go.mod
- Build: `./build.sh`
- The interpreter can be built from a fresh clone using the provided build script.

## Running it

| Command | What it does |
|---|---|
| `./run <file>` | Executes a program. Available from Lab 4. |
| `./run --tokenize <file>` | Prints the token stream. |
| `./run --parse <file>` | Prints the parsed tree. |
| `./run --eval <file>` | Evaluates each expression and prints its value. Available from Lab 3. |
| `./run` | Starts the REPL. |

Usage:

​```
./run [--tokenize | --parse | --eval] <file>
​```

Exactly one flag and one file path are accepted. Any other combination prints the usage message to stderr and exits with code 64.

Exit codes:

| Code | Meaning |
|---|---|
| 0 | accepted, or executed successfully |
| 64 | bad usage, or the file could not be read |
| 65 | static error (lexical or syntax error) |
| 70 | runtime error |

## File extension

`.viper` (this must match the `ext` field in every `tests/lab*/manifest.json`).

## Lexical structure

### Keywords

| Keyword | Token | Purpose |
|---|---|---|
| fang | LET | declares a variable |
| hiss | PRINT | outputs a value |
| if | IF | conditional statement |
| else | ELSE | alternative for a conditional statement |
| coil | FOR | for-style loop (initializer, condition, increment) |
| slither | WHILE | loop that repeats while a condition holds |
| true | TRUE | boolean literal equal to true |
| false | FALSE | boolean literal equal to false |
| dead | NONE | absence of a value |

Keyword lexemes follow a snake theme. Internal token types (LET, PRINT, FOR, WHILE, NONE) stay the same regardless of how the keyword is spelled, since the token type identifies the grammar category and the lexeme is just the surface spelling.

### Operators

Precedence is listed from lowest (1) to highest (6).

| Operator | Category | Operands | Associativity | Precedence |
|---|---|---|---|---|
| = | assignment | binary | right | 1 |
| == | equality | binary | left | 2 |
| != | equality | binary | left | 2 |
| < | comparison | binary | left | 3 |
| <= | comparison | binary | left | 3 |
| > | comparison | binary | left | 3 |
| >= | comparison | binary | left | 3 |
| + | arithmetic | binary | left | 4 |
| - | arithmetic | binary | left | 4 |
| * | arithmetic | binary | left | 5 |
| / | arithmetic | binary | left | 5 |
| - (negation) | arithmetic | unary | right | 6 |

Viper has no increment operator. Use `x = x + 1`.

### Literals

| Kind | Syntax | Produces |
|---|---|---|
| Number | `42`, `3.14` | number value (double-precision float) |
| String | `"hello"`, escapes supported | string value |
| Boolean | `true`, `false` | boolean value |
| Nil | `dead` | absence of a value |

Numeric literals may be integers or decimals.
A leading decimal point is not allowed: `.5`
Instead, the value must be written as: `0.5`
A trailing decimal point is also not allowed: `5.`
Instead, it must be written as: `5.0`

Strings may span multiple lines.
The following escape sequences are supported:
- `\n` newline
- `\"` quotation mark
- `\\` backslash

Any other escape sequence (e.g. `\t`) is a scan-time error. The rest of the string is discarded up to its closing quote and scanning resumes after it.

### Identifiers

- Start characters: letters and `_`
- Continue characters: letters, digits, and `_`
- Case-sensitive: yes
- Reserved keywords cannot be used as identifiers

### Comments

- Line comments: `//`
- Block comments: `/* ... */`
- Nesting: not supported
- Harness note: `comment_prefix` in `tests/lab*/manifest.json` is set to `//`.

Comments are scanned and discarded rather than emitted as tokens.

## Whitespace and termination

- Whitespace significant: no, spaces, tabs, and newlines are discarded by the scanner.
- Statement terminator: semicolon `;`
- Block delimiters: braces `{ }`
- Grouping delimiters: parentheses `( )`

## Token output format

Each token prints on its own line:

```
Token(type=<TYPE>, lexeme=<text>, literal=<value or null>, line=<n>)
```

Example for `fang x = 42;`:

```
Token(type=LET, lexeme=fang, literal=null, line=1)
Token(type=IDENTIFIER, lexeme=x, literal=null, line=1)
Token(type=EQUAL, lexeme==, literal=null, line=1)
Token(type=NUMBER, lexeme=42, literal=42, line=1)
Token(type=SEMICOLON, lexeme=;, literal=null, line=1)
Token(type=EOF, lexeme=, literal=null, line=1)
```

Only NUMBER and STRING tokens carry a literal value. Every other token, including `true`, `false`, and `dead`, prints `literal=null`.

Token types:

- Grouping and punctuation: LEFT_PAREN, RIGHT_PAREN, LEFT_BRACE, RIGHT_BRACE, SEMICOLON
- Arithmetic operators: PLUS, MINUS, STAR, SLASH
- - Assignment and comparison: EQUAL, EQUAL_EQUAL, BANG, BANG_EQUAL, LESS, LESS_EQUAL, GREATER, GREATER_EQUAL
- Literals and identifiers: IDENTIFIER, STRING, NUMBER
- Keywords: LET, PRINT, IF, ELSE, FOR, WHILE, TRUE, FALSE, NONE
- End marker: EOF

The format is frozen as of Lab 1; changes are recorded in the changelog.

## Grammar

```
program        → declaration* EOF ;

declaration    → varDecl | statement ;
varDecl        → "fang" IDENTIFIER ( "=" expression )? ";" ;

statement      → matched | unmatched ;

matched        → "if" "(" expression ")" matched "else" matched
               | forHead matched
               | whileHead matched
               | block
               | printStmt
               | exprStmt ;

unmatched      → "if" "(" expression ")" statement
               | "if" "(" expression ")" matched "else" unmatched
               | forHead unmatched
               | whileHead unmatched ;

forHead        → "coil" "(" ( varDecl | exprStmt | ";" )
                 expression? ";" expression? ")" ;
whileHead      → "slither" "(" expression ")" ;
block          → "{" declaration* "}" ;
printStmt      → "hiss" "(" expression ")" ";" ;
exprStmt       → expression ";" ;

expression     → assignment ;
assignment     → IDENTIFIER "=" assignment | equality ;
equality       → comparison ( ( "==" | "!=" ) comparison )* ;
comparison     → term ( ( ">" | ">=" | "<" | "<=" ) term )* ;
term           → factor ( ( "+" | "-" ) factor )* ;
factor         → unary ( ( "*" | "/" ) unary )* ;
unary          → ( "!" | "-" ) unary | primary ;
primary        → NUMBER | STRING | "true" | "false" | "dead"
               | IDENTIFIER | "(" expression ")" ;
```

> **Note on `statement`:** the `matched` / `unmatched` split removes the dangling
> else ambiguity. An `else` can only follow a then-branch that is itself matched,
> so it always attaches to the nearest preceding `if` that has no `else`. The
> parser implements this by consuming `else` as soon as a then-branch ends.

## Parse output format

```
(+ 1.0 (* 2.0 3.0))
```

This is the `--parse` output for `1 + 2 * 3;`.

The parser returns a single `Program` root holding all top-level statements. `--parse` prints exactly one line per top-level statement, in source order. Nested statements (inside a block, `if`, `slither`, or `coil`) print inline within their parent's line. An empty program prints nothing and exits 0.

An expression statement prints as just its expression, with no wrapper and no trailing `;`. For example, `1 + 2;` prints as `(+ 1.0 2.0)`, while `hiss(1 + 2);` prints as `(print (+ 1.0 2.0))`.

- Numbers print as: at least one decimal place (`2` prints as `2.0`, `3.14` prints as `3.14`)
- Strings print as: their contents, without quotes
- Booleans print as: `true` / `false`
- Nil (`dead`) prints as: `nil`
- Groupings print as: `(group <expr>)`
- Unary expressions print as: `(<op> <expr>)`, e.g. `(- 5.0)`
- Binary expressions print as: `(<op> <left> <right>)`
- Assignments print as: `(= <name> <value>)`
- Statements print as: `(var <name> <init>)`, `(print <expr>)`, `(block ...)`, `(if ...)`, `(slither ...)`, `(coil ...)`
- Expression statements print as: `<expr>` (no wrapper)

## Semantics

Defined in Lab 3.

### Values and types

[To be completed in Lab 3: runtime values and their host representation.]

### Value printing

- Numbers: [TBD]
- Nil: [TBD]
- Strings: [TBD]

### Truthiness

[TBD]

### Operator semantics

- Arithmetic: [TBD]
- `+` on strings: [TBD]
- Mixed types: [TBD]
- Comparison: [TBD]
- Equality across types: [TBD]
- Division by zero: [TBD]

### Scope and bindings

[TBD, Lab 4]

### Control flow and functions

[TBD, Lab 4 and Lab 5]

## Native functions

None yet.

## Errors and diagnostics

Message format:

```
[line 1] Error: Unexpected character '@'
[line 1] Error at end: Expect ')' after expression.
[line 1] Error at ';': Expect expression.
```

The first is a lexical error from the scanner. The others are syntax errors from the parser, which name the offending token (or `end` at end of file). Runtime error format will be added in Lab 3.

| Failure | Exit code |
|---|---|
| lexical error | 65 |
| syntax error | 65 |
| runtime error | 70 |

### Syntax error behavior

When the parser rejects a file, it reports every error to stderr and prints nothing to stdout, so a rejected file produces no partial parse output. The process exits with code 65.

The parser does not stop at the first error. After a syntax error it skips tokens until it reaches a likely statement boundary (just after a `;`, or just before a statement keyword such as `fang`, `hiss`, `if`, `slither`, `coil`, or a brace) and resumes there, so one file can report several errors. To avoid cascading noise, a follow-up error on the same line as the previous reported error is suppressed.

If the scanner reports a lexical error, the file is rejected before parsing starts, and the same stdout and exit code rules apply.

## Testing conventions

| Folder | Activity | Mode | Flag |
|---|---|---|---|
| tests/lab1 | Scanner | sidecar | `--tokenize` |
| tests/lab2 | Parser | sidecar | `--parse` |
| tests/lab3 | Evaluator | inline | `--eval` |
| tests/lab4 | Context | inline | none |
| tests/lab5 | Functions | inline | none |

Run locally with:

```bash
curl -sSL https://raw.githubusercontent.com/WhiteLicorice/cmsc-124-harness/v1.1/run_tests.py -o run_tests.py
./build.sh
python3 run_tests.py tests/lab1
python3 run_tests.py tests/lab2
```

## Sample code

```
fang x = 1 + 2 * 3;
hiss(x);
```

`./run --parse` output:

```
(var x (+ 1.0 (* 2.0 3.0)))
(print x)
```

## Design rationale

The language was designed to be dynamically typed and intentionally small enough to implement as an interpreter prototype within three months.

## Known limitations

- Evaluation is implemented from Lab 3; statements are executed from Lab 4.
- Functions arrive in Lab 5.

## Changelog

| Activity | What changed in the language |
|---|---|
| Lab 1 | Defined the language identity, file format, lexical structure, token vocabulary, comments, whitespace rules, numeric literal rules, string rules, scanner errors, and token output format. |
| Lab 1 (theme and loops) | Reskinned keyword lexemes to a viper theme (fang, hiss, coil, dead). Added the WHILE token and slither keyword for while-style loops. Added the PLUS_PLUS token for the ++ operator. Implemented translation of string escape sequences (\n, \", \\), with a scan error for any unsupported escape. |
| Lab 2 | Defined the grammar and operator precedence (assignment, equality, comparison, term, factor, unary). Added variable declarations, `hiss(...)` printing, blocks, if/else, `slither` and `coil` loops. Defined the `--parse` output format. The parser now returns a single `Program` root node. |
| Lab 2 | Defined the grammar and operator precedence (assignment, equality, comparison, term, factor, unary). Added variable declarations, `hiss(...)` printing, blocks, if/else, `slither` and `coil` loops. Defined the `--parse` output format. The parser now returns a single `Program` root node. Removed the `++` operator and the `PLUS_PLUS` token; `++` now scans as two `PLUS` tokens. |