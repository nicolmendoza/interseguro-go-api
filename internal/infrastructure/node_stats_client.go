package infrastructure

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	"interseguro/go-api/internal/domain"
)

type NodeStatsClient struct {
	baseURL string
	client  http.Client
}

func NewNodeStatsClient(baseURL string) NodeStatsClient {
	return NodeStatsClient{
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  http.Client{Timeout: 5 * time.Second},
	}
}

func (client NodeStatsClient) CalculateStats(matrices []domain.Matrix, bearerToken string) (json.RawMessage, int, error) {
	payload, _ := json.Marshal(map[string][]domain.Matrix{"matrices": matrices})
	req, err := http.NewRequest(http.MethodPost, client.baseURL+"/stats", bytes.NewReader(payload))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken)

	resp, err := client.client.Do(req)
	if err != nil {
		return nil, http.StatusBadGateway, errors.New("Node API is unavailable")
	}
	defer resp.Body.Close()

	var body json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, http.StatusBadGateway, errors.New("invalid Node API response")
	}

	return body, resp.StatusCode, nil
}
