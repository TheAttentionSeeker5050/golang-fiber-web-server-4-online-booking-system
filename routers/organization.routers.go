package routers

import (
	"example/web-server/controllers"

	"github.com/gofiber/fiber/v2"
)

// Organization defines the API endpoints
func OrganizationRoutes(group fiber.Router) {
	// CRUD for organizations
	group.Get("/", controllers.ListViewOrganizationsController)

	group.Get("/add", controllers.AddOrganizationController)

	group.Post("/add", controllers.AddOrganizationPostController)

	group.Get("/view/:id", controllers.DetailViewOrganizationController)

	group.Get("/edit/:id", controllers.EditOrganizationController)

	group.Post("/edit/:id", controllers.EditOrganizationPostController)

	group.Get("/delete/:id", controllers.DeleteOrganizationController)

	group.Post("/delete/:id", controllers.DeleteOrganizationPostController)
}
