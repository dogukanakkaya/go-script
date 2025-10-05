package evaluator

import "go-script/evaluator/builtins"

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
