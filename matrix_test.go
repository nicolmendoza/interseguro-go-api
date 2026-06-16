package main

import (
	"math"
	"testing"
)

func TestQRFactorization(t *testing.T) {
	matrix := Matrix{
		{12, -51, 4},
		{6, 167, -68},
		{-4, 24, -41},
	}

	q, r, err := QRFactorization(matrix)
	if err != nil {
		t.Fatalf("expected QR factorization, got error: %v", err)
	}

	assertClose(t, q[0][0], 0.8571428571)
	assertClose(t, q[1][0], 0.4285714286)
	assertClose(t, q[2][0], -0.2857142857)
	assertClose(t, r[0][0], 14)
	assertClose(t, r[0][1], 21)
	assertClose(t, r[0][2], -14)
}

func TestRotateClockwise(t *testing.T) {
	got := RotateClockwise(Matrix{
		{1, 2, 3},
		{4, 5, 6},
	})

	want := Matrix{
		{4, 1},
		{5, 2},
		{6, 3},
	}

	for i := range want {
		for j := range want[i] {
			if got[i][j] != want[i][j] {
				t.Fatalf("expected %v, got %v", want, got)
			}
		}
	}
}

func assertClose(t *testing.T, got float64, want float64) {
	t.Helper()
	if math.Abs(got-want) > 1e-8 {
		t.Fatalf("expected %v, got %v", want, got)
	}
}
