package mikkamakka

import (
	"fmt"
	"strings"
)

type Type int

const (
	Notype Type = iota
	Symbol
	Number
	String
	Bool
	Pair
	Nil
	Vector
	Struct
	Function
	Sys
	Error // true or false? turn into false
	Environment
)

type Val struct {
	typ   Type
	value interface{}
}

type (
	typeCheck func(*Val) *Val
	typeEq    func(*Val, *Val) *Val
)

func newVal(t Type, v interface{}) *Val {
	return &Val{t, v}
}

func typeString(t Type) string {
	switch t {
	case Symbol:
		return "symbol"
	case Number:
		return "number"
	case String:
		return "string"
	case Bool:
		return "bool"
	case Pair:
		return "pair"
	case Nil:
		return "nil"
	case Vector:
		return "vector"
	case Struct:
		return "struct"
	case Sys:
		return "sys"
	case Error:
		return "error"
	case Environment:
		return "environment"
	case Function:
		return "function"
	default:
		panic("invalid type")
	}
}

func is(v *Val, t Type) *Val {
	if v.typ == t {
		return True
	}

	return False
}

func unexpectedType(got Type, v *Val, expected ...Type) *Val {
	s := make([]string, len(expected))
	for i, e := range expected {
		s[i] = typeString(e)
	}

	msg := fmt.Sprintf(
		"expected: %s, got: %s, with value: %v",
		strings.Join(s, ", "),
		typeString(got), v.value)
	return Fatal(SysStringToString(msg))
}

func checkType(v *Val, expected ...Type) *Val {
	for _, t := range expected {
		if v.typ == t {
			return True
		}
	}

	return unexpectedType(v.typ, v, expected...)
}

func eqT(v []*Val, tc typeCheck, teq typeEq) *Val {
	if len(v) == 1 {
		return True
	}

	if tc(v[1]) == False {
		return False
	}

	if teq(v[0], v[1]) == False || eqT(v[1:], tc, teq) == False {
		return False
	}

	return True
}

func Eq(v ...*Val) *Val {
	if len(v) == 0 {
		return False
	}

	if IsNumber(v[0]) != False {
		return eqT(v, IsNumber, numberEq)
	}

	if IsString(v[0]) != False {
		return eqT(v, IsString, stringEq)
	}

	if IsSymbol(v[0]) != False {
		return eqT(v, IsSymbol, symbolEq)
	}

	if v[0] == v[1] {
		return True
	}

	return False
}
