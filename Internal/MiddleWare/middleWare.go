package middleware

import (
	config "marryo/Config"
	utils "marryo/Internal/Utils"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func MiddleWare() fiber.Handler {
	return func(c *fiber.Ctx) error {

		tokenStr := c.Cookies("access")
		if tokenStr == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "missing token",
			})
		}

		exists, err := config.Redis.Exists(config.Ctx, "blacklist:"+tokenStr).Result()
		if err != nil || exists == 1 {
			return c.Status(401).JSON(fiber.Map{
				"error": "blacklist token",
			})
		}

		token, err := utils.Parse(tokenStr)
		if err != nil || !token.Valid {
			return c.Status(401).JSON(fiber.Map{
				"error": "invalid token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.SendStatus(401)
		}

		idFloat, ok := claims["userID"].(float64)
		if !ok {
			return c.SendStatus(401)
		}

		userID := uint(idFloat)

		c.Locals("userID", userID)

		return c.Next()
	}
}
