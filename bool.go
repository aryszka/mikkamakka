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
