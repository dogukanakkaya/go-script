package print

import (
	"fmt"
	"go-script/internal"
)

// Builtin represents a built-in function
type Builtin struct {
	Name string
	Fn   func(args ...interface{}) interface{}
}

// Print is a built-in function that prints values to stdout
//
// Syntax: print(arg1, arg2, ...)
//
// Examples:
//
//	print("Hello, World!")           → prints: Hello, World!
//	print("x =", 42)                  → prints: x = 42
var Print = &Builtin{
	Name: "print",
	Fn: func(args ...interface{}) interface{} {
		for i, arg := range args {
			if i > 0 {
				fmt.Print(" ")
			}
			fmt.Print(internal.ToString(arg))
		}
		fmt.Println()
		return nil
	},
}
