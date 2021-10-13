package goat_mock

import (
	"fmt"
	"net/http"

	"github.com/andresmijares/goat-rest/core"
)

// Mock Provides a way to create mocks based on whatever
// convination needed to test
type Mock struct {
	Method string 
	URL string
	RequestBody string
	Error error
	ResponseBody string
	ResponseStatusCode int
}

// GetResponse gets a response object based on the mock configuration
func (m *Mock) GetResponse() (*core.Response, error) {
	if m.Error != nil {
		return nil, nil
	}

	response := core.Response{
		Status:     fmt.Sprintf("%d %s", m.ResponseStatusCode, http.StatusText(m.ResponseStatusCode)),
		StatusCode: m.ResponseStatusCode,
		Headers:    map[string][]string{},
		Body:       []byte{},
	}
	
	return &response, nil
	
}