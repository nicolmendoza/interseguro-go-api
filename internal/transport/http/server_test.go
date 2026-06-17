package httptransport

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"interseguro/go-api/internal/application"
	"interseguro/go-api/internal/domain"
)

type fakeStatsClient struct{}

func (fakeStatsClient) CalculateStats(_ []domain.Matrix, _ string) (json.RawMessage, int, error) {
	return json.RawMessage(`{"max":1,"min":0}`), http.StatusOK, nil
}

func TestQRRequiresJWT(t *testing.T) {
	app := newTestServer()

	req, _ := http.NewRequest(http.MethodPost, "/qr", bytes.NewBufferString(`{"matrix":[[1,0],[0,1]]}`))
	req.Header.Set("Content-Type", "application/json")
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusUnauthorized {
		t.Fatalf("expected 401, got %d", resp.StatusCode)
	}
}

func TestQRWithJWT(t *testing.T) {
	app := newTestServer()
	token := getTestToken(t, app)

	req, _ := http.NewRequest(http.MethodPost, "/qr", bytes.NewBufferString(`{"matrix":[[1,0],[0,1]]}`))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("request failed: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected 200, got %d", resp.StatusCode)
	}
}

func newTestServer() interface {
	Test(req *http.Request, msTimeout ...int) (*http.Response, error)
} {
	matrixService := application.NewMatrixService(fakeStatsClient{})
	return NewServer(Config{JWTSecret: "test-secret"}, matrixService)
}

func getTestToken(t *testing.T, app interface {
	Test(req *http.Request, msTimeout ...int) (*http.Response, error)
}) string {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, "/auth/token", bytes.NewBufferString(""))
	resp, err := app.Test(req)
	if err != nil {
		t.Fatalf("token request failed: %v", err)
	}

	var body struct {
		Token string `json:"token"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		t.Fatalf("could not decode token response: %v", err)
	}
	return body.Token
}
