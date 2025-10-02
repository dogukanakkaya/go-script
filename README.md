# Go-Script: A JavaScript-like Interpreter

A JavaScript interpreter written in Go from scratch, featuring
variables, functions, objects, loops, conditionals etc.

## Overview

**Key Components:**

- **Lexer** - Converts source code into tokens
- **Parser** - Builds an Abstract Syntax Tree (AST) from tokens
- **Evaluator** - Executes the AST and produces results

## How the Engine Works

### The Big Picture

When you run a JavaScript file or type code in the REPL, here's what happens:

```
Source Code → Lexer → Tokens → Parser → AST → Evaluator → Result
```

Let's walk through this process step by step with a concrete example:

**Input Code:**

```javascript
var x = 5 + 3;
print(x);
```

---

## Execution Flow

### Step 1: Source Code Input

**What happens:** The engine receives your source code as a string.

**Example:**

```
Input: "var x = 5 + 3;\nprint(x);"
```

This is just raw text at this point - the engine doesn't understand what it
means yet.

---

### Step 2: Lexical Analysis (Lexer)

**File:** `lexer/lexer.go`

**What happens:** The lexer scans through the source code character by character
and groups them into meaningful units called **tokens**.

**Process:**

1. Start at the first character
2. Recognize patterns (keywords, numbers, operators, etc.)
3. Create a token for each meaningful unit
4. Skip whitespace and comments
5. Move to the next character

**Example:** `var x = 5 + 3;`

The lexer reads this character by character:

| Characters | Token Produced                    | Token Type     |
| ---------- | --------------------------------- | -------------- |
| `var`      | `{Type: VAR, Literal: "var"}`     | Keyword        |
| (space)    | (skipped)                         | Whitespace     |
| `x`        | `{Type: IDENT, Literal: "x"}`     | Identifier     |
| (space)    | (skipped)                         | Whitespace     |
| `=`        | `{Type: ASSIGN, Literal: "="}`    | Operator       |
| (space)    | (skipped)                         | Whitespace     |
| `5`        | `{Type: NUMBER, Literal: "5"}`    | Number Literal |
| (space)    | (skipped)                         | Whitespace     |
| `+`        | `{Type: PLUS, Literal: "+"}`      | Operator       |
| (space)    | (skipped)                         | Whitespace     |
| `3`        | `{Type: NUMBER, Literal: "3"}`    | Number Literal |
| `;`        | `{Type: SEMICOLON, Literal: ";"}` | Delimiter      |

**Output (Token Stream):**

```
VAR "var"
IDENT "x"
ASSIGN "="
NUMBER "5"
PLUS "+"
NUMBER "3"
SEMICOLON ";"
```

**Key Insight:** The lexer converts meaningless characters into meaningful
tokens. It's like breaking a sentence into words.

---

### Step 3: Syntax Analysis (Parser)

**File:** `parser/parser.go`

**What happens:** The parser takes the stream of tokens and organizes them into
a tree structure called an **Abstract Syntax Tree (AST)** that represents the
program's structure.

**Process:**

1. Read tokens one at a time (with one-token lookahead)
2. Recognize syntax patterns (statements, expressions)
3. Build tree nodes that represent code structure
4. Handle operator precedence
5. Report syntax errors if the code is malformed

**Example:** Token stream from above

The parser recognizes patterns:

**Pattern Recognition:**

```
VAR IDENT ASSIGN <expression> SEMICOLON
```

↓ This matches a **Variable Declaration Statement**

The expression `5 + 3` is parsed as:

```
NUMBER PLUS NUMBER
```

↓ This matches an **Infix Expression** (binary operation)

**Output (Abstract Syntax Tree):**

```
Program
└── VarStatement
    ├── Name: "x"
    └── Value: InfixExpression
        ├── Left: NumberLiteral(5)
        ├── Operator: "+"
        └── Right: NumberLiteral(3)
```

**Visual Tree Structure:**

```
    VarStatement
    /           \
Name: "x"    InfixExpression
             /      |      \
       Number(5)   "+"   Number(3)
```

**Key Insight:** The parser understands the grammar of the language. It's like
understanding the grammatical structure of a sentence (subject, verb, object).

