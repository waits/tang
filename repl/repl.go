package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/waits/tang/eval"
	"github.com/waits/tang/lexer"
	"github.com/waits/tang/object"
	"github.com/waits/tang/parser"
)

const PROMPT = ">>> "

func Start(r io.Reader, w io.Writer) {
	scanner := bufio.NewScanner(r)
	env := object.NewEnv()

	for {
		fmt.Fprint(w, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(w, p.Errors())
			continue
		}

		evaluated := eval.Eval(program, env)
		if evaluated != nil {
			io.WriteString(w, evaluated.Inspect())
			io.WriteString(w, "\n")
		}
	}
}

func printParserErrors(w io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(w, "\t"+msg+"\n")
	}
}
