package golib

import (
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

// Subroute returns the subroute form a Gorilla Mux Subrouter. It uses the
// escaped form of the path, so a path like /files/he%2Fnk/more will keep the
// %2F.
func Subroute(req *http.Request) string {
	// Pln("route: %s", Route(req))
	// Pln("rout2: %s", url.PathEscape(Route(req)))
	// Pln("esc  : %s", req.URL.EscapedPath())
	return req.URL.EscapedPath()[len(Route(req)):]
}

// Route returns the route from a Gorilla Mux. It uses the escaped
// form of the path.
func Route(req *http.Request) string {
	route := mux.CurrentRoute(req)
	if route == nil {
		return req.URL.Path
	}
	routeRegex, err := route.GetPathRegexp()
	if err != nil {
		panic(err)
	}
	reg := regexp.MustCompile(routeRegex)
	escPath := req.URL.EscapedPath()
	match := reg.FindString(escPath)
	// Pln("Route: %s\nMatch: %v\ncut: %q", routeRegex, match, escPath[:len(match)])
	return escPath[:len(match)]
}
