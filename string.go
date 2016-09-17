package mikkamakka

import (
	"strconv"
	"strings"
	"unicode/utf8"
)

func SysStringToString(s string) *Val {
	return newVal(String, s)
}

func StringToSysString(s *Val) string {
	checkType(s, String)
	return s.value.(string)
}

func StringToBytes(s *Val) []byte {
	return []byte(StringToSysString(s))
}

func StringLen(s *Val) *Val {
	return SysIntToNumber(utf8.RuneCount(StringToBytes(s)))
}

func ByteLen(s *Val) *Val {
	checkType(s, String)
	return SysIntToNumber(len(StringToSysString(s)))
}

func AppendString(a ...*Val) *Val {
	var s []string
	for _, ai := range a {
		s = append(s, StringToSysString(ai))
	}

	return SysStringToString(strings.Join(s, ""))
}

func IsString(a *Val) *Val {
	return is(a, String)
}

func stringEq(left, right *Val) *Val {
	if StringToSysString(left) == StringToSysString(right) {
		return True
	}

	return False
}

func EscapeCompiled(a *Val) *Val {
	return SysStringToString(strconv.Quote(StringToSysString(a)))
}
