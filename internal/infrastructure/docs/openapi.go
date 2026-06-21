package docs

func CreateOpenAPISpec() map[string]any {
	return map[string]any{
		"openapi": "3.0.3",
		"info": map[string]any{
			"title":       "Interseguro Go API",
			"version":     "1.0.0",
			"description": "Documentacion Swagger para la API Go/Fiber de factorizacion QR, rotacion y analisis de matrices.",
		},
		"servers": []map[string]any{
			{
				"url":         "/",
				"description": "Go API",
			},
		},
		"tags": []map[string]any{
			{"name": "Health", "description": "Estado del servicio"},
			{"name": "Auth", "description": "Generacion de tokens JWT"},
			{"name": "Matrix", "description": "Factorizacion QR, rotacion y analisis"},
		},
		"components": map[string]any{
			"securitySchemes": map[string]any{
				"bearerAuth": map[string]any{
					"type":         "http",
					"scheme":       "bearer",
					"bearerFormat": "JWT",
				},
			},
			"schemas": map[string]any{
				"Matrix": map[string]any{
					"type": "array",
					"items": map[string]any{
						"type": "array",
						"items": map[string]any{
							"type": "number",
						},
					},
					"example": [][]float64{
						{12, -51, 4},
						{6, 167, -68},
						{-4, 24, -41},
					},
				},
				"MatrixRequest": map[string]any{
					"type":     "object",
					"required": []string{"matrix"},
					"properties": map[string]any{
						"matrix": map[string]any{"$ref": "#/components/schemas/Matrix"},
					},
				},
				"QRResult": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"q": map[string]any{"$ref": "#/components/schemas/Matrix"},
						"r": map[string]any{"$ref": "#/components/schemas/Matrix"},
					},
				},
				"RotationResult": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"rotated": map[string]any{"$ref": "#/components/schemas/Matrix"},
					},
				},
				"AnalyzeResult": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"q":     map[string]any{"$ref": "#/components/schemas/Matrix"},
						"r":     map[string]any{"$ref": "#/components/schemas/Matrix"},
						"stats": map[string]any{"type": "object"},
					},
				},
				"TokenResponse": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"token": map[string]any{"type": "string"},
					},
				},
				"ErrorResponse": map[string]any{
					"type": "object",
					"properties": map[string]any{
						"error": map[string]any{"type": "string"},
					},
				},
			},
		},
		"paths": map[string]any{
			"/health": map[string]any{
				"get": map[string]any{
					"tags":    []string{"Health"},
					"summary": "Consultar estado del servicio",
					"responses": map[string]any{
						"200": map[string]any{"description": "Servicio disponible"},
					},
				},
			},
			"/auth/token": map[string]any{
				"post": map[string]any{
					"tags":    []string{"Auth"},
					"summary": "Generar un token JWT",
					"responses": map[string]any{
						"200": map[string]any{
							"description": "JWT generado correctamente",
							"content": map[string]any{
								"application/json": map[string]any{
									"schema": map[string]any{"$ref": "#/components/schemas/TokenResponse"},
								},
							},
						},
					},
				},
			},
			"/qr":     protectedMatrixOperation("Calcular factorizacion QR", "QRResult"),
			"/rotate": protectedMatrixOperation("Rotar una matriz en sentido horario", "RotationResult"),
			"/analyze": map[string]any{
				"post": map[string]any{
					"tags":     []string{"Matrix"},
					"summary":  "Calcular QR y solicitar estadisticas a la API Node",
					"security": []map[string]any{{"bearerAuth": []string{}}},
					"requestBody": map[string]any{
						"required": true,
						"content": map[string]any{
							"application/json": map[string]any{
								"schema": map[string]any{"$ref": "#/components/schemas/MatrixRequest"},
							},
						},
					},
					"responses": map[string]any{
						"200": responseWithSchema("Resultado QR y estadisticas", "AnalyzeResult"),
						"400": map[string]any{"description": "Matriz invalida"},
						"401": map[string]any{"description": "JWT faltante o invalido"},
						"502": map[string]any{"description": "API Node no disponible"},
					},
				},
			},
		},
	}
}

func protectedMatrixOperation(summary string, responseSchema string) map[string]any {
	return map[string]any{
		"post": map[string]any{
			"tags":     []string{"Matrix"},
			"summary":  summary,
			"security": []map[string]any{{"bearerAuth": []string{}}},
			"requestBody": map[string]any{
				"required": true,
				"content": map[string]any{
					"application/json": map[string]any{
						"schema": map[string]any{"$ref": "#/components/schemas/MatrixRequest"},
					},
				},
			},
			"responses": map[string]any{
				"200": responseWithSchema("Operacion ejecutada correctamente", responseSchema),
				"400": map[string]any{"description": "Matriz invalida"},
				"401": map[string]any{"description": "JWT faltante o invalido"},
			},
		},
	}
}

func responseWithSchema(description string, schema string) map[string]any {
	return map[string]any{
		"description": description,
		"content": map[string]any{
			"application/json": map[string]any{
				"schema": map[string]any{"$ref": "#/components/schemas/" + schema},
			},
		},
	}
}
