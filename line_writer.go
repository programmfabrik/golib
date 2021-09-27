package golib

import (
	"bytes"
	"encoding/json"
	"io"
)

// JsonLineWriter implements io.Writer and send
// the data to an Encoder whenever a newline is written
// After each newline a "," is send to the underlying
// Writer
type JsonLineWriter struct {
	enc *json.Encoder
	w   io.Writer
	buf bytes.Buffer
}

func (lw JsonLineWriter) Flush() (err error) {
	if lw.buf.Len() == 0 {
		return nil
	}
	err = lw.enc.Encode(lw.buf.String())
	if err != nil {
		return err
	}
	lw.buf.Reset()
	return nil
}

func (lw JsonLineWriter) Write(data []byte) (n int, err error) {
	for _, b := range data {
		switch b {
		case '\n':
			// flush buffer as json
			lw.Flush()
			_, err = lw.w.Write([]byte{','})
			if err != nil {
				return n, err
			}
		default:
			_, err := lw.buf.Write([]byte{b})
			if err != nil {
				return n, err
			}
		}
		n++
	}
	return n, nil
}

// NewLineWriter encodes each newline as JSON followed by a ","
func NewJsonLineWriter(w io.Writer, enc *json.Encoder) *JsonLineWriter {
	return &JsonLineWriter{
		w:   w,
		enc: enc,
		buf: bytes.Buffer{},
	}
}
