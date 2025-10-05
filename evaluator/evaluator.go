package evaluator

import (
	"fmt"
	"go-script/ast"
	"go-script/evaluator/builtins"
	"go-script/internal"
)

// Function represents a runtime function value
// It captures the function's parameters, body, and closure environment
//
// Example: function(x, y) { return x + y; }
//
//	Creates: Function{
//	  Parameters: ["x", "y"],
//	  Body: BlockStatement{...},
//	  Env: <current environment>
//	}
type Function struct {
	Parameters []string
	Body       *ast.BlockStatement
	Env        *Environment
}

type Value = internal.Value
type Object = internal.Object
type Array = internal.Array
type ReturnValue = internal.ReturnValue

// Environment stores variables and handles scoping
// Each new scope (function, block) creates a new Environment with a parent
//
// Example: Global scope has variables like "print"
//
//	Function scope can access global variables through the parent chain
type Environment struct {
	store map[string]Value // Variables in this scope
	outer *Environment     // Parent scope (nil for global scope)
}

// Global environment with built-in functions such as JSON namespace
func NewGlobalEnvironment() *Environment {
	env := New(nil)

	jsonObj := make(Object)
	for name, builtin := range builtins.GetJSON() {
		jsonObj[name] = builtin
	}
	env.Set("JSON", jsonObj)

	return env
}

func New(outer *Environment) *Environment {
	return &Environment{
		store: make(map[string]Value),
		outer: outer,
	}
}

// Get retrieves a variable value from this scope or any parent scope
// Returns (value, true) if found, (nil, false) if not found
//
// Example: Looking up "x" in nested scopes
//  1. Check current scope
//  2. If not found, check parent scope
//  3. Continue up the chain until found or reach global scope
func (e *Environment) Get(name string) (Value, bool) {
	val, ok := e.store[name]
	if !ok && e.outer != nil {
		// Not in this scope, check parent scope
		return e.outer.Get(name)
	}
	return val, ok
}

func (e *Environment) Set(name string, val Value) Value {
	e.store[name] = val
	return val
}

// Update updates an existing variable by searching up the scope chain
// If the variable exists in a parent scope, it updates it there
// If it doesn't exist anywhere, it creates it in the current scope
//
// Example:
//
//	global has x = 5
//	In a nested scope: Update("x", 10) → updates x in global scope to 10
func (e *Environment) Update(name string, val Value) Value {
	// Check if variable exists in current scope
	if _, ok := e.store[name]; ok {
		e.store[name] = val
		return val
	}

	// Check if variable exists in parent scopes
	if e.outer != nil {
		if _, ok := e.outer.Get(name); ok {
			return e.outer.Update(name, val)
		}
	}

	// Variable doesn't exist anywhere, create it in current scope
	e.store[name] = val
	return val
}

// Eval is the main entry point for evaluation
// It takes an AST node and evaluates it in the given environment
//
// Example:
//
//	program := parser.ParseProgram()
//	env := New(nil)
//	result := Eval(program, env)
func Eval(node ast.Node, env *Environment) Value {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node, env)
	case *ast.VarStatement:
		return evalVarStatement(node, env)
	case *ast.ReturnStatement:
		return evalReturnStatement(node, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfStatement:
		return evalIfStatement(node, env)
	case *ast.WhileStatement:
		return evalWhileStatement(node, env)
	case *ast.NumberLiteral:
		return node.Value
	case *ast.StringLiteral:
		return node.Value
	case *ast.BooleanLiteral:
		return node.Value
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.PrefixExpression:
		return evalPrefixExpression(node, env)
	case *ast.InfixExpression:
		return evalInfixExpression(node, env)
	case *ast.AssignExpression:
		return evalAssignExpression(node, env)
	case *ast.FunctionLiteral:
		return &Function{
			Parameters: node.Parameters,
			Body:       node.Body,
			Env:        env,
		}
	case *ast.CallExpression:
		return evalCallExpression(node, env)
	case *ast.ObjectLiteral:
		return evalObjectLiteral(node, env)
	case *ast.ArrayLiteral:
		return evalArrayLiteral(node, env)
	case *ast.PropertyAccess:
		return evalPropertyAccess(node, env)
	case *ast.IndexExpression:
		return evalIndexExpression(node, env)
	}

	return nil
}

// evalProgram evaluates all statements in the program
// Returns the value of the last statement, or handles return statements
//
// Example: For program "var x = 5; x + 3;"
//  1. Evaluate "var x = 5" (stores x in environment)
//  2. Evaluate "x + 3" (returns 8)
//  3. Return 8 as final result
func evalProgram(program *ast.Program, env *Environment) Value {
	var result Value

	for _, statement := range program.Statements {
		result = Eval(statement, env)

		if returnValue, ok := result.(*ReturnValue); ok {
			return returnValue.Value
		}
	}

	return result
}

