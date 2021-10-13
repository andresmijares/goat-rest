package goat

import (
	"net/http"
	"time"
)

type Config interface {
	// SetHeaders Attach custom headers to the request
	SetHeaders(headers http.Header) Config
	// SetConnectionTimeout sets long do you wait for a request to response
	SetConnectionTimeout(timeout time.Duration) Config
	// SetRequestTimeout sets how long do we wait for a new connection until timeout
	SetResponseTimeout(timeout time.Duration) Config
	// SetMaxIdleConnections sets max connection in idle state
	SetMaxIdleConnections(connections int) Config
	// DisableTimeout Disables all timeouts, it's true by default
	DisableTimeouts(disable bool) Config
	// SetHttpClient allows to use a custom httpClient, when using a custom client, none of the internal configuration defaults are applied, ex: timeouts
	SetHttpClient(c *http.Client) Config
	// SetUserAgent allows setting a user agent for requests
	SetUserAgent(agent string) Config
	// Build Returns a HTTP Client interface, it should be you called after set all the custom configuration
	Create() Client
}

type config struct {
	maxIdleConnections int
	connectionTimeout time.Duration
	responseTimeout time.Duration
	disableTimeouts bool

	headers http.Header
	client *http.Client
	agent string
}

func New() Config {
	 config := &config{}
	 return config
}

// Build creates custom http client, it should be called after all custom configurations
// (if any)
func (c *config) Create() Client {
	client := &httpClient{
		config: c,
	}
	return client
}

// SetHttpClient allows to use a custom httpClient, when using a custom client,
// none of the internal configuration defaults are applied, ex: timeouts
func (c *config) SetHttpClient(client *http.Client) Config {
	c.client = client
	return c
}

// DisableTimeout Disable timeouts, it's disable by default
func (c *config) DisableTimeouts(disable bool) Config {
	c.disableTimeouts = disable
	return c
}

// SetHeaders Attach custom headers to the request
func (c *config) SetHeaders(headers http.Header) Config {
	c.headers = headers
	return c
}

// SetConnectionTimeout sets long do you wait for a request to response
func (c *config) SetConnectionTimeout(timeout time.Duration) Config {
	c.connectionTimeout = timeout
	return c
}

// SetRequestTimeout sets how long do we wait for a new connection until timeout
func (c *config) SetResponseTimeout(timeout time.Duration) Config {
	c.responseTimeout = timeout
	return c
}

// SetMaxIdleConnections sets max connection in idle state
func (c *config) SetMaxIdleConnections(connections int) Config {
	c.maxIdleConnections = connections
	return c
}

// SetUserAgent allows setting a user agent for requests
func (c *config) SetUserAgent(agent string) Config {
	c.agent = agent
	return c
}