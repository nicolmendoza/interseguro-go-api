package handlers

import (
	"interseguro/go-api/internal/infrastructure/docs"

	"github.com/gofiber/fiber/v2"
)

func OpenAPI(c *fiber.Ctx) error {
	return c.JSON(docs.CreateOpenAPISpec())
}

func Swagger(c *fiber.Ctx) error {
	c.Type("html", "utf-8")
	return c.SendString(`<!doctype html>
<html lang="es">
  <head>
    <meta charset="utf-8">
    <title>Interseguro Go API - Swagger</title>
    <link rel="stylesheet" href="https://unpkg.com/swagger-ui-dist@5/swagger-ui.css">
  </head>
  <body>
    <div id="swagger-ui"></div>
    <script src="https://unpkg.com/swagger-ui-dist@5/swagger-ui-bundle.js"></script>
    <script>
      window.onload = function () {
        window.ui = SwaggerUIBundle({
          url: "/openapi.json",
          dom_id: "#swagger-ui"
        });
      };
    </script>
  </body>
</html>`)
}