// evalVarStatement evaluates a variable declaration
//
// Examples:
//
//	"var x = 42;" → evaluates 42 and stores x = 42.0
//	"var sum = 5 + 3;" → evaluates 5 + 3 and stores sum = 8.0
func evalVarStatement(node *ast.VarStatement, env *Environment) Value {
	var val Value = nil

	if node.Value != nil {
		val = Eval(node.Value, env)
	}

	env.Set(node.Name, val)
	return val
}

// evalReturnStatement evaluates a return statement
// Wraps the return value so it can bubble up through nested scopes
//
// Example: "return 42;" → ReturnValue{Value: 42.0}
func evalReturnStatement(node *ast.ReturnStatement, env *Environment) Value {
	val := Eval(node.Value, env)
	return &ReturnValue{Value: val}
}

// evalBlockStatement evaluates a block of statements
// Creates a new scope for the block
//
// Example:
//
//	{
//	  var x = 5;
//	  var y = 10;
//	  x + y;
//	}
//	→ Creates new environment, evaluates statements, returns 15
func evalBlockStatement(block *ast.BlockStatement, env *Environment) Value {
	var result Value

	blockEnv := New(env)

	for _, statement := range block.Statements {
		result = Eval(statement, blockEnv)

		if returnValue, ok := result.(*ReturnValue); ok {
			return returnValue
		}
	}

	return result
}

// evalIfStatement evaluates a conditional statement
//
// Example:
//
//	if (x > 5) {
//	  print("big");
//	} else {
//	  print("small");
//	}
//	→ Evaluates condition, executes appropriate branch
func evalIfStatement(node *ast.IfStatement, env *Environment) Value {
	condition := Eval(node.Condition, env)

	if isTruthy(condition) {
		return Eval(node.Consequence, env)
	} else if node.Alternative != nil {
		return Eval(node.Alternative, env)
	}

	return nil
}

// evalWhileStatement evaluates a while loop
//
// Example:
//
//	var x = 0;
//	while (x < 5) {
//	  print(x);
//	  x = x + 1;
//	}
//	→ Loops while condition is true
func evalWhileStatement(node *ast.WhileStatement, env *Environment) Value {
	var result Value

	for {
		condition := Eval(node.Condition, env)
		if !isTruthy(condition) {
			break
		}

		result = Eval(node.Body, env)

		if _, ok := result.(*ReturnValue); ok {
			return result
		}
	}

	return result
}

// evalIdentifier looks up a variable's value in the environment
//
// Example: "x" → looks up x in environment, returns its value
func evalIdentifier(node *ast.Identifier, env *Environment) Value {
	val, ok := env.Get(node.Name)
	if !ok {
		return nil // Variable not found
	}
	return val
}

// evalPrefixExpression evaluates prefix operators (-, !)
//
// Examples:
//
//	"-5" → -5.0
//	"!true" → false
//	"-x" → negation of x's value
func evalPrefixExpression(node *ast.PrefixExpression, env *Environment) Value {
	right := Eval(node.Right, env)

	switch node.Operator {
	case "!":
		return !isTruthy(right)
	case "-":
		// Negation
		if num, ok := right.(float64); ok {
			return -num
		}
		return 0.0
	}

	return nil
}

// evalInfixExpression evaluates binary operators (+, -, *, /, ==, !=, <, >, etc.)
//
// Examples:
//
//	"5 + 3" → 8.0
//	"10 - 2" → 8.0
//	"4 * 3" → 12.0
//	"x == 5" → true or false
//	"hello" + " world" → "hello world"
func evalInfixExpression(node *ast.InfixExpression, env *Environment) Value {
	left := Eval(node.Left, env)
	right := Eval(node.Right, env)

	switch node.Operator {
	case "+":
		if leftStr, ok := left.(string); ok {
			return leftStr + internal.ToString(right)
		}
		if rightStr, ok := right.(string); ok {
			return internal.ToString(left) + rightStr
		}

		return toFloat(left) + toFloat(right)
	case "-":
		return toFloat(left) - toFloat(right)
	case "*":
		return toFloat(left) * toFloat(right)
	case "/":
		rightNum := toFloat(right)
		if rightNum == 0 {
			return 0.0 // division by zero
		}
		return toFloat(left) / rightNum
	case "==":
		return equals(left, right)
	case "!=":
		return !equals(left, right)
	case "<":
		return toFloat(left) < toFloat(right)
	case ">":
		return toFloat(left) > toFloat(right)
	case "<=":
		return toFloat(left) <= toFloat(right)
	case ">=":
		return toFloat(left) >= toFloat(right)
	}

	return nil
}

// evalAssignExpression evaluates an assignment
//
// Example: "x = 42" → updates x to 42, returns 42
// Note: Uses Update() to modify variables in parent scopes if they exist
func evalAssignExpression(node *ast.AssignExpression, env *Environment) Value {
	val := Eval(node.Value, env)
	env.Update(node.Name, val)
	return val
}

