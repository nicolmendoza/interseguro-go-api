package httptransport

import (
	"fmt"
	"strings"
	"time"

	"interseguro/go-api/internal/application"
	"interseguro/go-api/internal/domain"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	JWTSecret string
}

type MatrixRequest struct {
	Matrix domain.Matrix `json:"matrix"`
}

func NewServer(config Config, matrixService application.MatrixService) *fiber.App {
	app := fiber.New(fiber.Config{
		AppName: "interseguro-go-api",
	})
	app.Use(cors.New())

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"status": "ok", "service": "go-api"})
	})

	app.Post("/auth/token", func(c *fiber.Ctx) error {
		token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
			"sub": "technical-challenge-client",
			"exp": time.Now().Add(2 * time.Hour).Unix(),
		})

		signed, err := token.SignedString([]byte(config.JWTSecret))
		if err != nil {
			return writeError(c, fiber.StatusInternalServerError, "could not sign token")
		}

		return c.JSON(fiber.Map{"token": signed})
	})

	protected := app.Group("/", jwtMiddleware(config.JWTSecret))
	protected.Post("/qr", handleQR(matrixService))
	protected.Post("/rotate", handleRotate(matrixService))
	protected.Post("/analyze", handleAnalyze(matrixService))

	return app
}

func jwtMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		header := c.Get("Authorization")
		if !strings.HasPrefix(header, "Bearer ") {
			return writeError(c, fiber.StatusUnauthorized, "missing bearer token")
		}

		tokenText := strings.TrimPrefix(header, "Bearer ")
		token, err := jwt.Parse(tokenText, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !token.Valid {
			return writeError(c, fiber.StatusUnauthorized, "invalid token")
		}

		return c.Next()
	}
}

func handleQR(matrixService application.MatrixService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request, ok := parseMatrixRequest(c)
		if !ok {
			return nil
		}

		result, err := matrixService.Factorize(request.Matrix)
		if err != nil {
			return writeError(c, fiber.StatusBadRequest, err.Error())
		}

		return c.JSON(result)
	}
}

func handleRotate(matrixService application.MatrixService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request, ok := parseMatrixRequest(c)
		if !ok {
			return nil
		}

		rotated, err := matrixService.Rotate(request.Matrix)
		if err != nil {
			return writeError(c, fiber.StatusBadRequest, err.Error())
		}

		return c.JSON(fiber.Map{"rotated": rotated})
	}
}

func handleAnalyze(matrixService application.MatrixService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		request, ok := parseMatrixRequest(c)
		if !ok {
			return nil
		}

		result, statusCode, err := matrixService.Analyze(request.Matrix, c.Get("Authorization"))
		if err != nil {
			return writeError(c, statusCodeOrDefault(statusCode), err.Error())
		}
		if statusCode >= 400 {
			return c.Status(statusCode).Send(result.Stats)
		}

		return c.JSON(result)
	}
}

func parseMatrixRequest(c *fiber.Ctx) (MatrixRequest, bool) {
	var request MatrixRequest
	if err := c.BodyParser(&request); err != nil {
		_ = writeError(c, fiber.StatusBadRequest, "invalid JSON body")
		return MatrixRequest{}, false
	}
	return request, true
}

func writeError(c *fiber.Ctx, status int, message string) error {
	return c.Status(status).JSON(fiber.Map{"error": message})
}

func statusCodeOrDefault(statusCode int) int {
	if statusCode == 0 {
		return fiber.StatusBadRequest
	}
	return statusCode
}
