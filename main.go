package main

import (
	"fmt"
	"os"

	"pyro/repl"
)

func main() {
	fmt.Printf("Pyro 1.0.0\n") // TODO: Add a help message

	repl.Start(os.Stdin, os.Stdout)
}
