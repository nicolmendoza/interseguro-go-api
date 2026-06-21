package matrix

import (
	"encoding/json"

	matrixdomain "interseguro/go-api/internal/domain/matrix"
)

type QRResult struct {
	Q matrixdomain.Matrix `json:"q"`
	R matrixdomain.Matrix `json:"r"`
}

type RotationResult struct {
	Rotated matrixdomain.Matrix `json:"rotated"`
}

type AnalyzeResult struct {
	Q     matrixdomain.Matrix `json:"q"`
	R     matrixdomain.Matrix `json:"r"`
	Stats json.RawMessage     `json:"stats"`
}

type Service struct {
	statsClient StatsClient
}

func NewService(statsClient StatsClient) Service {
	return Service{statsClient: statsClient}
}

func (service Service) Factorize(input matrixdomain.Matrix) (QRResult, error) {
	q, r, err := matrixdomain.QRFactorization(input)
	if err != nil {
		return QRResult{}, err
	}
	return QRResult{Q: q, R: r}, nil
}

func (service Service) Rotate(input matrixdomain.Matrix) (RotationResult, error) {
	if err := input.ValidateRectangular(); err != nil {
		return RotationResult{}, err
	}
	return RotationResult{Rotated: matrixdomain.RotateClockwise(input)}, nil
}

func (service Service) Analyze(input matrixdomain.Matrix, bearerToken string) (AnalyzeResult, int, error) {
	result, err := service.Factorize(input)
	if err != nil {
		return AnalyzeResult{}, 0, err
	}

	stats, statusCode, err := service.statsClient.CalculateStats([]matrixdomain.Matrix{result.Q, result.R}, bearerToken)
	if err != nil {
		return AnalyzeResult{}, statusCode, err
	}

	return AnalyzeResult{Q: result.Q, R: result.R, Stats: stats}, statusCode, nil
}
