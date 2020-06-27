package repl

import (
	"fmt"

	"github.com/peterh/liner"

	"goto/eval"
	"goto/lexer"
	"goto/object"
	"goto/parser"
)

const (
	PS1 = ">> "
	PS2 = "... "
)

func Start() {

	term := liner.NewLiner()
	defer term.Close()

	env := object.NewEnvironment()
	code := ""
	prompt := PS1

	for {

		line, err := term.Prompt(prompt)

		if err != nil {
			fmt.Println("Aborted")
			break
		}

		code += line

		l := lexer.New(code)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			if line == "" {
				printParseErrors(p.Errors())
				code = ""
				prompt = PS1
			} else {
				term.AppendHistory(line)
				prompt = PS2
			}
			continue
		}

		term.AppendHistory(line)
		code = ""
		prompt = PS1
		result := eval.Eval(program, env)

		if result != nil {
			fmt.Println(result.Inspect())
		}
	}
}

func printParseErrors(errors []string) {
	for _, msg := range errors {
		fmt.Println("Error: ", msg)
	}
}
