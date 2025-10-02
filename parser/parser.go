package parser

import (
	"fmt"
	"go-script/ast"
	"go-script/lexer"
	"go-script/token"
	"strconv"
)

const (
	_           int = iota
	LOWEST          // Lowest precedence
	ASSIGN          // = (assignment, right-associative)
	EQUALS          // == or !=
	LESSGREATER     // < or > or <= or >=
	SUM             // + or -
	PRODUCT         // * or /
	PREFIX          // -x or !x
	CALL            // myFunction(x) or obj.property
)

var precedences = map[token.Type]int{
	token.ASSIGN: ASSIGN,
	token.EQ:     EQUALS,
	token.NEQ:    EQUALS,
	token.LT:     LESSGREATER,
	token.GT:     LESSGREATER,
	token.LTE:    LESSGREATER,
	token.GTE:    LESSGREATER,
	token.PLUS:   SUM,
	token.MINUS:  SUM,
	token.SLASH:  PRODUCT,
	token.STAR:   PRODUCT,
	token.LPAREN: CALL,
	token.DOT:    CALL,
}

type Parser struct {
	l            *lexer.Lexer // The lexer providing tokens
	currentToken token.Token  // Current token we're examining
	peekToken    token.Token  // Next token (for lookahead)
	errors       []string     // List of parsing errors
}

// New creates a new Parser for the given input source code
//
// Example usage:
//
//	p := parser.New("var x = 42;")
//	program := p.ParseProgram()
func New(input string) *Parser {
	l := lexer.New(input)
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	// Read two tokens so currentToken and peekToken are both set
	// This gives us one token of lookahead
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

// nextToken advances the parser to the next token
// currentToken becomes peekToken, and we read a new peekToken from the lexer
//
// Example: Parsing "var x"
//
//	Initial: current="var", peek="x"
//	After nextToken(): current="x", peek=EOF
func (p *Parser) nextToken() {
	p.currentToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) currentTokenIs(t token.Type) bool {
	return p.currentToken.Type == t
}

func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}

// expectPeek checks if the next token is of the expected type
// If yes, it advances to that token and returns true
// If no, it records an error and returns false
//
// Example: Parsing "var x = 5"
//
//	After seeing "var", we expectPeek(IDENT) to get "x"
//	After seeing "x", we expectPeek(ASSIGN) to get "="
func (p *Parser) expectPeek(t token.Type) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.errors = append(p.errors, fmt.Sprintf("expected next token to be %s, got %s instead", t, p.peekToken.Type))
	return false
}

func (p *Parser) getPrecedence(t token.Type) int {
	if p, ok := precedences[t]; ok {
		return p
	}
	return LOWEST
}

