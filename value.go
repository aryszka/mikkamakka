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

func unexpectedType(expected ...mtype) {
	s := make([]string, len(expected))
	for i, e := range expected {
		s[i] = typeString(e)
	}

	panic(fmt.Sprintf("expected: %s", strings.Join(s, ", ")))
}

func checkType(v *val, expected ...mtype) {
	for _, t := range expected {
		if v.mtype == t {
			return
		}
	}

	unexpectedType(expected...)
}
