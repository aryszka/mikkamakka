package mikkamakka

import (
	"fmt"
	"os"
)

func isError(a *val) *val {
	return is(a, merror)
}

func bisError(a []*val) *val {
	return isError(a[0])
}

func stringToError(a []*val) *val {
	checkType(a[0], mstring)
	return &val{merror, stringVal(a[0])}
}

func errorStringRaw(a *val) string {
	switch v := a.value.(type) {
	case error:
		return v.Error()
	case string:
		return v
	default:
		return "unknown error"
	}
}

func errorString(a *val) *val {
	return fromString(errorStringRaw(a))
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
	checkType(e, merror)
	return fromString(fmt.Sprintf("<error:%s>", errorStringRaw(e)))
}

func IsError(a *Val) *Val {
	return (*Val)(isError((*val)(a)))
}

func Fatal(a *Val) *Val {
	return (*Val)(fatal((*val)(a)))
}