// ParseProgram is the entry point for parsing
// It parses the entire program and returns the root AST node
//
// Example: For input "var x = 5; var y = 10;"
//
//	Returns: Program{
//	  Statements: [
//	    VarStatement{Name: "x", Value: NumberLiteral{5}},
//	    VarStatement{Name: "y", Value: NumberLiteral{10}}
//	  ]
//	}
func (p *Parser) ParseProgram() *ast.Program {
	program := &ast.Program{}
	program.Statements = []ast.Statement{}

	// Keep parsing statements until we reach EOF
	for !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.currentToken.Type {
	case token.VAR, token.LET: // todo: for now var and let are treated the same and const is not implemented
		return p.parseVarStatement()
	case token.RETURN:
		return p.parseReturnStatement()
	case token.IF:
		return p.parseIfStatement()
	case token.WHILE:
		return p.parseWhileStatement()
	case token.LBRACE:
		return p.parseBlockStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseVarStatement parses a variable declaration
//
// Syntax: var <identifier> = <expression>;
//
//	let <identifier> = <expression>;
//
// Examples:
//
//	"var x = 42;" → VarStatement{Name: "x", Value: NumberLiteral{42}}
//	"let name = "John";" → VarStatement{Name: "name", Value: StringLiteral{"John"}}
func (p *Parser) parseVarStatement() *ast.VarStatement {
	stmt := &ast.VarStatement{}

	// Expect an identifier after 'var'
	if !p.expectPeek(token.IDENT) {
		return nil
	}

	stmt.Name = p.currentToken.Literal

	// Check if there's an initialization (= value)
	if p.peekTokenIs(token.ASSIGN) {
		p.nextToken() // consume identifier
		p.nextToken() // consume =

		stmt.Value = p.parseExpression(LOWEST)
	}

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseReturnStatement parses a return statement
//
// Syntax: return <expression>;
//
// Examples:
//
//	"return 42;" → ReturnStatement{Value: NumberLiteral{42}}
//	"return x + 5;" → ReturnStatement{Value: InfixExpression{...}}
func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	stmt := &ast.ReturnStatement{}

	p.nextToken() // move past 'return'

	stmt.Value = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseIfStatement parses an if statement with optional else
//
// Syntax: if (condition) { ... } else { ... }
//
// Example:
//
//	"if (x > 5) { print(x); } else { print("small"); }"
//	→ IfStatement{
//	    Condition: InfixExpression{...},
//	    Consequence: BlockStatement{...},
//	    Alternative: BlockStatement{...}
//	  }
func (p *Parser) parseIfStatement() *ast.IfStatement {
	stmt := &ast.IfStatement{}

	// Expect '(' after 'if'
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken() // move past '('
	stmt.Condition = p.parseExpression(LOWEST)

	// Expect ')' after condition
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// Parse the consequence block
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Consequence = p.parseBlockStatement()

	// Check for else clause
	if p.peekTokenIs(token.ELSE) {
		p.nextToken() // consume else
		p.nextToken() // move to next token

		// Else can be followed by another if (else if) or a block
		if p.currentTokenIs(token.IF) {
			stmt.Alternative = p.parseIfStatement()
		} else if p.currentTokenIs(token.LBRACE) {
			stmt.Alternative = p.parseBlockStatement()
		}
	}

	return stmt
}

// parseWhileStatement parses a while loop
//
// Syntax: while (condition) { ... }
//
// Example:
//
//	"while (x < 10) { x = x + 1; }"
//	→ WhileStatement{
//	    Condition: InfixExpression{...},
//	    Body: BlockStatement{...}
//	  }
func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	stmt := &ast.WhileStatement{}

	// Expect '(' after 'while'
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	p.nextToken() // move past '('
	stmt.Condition = p.parseExpression(LOWEST)

	// Expect ')' after condition
	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	// Parse the body block
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	stmt.Body = p.parseBlockStatement()

	return stmt
}

// parseBlockStatement parses a block of statements { ... }
//
// Example:
//
//	"{ var x = 5; print(x); }"
//	→ BlockStatement{
//	    Statements: [
//	      VarStatement{...},
//	      ExpressionStatement{...}
//	    ]
//	  }
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{}
	block.Statements = []ast.Statement{}

	p.nextToken() // move past '{'

	// Parse statements until we hit '}' or EOF
	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// parseExpressionStatement parses an expression as a statement
//
// Examples:
//
//	"x + 5;" → ExpressionStatement{Expression: InfixExpression{...}}
//	"print("hello");" → ExpressionStatement{Expression: CallExpression{...}}
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseExpression is the core of the Pratt parser
// It handles precedence automatically
//
// Example: "2 + 3 * 4"
//  1. Parse 2 (NumberLiteral)
//  2. See + (precedence 4)
//  3. Parse right side with precedence 4
//  4. See 3 (NumberLiteral)
//  5. See * (precedence 6 > 4)
//  6. Parse 3 * 4 first (precedence rules)
//  7. Return 2 + (3 * 4)
func (p *Parser) parseExpression(precedence int) ast.Expression {
	var leftExp ast.Expression

	switch p.currentToken.Type {
	case token.IDENT:
		leftExp = p.parseIdentifier()
	case token.NUMBER:
		leftExp = p.parseNumberLiteral()
	case token.STRING:
		leftExp = p.parseStringLiteral()
	case token.TRUE, token.FALSE:
		leftExp = p.parseBooleanLiteral()
	case token.BANG, token.MINUS:
		leftExp = p.parsePrefixExpression()
	case token.LPAREN:
		leftExp = p.parseGroupedExpression()
	case token.FUNC:
		leftExp = p.parseFunctionLiteral()
	case token.LBRACE:
		leftExp = p.parseObjectLiteral()
	default:
		p.noPrefixParseFnError(p.currentToken.Type)
		return nil
	}

	// Parse infix expressions (operators between expressions)
	// Continue while the next operator has higher precedence
	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.getPrecedence(p.peekToken.Type) {
		switch p.peekToken.Type {
		case token.PLUS, token.MINUS, token.STAR, token.SLASH,
			token.EQ, token.NEQ, token.LT, token.GT, token.LTE, token.GTE:
			p.nextToken()
			leftExp = p.parseInfixExpression(leftExp)
		case token.LPAREN:
			p.nextToken()
			leftExp = p.parseCallExpression(leftExp)
		case token.DOT:
			p.nextToken()
			leftExp = p.parsePropertyAccess(leftExp)
		case token.ASSIGN:
			p.nextToken()
			leftExp = p.parseAssignExpression(leftExp)
		default:
			return leftExp
		}
	}

	return leftExp
}

// parseIdentifier parses an identifier (variable/function name)
//
// Example: "myVar" → Identifier{Name: "myVar"}
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Name: p.currentToken.Literal}
}

