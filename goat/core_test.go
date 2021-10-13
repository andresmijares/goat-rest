package goat

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/andresmijares/goat-rest/mime"
)

func TestSetRequestHeaders(t *testing.T) {
	t.Run("Custom headers are added", func(t *testing.T) {
		cfg := config{}
		client := httpClient{
			config: &cfg,
		}

		headers := make(http.Header)
		headers.Set("Content-Type", "application/json")
		headers.Set("X-Correlation-Id", "ABC123")
		client.config.headers = headers

		requestHeaders := make(http.Header)
		requestHeaders.Set("X-AWS-XRAY-ID", "AWS12345")

		fullHeaders := client.setHeaders(requestHeaders)
		if len(fullHeaders) != 3 {
			t.Errorf("Number of headers dont match")
		}
	})

	t.Run("Custom Agent is set if present", func(t *testing.T) {
		customAgent := "Custom-Agent-Id"
		cfg := config{}
		client := httpClient{
			config: &cfg,
		}

		cfg.agent = customAgent
		requestHeaders := make(http.Header)
		requestHeaders.Set("X-AWS-XRAY-ID", "AWS12345")

		headers := client.setHeaders(requestHeaders)

		if headers.Get(mime.HeaderUserAgent) != customAgent {
			t.Errorf("Custom Agent does't match")
		}
	})

	t.Run("Custom Agent is NOT set not passed", func(t *testing.T) {
		expectAgent := "Custom-Agent-Id"
		configHeader := make(http.Header)
		configHeader.Add(mime.HeaderUserAgent, expectAgent)

		cfg := config{
			headers: configHeader,
		}
		client := httpClient{
			config: &cfg,
		}

		wrongAgent := "After-Set-Agent"
		cfg.agent = wrongAgent
		requestHeaders := make(http.Header)
		requestHeaders.Set("Random-Header", "true")
		headers := client.setHeaders(requestHeaders)

		if headers.Get(mime.HeaderUserAgent) != expectAgent {
			t.Errorf("Custom Agent doesnt match config agent")
		}

		if headers.Get(mime.HeaderUserAgent) == wrongAgent {
			t.Errorf("Custom Agent is overwriting config value")
		}
	})
}

func TestGetRequestBody(t *testing.T) {
	client := httpClient{}

	t.Run("NoBodyNilResponse", func(t *testing.T) {
		nilBody, err := client.getRequestBody("", nil)
		if err != nil {
			t.Errorf("it should return nil error")
		}

		if nilBody != nil {
			t.Errorf("it should return a nil body")
		}
	})

	t.Run("BodyIsJSON", func(t *testing.T) {
		requestBody := []string{"a", "b"}
		jsonBody, err := client.getRequestBody("application/json", requestBody)
		if err != nil {
			t.Errorf("it should return nil error")
		}

		if string(jsonBody) != `["a","b"]` {
			t.Errorf("it should return a json body")
		}
	})

	t.Run("BodyIsXML", func(t *testing.T) {
		requestBody := []string{"string", "sample"}
		xmlBody, err := client.getRequestBody("application/xml", requestBody)
		if err != nil {
			t.Errorf("it should return nil error")
		}

		if string(xmlBody) != `<string>string</string><string>sample</string>` {
			t.Errorf("it should return a xml body")
		}
	})

	t.Run("BodyIsJSONByDefault", func(t *testing.T) {
		requestBody := []string{"c"}
		customBody, err := client.getRequestBody("application/customType", requestBody)
		if err != nil {
			t.Errorf("it should return nil error")
		}

		if string(customBody) != `["c"]` {
			t.Errorf("it should return a json given a custom application type")
		}
	})
}

func TestSetMaxIdleConnections(t *testing.T) {
	cfg := config{}
	client := httpClient{
		config: &cfg,
	}

	t.Run("SetMaxIdleConnectionsDefault", func(t *testing.T) {
		if client.getMaxIdleConnections() != defaultMaxIdleConnections {
			t.Errorf("geMaxIdleConnections default doesnt match")
		}
	})

	t.Run("SetMaxIdleConnections", func(t *testing.T) {
		expectValue := 1
		cfg.maxIdleConnections = expectValue
		if client.getMaxIdleConnections() != expectValue {
			t.Errorf("geMaxIdleConnections doesnt match")
		}
	})
}

func TestSetResponseTimeout(t *testing.T) {
	cfg := config{}
	client := httpClient{
		config: &cfg,
	}

	t.Run("SetResponseTimeoutDefault", func(t *testing.T) {
		if client.getResponseTimeout() != defaultTimeout {
			t.Errorf("getResponseTimeout default doesnt match")
		}
	})

	t.Run("SetResponseTimeout", func(t *testing.T) {
		expectValue := time.Second * 5
		cfg.responseTimeout = expectValue
		if client.getResponseTimeout() != expectValue {
			t.Errorf("getConnectionTimeout doesnt match")
		}
	})

	t.Run("SetDisableTimeout", func(t *testing.T) {
		expectValue := time.Second * 5
		cfg.responseTimeout = expectValue
		cfg.disableTimeouts = true
		fmt.Print("disabled", client.getResponseTimeout())
		if client.getResponseTimeout() != 0 {
			t.Errorf("Disable Timeout doesnt work")
		}
	})
}

func TestSetConnectionTimeout(t *testing.T) {
	cfg := config{}
	client := httpClient{
		config: &cfg,
	}

	t.Run("SetConnectionTimeoutDefault", func(t *testing.T) {
		if client.getConnectionTimeout() != defaultConnectionTimeout {
			t.Errorf("getConnectionTimeout default doesnt match")
		}
	})

	t.Run("SetConnectionTimeout", func(t *testing.T) {
		expectValue := time.Second * 5
		cfg.connectionTimeout = expectValue
		if client.getConnectionTimeout() != expectValue {
			t.Errorf("getConnectionTimeout doesnt match")
		}
	})

	t.Run("SetDisableTimeout", func(t *testing.T) {
		expectValue := time.Second * 5
		cfg.connectionTimeout = expectValue
		cfg.disableTimeouts = true
		if client.getConnectionTimeout() != 0 {
			t.Errorf("Disable Connection doesnt work")
		}
	})
}

func TestDisableAllTimeouts(t *testing.T) {
	cfg := config{}
	client := httpClient{
		config: &cfg,
	}

	defaultValue := time.Second * 5
	cfg.connectionTimeout = defaultValue
	cfg.responseTimeout = defaultValue

	cfg.disableTimeouts = true

	if client.getConnectionTimeout() != 0 {
		t.Errorf("Disable Connection doesnt work on getConnectionTimeout")
	}

	if client.getResponseTimeout() != 0 {
		t.Errorf("Disable Connection doesnt work on getResponseTimeout")
	}

}

// func TestDo(t *testing.T) {
// 	client := httpClient{}

// 	_, err := client.do(http.MethodGet, "http://localhost:80", nil, nil)
// 	fmt.Print(err.Error())

// }
