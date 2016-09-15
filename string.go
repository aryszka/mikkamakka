package mikkamakka

import "strconv"

type str struct {
	sys string
}

func fromBytes(b []byte) *Val {
	return &Val{mstring, &str{string(b)}}
}

func fromString(s string) *Val {
	return &Val{mstring, &str{s}}
}

func stringVal(s *Val) string {
	return s.value.(*str).sys
}

func byteVal(s *Val) []byte {
	return []byte(s.value.(*str).sys)
}

func appendString(a ...*Val) *Val {
	var b []byte
	for _, ai := range a {
		checkType(ai, mstring)
		b = append(b, byteVal(ai)...)
	}

	return fromBytes(b)
}

func stringLength(s *Val) *Val {
	checkType(s, mstring)
	return NumberFromRawInt(len(s.value.(*str).sys))
}

func isString(a *Val) *Val {
	return is(a, mstring)
}

func seq(left, right *Val) *Val {
	if stringVal(left) == stringVal(right) {
		return True
	}

	return False
}

func escapeCompiled(a *Val) *Val {
	checkType(a, mstring)
	return fromString(strconv.Quote(stringVal(a)))
}

func FromString(s string) *Val {
	return (*Val)(fromString(s))
}
