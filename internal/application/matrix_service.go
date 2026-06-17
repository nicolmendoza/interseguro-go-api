package application

import (
	"encoding/json"

	"interseguro/go-api/internal/domain"
)

type StatsClient interface {
	CalculateStats(matrices []domain.Matrix, bearerToken string) (json.RawMessage, int, error)
}

type QRResult struct {
	Q domain.Matrix `json:"q"`
	R domain.Matrix `json:"r"`
}

type AnalyzeResult struct {
	Q     domain.Matrix   `json:"q"`
	R     domain.Matrix   `json:"r"`
	Stats json.RawMessage `json:"stats"`
}

type MatrixService struct {
	statsClient StatsClient
}

func NewMatrixService(statsClient StatsClient) MatrixService {
	return MatrixService{statsClient: statsClient}
}

func (service MatrixService) Factorize(matrix domain.Matrix) (QRResult, error) {
	q, r, err := domain.QRFactorization(matrix)
	if err != nil {
		return QRResult{}, err
	}
	return QRResult{Q: q, R: r}, nil
}

func (service MatrixService) Rotate(matrix domain.Matrix) (domain.Matrix, error) {
	if err := matrix.ValidateRectangular(); err != nil {
		return nil, err
	}
	return domain.RotateClockwise(matrix), nil
}

func (service MatrixService) Analyze(matrix domain.Matrix, bearerToken string) (AnalyzeResult, int, error) {
	result, err := service.Factorize(matrix)
	if err != nil {
		return AnalyzeResult{}, 0, err
	}

	stats, statusCode, err := service.statsClient.CalculateStats([]domain.Matrix{result.Q, result.R}, bearerToken)
	if err != nil {
		return AnalyzeResult{}, statusCode, err
	}

	return AnalyzeResult{Q: result.Q, R: result.R, Stats: stats}, statusCode, nil
}
