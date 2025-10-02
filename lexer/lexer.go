package lexer

import (
	"unicode"

	"go-script/token"
)

type Lexer struct {
	input    string // The source code
	position int    // Current position in input (points to current char)
	ch       byte   // Current character under examination
}

func New(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // Initialize, read the first character
	return l
}

// readChar advances the lexer to the next character in the input.
// It updates both the position and the current character (ch).
// When we reach the end of input, ch is set to 0 (null byte) to signal EOF.
//
// Example: For input "var"
//
//	Initial state: position=0, ch='v'
//	After readChar(): position=1, ch='a'
//	After readChar(): position=2, ch='r'
//	After readChar(): position=3, ch=0 (EOF)
func (l *Lexer) readChar() {
	if l.position >= len(l.input) {
		l.ch = 0 // 0 represents EOF (end of file)
	} else {
		l.ch = l.input[l.position]
	}
	l.position++
}

// peekChar looks ahead at the next character WITHOUT advancing the position.
// This is useful for two-character operators like "==", "!=", "<=", etc.
//
// Example: For input "==" at position 0
//
//	ch = '='
//	peekChar() returns '=' (the second =)
//	position stays at 0 (we only looked, didn't move)
func (l *Lexer) peekChar() byte {
	if l.position >= len(l.input) {
		return 0
	}
	return l.input[l.position]
}

// NextToken reads the next token from the input and returns it.
// This is the main method of the lexer - it's called repeatedly to get all tokens.
//
// Process:
//  1. Skip any whitespace (spaces, tabs, newlines)
//  2. Examine the current character
//  3. Determine what kind of token it starts
//  4. Read the complete token
//  5. Return the token
//
// Example: For input "var x = 5;"
//
//	Output: [
//	  {VAR, "var"},
//	  {IDENT, "x"},
//	  {ASSIGN, "="},
//	  {NUMBER, "5"},
//	  {SEMICOLON, ";"}
//	]
func (l *Lexer) NextToken() token.Token {
	var tok token.Token

	l.skipWhitespace()

	// Examine the current character and create appropriate token
	switch l.ch {
	case 0: // end of the input
		tok = token.Token{Type: token.EOF, Literal: ""}
	case '=':
		// Could be '=' (assignment) or '==' (equality check)
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.ASSIGN, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.NEQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.BANG, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.LTE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.LT, l.ch)
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = token.Token{Type: token.GTE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(token.GT, l.ch)
		}
	case '+':
		tok = newToken(token.PLUS, l.ch)
	case '-':
		tok = newToken(token.MINUS, l.ch)
	case '*':
		tok = newToken(token.STAR, l.ch)
	case '/':
		if l.peekChar() == '/' {
			// It's a comment (//...), skip until end of line
			l.skipComment()
			return l.NextToken() // get the next real token
		}
		tok = newToken(token.SLASH, l.ch)
	case '(':
		tok = newToken(token.LPAREN, l.ch)
	case ')':
		tok = newToken(token.RPAREN, l.ch)
	case '{':
		tok = newToken(token.LBRACE, l.ch)
	case '}':
		tok = newToken(token.RBRACE, l.ch)
	case ',':
		tok = newToken(token.COMMA, l.ch)
	case ';':
		tok = newToken(token.SEMICOLON, l.ch)
	case ':':
		tok = newToken(token.COLON, l.ch)
	case '.':
		tok = newToken(token.DOT, l.ch)
	case '"':
		tok.Type = token.STRING
		tok.Literal = l.readString()
	default:
		// Not a single-character token, check for multi-character tokens
		if isLetter(l.ch) {
			// It's an identifier or keyword
			tok.Literal = l.readIdentifier()
			tok.Type = token.LookupIdent(tok.Literal)
			return tok
		} else if isDigit(l.ch) {
			tok.Type = token.NUMBER
			tok.Literal = l.readNumber()
			return tok 
		} else {
			// Unknown character - create an ILLEGAL token
			tok = newToken(token.ILLEGAL, l.ch)
		}
	}

	l.readChar() // Continue to next character
	return tok
}

func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

func (l *Lexer) skipComment() {
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}

// readIdentifier reads an identifier (variable/function name) or keyword.
// Identifiers can contain letters, digits, and underscores, but must start
// with a letter or underscore.
//
// Example inputs and outputs:
//
//	"var" → "var" (will be recognized as keyword by LookupIdent)
//	"camelCaseIdentifier" → "camelCaseIdentifier"
func (l *Lexer) readIdentifier() string {
	startPos := l.position - 1 // -1 because we've already read first char
	// Keep reading while we see letters, digits, or underscores
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	return l.input[startPos : l.position-1]
}

// readNumber reads a numeric literal.
// Supports both integers and floating-point numbers.
//
// Example inputs and outputs:
//
//	"42" → "42" (integer)
//	"3.14" → "3.14" (float)
//
// Note: doesn't handle scientific notation: 1e10, 1_000_000 etc.
func (l *Lexer) readNumber() string {
	startPos := l.position - 1 // -1 because we've already read first digit
	// Keep reading digits and decimal points
	for isDigit(l.ch) || l.ch == '.' {
		l.readChar()
	}
	return l.input[startPos : l.position-1]
}

// readString reads a string literal between double quotes.
// It reads until it finds the closing quote or reaches EOF.
//
// Example: For input "hello world"
//
//	Input with quotes: "hello world"
//	Returns: hello world (without quotes)
//
// Note: doesn't handle escape sequences: \n, \t, \"
func (l *Lexer) readString() string {
	startPos := l.position // position is already past the opening "
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[startPos : l.position-1]
}

// Example:
//
//	isLetter('a') → true
//	isLetter('Z') → true
//	isLetter('_') → true
//	isLetter('5') → false
func isLetter(ch byte) bool {
	return unicode.IsLetter(rune(ch)) || ch == '_'
}


// Example:
//
//	isDigit('5') → true
//	isDigit('0') → true
//	isDigit('a') → false
func isDigit(ch byte) bool {
	return unicode.IsDigit(rune(ch))
}

// Example:
//
//	newToken(token.PLUS, '+') → {PLUS, "+"}
//	newToken(token.LPAREN, '(') → {LPAREN, "("}
func newToken(tokenType token.Type, ch byte) token.Token {
	return token.Token{Type: tokenType, Literal: string(ch)}
}
