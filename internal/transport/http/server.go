package httptransport

import (
	matrixapp "interseguro/go-api/internal/application/matrix"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type Config struct {
	JWTSecret string
}

func NewServer(config Config, matrixService matrixapp.Service) *fiber.App {

	app := fiber.New(fiber.Config{
		AppName: "interseguro-go-api",
	})

	app.Use(cors.New())

	registerRoutes(app, config, matrixService)

	return app
}
