package handlers

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func Token(jwtSecret string) fiber.Handler {
	return func(c *fiber.Ctx) error {

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": "technical-challenge-client",
			"exp": time.Now().Add(2 * time.Hour).Unix(),
		})

		signed, err := token.SignedString([]byte(jwtSecret))
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "no se pudo firmar el token"})
		}

		return c.JSON(fiber.Map{"token": signed})
	}
}
