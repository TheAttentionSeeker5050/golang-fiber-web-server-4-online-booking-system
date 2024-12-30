package middlewares

import (
	"github.com/gofiber/fiber/v2"
)

// user is authorized middleware
func UserIsAuthorized(c *fiber.Ctx) error {

	// get request cookie Authenticated
	cookie := c.Cookies("Authenticated", "false")

	// Check if the cookie indicates the user is authorized
	isAuthorized := cookie == "true"

	// Set a local variable for the request context
	c.Locals("userAuthorized", isAuthorized)

	// Go to next middleware:
	return c.Next()
}
