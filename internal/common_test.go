package internal

import "testing"

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
		result := ToString(tt.input)
		if result != tt.expected {
			t.Errorf("ToString(%v) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestToStringObject(t *testing.T) {
	obj := map[string]interface{}{
		"name": "Alice",
		"age":  30.0,
	}

	result := ToString(obj)
	// map iteration order is not constant, so we check for both possible orderings
	valid1 := "{name: Alice, age: 30}"
	valid2 := "{age: 30, name: Alice}"

	if result != valid1 && result != valid2 {
		t.Errorf("ToString(object) = %q, expected %q or %q", result, valid1, valid2)
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
		result := ToString(tt.input)
		if result != tt.expected {
			t.Errorf("ToString(%v) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}
