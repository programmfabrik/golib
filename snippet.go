package golib

import (
	"bufio"
	"bytes"
	"fmt"
	"strconv"
)

// Snippet returns lines from lineNo-plusMinus to lineNo+plusMinus, prefixed with the line number
func Snippet(bts []byte, lineNo, plusMinus int) (snipped []byte) {
	snipped = []byte{}
	scanner := bufio.NewScanner(bytes.NewReader(bts))
	line := 0
	var from, to, pad int
	if lineNo > -1 && plusMinus > -1 {
		from = lineNo - plusMinus
		to = lineNo + plusMinus
		fromLen := len(strconv.Itoa(from))
		toLen := len(strconv.Itoa(to))
		pad = fromLen
		if toLen > fromLen {
			pad = toLen
		}
	} else {
		pad = 5
		from = -1
		to = -1
	}

	for scanner.Scan() {
		line++
		if to > -1 && line > to {
			break
		}
		if from == -1 || line >= from && line <= to {
			gt := " "
			if line == lineNo {
				gt = "➜"
			}
			snipped = append(snipped, []byte(fmt.Sprintf("%s%"+strconv.Itoa(pad+1)+"d: %s\n", gt, line, scanner.Text()))...)
		}
	}
	return snipped
}
