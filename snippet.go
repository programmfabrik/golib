package lib

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
	from := lineNo - plusMinus
	to := lineNo + plusMinus

	fromLen := len(strconv.Itoa(from))
	toLen := len(strconv.Itoa(to))

	pad := fromLen
	if toLen > fromLen {
		pad = toLen
	}

	for scanner.Scan() {
		line++
		if line > to {
			break
		}
		if line >= from && line <= to {
			gt := " "
			if line == lineNo {
				gt = "➜"
			}
			snipped = append(snipped, []byte(fmt.Sprintf("%s%"+strconv.Itoa(pad+1)+"d: %s\n", gt, line, scanner.Text()))...)
		}
	}
	return snipped
}
