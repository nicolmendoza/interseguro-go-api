package handlers

import (
	matrixapp "interseguro/go-api/internal/application/matrix"
	"interseguro/go-api/internal/transport/http/dto"

	"github.com/gofiber/fiber/v2"
)

func QR(matrixService matrixapp.Service) fiber.Handler {

	return func(c *fiber.Ctx) error {

		request, ok := parseMatrixRequest(c)

		if !ok {
			return nil
		}

		result, err := matrixService.Factorize(request.Matrix)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(result)
	}
}

func Rotate(matrixService matrixapp.Service) fiber.Handler {

	return func(c *fiber.Ctx) error {
		request, ok := parseMatrixRequest(c)

		if !ok {
			return nil
		}

		result, err := matrixService.Rotate(request.Matrix)

		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
		}

		return c.JSON(result)
	}
}

func Analyze(matrixService matrixapp.Service) fiber.Handler {

	return func(c *fiber.Ctx) error {
		request, ok := parseMatrixRequest(c)
		if !ok {
			return nil
		}

		result, statusCode, err := matrixService.Analyze(request.Matrix, c.Get("Authorization"))

		if err != nil {
			return c.Status(statusCodeOrDefault(statusCode)).JSON(fiber.Map{"error": err.Error()})
		}

		if statusCode >= 400 {
			return c.Status(statusCode).Send(result.Stats)
		}

		return c.JSON(result)
	}
}

func parseMatrixRequest(c *fiber.Ctx) (dto.MatrixRequest, bool) {
	var request dto.MatrixRequest

	if err := c.BodyParser(&request); err != nil {
		_ = c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "el cuerpo de la solicitud debe ser un JSON valido"})

		return dto.MatrixRequest{}, false
	}
	return request, true
}

func statusCodeOrDefault(statusCode int) int {
	if statusCode == 0 {
		return fiber.StatusBadRequest
	}
	return statusCode
}
