package middleware

import "github.com/gofiber/fiber/v2"



func RequirePermission(permission string) fiber.Handler {

	return func(c *fiber.Ctx) error {

		userPermissions := c.Locals("permissions").([]string)

		for _, p := range userPermissions {
			if p == permission {
				return c.Next()
			}
		}

		return c.Status(403).JSON(fiber.Map{
			"error": "Access denied",
		})
	}
}