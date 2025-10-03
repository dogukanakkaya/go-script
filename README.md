# Go-Script: A JavaScript Interpreter

A JavaScript interpreter written in Go from scratch, featuring variables,
functions, objects, loops, conditionals, and built-in functions
like `print`, `fetch`, and `JSON` utilities.

## Overview

**Key Components:**

- **Lexer** - Converts source code into tokens
- **Parser** - Builds an Abstract Syntax Tree (AST) from tokens
- **AST** - Tree representation of code structure
- **Evaluator** - Executes the AST and produces results
- **Built-ins** - Native functions (print, fetch, setTimeout, JSON)

## How the Engine Works

### The Big Picture

When you run JavaScript code, it goes through a multi-stage pipeline:

```
Source Code → Lexer → Tokens → Parser → AST → Evaluator → Result
```

### Simple Example

**Input Code:**

```javascript
var x = 5 + 3;
print(x);
```

---

## Execution Pipeline

### Stage 1: Lexical Analysis (Lexer)

**Location:** `lexer/lexer.go`

Converts raw text into tokens:

```
"var x = 5 + 3;" → [VAR, IDENT(x), ASSIGN, NUMBER(5), PLUS, NUMBER(3), SEMICOLON]
```

**What it does:**

- Scans character by character
- Groups characters into meaningful tokens (keywords, identifiers, operators,
  literals)
- Skips whitespace and comments
- Recognizes patterns (numbers, strings, identifiers)

---

### Stage 2: Syntax Analysis (Parser)

**Location:** `parser/parser.go`

Builds an Abstract Syntax Tree (AST) from tokens:

```
[Tokens] → AST Tree Structure
```

**What it does:**

- Reads tokens sequentially
- Recognizes grammar patterns (statements, expressions)
- Handles operator precedence (`*` before `+`)
- Creates tree nodes representing code structure
- Reports syntax errors

**Example AST:**

```
VarStatement
├── Name: "x"
└── Value: InfixExpression
    ├── Left: NumberLiteral(5)
    ├── Operator: "+"
    └── Right: NumberLiteral(3)
```

---

### Stage 3: Evaluation (Evaluator)

**Location:** `evaluator/evaluator.go`

Executes the AST and produces results:

**What it does:**

- Traverses AST tree recursively
- Maintains runtime environment (variables, functions)
- Handles scoping and closures
- Executes statements and evaluates expressions
- Manages control flow (if/while/return)

**For our example:**

1. Evaluate `5 + 3` → `8.0`
2. Store in environment: `x = 8.0`
3. Call built-in `print(x)`
4. Output: `8`

---

## Language Features

### Variables & Functions

```javascript
var x = 10;
var add = function (a, b) {
    return a + b;
};
print(add(5, 3)); // 8
```

### Closures

Functions capture their defining environment:

```javascript
var makeCounter = function () {
    var count = 0;
    return function () {
        count = count + 1;
        return count;
    };
};
var counter = makeCounter();
print(counter()); // 1
print(counter()); // 2
```

### Objects

```javascript
var person = { name: "Alice", age: 30 };
print(person.name); // Alice
```

### Control Flow

```javascript
if (x > 5) { print("x is greater than 5"); }
while (i < 10) { i = i + 1; }
```

---

## Built-in Functions

### Core Built-ins

**Location:** `evaluator/builtins/`

The interpreter includes several built-in functions, organized into separate
packages:

#### `print()` - Console Output

**Package:** `evaluator/builtins/print/`

```javascript
print("Hello", 42, true); // Hello 42 true
```

#### `fetch()` - HTTP Requests

**Package:** `evaluator/builtins/fetch/`

```javascript
let response = fetch("https://jsonplaceholder.typicode.com/posts")
print("Body:", response.body)
```

#### `JSON.stringify()` / `JSON.parse()`

**Package:** `evaluator/builtins/json/`

```javascript
var obj = { name: "Dogukan", age: 25 };
var json = JSON.stringify(obj);
print(json); // {"age":25,"name":"Dogukan"}

var parsed = JSON.parse(json);
print(parsed.name); // Dogukan
```

---

## Usage

### REPL Mode (Interactive)
Note: REPL mode does not have context memory between commands. If you define a variable, it won't persist in the next command.

```bash
go run main.go
```

### File Execution Mode

```bash
go run main.go script.js
```

### Testing

Run all tests:

```bash
go test ./...
```

The project includes comprehensive test coverage across all components.

---

## Project Structure

```
go-script/
├── main.go                    # Entry point (REPL and file execution)
├── internal/
│   ├── common.go              # Shared utilities
│   └── common_test.go
├── token/
│   ├── token.go               # Token type definitions
│   └── token_test.go
├── lexer/
│   ├── lexer.go               # Lexical analyzer
│   └── lexer_test.go
├── parser/
│   ├── parser.go              # Syntax analyzer (AST builder)
│   └── parser_test.go
├── ast/
│   ├── ast.go                 # AST node definitions
│   └── ast_test.go
└── evaluator/
    ├── evaluator.go           # Runtime evaluator
    ├── evaluator_test.go
    └── builtins/              # Built-in functions
        ├── builtins.go        # Registry and exports
        ├── builtins_test.go
        ├── print/
        │   ├── print.go       # print() builtin
        │   └── print_test.go
        ├── fetch/
        │   ├── fetch.go       # fetch() builtin
        │   └── fetch_test.go
        ├── json/
        │   ├── json.go        # JSON.stringify/parse
        │   └── json_test.go
```