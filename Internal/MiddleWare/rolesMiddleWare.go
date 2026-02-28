package middleware

import (
	// "log"
	config "marryo/Config"
	utils "marryo/Internal/Utils"

	// "strconv"

	"strings"

	"github.com/gofiber/fiber/v2"
)

func RolesMiddleWare(Roles ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		tokenStr := c.Cookies("access")
		if tokenStr == "" {
			return c.Status(401).JSON(fiber.Map{
				"error": "missing token",
			})
		}

		exists, err := config.Redis.Exists(config.Ctx, "blacklist:"+tokenStr).Result()
		if err != nil || exists > 0 {
			return c.Status(401).JSON(fiber.Map{
				"error": "token is blacklisted",
			})
		}

		claims, err := utils.VerifyToken(tokenStr)
		if err != nil {
			return c.SendStatus(401)
		}

		// if claims.Role != consts.Admin {
		// 	return c.Status(401).JSON(fiber.Map{
		// 		"error": "access denied",
		// 	})
		// }

		authorized := false
		for _, role := range Roles {
			if strings.EqualFold(claims.Role, role) {
				authorized = true
				break
			}
		}
		if !authorized {
			return c.Status(403).JSON(fiber.Map{
				"error": "access denied: insufficient permissions or role mismatch",
			})
		}

		c.Locals("userID", claims.UserID)

		return c.Next()
	}
}