// parseNumberLiteral parses a numeric literal
//
// Examples:
//
//	"42" → NumberLiteral{Value: 42.0}
//	"3.14" → NumberLiteral{Value: 3.14}
func (p *Parser) parseNumberLiteral() ast.Expression {
	lit := &ast.NumberLiteral{}

	value, err := strconv.ParseFloat(p.currentToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("could not parse %q as number", p.currentToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

// parseStringLiteral parses a string literal
//
// Example: "hello" → StringLiteral{Value: "hello"}
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Value: p.currentToken.Literal}
}

// parseBooleanLiteral parses a boolean literal
//
// Examples:
//
//	"true" → BooleanLiteral{Value: true}
//	"false" → BooleanLiteral{Value: false}
func (p *Parser) parseBooleanLiteral() ast.Expression {
	return &ast.BooleanLiteral{Value: p.currentTokenIs(token.TRUE)}
}

// parsePrefixExpression parses a prefix operator expression
//
// Examples:
//
//	"-5" → PrefixExpression{Operator: "-", Right: NumberLiteral{5}}
//	"!true" → PrefixExpression{Operator: "!", Right: BooleanLiteral{true}}
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Operator: p.currentToken.Literal,
	}

	p.nextToken()

	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// parseInfixExpression parses a binary operator expression
//
// Examples:
//
//	"5 + 3" → InfixExpression{Left: NumberLiteral{5}, Op: "+", Right: NumberLiteral{3}}
//	"x * 2" → InfixExpression{Left: Identifier{"x"}, Op: "*", Right: NumberLiteral{2}}
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Operator: p.currentToken.Literal,
		Left:     left,
	}

	precedence := p.getPrecedence(p.currentToken.Type)
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// parseGroupedExpression parses an expression in parentheses
//
// Example: "(2 + 3)" → InfixExpression{...}
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken() // move past '('

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// parseFunctionLiteral parses a function definition
//
// Syntax: function(param1, param2) { ... }
//
// Example:
//
//	"function(x, y) { return x + y; }"
//	→ FunctionLiteral{
//	    Parameters: ["x", "y"],
//	    Body: BlockStatement{...}
//	  }
func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{}

	// Expect '(' after 'function'
	if !p.expectPeek(token.LPAREN) {
		return nil
	}

	lit.Parameters = p.parseFunctionParameters()

	// Expect function body
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	lit.Body = p.parseBlockStatement()

	return lit
}

