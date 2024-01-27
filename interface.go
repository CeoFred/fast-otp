package fastotp

import (
	"net/http"
)

// HttpClient is the interface for the HTTP client.
type HttpClient interface {
	Get(id string) (*http.Response, error)
	Post(endpoint string, payload interface{}) (*http.Response, error)
}
