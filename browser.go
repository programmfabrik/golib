package golib

import (
	"net/url"
	"strings"
)

func IsDisplayableInBrowser(e string) bool {
	ext := strings.ToLower(e)
	for _, extOk := range []string{"jpg", "jpeg", "png", "gif", "webp", "svg", "bmp"} {
		if extOk == ext {
			return true
		}
	}
	return false
}

func ContentDisposition(disposition string, filename string) (key, value string) {
	if disposition == "" {
		disposition = "inline"
	}
	if filename != "" {
		// We use path escape here, so that " " is not changed into "+"
		return "Content-Disposition", disposition + "; filename=\"" + strings.ReplaceAll(filename, `"`, `_`) +
			"\"; filename*=UTF-8''" + url.PathEscape(filename)
	} else {
		return "Content-Disposition", disposition
	}
}
