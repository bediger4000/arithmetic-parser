package parser

import (
	"arithmetic-parser/lexer"
	"arithmetic-parser/tree"
	"fmt"
)

/*
expr -> term   {add-op term}
term -> factor {mult-op factor}
factor -> '(' expr ')' | NUMBER
add-op -> '+'|'-'
mult-op -> '*'|'/'|'%'
*/

/*
* One parse method per non-terminal symbol
* A non-terminal symbol on the RHS of a rewrite rule
  leads to a call to the parse method for that non-terminal
* Terminal symbol on the RHS of a rewrite rule leads to "consuming"
  that token from the input token string
* | in the grammar leads to "if-else" in the parser
* {...} in the grammar leads to "while" in the parser
*/

type Parser struct {
	lexer *lexer.Lexer
}

func (p *Parser) Parse() *tree.Node {
	return p.expr()
}

func (p *Parser) expr() *tree.Node {
	node := p.term()
	for kind, lexeme := p.lexer.NextToken(); kind == lexer.ADD_OP; kind, lexeme = p.lexer.NextToken() {
		tmp := tree.NewNode(kind, lexeme)
		p.lexer.Consume()
		tmp.Left = node
		node = tmp
		node.Right = p.term()
	}
	return node
}

func (p *Parser) term() *tree.Node {
	node := p.factor()
	for kind, lexeme := p.lexer.NextToken(); kind == lexer.MULT_OP; kind, lexeme = p.lexer.NextToken() {
		tmp := tree.NewNode(kind, lexeme)
		p.lexer.Consume()
		tmp.Left = node
		node = tmp
		node.Right = p.factor()
	}
	return node

}

func (p *Parser) factor() *tree.Node {
	kind, lexeme := p.lexer.NextToken()
	switch kind {
	case lexer.CONSTANT:
		return tree.NewNode(kind, lexeme)
		p.lexer.Consume()
	case lexer.LPAREN:
		expr := p.expr()
		kind, lexeme = p.lexer.NextToken()
		if kind != lexer.RPAREN {
			fmt.Printf("Wanted an RPAREN, got %v: %q\n", kind, lexeme)
		}
		p.lexer.Consume()
		return expr
	default:
		fmt.Printf("Wanted a CONSTANT or LPAREN, got %v: %q\n", kind, lexeme)
	}
	return nil
}

func (p *Parser) addOp() *tree.Node {
	kind, lexeme := p.lexer.NextToken()
	if kind != lexer.ADD_OP {
		fmt.Printf("Wanted an ADD-OP, got %v: %q\n", kind, lexeme)
	}
	return tree.NewNode(kind, lexeme)
}

func (p *Parser) multOp() *tree.Node {
	kind, lexeme := p.lexer.NextToken()
	if kind != lexer.MULT_OP {
		fmt.Printf("Wanted an MULT-OP, got %v: %q\n", kind, lexeme)
	}
	return tree.NewNode(kind, lexeme)
}

func NewParser(lxr *lexer.Lexer) *Parser {
	return &Parser{lexer: lxr}
}