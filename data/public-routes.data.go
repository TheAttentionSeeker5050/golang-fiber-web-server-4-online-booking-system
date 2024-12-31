package data

// list of public routes
var publicRoutes = []string{
	"/auth/login",
	"/auth/register",
	"/auth/logout",
	"/",
	"/auth/google",
	"/auth/google/callback",
	"/contact",
	"about",
}

// function verify if the route is public
func IsPublicRoute(route string) bool {
	for _, publicRoute := range publicRoutes {
		if publicRoute == route {
			return true
		}
	}
	return false
}
