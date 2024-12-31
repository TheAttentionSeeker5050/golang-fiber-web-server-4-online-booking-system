package utils

import "github.com/gofiber/fiber/v2"

// RenderTemplate merges global context (e.g., userAuthorized) with page-specific data
func CustomRenderTemplate(c *fiber.Ctx, templateName string, pageContext fiber.Map) error {
	// Get the userAuthorized value from Locals
	userAuthorized, ok := c.Locals("userAuthorized").(bool)
	if !ok {
		userAuthorized = false // Default to false
	}

	// Merge userAuthorized into the context
	if pageContext == nil {
		pageContext = fiber.Map{}
	}
	pageContext["userAuthorized"] = userAuthorized

	// Render the template
	return c.Render(templateName, pageContext)
}

// the path to the layout file
// errorMsg is optional
func GetFiberRenderMappingsAuthForms(email string, password string, errorMsgs *[]string, success bool) fiber.Map {
	var resultFiberMap fiber.Map = fiber.Map{
		"Errors":   []string{},
		"Email":    email,
		"Password": password,
		"Success":  success,
	}

	// if the size of the errorMsgs is 0
	if errorMsgs == nil {
		if !success {
			resultFiberMap["Errors"] = []string{"Internal server error"}
		}
	} else {
		resultFiberMap["Errors"] = *errorMsgs
	}

	return resultFiberMap
}
