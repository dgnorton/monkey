package parser_test

import (
	"testing"

	"github.com/dgnorton/monkey/parser"
)

func TestParse(t *testing.T) {
	prog, err := parser.Parse("let x = 5;")
	if err != nil {
		t.Fatal(err)
	}
	_ = prog
	//fmt.Println(prog.TokenLiteral())
}
