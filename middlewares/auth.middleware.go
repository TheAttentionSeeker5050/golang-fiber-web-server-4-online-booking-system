package middlewares

import (
	"example/web-server/controllers"
	"example/web-server/data"
	"example/web-server/utils"

	"github.com/gofiber/fiber/v2"
)

// user is authorized middleware
func UserIsAuthorized(c *fiber.Ctx) error {
	// get the provider and the token from the cookie
	provider := c.Cookies("TokenProvider", "")
	token := c.Cookies("AccessToken", "")

	c.Locals("userAuthorized", false)

	// if the provider of the token is Google
	if provider == data.AUTH_PROVIDER_GOOGLE {
		// verify the token
		claims, err := utils.VerifyGoogleToken(token)
		if err == nil {
			// add user id to local context, this will only be available in this request, and will not be passed to the next request or the template
			c.Locals("userID", claims.ID)
			c.Locals("userAuthorized", true)
		}
	} else if provider == data.AUTH_PROVIDER_LOCAL {
		// verify the token
		valid, err := utils.VerifyLocalAuthJWTToken(token)
		if err == nil && valid {
			claims, err := utils.GetLocalAuthJWTTokenClaims(token)
			if err == nil {
				// add user id to local context, this will only be available in this request, and will not be passed to the next request or the template
				c.Locals("userID", claims.ID)
				c.Locals("userAuthorized", true)
			}
		}
	}

	if !data.IsPublicRoute(c.Path()) && c.Locals("userAuthorized") == false {
		// deny access
		return controllers.UnauthorizedErrorPageController(c)
	}

	// Go to next middleware:
	return c.Next()
}
