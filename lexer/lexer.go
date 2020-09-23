package lexer

import (
	"unicode"
)

// TokenType tells parser what the lexer thinks
// the category for this token is.
type TokenType int

// EOF and others: all the types of tokens
// ADD_OP, MULT_OP, EXP_OP aren't individual operation
// tokens, but rather levels of precedence. The lexeme
// for a MULT_OP will be "*" or "/" or "%", which all have
// the same level of precedence.
const (
	EOF      TokenType = iota
	ADD_OP   TokenType = iota
	MULT_OP  TokenType = iota
	EXP_OP   TokenType = iota
	PLUS     TokenType = iota
	MINUS    TokenType = iota
	MULT     TokenType = iota
	DIV      TokenType = iota
	REM      TokenType = iota
	EXP      TokenType = iota
	CONSTANT TokenType = iota
	LPAREN   TokenType = iota
	RPAREN   TokenType = iota
	POSITIVE TokenType = iota
	NEGATIVE TokenType = iota
	EOL      TokenType = iota
)

func (t TokenType) String() string {
	switch t {
	case ADD_OP:
		return "ADD_OP"
	case MULT_OP:
		return "MULT_OP"
	case CONSTANT:
		return "CONSTANT"
	case LPAREN:
		return "LPAREN"
	case RPAREN:
		return "RPAREN"
	case EOL:
		return "EOL"
	case EOF:
		return "EOF"
	}
	return "unknown"
}

type item struct {
	kind   TokenType
	lexeme string
}

// Lexer instances hold information needed to break
// a string into arithmetic expression tokens.
type Lexer struct {
	input       []rune
	start       int
	pos         int
	items       chan item
	currentItem item
	consumed    bool
}

type stateFn func(*Lexer) stateFn

// Lex creates a new ready-to-go Lexer instance.
// Runs a goroutine in the background.
func Lex(input string) *Lexer {
	l := &Lexer{
		input:    []rune(input),
		items:    make(chan item),
		consumed: true,
	}
	go l.run()
	return l
}

// NextToken called by parser to retrieve whatever
// the lexer thinks is the next token.
func (l *Lexer) NextToken() (TokenType, string) {
	if l.consumed {
		l.currentItem = <-l.items
		l.consumed = false
	}
	return l.currentItem.kind, l.currentItem.lexeme
}

// Consume called by parser when it has found a place in the parse tree for the
// token. Parse can and does call NextToken() repeatedly to find out the
// token's type.
func (l *Lexer) Consume() {
	l.consumed = true
}

func (l *Lexer) run() {
	for state := lexWhiteSpace; state != nil; {
		state = state(l)
	}
	close(l.items)
}

func lexWhiteSpace(l *Lexer) stateFn {
	for _, r := range l.input[l.start:] {
		switch r {
		case ' ', '"', '\'', '\t':
			l.pos++
			l.start++
		default:
			return l.nextStateFn()
		}
	}
	return nil
}

func (l *Lexer) nextStateFn() stateFn {
	if l.pos >= len(l.input) {
		return lexEOF
	}
	switch l.input[l.pos] {
	case '(':
		return lexLeftParen
	case ')':
		return lexRightParen
	case '/':
		return lexSlash
	case '-':
		return lexMinus
	case '+':
		return lexPlus
	case '*':
		return lexStar
	case '^':
		return lexExp
	case '%':
		return lexMod
	case '\n':
		return lexEOL
	default:
		if unicode.IsDigit(l.input[l.pos]) {
			return lexNumber
		}
		return lexWhiteSpace
	}
}

func (l *Lexer) emit(t TokenType) {
	l.items <- item{t, string(l.input[l.start:l.pos])}
	l.start = l.pos
}

func lexEOF(l *Lexer) stateFn {
	return nil
}

func lexNumber(l *Lexer) stateFn {
	for l.pos < len(l.input) && unicode.IsDigit(rune(l.input[l.pos])) {
		l.pos++
	}
	l.emit(CONSTANT)
	return l.nextStateFn()
}

func lexLeftParen(l *Lexer) stateFn {
	l.pos++
	l.emit(LPAREN)
	return l.nextStateFn()
}

func lexRightParen(l *Lexer) stateFn {
	l.pos++
	l.emit(RPAREN)
	return l.nextStateFn()
}

func lexPlus(l *Lexer) stateFn {
	l.pos++
	l.emit(ADD_OP)
	return l.nextStateFn()
}

func lexMinus(l *Lexer) stateFn {
	l.pos++
	l.emit(ADD_OP)
	return l.nextStateFn()
}

func lexSlash(l *Lexer) stateFn {
	l.pos++
	l.emit(MULT_OP)
	return l.nextStateFn()
}

func lexExp(l *Lexer) stateFn {
	l.pos++
	l.emit(EXP_OP)
	return l.nextStateFn()
}

func lexStar(l *Lexer) stateFn {
	l.pos++
	l.emit(MULT_OP)
	return l.nextStateFn()
}

func lexMod(l *Lexer) stateFn {
	l.pos++
	l.emit(MULT_OP)
	return l.nextStateFn()
}

func lexEOL(l *Lexer) stateFn {
	l.pos++
	l.emit(EOL)
	return l.nextStateFn()
}
