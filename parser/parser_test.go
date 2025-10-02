package parser

import (
	"go-script/ast"
	"testing"
)

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}

func TestVarStatements(t *testing.T) {
	input := `
var x = 5;
var y = 10;
var foobar = 838383;
`

	p := New(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	tests := []struct {
		expectedIdentifier string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, tt := range tests {
		stmt := program.Statements[i]
		if !testVarStatement(t, stmt, tt.expectedIdentifier) {
			return
		}
	}
}

func testVarStatement(t *testing.T, s ast.Statement, name string) bool {
	varStmt, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("s not *ast.VarStatement. got=%T", s)
		return false
	}

	if varStmt.Name != name {
		t.Errorf("varStmt.Name not '%s'. got=%s", name, varStmt.Name)
		return false
	}

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
return 5;
return 10;
return 993322;
`

	p := New(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d",
			len(program.Statements))
	}

	for _, stmt := range program.Statements {
		returnStmt, ok := stmt.(*ast.ReturnStatement)
		if !ok {
			t.Errorf("stmt not *ast.ReturnStatement. got=%T", stmt)
			continue
		}
		if returnStmt.Value == nil {
			t.Errorf("returnStmt.Value is nil")
		}
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	p := New(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Name != "foobar" {
		t.Errorf("ident.Name not %s. got=%s", "foobar", ident.Name)
	}
}

func TestNumberLiteralExpression(t *testing.T) {
	input := "5;"

	p := New(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	literal, ok := stmt.Expression.(*ast.NumberLiteral)
	if !ok {
		t.Fatalf("exp not *ast.NumberLiteral. got=%T", stmt.Expression)
	}
	if literal.Value != 5 {
		t.Errorf("literal.Value not %f. got=%f", 5.0, literal.Value)
	}
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	p := New(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

func TestBooleanExpression(t *testing.T) {
	tests := []struct {
		input         string
		expectedValue bool
	}{
		{"true;", true},
		{"false;", false},
	}

	for _, tt := range tests {
		p := New(tt.input)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program has not enough statements. got=%d",
				len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		boolean, ok := stmt.Expression.(*ast.BooleanLiteral)
		if !ok {
			t.Fatalf("exp not *ast.BooleanLiteral. got=%T", stmt.Expression)
		}
		if boolean.Value != tt.expectedValue {
			t.Errorf("boolean.Value not %t. got=%t", tt.expectedValue, boolean.Value)
		}
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5.0},
		{"-15;", "-", 15.0},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		p := New(tt.input)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  float64
		operator   string
		rightValue float64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, tt := range infixTests {
		p := New(tt.input)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
				1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
				program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s",
				tt.operator, exp.Operator)
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"-a * b", "((-a) * b)"},
		{"!-a", "(!(-a))"},
		{"a + b + c", "((a + b) + c)"},
		{"a + b - c", "((a + b) - c)"},
		{"a * b * c", "((a * b) * c)"},
		{"a * b / c", "((a * b) / c)"},
		{"a + b / c", "(a + (b / c))"},
		{"a + b * c + d / e - f", "(((a + (b * c)) + (d / e)) - f)"},
		{"3 + 4 * 5 == 3 * 1 + 4 * 5", "((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))"},
		{"1 + (2 + 3) + 4", "((1 + (2 + 3)) + 4)"},
		{"(5 + 5) * 2", "((5 + 5) * 2)"},
		{"2 / (5 + 5)", "(2 / (5 + 5))"},
	}

	for _, tt := range tests {
		p := New(tt.input)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.Statements[0].(*ast.ExpressionStatement).Expression
		if actual == nil {
			t.Errorf("expression is nil")
			continue
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	p := New(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Condition == nil {
		t.Errorf("stmt.Condition is nil")
	}

	if stmt.Consequence == nil {
		t.Errorf("stmt.Consequence is nil")
	}

	if len(stmt.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(stmt.Consequence.Statements))
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	p := New(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.IfStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Alternative == nil {
		t.Errorf("stmt.Alternative was nil")
	}
}

func TestWhileStatement(t *testing.T) {
	input := `while (x < 10) { x = x + 1; }`

	p := New(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.WhileStatement. got=%T",
			program.Statements[0])
	}

	if stmt.Condition == nil {
		t.Errorf("stmt.Condition is nil")
	}

	if stmt.Body == nil {
		t.Errorf("stmt.Body is nil")
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `function(x, y) { x + y; }`

	p := New(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
			stmt.Expression)
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	if function.Parameters[0] != "x" {
		t.Fatalf("parameter is not 'x'. got=%s", function.Parameters[0])
	}

	if function.Parameters[1] != "y" {
		t.Fatalf("parameter is not 'y'. got=%s", function.Parameters[1])
	}

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	p := New(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}
}

func TestObjectLiteralParsing(t *testing.T) {
	input := `var obj = { name: "John", age: 30 };`

	p := New(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.VarStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.VarStatement. got=%T",
			program.Statements[0])
	}

	objLit, ok := stmt.Value.(*ast.ObjectLiteral)
	if !ok {
		t.Fatalf("stmt.Value is not ast.ObjectLiteral. got=%T",
			stmt.Value)
	}

	if len(objLit.Pairs) != 2 {
		t.Fatalf("object literal has wrong number of pairs. got=%d",
			len(objLit.Pairs))
	}

	if _, ok := objLit.Pairs["name"]; !ok {
		t.Errorf("object literal missing 'name' key")
	}

	if _, ok := objLit.Pairs["age"]; !ok {
		t.Errorf("object literal missing 'age' key")
	}
}

func TestPropertyAccessParsing(t *testing.T) {
	input := "person.name"

	p := New(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	propAccess, ok := stmt.Expression.(*ast.PropertyAccess)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.PropertyAccess. got=%T",
			stmt.Expression)
	}

	if propAccess.Property != "name" {
		t.Errorf("propAccess.Property not 'name'. got=%s", propAccess.Property)
	}
}

func TestAssignmentParsing(t *testing.T) {
	input := "x = 5;"

	p := New(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	assign, ok := stmt.Expression.(*ast.AssignExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.AssignExpression. got=%T",
			stmt.Expression)
	}

	if assign.Name != "x" {
		t.Errorf("assign.Name not 'x'. got=%s", assign.Name)
	}
}

func TestBlockStatement(t *testing.T) {
	input := `{ var x = 5; var y = 10; }`

	p := New(input)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d",
			len(program.Statements))
	}

	block, ok := program.Statements[0].(*ast.BlockStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.BlockStatement. got=%T",
			program.Statements[0])
	}

	if len(block.Statements) != 2 {
		t.Errorf("block.Statements does not contain 2 statements. got=%d",
			len(block.Statements))
	}
}
