package main

type str struct {
	sys string
}

func fromBytes(b []byte) *val {
	return &val{mstring, &str{string(b)}}
}

func fromString(s string) *val {
	return &val{mstring, &str{s}}
}

func stringVal(s *val) string {
	return s.value.(*str).sys
}

func byteVal(s *val) []byte {
	return []byte(s.value.(*str).sys)
}

func appendString(a ...*val) *val {
	var b []byte
	for _, ai := range a {
		checkType(ai, mstring)
		b = append(b, byteVal(ai)...)
	}

	return fromBytes(b)
}

func stringLength(s *val) *val {
	checkType(s, mstring)
	return fromInt(len(s.value.(*str).sys))
}

func isString(a *val) *val {
	return is(a, mstring)
}
