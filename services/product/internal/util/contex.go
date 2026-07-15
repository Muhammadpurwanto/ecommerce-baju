package util

import (
	"github.com/gofiber/fiber/v2"
)

func GetUserRoleFromContext(c *fiber.Ctx) string {
	return c.Get("X-User-Role", "customer") // Dari Gateway
}

func IsAdmin(c *fiber.Ctx) bool {
	return GetUserRoleFromContext(c) == "admin"
}
