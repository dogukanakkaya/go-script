package builtins

import (
	"testing"
)

func TestGet(t *testing.T) {
	tests := []struct {
		name     string
		expected bool
	}{
		{"print", true},
		{"fetch", true},
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

func TestBuiltinRegistry(t *testing.T) {
	expectedBuiltins := []string{"print", "fetch"}

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

func TestGetJSON(t *testing.T) {
	jsonNamespace := GetJSON()

	if jsonNamespace == nil {
		t.Fatal("GetJSON() returned nil")
	}

	// Check that stringify and parse exist
	if jsonNamespace["stringify"] == nil {
		t.Error("JSON.stringify should be defined")
	}
	if jsonNamespace["parse"] == nil {
		t.Error("JSON.parse should be defined")
	}

	// Check names
	if jsonNamespace["stringify"].Name != "stringify" {
		t.Errorf("Expected stringify, got %s", jsonNamespace["stringify"].Name)
	}
	if jsonNamespace["parse"].Name != "parse" {
		t.Errorf("Expected parse, got %s", jsonNamespace["parse"].Name)
	}
}
