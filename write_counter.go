package lib

import (
	"io"
	"sync/atomic"
)

// WriterCounter is counter for io.Writer
type WriterCounter struct {
	io.Writer
	count int64
}

// NewWriterCounter function create new WriterCounter
func NewWriterCounter(w io.Writer) *WriterCounter {
	return &WriterCounter{
		Writer: w,
	}
}

func (counter *WriterCounter) Write(buf []byte) (int, error) {
	n, err := counter.Writer.Write(buf)
	atomic.AddInt64(&counter.count, int64(n))
	return n, err
}

// Count function return counted bytes
func (counter *WriterCounter) Count() int64 {
	return atomic.LoadInt64(&counter.count)
}
