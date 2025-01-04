package routers

import (
	"example/web-server/controllers"

	"github.com/gofiber/fiber/v2"
)

func DashboardRoutes(group fiber.Router) {
	group.Get("/", controllers.ViewDashboardController)
}
