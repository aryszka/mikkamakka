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

func appendString(left, right *val) *val {
	checkType(left, mstring)
	checkType(right, mstring)
	return &val{mstring, &str{left.value.(*str).sys + right.value.(*str).sys}}
}

func stringLength(s *val) *val {
	checkType(s, mstring)
	return fromInt(len(s.value.(*str).sys))
}
