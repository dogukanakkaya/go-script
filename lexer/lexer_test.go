package lexer

import (
	"go-script/token"
	"testing"
)

func TestNextToken_BasicTokens(t *testing.T) {
	input := `=+(){},;`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.ASSIGN, "="},
		{token.PLUS, "+"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RBRACE, "}"},
		{token.COMMA, ","},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_SimpleProgram(t *testing.T) {
	input := `var x = 5;
var y = 10;
var add = function(a, b) {
	return a + b;
};
`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.VAR, "var"},
		{token.IDENT, "y"},
		{token.ASSIGN, "="},
		{token.NUMBER, "10"},
		{token.SEMICOLON, ";"},
		{token.VAR, "var"},
		{token.IDENT, "add"},
		{token.ASSIGN, "="},
		{token.FUNC, "function"},
		{token.LPAREN, "("},
		{token.IDENT, "a"},
		{token.COMMA, ","},
		{token.IDENT, "b"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.RETURN, "return"},
		{token.IDENT, "a"},
		{token.PLUS, "+"},
		{token.IDENT, "b"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Operators(t *testing.T) {
	input := `+ - * / ! == != < > <= >=`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.PLUS, "+"},
		{token.MINUS, "-"},
		{token.STAR, "*"},
		{token.SLASH, "/"},
		{token.BANG, "!"},
		{token.EQ, "=="},
		{token.NEQ, "!="},
		{token.LT, "<"},
		{token.GT, ">"},
		{token.LTE, "<="},
		{token.GTE, ">="},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Keywords(t *testing.T) {
	input := `var function if else while return true false`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.FUNC, "function"},
		{token.IF, "if"},
		{token.ELSE, "else"},
		{token.WHILE, "while"},
		{token.RETURN, "return"},
		{token.TRUE, "true"},
		{token.FALSE, "false"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Identifiers(t *testing.T) {
	input := `myVar _private camelCase snake_case var123`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IDENT, "myVar"},
		{token.IDENT, "_private"},
		{token.IDENT, "camelCase"},
		{token.IDENT, "snake_case"},
		{token.IDENT, "var123"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Numbers(t *testing.T) {
	input := `42 3.14 0.5 100.99 0`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.NUMBER, "42"},
		{token.NUMBER, "3.14"},
		{token.NUMBER, "0.5"},
		{token.NUMBER, "100.99"},
		{token.NUMBER, "0"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Strings(t *testing.T) {
	input := `"hello" "world" "hello world" "" "123"`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.STRING, "hello"},
		{token.STRING, "world"},
		{token.STRING, "hello world"},
		{token.STRING, ""},
		{token.STRING, "123"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Comments(t *testing.T) {
	input := `var x = 5; // this is a comment
var y = 10; // another comment
// full line comment
var z = 15;`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.VAR, "var"},
		{token.IDENT, "y"},
		{token.ASSIGN, "="},
		{token.NUMBER, "10"},
		{token.SEMICOLON, ";"},
		{token.VAR, "var"},
		{token.IDENT, "z"},
		{token.ASSIGN, "="},
		{token.NUMBER, "15"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_IfStatement(t *testing.T) {
	input := `if (x > 5) { print(x); } else { print("small"); }`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IF, "if"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.GT, ">"},
		{token.NUMBER, "5"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "print"},
		{token.LPAREN, "("},
		{token.IDENT, "x"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.ELSE, "else"},
		{token.LBRACE, "{"},
		{token.IDENT, "print"},
		{token.LPAREN, "("},
		{token.STRING, "small"},
		{token.RPAREN, ")"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_ObjectLiteral(t *testing.T) {
	input := `{ name: "John", age: 30 }`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.LBRACE, "{"},
		{token.IDENT, "name"},
		{token.COLON, ":"},
		{token.STRING, "John"},
		{token.COMMA, ","},
		{token.IDENT, "age"},
		{token.COLON, ":"},
		{token.NUMBER, "30"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_PropertyAccess(t *testing.T) {
	input := `person.name obj.method()`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.IDENT, "person"},
		{token.DOT, "."},
		{token.IDENT, "name"},
		{token.IDENT, "obj"},
		{token.DOT, "."},
		{token.IDENT, "method"},
		{token.LPAREN, "("},
		{token.RPAREN, ")"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_WhileLoop(t *testing.T) {
	input := `while (i < 10) { i = i + 1; }`

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.WHILE, "while"},
		{token.LPAREN, "("},
		{token.IDENT, "i"},
		{token.LT, "<"},
		{token.NUMBER, "10"},
		{token.RPAREN, ")"},
		{token.LBRACE, "{"},
		{token.IDENT, "i"},
		{token.ASSIGN, "="},
		{token.IDENT, "i"},
		{token.PLUS, "+"},
		{token.NUMBER, "1"},
		{token.SEMICOLON, ";"},
		{token.RBRACE, "}"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}

func TestNextToken_Whitespace(t *testing.T) {
	input := `   var    x   =   5   ;   `

	tests := []struct {
		expectedType    token.Type
		expectedLiteral string
	}{
		{token.VAR, "var"},
		{token.IDENT, "x"},
		{token.ASSIGN, "="},
		{token.NUMBER, "5"},
		{token.SEMICOLON, ";"},
		{token.EOF, ""},
	}

	l := New(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q",
				i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q",
				i, tt.expectedLiteral, tok.Literal)
		}
	}
}
