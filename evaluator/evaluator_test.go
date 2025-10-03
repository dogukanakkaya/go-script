package evaluator

import (
	"go-script/parser"
	"testing"
)

func testEval(input string) Value {
	p := parser.New(input)
	program := p.ParseProgram()
	env := NewGlobalEnvironment()
	return Eval(program, env)
}

func TestEvalNumberLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5;", 5},
		{"10;", 10},
		{"42;", 42},
		{"3.14;", 3.14},
		{"0;", 0},
		{"-5;", -5},
		{"-10.5;", -10.5},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		if num, ok := result.(float64); !ok || num != tt.expected {
			t.Errorf("For input %q: expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestEvalStringLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"hello";`, "hello"},
		{`"world";`, "world"},
		{`"";`, ""},
		{`"Hello, World!";`, "Hello, World!"},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		if str, ok := result.(string); !ok || str != tt.expected {
			t.Errorf("For input %q: expected %q, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestEvalBooleanLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		if b, ok := result.(bool); !ok || b != tt.expected {
			t.Errorf("For input %q: expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestEvalPrefixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected Value
	}{
		{"!true;", false},
		{"!false;", true},
		{"!!true;", true},
		{"!!false;", false},
		{"-5;", -5.0},
		{"-10;", -10.0},
		{"--5;", 5.0},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		switch exp := tt.expected.(type) {
		case bool:
			if b, ok := result.(bool); !ok || b != exp {
				t.Errorf("For input %q: expected %v, got %v", tt.input, exp, result)
			}
		case float64:
			if num, ok := result.(float64); !ok || num != exp {
				t.Errorf("For input %q: expected %v, got %v", tt.input, exp, result)
			}
		}
	}
}

func TestEvalInfixExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"5 + 5;", 10},
		{"5 - 3;", 2},
		{"4 * 3;", 12},
		{"10 / 2;", 5},
		{"2 + 3 * 4;", 14},
		{"(2 + 3) * 4;", 20},
		{"10 - 2 - 3;", 5},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		if num, ok := result.(float64); !ok || num != tt.expected {
			t.Errorf("For input %q: expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestEvalComparisonOperators(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"5 == 5;", true},
		{"5 != 5;", false},
		{"5 == 3;", false},
		{"5 != 3;", true},
		{"5 > 3;", true},
		{"5 < 3;", false},
		{"5 >= 5;", true},
		{"5 <= 5;", true},
		{"3 < 5;", true},
		{"3 > 5;", false},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		if b, ok := result.(bool); !ok || b != tt.expected {
			t.Errorf("For input %q: expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestEvalStringConcatenation(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{`"Hello" + " " + "World";`, "Hello World"},
		{`"Number: " + 42;`, "Number: 42"},
		{`42 + " is the answer";`, "42 is the answer"},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		if str, ok := result.(string); !ok || str != tt.expected {
			t.Errorf("For input %q: expected %q, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestEvalVarStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"var x = 5; x;", 5},
		{"var x = 5; var y = 10; x + y;", 15},
		{"var x = 5; var y = x + 5; y;", 10},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		if num, ok := result.(float64); !ok || num != tt.expected {
			t.Errorf("For input %q: expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestEvalAssignments(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"var x = 5; x = 10; x;", 10},
		{"var x = 5; x = x + 5; x;", 10},
		{"var x = 1; var y = 2; x = y; x;", 2},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		if num, ok := result.(float64); !ok || num != tt.expected {
			t.Errorf("For input %q: expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestEvalIfExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected Value
	}{
		{"if (true) { 10; }", 10.0},
		{"if (false) { 10; }", nil},
		{"if (1 < 2) { 10; }", 10.0},
		{"if (1 > 2) { 10; }", nil},
		{"if (1 > 2) { 10; } else { 20; }", 20.0},
		{"if (1 < 2) { 10; } else { 20; }", 10.0},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		switch exp := tt.expected.(type) {
		case float64:
			if num, ok := result.(float64); !ok || num != exp {
				t.Errorf("For input %q: expected %v, got %v", tt.input, exp, result)
			}
		case nil:
			if result != nil {
				t.Errorf("For input %q: expected nil, got %v", tt.input, result)
			}
		}
	}
}

func TestEvalWhileLoops(t *testing.T) {
	input := `
		var x = 0;
		var sum = 0;
		while (x < 5) {
			sum = sum + x;
			x = x + 1;
		}
		sum;
	`
	result := testEval(input)
	expected := 10.0

	if num, ok := result.(float64); !ok || num != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestEvalReturnStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"return 10;", 10},
		{"return 10; 9;", 10},
		{"return 2 * 5; 9;", 10},
		{"9; return 2 * 5; 9;", 10},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		if num, ok := result.(float64); !ok || num != tt.expected {
			t.Errorf("For input %q: expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestEvalFunctionLiterals(t *testing.T) {
	input := "function(x) { x + 2; };"
	result := testEval(input)

	fn, ok := result.(*Function)
	if !ok {
		t.Fatalf("Expected *Function, got %T", result)
	}

	if len(fn.Parameters) != 1 {
		t.Errorf("Expected 1 parameter, got %d", len(fn.Parameters))
	}

	if fn.Parameters[0] != "x" {
		t.Errorf("Expected parameter 'x', got %q", fn.Parameters[0])
	}
}

func TestEvalFunctionCalls(t *testing.T) {
	tests := []struct {
		input    string
		expected float64
	}{
		{"var add = function(x, y) { x + y; }; add(5, 3);", 8},
		{"var double = function(x) { x * 2; }; double(5);", 10},
		{"var identity = function(x) { x; }; identity(42);", 42},
		{"function(x) { x + 1; }(5);", 6},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		if num, ok := result.(float64); !ok || num != tt.expected {
			t.Errorf("For input %q: expected %v, got %v", tt.input, tt.expected, result)
		}
	}
}

func TestEvalClosures(t *testing.T) {
	input := `
		var makeAdder = function(x) {
			function(y) { x + y; };
		};
		var addFive = makeAdder(5);
		addFive(3);
	`
	result := testEval(input)
	expected := 8.0

	if num, ok := result.(float64); !ok || num != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestEvalRecursion(t *testing.T) {
	input := `
		var factorial = function(n) {
			if (n == 0) {
				return 1;
			}
			return n * factorial(n - 1);
		};
		factorial(5);
	`
	result := testEval(input)
	expected := 120.0

	if num, ok := result.(float64); !ok || num != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestEvalObjectLiterals(t *testing.T) {
	input := `var person = { name: "John", age: 30 }; person;`
	result := testEval(input)

	obj, ok := result.(Object)
	if !ok {
		t.Fatalf("Expected Object, got %T", result)
	}

	if len(obj) != 2 {
		t.Errorf("Expected 2 properties, got %d", len(obj))
	}

	name, exists := obj["name"]
	if !exists {
		t.Error("Property 'name' not found")
	}
	if str, ok := name.(string); !ok || str != "John" {
		t.Errorf("Expected name 'John', got %v", name)
	}

	age, exists := obj["age"]
	if !exists {
		t.Error("Property 'age' not found")
	}
	if num, ok := age.(float64); !ok || num != 30 {
		t.Errorf("Expected age 30, got %v", age)
	}
}

func TestEvalArrayLiterals(t *testing.T) {
	input := `var arr = [1, 2, 3]; arr;`
	result := testEval(input)

	arr, ok := result.(Array)
	if !ok {
		t.Fatalf("Expected Array, got %T", result)
	}

	if len(arr) != 3 {
		t.Errorf("Expected array of length 3, got %d", len(arr))
	}

	expectedValues := []float64{1, 2, 3}
	for i, expected := range expectedValues {
		if num, ok := arr[i].(float64); !ok || num != expected {
			t.Errorf("At index %d: expected %v, got %v", i, expected, arr[i])
		}
	}
}

func TestEvalPropertyAccess(t *testing.T) {
	tests := []struct {
		input    string
		expected Value
	}{
		{`var obj = { x: 10 }; obj.x;`, 10.0},
		{`var obj = { name: "Alice" }; obj.name;`, "Alice"},
		{`var obj = { a: 1, b: 2 }; obj.a + obj.b;`, 3.0},
	}

	for _, tt := range tests {
		result := testEval(tt.input)
		switch exp := tt.expected.(type) {
		case float64:
			if num, ok := result.(float64); !ok || num != exp {
				t.Errorf("For input %q: expected %v, got %v", tt.input, exp, result)
			}
		case string:
			if str, ok := result.(string); !ok || str != exp {
				t.Errorf("For input %q: expected %q, got %v", tt.input, exp, result)
			}
		}
	}
}

func TestEnvironmentGet(t *testing.T) {
	env := New(nil)
	env.Set("x", 42.0)

	val, ok := env.Get("x")
	if !ok {
		t.Error("Variable 'x' should exist")
	}
	if num, ok := val.(float64); !ok || num != 42.0 {
		t.Errorf("Expected 42.0, got %v", val)
	}

	_, ok = env.Get("y")
	if ok {
		t.Error("Variable 'y' should not exist")
	}
}

func TestEnvironmentSet(t *testing.T) {
	env := New(nil)

	env.Set("x", 10.0)
	val, ok := env.Get("x")
	if !ok || val.(float64) != 10.0 {
		t.Error("Setting variable failed")
	}

	env.Set("x", 20.0)
	val, ok = env.Get("x")
	if !ok || val.(float64) != 20.0 {
		t.Error("Updating variable failed")
	}
}

func TestEnvironmentScoping(t *testing.T) {
	outer := New(nil)
	outer.Set("x", 10.0)

	inner := New(outer)
	inner.Set("y", 20.0)

	// Inner scope can access outer variable
	val, ok := inner.Get("x")
	if !ok || val.(float64) != 10.0 {
		t.Error("Inner scope should access outer variable")
	}

	// Inner scope has its own variable
	val, ok = inner.Get("y")
	if !ok || val.(float64) != 20.0 {
		t.Error("Inner scope should have its own variable")
	}

	// Outer scope cannot access inner variable
	_, ok = outer.Get("y")
	if ok {
		t.Error("Outer scope should not access inner variable")
	}
}

func TestEnvironmentUpdate(t *testing.T) {
	outer := New(nil)
	outer.Set("x", 10.0)

	inner := New(outer)
	inner.Update("x", 20.0)

	// Update should modify outer scope
	val, ok := outer.Get("x")
	if !ok || val.(float64) != 20.0 {
		t.Error("Update should modify variable in outer scope")
	}
}

func TestIsTruthy(t *testing.T) {
	tests := []struct {
		value    Value
		expected bool
	}{
		{true, true},
		{false, false},
		{0.0, false},
		{1.0, true},
		{42.0, true},
		{"", false},
		{"hello", true},
		{nil, false},
	}

	for _, tt := range tests {
		result := isTruthy(tt.value)
		if result != tt.expected {
			t.Errorf("isTruthy(%v) = %v, expected %v", tt.value, result, tt.expected)
		}
	}
}

func TestToFloat(t *testing.T) {
	tests := []struct {
		value    Value
		expected float64
	}{
		{42.0, 42.0},
		{3.14, 3.14},
		{true, 1.0},
		{false, 0.0},
		{nil, 0.0},
		{"5", 5.0},
		{"3.14", 3.14},
	}

	for _, tt := range tests {
		result := toFloat(tt.value)
		if result != tt.expected {
			t.Errorf("toFloat(%v) = %v, expected %v", tt.value, result, tt.expected)
		}
	}
}

func TestEquals(t *testing.T) {
	tests := []struct {
		a        Value
		b        Value
		expected bool
	}{
		{5.0, 5.0, true},
		{5.0, 3.0, false},
		{"hello", "hello", true},
		{"hello", "world", false},
		{true, true, true},
		{true, false, false},
		{nil, nil, true},
		{5.0, "5", false},
		{nil, 0.0, false},
	}

	for _, tt := range tests {
		result := equals(tt.a, tt.b)
		if result != tt.expected {
			t.Errorf("equals(%v, %v) = %v, expected %v", tt.a, tt.b, result, tt.expected)
		}
	}
}

func TestBlockStatementScoping(t *testing.T) {
	input := `
		var x = 10;
		{
			var y = 20;
			x = x + y;
		}
		x;
	`
	result := testEval(input)
	expected := 30.0

	if num, ok := result.(float64); !ok || num != expected {
		t.Errorf("Expected %v, got %v", expected, result)
	}
}

func TestDivisionByZero(t *testing.T) {
	input := "10 / 0;"
	result := testEval(input)

	if num, ok := result.(float64); !ok || num != 0.0 {
		t.Errorf("Expected 0.0 for division by zero, got %v", result)
	}
}
