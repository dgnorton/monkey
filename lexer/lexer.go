package lexer

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

// Lexer lexes / tokenizes Monkey language.
type Lexer struct {
	filename string
	r        *bufio.Reader
	closer   io.Closer

	line int
	col  int

	currune rune
	nxtrune rune

	curtok *Token
	nxttok *Token
}

// New returns a new instance of a Monkey language lexer.
func New(filename string, r io.Reader) *Lexer {
	closer, _ := r.(io.Closer)
	return &Lexer{
		filename: filename,
		r:        bufio.NewReader(r),
		closer:   closer,
		line:     1,
		col:      1,
	}
}

// Open opens a Monkey language script file and returns a lexer.
func Open(filename string) (*Lexer, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	return New(filename, f), nil
}

// Close closes the lexer and underlying reader if it supports closing.
func (l *Lexer) Close() error {
	if l.closer != nil {
		return l.closer.Close()
	}
	return nil
}

// Next returns the next Token from the input.
func (l *Lexer) Next() (*Token, error) {
	if l.nxttok != nil {
		l.curtok = l.nxttok
		l.nxttok = nil
		return l.curtok, nil
	}

	tok, err := l.readTok()
	if err != nil {
		return nil, err
	}

	l.curtok = tok

	return tok, nil
}

// Peak returns the next Token without reading past it.
func (l *Lexer) Peak() (*Token, error) {
	if l.nxttok != nil {
		return l.nxttok, nil
	}

	tok, err := l.readTok()
	if err != nil {
		return nil, err
	}

	l.nxttok = tok

	return tok, nil
}

// readTok reads in the next token from input.
func (l *Lexer) readTok() (*Token, error) {
	if err := l.skipSpace(); err != nil {
		if err != io.EOF {
			return nil, l.lexErr(err)
		}
		return l.newTok(EOF, "")
	}

	r, err := l.peakRune()
	if err != nil {
		if err != io.EOF {
			return nil, l.lexErr(err)
		}
		return l.newTok(EOF, "")
	}

	switch r {
	case ';':
		r, _ = l.readRune()
		return l.newTok(SEMICOLON, ";")
	case '=':
		r, _ = l.readRune()
		return l.newTok(ASSIGN, "=")
	case '+':
		r, _ = l.readRune()
		return l.newTok(PLUS, "+")
	case '(':
		r, _ = l.readRune()
		return l.newTok(LPAREN, "(")
	case ')':
		r, _ = l.readRune()
		return l.newTok(RPAREN, ")")
	case '{':
		r, _ = l.readRune()
		return l.newTok(LBRACE, "{")
	case '}':
		r, _ = l.readRune()
		return l.newTok(RBRACE, "}")
	case '[':
		r, _ = l.readRune()
		return l.newTok(LSQUARE, "[")
	case ']':
		r, _ = l.readRune()
		return l.newTok(RSQUARE, "]")
	case ',':
		r, _ = l.readRune()
		return l.newTok(COMMA, ",")
	default:
		if isLetter(r) {
			return l.readIdentTok()
		} else if isDigit(r) {
			return l.readNumTok()
		}
		return l.newTok(ILLEGAL, string(r))
	}
}

// readIdentTok reads and returns an identifier token.
func (l *Lexer) readIdentTok() (*Token, error) {
	startCol := l.col
	var sb strings.Builder

	for {
		r, err := l.peakRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, l.lexErr(err)
		}

		if !isLetter(r) && !isDigit(r) {
			break
		}

		r, _ = l.readRune()
		if _, err := sb.WriteRune(r); err != nil {
			return nil, l.lexErr(err)
		}
	}

	ident := sb.String()
	tokType := lookupIdentType(ident)

	tok, _ := l.newTok(tokType, ident)
	tok.Col = startCol

	return tok, nil
}

// readNumTok reads and returns any supported type of number token.
func (l *Lexer) readNumTok() (*Token, error) {
	// TODO: add float support?
	return l.readIntTok()
}

// readIntTok reads and returns an integer token.
func (l *Lexer) readIntTok() (*Token, error) {
	var sb strings.Builder

	for {
		r, err := l.peakRune()
		if err != nil {
			if err == io.EOF {
				break
			}
			return nil, l.lexErr(err)
		}

		if !isDigit(r) {
			break
		}

		r, _ = l.readRune()
		if _, err := sb.WriteRune(r); err != nil {
			return nil, l.lexErr(err)
		}
	}

	i, err := strconv.Atoi(sb.String())
	if err != nil {
		return nil, err
	}

	tok, _ := l.newTok(INT, sb.String())
	tok.Int = i

	return tok, nil
}

// newToken returns a new Token.
func (l *Lexer) newTok(t TokenType, s string) (*Token, error) {
	var err error
	if t == ILLEGAL {
		err = l.lexErr(fmt.Errorf("invalid token: %s", s))
	}
	return NewToken(t, l.filename, l.line, l.col-1, s), err
}

