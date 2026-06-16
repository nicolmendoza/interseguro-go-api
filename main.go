package main

import (
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
)

func main() {
	app := NewServer(Config{
		JWTSecret: getEnv("JWT_SECRET", "interseguro-secret"),
		NodeURL:   getEnv("NODE_API_URL", "http://localhost:3001"),
	})

	port := getEnv("PORT", "3000")
	log.Fatal(app.Listen(":" + port))
}

func getEnv(key string, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func writeError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"error": message})
}
