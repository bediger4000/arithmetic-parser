package tree

// Parse tree - a binary tree of objects of type Node,
// and associated utility functions and methods.

import (
	"bytes"
	"fmt"
	"io"

	"arithmetic-parser/lexer"
	"arithmetic-parser/value"
)

// Node has all elements exported, everything reaches inside instances
// of Node to find things out, or to change Left and Right. Private
// elements would cost me gross ol' getter and setter boilerplate.
type Node struct {
	Op     lexer.TokenType
	Lexeme string
	Left   *Node
	Right  *Node
}

// NewNode creates interior nodes of a parse tree, which will
// all have a +, -, *, / operator associated
// The lexer ADD_OP, MULT_OP, EXP_OP constants exist to
// have precedence levels with more than one operation at
// each precedence.
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
		n.Lexeme = lexeme
	}
	return &n
}

// UnaryNode handles "-something" and "+something" situtations.
// It returns "something" in "+something" cases,
// but sets up a "0 - something" sub-tree for unary negation.
func UnaryNode(unaryOp string, factor *Node) *Node {
	if unaryOp == "+" {
		return factor
	}
	return &Node{Op: lexer.MINUS, Left: &Node{Op: lexer.CONSTANT, Lexeme: "0"}, Right: factor}
}

// Eval recursively traverses a parse tree for an arithmetic expression.
// It uses type value.Value to do the numerical evaluation.
func (p *Node) Eval() value.Value {
	if p.Op == lexer.CONSTANT {
		return value.NewValue(p.Lexeme)
	}
	left := p.Left.Eval()
	right := p.Right.Eval()
	return left.BinaryOp(p.Op, right)
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
	case lexer.CONSTANT, lexer.POSITIVE, lexer.NEGATIVE:
		oper = 0
	}
	if oper != 0 {
		fmt.Fprintf(w, " %c ", oper)
	}

	if p.Op == lexer.CONSTANT {
		fmt.Fprintf(w, "%s", p.Lexeme)
	}

	if p.Op == lexer.NEGATIVE {
		fmt.Fprint(w, "-")
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

func (p *Node) String() string {
	var sb bytes.Buffer
	p.Print(&sb)
	return sb.String()
}

func (p *Node) graphNode(w io.Writer) {

	var label string

	switch p.Op {
	case lexer.CONSTANT:
		label = fmt.Sprintf("%s", p.Lexeme)
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
	case lexer.EXP:
		label = "^"
	case lexer.NEGATIVE:
		label = "~"
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
