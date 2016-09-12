package mikkamakka

var (
	vfalse = &val{mbool, false}
	vtrue  = &val{mbool, true}
)

func bfromString(s string) *val {
	switch s {
	case "true":
		return vtrue
	case "false":
		return vfalse
	default:
		return invalidToken
	}
}

func tryBoolFromString(s *val) *val {
	checkType(s, mstring)
	return bfromString(stringVal(s))
}

func btryBoolFromString(a []*val) *val {
	return tryBoolFromString(a[0])
}

func boolToString(b *val) *val {
	if b == vtrue {
		return fromString("true")
	}

	return fromString("false")
}

func bboolToString(a []*val) *val {
	return boolToString(a[0])
}

func isBool(a *val) *val {
	if a.mtype == mbool {
		return vtrue
	}

	return vfalse
}

func bisBool(a []*val) *val {
	return isBool(a[0])
}

func and(v ...*val) *val {
	if len(v) == 0 {
		return vtrue
	}

	if len(v) == 1 || v[0] == vfalse {
		return v[0]
	}

	return and(v[1:]...)
}

func band(v []*val) *val {
	return and(v...)
}

func or(v ...*val) *val {
	if len(v) == 0 {
		return vfalse
	}

	if len(v) == 1 || v[0] != vfalse {
		return v[0]
	}

	return or(v[1:]...)
}

func bor(v []*val) *val {
	return or(v...)
}

func not(a []*val) *val {
	checkType(a[0], mbool)
	if a[0] == vfalse {
		return vtrue
	}

	return vfalse
}

var Vfalse = (*Val)(vfalse)
var Vtrue = (*Val)(vtrue)

func BfromString(s string) *Val {
	return (*Val)(bfromString(s))
}
