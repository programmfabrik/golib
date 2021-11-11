package golib

import (
	"bytes"
	"io"

	"github.com/gabriel-vasile/mimetype"
)

// MimeTypeReader returns the MIME type of input and a new reader
// containing the whole data from input.
func MimeTypeReader(input io.Reader) (mimeType string, recycled io.Reader, err error) {
	// header will store the bytes mimetype uses for detection.
	header := bytes.NewBuffer(nil)

	// After DetectReader, the data read from input is copied into header.
	mtype, err := mimetype.DetectReader(io.TeeReader(input, header))

	// Concatenate back the header to the rest of the file.
	// recycled now contains the complete, original data.
	recycled = io.MultiReader(header, input)

	return mtype.String(), recycled, err
}
