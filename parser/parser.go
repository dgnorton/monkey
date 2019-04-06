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

// Parser is a Monkey language parser.
type Parser struct {
	lex *lexer.Lexer
}

// New returns a new Parser.
func New(l *lexer.Lexer) *Parser {
	return &Parser{
		lex: l,
	}
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
		case lexer.EOF:
			return prog, nil
		default:
			err = fmt.Errorf("invalid token: %s", tok.String)
		}

		if err != nil {
			return nil, err
		}

		if stmt != nil {
			prog.AddStmt(stmt)
		}
	}
}

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
	value, err := p.expr()
	if err != nil {
		return nil, err
	}

	// ";"
	_, err = p.requireTok(lexer.SEMICOLON)
	if err != nil {
		return nil, err
	}

	return ast.NewLetStmt(letTok, name, value), nil
}

func (p *Parser) expr() (ast.Expression, error) {
	for {
		tok, err := p.lex.Peek()
		if err != nil {
			return nil, p.parseErr(tok, err)
		}
		if tok.Type == lexer.SEMICOLON || tok.EOF() {
			return nil, nil
		}

		p.lex.Next()
	}
}

func (p *Parser) identExpr() (*ast.IdentExpr, error) {
	tok, err := p.requireTok(lexer.IDENT)
	if err != nil {
		return nil, err
	}
	return ast.NewIdentExpr(tok), nil
}

func (p *Parser) requireTok(expType lexer.TokenType) (*lexer.Token, error) {
	tok, err := p.lex.Next()
	if err != nil {
		return nil, err
	}

	if tok.Type != expType {
		return nil, p.parseErr(tok, err)
	}

	return tok, nil
}

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
	return fmt.Sprintf("%s|%d col %d| %s", e.Tok.File, e.Tok.Line, e.Tok.Col, e.Err)
}
