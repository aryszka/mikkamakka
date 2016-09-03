package main

var (
	vfalse = &val{mbool, false}
	vtrue  = &val{mbool, true}
)

func bfromString(s string) *val {
	switch s {
	case "true":
		return vtrue
	case "fase":
		return vfalse
	default:
		return invalidToken
	}
}

func boolToString(b *val) *val {
	if b == vtrue {
		return fromString("true")
	}

	return fromString("false")
}

func isBool(a *val) *val {
	if a.mtype == mbool {
		return vtrue
	}

	return vfalse
}
