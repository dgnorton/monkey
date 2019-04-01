package lexer_test

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/dgnorton/monkey/lexer"
)

func TestLexer(t *testing.T) {
	code := `let add = fn(丢, b) {
  return 丢 + b;
};`

	dir, file := mustWriteTempFile("", code, t)
	defer os.RemoveAll(dir)

	exp := []*lexer.Token{
		&lexer.Token{
			Type:   lexer.LET,
			File:   file,
			Line:   1,
			Col:    1,
			String: "let",
		},
		&lexer.Token{
			Type:   lexer.IDENT,
			File:   file,
			Line:   1,
			Col:    5,
			String: "add",
		},
		&lexer.Token{
			Type:   lexer.ASSIGN,
			File:   file,
			Line:   1,
			Col:    9,
			String: "=",
		},
		&lexer.Token{
			Type:   lexer.FN,
			File:   file,
			Line:   1,
			Col:    11,
			String: "fn",
		},
		&lexer.Token{
			Type:   lexer.LPAREN,
			File:   file,
			Line:   1,
			Col:    13,
			String: "(",
		},
		&lexer.Token{
			Type:   lexer.IDENT,
			File:   file,
			Line:   1,
			Col:    14,
			String: "丢",
		},
		&lexer.Token{
			Type:   lexer.COMMA,
			File:   file,
			Line:   1,
			Col:    15,
			String: ",",
		},
		&lexer.Token{
			Type:   lexer.IDENT,
			File:   file,
			Line:   1,
			Col:    17,
			String: "b",
		},
		&lexer.Token{
			Type:   lexer.RPAREN,
			File:   file,
			Line:   1,
			Col:    18,
			String: ")",
		},
		&lexer.Token{
			Type:   lexer.LBRACE,
			File:   file,
			Line:   1,
			Col:    20,
			String: "{",
		},
		&lexer.Token{
			Type:   lexer.RETURN,
			File:   file,
			Line:   2,
			Col:    3,
			String: "return",
		},
		&lexer.Token{
			Type:   lexer.IDENT,
			File:   file,
			Line:   2,
			Col:    10,
			String: "丢",
		},
		&lexer.Token{
			Type:   lexer.PLUS,
			File:   file,
			Line:   2,
			Col:    12,
			String: "+",
		},
		&lexer.Token{
			Type:   lexer.IDENT,
			File:   file,
			Line:   2,
			Col:    14,
			String: "b",
		},
		&lexer.Token{
			Type:   lexer.SEMICOLON,
			File:   file,
			Line:   2,
			Col:    15,
			String: ";",
		},
		&lexer.Token{
			Type:   lexer.RBRACE,
			File:   file,
			Line:   3,
			Col:    1,
			String: "}",
		},
		&lexer.Token{
			Type:   lexer.SEMICOLON,
			File:   file,
			Line:   3,
			Col:    2,
			String: ";",
		},
		&lexer.Token{
			Type:   lexer.EOF,
			File:   file,
			Line:   3,
			Col:    2,
			String: "",
		},
	}

	lex, err := lexer.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	defer lex.Close()

	for i := 0; ; i++ {
		tok, err := lex.Next()
		if err != nil {
			t.Fatal(err)
		}

		if !reflect.DeepEqual(exp[i], tok) {
			t.Fatalf("tokens don't match:\nexp:\n\t%v\ngot:\n\t%v", exp[i], tok)
		}

		if tok.EOF() {
			break
		}
	}
}

func mustTempDir(t *testing.T) string {
	t.Helper()
	dir, err := ioutil.TempDir("", "monkey_lexer")
	if err != nil {
		panic(err)
	}
	return dir
}

func mustWriteTempFile(dir, s string, t *testing.T) (string, string) {
	t.Helper()

	if dir == "" {
		dir = mustTempDir(t)
	}

	f, err := ioutil.TempFile(dir, "")
	if err != nil {
		panic(err)
	}
	defer f.Close()

	if _, err = f.Write([]byte(s)); err != nil {
		panic(err)
	}

	return dir, f.Name()
}
