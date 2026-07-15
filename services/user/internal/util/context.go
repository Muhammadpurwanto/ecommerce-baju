package util

import (
	"github.com/gofiber/fiber/v2"
)

// GetUserIDFromContext extracts user ID from Fiber context (set by API Gateway via header)
func GetUserIDFromContext(c *fiber.Ctx) (string, error) {
	userIDStr := c.Get("X-User-ID")
	if userIDStr == "" {
		return "", fiber.ErrUnauthorized
	}

	return userIDStr, nil
}

// GetUserRoleFromContext extracts user role from Fiber context
func GetUserRoleFromContext(c *fiber.Ctx) string {
	return c.Get("X-User-Role", "customer")
}

// IsAdmin checks if the current user is an admin
func IsAdmin(c *fiber.Ctx) bool {
	return GetUserRoleFromContext(c) == "admin"
}
