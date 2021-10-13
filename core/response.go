package core

import (
	"encoding/json"
	"net/http"
)

type Response struct {
	Status string
	StatusCode int
	Headers http.Header
	Body []byte
}

// Bytes returns bytes representation of the response
func (r *Response) Bytes() []byte {
	return r.Body
}

// String returns string representation of the response
func (r *Response) String() string {
	return string(r.Body)
}

// UnmarshalJson used to parse custom structs with the response
func (r *Response) UnmarshalJson(target interface{}) error {
	return json.Unmarshal(r.Bytes(), target)
}