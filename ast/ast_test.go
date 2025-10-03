package ast

import "testing"

func TestProgramCreation(t *testing.T) {
	program := &Program{
		Statements: []Statement{},
	}

	if program.Statements == nil {
		t.Error("Program.Statements should not be nil")
	}

	if len(program.Statements) != 0 {
		t.Errorf("New program should have 0 statements, got %d", len(program.Statements))
	}
}

func TestVarStatementCreation(t *testing.T) {
	stmt := &VarStatement{
		Name:  "x",
		Value: &NumberLiteral{Value: 5},
	}

	if stmt.Name != "x" {
		t.Errorf("VarStatement.Name should be 'x', got '%s'", stmt.Name)
	}

	numLit, ok := stmt.Value.(*NumberLiteral)
	if !ok {
		t.Errorf("VarStatement.Value should be *NumberLiteral, got %T", stmt.Value)
	}

	if numLit.Value != 5 {
		t.Errorf("NumberLiteral.Value should be 5, got %f", numLit.Value)
	}
}

func TestReturnStatementCreation(t *testing.T) {
	stmt := &ReturnStatement{
		Value: &NumberLiteral{Value: 42},
	}

	numLit, ok := stmt.Value.(*NumberLiteral)
	if !ok {
		t.Errorf("ReturnStatement.Value should be *NumberLiteral, got %T", stmt.Value)
	}

	if numLit.Value != 42 {
		t.Errorf("NumberLiteral.Value should be 42, got %f", numLit.Value)
	}
}

func TestExpressionStatementCreation(t *testing.T) {
	stmt := &ExpressionStatement{
		Expression: &Identifier{Name: "myVar"},
	}

	ident, ok := stmt.Expression.(*Identifier)
	if !ok {
		t.Errorf("ExpressionStatement.Expression should be *Identifier, got %T", stmt.Expression)
	}

	if ident.Name != "myVar" {
		t.Errorf("Identifier.Name should be 'myVar', got '%s'", ident.Name)
	}
}

func TestBlockStatementCreation(t *testing.T) {
	block := &BlockStatement{
		Statements: []Statement{
			&VarStatement{Name: "x", Value: &NumberLiteral{Value: 5}},
			&VarStatement{Name: "y", Value: &NumberLiteral{Value: 10}},
		},
	}

	if len(block.Statements) != 2 {
		t.Errorf("BlockStatement should have 2 statements, got %d", len(block.Statements))
	}
}

func TestIfStatementCreation(t *testing.T) {
	stmt := &IfStatement{
		Condition: &BooleanLiteral{Value: true},
		Consequence: &BlockStatement{
			Statements: []Statement{},
		},
		Alternative: nil,
	}

	boolLit, ok := stmt.Condition.(*BooleanLiteral)
	if !ok {
		t.Errorf("IfStatement.Condition should be *BooleanLiteral, got %T", stmt.Condition)
	}

	if !boolLit.Value {
		t.Error("BooleanLiteral.Value should be true")
	}

	if stmt.Consequence == nil {
		t.Error("IfStatement.Consequence should not be nil")
	}
}

func TestWhileStatementCreation(t *testing.T) {
	stmt := &WhileStatement{
		Condition: &BooleanLiteral{Value: true},
		Body: &BlockStatement{
			Statements: []Statement{},
		},
	}

	if stmt.Condition == nil {
		t.Error("WhileStatement.Condition should not be nil")
	}

	if stmt.Body == nil {
		t.Error("WhileStatement.Body should not be nil")
	}
}

func TestIdentifierCreation(t *testing.T) {
	ident := &Identifier{Name: "foobar"}

	if ident.Name != "foobar" {
		t.Errorf("Identifier.Name should be 'foobar', got '%s'", ident.Name)
	}
}

func TestNumberLiteralCreation(t *testing.T) {
	tests := []struct {
		value float64
	}{
		{42.0},
		{3.14},
		{0.0},
		{-5.0},
	}

	for _, tt := range tests {
		lit := &NumberLiteral{Value: tt.value}
		if lit.Value != tt.value {
			t.Errorf("NumberLiteral.Value should be %f, got %f", tt.value, lit.Value)
		}
	}
}

func TestStringLiteralCreation(t *testing.T) {
	tests := []string{"hello", "world", "", "123"}

	for _, str := range tests {
		lit := &StringLiteral{Value: str}
		if lit.Value != str {
			t.Errorf("StringLiteral.Value should be '%s', got '%s'", str, lit.Value)
		}
	}
}

func TestBooleanLiteralCreation(t *testing.T) {
	trueLit := &BooleanLiteral{Value: true}
	if !trueLit.Value {
		t.Error("BooleanLiteral.Value should be true")
	}

	falseLit := &BooleanLiteral{Value: false}
	if falseLit.Value {
		t.Error("BooleanLiteral.Value should be false")
	}
}

