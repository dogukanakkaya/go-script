package builtins

import (
	"bytes"
	"io"
	"os"
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"print", true},
		{"nonexistent", false},
		{"", false},
	}

	for _, tt := range tests {
		builtin, ok := Get(tt.name)
		if ok != tt.expected {
			t.Errorf("Get(%q) returned ok=%v, expected %v", tt.name, ok, tt.expected)
		}
		if ok && builtin == nil {
			t.Errorf("Get(%q) returned nil builtin", tt.name)
		}
		if ok && builtin.Name != tt.name {
			t.Errorf("Get(%q) returned builtin with name %q", tt.name, builtin.Name)
		}
	}
}

func TestPrintFunction(t *testing.T) {
	// Capture stdout
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Print.Fn("Hello", 42, true)

	// Restore stdout
	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	expected := "Hello 42 true\n"
	if output != expected {
		t.Errorf("print() output = %q, expected %q", output, expected)
	}
}

func TestPrintFunctionSingleArg(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Print.Fn("Test")

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	expected := "Test\n"
	if output != expected {
		t.Errorf("print() output = %q, expected %q", output, expected)
	}
}

func TestPrintFunctionNoArgs(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Print.Fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	expected := "\n"
	if output != expected {
		t.Errorf("print() output = %q, expected %q", output, expected)
	}
}

func TestToString(t *testing.T) {
	tests := []struct {
		input    interface{}
		expected string
	}{
		{42.0, "42"},
		{3.14, "3.14"},
		{"hello", "hello"},
		{true, "true"},
		{false, "false"},
		{nil, "nil"},
		{0.0, "0"},
		{-5.0, "-5"},
	}

	for _, tt := range tests {
		result := toString(tt.input)
		if result != tt.expected {
			t.Errorf("toString(%v) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestToStringObject(t *testing.T) {
	obj := map[string]interface{}{
		"name": "Alice",
		"age":  30.0,
	}

	result := toString(obj)
	// map iteration order is not constant, so we check for both possible orderings
	valid1 := "{name: Alice, age: 30}"
	valid2 := "{age: 30, name: Alice}"

	if result != valid1 && result != valid2 {
		t.Errorf("toString(object) = %q, expected %q or %q", result, valid1, valid2)
	}
}

func TestPrintWithDifferentTypes(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	Print.Fn("String:", 123.0, true, false, nil)

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	output := buf.String()

	expected := "String: 123 true false nil\n"
	if output != expected {
		t.Errorf("print() output = %q, expected %q", output, expected)
	}
}

func TestBuiltinRegistry(t *testing.T) {
	expectedBuiltins := []string{"print"}

	for _, name := range expectedBuiltins {
		builtin, ok := Get(name)
		if !ok {
			t.Errorf("Expected builtin %q to be registered", name)
		}
		if builtin == nil {
			t.Errorf("Builtin %q is nil", name)
		}
		if builtin.Fn == nil {
			t.Errorf("Builtin %q has nil function", name)
		}
	}
}

func TestPrintReturnsNil(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	result := Print.Fn("test")

	w.Close()
	os.Stdout = old
	io.Copy(io.Discard, r)

	if result != nil {
		t.Errorf("print() should return nil, got %v", result)
	}
}

func TestToStringWithFloat(t *testing.T) {
	tests := []struct {
		input    float64
		expected string
	}{
		{0.0, "0"},
		{1.0, "1"},
		{-1.0, "-1"},
		{42.0, "42"},
		{3.14, "3.14"},
		{-3.14, "-3.14"},
		{100.0, "100"},
	}

	for _, tt := range tests {
		result := toString(tt.input)
		if result != tt.expected {
			t.Errorf("toString(%v) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}
