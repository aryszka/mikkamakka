package mikkamakka

import (
	"fmt"
	"os"
)

func IsError(a *Val) *Val {
	return is(a, merror)
}

func ErrorFromRawString(s string) *Val {
	return &Val{merror, s}
}

func ErrorFromSysError(err error) *Val {
	return &Val{merror, err}
}

func StringToError(a *Val) *Val {
	return ErrorFromRawString(RawString(a))
}

func RawErrorString(a *Val) string {
	checkType(a, merror)
	switch v := a.value.(type) {
	case error:
		return v.Error()
	case string:
		return v
	default:
		return "unknown error"
	}
}

func ErrorToString(a *Val) *Val {
	return StringFromRaw(fmt.Sprintf("<error:%s>", RawErrorString(a)))
}

func Fatal(a *Val) *Val {
	// panic(errorString(a))

	switch a.mtype {
	case mstring:
		Fwrite(Stderr(), a)
	case merror:
		Fwrite(Stderr(), ErrorToString(a))
	}

	println()
	os.Exit(-1)
	return a
}
