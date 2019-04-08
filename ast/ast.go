package ast

import (
	"strings"

	"github.com/dgnorton/monkey/lexer"
)

// Node represents a node in the AST. All nodes implement this interface.
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement represents a statement in the AST. All statement nodes
// implement this interface.
type Statement interface {
	Node
	statement()
}

// Expression represents an expression in the AST. All expressions nodes
// implement this interface.
type Expression interface {
	Node
	expression()
}

// Program is the top level node in the AST.
type Program struct {
	Statements []Statement
}

// NewProgram creates a new program.
func NewProgram() *Program {
	return &Program{
		Statements: []Statement{},
	}
}

// String returns a string representation of the program.
func (node *Program) String() string {
	var sb strings.Builder

	for _, stmt := range node.Statements {
		sb.WriteString(stmt.String() + "\n")
	}

	return sb.String()
}

// AddStmt adds a statement to the program.
func (node *Program) AddStmt(stmt Statement) {
	node.Statements = append(node.Statements, stmt)
}

func (node *Program) TokenLiteral() string {
	if len(node.Statements) == 0 {
		return ""
	}
	return node.Statements[0].TokenLiteral()
}

// LetStmt is a let statement node.
type LetStmt struct {
	Token *lexer.Token
	Name  *IdentExpr
	Value Expression
}

// NewLetStmt returns a new let statement.
func NewLetStmt(t *lexer.Token, name *IdentExpr, value Expression) *LetStmt {
	return &LetStmt{
		Token: t,
		Name:  name,
		Value: value,
	}
}

func (stmt *LetStmt) statement()           {}
func (stmt *LetStmt) TokenLiteral() string { return stmt.Token.String }

func (stmt *LetStmt) String() string {
	var sb strings.Builder

	sb.WriteString(stmt.Token.String + " ")
	sb.WriteString(stmt.Name.String())
	sb.WriteString(" = ")

	if stmt.Value != nil {
		sb.WriteString(stmt.Value.String())
	} else {
		sb.WriteString("<nil-expr>")
	}

	sb.WriteString(";")

	return sb.String()
}

// ReturnStmt represents a return statement.
type ReturnStmt struct {
	Token *lexer.Token
	Value Expression
}

// NewReturnStmt creates a new return statement.
func NewReturnStmt(t *lexer.Token, e Expression) *ReturnStmt {
	return &ReturnStmt{
		Token: t,
		Value: e,
	}
}

func (stmt *ReturnStmt) statement()           {}
func (stmt *ReturnStmt) TokenLiteral() string { return stmt.Token.String }

func (stmt *ReturnStmt) String() string {
	var sb strings.Builder

	sb.WriteString(stmt.Token.String + " ")

	if stmt.Value != nil {
		sb.WriteString(stmt.Value.String())
	} else {
		sb.WriteString("<nil-expr-stmt>")
	}

	sb.WriteString(";")

	return sb.String()
}

// ExprStmt represents a standalone expression used like a statement.
type ExprStmt struct {
	Token *lexer.Token
	Expr  Expression
}

// NewExprStmt creates a new expression statement.
func NewExprStmt(t *lexer.Token, e Expression) *ExprStmt {
	return &ExprStmt{
		Token: t,
		Expr:  e,
	}
}

func (stmt *ExprStmt) statement()           {}
func (stmt *ExprStmt) TokenLiteral() string { return stmt.Token.String }

func (stmt *ExprStmt) String() string {
	if stmt.Expr != nil {
		return stmt.Expr.String() + ";"
	}
	return "<nil-expr-stmt>;"
}

// IdentExpr is an identifier expression. There are places where
// identifiers are used as statements but this will be used in both.
type IdentExpr struct {
	Token *lexer.Token
	Value string
}

// NewIdentExpr returns a new identifier expression.
func NewIdentExpr(t *lexer.Token) *IdentExpr {
	return &IdentExpr{
		Token: t,
		Value: t.String,
	}
}

func (expr *IdentExpr) expression()          {}
func (expr *IdentExpr) TokenLiteral() string { return expr.Token.String }
func (expr *IdentExpr) String() string       { return expr.Value }

// IntLiteralExpr is an integer literal.
type IntLiteralExpr struct {
	Token *lexer.Token
	Value int64
}

// NewIntLiteralExpr creates a new integer literal expression.
func NewIntLiteralExpr(t *lexer.Token) *IntLiteralExpr {
	return &IntLiteralExpr{
		Token: t,
		Value: t.Int,
	}
}

func (expr *IntLiteralExpr) expression()          {}
func (expr *IntLiteralExpr) TokenLiteral() string { return expr.Token.String }
func (expr *IntLiteralExpr) String() string       { return expr.Token.String }

// PrefixExpr is an expression with a prefix operator.
type PrefixExpr struct {
	Token    *lexer.Token
	Operator string
	Expr     Expression
}

// NewPrefixExpr creates a new prefix expression.
func NewPrefixExpr(t *lexer.Token, operator string, e Expression) *PrefixExpr {
	return &PrefixExpr{
		Token:    t,
		Operator: operator,
		Expr:     e,
	}
}

func (expr *PrefixExpr) expression()          {}
func (expr *PrefixExpr) TokenLiteral() string { return expr.Token.String }

func (expr *PrefixExpr) String() string {
	var sb strings.Builder
	sb.WriteString(expr.Operator)
	sb.WriteString(expr.Expr.String())
	return sb.String()
}
