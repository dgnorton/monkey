package repl_test

import (
	"io"
	"strings"
	"testing"
	"time"

	"github.com/dgnorton/monkey/repl"
)

func TestREPL(t *testing.T) {
	stdin := strings.NewReader("let x = 42;\n")
	stdoutr, stdoutw := io.Pipe()
	_, stderrw := io.Pipe()

	// Start the Monkey REPL with some input.
	mky := repl.New(stdin, stdoutw, stderrw)
	cancel, _ := mky.Start()
	defer cancel()

	outch := make(chan byte)
	errch := make(chan error)
	go func() {
		var buf [1]byte
		for {
			if _, err := stdoutr.Read(buf[:]); err != nil {
				if err != io.EOF {
					errch <- err
				}
				return
			}
			outch <- buf[0]
		}
	}()

	exp := `>> &{LET  1 1 let 0}
&{IDENT  1 5 x 0}
&{ASSIGN  1 7 = 0}
&{INT  1 9 42 42}
&{SEMICOLON  1 11 ; 0}
&{EOF  2 0  0}
>> `

	// Wait for REPL results.
	var got strings.Builder
loop:
	for {
		select {
		case err := <-errch:
			if err != nil {
				t.Fatal(err)
			}
			break loop
		case b := <-outch:
			got.Write([]byte{b})
			if got.String() == exp {
				break loop
			}
		case <-time.After(time.Second * 1):
			if got.String() != exp {
				t.Fatalf("\nexp: %s\ngot: %s", exp, got.String())
			}
			t.Fatal("timed out")
		}
	}

}
