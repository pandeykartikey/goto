package main

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/pandeykartikey/goto/eval"
	"github.com/pandeykartikey/goto/lexer"
	"github.com/pandeykartikey/goto/object"
	"github.com/pandeykartikey/goto/parser"
	"github.com/pandeykartikey/goto/repl"
)

func main() {

	if len(os.Args) > 2 {
		fmt.Println("Usage:", os.Args[0], "[FILE]")
		return
	}

	if len(os.Args) == 2 {
		code, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		l := lexer.New(string(code))
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			p.PrintParseErrors()
			return
		}
		env := object.NewEnvironment()

		result := eval.Eval(program, env)

		if result != nil {
			fmt.Println(result.Inspect())
		}
	} else {
		fmt.Printf("GoTo 0.1.0\n")

		repl.Start()
	}
}
