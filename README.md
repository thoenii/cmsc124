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
| `./run <file>` | Executes a program: statements run in order, and only `hiss` produces output. Available from Lab 4. |
| `./run --tokenize <file>` | Prints the token stream. |
| `./run --parse <file>` | Prints the parsed tree. |
| `./run --eval <file>` | Evaluates each expression and prints its value. Available from Lab 3. |
| `./run` | Starts the REPL. Variables survive between lines. |

Usage:

```
./run <file>
./run [--tokenize | --parse | --eval] <file>
./run
```

With a flag, exactly one flag and one file path are accepted. Without a flag, exactly one file path is accepted. With no arguments, the REPL starts. Any other combination prints the usage message to stderr and exits with code 64.

Exit codes:

| Code | Meaning |
|---|---|
| 0 | accepted, or executed successfully |
| 64 | bad usage, or the file could not be read |
| 65 | static error (lexical or syntax error). Nothing was executed. |
| 70 | runtime error. The program started and then failed. |

Diagnostics go to stderr. Program output goes to stdout.

### Script mode and REPL

| Mode | What gets printed |
|---|---|
| `./run <file>` (script mode) | Only `hiss` prints. An expression statement such as `1 + 2;` is evaluated and its value is thrown away, so nothing is shown. |
| REPL | A bare expression typed at the prompt prints its value. `hiss(...)` prints once, not twice. |
| `./run --eval <file>` | Unchanged from Lab 3: each expression's value is printed. |

Script mode and the REPL share one evaluator. File mode creates one interpreter, runs the program, and exits with the right code. The REPL creates one interpreter and feeds it lines, so its variables persist.

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
- Assignment and comparison: EQUAL, EQUAL_EQUAL, BANG, BANG_EQUAL, LESS, LESS_EQUAL, GREATER, GREATER_EQUAL
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

### Where each form is legal

| Form | Where it is allowed |
|---|---|
| `fang` declaration | Only as a `declaration`: at the top level of a program, directly inside a block, or as the initializer slot of a `coil` header. It is **not** a statement, so it cannot be the body of an `if`, `slither`, or `coil` without a block. `if (x) fang y = 1;` is a syntax error. |
| `hiss(...)` | Anywhere a statement is allowed. It needs parentheses and a `;`. |
| Expression statement | Anywhere a statement is allowed. It needs a `;`. |
| Block `{ ... }` | Anywhere a statement is allowed. It nests to any depth. Each block is its own scope. |
| Assignment `x = value` | Anywhere an expression is allowed. It is right-associative, so `a = b = 3;` works, and it produces the assigned value. The left side must be a plain variable name. `1 + 2 = 3;` is a syntax error. |
| Variable read | Anywhere an expression is allowed. |

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

> `--parse` is a debug dump of the tree and its format has been frozen since Lab 2, which is why it prints `2.0` and `nil`. Runtime output (`--eval`, `hiss`, the REPL) follows the rules under **Value printing** below instead.

## Semantics

Defined in Lab 3 and Lab 4.

### Values and types

Viper is dynamically typed: a value's type is only known at runtime. There are four kinds of value:

| Viper value | Go representation |
|---|---|
| number | `float64` |
| string | `string` |
| boolean | `bool` |
| `dead` (nil) | Go `nil` |

The host type is `type Value = any` (in `evaluator/value.go`). Go has no single type that holds all four, so `any` is the uniform storage, a type switch (`v.(type)`) is the runtime interrogation, and a type assertion (`v.(float64)`) converts a value back to a host number for arithmetic.

There is **one numeric type**, a double-precision float. There is no separate integer type, so `7 / 2` is `3.5`.

### Value printing

All runtime output goes through a single function, `evaluator.Stringify`, so `--eval`, `hiss` and the REPL always format a value the same way.

- Numbers: a trailing `.0` is trimmed, so `2 + 3` prints `5` and `2.5` prints `2.5` (same as jlox).
- Nil: prints `dead`, the language's own spelling. (`--parse` prints `nil` because that format was frozen in Lab 2.)
- Strings: print without quotes, so `"hi" + " there"` prints `hi there`.
- Booleans: print `true` / `false`.

