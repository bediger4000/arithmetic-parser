package main

import (
	"arithmetic-parser/lexer"
	"arithmetic-parser/parser"
	"fmt"
	"os"
)

func main() {
	dotFile := false
	str := os.Args[1]
	if str == "-g" {
		dotFile = true
		str = os.Args[2]
	}

	lxr := lexer.Lex(str)
	psr := parser.NewParser(lxr)

	tree := psr.Parse()

	if dotFile {
		tree.GraphNode(os.Stdout)
	} else {
		fmt.Printf("Reconstituted expression: %q\n", tree)
		fmt.Printf("/* %s */\n", tree.Eval())
	}
}
