package weavychat

import (
	"context"
	"net/http"
)

// HttpClient is the interface for the HTTP client.
type HttpClient interface {
	Get(ctx context.Context, id string) (*http.Response, error)
	Post(ctx context.Context, endpoint string, payload interface{}) (*http.Response, error)
	Delete(ctx context.Context, endpoint string, payload interface{}) (*http.Response, error)
}
