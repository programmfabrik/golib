package golib

import (
	"net/url"
	"slices"
	"strings"
)

// List of extensions a browser can display
var ExtDisplayableInBrowser = []string{"jpg", "jpeg", "png", "gif", "webp", "svg", "bmp"}

func IsDisplayableInBrowser(e string) bool {
	return slices.Contains(ExtDisplayableInBrowser, strings.ToLower(e))
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
