package mikkamakka

import (
	"fmt"
	"os"
)

func isError(a *Val) *Val {
	return is(a, merror)
}

func stringToError(a *Val) *Val {
	checkType(a, mstring)
	return &Val{merror, stringVal(a)}
}

func errorStringRaw(a *Val) string {
	switch v := a.value.(type) {
	case error:
		return v.Error()
	case string:
		return v
	default:
		return "unknown error"
	}
}

func errorString(a *Val) *Val {
	return fromString(errorStringRaw(a))
}

func fatal(a *Val) *Val {
	// panic(errorString(a))

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

func estring(e *Val) *Val {
	checkType(e, merror)
	return fromString(fmt.Sprintf("<error:%s>", errorStringRaw(e)))
}

func IsError(a *Val) *Val {
	return isError(a)
}

func Fatal(a *Val) *Val {
	return fatal(a)
}
