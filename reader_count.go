package golib

import (
	"io"
)

// It simply forwards the Read() call, while displaying
// the results from individual calls to it.
type ReadCounter struct {
	io.Reader
	total  uint64 // Total # of bytes transferred
	prefix string
}

func NewReadCounter(prefix string, reader io.Reader) *ReadCounter {
	rc := &ReadCounter{Reader: reader}
	rc.prefix = prefix
	return rc
}

// Read 'overrides' the underlying io.Reader's Read method.
// This is the one that will be called by io.Copy(). We simply
// use it to keep track of byte counts and then forward the call.
func (pt *ReadCounter) Read(p []byte) (int, error) {
	n, err := pt.Reader.Read(p)
	pt.total += uint64(n)
	return n, err
}
