package middlewares

import (
	"example/web-server/data"
	"example/web-server/models"
	"math/rand"

	"github.com/gofiber/fiber/v2"
)

// user is authorized middleware
func UserIsAuthorized(c *fiber.Ctx) error {
	// get the provider and the token from the cookie
	provider := c.Cookies("TokenProvider", "")
	token := c.Cookies("AccessToken", "")

	// if the provider of the token is Google
	if provider == "Google" {
		// verify the token
		claims, err := models.VerifyGoogleToken(token)
		if err != nil {

			// if the token is invalid, check if the route is public
			if !data.IsPublicRoute(c.Path()) {
				// deny access
				return c.SendStatus(fiber.StatusUnauthorized)
			}
		}

		// add user id to local context, this will only be available in this request, and will not be passed to the next request or the template
		c.Locals("userID", claims.ID)
		c.Locals("userAuthorized", true)

	} else {
		// Set a local variable for the request context
		c.Locals("userAuthorized", rand.Intn(2) == 1)
	}

	// Go to next middleware:
	return c.Next()
}
