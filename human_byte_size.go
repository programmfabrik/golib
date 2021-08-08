package lib

import "fmt"

func HumanByteSize(b uint64) string {
	//Source https://yourbasic.org/golang/formatting-byte-size-to-human-readable-format/
	const unit = 1000
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	return fmt.Sprintf("%.1f %cB",
		float64(b)/float64(div), "kMGTPE"[exp])
}
