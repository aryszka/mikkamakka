package mikkamakka

import "errors"

type readerFail struct{}

var errTestReadFailed = errors.New("test read failed")

func (fr readerFail) Read([]byte) (int, error) {
	return 0, errTestReadFailed
}

func (fr readerFail) Write([]byte) (int, error) {
	return 0, errTestReadFailed
}

func (fr readerFail) Close([]byte) error {
	return errTestReadFailed
}

func FailingIO() *Val {
	return &Val{sys, &file{sys: readerFail{}}}
}
