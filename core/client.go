package core

import "net/http"

// HttpClient the only purpose for this interface is to be able to mock
// tests when using a mock server instead of using the client's method directly
type HttpClient interface {
	Do(request *http.Request) (*http.Response, error)
}
