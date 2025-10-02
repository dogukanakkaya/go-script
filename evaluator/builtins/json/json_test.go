package json

import (
	"testing"
)

func TestJSONStringify(t *testing.T) {
	tests := []struct {
		name     string
		input    interface{}
		expected string
	}{
		{
			name:     "simple object",
			input:    map[string]interface{}{"name": "Alice", "age": float64(30)},
			expected: `{"age":30,"name":"Alice"}`,
		},
		{
			name:     "string",
			input:    "hello",
			expected: `"hello"`,
		},
		{
			name:     "number",
			input:    float64(42),
			expected: `42`,
		},
		{
			name:     "boolean true",
			input:    true,
			expected: `true`,
		},
		{
			name:     "boolean false",
			input:    false,
			expected: `false`,
		},
		{
			name:     "null",
			input:    nil,
			expected: `null`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Stringify.Fn(tt.input)
			
			// Check if result is an error
			if errMap, ok := result.(map[string]interface{}); ok {
				if _, hasError := errMap["error"]; hasError {
					t.Fatalf("Stringify returned error: %v", errMap["error"])
				}
			}
			
			resultStr, ok := result.(string)
			if !ok {
				t.Fatalf("Expected string result, got %T", result)
			}
			
			// For objects, check both possible orderings (map iteration is non-deterministic)
			if tt.name == "simple object" {
				alt := `{"name":"Alice","age":30}`
				if resultStr != tt.expected && resultStr != alt {
					t.Errorf("Expected %q or %q, got %q", tt.expected, alt, resultStr)
				}
			} else {
				if resultStr != tt.expected {
					t.Errorf("Expected %q, got %q", tt.expected, resultStr)
				}
			}
		})
	}
}

func TestJSONStringifyNoArgs(t *testing.T) {
	result := Stringify.Fn()
	
	errMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected error map, got %T", result)
	}
	
	errorMsg, hasError := errMap["error"]
	if !hasError {
		t.Errorf("Expected error field")
	}
	
	if errorMsg != "JSON.stringify requires exactly 1 argument" {
		t.Errorf("Expected specific error message, got %v", errorMsg)
	}
}

func TestJSONStringifyTooManyArgs(t *testing.T) {
	result := Stringify.Fn("arg1", "arg2")
	
	errMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected error map, got %T", result)
	}
	
	errorMsg, hasError := errMap["error"]
	if !hasError {
		t.Errorf("Expected error field")
	}
	
	if errorMsg != "JSON.stringify requires exactly 1 argument" {
		t.Errorf("Expected specific error message, got %v", errorMsg)
	}
}

func TestJSONParse(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		checkResult func(t *testing.T, result interface{})
	}{
		{
			name:  "simple object",
			input: `{"name":"Alice","age":30}`,
			checkResult: func(t *testing.T, result interface{}) {
				obj, ok := result.(map[string]interface{})
				if !ok {
					t.Fatalf("Expected map[string]interface{}, got %T", result)
				}
				if obj["name"] != "Alice" {
					t.Errorf("Expected name=Alice, got %v", obj["name"])
				}
				if obj["age"] != float64(30) {
					t.Errorf("Expected age=30, got %v", obj["age"])
				}
			},
		},
		{
			name:  "string",
			input: `"hello"`,
			checkResult: func(t *testing.T, result interface{}) {
				str, ok := result.(string)
				if !ok {
					t.Fatalf("Expected string, got %T", result)
				}
				if str != "hello" {
					t.Errorf("Expected 'hello', got %q", str)
				}
			},
		},
		{
			name:  "number",
			input: `42`,
			checkResult: func(t *testing.T, result interface{}) {
				num, ok := result.(float64)
				if !ok {
					t.Fatalf("Expected float64, got %T", result)
				}
				if num != 42 {
					t.Errorf("Expected 42, got %v", num)
				}
			},
		},
		{
			name:  "boolean true",
			input: `true`,
			checkResult: func(t *testing.T, result interface{}) {
				b, ok := result.(bool)
				if !ok {
					t.Fatalf("Expected bool, got %T", result)
				}
				if !b {
					t.Errorf("Expected true, got %v", b)
				}
			},
		},
		{
			name:  "boolean false",
			input: `false`,
			checkResult: func(t *testing.T, result interface{}) {
				b, ok := result.(bool)
				if !ok {
					t.Fatalf("Expected bool, got %T", result)
				}
				if b {
					t.Errorf("Expected false, got %v", b)
				}
			},
		},
		{
			name:  "null",
			input: `null`,
			checkResult: func(t *testing.T, result interface{}) {
				if result != nil {
					t.Errorf("Expected nil, got %v", result)
				}
			},
		},
		{
			name:  "array",
			input: `[1,2,3]`,
			checkResult: func(t *testing.T, result interface{}) {
				arr, ok := result.([]interface{})
				if !ok {
					t.Fatalf("Expected []interface{}, got %T", result)
				}
				if len(arr) != 3 {
					t.Errorf("Expected array length 3, got %d", len(arr))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Parse.Fn(tt.input)
			
			// Check if result is an error
			if errMap, ok := result.(map[string]interface{}); ok {
				if _, hasError := errMap["error"]; hasError {
					t.Fatalf("Parse returned error: %v", errMap["error"])
				}
			}
			
			tt.checkResult(t, result)
		})
	}
}

func TestJSONParseInvalidJSON(t *testing.T) {
	result := Parse.Fn(`{invalid json}`)
	
	errMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected error map, got %T", result)
	}
	
	_, hasError := errMap["error"]
	if !hasError {
		t.Errorf("Expected error field for invalid JSON")
	}
}

func TestJSONParseNoArgs(t *testing.T) {
	result := Parse.Fn()
	
	errMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected error map, got %T", result)
	}
	
	errorMsg, hasError := errMap["error"]
	if !hasError {
		t.Errorf("Expected error field")
	}
	
	if errorMsg != "JSON.parse requires exactly 1 argument" {
		t.Errorf("Expected specific error message, got %v", errorMsg)
	}
}

func TestJSONParseTooManyArgs(t *testing.T) {
	result := Parse.Fn("arg1", "arg2")
	
	errMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected error map, got %T", result)
	}
	
	errorMsg, hasError := errMap["error"]
	if !hasError {
		t.Errorf("Expected error field")
	}
	
	if errorMsg != "JSON.parse requires exactly 1 argument" {
		t.Errorf("Expected specific error message, got %v", errorMsg)
	}
}

func TestJSONParseNonString(t *testing.T) {
	result := Parse.Fn(float64(42))
	
	errMap, ok := result.(map[string]interface{})
	if !ok {
		t.Fatalf("Expected error map, got %T", result)
	}
	
	_, hasError := errMap["error"]
	if !hasError {
		t.Errorf("Expected error field for non-string argument")
	}
}

func TestJSONNamespace(t *testing.T) {
	// Check that JSON namespace has both methods
	if JSON["stringify"] == nil {
		t.Errorf("JSON.stringify should be defined")
	}
	if JSON["parse"] == nil {
		t.Errorf("JSON.parse should be defined")
	}
	
	// Check that they have the correct names
	if JSON["stringify"].Name != "stringify" {
		t.Errorf("Expected name 'stringify', got %q", JSON["stringify"].Name)
	}
	if JSON["parse"].Name != "parse" {
		t.Errorf("Expected name 'parse', got %q", JSON["parse"].Name)
	}
}