func TestPrefixExpressionCreation(t *testing.T) {
	expr := &PrefixExpression{
		Operator: "-",
		Right:    &NumberLiteral{Value: 5},
	}

	if expr.Operator != "-" {
		t.Errorf("PrefixExpression.Operator should be '-', got '%s'", expr.Operator)
	}

	numLit, ok := expr.Right.(*NumberLiteral)
	if !ok {
		t.Errorf("PrefixExpression.Right should be *NumberLiteral, got %T", expr.Right)
	}

	if numLit.Value != 5 {
		t.Errorf("NumberLiteral.Value should be 5, got %f", numLit.Value)
	}
}

func TestInfixExpressionCreation(t *testing.T) {
	expr := &InfixExpression{
		Left:     &NumberLiteral{Value: 5},
		Operator: "+",
		Right:    &NumberLiteral{Value: 3},
	}

	if expr.Operator != "+" {
		t.Errorf("InfixExpression.Operator should be '+', got '%s'", expr.Operator)
	}

	leftLit, ok := expr.Left.(*NumberLiteral)
	if !ok {
		t.Errorf("InfixExpression.Left should be *NumberLiteral, got %T", expr.Left)
	}
	if leftLit.Value != 5 {
		t.Errorf("Left NumberLiteral.Value should be 5, got %f", leftLit.Value)
	}

	rightLit, ok := expr.Right.(*NumberLiteral)
	if !ok {
		t.Errorf("InfixExpression.Right should be *NumberLiteral, got %T", expr.Right)
	}
	if rightLit.Value != 3 {
		t.Errorf("Right NumberLiteral.Value should be 3, got %f", rightLit.Value)
	}
}

func TestAssignExpressionCreation(t *testing.T) {
	expr := &AssignExpression{
		Name:  "x",
		Value: &NumberLiteral{Value: 42},
	}

	if expr.Name != "x" {
		t.Errorf("AssignExpression.Name should be 'x', got '%s'", expr.Name)
	}

	numLit, ok := expr.Value.(*NumberLiteral)
	if !ok {
		t.Errorf("AssignExpression.Value should be *NumberLiteral, got %T", expr.Value)
	}
	if numLit.Value != 42 {
		t.Errorf("NumberLiteral.Value should be 42, got %f", numLit.Value)
	}
}

func TestFunctionLiteralCreation(t *testing.T) {
	fn := &FunctionLiteral{
		Parameters: []string{"x", "y"},
		Body: &BlockStatement{
			Statements: []Statement{},
		},
	}

	if len(fn.Parameters) != 2 {
		t.Errorf("FunctionLiteral should have 2 parameters, got %d", len(fn.Parameters))
	}

	if fn.Parameters[0] != "x" {
		t.Errorf("First parameter should be 'x', got '%s'", fn.Parameters[0])
	}

	if fn.Parameters[1] != "y" {
		t.Errorf("Second parameter should be 'y', got '%s'", fn.Parameters[1])
	}

	if fn.Body == nil {
		t.Error("FunctionLiteral.Body should not be nil")
	}
}

func TestCallExpressionCreation(t *testing.T) {
	call := &CallExpression{
		Function: &Identifier{Name: "add"},
		Arguments: []Expression{
			&NumberLiteral{Value: 5},
			&NumberLiteral{Value: 3},
		},
	}

	fnIdent, ok := call.Function.(*Identifier)
	if !ok {
		t.Errorf("CallExpression.Function should be *Identifier, got %T", call.Function)
	}
	if fnIdent.Name != "add" {
		t.Errorf("Function name should be 'add', got '%s'", fnIdent.Name)
	}

	if len(call.Arguments) != 2 {
		t.Errorf("CallExpression should have 2 arguments, got %d", len(call.Arguments))
	}
}

func TestObjectLiteralCreation(t *testing.T) {
	obj := &ObjectLiteral{
		Pairs: map[string]Expression{
			"name": &StringLiteral{Value: "John"},
			"age":  &NumberLiteral{Value: 30},
		},
	}

	if len(obj.Pairs) != 2 {
		t.Errorf("ObjectLiteral should have 2 pairs, got %d", len(obj.Pairs))
	}

	nameVal, ok := obj.Pairs["name"]
	if !ok {
		t.Error("ObjectLiteral should have 'name' key")
	}

	strLit, ok := nameVal.(*StringLiteral)
	if !ok {
		t.Errorf("'name' value should be *StringLiteral, got %T", nameVal)
	}
	if strLit.Value != "John" {
		t.Errorf("'name' value should be 'John', got '%s'", strLit.Value)
	}

	ageVal, ok := obj.Pairs["age"]
	if !ok {
		t.Error("ObjectLiteral should have 'age' key")
	}

	numLit, ok := ageVal.(*NumberLiteral)
	if !ok {
		t.Errorf("'age' value should be *NumberLiteral, got %T", ageVal)
	}
	if numLit.Value != 30 {
		t.Errorf("'age' value should be 30, got %f", numLit.Value)
	}
}

