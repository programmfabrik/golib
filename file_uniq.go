package golib

import (
	"fmt"
	"path/filepath"
	"strings"
)

type UniqueFilename map[string]bool

func NewUniqueFilename() (uf *UniqueFilename) {
	return &UniqueFilename{}
}

func (uf *UniqueFilename) Add(fn string) string {
	if (*uf)[strings.ToLower(fn)] {
		// fn already exists, find a new path
		i := 1
		for {
			ext := filepath.Ext(fn)
			newPath := fn[:len(fn)-len(ext)] + "_" + fmt.Sprintf("%05d", i) + ext
			if !(*uf)[strings.ToLower(newPath)] {
				// found a new path
				fn = newPath
				break
			}
			i++
		}
	}
	(*uf)[strings.ToLower(fn)] = true
	return fn
}
