package lib

import (
	"fmt"
	"time"

	"github.com/pkg/errors"
)

var APITimeFormats = []string{
	"2006-01-02T15:04:05-0700",
	time.RFC3339,
}

func FormatTime(t time.Time) string {
	return t.Format(APITimeFormats[0])
}

func ParseTime(s string) (time.Time, error) {
	for _, tf := range APITimeFormats {
		t, err := time.Parse(tf, s)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, errors.Errorf("Unable to ParseTime %q to any known format.", s)
}

func FormatDuration(d time.Duration) string {
	if d == 0 {
		return ""
	}
	h := int(d.Hours())
	m := int(d.Minutes()) - h*60
	s := int(d.Seconds()) - m*60 - h*60*60
	ms := int(d.Milliseconds()) - s*1000 - m*60*1000 - h*60*60*1000
	switch {
	case h > 0:
		return fmt.Sprintf("%dh%dm%ds", h, m, s)
	case m > 0:
		return fmt.Sprintf("%dm%ds", m, s)
	case s > 0:
		return fmt.Sprintf("%ds", s)
	default:
		return fmt.Sprintf("%dms", ms)
	}
}
