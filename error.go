package main

import "os"

func isError(a *val) *val {
	return is(a, merror)
}

func bisError(a []*val) *val {
	return isError(a[0])
}

func errorString(a *val) *val {
	switch v := a.value.(type) {
	case error:
		return fromString(v.Error())
	case string:
		return fromString(v)
	default:
		return fromString("unknown error")
	}
}

func fatal(a *val) *val {
	switch a.mtype {
	case mstring:
		fwrite(stderr(), a)
	case merror:
		fwrite(stderr(), errorString(a))
	}

	println()
	os.Exit(-1)
	return a
}

func estring(e *val) *val {
	return fromString("<error>")
}
