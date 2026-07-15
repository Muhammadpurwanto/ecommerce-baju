package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"

	"github.com/Muhammadpurwanto/ecommerce-baju/services/gateway/config"
)

func Protected(cfg *config.Config, rdb *redis.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Missing or invalid Authorization header",
			})
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")

		// Check if token is blacklisted in Redis
		if rdb != nil {
			ctx := c.Context()
			key := fmt.Sprintf("blacklist:%s", tokenString)
			val, err := rdb.Get(ctx, key).Result()
			if err == nil && val == "revoked" {
				return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
					"success": false,
					"message": "Token has been revoked (logged out)",
				})
			}
		}

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWTSecret), nil
		})

		if err != nil || !token.Valid {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid or expired token",
			})
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"success": false,
				"message": "Invalid token claims",
			})
		}

		userID := fmt.Sprintf("%v", claims["user_id"])
		role := fmt.Sprintf("%v", claims["role"])

		// Inject header ke request yang akan di-proxy
		c.Request().Header.Set("X-User-ID", userID)
		c.Request().Header.Set("X-User-Role", role)

		return c.Next()
	}
}

// Opsional: Untuk rute public, kita hapus header internal agar tidak dipalsukan user dari luar
func StripInternalHeaders() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Request().Header.Del("X-User-ID")
		c.Request().Header.Del("X-User-Role")
		return c.Next()
	}
}
