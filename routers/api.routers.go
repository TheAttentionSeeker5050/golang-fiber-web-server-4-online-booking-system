package routers

import "github.com/gofiber/fiber/v2"

// APIRoutes defines the API endpoints
func APIRoutes(group fiber.Router) {
	// /foo endpoint
	group.Get("/foo", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "bar",
		})
	})

	// /hello endpoint
	group.Get("/hello", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"message": "world",
		})
	})
}
