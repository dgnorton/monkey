package lexer_test

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"

	"github.com/dgnorton/monkey/lexer"
)

func TestLexer_Next(t *testing.T) {
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
			Type:   lexer.ADD,
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

func TestLexer_Peak(t *testing.T) {
	code := `let add`

	dir, file := mustWriteTempFile("", code, t)
	defer os.RemoveAll(dir)

	tok1 := &lexer.Token{
		Type:   lexer.LET,
		File:   file,
		Line:   1,
		Col:    1,
		String: "let",
	}

	tok2 := &lexer.Token{
		Type:   lexer.IDENT,
		File:   file,
		Line:   1,
		Col:    5,
		String: "add",
	}

	tok3 := &lexer.Token{
		Type:   lexer.EOF,
		File:   file,
		Line:   1,
		Col:    7,
		String: "",
	}

	lex, err := lexer.Open(file)
	if err != nil {
		t.Fatal(err)
	}
	defer lex.Close()

	// Test peaking the first token in the input.
	tok, err := lex.Peak()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(tok, tok1) {
		t.Fatalf("tokens don't match:\n\texp: %v\n\tgot: %v", tok1, tok)
	}

	// Read the first token from input and make sure it's the same
	// as we peaked before.
	if tok, err = lex.Next(); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(tok, tok1) {
		t.Fatalf("tokens don't match:\n\texp: %v\n\tgot: %v", tok1, tok)
	}

	// Peak the second token in the input.
	tok, err = lex.Peak()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(tok, tok2) {
		t.Fatalf("tokens don't match:\n\texp: %v\n\tgot: %v", tok2, tok)
	}

	// Peak the second token again and make sure it hasn't changed
	// since peaking it before.
	tok, err = lex.Peak()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(tok, tok2) {
		t.Fatalf("tokens don't match:\n\texp: %v\n\tgot: %v", tok2, tok)
	}

	// Read the second token in the input and make sure it's the same
	// as the token peaked before.
	if tok, err = lex.Next(); err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(tok, tok2) {
		t.Fatalf("tokens don't match:\n\texp: %v\n\tgot: %v", tok2, tok)
	}

	// Peak the last (EOF) token from the input.
	tok, err = lex.Peak()
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(tok, tok3) {
		t.Fatalf("tokens don't match:\n\texp: %v\n\tgot: %v", tok3, tok)
	}
}

func TestLexer_Operators(t *testing.T) {
	code := `=+-*/!<>`

	dir, file := mustWriteTempFile("", code, t)
	defer os.RemoveAll(dir)

	exps := []*lexer.Token{
		&lexer.Token{
			Type:   lexer.ASSIGN,
			File:   file,
			Line:   1,
			Col:    1,
			String: "=",
		},
		&lexer.Token{
			Type:   lexer.ADD,
			File:   file,
			Line:   1,
			Col:    2,
			String: "+",
		},
		&lexer.Token{
			Type:   lexer.SUB,
			File:   file,
			Line:   1,
			Col:    3,
			String: "-",
		},
		&lexer.Token{
			Type:   lexer.MUL,
			File:   file,
			Line:   1,
			Col:    4,
			String: "*",
		},
		&lexer.Token{
			Type:   lexer.DIV,
			File:   file,
			Line:   1,
			Col:    5,
			String: "/",
		},
		&lexer.Token{
			Type:   lexer.NOT,
			File:   file,
			Line:   1,
			Col:    6,
			String: "!",
		},
		&lexer.Token{
			Type:   lexer.LT,
			File:   file,
			Line:   1,
			Col:    7,
			String: "<",
		},
		&lexer.Token{
			Type:   lexer.GT,
			File:   file,
			Line:   1,
			Col:    8,
			String: ">",
		},
		&lexer.Token{
			Type:   lexer.EOF,
			File:   file,
			Line:   1,
			Col:    8,
			String: "",
		},
	}

	lex, err := lexer.Open(file)
	if err != nil {
		t.Fatal(err)
	}

	gots := make([]*lexer.Token, 0, len(exps))
	for i := 0; ; i++ {
		got, err := lex.Next()
		if err != nil {
			t.Fatal(err)
		}
		gots = append(gots, got)
		if !reflect.DeepEqual(exps[i], got) {
			t.Fatalf("tokens don't match:\n\texp: %v\n\tgot: %v", exps[i], got)
		}

		if got.EOF() {
			break
		}
	}

	if len(gots) != len(exps) {
		t.Fatalf("exp %d tokens, got %d", len(exps), len(gots))
	}
}

func TestLexer_Keywords(t *testing.T) {
	code := `fn let true false if else return`

	dir, file := mustWriteTempFile("", code, t)
	defer os.RemoveAll(dir)

	exps := []*lexer.Token{
		&lexer.Token{
			Type:   lexer.FN,
			File:   file,
			Line:   1,
			Col:    1,
			String: "fn",
		},
		&lexer.Token{
			Type:   lexer.LET,
			File:   file,
			Line:   1,
			Col:    4,
			String: "let",
		},
		&lexer.Token{
			Type:   lexer.TRUE,
			File:   file,
			Line:   1,
			Col:    8,
			String: "true",
		},
		&lexer.Token{
			Type:   lexer.FALSE,
			File:   file,
			Line:   1,
			Col:    13,
			String: "false",
		},
		&lexer.Token{
			Type:   lexer.IF,
			File:   file,
			Line:   1,
			Col:    19,
			String: "if",
		},
		&lexer.Token{
			Type:   lexer.ELSE,
			File:   file,
			Line:   1,
			Col:    22,
			String: "else",
		},
		&lexer.Token{
			Type:   lexer.RETURN,
			File:   file,
			Line:   1,
			Col:    27,
			String: "return",
		},
		&lexer.Token{
			Type:   lexer.EOF,
			File:   file,
			Line:   1,
			Col:    32,
			String: "",
		},
	}

	lex, err := lexer.Open(file)
	if err != nil {
		t.Fatal(err)
	}

	gots := make([]*lexer.Token, 0, len(exps))
	for i := 0; ; i++ {
		got, err := lex.Next()
		if err != nil {
			t.Fatal(err)
		}
		gots = append(gots, got)
		if !reflect.DeepEqual(exps[i], got) {
			t.Fatalf("tokens don't match:\n\texp: %v\n\tgot: %v", exps[i], got)
		}

		if got.EOF() {
			break
		}
	}

	if len(gots) != len(exps) {
		t.Fatalf("exp %d tokens, got %d", len(exps), len(gots))
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
