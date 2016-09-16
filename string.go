package mikkamakka

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

func StringFromRaw(s string) *Val {
	return &Val{mstring, s}
}

func RawString(s *Val) string {
	checkType(s, mstring)
	return s.value.(string)
}

func RawBytes(s *Val) []byte {
	return []byte(RawString(s))
}

func StringLen(s *Val) *Val {
	checkType(s, mstring)
	return NumberFromRawInt(utf8.RuneCount(RawBytes(s)))
}

func ByteLen(s *Val) *Val {
	checkType(s, mstring)
	return NumberFromRawInt(len(RawString(s)))
}

func AppendString(a ...*Val) *Val {
	var s []string
	for _, ai := range a {
		s = append(s, RawString(ai))
	}

	return StringFromRaw(strings.Join(s, ""))
}

func IsString(a *Val) *Val {
	return is(a, mstring)
}

func stringEq(left, right *Val) *Val {
	if RawString(left) == RawString(right) {
		return True
	}

	return False
}

func EscapeCompiled(a *Val) *Val {
	checkType(a, mstring)
	return StringFromRaw(strconv.Quote(RawString(a)))
}
