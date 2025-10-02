package token

type Type string

type Token struct {
	Type    Type
	Literal string
}

const (
	// Special tokens
	EOF     Type = "EOF"
	ILLEGAL Type = "ILLEGAL"

	// Identifiers and literals
	IDENT  Type = "IDENT"
	NUMBER Type = "NUMBER"
	STRING Type = "STRING"

	// Operators - used for mathematical and logical operations
	ASSIGN Type = "="
	PLUS   Type = "+"
	MINUS  Type = "-"
	STAR   Type = "*"
	SLASH  Type = "/"
	BANG   Type = "!"
	DOT    Type = "."

	// Comparison operators - used for comparing values
	EQ  Type = "=="
	NEQ Type = "!="
	LT  Type = "<"
	GT  Type = ">"
	LTE Type = "<="
	GTE Type = ">="

	// Delimiters - used to group and separate code elements
	LPAREN    Type = "("
	RPAREN    Type = ")"
	LBRACE    Type = "{"
	RBRACE    Type = "}"
	COMMA     Type = ","
	SEMICOLON Type = ";"
	COLON     Type = ":"

	// Keywords - reserved words with special meaning
	VAR    Type = "var"
	LET    Type = "let"
	FUNC   Type = "function"
	IF     Type = "if"
	ELSE   Type = "else"
	WHILE  Type = "while"
	RETURN Type = "return"
	TRUE   Type = "true"
	FALSE  Type = "false"
)

// Example: When the lexer sees "var", it checks this map and returns TokVar
//
//	When it sees "myVariable", it's not in the map, so it returns TokIdent
var keywords = map[string]Type{
	"var":      VAR,
	"let":      LET,
	"function": FUNC,
	"if":       IF,
	"else":     ELSE,
	"while":    WHILE,
	"return":   RETURN,
	"true":     TRUE,
	"false":    FALSE,
}

// LookupIdent checks if an identifier is a keyword.
// If it is, it returns the keyword's token type.
// If not, it returns IDENT (meaning it's a regular identifier/variable name).
//
// Example:
//
//	LookupIdent("var")   -> VAR (it's a keyword)
//	LookupIdent("myVar") -> IDENT (it's a regular identifier)
//	LookupIdent("if")    -> IF (it's a keyword)
func LookupIdent(ident string) Type {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
