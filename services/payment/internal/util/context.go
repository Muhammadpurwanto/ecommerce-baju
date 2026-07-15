package util

import (
	"github.com/gofiber/fiber/v2"
)

func GetUserIDFromContext(c *fiber.Ctx) (string, error) {
	userIDStr := c.Get("X-User-ID")
	if userIDStr == "" {
		return "", fiber.ErrUnauthorized
	}

	return userIDStr, nil
}
