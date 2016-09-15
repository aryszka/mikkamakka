package mikkamakka

import (
	"bytes"
	"io"
	"os"
)

const (
	incompatibleResource = "incompatible resource"
	discardedResource    = "discarded resource"
)

type file struct {
	sys      interface{}
	buf      []byte
	err      error
	original *Val
}

var Eof = &Val{merror, "EOF"}

func isSys(a *Val) *Val {
	return is(a, sys)
}

func fopen(a *Val) *Val {
	checkType(a, mstring)

	f, err := os.Open(stringVal(a))
	if err != nil {
		return &Val{merror, err}
	}

	return &Val{sys, &file{sys: f}}
}

func stdin() *Val {
	return &Val{sys, &file{sys: os.Stdin}}
}

func stdout() *Val {
	return &Val{sys, &file{sys: os.Stdout}}
}

func stderr() *Val {
	return &Val{sys, &file{sys: os.Stderr}}
}

func fstate(f *Val) *Val {
	checkType(f, sys)

	ft, ok := f.value.(*file)
	if !ok {
		panic(incompatibleResource)
	}

	if ft.sys == nil {
		panic(discardedResource)
	}

	if ft.err != nil && ft.err != io.EOF {
		return &Val{merror, ft.err}
	}

	if ft.err == io.EOF && len(ft.buf) == 0 {
		return Eof
	}

	return fromBytes(ft.buf)
}

func fread(f *Val, n *Val) *Val {
	checkType(f, sys)
	checkType(n, number)

	ft, ok := f.value.(*file)
	if !ok {
		panic(incompatibleResource)
	}

	if ft.sys == nil {
		panic(discardedResource)
	}

	sv := ft.sys
	r, ok := sv.(io.Reader)
	if !ok {
		panic(incompatibleResource)
	}

	ft.sys = nil

	ft.buf = make([]byte, intVal(n))
	rn, err := r.Read(ft.buf)
	ft.buf = ft.buf[:rn]

	return &Val{sys, &file{sys: sv, buf: ft.buf, err: err, original: f}}
}

func fwrite(f *Val, s *Val) *Val {
	checkType(f, sys)
	checkType(s, mstring)

	ft, ok := f.value.(*file)
	if !ok {
		panic(incompatibleResource)
	}

	if ft.sys == nil {
		panic(discardedResource)
	}

	sv := ft.sys
	w, ok := sv.(io.Writer)
	if !ok {
		panic(incompatibleResource)
	}

	ft.sys = nil

	_, err := w.Write(byteVal(s))
	return &Val{sys, &file{sys: sv, err: err, original: f}}
}

func fclose(f *Val) *Val {
	checkType(f, sys)

	ft, ok := f.value.(*file)
	if !ok {
		panic(incompatibleResource)
	}

	if ft.sys == nil {
		panic(discardedResource)
	}

	sv := ft.sys
	c, ok := sv.(io.Closer)
	var err error
	if ok {
		err = c.Close()
	}

	return &Val{sys, &file{sys: sv, err: err, original: f}}
}

func sstring(f *Val) *Val {
	checkType(f, sys)
	return fromString("<file>")
}

func buffer() *Val {
	return &Val{sys, &file{sys: bytes.NewBuffer(nil)}}
}

func derivedObject(o, from *Val) *Val {
	checkType(o, sys)
	checkType(from, sys)

	if o.value.(*file).original == from {
		return True
	}

	if from.value.(*file).original == nil {
		return False
	}

	return derivedObject(o, from.value.(*file).original)
}

func argv() *Val {
	argv := Nil
	for i := len(os.Args) - 1; i >= 0; i-- {
		argv = Cons(fromString(os.Args[i]), argv)
	}

	return argv
}

func Stdin() *Val {
	return stdin()
}

func Stdout() *Val {
	return stdout()
}
