package parser_test

import (
	"reflect"
	"strings"
	"testing"

	"github.com/dgnorton/monkey/parser"
)

func TestParser_LetStmt(t *testing.T) {
	code := "let x = 5;"
	prog, err := parser.Parse(code)
	if err != nil {
		t.Fatal(err)
	}

	EQ(1, len(prog.Statements), t)
	EQ(code, prog.String(), t)
}

func TestParser_ReturnStmt(t *testing.T) {
	code := "return 5;"
	prog, err := parser.Parse(code)
	if err != nil {
		t.Fatal(err)
	}

	EQ(1, len(prog.Statements), t)
	EQ(code, prog.String(), t)
}

func TestParser_ExprStmt(t *testing.T) {
	code := `foo;
5;
-bar;
!baz;
`

	prog, err := parser.Parse(code)
	if err != nil {
		t.Fatal(err)
	}

	expNumStmts := len(strings.Split(code, "\n")) - 1

	EQ(expNumStmts, len(prog.Statements), t)
	EQ(code, prog.String(), t)
}

func EQ(exp, got interface{}, t *testing.T) {
	t.Helper()

	sameKind(exp, got, t)

	switch exp.(type) {
	case string, int, int64, uint, uint64:
		if exp != got {
			t.Fatalf("\nexp: %v\ngot: %v", exp, got)
		}
	default:
		t.Fatalf("don't know how to test type: %T", exp)
	}
}

func sameKind(exp, got interface{}, t *testing.T) {
	t.Helper()

	expTyp := reflect.TypeOf(exp)
	gotTyp := reflect.TypeOf(got)
	if expTyp.Kind() != gotTyp.Kind() {
		t.Fatalf("\nexp type: %T\ngot type: %T", exp, got)
	}
}
