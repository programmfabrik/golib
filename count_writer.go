package golib

import (
	"errors"
	"time"
)

type CountWriter struct {
	// This one needs to be a poiner, otherwise it gets reset in every Write call
	bytes int64
	limit int64
	// lastPrint int64
	start time.Time
}

func NewCountWriter() *CountWriter {
	return &CountWriter{}
}

func NewLimitCountWriter(limit int64) *CountWriter {
	return &CountWriter{start: time.Now(), limit: limit}
}

var LimitExceeded = errors.New("Limit exceeded")

func (cbw *CountWriter) Write(p []byte) (n int, err error) {
	n = len(p)
	cbw.bytes += int64(n)
	if cbw.limit > 0 && cbw.bytes > cbw.limit {
		return n, LimitExceeded
	}
	// if cbw.lastPrint == 0 || cbw.lastPrint+1_000_000 < cbw.bytes {
	// 	Pln("writing %s  bytes now %s %s/sec.", HumanByteSize(uint64(cbw.bytes-cbw.lastPrint)), HumanByteSize(uint64(cbw.bytes)),
	// 		HumanByteSize(uint64(float64(cbw.bytes)/time.Since(cbw.start).Seconds())))
	// 	cbw.lastPrint = cbw.bytes
	// }
	return n, nil
}

func (cbw CountWriter) Count() int64 {
	return cbw.bytes
}
