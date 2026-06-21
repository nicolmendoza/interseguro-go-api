package matrix

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
