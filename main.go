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

	fmt.Printf("Reconstituted expression: %q\n", tree)
	fmt.Printf("/* %d */\n", tree.Eval().Const)

	if dotFile {
		tree.GraphNode(os.Stdout)
	}
}
