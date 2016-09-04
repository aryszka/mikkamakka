package main

import (
	"io"
	"os"
)

const (
	incompatibleResource = "incompatible resource"
	discardedResource    = "discarded resource"
)

type file struct {
	sys *os.File
	buf []byte
	err error
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

func stdin() *val {
	return &val{sys, &file{sys: os.Stdin}}
}

func stdout() *val {
	return &val{sys, &file{sys: os.Stdout}}
}

func stderr() *val {
	return &val{sys, &file{sys: os.Stderr}}
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

	var sysFile *os.File
	sysFile, ft.sys = ft.sys, nil

	ft.buf = make([]byte, intVal(n))
	rn, err := sysFile.Read(ft.buf)
	ft.buf = ft.buf[:rn]

	return &val{sys, &file{sys: sysFile, buf: ft.buf, err: err}}
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

	var sysFile *os.File
	sysFile, ft.sys = ft.sys, nil

	_, err := sysFile.Write(byteVal(s))
	return &val{sys, &file{sys: sysFile, err: err}}
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

	var sysFile *os.File
	sysFile, ft.sys = ft.sys, nil

	return &val{sys, &file{sys: sysFile, err: sysFile.Close()}}
}

func sstring(f *val) *val {
	checkType(f, sys)
	return fromString("<file>")
}
