# Interseguro Go API

API REST construida con Go y Fiber.

Su responsabilidad es recibir una matriz, validar la entrada, calcular factorizacion QR, rotar la matriz y coordinar con la API Node para obtener estadisticas sobre las matrices resultantes.

## Requisitos

- Go 1.22+
- Node.js 20+, solo si se desea ejecutar scripts auxiliares con npm
- Docker

La forma recomendada de levantar el proyecto es con Docker, porque el reto solicita contenerizar las aplicaciones y facilita ejecutar el servicio con configuracion consistente.

## Produccion

Servicio desplegado en Google Cloud Run:

```txt
Go API:   https://interseguro-go-api-745150536858.europe-west1.run.app
```

## Despliegue en Google Cloud

Este servicio fue desplegado en Google Cloud Run usando un contenedor Docker:

- `interseguro-go-api`: API Go/Fiber para QR, rotacion y orquestacion.

El servicio se construye como imagen Docker y se despliega de forma independiente en Cloud Run. El codigo fuente esta versionado en GitHub y Cloud Build queda conectado al repositorio para compilar la imagen y publicar una nueva revision cuando se suben cambios.

Las variables de entorno se configuran desde Cloud Run, no quedan hardcodeadas en el codigo fuente. En produccion la API Go usa:

```txt
JWT_SECRET=<configurado en Cloud Run>
NODE_API_URL=https://interseguro-node-api-745150536858.europe-west1.run.app
```

`PORT` no se define manualmente en Cloud Run porque la plataforma lo inyecta automaticamente en el contenedor.

## Variables de entorno

```txt
PORT
JWT_SECRET
NODE_API_URL
```

## Ejecutar con Docker

Crear red compartida:

```bash
docker network create interseguro-net
```

Construir imagen:

```bash
docker build -t interseguro-go-api .
```

Ejecutar contenedor:

```bash
docker run --rm --name interseguro-go-api \
  --network interseguro-net \
  -e PORT=3000 \
  -e JWT_SECRET=local-development-secret \
  -e NODE_API_URL=http://interseguro-node-api:8080 \
  -p 3000:3000 \
  interseguro-go-api
```

La API queda en:

```txt
http://localhost:3000
```

## Tests

```bash
npm test
```

Cobertura actual:

- `domain/matrix/matrix_test.go`: QR y rotacion.
- `infrastructure/statsclient/node_stats_client_test.go`: cliente HTTP hacia Node con servidor falso.
- `transport/http/server_test.go`: endpoints protegidos con JWT.

## Arquitectura hexagonal

La API usa una arquitectura hexagonal ligera.

```txt
internal/
  domain/
    matrix/
      matrix.go
      validation.go
      qr.go
      rotation.go
      matrix_test.go

  application/
    matrix/
      service.go
      ports.go

  infrastructure/
    config/
      config.go
    statsclient/
      node_stats_client.go
      node_stats_client_test.go

  transport/
    http/
      server.go
      routes.go
      dto/
      handlers/
      middleware/
      server_test.go

main.go
```

Responsabilidades:

- `domain`: reglas puras de matriz, QR, rotacion y validaciones.
- `application`: caso de uso que orquesta dominio y cliente de estadisticas.
- `infrastructure`: variables de entorno y cliente HTTP hacia Node.
- `transport/http`: adaptador HTTP con Fiber, rutas, handlers, DTOs y JWT.

## Endpoints

```txt
GET  /health
POST /auth/token
POST /qr
POST /rotate
POST /analyze
```

Los endpoints `/qr`, `/rotate` y `/analyze` requieren JWT:

```txt
Authorization: Bearer <token>
```

## Estructura de datos

```go
type Matrix [][]float64

type MatrixRequest struct {
    Matrix Matrix `json:"matrix"`
}

type QRResult struct {
    Q Matrix `json:"q"`
    R Matrix `json:"r"`
}

type RotationResult struct {
    Rotated Matrix `json:"rotated"`
}

type AnalyzeResult struct {
    Q     Matrix          `json:"q"`
    R     Matrix          `json:"r"`
    Stats json.RawMessage `json:"stats"`
}
```

Ejemplo request:

```json
{
  "matrix": [
    [12, -51, 4],
    [6, 167, -68],
    [-4, 24, -41]
  ]
}
```

Validaciones:

- La matriz debe tener al menos una fila.
- La matriz debe tener al menos una columna.
- Todas las filas deben tener la misma cantidad de columnas.
- Los valores deben ser numeros finitos.
- Para QR, `filas >= columnas`.
- Para QR con Gram-Schmidt, las columnas deben ser linealmente independientes.

## Flujo principal

```txt
Frontend
  -> Go API /analyze
      -> valida matriz
      -> calcula Q y R
      -> llama Node API /stats con Q y R
      -> devuelve q, r y stats

Frontend
  -> Go API /rotate
      -> devuelve matriz rotada
```
