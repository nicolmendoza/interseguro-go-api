package dto

import matrixdomain "interseguro/go-api/internal/domain/matrix"

type MatrixRequest struct {
	Matrix matrixdomain.Matrix `json:"matrix"`
}
