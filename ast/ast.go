package ast

import (
	"github.com/dgnorton/monkey/lexer"
)

// Node represents a node in the AST. All nodes implement this interface.
type Node interface {
	TokenLiteral() string
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

// NewProgram creates a new Program.
func NewProgram() *Program {
	return &Program{
		Statements: []Statement{},
	}
}

// AddStmt adds a statement to the program.
func (p *Program) AddStmt(stmt Statement) {
	p.Statements = append(p.Statements, stmt)
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) == 0 {
		return ""
	}
	return p.Statements[0].TokenLiteral()
}

// LetStmt is a let statement node.
type LetStmt struct {
	Token *lexer.Token
	Name  *IdentExpr
	Value Expression
}

// NewLetStmt returns a new LetStmt.
func NewLetStmt(t *lexer.Token, name *IdentExpr, value Expression) *LetStmt {
	return &LetStmt{
		Token: t,
		Name:  name,
		Value: value,
	}
}

func (stmt *LetStmt) statement()           {}
func (stmt *LetStmt) TokenLiteral() string { return stmt.Token.String }

// IdentExpr is an identifier expression. There are places where
// identifiers are used as statements but this will be used in both.
type IdentExpr struct {
	Token *lexer.Token
	Value string
}

// NewIdentExpr returns a new IdentExpr.
func NewIdentExpr(t *lexer.Token) *IdentExpr {
	return &IdentExpr{
		Token: t,
		Value: t.String,
	}
}

func (expr *IdentExpr) expression()          {}
func (expr *IdentExpr) TokenLiteral() string { return expr.Token.String }
