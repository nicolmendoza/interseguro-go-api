package statsclient

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"time"

	matrixdomain "interseguro/go-api/internal/domain/matrix"
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

func (client NodeStatsClient) CalculateStats(matrices []matrixdomain.Matrix, bearerToken string) (json.RawMessage, int, error) {

	payload, err := json.Marshal(map[string][]matrixdomain.Matrix{"matrices": matrices})
	if err != nil {
		return nil, http.StatusInternalServerError, errors.New("no se pudo preparar la solicitud para la API de Node")
	}

	req, err := http.NewRequest(http.MethodPost, client.baseURL+"/stats", bytes.NewReader(payload))
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", bearerToken)
	resp, err := client.client.Do(req)
	if err != nil {
		return nil, http.StatusBadGateway, errors.New("la API de Node no esta disponible")
	}

	defer resp.Body.Close()
	var body json.RawMessage
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return nil, http.StatusBadGateway, errors.New("la API de Node devolvio una respuesta invalida")
	}

	return body, resp.StatusCode, nil
}
