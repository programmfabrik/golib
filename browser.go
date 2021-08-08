package lib

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
		return "Content-Disposition", disposition + "; filename*=UTF-8''" + url.QueryEscape(filename)
	} else {
		return "Content-Disposition", disposition
	}
}
