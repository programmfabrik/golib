package golib

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

// FileMimeType returns the content-type of a file
func FileMimeType(fn string) (mt string, err error) {
	f, err := os.Open(fn)
	if err != nil {
		return "", err
	}
	defer f.Close()
	data := make([]byte, 512)
	_, err = f.Read(data)
	if err != nil {
		return "", err
	}
	mt = http.DetectContentType(data)
	return mt, nil
}

// ReadGzipFile opens file fn and returns its unpacked bytes.
func ReadGzipFile(fn string) (bs []byte, err error) {
	buf := bytes.Buffer{}
	f, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r, err := gzip.NewReader(f)
	if err != nil {
		return nil, err
	}
	_, err = io.Copy(&buf, r)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

// FilepathAbs returns the file's absolute path, relative
// to the startDir
func FilepathAbs(startDir, path string) (pathAbs string) {
	if filepath.IsAbs(path) {
		return path
	}
	return filepath.Join(startDir, path)
}
