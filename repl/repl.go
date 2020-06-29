package repl

import (
	"fmt"

	"github.com/peterh/liner"

	"github.com/pandeykartikey/goto/eval"
	"github.com/pandeykartikey/goto/lexer"
	"github.com/pandeykartikey/goto/object"
	"github.com/pandeykartikey/goto/parser"
)

const (
	PS1 = ">> "
	PS2 = "... "
)

func Start() {

	term := liner.NewLiner()
	defer term.Close()
	term.SetCtrlCAborts(true)

	env := object.NewEnvironment()
	code := ""
	prompt := PS1

	for {

		line, err := term.Prompt(prompt)

		if err != nil {
			if err == liner.ErrPromptAborted {
				code = ""
				prompt = PS1
				continue
			} else {
				fmt.Println("Aborted")
			}
			break
		}

		if line == "exit" {
			break
		}

		code += "\n" + line

		l := lexer.New(code)
		p := parser.New(l)

		program := p.ParseProgram()

		if len(p.Errors()) != 0 {
			if line == "" {
				p.PrintParseErrors()
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
