package weavychat

import (
	"context"
	"net/http"
)

type mockedHTTPClient struct {
	GetFunc  func(ctx context.Context, id string) (*http.Response, error)
	PostFunc func(ctx context.Context, endpoint string, payload interface{}) (*http.Response, error)
}

func (m mockedHTTPClient) Get(ctx context.Context, id string) (*http.Response, error) {
	return m.GetFunc(ctx, id)
}

func (m mockedHTTPClient) Post(ctx context.Context, endpoint string, payload interface{}) (*http.Response, error) {
	return m.PostFunc(ctx, endpoint, payload)
}
