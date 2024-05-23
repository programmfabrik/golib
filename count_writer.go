package golib

import "fmt"

type CountWriter struct {
	// This one needs to be a poiner, otherwise it gets reset in every Write call
	bytes *int64
	limit int64
}

func NewCountWriter() CountWriter {
	return CountWriter{Int64Ref(0), 0}
}

func NewLimitCountWriter(limit int64) CountWriter {
	return CountWriter{Int64Ref(0), limit}
}

var LimitExceeded = fmt.Errorf("Limit exceeded")

func (cbw CountWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	*cbw.bytes += int64(n)
	if cbw.limit > 0 && *cbw.bytes > cbw.limit {
		return n, LimitExceeded
	}
	return n, nil
}

func (cbw CountWriter) Count() int64 {
	if cbw.bytes == nil {
		return 0
	}
	return *cbw.bytes
}
