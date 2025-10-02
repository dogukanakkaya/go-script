package builtins

import (
	"fmt"
)

type Builtin struct {
	Name string
	Fn   func(args ...interface{}) interface{}
}

var builtins = map[string]*Builtin{
	"print": Print,
}

func Get(name string) (*Builtin, bool) {
	builtin, ok := builtins[name]
	return builtin, ok
}

// Print is a built-in function that prints values to stdout
//
// Syntax: print(arg1, arg2, ...)
//
// Examples:
//
//	print("Hello, World!")           → prints: Hello, World!
//	print("x =", 42)                  → prints: x = 42
//
var Print = &Builtin{
	Name: "print",
	Fn: func(args ...interface{}) interface{} {
		for i, arg := range args {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(toString(arg))
		}
		fmt.Println()
		return nil
	},
}

func toString(val interface{}) string {
	if val == nil {
		return "nil"
	}

	switch v := val.(type) {
	case string:
		return v
	case float64:
		// Format integers without decimal point
		if v == float64(int64(v)) {
			return fmt.Sprintf("%d", int64(v))
		}
		return fmt.Sprintf("%v", v)
	case bool:
		if v {
			return "true"
		}
		return "false"
	case map[string]interface{}:
		// Format object as {key: value, ...}
		result := "{"
		first := true
		for k, val := range v {
			if !first {
				result += ", "
			}
			result += k + ": " + toString(val)
			first = false
		}
		result += "}"
		return result
	default:
		// Handle other types (functions, etc.)
		return fmt.Sprintf("%v", v)
	}
}
