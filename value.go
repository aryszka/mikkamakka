package main

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
	procedure
)

type val struct {
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
	case procedure:
		return "procedure"
	default:
		panic("invalid type")
	}
}

func is(v *val, t mtype) *val {
	if v.mtype == t {
		return vtrue
	}

	return vfalse
}

func unexpectedType(got mtype, expected ...mtype) {
	s := make([]string, len(expected))
	for i, e := range expected {
		s[i] = typeString(e)
	}

	panic(fmt.Sprintf("expected: %s, got: %s", strings.Join(s, ", "), typeString(got)))
}

func checkType(v *val, expected ...mtype) {
	for _, t := range expected {
		if v.mtype == t {
			return
		}
	}

	unexpectedType(v.mtype, expected...)
}

func eq(v ...*val) *val {
	if len(v) == 0 {
		return vfalse
	}

	if len(v) == 1 {
		return vtrue
	}

	a := v[0]
	b := v[1]

	switch {
	case isNumber(a) != vfalse && isNumber(b) != vfalse:
		return and(neq(a, b), eq(v[1:]...))
	case isString(a) != vfalse && isString(b) != vfalse:
		return and(seq(a, b), eq(v[1:]...))
	case isSymbol(a) != vfalse && isSymbol(b) != vfalse:
		return and(smeq(a, b), eq(v[1:]...))
	default:
		if a != b {
			return vfalse
		}

		return eq(v[1:]...)
	}
}

func beq(v []*val) *val {
	return eq(v...)
}