// parseFunctionParameters parses the parameter list of a function
//
// Example: "(a, b, c)" → ["a", "b", "c"]
func (p *Parser) parseFunctionParameters() []string {
	identifiers := []string{}

	// Empty parameter list
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken() // move to first parameter

	identifiers = append(identifiers, p.currentToken.Literal)

	// Parse remaining parameters
	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // consume comma
		p.nextToken() // move to next parameter
		identifiers = append(identifiers, p.currentToken.Literal)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// parseCallExpression parses a function call
//
// Examples:
//
//	"add(5, 3)" → CallExpression{Function: Identifier{"add"}, Arguments: [...]}
//	"print("hello")" → CallExpression{Function: Identifier{"print"}, Arguments: [...]}
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Function: function}
	exp.Arguments = p.parseCallArguments()
	return exp
}

// parseCallArguments parses the argument list of a function call
//
// Example: "(5, 3, x)" → [NumberLiteral{5}, NumberLiteral{3}, Identifier{"x"}]
func (p *Parser) parseCallArguments() []ast.Expression {
	args := []ast.Expression{}

	// Empty argument list
	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return args
	}

	p.nextToken() // move to first argument
	args = append(args, p.parseExpression(LOWEST))

	// Parse remaining arguments
	for p.peekTokenIs(token.COMMA) {
		p.nextToken() // consume comma
		p.nextToken() // move to next argument
		args = append(args, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return args
}

// parseObjectLiteral parses an object literal
//
// Example:
//
//	"{ name: "John", age: 30 }"
//	→ ObjectLiteral{
//	    Pairs: {
//	      "name": StringLiteral{"John"},
//	      "age": NumberLiteral{30}
//	    }
//	  }
func (p *Parser) parseObjectLiteral() ast.Expression {
	obj := &ast.ObjectLiteral{}
	obj.Pairs = make(map[string]ast.Expression)

	p.nextToken() // move past '{'

	for !p.currentTokenIs(token.RBRACE) && !p.currentTokenIs(token.EOF) {
		// Parse key (can be identifier or string)
		var key string
		if p.currentTokenIs(token.IDENT) || p.currentTokenIs(token.STRING) {
			key = p.currentToken.Literal
		} else {
			return nil
		}

		// Expect ':' after key
		if !p.expectPeek(token.COLON) {
			return nil
		}

		p.nextToken() // move to value
		value := p.parseExpression(LOWEST)
		obj.Pairs[key] = value

		// Check for comma (more properties) or closing brace
		if p.peekTokenIs(token.COMMA) {
			p.nextToken() // consume comma
			p.nextToken() // move to next property
		} else if p.peekTokenIs(token.RBRACE) {
			p.nextToken() // consume closing brace
			break
		}
	}

	return obj
}

// parsePropertyAccess parses object property access
//
// Example:
//
//	"person.name" → PropertyAccess{Object: Identifier{"person"}, Property: "name"}
func (p *Parser) parsePropertyAccess(object ast.Expression) ast.Expression {
	exp := &ast.PropertyAccess{Object: object}

	if !p.expectPeek(token.IDENT) {
		return nil
	}

	exp.Property = p.currentToken.Literal

	return exp
}

// parseAssignExpression parses an assignment expression
//
// Example:
//
//	"x = 5" → AssignExpression{Name: "x", Value: NumberLiteral{5}}
func (p *Parser) parseAssignExpression(left ast.Expression) ast.Expression {
	// Assignment only works with identifiers on the left
	ident, ok := left.(*ast.Identifier)
	if !ok {
		p.errors = append(p.errors, "invalid assignment target")
		return nil
	}

	exp := &ast.AssignExpression{Name: ident.Name}

	p.nextToken() // move past '='
	exp.Value = p.parseExpression(LOWEST)

	return exp
}

// noPrefixParseFnError records an error when we can't parse a prefix expression
func (p *Parser) noPrefixParseFnError(t token.Type) {
	msg := fmt.Sprintf("no prefix parse function for %s found", t)
	p.errors = append(p.errors, msg)
}
