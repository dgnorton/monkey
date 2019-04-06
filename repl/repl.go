package repl

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/dgnorton/monkey/lexer"
)

type REPL struct {
	stdin  io.Reader
	stdout io.Writer
	stderr io.Writer
	prompt string
}

func New(stdin io.Reader, stdout, stderr io.Writer) *REPL {
	return &REPL{
		stdin:  stdin,
		stdout: stdout,
		stderr: stderr,
		prompt: ">> ",
	}
}

func (r *REPL) Start() (cancel func(), done chan struct{}) {
	stop := make(chan struct{})
	done = make(chan struct{})
	go r.loop(stop, done)
	cancel = func() { close(stop); <-done }
	return cancel, done
}

func (r *REPL) loop(stop chan struct{}, done chan struct{}) {
	defer close(done)
	for {
		// See if we're supposed to exit.
		select {
		case <-stop:
			return
		default:
		}

		// Read
		code, err := r.read()
		if err != nil {
			if err == io.EOF {
				return
			}

			fmt.Fprintln(r.stdout, err.Error())
			continue
		}

		// Eval
		results := r.eval(code, stop)

		// Print
		for result := range results {
			fmt.Fprintln(r.stdout, result)
		}
	}
}

func (r *REPL) read() (string, error) {
	fmt.Fprint(r.stdout, r.prompt)
	return bufio.NewReader(r.stdin).ReadString('\n')
}

func (r *REPL) eval(code string, stop chan struct{}) chan string {
	ch := make(chan string)
	go func() {
		defer close(ch)
		lex := lexer.New("", strings.NewReader(code))
		for {
			// See if we're supposed to exit.
			select {
			case <-stop:
			default:
			}

			tok, err := lex.Next()
			if err != nil {
				ch <- err.Error()
				return
			}

			ch <- fmt.Sprint(tok)

			if tok.EOF() {
				break
			}
		}
	}()
	return ch
}

func (r *REPL) print(result string) error {
	return nil
}
