package main

import (
	"fmt"
	"os"

	"goto/repl"
)

func main() {
	fmt.Printf("GoTo 1.0.0\n") // TODO: Add a help message

	repl.Start(os.Stdin, os.Stdout)
}
