package token

import "testing"

func TestTokenTypeConstants(t *testing.T) {
	tests := []struct {
		tokenType Type
		expected  string
	}{
		{EOF, "EOF"},
		{ILLEGAL, "ILLEGAL"},
		{IDENT, "IDENT"},
		{NUMBER, "NUMBER"},
		{STRING, "STRING"},
		{ASSIGN, "="},
		{PLUS, "+"},
		{MINUS, "-"},
		{STAR, "*"},
		{SLASH, "/"},
		{BANG, "!"},
		{DOT, "."},
		{EQ, "=="},
		{NEQ, "!="},
		{LT, "<"},
		{GT, ">"},
		{LTE, "<="},
		{GTE, ">="},
		{LPAREN, "("},
		{RPAREN, ")"},
		{LBRACE, "{"},
		{RBRACE, "}"},
		{LBRACKET, "["},
		{RBRACKET, "]"},
		{COMMA, ","},
		{SEMICOLON, ";"},
		{COLON, ":"},
		{VAR, "var"},
		{FUNC, "function"},
		{IF, "if"},
		{ELSE, "else"},
		{WHILE, "while"},
		{RETURN, "return"},
		{TRUE, "true"},
		{FALSE, "false"},
	}

	for _, tt := range tests {
		if string(tt.tokenType) != tt.expected {
			t.Errorf("Token type mismatch: expected %q, got %q", tt.expected, string(tt.tokenType))
		}
	}
}

func TestTokenCreation(t *testing.T) {
	tok := Token{
		Type:    NUMBER,
		Literal: "42",
	}

	if tok.Type != NUMBER {
		t.Errorf("Expected token type NUMBER, got %q", tok.Type)
	}

	if tok.Literal != "42" {
		t.Errorf("Expected literal '42', got %q", tok.Literal)
	}
}

func TestLookupIdent_Keywords(t *testing.T) {
	tests := []struct {
		input    string
		expected Type
	}{
		{"var", VAR},
		{"let", LET},
		{"function", FUNC},
		{"if", IF},
		{"else", ELSE},
		{"while", WHILE},
		{"return", RETURN},
		{"true", TRUE},
		{"false", FALSE},
	}

	for _, tt := range tests {
		result := LookupIdent(tt.input)
		if result != tt.expected {
			t.Errorf("LookupIdent(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestLookupIdent_Identifiers(t *testing.T) {
	tests := []string{
		"myVariable",
		"x",
		"counter",
		"userName",
		"calculate",
		"VAR", // Case sensitive
		"FUNCTION",
		"If",
	}

	for _, ident := range tests {
		result := LookupIdent(ident)
		if result != IDENT {
			t.Errorf("LookupIdent(%q) should return IDENT, got %q", ident, result)
		}
	}
}

func TestLookupIdent_CaseSensitive(t *testing.T) {
	tests := []struct {
		input    string
		expected Type
	}{
		{"var", VAR},
		{"Var", IDENT},
		{"VAR", IDENT},
		{"if", IF},
		{"If", IDENT},
		{"IF", IDENT},
		{"true", TRUE},
		{"True", IDENT},
		{"TRUE", IDENT},
	}

	for _, tt := range tests {
		result := LookupIdent(tt.input)
		if result != tt.expected {
			t.Errorf("LookupIdent(%q) = %q, expected %q", tt.input, result, tt.expected)
		}
	}
}

func TestAllKeywordsInMap(t *testing.T) {
	expectedKeywords := []string{
		"var", "let", "function", "if", "else", "while", "return", "true", "false",
	}

	for _, keyword := range expectedKeywords {
		if _, exists := keywords[keyword]; !exists {
			t.Errorf("Keyword %q not found in keywords map", keyword)
		}
	}
}

func TestKeywordsMapSize(t *testing.T) {
	expectedSize := 9 // var, let, function, if, else, while, return, true, false

	if len(keywords) != expectedSize {
		t.Errorf("Expected %d keywords in map, got %d", expectedSize, len(keywords))
	}
}

func TestTokenTypeAsString(t *testing.T) {
	tok := Token{Type: PLUS, Literal: "+"}

	typeAsString := string(tok.Type)
	if typeAsString != "+" {
		t.Errorf("Expected '+', got %q", typeAsString)
	}
}
