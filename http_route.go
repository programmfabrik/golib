package lib

import (
	"net/http"
	"regexp"

	"github.com/gorilla/mux"
)

// Subroute returns the subroute form a Gorilla Mux
// Subrouter
func Subroute(req *http.Request) string {
	return req.URL.Path[len(Route(req)):]
}

// Route returns the route from a Gorilla Mux
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
	match := reg.FindString(req.URL.Path)
	// logrus.Debugf("Route: %s %v", routeRegex, match)
	return req.URL.Path[:len(match)]
}
