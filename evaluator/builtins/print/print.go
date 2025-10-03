package print

import (
	"fmt"
	"go-script/internal"
)

// Print is a built-in function that prints values to stdout
//
// Syntax: print(arg1, arg2, ...)
//
// Examples:
//
//	print("Hello, World!")           → prints: Hello, World!
//	print("x =", 42)                  → prints: x = 42
var Print = &internal.Builtin{
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
