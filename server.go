package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	JWTSecret string
	NodeURL   string
}

type MatrixRequest struct {
	Matrix Matrix `json:"matrix"`
}

type QRResponse struct {
	Q Matrix `json:"q"`
	R Matrix `json:"r"`
}

type AnalyzeResponse struct {
	Q     Matrix          `json:"q"`
	R     Matrix          `json:"r"`
	Stats json.RawMessage `json:"stats"`
}

func NewServer(config Config) *fiber.App {
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
	protected.Post("/qr", handleQR)
	protected.Post("/rotate", handleRotate)
	protected.Post("/analyze", handleAnalyze(config))

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

func handleQR(c *fiber.Ctx) error {
	var request MatrixRequest
	if err := c.BodyParser(&request); err != nil {
		return writeError(c, fiber.StatusBadRequest, "invalid JSON body")
	}
	if err := request.Matrix.ValidateRectangular(); err != nil {
		return writeError(c, fiber.StatusBadRequest, err.Error())
	}

	q, r, err := QRFactorization(request.Matrix)
	if err != nil {
		return writeError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(QRResponse{Q: q, R: r})
}

func handleRotate(c *fiber.Ctx) error {
	var request MatrixRequest
	if err := c.BodyParser(&request); err != nil {
		return writeError(c, fiber.StatusBadRequest, "invalid JSON body")
	}
	if err := request.Matrix.ValidateRectangular(); err != nil {
		return writeError(c, fiber.StatusBadRequest, err.Error())
	}

	return c.JSON(fiber.Map{"rotated": RotateClockwise(request.Matrix)})
}

func handleAnalyze(config Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var request MatrixRequest
		if err := c.BodyParser(&request); err != nil {
			return writeError(c, fiber.StatusBadRequest, "invalid JSON body")
		}
		if err := request.Matrix.Validate(); err != nil {
			return writeError(c, fiber.StatusBadRequest, err.Error())
		}

		q, r, err := QRFactorization(request.Matrix)
		if err != nil {
			return writeError(c, fiber.StatusBadRequest, err.Error())
		}

		payload, _ := json.Marshal(fiber.Map{"matrices": []Matrix{q, r}})
		req, err := http.NewRequest(http.MethodPost, strings.TrimRight(config.NodeURL, "/")+"/stats", bytes.NewReader(payload))
		if err != nil {
			return writeError(c, fiber.StatusInternalServerError, "could not create Node API request")
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", c.Get("Authorization"))

		client := http.Client{Timeout: 5 * time.Second}
		resp, err := client.Do(req)
		if err != nil {
			return writeError(c, fiber.StatusBadGateway, "Node API is unavailable")
		}
		defer resp.Body.Close()

		var stats json.RawMessage
		if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
			return writeError(c, fiber.StatusBadGateway, "invalid Node API response")
		}
		if resp.StatusCode >= 400 {
			return c.Status(resp.StatusCode).Send(stats)
		}

		return c.JSON(AnalyzeResponse{Q: q, R: r, Stats: stats})
	}
}
