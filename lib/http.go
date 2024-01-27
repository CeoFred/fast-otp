package httpclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

var (
	FastOTPClient = &http.Client{
		Timeout: time.Duration(5) * time.Second,
		Transport: &http.Transport{
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 100,
			IdleConnTimeout:     30 * time.Second,
		},
	}
)

// APIClient is a wrapper for making HTTP requests to the fastotp API.
type APIClient struct {
	baseURL string
	apiKey  string
	client  *http.Client
}

// NewAPIClient creates a new instance of APIClient.
func NewAPIClient(baseURL, apiKey string) *APIClient {
	return &APIClient{
		baseURL: baseURL,
		apiKey:  apiKey,
		client:  FastOTPClient,
	}
}

// Post sends a POST request to the specified endpoint with the given payload.
func (c *APIClient) Post(endpoint string, payload interface{}) (*http.Response, error) {
	url := c.baseURL + endpoint

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	return c.client.Do(req)
}

// Get sends a GET request to the specified endpoint, appending id as a path parameter
func (c *APIClient) Get(id string) (*http.Response, error) {
	url := fmt.Sprintf("%s/%s", c.baseURL, id)

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", c.apiKey)

	return c.client.Do(req)
}
