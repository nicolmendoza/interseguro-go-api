package statsclient

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	matrixdomain "interseguro/go-api/internal/domain/matrix"
)

func TestNodeStatsClientCalculateStats(t *testing.T) {
	expectedToken := "Bearer test-token"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodPost {
			t.Fatalf("se esperaba POST, se obtuvo %s", r.Method)
		}

		if r.URL.Path != "/stats" {
			t.Fatalf("se esperaba /stats, se obtuvo %s", r.URL.Path)
		}

		if r.Header.Get("Authorization") != expectedToken {
			t.Fatalf("se esperaba header Authorization %q, se obtuvo %q", expectedToken, r.Header.Get("Authorization"))
		}

		if r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("se esperaba Content-Type JSON, se obtuvo %q", r.Header.Get("Content-Type"))
		}

		var body struct {
			Matrices []matrixdomain.Matrix `json:"matrices"`
		}

		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Fatalf("no se pudo decodificar el cuerpo de la solicitud: %v", err)
		}

		if len(body.Matrices) != 1 || body.Matrices[0][0][0] != 1 {
			t.Fatalf("payload de matrices inesperado: %v", body.Matrices)
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"max":1,"min":0}`))
	}))
	defer server.Close()

	client := NewNodeStatsClient(server.URL)
	body, statusCode, err := client.CalculateStats([]matrixdomain.Matrix{{{1, 0}, {0, 1}}}, expectedToken)

	if err != nil {
		t.Fatalf("se esperaba una respuesta exitosa, se obtuvo error: %v", err)
	}

	if statusCode != http.StatusOK {
		t.Fatalf("se esperaba estado 200, se obtuvo %d", statusCode)
	}

	if string(body) != `{"max":1,"min":0}` {
		t.Fatalf("cuerpo de respuesta inesperado: %s", body)
	}
}
