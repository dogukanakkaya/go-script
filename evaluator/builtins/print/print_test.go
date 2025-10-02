package print

import (
	"bytes"
	"io"
	"os"
	"testing"
)

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