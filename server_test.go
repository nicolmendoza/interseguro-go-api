package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"
)

func TestQRRequiresJWT(t *testing.T) {
	app := NewServer(Config{JWTSecret: "test-secret", NodeURL: "http://localhost:3001"})

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
	app := NewServer(Config{JWTSecret: "test-secret", NodeURL: "http://localhost:3001"})
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

func getTestToken(t *testing.T, app interface {
	Test(req *http.Request, msTimeout ...int) (*http.Response, error)
}) string {
	t.Helper()
	req, _ := http.NewRequest(http.MethodPost, "/auth/token", nil)
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
