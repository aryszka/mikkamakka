package main

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
	original *val
}

var eof = &val{merror, "EOF"}

func isSys(a *val) *val {
	return is(a, sys)
}

func fopen(a *val) *val {
	checkType(a, mstring)

	f, err := os.Open(a.value.(string))
	if err != nil {
		return &val{merror, err}
	}

	return &val{sys, &file{sys: f}}
}

func bfopen(a []*val) *val {
	return fopen(a[0])
}

func stdin() *val {
	return &val{sys, &file{sys: os.Stdin}}
}

func bstdin([]*val) *val {
	return stdin()
}

func stdout() *val {
	return &val{sys, &file{sys: os.Stdout}}
}

func bstdout([]*val) *val {
	return stdout()
}

func stderr() *val {
	return &val{sys, &file{sys: os.Stderr}}
}

func bstderr([]*val) *val {
	return stderr()
}

func fstate(f *val) *val {
	checkType(f, sys)

	ft, ok := f.value.(*file)
	if !ok {
		panic(incompatibleResource)
	}

	if ft.sys == nil {
		panic(discardedResource)
	}

	if ft.err != nil && ft.err != io.EOF {
		return &val{merror, ft.err}
	}

	if ft.err == io.EOF && len(ft.buf) == 0 {
		return eof
	}

	return fromBytes(ft.buf)
}

func bfstate(a []*val) *val {
	return fstate(a[0])
}

func fread(f *val, n *val) *val {
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

	return &val{sys, &file{sys: sv, buf: ft.buf, err: err, original: f}}
}

func bfread(a []*val) *val {
	return fread(a[0], a[1])
}

func fwrite(f *val, s *val) *val {
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
	return &val{sys, &file{sys: sv, err: err, original: f}}
}

func bfwrite(a []*val) *val {
	return fwrite(a[0], a[1])
}

func fclose(f *val) *val {
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

	return &val{sys, &file{sys: sv, err: err, original: f}}
}

func bfclose(a []*val) *val {
	return fclose(a[0])
}

func sstring(f *val) *val {
	checkType(f, sys)
	return fromString("<file>")
}

func bsstring(a []*val) *val {
	return sstring(a[0])
}

func buffer() *val {
	return &val{sys, &file{sys: bytes.NewBuffer(nil)}}
}

func bbuffer([]*val) *val {
	return buffer()
}

func derivedObject(a []*val) *val {
	checkType(a[0], sys)
	checkType(a[1], sys)

	if a[0].value.(*file).original == a[1] {
		return vtrue
	}

	if a[1].value.(*file).original == nil {
		return vfalse
	}

	return derivedObject([]*val{a[0], a[1].value.(*file).original})
}
