package golib

import "net/http"

// Uses stdlib http.Redirect but sets cache header
func HttpRedirect(rw http.ResponseWriter, req *http.Request, url string, code int) {
	rw.Header().Set("cache-control", "no-store")
	http.Redirect(rw, req, url, code)
}
