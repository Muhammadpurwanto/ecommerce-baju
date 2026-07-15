package util

import (
	"errors"

	"github.com/gofiber/fiber/v2"
)

func GetUserIDFromContext(c *fiber.Ctx) (string, error) {
	// Ambil dari header yang di-inject oleh middleware auth
	userID := c.Get("X-User-ID")
	if userID == "" {
		return "", errors.New("user ID not found in context")
	}
	return userID, nil
}

func IsAdmin(c *fiber.Ctx) bool {
	role := c.Get("X-User-Role")
	return role == "admin"
}
