package httptransport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	matrixapp "interseguro/go-api/internal/application/matrix"
	matrixdomain "interseguro/go-api/internal/domain/matrix"
)

type fakeStatsClient struct{}

func (fakeStatsClient) CalculateStats(_ []matrixdomain.Matrix, _ string) (json.RawMessage, int, error) {
	return json.RawMessage(`{"max":1,"min":0}`), http.StatusOK, nil
}

func TestQRRequiresJWT(t *testing.T) {

	app := newTestServer()
	req := newTestRequest(t, http.MethodPost, "/qr", `{"matrix":[[1,0],[0,1]]}`)
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("la solicitud fallo: %v", err)
	}

	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("se esperaba estado 401, se obtuvo %d", resp.StatusCode)
	}
}

func TestQRWithJWT(t *testing.T) {

	app := newTestServer()
	token := getTestToken(t, app)
	req := newTestRequest(t, http.MethodPost, "/qr", `{"matrix":[[1,0],[0,1]]}`)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("la solicitud fallo: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("se esperaba estado 200, se obtuvo %d", resp.StatusCode)
	}
}

func newTestServer() interface {
	Test(req *http.Request, msTimeout ...int) (*http.Response, error)
} {

	matrixService := matrixapp.NewService(fakeStatsClient{})
	return NewServer(Config{JWTSecret: "test-secret"}, matrixService)
}

func getTestToken(t *testing.T, app interface {
	Test(req *http.Request, msTimeout ...int) (*http.Response, error)
}) string {

	t.Helper()
	req := newTestRequest(t, http.MethodPost, "/auth/token", "")
	resp, err := app.Test(req)

	if err != nil {
		t.Fatalf("la solicitud del token fallo: %v", err)
	}

	var body struct {
		Token string `json:"token"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("no se pudo decodificar la respuesta del token: %v", err)
	}

	return body.Token
}

func newTestRequest(t *testing.T, method string, path string, body string) *http.Request {

	t.Helper()
	req, err := http.NewRequest(method, path, bytes.NewBufferString(body))
	if err != nil {
		t.Fatalf("no se pudo crear la solicitud: %v", err)
	}

	return req
}