### Truthiness

`false` and `dead` are falsey. Everything else is truthy, including `0` and the empty string `""` (the same rule Lox uses). So `!dead` is `true`, `!false` is `true`, and `!0`, `!""` and `!"hello"` are all `false`.

### Operator semantics

- Arithmetic: `-`, `*`, `/` and unary `-` work on numbers only. Any other operand type is a runtime error. There is no implicit conversion, so arithmetic with a string is an error.
- `+` on strings: `+` works on two numbers (addition) or two strings (concatenation). Mixed operands, like `"a" + 1`, are a runtime error: `Operands must be two numbers or two strings.`
- Mixed types: there is no coercion anywhere. A number never becomes a string, and a string never becomes a number.
- Comparison: `<`, `<=`, `>`, `>=` work on numbers only and return a boolean. Comparing a string with a number (or anything that isn't two numbers) is a runtime error.
- Equality across types: `==` and `!=` work on any pair of values and never raise an error. Values of different types are never equal, so `1 == "1"` is `false` and `1 != "1"` is `true`. Strings compare by content, so `"ab" == "ab"` is `true`. `dead == dead` is `true`.
- Division by zero: a runtime error (`Division by zero.`), exit code 70. This also covers `0 / 0`. We chose to reject it, instead of inheriting IEEE 754's `Infinity`, so that no infinity or NaN value ever exists in the language, and `Stringify` never has to print one.
- Evaluation order: left to right. For a binary expression the left operand is evaluated first, then the right, then the operator is applied. Operands are always evaluated before their types are checked.
- Type checks: every operator that needs numbers uses one shared helper, so the check and its error message are written once.

### Statements and output

- `fang` declares a variable (see **Scope and bindings**).
- `hiss(expr);` evaluates the expression and prints its value on its own line, using `Stringify`.
- An expression statement evaluates its expression and discards the result. This matters once assignment exists, because `x = 5;` does its work through a side effect.
- A block runs its declarations in order inside a new scope.
- Statements run top to bottom. In file mode, output printed before a runtime error stays printed.

### Scope and bindings

Defined in Lab 4.

**Environments.** Bindings live in environments. An environment holds a table of name to value and a link to the environment that encloses it. The global environment has no enclosing link. Entering a block creates a new environment linked to the current one. Leaving the block restores the previous one, even if a runtime error happened inside the block.

**Scope rule: lexical (static).** The links follow source nesting: a block's environment is enclosed by the environment the block is written inside. Blocks alone cannot tell lexical and dynamic scope apart. Function calls will, and that decision is recorded in Lab 5.

**Which binding wins.** When a name exists in several environments, the innermost one wins. A lookup checks the current environment first, then follows the enclosing links outward. If no environment has the name, it is an error.

**The three operations.**

| Operation | What it does |
|---|---|
| declare (`fang x = v;`) | Creates the binding in the **current** environment only. It never looks outward. An inner declaration shadows an outer one without changing it. |
| read (`x`) | Walks outward from the current environment and returns the first match. |
| assign (`x = v`) | Walks outward and changes the first existing match. It **never creates** a binding. |

Declaring creates and assigning updates. This is why `fang x` inside a block shadows the outer `x`, while `x = 1` inside a block with no local `x` changes the outer one. A misspelled name in an assignment is an error, not a new variable.

**Redeclaration.** Declaring a name that already exists in the same scope is allowed. The new declaration overwrites the old value in that scope.

**Uninitialized variables.** `fang x;` declares `x` without a value. Reading `x` before it has been given a value is a runtime error (decided in Lab 3). Assigning to it gives it a value.

**Undefined variables.** Reading or assigning a name that was never declared is a **runtime** error: the message has a line number on stderr, and file mode exits 70. In the REPL the error is printed and the prompt returns. We chose runtime over a static check so that recursive and mutually recursive functions can work in Lab 5 without a separate resolution pass.

**Visibility, storage, and mutability.**

| Binding | Becomes visible | Storage begins | Storage ends | Value or type can change? |
|---|---|---|---|---|
| global `fang` | after its declaration finishes running | when the declaration runs | when the program (or REPL session) ends | value: yes, by assignment. Type: yes (dynamically typed). |
| `fang` inside a block | after its declaration finishes running | when the declaration runs | when the block ends | value: yes, by assignment. Type: yes (dynamically typed). |

Scope (where a name can be seen) and lifetime (how long its storage exists) are the same for blocks, since Viper has no functions yet. Block variables are stack-dynamic: they last for one run of the block. Global variables last for the whole program. Lab 5 revisits this if functions can keep a scope alive after it ends. The interpreter stores environments in a Go map, but that is an implementation detail and not something the language promises.

### Control flow and functions

[TBD, Lab 5]

## Native functions

None yet.

## Errors and diagnostics

Message format:

```
[line 1] Error: Unexpected character '@'
[line 1] Error at end: Expect ')' after expression.
[line 1] Error at ';': Expect expression.
[line 2] Error at 'print': Expect ';' after variable declaration.
[line 3] Error at '=': Invalid assignment target.
[line 1] Runtime error: Operands must be numbers.
[line 1] Runtime error: Undefined variable 'a'.
```

The first is a lexical error from the scanner. The next four are syntax errors from the parser, which name the offending token (or `end` at end of file). The last two are runtime errors from the evaluator.

| Failure | Exit code |
|---|---|
| lexical error | 65 |
| syntax error | 65 |
| runtime error | 70 |

### Syntax error behavior

When the parser rejects a file, it reports every error to stderr and prints nothing to stdout, so a rejected file produces no partial parse output. **Nothing is executed.** The process exits with code 65.

The parser does not stop at the first error. After a syntax error it skips tokens until it reaches a likely statement boundary (just after a `;`, or just before a statement keyword such as `fang`, `hiss`, `if`, `slither`, `coil`, or a brace) and resumes there, so one file can report several errors. Synchronization always consumes at least one token, so it cannot loop forever. To avoid cascading noise, a follow-up error on the same line as the previous reported error is suppressed.

Syntax errors added in Lab 4:

| Situation | Message |
|---|---|
| missing `;` after a variable declaration | `Expect ';' after variable declaration.` |
| missing `;` after an expression statement | `Expect ';' after expression.` |
| missing `;` after a `hiss` statement | `Expect ';' after value.` |
| missing `}` at the end of a block | `Expect '}' after block.` |
| left side of `=` is not a variable name | `Invalid assignment target.` |

If the scanner reports a lexical error, the file is rejected before parsing starts, and the same stdout and exit code rules apply.

### Runtime error behavior

A runtime error is a program that was accepted by the scanner and parser but hit a bad value while running. The evaluator stops and returns the error. It carries the offending token, so the message always includes a line number:

```
[line 3] Runtime error: Operands must be numbers.
```

Messages are written to stderr:

| Situation | Message |
|---|---|
| `-`, `*`, `/`, `<`, `<=`, `>`, `>=` on non-numbers | `Operands must be numbers.` |
| unary `-` on a non-number | `Operand must be a number.` |
| `+` on anything but two numbers or two strings | `Operands must be two numbers or two strings.` |
| division by zero | `Division by zero.` |
| reading or assigning a name that was never declared | `Undefined variable 'x'.` |
| reading a variable declared without a value and never assigned | `Variable 'x' is not initialized.` |

In file mode, execution stops at the first runtime error. Output printed before the error stays printed, and the process exits 70. In the REPL, the error is printed and the prompt returns. The REPL's variables survive the error, and the active scope is restored, so a runtime error inside a block does not leave the session stuck inside that block. Both modes share the same evaluator and differ only in what they do with an error.

A host-language panic or stack trace must never reach the user. Wrong programs produce a runtime error, not a crash.

## Testing conventions

| Folder | Activity | Mode | Flag |
|---|---|---|---|
| tests/lab0 | Hello program | sidecar | none |
| tests/lab1 | Scanner | sidecar | `--tokenize` |
| tests/lab2 | Parser | sidecar | `--parse` |
| tests/lab3 | Evaluator | inline | `--eval` |
| tests/lab4 | Context | inline | none |
| tests/lab5 | Functions | inline | none |

`tests/lab0/hello.viper` is now a real Viper program that prints what `hello.expected` says.

Run locally with:

```bash
curl -sSL https://raw.githubusercontent.com/WhiteLicorice/cmsc-124-harness/v1.1/run_tests.py -o run_tests.py
./build.sh
python3 run_tests.py tests/lab0
python3 run_tests.py tests/lab1
python3 run_tests.py tests/lab2
python3 run_tests.py tests/lab3
python3 run_tests.py tests/lab4
```

Inline tests (Lab 3 onward) put the expected result in a comment on the same line as the code. Each error test file holds exactly one error annotation. Expected values are worked out by hand, not copied from the interpreter's output.

Inline annotations:

| Annotation | Checks |
|---|---|
| `// expect: value` | the next stdout line, in order |
| `// expect nothing` | the file intentionally prints nothing |
| `// expect runtime error: message` | a substring of stderr, and exit 70 |
| `// expect error: message` | a substring of stderr, and exit 65 |

Annotate the message text only, never the `[line N]` prefix.

Lab 4 tests cover: shadowing, nested assignment through two scopes, undefined read, undefined assignment, invalid assignment target, missing `;`, several syntax errors in one file, uninitialized read, and a REPL-style check that a runtime error inside a block does not break the scope afterward.

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

Scope example:

```
fang x = "outer";
{
  fang x = "inner";
  hiss(x);
}
hiss(x);
```

`./run` output:

```
inner
outer
```

## Design rationale

The language was designed to be dynamically typed and intentionally small enough to implement as an interpreter prototype within three months.

Lab 3 decisions favor strictness over convenience: no implicit conversions, `+` never mixes types, and division by zero is an error. A program that does something surprising fails loudly with a line number instead of quietly producing a value.

Lab 4 continues this. Assignment never creates a variable, so a typo is an error instead of a silent new variable. Undefined names are caught at runtime so that Lab 5 functions can call each other without a separate resolution pass.

## Known limitations

- `if`, `slither` and `coil` are parsed but not executed until Lab 5.
- Functions arrive in Lab 5.
- `--eval` keeps its Lab 3 behavior and does not run declarations or blocks. Use plain `./run` for full programs.

## Changelog

| Activity | What changed in the language |
|---|---|
| Lab 1 | Defined the language identity, file format, lexical structure, token vocabulary, comments, whitespace rules, numeric literal rules, string rules, scanner errors, and token output format. |
| Lab 1 (theme and loops) | Reskinned keyword lexemes to a viper theme (fang, hiss, coil, dead). Added the WHILE token and slither keyword for while-style loops. Added the PLUS_PLUS token for the ++ operator. Implemented translation of string escape sequences (\n, \", \\), with a scan error for any unsupported escape. |
| Lab 2 | Defined the grammar and operator precedence (assignment, equality, comparison, term, factor, unary). Added variable declarations, `hiss(...)` printing, blocks, if/else, `slither` and `coil` loops. Defined the `--parse` output format. The parser now returns a single `Program` root node. Removed the `++` operator and the `PLUS_PLUS` token; `++` now scans as two `PLUS` tokens. |
| Lab 3 | Added the evaluator and runtime values (number, string, boolean, `dead`). Defined value printing (trimmed numbers, `dead`, unquoted strings), truthiness (only `false` and `dead` are falsey), strict typing for operators with no implicit conversion, `+` on two numbers or two strings only, equality across any types (mismatched types are unequal), left-to-right evaluation, and division by zero as a runtime error. Added the runtime error format and exit code 70. `hiss(...)` now prints in `--eval` and the REPL. Added `--eval`. |
| Lab 4 | Plain `./run <file>` now executes programs: statements run in order and only `hiss` prints. Added environments, `fang` declarations, assignment, and blocks with lexical scope (innermost binding wins). Declaring creates a binding in the current scope; assigning only updates an existing one. Redeclaring in the same scope is allowed. Undefined names are runtime errors (exit 70). Added syntax errors for missing `;` and `}` and for invalid assignment targets, with synchronization and exit 65 before anything runs. The REPL keeps variables across lines and recovers after runtime errors. `tests/lab0/hello.viper` is now a real program. |