// evalCallExpression evaluates a function call
//
// Examples:
//
//	"add(5, 3)" → calls add function with arguments [5, 3]
//	"print("hello")" → calls builtin print function
//	"JSON.stringify(obj)" → calls JSON.stringify builtin
func evalCallExpression(node *ast.CallExpression, env *Environment) Value {
	// Check for builtin functions
	if ident, ok := node.Function.(*ast.Identifier); ok {
		if builtin, ok := builtins.Get(ident.Name); ok {
			args := []interface{}{}
			for _, arg := range node.Arguments {
				args = append(args, Eval(arg, env))
			}
			return builtin.Fn(args...)
		}
	}

	function := Eval(node.Function, env)

	// Check if it's a builtin from property access (like JSON.stringify)
	if builtin, ok := function.(*internal.Builtin); ok {
		args := []interface{}{}
		for _, arg := range node.Arguments {
			args = append(args, Eval(arg, env))
		}
		return builtin.Fn(args...)
	}

	fn, ok := function.(*Function)
	if !ok {
		return nil
	}

	args := []Value{}
	for _, arg := range node.Arguments {
		args = append(args, Eval(arg, env))
	}

	// Create new environment for function execution
	// Parent is the function's closure environment (where it was defined)
	fnEnv := New(fn.Env)

	// Bind parameters to argument values
	for i, param := range fn.Parameters {
		if i < len(args) {
			fnEnv.Set(param, args[i])
		} else {
			fnEnv.Set(param, nil) // Unspecified parameters are nil
		}
	}

	result := Eval(fn.Body, fnEnv)

	if returnValue, ok := result.(*ReturnValue); ok {
		return returnValue.Value
	}

	return result
}

// evalObjectLiteral evaluates an object literal
//
// Example:
//
//	{ name: "John", age: 30 }
//	→ Object{"name": "John", "age": 30.0}
func evalObjectLiteral(node *ast.ObjectLiteral, env *Environment) Value {
	obj := make(Object)

	for key, valueNode := range node.Pairs {
		value := Eval(valueNode, env)
		obj[key] = value
	}

	return obj
}

// evalArrayLiteral evaluates an array literal
//
// Example:
//
//	[1, 2, 3]
//	→ []Value{1.0, 2.0, 3.0}
func evalArrayLiteral(node *ast.ArrayLiteral, env *Environment) Value {
	elements := make(Array, 0)

	for _, elemNode := range node.Elements {
		elem := Eval(elemNode, env)
		elements = append(elements, elem)
	}

	return elements
}

// evalIndexExpression evaluates array or object index access
//
// Example:
//
//	arr[0] → gets first element of array
//	obj["key"] → gets "key" property of object
func evalIndexExpression(node *ast.IndexExpression, env *Environment) Value {
	left := Eval(node.Left, env)
	if left == nil {
		return nil
	}

	index := Eval(node.Index, env)
	if index == nil {
		return nil
	}

	// Handle Array type
	if arr, ok := left.(Array); ok {
		idx, ok := index.(float64)
		if !ok {
			return nil
		}
		i := int(idx)
		if i < 0 || i >= len(arr) {
			return nil // undefined behavior for out of bounds
		}
		return arr[i]
	}

	// Handle Object type (string indexing)
	if obj, ok := left.(Object); ok {
		key, ok := index.(string)
		if !ok {
			return nil
		}
		return obj[key]
	}

	return nil
}

// evalPropertyAccess evaluates object property access
//
// Example:
//
//	person.name → looks up "name" property in person object
//	obj.x → looks up "x" property in obj
func evalPropertyAccess(node *ast.PropertyAccess, env *Environment) Value {
	object := Eval(node.Object, env)

	// Handle Array type - support array properties
	if arr, ok := object.(Array); ok {
		return GetArrayProperty(arr, node.Property)
	}

	// Handle Object type (map[string]Value)
	if obj, ok := object.(Object); ok {
		return obj[node.Property]
	}

	return nil
}

func isTruthy(val Value) bool {
	if val == nil {
		return false
	}

	switch v := val.(type) {
	case bool:
		return v
	case float64:
		return v != 0
	case string:
		return v != ""
	default:
		return true // Objects or functions are truthy by default
	}
}

func toFloat(val Value) float64 {
	if val == nil {
		return 0
	}

	switch v := val.(type) {
	case float64:
		return v
	case bool:
		if v {
			return 1
		}
		return 0
	case string:
		// Try to parse string as number
		var num float64
		fmt.Sscanf(v, "%f", &num)
		return num
	default:
		return 0
	}
}

func equals(a, b Value) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}

	switch aVal := a.(type) {
	case float64:
		if bVal, ok := b.(float64); ok {
			return aVal == bVal
		}
	case string:
		if bVal, ok := b.(string); ok {
			return aVal == bVal
		}
	case bool:
		if bVal, ok := b.(bool); ok {
			return aVal == bVal
		}
	}

	return false
}
