package mikkamakka

import (
	"bytes"
	"io"
	"os"
)

const (
	IncompatibleResource = "incompatible resource"
	DiscardedResource    = "discarded resource"
)

type file struct {
	sys interface{}
	buf []byte
	err *Val
}

var Eof = ErrorFromRawString("EOF")

func IsSys(a *Val) *Val {
	return is(a, sys)
}

func Fopen(a *Val) *Val {
	f, err := os.Open(RawString(a))
	if err != nil {
		return ErrorFromSysError(err)
	}

	return &Val{sys, &file{sys: f}}
}

func Stdin() *Val {
	return &Val{sys, &file{sys: os.Stdin}}
}

func Stdout() *Val {
	return &Val{sys, &file{sys: os.Stdout}}
}

func Stderr() *Val {
	return &Val{sys, &file{sys: os.Stderr}}
}

func Fstate(f *Val) *Val {
	checkType(f, sys)

	ft, ok := f.value.(*file)
	if !ok {
		panic(IncompatibleResource)
	}

	if ft.sys == nil {
		panic(DiscardedResource)
	}

	if ft.err != nil {
		return ft.err
	}

	return StringFromRaw(string(ft.buf))
}

func sysErr(err error) *Val {
	var v *Val
	if err == io.EOF {
		v = Eof
	} else if err != nil {
		v = ErrorFromSysError(err)
	}

	return v
}

func Fread(f *Val, n *Val) *Val {
	checkType(f, sys)

	ft, ok := f.value.(*file)
	if !ok {
		panic(IncompatibleResource)
	}

	if ft.sys == nil {
		panic(DiscardedResource)
	}

	sv := ft.sys
	r, ok := sv.(io.Reader)
	if !ok {
		panic(IncompatibleResource)
	}

	ft.sys = nil

	buf := make([]byte, RawInt(n))
	rn, err := r.Read(buf)
	buf = buf[:rn]
	serr := sysErr(err)

	return &Val{sys, &file{sys: sv, buf: buf, err: serr}}
}

func Fwrite(f *Val, s *Val) *Val {
	checkType(f, sys)

	ft, ok := f.value.(*file)
	if !ok {
		panic(IncompatibleResource)
	}

	if ft.sys == nil {
		panic(DiscardedResource)
	}

	sv := ft.sys
	w, ok := sv.(io.Writer)
	if !ok {
		panic(IncompatibleResource)
	}

	ft.sys = nil

	_, err := w.Write(RawBytes(s))
	serr := sysErr(err)
	return &Val{sys, &file{sys: sv, err: serr}}
}

func Fclose(f *Val) *Val {
	checkType(f, sys)

	ft, ok := f.value.(*file)
	if !ok {
		panic(IncompatibleResource)
	}

	if ft.sys == nil {
		panic(DiscardedResource)
	}

	sv := ft.sys
	c, ok := sv.(io.Closer)
	var serr *Val
	if ok {
		err := c.Close()
		serr = sysErr(err)
	}

	return &Val{sys, &file{sys: sv, err: serr}}
}

func SysToString(f *Val) *Val {
	checkType(f, sys)
	return StringFromRaw("<file>")
}

func Buffer() *Val {
	return &Val{sys, &file{sys: bytes.NewBuffer(nil)}}
}

func Argv() *Val {
	argv := Nil
	for i := len(os.Args) - 1; i >= 0; i-- {
		argv = Cons(StringFromRaw(os.Args[i]), argv)
	}

	return argv
}
