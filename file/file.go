package file

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/waits/tang/eval"
	"github.com/waits/tang/lexer"
	"github.com/waits/tang/object"
	"github.com/waits/tang/parser"
)

func Exec(w io.Writer, path string, args []string) {
	f, err := os.Open(path)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	text, err := ioutil.ReadAll(f)
	if err != nil {
		fmt.Fprintln(w, err)
		return
	}

	l := lexer.New(string(text))
	p := parser.New(l)

	program := p.ParseProgram()
	if len(p.Errors()) != 0 {
		printParserErrors(w, p.Errors())
		return
	}

	env := object.NewEnv()
	output := eval.Eval(program, env)
	if output != nil {
		fmt.Fprintln(w, output.Inspect())
	}
}

func printParserErrors(w io.Writer, errors []string) {
	for _, msg := range errors {
		fmt.Fprintf(w, "\t%s\n", msg)
	}
}