// readRune returns the next rune.
func (l *Lexer) readRune() (rune, error) {
	if l.nxtrune != 0 {
		l.currune = l.nxtrune
		l.nxtrune = 0
	} else {
		r, _, err := l.r.ReadRune()
		if err != nil {
			return 0, err
		}

		l.currune = r
	}

	l.col++
	if l.currune == '\n' {
		l.line++
		l.col = 1
	}

	return l.currune, nil
}

// peakRune returns the next rune that would be returned by readRune() but
// leaves it in the input stream.
func (l *Lexer) peakRune() (rune, error) {
	if l.nxtrune != 0 {
		return l.nxtrune, nil
	}

	r, _, err := l.r.ReadRune()
	if err != nil {
		return 0, err
	}

	l.nxtrune = r

	return r, nil
}

// skipSpace advances the lexer past whitespace.
func (l *Lexer) skipSpace() error {
	for {
		r, err := l.peakRune()
		if err != nil {
			return err
		}

		if !unicode.IsSpace(r) {
			return nil
		}

		r, _ = l.readRune()
	}
}

// lexErr returns a lexer error.
func (l *Lexer) lexErr(err error) error {
	return &Error{
		Err:  err,
		File: l.filename,
		Line: l.line,
		Col:  l.col,
	}
}

// TokenType is used by the Token struct to distinguish which type of token
// it holds.
type TokenType int

const (
	ILLEGAL TokenType = iota
	EOF

	// Identifiers and literals
	IDENT
	INT

	// Operators
	ASSIGN // '='
	PLUS   // '+'

	// Delimeters
	SEMICOLON // ';'
	LPAREN    // '('
	RPAREN    // ')'
	LBRACE    // '{'
	RBRACE    // '}'
	LSQUARE   // '['
	RSQUARE   // ']'
	COMMA     // ','

	// Keywords
	FN
	LET
	RETURN
)

// String returns a string representation of the token type.
func (t TokenType) String() string {
	switch t {
	case ILLEGAL:
		return "ILLEGAL"
	case EOF:
		return "EOF"
	case IDENT:
		return "IDENT"
	case INT:
		return "INT"
	case ASSIGN:
		return "ASSIGN"
	case PLUS:
		return "PLUS"
	case SEMICOLON:
		return "SEMICOLON"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case LBRACE:
		return "LBRACE"
	case RBRACE:
		return "RBRACE"
	case LSQUARE:
		return "LSQUARE"
	case RSQUARE:
		return "RSQUARE"
	case COMMA:
		return "COMMA"
	case FN:
		return "FN"
	case LET:
		return "LET"
	case RETURN:
		return "RETURN"
	default:
		return "INVALID TOKEN TYPE"
	}
}

// Token represents a single token in a Monkey language script.
type Token struct {
	Type   TokenType
	File   string
	Line   int
	Col    int
	String string
	Int    int
}

// NewToken returns a new Token with only the String value set.
func NewToken(t TokenType, filename string, line, col int, s string) *Token {
	return &Token{
		Type:   t,
		File:   filename,
		Line:   line,
		Col:    col,
		String: s,
	}
}

// NewIntToken returns a new Token with both the String and equivalent Int
// values set.
func NewIntToken(filename string, line, col int, s string) (*Token, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return nil, err
	}

	return &Token{
		Type:   INT,
		File:   filename,
		Line:   line,
		Col:    col,
		String: s,
		Int:    i,
	}, nil
}

// EOF returns true if the token is an EOF token.
func (t *Token) EOF() bool { return t.Type == EOF }

// Error represents a lexer error.
type Error struct {
	Err  error
	File string
	Line int
	Col  int
}

// Error returns a string representation of the error.
func (e *Error) Error() string {
	return fmt.Sprintf("%s|%d col %d| %s", e.File, e.Line, e.Col, e.Err)
}

// isLetter returns true if the rune is a valid identifier character.
func isLetter(r rune) bool {
	return 'a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || r == '_' || r >= utf8.RuneSelf && unicode.IsLetter(r)
}

// isDigit returns true if the rune is a valid numeric digit.
func isDigit(r rune) bool {
	return '0' <= r && r <= '9' || r >= utf8.RuneSelf && unicode.IsDigit(r)
}

// keywords is a map of Monkey language keywords to token types.
var keywords = map[string]TokenType{
	"fn":     FN,
	"let":    LET,
	"return": RETURN,
}

// lookupIdentType returns the keyword token type for the identifier or IDENT
// if the identifier isn't a Monkey language keyword.
func lookupIdentType(ident string) TokenType {
	if tokType, ok := keywords[ident]; ok {
		return tokType
	}
	return IDENT
}
