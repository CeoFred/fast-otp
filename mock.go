package fastotp

import (
	"net/http"
)

type mockedHTTPClient struct {
	GetFunc  func(id string) (*http.Response, error)
	PostFunc func(endpoint string, payload interface{}) (*http.Response, error)
}

func (m mockedHTTPClient) Get(id string) (*http.Response, error) {
	return m.GetFunc(id)
}

func (m mockedHTTPClient) Post(endpoint string, payload interface{}) (*http.Response, error) {
	return m.PostFunc(endpoint, payload)
}
