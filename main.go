package main

import (
	"arithmetic-parser/lexer"
	"arithmetic-parser/parser"
	"fmt"
	"os"
)

func main() {
	lxr := lexer.Lex("bob", os.Args[1])
	psr := parser.NewParser(lxr)

	tree := psr.Parse()

	fmt.Printf("Reconstituted expression: %s\n", tree)
}
