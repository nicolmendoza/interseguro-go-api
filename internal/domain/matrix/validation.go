package matrix

import (
	"errors"
	"fmt"
	"math"
)

func (matrix Matrix) ValidateForQR() error {

	if err := matrix.ValidateRectangular(); err != nil {
		return err
	}

	if len(matrix) < len(matrix[0]) {
		return errors.New("la factorizacion QR requiere que la cantidad de filas sea mayor o igual a la cantidad de columnas")
	}

	return nil
}

func (matrix Matrix) ValidateRectangular() error {

	if len(matrix) == 0 {
		return errors.New("la matriz debe contener al menos una fila")
	}

	if len(matrix[0]) == 0 {
		return errors.New("la matriz debe contener al menos una columna")
	}

	columns := len(matrix[0])
	for i, row := range matrix {
		if len(row) != columns {
			return fmt.Errorf("la matriz debe ser rectangular: la fila %d tiene una longitud diferente", i)
		}
		for _, value := range row {
			if math.IsNaN(value) || math.IsInf(value, 0) {
				return errors.New("los valores de la matriz deben ser numeros finitos")
			}
		}
	}

	return nil
}