**Operator Precedence:** The parser handles precedence correctly. For example:

```javascript
2 + 3 * 4;
```

Is parsed as:

```
  +
 / \
2   *
   / \
  3   4
```

Not as:

```
    *
   / \
  +   4
 / \
2   3
```

This ensures `3 * 4` is calculated first, giving `2 + 12 = 14`, not
`5 * 4 = 20`.

---

### Step 4: Semantic Analysis & Evaluation (Evaluator)

**File:** `evaluator/evaluator.go`

**What happens:** The evaluator traverses the AST tree and executes each node,
maintaining runtime state (variables, functions) in an **Environment**.

**Process:**

1. Start at the root of the AST
2. Visit each node in order
3. For statements: execute them (create variables, control flow, etc.)
4. For expressions: evaluate them and return a value
5. Maintain environment (variable storage) with scoping
6. Handle function calls with closures

**Example:** Evaluating our AST

**Step 4.1: Evaluate VarStatement**

```
VarStatement: Name = "x", Value = InfixExpression(5 + 3)
```

First, evaluate the value expression:

**Step 4.2: Evaluate InfixExpression**

```
InfixExpression: Left = 5, Operator = "+", Right = 3
```

- Evaluate Left → Returns `5.0` (float64)
- Evaluate Right → Returns `3.0` (float64)
- Apply operator `+` → Returns `8.0`

**Step 4.3: Store in Environment**

```
Environment.Set("x", 8.0)
```

The environment now looks like:

```
{
  "x": 8.0
}
```

**Key Insight:** The evaluator brings the code to life. It's like actually
performing the actions described in the sentence.

---

### Step 5: Complete Example with Multiple Statements

Let's trace a more complex example:

**Input Code:**

```javascript
var add = function(a, b) {
    return a + b;
};
var result = add(5, 3);
print(result);
```

#### **Phase 1: Lexer Output**

```
VAR "var"
IDENT "add"
ASSIGN "="
FUNC "func"
LPAREN "("
IDENT "a"
COMMA ","
IDENT "b"
RPAREN ")"
LBRACE "{"
RETURN "return"
IDENT "a"
PLUS "+"
IDENT "b"
SEMICOLON ";"
RBRACE "}"
SEMICOLON ";"
... (continues for rest of code)
```

#### **Phase 2: Parser Output (AST)**

```
Program
├── Statement 1: VarStatement
│   ├── Name: "add"
│   └── Value: FunctionLiteral
│       ├── Parameters: ["a", "b"]
│       └── Body: BlockStatement
│           └── ReturnStatement
│               └── Value: InfixExpression
│                   ├── Left: Identifier("a")
│                   ├── Operator: "+"
│                   └── Right: Identifier("b")
│
├── Statement 2: VarStatement
│   ├── Name: "result"
│   └── Value: CallExpression
│       ├── Function: Identifier("add")
│       └── Arguments: [NumberLiteral(5), NumberLiteral(3)]
│
└── Statement 3: ExpressionStatement
    └── CallExpression
        ├── Function: Identifier("print")
        └── Arguments: [Identifier("result")]
```

#### **Phase 3: Evaluator Execution**

**Step 1: Evaluate first VarStatement**

- Create a Function object
- Store in environment: `{ "add": Function(...) }`

**Step 2: Evaluate second VarStatement**

- Need to evaluate CallExpression `add(5, 3)`
- Look up `add` in environment → Find Function
- Evaluate arguments: `5` → `5.0`, `3` → `3.0`
- Create new environment for function execution
- Bind parameters: `a = 5.0`, `b = 3.0`
- Execute function body:
  - Evaluate ReturnStatement
  - Evaluate InfixExpression: `5.0 + 3.0` → `8.0`
  - Return `8.0`
- Store in environment: `{ "add": Function(...), "result": 8.0 }`

**Step 3: Evaluate ExpressionStatement (print call)**

- Call builtin `print` function
- Look up `result` → `8.0`
- Output to console: `8`

**Final Output:**

```
8
```

---

### Step 6: Advanced Features

#### **Scoping and Closures**

**Example:**

