# Arithmetic Expressions, lexer and parser

[![Go Report Card](https://goreportcard.com/badge/github.com/bediger4000/arithmetic-parser)](https://goreportcard.com/report/github.com/bediger4000/arithmetic-parser)

Another cut at a Golang algebraic-order-of-operations arithmetic
expression evaluator.

Compare to [my earlier incarnation](https://github.com/bediger4000/arithmetic-expressions)
which was more-or-less a Go transliteration of a C program.

This version does share most of the implementation of the parse tree
with my earlier incarnation.
It has an idiomatic Go arithmetic evaluation.

## Daily Coding Problem: Problem #974 [Hard]

This repo is a good, albeit over-elaborate, solution to a Daily Coding Problem.

---
This problem was asked by Facebook.

Given a string consisting of parentheses,
single digits,
and positive and negative signs,
convert the string into a mathematical expression to obtain the answer.

Don't use eval or a similar built-in parser.

For example, given `-1 + (2 + 3)`, you should return 4.

---

### Interview Analysis

Darn tootin, this is a "[Hard]" problem.
It requires lexing, parsing (a grammar),
and code to evaluate the parser's output.

The problem statement is phrased in a way to make the lexing easier,
in that number representations only have a single digit.
The solution also doesn't have to account for multiplication and division.
The example does have a unary minus sign,
and that complicates the grammar.

This can't possibly be a whiteboard question.
There's too much to it - lexing, a grammar,
and evaluation of a parse tree.
The only way to "solve" this as a whiteboard question would be to talk
over the different pieces,
offering a general design rather than a specific implementation.

It could be a take home question.
It's not out of the realm of possibility to do this in 3 or 4 hours,
although the solution wouldn't be general or robust.
The candidate would almost have to use a parser generator (like yacc, or bison)
unless they already had something like the recursive descent parser
in this repo on hand.

I'm not too sure what an interviewer could expect from this question,
other than a lot of talk about design.
Implementation code would be extensive, and thus hard to completely review.
If the interviewer just looks for a design,
discussion of lexing and parsing pecularities like the unary minus would be a good sign.
Discussion of implementation alternatives (hand-written vs generated lexer and parser),
and how to do evaluation of parsed code, including error handling,
might be things to look for.

## Build

    $ cd $GOPATH/src
    $ git clone https://github.com/bediger4000/arithmetic-parser.git
    $ cd arithmetic-parser
    $ go build arithmetic-parser

`arithmetic-parser` parses and prints a single arithmetic expression,
passed to `arithmetic-parser` as a command line arguments:

    $ ./arithmetic-expressions '1 + 3*4'
    Reconstituted expression: "1 + (3 * 4)"
    /* 13 */

With a `-g` command line flag,
`arithmetic-parser` prints a [GraphViz](http://graphviz.org/) `dot` format
representation of the parse tree,
evaluates the parse tree and prints the value.
You would do something like this:

    $ ./arithmetic-expressions -g '1 + 3*4' > x.dot
    $ dot -Tpng -o x.png x.dot
    $ feh x.png


## Lexer

I did the lexer based on a [Rob Pike talk](https://www.youtube.com/watch?v=HxaD_trXwRE)

The lexer runs in its own goroutine,
using a channel to give tokens and token types to the parser.
Using a lexer struct looks like some object oriented
garbage, but under the hood, it's asynchronous with the parser.

Following along with Rob Pike and understanding his lexer design
was my motivation for this project.

## Parser

I did a recursive descent parser,
using [this Ohio State class handout](http://web.cse.ohio-state.edu/software/2231/web-sw2/extras/slides/27.Recursive-Descent-Parsing.pdf)
as a guide.

The grammar looks like this:

    expr     ->  term   {add-op term}
    term     ->  spork  {mult-op spork}
    spork    ->  factor {exp-op factor}
    factor   ->  '(' expr ')' | add-op factor | NUMBER
    add-op   ->  '+'|'-'
    mult-op  ->  '*'|'/'|'%'
    exp-op   ->  '^'

Punctuation (parentheses), operation signs and numbers
are terminal symbols.

I use "spork" as a name to get an extra level of precedence.
There's got to be an official name for this precedence.

I added "%" (for remainder/modulo),
'^' (for exponentiation)
and allowed unary positive and negative operators,
to customize the exercize.

"CFG" below abbreviates "context free grammar".
I was able to follow their rules to write the code:

* One parse method per non-terminal symbol
* A non-terminal symbol on the right-hand side of a rewrite rule leads
  to a call to the parse method for that non-terminal
* A terminal symbol on the right-hand side of a rewrite rule leads to
  "consuming" that token from the input token string
* | in the CFG leads to "if-else" in the parser
* {...}in the CFG leads to *while* in the parser

The grammar has to be correct for this sort of semi-mechanical
coding to work.
The tricky part was realizing that the parser needed to see what
type of token it had,
but not "use up" that token every time the parser looked at its type,
or the lexeme itself.
The "consuming" note in the The Ohio State handout
didn't make sense until I realized that.

## Expression evaluation

I borrowed an interface from a [2010 Google I/O talk](https://blog.golang.org/io2010)
by Rob Pike and Russ Cox
to do the actual arithmetic.

The interface looks like this:

```go
type Value interface {
    BinaryOp(op lexer.TokenType, y Value) Value
    String() string
}
```

It has an integer arithmetic implementation,
and an error holder implementation.
Package `tree` creates new `Value` instances and
calls `BinaryOp()` on them.
This simplifies `tree.Node.Eval()` immensely,
and separates arithmetic from parse tree.
Because there's an error holder implementation,
reporting run-time problems like divide-by-zero
becomes much easier.
at the cost of moving that code into package `value`.
