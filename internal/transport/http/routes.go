package httptransport

import (
	matrixapp "interseguro/go-api/internal/application/matrix"
	"interseguro/go-api/internal/transport/http/handlers"
	"interseguro/go-api/internal/transport/http/middleware"

	"github.com/gofiber/fiber/v2"
)

func registerRoutes(app *fiber.App, config Config, matrixService matrixapp.Service) {

	app.Get("/health", handlers.Health)
	app.Post("/auth/token", handlers.Token(config.JWTSecret))

	protected := app.Group("/", middleware.JWT(config.JWTSecret))
	protected.Post("/qr", handlers.QR(matrixService))
	protected.Post("/rotate", handlers.Rotate(matrixService))
	protected.Post("/analyze", handlers.Analyze(matrixService))
}
