package routers

import (
	"example/web-server/controllers"
	"example/web-server/utils"

	"github.com/gofiber/fiber/v2"
)

// APIRoutes defines the API endpoints
func AuthRoutes(group fiber.Router) {

	// /login endpoint
	group.Post("/login", controllers.AuthLoginController)

	// /register endpoint
	group.Post("/register", controllers.AuthRegisterController)

	// /logout POST endpoint
	group.Post("/logout", controllers.AuthLogoutController)

	// Login page route
	group.Get("/login", func(c *fiber.Ctx) error {
		argumentsMap := &fiber.Map{
			"Title": "Login",
		}

		return utils.CustomRenderTemplate(c, "auth/login", *argumentsMap)
	})

	// Register page route
	group.Get("/register", func(c *fiber.Ctx) error {
		argumentsMap := &fiber.Map{
			"Title": "Register",
		}

		return utils.CustomRenderTemplate(c, "auth/register", *argumentsMap)
	})

	// Logout route (can include logout logic later)
	group.Get("/logout", func(c *fiber.Ctx) error {
		argumentsMap := &fiber.Map{
			"Title": "Logout",
		}

		return utils.CustomRenderTemplate(c, "auth/logout", *argumentsMap)
	})

	// Logout route (can include logout logic later)
	group.Get("/logout-success", func(c *fiber.Ctx) error {
		argumentsMap := &fiber.Map{
			"Title": "Logout Success",
		}

		return utils.CustomRenderTemplate(c, "auth/logout-success", *argumentsMap)
	})
}
