package value

import (
	"arithmetic-parser/lexer"
	"fmt"
	"strconv"
)

// Value interface allows parse tree evaluation to return
// an error all the way up the call stack to the user.
// There's an integer type and an error type that fit this interface.
type Value interface {
	BinaryOp(op lexer.TokenType, y Value) Value
	String() string
}

// NewValue creates an instance of type Int if possible, which fits Value
// interface. Otherwise it creates an Error instance, which will end up the
// result of an evaluation.
func NewValue(lit string) Value {
	x, err := strconv.Atoi(lit)
	if err == nil {
		return Int(x)
	}
	return Error(fmt.Sprintf("illegal literal '%s'", lit))
}

// Int implements Value interface for integer arithmetic.
type Int int

func (x Int) String() string { return strconv.Itoa(int(x)) }

// BinaryOp implements integer arithmetic for type Int
func (x Int) BinaryOp(op lexer.TokenType, y Value) Value {
	switch y := y.(type) {
	case Int:
		switch op {
		case lexer.PLUS:
			return x + y
		case lexer.MINUS:
			return x - y
		case lexer.MULT:
			return x * y
		case lexer.DIV:
			if y == 0 {
				return Error(fmt.Sprintf("division by zero: '%v / %v'", x, y))
			}
			return x / y
		// case lexer.EXP: // not implemented, falls through to illegal op
		case lexer.REM:
			if y == 0 {
				return Error(fmt.Sprintf("modulo of zero: '%v %% %v'", x, y))
			}
			return x % y
		}
	case Error:
		return y
	}
	return Error(fmt.Sprintf("illegal op: '%v %s %v'", x, op, y))
}

// Error implements Value interface for pushing errors up the stack
type Error string

func (e Error) String() string {
	return string(e)
}

// BinaryOp makes Xerror instances into fit interface Value
func (e Error) BinaryOp(op lexer.TokenType, y Value) Value {
	return e
}
