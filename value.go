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
	mbool
	pair
	mnil
	vector
	mstruct
	function
	sys
	merror // true or false? turn into false
	environment
)

type Val struct {
	mtype mtype
	value interface{}
}

type (
	typeCheck func(*Val) *Val
	typeEq    func(*Val, *Val) *Val
)

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

func unexpectedType(got mtype, v *Val, expected ...mtype) *Val {
	s := make([]string, len(expected))
	for i, e := range expected {
		s[i] = typeString(e)
	}

	msg := fmt.Sprintf(
		"expected: %s, got: %s, with value: %v",
		strings.Join(s, ", "),
		typeString(got), v.value)
	return fatal(fromString(msg))
}

func checkType(v *Val, expected ...mtype) *Val {
	for _, t := range expected {
		if v.mtype == t {
			return True
		}
	}

	return unexpectedType(v.mtype, v, expected...)
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

	if isString(v[0]) != False {
		return eqT(v, isString, seq)
	}

	if isSymbol(v[0]) != False {
		return eqT(v, isSymbol, smeq)
	}

	if v[0] == v[1] {
		return True
	}

	return False
}
