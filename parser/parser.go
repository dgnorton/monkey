package parser

import (
	"fmt"
	"strings"

	"github.com/dgnorton/monkey/ast"
	"github.com/dgnorton/monkey/lexer"
)

// Parse parses a string and returns an AST.
func Parse(code string) (*ast.Program, error) {
	l := lexer.New("", strings.NewReader(code))
	p := New(l)
	return p.Parse()
}

// ParseFile parses a file and returns an AST.
func ParseFile(filename string) (*ast.Program, error) {
	p, err := Open(filename)
	if err != nil {
		return nil, err
	}
	return p.Parse()
}

// Types of expression parsing funcs.
type prefixParseFn func() (ast.Expression, error)
type infixParseFn func(ast.Expression) (ast.Expression, error)

// Parser is a Monkey language parser.
type Parser struct {
	lex *lexer.Lexer

	// Parsing function lookup arrays, indexed by lexer.TokenType.
	prefixParseFns []prefixParseFn
	infixParseFns  []infixParseFn
}

// New returns a new Parser.
func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		lex:            l,
		prefixParseFns: make([]prefixParseFn, lexer.MAXTOKTYPE),
		infixParseFns:  make([]infixParseFn, lexer.MAXTOKTYPE),
	}

	p.prefixParseFns[lexer.IDENT] = func() (ast.Expression, error) { return p.identExpr() }
	p.prefixParseFns[lexer.INT] = func() (ast.Expression, error) { return p.intLiteralExpr() }
	p.prefixParseFns[lexer.SUB] = func() (ast.Expression, error) { return p.prefixExpr() }
	p.prefixParseFns[lexer.NOT] = func() (ast.Expression, error) { return p.prefixExpr() }

	return p
}

// Open opens a file and returns a parser for it.
func Open(filename string) (*Parser, error) {
	l, err := lexer.Open(filename)
	if err != nil {
		return nil, err
	}
	return New(l), nil
}

// Parse parses the output of its lexer and returns an AST.
func (p *Parser) Parse() (*ast.Program, error) {
	prog := ast.NewProgram()
	for {
		tok, err := p.lex.Peek()
		if err != nil {
			return nil, err
		}

		var stmt ast.Statement

		switch tok.Type {
		case lexer.LET:
			stmt, err = p.letStmt()
		case lexer.RETURN:
			stmt, err = p.returnStmt()
		case lexer.EOF:
			return prog, nil
		default:
			stmt, err = p.exprStmt()
			//err = p.parseErr(tok, fmt.Errorf("invalid token: %s", tok.String))
		}

		if err != nil {
			return nil, err
		}

		if stmt != nil {
			prog.AddStmt(stmt)
		}
	}
}

// letStmt parses a let statement.
func (p *Parser) letStmt() (*ast.LetStmt, error) {
	// "let"
	letTok, err := p.requireTok(lexer.LET)
	if err != nil {
		return nil, err
	}

	// Name identifier
	name, err := p.identExpr()

	// "="
	_, err = p.requireTok(lexer.ASSIGN)
	if err != nil {
		return nil, err
	}

	// Expression
	value, err := p.expr(LOWEST)
	if err != nil {
		return nil, err
	}

	// ";"
	_, err = p.optionalTok(lexer.SEMICOLON)
	if err != nil {
		return nil, err
	}

	return ast.NewLetStmt(letTok, name, value), nil
}

// returnStmt parses a return statement.
func (p *Parser) returnStmt() (*ast.ReturnStmt, error) {
	// "return"
	returnTok, err := p.requireTok(lexer.RETURN)
	if err != nil {
		return nil, err
	}

	// Expression
	value, err := p.expr(LOWEST)
	if err != nil {
		return nil, err
	}

	// ";"
	_, err = p.optionalTok(lexer.SEMICOLON)
	if err != nil {
		return nil, err
	}

	return ast.NewReturnStmt(returnTok, value), nil
}

