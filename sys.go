package mikkamakka

import (
	"bytes"
	"io"
	"os"
)

var (
	IncompatibleResource = SysStringToError("incompatible resource")
	DiscardedResource    = SysStringToError("discarded resource")
)

type file struct {
	sys interface{}
	buf []byte
	err *Val
}

var Eof = SysStringToError("EOF")

func IsSys(a *Val) *Val {
	return is(a, Sys)
}

func Fopen(a *Val) *Val {
	f, err := os.Open(StringToSysString(a))
	if err != nil {
		return SysErrorToError(err)
	}

	return newVal(Sys, &file{sys: f})
}

func Stdin() *Val {
	return newVal(Sys, &file{sys: os.Stdin})
}

func Stdout() *Val {
	return newVal(Sys, &file{sys: os.Stdout})
}

func Stderr() *Val {
	return newVal(Sys, &file{sys: os.Stderr})
}

func Fstate(f *Val) *Val {
	checkType(f, Sys)

	ft, ok := f.value.(*file)
	if !ok {
		Fatal(IncompatibleResource)
	}

	if ft.sys == nil {
		Fatal(DiscardedResource)
	}

	if ft.err != nil {
		return ft.err
	}

	return SysStringToString(string(ft.buf))
}

func sysErr(err error) *Val {
	var v *Val
	if err == io.EOF {
		v = Eof
	} else if err != nil {
		v = SysErrorToError(err)
	}

	return v
}

func Fread(f *Val, n *Val) *Val {
	checkType(f, Sys)

	ft, ok := f.value.(*file)
	if !ok {
		Fatal(IncompatibleResource)
	}

	if ft.sys == nil {
		Fatal(DiscardedResource)
	}

	sv := ft.sys
	r, ok := sv.(io.Reader)
	if !ok {
		Fatal(IncompatibleResource)
	}

	ft.sys = nil

	buf := make([]byte, NumberToSysInt(n))
	rn, err := r.Read(buf)
	buf = buf[:rn]
	serr := sysErr(err)

	return newVal(Sys, &file{sys: sv, buf: buf, err: serr})
}

func Fwrite(f *Val, s *Val) *Val {
	checkType(f, Sys)

	ft, ok := f.value.(*file)
	if !ok {
		Fatal(IncompatibleResource)
	}

	if ft.sys == nil {
		Fatal(DiscardedResource)
	}

	sv := ft.sys
	w, ok := sv.(io.Writer)
	if !ok {
		Fatal(IncompatibleResource)
	}

	ft.sys = nil

	_, err := w.Write(StringToBytes(s))
	serr := sysErr(err)
	return newVal(Sys, &file{sys: sv, err: serr})
}

func Fclose(f *Val) *Val {
	checkType(f, Sys)

	ft, ok := f.value.(*file)
	if !ok {
		Fatal(IncompatibleResource)
	}

	if ft.sys == nil {
		Fatal(DiscardedResource)
	}

	sv := ft.sys
	c, ok := sv.(io.Closer)
	var serr *Val
	if ok {
		err := c.Close()
		serr = sysErr(err)
	}

	return newVal(Sys, &file{sys: sv, err: serr})
}

func SysToString(f *Val) *Val {
	checkType(f, Sys)
	return SysStringToString("<file>")
}

func Buffer() *Val {
	return newVal(Sys, &file{sys: bytes.NewBuffer(nil)})
}

func Argv() *Val {
	argv := NilVal
	for i := len(os.Args) - 1; i >= 0; i-- {
		argv = Cons(SysStringToString(os.Args[i]), argv)
	}

	return argv
}

func Exit(n *Val) *Val {
	os.Exit(NumberToSysInt(n))
	return n
}
