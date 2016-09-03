package main

import (
	"fmt"
	"os"
)

func mpanic(a *val) *val {
	panic(fmt.Errorf("%v", a.value))
}

func isError(a *val) *val {
	return is(a, merror)
}

func errorString(a *val) *val {
	return fromString(a.value.(error).Error())
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