// exprStmt parses an expression statement.
func (p *Parser) exprStmt() (*ast.ExprStmt, error) {
	// Peek first token of the expresson.
	tok, err := p.lex.Peek()
	if err != nil {
		return nil, p.parseErr(tok, err)
	}
	fmt.Printf("exprStmt: %v\n", tok)

	// Expression
	expr, err := p.expr(LOWEST)
	if err != nil {
		return nil, err
	}

	// ";"
	_, err = p.optionalTok(lexer.SEMICOLON)
	if err != nil {
		return nil, err
	}

	return ast.NewExprStmt(tok, expr), nil
}

type precedence int

const (
	_ precedence = iota
	LOWEST
	EQUALS      // ==
	LESSGREATER // > or <
	SUM         // +
	PRODUCT     // *
	PREFIX      // -x or !x
)

// expr parses an expression.
func (p *Parser) expr(pr precedence) (ast.Expression, error) {
	tok, err := p.lex.Peek()
	if err != nil {
		return nil, err
	}
	fmt.Printf("expr: %v\n", tok)

	exprFn := p.prefixParseFns[tok.Type]
	if exprFn == nil {
		return nil, nil
	}

	return exprFn()
}

// identExpr parses an identifier expression.
func (p *Parser) identExpr() (*ast.IdentExpr, error) {
	tok, err := p.requireTok(lexer.IDENT)
	if err != nil {
		return nil, err
	}
	fmt.Printf("identExpr: %v\n", tok)
	return ast.NewIdentExpr(tok), nil
}

// intLiteralExpr parses an integer literal expression.
func (p *Parser) intLiteralExpr() (*ast.IntLiteralExpr, error) {
	tok, err := p.requireTok(lexer.INT)
	if err != nil {
		return nil, err
	}
	fmt.Printf("intLiteralExpr: %v\n", tok)
	return ast.NewIntLiteralExpr(tok), nil
}

// prefixExpr parses a prefix expression
func (p *Parser) prefixExpr() (*ast.PrefixExpr, error) {
	tok, err := p.requireTok(lexer.NOT, lexer.SUB)
	if err != nil {
		return nil, err
	}

	expr, err := p.expr(PREFIX)
	if err != nil {
		return nil, err
	}

	fmt.Printf("prefixExpr: %v\n", tok)
	return ast.NewPrefixExpr(tok, tok.String, expr), nil
}

// requireTok returns the next token if it matches one of the expected token types.
// If the type does not match, an error is returned.
func (p *Parser) requireTok(expTypes ...lexer.TokenType) (*lexer.Token, error) {
	tok, err := p.lex.Next()
	if err != nil {
		return nil, err
	}

	for _, expType := range expTypes {
		if tok.Type == expType {
			return tok, nil
		}
	}

	err = fmt.Errorf("expected '!' or '-' but found: %s", tok.String)
	return nil, p.parseErr(tok, err)
}

// optionalTok returns the next token if it matches the expected token type.
// If the type does not match, nil is returned for the token and error.
func (p *Parser) optionalTok(expType lexer.TokenType) (*lexer.Token, error) {
	tok, err := p.lex.Peek()
	if err != nil {
		return nil, err
	}

	if tok.Type != expType {
		return nil, nil
	}

	p.lex.Next()

	return tok, nil
}

// parseErr returns a parser error for the given token.
func (p *Parser) parseErr(tok *lexer.Token, err error) *Error {
	return &Error{
		Err: err,
		Tok: tok,
	}
}

// Error represents a parse error.
type Error struct {
	Err error
	Tok *lexer.Token
}

// Error returns a string representation of the error.
func (e *Error) Error() string {
	filename := e.Tok.File
	if filename == "" {
		filename = "<no-file>"
	}
	return fmt.Sprintf("%s|%d col %d| %s", filename, e.Tok.Line, e.Tok.Col, e.Err)
}
