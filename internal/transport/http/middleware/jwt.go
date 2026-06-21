package middleware

import (
	"fmt"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func JWT(secret string) fiber.Handler {

	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")

		if !strings.HasPrefix(header, "Bearer ") {

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "falta el token bearer"})
		}

		tokenText := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenText, func(token *jwt.Token) (any, error) {

			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {

				return nil, fmt.Errorf("metodo de firma inesperado")
			}

			return []byte(secret), nil
		})
		if err != nil || !token.Valid {

			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token invalido"})
		}

		return c.Next()
	}
}
