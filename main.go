package main

import (
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/pandeykartikey/goto/eval"
	"github.com/pandeykartikey/goto/lexer"
	"github.com/pandeykartikey/goto/object"
	"github.com/pandeykartikey/goto/parser"
	"github.com/pandeykartikey/goto/repl"
)

var filename = flag.String("f", "", "file to run")

func main() {
	flag.Parse()

	if *filename != "" {
		code, err := ioutil.ReadFile(*filename)
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
		fmt.Printf("GoTo 1.0.0\n")

		repl.Start()
	}
}
