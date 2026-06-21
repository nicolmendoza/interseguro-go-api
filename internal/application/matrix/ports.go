package matrix

import (
	"encoding/json"

	matrixdomain "interseguro/go-api/internal/domain/matrix"
)

type StatsClient interface {
	CalculateStats(matrices []matrixdomain.Matrix, bearerToken string) (json.RawMessage, int, error)
}
