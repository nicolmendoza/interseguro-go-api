package main

import (
	"errors"
	"fmt"
	"math"
)

type Matrix [][]float64

const epsilon = 1e-10

func (matrix Matrix) Validate() error {
	if err := matrix.ValidateRectangular(); err != nil {
		return err
	}

	if len(matrix) < len(matrix[0]) {
		return errors.New("QR factorization expects rows >= columns for this compact implementation")
	}

	return nil
}

func (matrix Matrix) ValidateRectangular() error {
	if len(matrix) == 0 {
		return errors.New("matrix must contain at least one row")
	}
	if len(matrix[0]) == 0 {
		return errors.New("matrix must contain at least one column")
	}

	columns := len(matrix[0])
	for i, row := range matrix {
		if len(row) != columns {
			return fmt.Errorf("matrix must be rectangular: row %d has a different length", i)
		}
		for _, value := range row {
			if math.IsNaN(value) || math.IsInf(value, 0) {
				return errors.New("matrix values must be finite numbers")
			}
		}
	}

	return nil
}

func QRFactorization(a Matrix) (Matrix, Matrix, error) {
	if err := a.Validate(); err != nil {
		return nil, nil, err
	}

	rows := len(a)
	columns := len(a[0])
	q := zeroMatrix(rows, columns)
	r := zeroMatrix(columns, columns)

	for j := 0; j < columns; j++ {
		v := column(a, j)
		for i := 0; i < j; i++ {
			qi := column(q, i)
			r[i][j] = dot(qi, v)
			for k := 0; k < rows; k++ {
				v[k] -= r[i][j] * qi[k]
			}
		}

		norm := vectorNorm(v)
		if norm < epsilon {
			return nil, nil, errors.New("matrix columns must be linearly independent")
		}

		r[j][j] = norm
		for k := 0; k < rows; k++ {
			q[k][j] = round(v[k] / norm)
		}
	}

	return q.rounded(), r.rounded(), nil
}

func RotateClockwise(matrix Matrix) Matrix {
	rows := len(matrix)
	columns := len(matrix[0])
	rotated := zeroMatrix(columns, rows)

	for row := 0; row < rows; row++ {
		for col := 0; col < columns; col++ {
			rotated[col][rows-1-row] = matrix[row][col]
		}
	}

	return rotated
}

func zeroMatrix(rows int, columns int) Matrix {
	matrix := make(Matrix, rows)
	for i := range matrix {
		matrix[i] = make([]float64, columns)
	}
	return matrix
}

func column(matrix Matrix, index int) []float64 {
	values := make([]float64, len(matrix))
	for i := range matrix {
		values[i] = matrix[i][index]
	}
	return values
}

func dot(a []float64, b []float64) float64 {
	total := 0.0
	for i := range a {
		total += a[i] * b[i]
	}
	return total
}

func vectorNorm(values []float64) float64 {
	return math.Sqrt(dot(values, values))
}

func (matrix Matrix) rounded() Matrix {
	for i := range matrix {
		for j := range matrix[i] {
			matrix[i][j] = round(matrix[i][j])
		}
	}
	return matrix
}

func round(value float64) float64 {
	if math.Abs(value) < epsilon {
		return 0
	}
	return math.Round(value*1e10) / 1e10
}
