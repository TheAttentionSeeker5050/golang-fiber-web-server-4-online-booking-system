package routers

import (
	"example/web-server/controllers"
	"example/web-server/utils"

	"github.com/gofiber/fiber/v2"
)

func ProfileRoutes(group fiber.Router) {
	// Profile page route
	group.Get("/", controllers.ViewProfileDataController)
	group.Get("/edit", controllers.EditProfileController)
	group.Post("/edit", func(c *fiber.Ctx) error {
		return c.SendString("Profile updated")
	})
	group.Get("/edit-password", func(c *fiber.Ctx) error {
		return utils.CustomRenderTemplate(c, "profile/edit-password", fiber.Map{
			"Title": "Edit Password",
		})
	})
	group.Post("/edit-password", func(c *fiber.Ctx) error {
		return c.SendString("Password updated")
	})
}
