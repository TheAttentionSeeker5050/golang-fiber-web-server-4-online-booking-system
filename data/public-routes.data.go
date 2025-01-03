package data

import "regexp"

// list of public routes
var publicRoutes = []string{
	// all the routes starting with /auth and /auth/ should be public
	"/auth/*",
	// these should be exact match: /, /contact, /about
	"^/$",
	"^/contact$",
	"^/about$",
	// public routes
	"/public/*",
}

// function verify if the route is public
func IsPublicRoute(route string) bool {
	for _, publicRoute := range publicRoutes {
		// make the condition be regex using regex comparison
		regexpRoute := regexp.MustCompile(publicRoute)
		if regexpRoute.MatchString(route) {
			return true
		}
	}
	return false
}
