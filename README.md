# [Language name]

## Creators

- Marie Toney Fay S. Gelvezon (@thoenii)
- Aleighia Keith L. Reyes (@akreyes9)

## Overview

Viper is a small, dynamically typed programming language designed as an educational interpreter project for CMSC 124. It is intended for students and beginner programmers who want a simple and familiar syntax for writing basic programs involving variables, expressions, conditions, loops, and output.

## Host language and build

- Host language: Go 1.26.7
- Version metadata: go.mod
- Build: ./build.sh
- The interpreter can be built from a fresh clone using the provided build script. 

## Running it

| Command | What it does |
|---|---|
| `./run <file>` | [Executes a program. Available from Lab 4.] |
| `./run --tokenize <file>` | [Prints the token stream.] |
| `./run --parse <file>` | [Prints the parsed tree.] |
| `./run --eval <file>` | [Evaluates each expression and prints its value.] |
| `./run` | [Starts the REPL.] |


Exit codes: 0 for successful execution, 65 for static errors (including lexical or syntax errors), and 70 for runtime errors.

## File extension

`.src` The file extension must match the ext field in every tests.

## Lexical structure

### Keywords

| Keyword | Purpose |
|---|---|
| let | declares a variable |
| print | outputs a value |
| if | conditional statement |
| else | alternative for a conditional statement |
| for | creates a loop |
| true | boolean literal equal to true |
| false | boolean literal equal to false |
| none | absence of values |


### Operators


| Operator | Category | Operands | Associativity | Precedence |
|---|---|---|---|---|
| = | assignment | binary | right | [TBD] |
| == | comparison | binary | left | [TBD] |
| != | comparison | binary | left | [TBD] |
| < | comparison | binary | left | [TBD] |
| <= | comparison | binary | left | [TBD] |
| > | comparison | binary | left | [TBD] |
| >= | comparison | binary | left | [TBD] |
| + | arithmetic | binary | left | [TBD] |
| - | arithmetic | binary | left | [TBD] |
| * | arithmetic | binary | left | [TBD] |
| / | arithmetic | binary | left | [TBD] |
| ++ | increment | unary | left | [TBD] |

### Literals


| Kind | Syntax | Produces |
|---|---|---|
| Number | [e.g. 42, 3.14] | [int, float, number value] |
| String | [e.g. "hello", escapes supported] | [string value] |
| Boolean | [true, false] | [boolean value] |
| Nil | [spelling] | [absence of values] |

Numeric literals may be integers or decimals.
A leading decimal point is not allowed: .5
Instead, the value must be written as: 0.5
A trailing decimal point is also not allowed: 5.
Instead, it must be written as: 5.0

Strings may span multiple lines.
The ff. escape sequences are supported:
\n newline
\” quotation mark
\\ backlash

### Identifiers

- Start characters: alphabets and _
- Continue characters: alphabets, digits, and _
- Case-sensitive: yes
- Reserved keywords cannot be used as identifiers

### Comments

- Line comments: //
- Block comments: /* … */
- Nesting: not supported
- [Harness note: comment_prefix in tests/lab*/manifest.json is set to the
  token above.]

Comments are scanned and discarded rather than emitted as tokens.

## Whitespace and termination

- Whitespace significant: no, spaces and tabs are discarded by the scanner.
- Statement terminator: semicolon ;
- Block delimiters: braces { }
- Grouping delimiters: parentheses ( )

## Token output format

```
Example:

LET let null 1
IDENTIFIER x null 1
EQUAL = null 1
NUMBER 42 42 1
SEMICOLON ; null 1
EOF  null 1 
```
Tokens 

Grouping and punctuation
-LEFT_PAREN
-RIGHT_PAREN
-LEFT_BRACE
-RIGHT_BRACE
-SEMICOLON
Arithmetic operators
-PLUS
-MINUS
-STAR
-SLASH
-PLUS_PLUS
Assignment and comparison operators
-EQUAL
-EQUAL_EQUAL
-BANG_EQUAL
-LESS
-LESS_EQUAL
-GREATER
-GREATER_EQUAL
Literals and identifiers
-IDENTIFIER
-STRING
-NUMBER
Keywords
-LET
-PRINT
-IF
-ELSE
-FOR
-TRUE
-FALSE
-NONE
End marker
-EOF
[What each field means. Frozen as of Lab 1; changes are recorded in the
changelog.]

## Grammar

```
[Your complete context-free grammar, current as of the latest activity.
Unambiguous, with precedence and associativity encoded in rule structure.]
```

## Parse output format

```
[one line of real --parse output, e.g. (+ 1.0 (* 2.0 3.0))]
```

- Groupings print as: [form]
- Numbers print as: [form]

## Semantics

### Values and types

[What runtime values exist, and how they are represented in the host
language.]

### Value printing

- Numbers: [e.g. 5 rather than 5.0]
- Nil: [spelling]
- Strings: [with or without quotes]

### Truthiness

[The complete rule. Which values are false in a condition; everything else is
true.]

### Operator semantics

- Arithmetic: [accepted operand types]
- `+` on strings: [concatenation, error, or coercion]
- Mixed types: [what happens]
- Comparison: [accepted operand types]
- Equality across types: [false, or an error]
- Division by zero: [value produced, or runtime error]

### Scope and bindings

- Redeclaration in the same scope: [allowed or an error]
- Uninitialized variable holds: [value]
- Shadowing: [behavior]
- Undefined name: [static error with exit 65, or runtime error with exit 70]

### Control flow and functions

- Logical operators return: [booleans, or the operand]
- Dangling else binds to: [which if]
- Closure capture of a loop variable: [per iteration, or shared]
- Function with no return statement produces: [value]
- Arity mismatch: [message and exit code]

## Native functions


| Name | Arguments | Returns | Notes |
|---|---|---|---|
| [name] | [count and types] | [type] | [caveats] |


## Errors and diagnostics

Message format:

```
[one real static error]
[one real runtime error]
```


| Failure | Exit code |
|---|---|
| [lexical error] | 65 |
| [syntax error] | 65 |
| [runtime error] | 70 |


## Testing conventions


| Folder | Activity | Mode | Flag |
|---|---|---|---|
| tests/lab1 | Scanner | sidecar | `--tokenize` |
| tests/lab2 | Parser | sidecar | `--parse` |
| tests/lab3 | Evaluator | inline | `--eval` |
| tests/lab4 | Context | inline | none |
| tests/lab5 | Functions | inline | none |


```
[specific tests]...
```

Run locally with:

```bash
curl -sSL https://raw.githubusercontent.com/WhiteLicorice/cmsc-124-harness/v1.1/run_tests.py -o run_tests.py
./build.sh
python3 run_tests.py tests/lab1
```

## Sample code

```
[a short program]
```

Output:

```
[its output]
```

## Design rationale

The language was designed to be dynamically typed and intentionally small enough to implement as an interpreter prototype within three months. 

## Known limitations

- The language prototype currently supports only the lexical structure defined in lab 1.
- Parsing will be implemented in Lab 2.
- Further functions will be implemented til the next Laboratory activities.

## Changelog


| Activity | What changed in the language |
|---|---|
| Lab 1 | Defined the language identity, file format, lexical structure, token vocabulary, comments, whitespace rules, numeric literal rules, string rules, scanner errors, and token output format.  |

