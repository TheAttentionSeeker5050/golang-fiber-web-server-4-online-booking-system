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
	group.Post("/edit", controllers.EditProfilePostController)

	group.Get("/edit-password", func(c *fiber.Ctx) error {
		return utils.CustomRenderTemplate(c, "profile/edit-password", fiber.Map{
			"Title": "Edit Password",
		})
	})
	group.Post("/edit-password", controllers.EditPasswordPostController)
}
