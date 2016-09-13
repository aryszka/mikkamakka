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

func bfopen(a []*Val) *Val {
	return fopen(a[0])
}

func stdin() *Val {
	return &Val{sys, &file{sys: os.Stdin}}
}

func bstdin([]*Val) *Val {
	return stdin()
}

func stdout() *Val {
	return &Val{sys, &file{sys: os.Stdout}}
}

func bstdout([]*Val) *Val {
	return stdout()
}

func stderr() *Val {
	return &Val{sys, &file{sys: os.Stderr}}
}

func bstderr([]*Val) *Val {
	return stderr()
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

func bfstate(a []*Val) *Val {
	return fstate(a[0])
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

func bfread(a []*Val) *Val {
	return fread(a[0], a[1])
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

func bfwrite(a []*Val) *Val {
	return fwrite(a[0], a[1])
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

func bfclose(a []*Val) *Val {
	return fclose(a[0])
}

func sstring(f *Val) *Val {
	checkType(f, sys)
	return fromString("<file>")
}

func bsstring(a []*Val) *Val {
	return sstring(a[0])
}

func buffer() *Val {
	return &Val{sys, &file{sys: bytes.NewBuffer(nil)}}
}

func bbuffer([]*Val) *Val {
	return buffer()
}

func derivedObject(a []*Val) *Val {
	checkType(a[0], sys)
	checkType(a[1], sys)

	if a[0].value.(*file).original == a[1] {
		return True
	}

	if a[1].value.(*file).original == nil {
		return False
	}

	return derivedObject([]*Val{a[0], a[1].value.(*file).original})
}

func argv([]*Val) *Val {
	argv := Nil
	for i := len(os.Args) - 1; i >= 0; i-- {
		argv = cons(fromString(os.Args[i]), argv)
	}

	return argv
}

func Stdin() *Val {
	return stdin()
}

func Stdout() *Val {
	return stdout()
}
