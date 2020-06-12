package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/adrian83/monkey/pkg/evaluator"
	"github.com/adrian83/monkey/pkg/lexer"
	"github.com/adrian83/monkey/pkg/object"
	"github.com/adrian83/monkey/pkg/parser"
)

const (
	prompt    = ">>"
	lineBreak = "\n"
)

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	env := object.NewEnvironment()

	for {
		fmt.Printf(prompt)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		l := lexer.New(line)
		p := parser.New(l)

		program, err := p.ParseProgram()
		if err != nil {
			io.WriteString(out, err.Error())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, lineBreak)
		}
	}
}
