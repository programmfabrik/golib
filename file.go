package golib

import (
	"bytes"
	"compress/gzip"
	"fmt"
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

// AbsPathExecutable prefixes p with the path to the currently running
// executable if p is not absolute.
func AbsPathExecutable(p string) (string, error) {
	if filepath.IsAbs(p) {
		return p, nil
	}
	// prefix the dir of the executable
	execFile, err := os.Executable()
	if err != nil {
		return "", fmt.Errorf("Unable to get path of executable: %w", err)
	}
	p, err = filepath.Abs(filepath.Join(filepath.Dir(execFile), p))
	if err != nil {
		return "", fmt.Errorf("Unable to get path of executable: %w", err)
	}
	return p, nil
}
