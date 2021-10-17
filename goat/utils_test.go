package goat

import (
	"net/http"
	"testing"
)

func TestGetHeaders(t *testing.T) {

	t.Run("TestGetHeaders with length", func(t *testing.T) {
		headers := make(http.Header)
		headers.Add("Foo", "Bar")

		h := getHeaders(headers)

		for k, _ := range h {
			if k != "Foo" {
				t.Errorf("invalid header")
			}
		}

	})

	t.Run("TestGetHeaders with no length", func(t *testing.T) {
		h := getHeaders()
		if len(h) != 0 {
			t.Errorf("invalid number of headers")
		}
	})
}
