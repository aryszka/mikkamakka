package mikkamakka

var (
	False = &Val{mbool, false}
	True  = &Val{mbool, true}
)

func bfromString(s string) *Val {
	switch s {
	case "true":
		return True
	case "false":
		return False
	default:
		return invalidToken
	}
}

func tryBoolFromString(s *Val) *Val {
	checkType(s, mstring)
	return bfromString(stringVal(s))
}

func btryBoolFromString(a []*Val) *Val {
	return tryBoolFromString(a[0])
}

func boolToString(b *Val) *Val {
	if b == True {
		return fromString("true")
	}

	return fromString("false")
}

func bboolToString(a []*Val) *Val {
	return boolToString(a[0])
}

func isBool(a *Val) *Val {
	if a.mtype == mbool {
		return True
	}

	return False
}

func bisBool(a []*Val) *Val {
	return isBool(a[0])
}

func and(v ...*Val) *Val {
	if len(v) == 0 {
		return True
	}

	if len(v) == 1 || v[0] == False {
		return v[0]
	}

	return and(v[1:]...)
}

func band(v []*Val) *Val {
	return and(v...)
}

func or(v ...*Val) *Val {
	if len(v) == 0 {
		return False
	}

	if len(v) == 1 || v[0] != False {
		return v[0]
	}

	return or(v[1:]...)
}

func bor(v []*Val) *Val {
	return or(v...)
}

func not(a []*Val) *Val {
	checkType(a[0], mbool)
	if a[0] == False {
		return True
	}

	return False
}
