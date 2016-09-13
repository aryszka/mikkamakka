package mikkamakka

import (
	"fmt"
	"strings"
)

type mtype int

const (
	notype mtype = iota
	symbol
	number
	mstring
	mbool // sure that needed? better needed
	pair
	mnil
	vector
	mstruct
	sys
	merror // true or false?
	environment
	function
)

type Val struct {
	mtype mtype
	value interface{}
}

func typeString(t mtype) string {
	switch t {
	case symbol:
		return "symbol"
	case number:
		return "number"
	case mstring:
		return "string"
	case mbool:
		return "bool"
	case pair:
		return "pair"
	case mnil:
		return "nil"
	case vector:
		return "vector"
	case mstruct:
		return "struct"
	case sys:
		return "sys"
	case merror:
		return "error"
	case environment:
		return "environment"
	case function:
		return "function"
	default:
		panic("invalid type")
	}
}

func is(v *Val, t mtype) *Val {
	if v.mtype == t {
		return True
	}

	return False
}

func unexpectedType(got mtype, expected ...mtype) {
	s := make([]string, len(expected))
	for i, e := range expected {
		s[i] = typeString(e)
	}

	panic(fmt.Sprintf("expected: %s, got: %s", strings.Join(s, ", "), typeString(got)))
}

func checkType(v *Val, expected ...mtype) {
	for _, t := range expected {
		if v.mtype == t {
			return
		}
	}

	unexpectedType(v.mtype, expected...)
}

func eq(v ...*Val) *Val {
	if len(v) == 0 {
		return False
	}

	if len(v) == 1 {
		return True
	}

	a := v[0]
	b := v[1]

	switch {
	case isNumber(a) != False && isNumber(b) != False:
		return and(neq(a, b), eq(v[1:]...))
	case isString(a) != False && isString(b) != False:
		return and(seq(a, b), eq(v[1:]...))
	case isSymbol(a) != False && isSymbol(b) != False:
		return and(smeq(a, b), eq(v[1:]...))
	default:
		if a != b {
			return False
		}

		return eq(v[1:]...)
	}
}

func beq(v []*Val) *Val {
	return eq(v...)
}
