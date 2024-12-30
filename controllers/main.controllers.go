package controllers

import (
	"example/web-server/utils"

	"github.com/gofiber/fiber/v2"
)

func StaticPageController(c *fiber.Ctx, template string, argumentsMap *fiber.Map) error {
	if argumentsMap == nil {
		argumentsMap = &fiber.Map{}
	}

	// add the Authenticated flag to the argumentsMap using http headers
	(*argumentsMap)["IsAuthenticated"] = string(c.Request().Header.Peek("Authenticated")) == "true"

	return utils.CustomRenderTemplate(c, template, *argumentsMap)
}
