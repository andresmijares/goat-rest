package goat

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/andresmijares/goat-rest/core"
	"github.com/andresmijares/goat-rest/goat_mock"
	"github.com/andresmijares/goat-rest/mime"
)

var (
	defaultMaxIdleConnections = 5
	defaultTimeout            = 5 * time.Second
	defaultConnectionTimeout  = 5 * time.Second

	errInvalidRequest = errors.New("unable to perform request")
)

func (c *httpClient) do(method string, url string, headers http.Header, body interface{}) (*core.Response, error) {
	allHeaders := c.setHeaders(headers)

	requestBody, err := c.getRequestBody(headers.Get(mime.HeaderContentType), body)
	if err != nil {
		return nil, err
	}

	request, err := http.NewRequest(method, url, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, errInvalidRequest
	}

	request.Header = allHeaders

	c.client = c.createHttpClient()

	response, err := c.client.Do(request)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	responseBody, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return nil, err
	}

	return &core.Response{
		Status:     response.Status,
		StatusCode: response.StatusCode,
		Headers:    response.Header,
		Body:       responseBody,
	}, nil
}

func (c *httpClient) setHeaders(requestHeader http.Header) http.Header {
	h := make(http.Header)

	// commons headers
	for header, value := range c.config.headers {
		if len(value) > 0 {
			h.Set(header, value[0])
		}
	}

	// custom headers
	for header, value := range requestHeader {
		if len(value) > 0 {
			h.Set(header, value[0])
		}
	}

	// support custom Agent if not already created
	if c.config.agent != "" {
		if h.Get(mime.HeaderUserAgent) != "" {
			return h
		}
		h.Set(mime.HeaderUserAgent, c.config.agent)
	}

	return h
}

func (c *httpClient) createHttpClient() core.HttpClient {
	// Enables support for mocked server is enabled
	if goat_mock.MockupServer.IsEnabled() {
		return goat_mock.MockupServer.GetClient()
	}

	// ensures the client is instanced only one
	// even if using multiple goroutines
	c.clientOnce.Do(func() {
		if c.config.client != nil {
			// consumer has its own client
			c.client = c.config.client
			return
		}

		c.client = &http.Client{
			Timeout: c.getConnectionTimeout() + c.getResponseTimeout(), // includes the whole rountrip timeout, zero means no timeout
			Transport: &http.Transport{
				MaxIdleConnsPerHost:   c.getMaxIdleConnections(), // Max connection in idle state
				ResponseHeaderTimeout: c.getResponseTimeout(),    // how long do we wait for a request to response
				DialContext: (&net.Dialer{
					Timeout: c.getConnectionTimeout(),
				}).DialContext, // how long do we wait for a new connection until timeout
			},
		}
	})

	return c.client
}

func (c *httpClient) getRequestBody(contentType string, body interface{}) ([]byte, error) {
	if body == nil {
		return nil, nil
	}

	switch strings.ToLower(contentType) {
	case mime.ApplicationTypeJSON:
		return json.Marshal(body)
	case mime.ApplicationTypeXML:
		return xml.Marshal(body)
	default:
		return json.Marshal(body)
	}
}

/***
 The following private methods, only set defaults for the client configuration

 Methods name should match the action plus prop name:
 	ex getMaxIdleConnections -> would get default for maxIdleConnections
***/
func (c *httpClient) getMaxIdleConnections() int {
	if c.config.maxIdleConnections > 0 {
		return c.config.maxIdleConnections
	}
	return defaultMaxIdleConnections
}

func (c *httpClient) getResponseTimeout() time.Duration {
	// allows to disable timeouts if requested
	if c.config.disableTimeouts {
		return 0
	}

	if c.config.responseTimeout > 0 {
		return c.config.responseTimeout
	}

	return defaultTimeout
}

func (c *httpClient) getConnectionTimeout() time.Duration {
	// allows to disable timeouts if requested
	if c.config.disableTimeouts {
		return 0
	}

	if c.config.connectionTimeout > 0 {
		return c.config.connectionTimeout
	}

	return defaultConnectionTimeout
}
