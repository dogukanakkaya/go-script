package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"go-script/environment"
	"go-script/evaluator"
	"go-script/parser"
)

func runCode(code string) bool {
	env := environment.NewGlobalEnvironment()

	p := parser.New(code)
	program := p.ParseProgram()

	if len(p.Errors()) > 0 {
		fmt.Println("Parser errors:")
		for _, msg := range p.Errors() {
			fmt.Printf("  %s\n", msg)
		}
		return false
	}

	evaluator.Eval(program, env)
	return true
}

func runFile(filename string) {
	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file '%s': %v\n", filename, err)
		os.Exit(1)
	}

	code := string(content)
	if !runCode(code) {
		os.Exit(1)
	}
}

func runREPL() {
	fmt.Println("╔═══════════════════════════════════════════════════════╗")
	fmt.Println("║   go-script REPL - JavaScript Interpreter             ║")
	fmt.Println("║   Type JavaScript code and press Enter to execute     ║")
	fmt.Println("║   Type 'exit' or press Ctrl+C to quit                 ║")
	fmt.Println("╚═══════════════════════════════════════════════════════╝")
	fmt.Println()

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print(">> ")

		if !scanner.Scan() {
			break
		}

		line := scanner.Text()

		if strings.TrimSpace(line) == "" {
			continue
		}

		if strings.TrimSpace(line) == "exit" {
			fmt.Println("Goodbye!")
			return
		}

		runCode(line)
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
	}
}

func printUsage() {
	fmt.Println("go-script - JavaScript Interpreter")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  go run main.go         # Start REPL")
	fmt.Println("  go run main.go test.js # Run with go run")
	fmt.Println()
}

func main() {
	args := os.Args[1:]

	if len(args) == 0 {
		runREPL()
		return
	}

	if args[0] == "-h" || args[0] == "--help" {
		printUsage()
		return
	}

	filename := args[0]
	runFile(filename)
}
