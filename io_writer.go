package golib

import (
	"io"
	"log"
)

type IOWriter func([]byte) (int, error)

func (iow IOWriter) Write(bts []byte) (n int, err error) {
	return iow(bts)
}

func IOWriterFromLogger(l *log.Logger) io.Writer {
	return IOWriter(func(bts []byte) (int, error) {
		l.Print(string(bts))
		return len(bts), nil
	})
}
