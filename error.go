package mikkamakka

import (
	"fmt"
	"os"
)

func IsError(a *Val) *Val {
	return is(a, Error)
}

func SysStringToError(s string) *Val {
	return newVal(Error, s)
}

func SysErrorToError(err error) *Val {
	return newVal(Error, err)
}

func StringToError(a *Val) *Val {
	return SysStringToError(StringToSysString(a))
}

func ErrorToSysString(a *Val) string {
	checkType(a, Error)
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
	return SysStringToString(fmt.Sprintf("<error:%s>", ErrorToSysString(a)))
}

func Fatal(a *Val) *Val {
	// panic(errorString(a))

	switch a.typ {
	case String:
		Fwrite(Stderr(), a)
	case Error:
		Fwrite(Stderr(), ErrorToString(a))
	}

	println()
	os.Exit(-1)
	return a
}