func TestArrayLiteralCreation(t *testing.T) {
	arr := &ArrayLiteral{
		Elements: []Expression{
			&StringLiteral{Value: "hello"},
			&NumberLiteral{Value: 5},
			&BooleanLiteral{Value: true},
		},
	}

	if len(arr.Elements) != 3 {
		t.Errorf("ArrayLiteral should have 3 elements, got %d", len(arr.Elements))
	}

	strLit, ok := arr.Elements[0].(*StringLiteral)
	if !ok {
		t.Errorf("First element should be *StringLiteral, got %T", arr.Elements[0])
	}
	if strLit.Value != "hello" {
		t.Errorf("First element value should be 'hello', got '%s'", strLit.Value)
	}

	numLit, ok := arr.Elements[1].(*NumberLiteral)
	if !ok {
		t.Errorf("Second element should be *NumberLiteral, got %T", arr.Elements[1])
	}
	if numLit.Value != 5 {
		t.Errorf("Second element value should be 42, got %f", numLit.Value)
	}

	boolLit, ok := arr.Elements[2].(*BooleanLiteral)
	if !ok {
		t.Errorf("Third element should be *BooleanLiteral, got %T", arr.Elements[2])
	}
	if !boolLit.Value {
		t.Error("Third element value should be true")
	}
}

func TestPropertyAccessCreation(t *testing.T) {
	propAccess := &PropertyAccess{
		Object:   &Identifier{Name: "person"},
		Property: "name",
	}

	objIdent, ok := propAccess.Object.(*Identifier)
	if !ok {
		t.Errorf("PropertyAccess.Object should be *Identifier, got %T", propAccess.Object)
	}
	if objIdent.Name != "person" {
		t.Errorf("Object name should be 'person', got '%s'", objIdent.Name)
	}

	if propAccess.Property != "name" {
		t.Errorf("Property should be 'name', got '%s'", propAccess.Property)
	}
}

func TestInterfaceImplementation(t *testing.T) {
	var _ Statement = (*VarStatement)(nil)
	var _ Statement = (*ReturnStatement)(nil)
	var _ Statement = (*ExpressionStatement)(nil)
	var _ Statement = (*BlockStatement)(nil)
	var _ Statement = (*IfStatement)(nil)
	var _ Statement = (*WhileStatement)(nil)

	var _ Expression = (*Identifier)(nil)
	var _ Expression = (*NumberLiteral)(nil)
	var _ Expression = (*StringLiteral)(nil)
	var _ Expression = (*BooleanLiteral)(nil)
	var _ Expression = (*PrefixExpression)(nil)
	var _ Expression = (*InfixExpression)(nil)
	var _ Expression = (*AssignExpression)(nil)
	var _ Expression = (*FunctionLiteral)(nil)
	var _ Expression = (*CallExpression)(nil)
	var _ Expression = (*ObjectLiteral)(nil)
	var _ Expression = (*PropertyAccess)(nil)
}

func TestComplexAST(t *testing.T) {
	// Represents: var add = function(a, b) { return a + b; };
	program := &Program{
		Statements: []Statement{
			&VarStatement{
				Name: "add",
				Value: &FunctionLiteral{
					Parameters: []string{"a", "b"},
					Body: &BlockStatement{
						Statements: []Statement{
							&ReturnStatement{
								Value: &InfixExpression{
									Left:     &Identifier{Name: "a"},
									Operator: "+",
									Right:    &Identifier{Name: "b"},
								},
							},
						},
					},
				},
			},
		},
	}

	if len(program.Statements) != 1 {
		t.Fatalf("Program should have 1 statement, got %d", len(program.Statements))
	}

	varStmt, ok := program.Statements[0].(*VarStatement)
	if !ok {
		t.Fatalf("Statement should be *VarStatement, got %T", program.Statements[0])
	}

	if varStmt.Name != "add" {
		t.Errorf("Variable name should be 'add', got '%s'", varStmt.Name)
	}

	fnLit, ok := varStmt.Value.(*FunctionLiteral)
	if !ok {
		t.Fatalf("Value should be *FunctionLiteral, got %T", varStmt.Value)
	}

	if len(fnLit.Parameters) != 2 {
		t.Errorf("Function should have 2 parameters, got %d", len(fnLit.Parameters))
	}

	if len(fnLit.Body.Statements) != 1 {
		t.Errorf("Function body should have 1 statement, got %d", len(fnLit.Body.Statements))
	}
}
