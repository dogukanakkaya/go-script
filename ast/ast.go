package ast

type Node interface{}

type Statement interface {
	Node
	statementNode()
}

type Expression interface {
	Node
	expressionNode()
}

// Program is the root node of the AST.
// It contains all the top-level statements in the source code.
//
// Example: For source code:
//
//	var x = 5;
//	var y = 10;
//	print(x + y);
//
// The Program node contains 3 statements: 2 VarStatements and 1 ExpressionStatement
type Program struct {
	Statements []Statement
}

type VarStatement struct {
	Name  string
	Value Expression
}

// for compile-time type safety
func (vs *VarStatement) statementNode() {}

type ReturnStatement struct {
	Value Expression
}

func (rs *ReturnStatement) statementNode() {}

type ExpressionStatement struct {
	Expression Expression
}

func (es *ExpressionStatement) statementNode() {}

type BlockStatement struct {
	Statements []Statement
}

func (bs *BlockStatement) statementNode() {}

type IfStatement struct {
	Condition   Expression
	Consequence *BlockStatement // The block to execute if condition is true
	Alternative Statement       // The else block (can be nil, another IfStatement, or BlockStatement)
}

func (is *IfStatement) statementNode() {}

type WhileStatement struct {
	Condition Expression
	Body      *BlockStatement
}

func (ws *WhileStatement) statementNode() {}

type Identifier struct {
	Name string
}

func (i *Identifier) expressionNode() {}

type NumberLiteral struct {
	Value float64
}

func (nl *NumberLiteral) expressionNode() {}

type StringLiteral struct {
	Value string
}

func (sl *StringLiteral) expressionNode() {}

type BooleanLiteral struct {
	Value bool
}

func (bl *BooleanLiteral) expressionNode() {}

type PrefixExpression struct {
	Operator string     // The operator: "-" (negation) or "!" (logical NOT)
	Right    Expression // The operand expression
}

func (pe *PrefixExpression) expressionNode() {}

type InfixExpression struct {
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode() {}

type AssignExpression struct {
	Name  string
	Value Expression
}

func (ae *AssignExpression) expressionNode() {}

type FunctionLiteral struct {
	Parameters []string
	Body       *BlockStatement
}

func (fl *FunctionLiteral) expressionNode() {}

type CallExpression struct {
	Function  Expression
	Arguments []Expression
}

func (ce *CallExpression) expressionNode() {}

type ObjectLiteral struct {
	Pairs map[string]Expression
}

func (ol *ObjectLiteral) expressionNode() {}

type ArrayLiteral struct {
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode() {}

type PropertyAccess struct {
	Object   Expression
	Property string
}

func (pa *PropertyAccess) expressionNode() {}

type IndexExpression struct {
	Left  Expression // The array or object being indexed
	Index Expression // The index value
}

func (ie *IndexExpression) expressionNode() {}
