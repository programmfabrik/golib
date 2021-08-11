package lib

import (
	"bufio"
	"bytes"
	"fmt"
)

// Snippet returns lines from-to, prefixed with the line number
func Snippet(bts []byte, from, to int) (snipped []byte) {
	snipped = []byte{}
	scanner := bufio.NewScanner(bytes.NewReader(bts))
	line := 0
	for scanner.Scan() {
		line++
		if line > to {
			break
		}
		if line >= from && line <= to {
			snipped = append(snipped, []byte(fmt.Sprintf("%5d: %s\n", line, scanner.Text()))...)
		}
	}
	return snipped
}
