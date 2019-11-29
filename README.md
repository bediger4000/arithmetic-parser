# Arithmetic Expressions, lexer and parser

Another cut at a Golang algebraic-order-of-operations arithmetic
expression evaluator.

Compare to [my earlier incarnation](https://github.com/bediger4000/arithmetic-expressions)
which was more-or-less a Go transliteration of a C program.

This version does share the implementation of the parse tree
with my earlier incarnation.

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
was my motivation for this little project.

## Parser

I did a recursive descent parser,
using [this Ohio State class handout](http://web.cse.ohio-state.edu/software/2231/web-sw2/extras/slides/27.Recursive-Descent-Parsing.pdf)
as a guide.

The grammar looks like this:

    expr     ->  term   {add-op term}
    term     ->  spork  {mult-op spork}
    spork    ->  factor {exp-op factor}
    factor   ->  '(' expr ')' | '-' factor | NUMBER
    add-op   ->  '+'|'-'
    mult-op  ->  '*'|'/'|'%'
    exp-op   ->  '^'

I added "%" (for remainder/modulo),
'^' (for exponentiation)
and allowed negative numbers,
to customize the exercize.

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
