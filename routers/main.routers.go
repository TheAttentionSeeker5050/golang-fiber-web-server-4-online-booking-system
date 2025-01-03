package routers

import (
	"example/web-server/controllers"
	"example/web-server/utils"

	"github.com/gofiber/fiber/v2"
)

// import the fiber package

// HomeRoutes defines the routes for the main application
func HomeRoutes(group fiber.Router) {

	// use favicon
	group.Static("/favicon.svg", "./public/favicon.svg")
	// Home page route
	group.Get("/", func(c *fiber.Ctx) error {

		return controllers.StaticPageController(c, "index", &fiber.Map{
			"Title": "Home",
		})
	})

	// About page route
	group.Get("/about", func(c *fiber.Ctx) error {
		return controllers.StaticPageController(c, "pages/about", &fiber.Map{
			"Title": "About",
		})
	})

	// Contact page route
	group.Get("/contact", func(c *fiber.Ctx) error {
		return controllers.StaticPageController(c, "pages/contact", &fiber.Map{
			"Title": "Contact",
		})
	})

	// // Profile page route
	// group.Get("/profile", controllers.ViewProfileDataController)

	// Dashboard page route
	group.Get("/dashboard", func(c *fiber.Ctx) error {
		argumentsMap := &fiber.Map{
			"Title": "Dashboard",
		}

		// add the Authenticated flag to the argumentsMap using http headers
		(*argumentsMap)["IsAuthenticated"] = string(c.Request().Header.Peek("Authenticated")) == "true"

		return utils.CustomRenderTemplate(c, "dashboard", *argumentsMap)
	})
}
