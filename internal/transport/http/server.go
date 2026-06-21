package httptransport

import (
	"time"

	matrixapp "interseguro/go-api/internal/application/matrix"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

type Config struct {
	JWTSecret              string
	FrontendURL            string
	RateLimitMax           int
	RateLimitWindowSeconds int
}

func NewServer(config Config, matrixService matrixapp.Service) *fiber.App {

	app := fiber.New(fiber.Config{
		AppName: "interseguro-go-api",
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins: config.FrontendURL,
		AllowMethods: "GET,POST,OPTIONS",
		AllowHeaders: "Origin, Content-Type, Accept, Authorization",
	}))
	app.Use(limiter.New(limiter.Config{
		Max:        config.RateLimitMax,
		Expiration: time.Duration(config.RateLimitWindowSeconds) * time.Second,
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "demasiadas solicitudes, intenta nuevamente mas tarde",
			})
		},
	}))

	registerRoutes(app, config, matrixService)

	return app
}