```javascript
var makeCounter = function() {
    var count = 0;
    return function() {
        count = count + 1;
        return count;
    };
};

var counter = makeCounter();
print(counter()); // Output: 1
print(counter()); // Output: 2
print(counter()); // Output: 3
```

**How it works:**

1. **Call `makeCounter()`:**
   - Creates environment E1
   - Sets `count = 0` in E1
   - Creates inner function that captures E1 (closure)
   - Returns inner function

2. **First `counter()` call:**
   - Inner function executes with access to E1
   - Reads `count` from E1 → `0`
   - Updates `count` to `1` in E1
   - Returns `1`

3. **Second `counter()` call:**
   - Inner function still has access to E1
   - Reads `count` from E1 → `1`
   - Updates `count` to `2` in E1
   - Returns `2`

**Key Insight:** The inner function "remembers" the environment where it was
created, even after `makeCounter` has returned.

#### **Object Property Access**

**Example:**

```javascript
var person = {
    name: "Alice",
    age: 30,
};
print(person.name); // Output: Alice
```

**How it works:**

1. **Parse object literal:** Creates ObjectLiteral node with pairs
2. **Evaluate object literal:** Creates runtime Object (hash map)
   ```
   Object {
     "name": "Alice",
     "age": 30.0
   }
   ```
3. **Parse property access:** Creates PropertyAccess node
4. **Evaluate property access:** Looks up key in object hash map

#### **While Loops**

**Example:**

```javascript
var i = 0;
var sum = 0;
while (i < 5) {
    sum = sum + i;
    i = i + 1;
}
print(sum); // Output: 10
```

**Execution flow:**

1. **Initialize:** `i = 0`, `sum = 0`
2. **Iteration 1:** Condition `0 < 5` → true, execute body, `sum = 0`, `i = 1`
3. **Iteration 2:** Condition `1 < 5` → true, execute body, `sum = 1`, `i = 2`
4. **Iteration 3:** Condition `2 < 5` → true, execute body, `sum = 3`, `i = 3`
5. **Iteration 4:** Condition `3 < 5` → true, execute body, `sum = 6`, `i = 4`
6. **Iteration 5:** Condition `4 < 5` → true, execute body, `sum = 10`, `i = 5`
7. **Check:** Condition `5 < 5` → false, exit loop
8. **Result:** `sum = 10`

---

## Usage

### REPL Mode (Interactive)

```bash
go run main.go
```

Then type code interactively:

```
> var x = 10;
> var y = 20;
> x + y
30
> func greet(name) { return "Hello, " + name; }
> greet("World")
Hello, World
```

### File Execution Mode

Create a file `example.js`:

```javascript
var factorial = function(n) {
    if (n == 0) {
        return 1;
    }
    return n * factorial(n - 1);
};

print("Factorial of 5 is: " + factorial(5));
```

Run it:

```bash
go run main.go example.js
```

Output:

```
Factorial of 5 is: 120
```

---

## Testing

The project includes comprehensive test coverage (90 tests total):

### Run All Tests

```bash
go test ./...
```

## Project Structure

```
go-script/
├── main.go                 # Entry point (REPL and file execution)
├── token/
│   ├── token.go           # Token type definitions
│   └── token_test.go      # Token tests
├── lexer/
│   ├── lexer.go           # Lexical analyzer
│   └── lexer_test.go      # Lexer tests
├── parser/
│   ├── parser.go          # Syntax analyzer
│   └── parser_test.go     # Parser tests
├── ast/
│   ├── ast.go             # AST node definitions
│   └── ast_test.go        # AST tests
├── evaluator/
│   ├── evaluator.go       # Runtime evaluator
│   └── evaluator_test.go  # Evaluator tests
│   └── builtins/
│       ├── builtins.go       # Built-in functions (e.g., print)
│       └── builtins_test.go  # Built-in functions tests
└── examples/
    ├── lexer_demo.go      # Lexer demonstration
    ├── parser_demo.go     # Parser demonstration
    └── evaluator_demo.go  # Full pipeline demonstration
```

---

## Building and Running

### Build the project

```bash
go build -o go-script
```

### Run the executable

```bash
./go-script                # REPL mode
./go-script myfile.js      # Execute file
```