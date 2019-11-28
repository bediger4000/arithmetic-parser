package tree

// Parse tree - a binary tree of objects of type Node,
// and associated utility functions and methods.

import (
	"arithmetic-parser/lexer"
	"bytes"
	"fmt"
	"io"
	"strconv"
)

// Node has all elements exported, everything reaches inside instances
// of Node to find things out, or to change Left and Right. Private
// elements would cost me gross ol' getter and setter boilerplate.
type Node struct {
	Op    lexer.TokenType
	Const int
	Left  *Node
	Right *Node
}

// NewNode creates interior nodes of a parse tree, which will
// all have a +, -, *, / operator associated
func NewNode(op lexer.TokenType, lexeme string) *Node {
	var n Node
	switch op {
	case lexer.ADD_OP, lexer.MULT_OP, lexer.EXP_OP:
		switch lexeme {
		case "+":
			n.Op = lexer.PLUS
		case "-":
			n.Op = lexer.MINUS
		case "*":
			n.Op = lexer.MULT
		case "/":
			n.Op = lexer.DIV
		case "^":
			n.Op = lexer.EXP
		case "%":
			n.Op = lexer.REM
		}
	case lexer.CONSTANT:
		n.Op = lexer.CONSTANT
		n.Const, _ = strconv.Atoi(lexeme)
	}
	return &n
}

func (p *Node) Eval() *Node {
	switch p.Op {
	case lexer.CONSTANT:
		return p
	case lexer.PLUS:
		return &Node{Const: p.Left.Eval().Const + p.Right.Eval().Const}
	case lexer.MINUS:
		return &Node{Const: p.Left.Eval().Const - p.Right.Eval().Const}
	case lexer.DIV:
		return &Node{Const: p.Left.Eval().Const / p.Right.Eval().Const}
	case lexer.MULT:
		return &Node{Const: p.Left.Eval().Const * p.Right.Eval().Const}
	case lexer.REM:
		return &Node{Const: p.Left.Eval().Const % p.Right.Eval().Const}
	case lexer.EXP:
		exponent := p.Right.Eval().Const
		if exponent == 0 {
			return &Node{Const: 1}
		}
		// This is wrong
		if exponent < 0 {
			return &Node{Const: 0}
		}
		base := p.Left.Eval().Const
		answer := 1
		for ; exponent > 0; exponent-- {
			answer *= base
		}
		return &Node{Const: answer}
		// what about fractional powers (1/2 == square root)
	}
	return nil
}

// Print puts a human-readable, nicely formatted string representation
// of a parse tree onto the io.Writer, w.  Essentially just an in-order
// traversal of a binary tree, with accommodating a few oddities, like
// parenthesization, and the "~" (not) operator being a prefix.
func (p *Node) Print(w io.Writer) {

	if p.Left != nil {
		printParen := false
		if p.Left.Op != lexer.CONSTANT {
			fmt.Fprintf(w, "(")
			printParen = true
		}
		p.Left.Print(w)
		if printParen {
			fmt.Fprintf(w, ")")
		}
	}

	var oper rune
	switch p.Op {
	case lexer.MULT:
		oper = '*'
	case lexer.DIV:
		oper = '/'
	case lexer.REM:
		oper = '%'
	case lexer.EXP:
		oper = '^'
	case lexer.PLUS:
		oper = '+'
	case lexer.MINUS:
		oper = '-'
	case lexer.CONSTANT:
		oper = 0
	}
	if oper != 0 {
		fmt.Fprintf(w, " %c ", oper)
	}

	if p.Op == lexer.CONSTANT {
		fmt.Fprintf(w, "%d", p.Const)
	}

	if p.Right != nil {
		printParen := false
		if p.Right.Op != lexer.CONSTANT {
			fmt.Fprintf(w, "(")
			printParen = true
		}
		p.Right.Print(w)
		if printParen {
			fmt.Fprintf(w, ")")
		}
	}
}

// ExpressionToString creates a Golang string with a human readable
// representation of a parse tree in it.
func ExpressionToString(root *Node) string {
	var sb bytes.Buffer
	root.Print(&sb)
	return sb.String()
}

func (p *Node) String() string {
	return ExpressionToString(p)
}

func (p *Node) graphNode(w io.Writer) {

	var label string

	switch p.Op {
	case lexer.CONSTANT:
		label = fmt.Sprintf("%d", p.Const)
	case lexer.MINUS:
		label = "-"
	case lexer.PLUS:
		label = "+"
	case lexer.DIV:
		label = "/"
	case lexer.MULT:
		label = "*"
	case lexer.REM:
		label = "%"
	}

	fmt.Fprintf(w, "n%p [label=\"%s\"];\n", p, label)

	if p.Left != nil {
		p.Left.graphNode(w)
		fmt.Fprintf(w, "n%p -> n%p;\n", p, p.Left)
	}
	if p.Right != nil {
		p.Right.graphNode(w)
		fmt.Fprintf(w, "n%p -> n%p;\n", p, p.Right)
	}
}

// GraphNode puts a dot-format text representation of
// a parse tree on w io.Writer.
func (p *Node) GraphNode(w io.Writer) {
	fmt.Fprintf(w, "digraph g {\n")
	p.graphNode(w)
	fmt.Fprintf(w, "}\n")
}
