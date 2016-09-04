package main

import (
	"fmt"
	"os"
)

func mpanic(a *val) *val {
	panic(fmt.Errorf("%v", a.value))
	return a
}

func isError(a *val) *val {
	return is(a, merror)
}

func errorString(a *val) *val {
	switch v := a.value.(type) {
	case error:
		return fromString(v.Error())
	case string:
		return fromString(v)
	default:
		return mpanic(a)
	}
}

func fatal(a *val) {
	switch a.mtype {
	case mstring:
		fwrite(stderr(), a)
	case merror:
		fwrite(stderr(), errorString(a))
	}
	os.Exit(-1)
}

func estring(e *val) *val {
	return fromString("<error>")
}
