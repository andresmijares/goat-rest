package goat

import (
	"errors"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/andresmijares/goat-rest/goat_mock"
	"github.com/andresmijares/goat-rest/mime"
)

func TestDo(t *testing.T) {
	// pass unmarshable interface
	t.Run("TestRequestBodyError", func(t *testing.T) {
		cfg := config{}
		client := httpClient{
			config: &cfg,
		}

		unmarshable := map[string]interface{}{
			"foo": make(chan int),
		}
		_, err := client.do(http.MethodGet, "http://localhost", nil, unmarshable)
		if err == nil {
			t.Errorf("Should return malformed marshal json")
		}
	})

	// pass a weird method name, it has to land into validateHttp method to throw
	t.Run("TestUnableToPerformRequest", func(t *testing.T) {
		cfg := config{}
		client := httpClient{
			config: &cfg,
		}

		_, err := client.do("[0:0]", "http://localhost", nil, nil)
		if err != errInvalidRequest {
			t.Errorf("Should return invalid http method")
		}
	})

	// connection refused
	t.Run("TestRequestResponseError", func(t *testing.T) {
		cfg := config{}
		client := httpClient{
			config: &cfg,
		}

		_, err := client.do(http.MethodPost, "http://localhost", nil, nil)
		if !strings.Contains(err.Error(), "connection refused") {
			t.Errorf("Should return connection refuse")
		}
	})

	// bad client response
	t.Run("TestRequestUnableToDecodeJsonResponse", func(t *testing.T) {
		expectedResponse := `{"foo": "bar"}`
		url := "http://127.0.0.1"
		sendBody := `{"foo":"bar"}`
		goat_mock.MockupServer.Start()
		goat_mock.MockupServer.Flush()
		goat_mock.MockupServer.Add(goat_mock.Mock{
			Method:             http.MethodPost,
			URL:                url,
			RequestBody:        sendBody,
			Error:              errors.New("Invalid client response"),
			ResponseBody:       expectedResponse,
			ResponseStatusCode: 200,
		})
		cfg := config{}

		client := httpClient{
			config: &cfg,
		}

		payload := map[string]interface{}{
			"foo": "bar",
		}

		_, err := client.do(http.MethodPost, url, nil, payload)
		if err.Error() != "Invalid client response" {
			t.Errorf("Should return bad response from client")
		}
		goat_mock.MockupServer.Stop()
	})

	// unparsable json
	// t.Run("TestRequestUnableToDecodeJsonResponse", func(t *testing.T) {
	// 	expectedResponse := `...`
	// 	url := "http://127.0.0.2"
	// 	sendBody := `{"foo":"bar"}`
	// 	goat_mock.MockupServer.Start()
	// 	goat_mock.MockupServer.Flush()
	// 	goat_mock.MockupServer.Add(goat_mock.Mock{
	// 		Method:             http.MethodPost,
	// 		URL:                url,
	// 		RequestBody:        sendBody,
	// 		Error:              nil,
	// 		ResponseBody:       expectedResponse,
	// 		ResponseStatusCode: 200,
	// 	})
	// 	cfg := config{}

	// 	client := httpClient{
	// 		config: &cfg,
	// 	}

	// 	payload := map[string]interface{}{
	// 		"foo": "bar",
	// 	}

	// 	_, err := client.do(http.MethodPost, url, nil, payload)
	// 	fmt.Print("err", err)
	// 	if err.Error() == "invalid json response" {
	// 		t.Errorf("Should return unparsable response")
	// 	}
	// 	goat_mock.MockupServer.Stop()
	// })

	t.Run("TestNoError", func(t *testing.T) {
		expectedResponse := `{"foo": "bar"}`
		url := "http://127.0.0.1"
		sendBody := `{"foo":"bar"}`
		goat_mock.MockupServer.Start()
		goat_mock.MockupServer.Flush()
		goat_mock.MockupServer.Add(goat_mock.Mock{
			Method:             http.MethodPost,
			URL:                url,
			RequestBody:        sendBody,
			Error:              nil,
			ResponseBody:       expectedResponse,
			ResponseStatusCode: 200,
		})
		cfg := config{}

		client := httpClient{
			config: &cfg,
		}

		payload := map[string]interface{}{
			"foo": "bar",
		}

		resp, err := client.do(http.MethodPost, url, nil, payload)
		if err != nil {
			t.Errorf("Should return nil error")
		}

		if resp.StatusCode != 200 {
			t.Errorf("Should return status code 200")
		}
		goat_mock.MockupServer.Stop()
	})
}

func TestCreateHTTPClient(t *testing.T) {

	t.Run("Use Mock Server", func(t *testing.T) {
		client := httpClient{}
		goat_mock.MockupServer.Start()
		client.createHttpClient()

		// ! TODO: revisit this test after examples
		if goat_mock.MockupServer.IsEnabled() != true {
			t.Errorf("Mock server should be enabled")
		}
		goat_mock.MockupServer.Stop()
	})

	t.Run("Returns custom client", func(t *testing.T) {
		expected := "Test-Client"
		customHeader := make(http.Header)
		customHeader.Add(expected, "true")
		cfg := config{
			client:  &http.Client{},
			headers: customHeader,
		}
		client := httpClient{
			config: &cfg,
		}

		header := client.setHeaders(make(http.Header))
		client.createHttpClient()
		if header.Get(expected) != "true" {
			t.Errorf("Custom Client not returned")
		}

	})

	t.Run("Returns client with config", func(t *testing.T) {
		expected := "Default-Client"

		cfg := config{}
		client := httpClient{
			config: &cfg,
		}

		// TODO: add another method to the interface to tag the client and get a meaningful test answer
		headers := make(http.Header)
		headers.Set(expected, "true")
		header := client.setHeaders(headers)
		client.createHttpClient()
		if header.Get(expected) != "true" {
			t.Errorf("Default Client not returned")
		}

	})
}

func TestSetRequestHeaders(t *testing.T) {
	t.Run("TestCustomHeaders", func(t *testing.T) {
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

	t.Run("TestCustomAgent", func(t *testing.T) {
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

	t.Run("TestNotCustomAgent", func(t *testing.T) {
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

	t.Run("TestNoBodyNilResponse", func(t *testing.T) {
		nilBody, err := client.getRequestBody("", nil)
		if err != nil {
			t.Errorf("it should return nil error")
		}

		if nilBody != nil {
			t.Errorf("it should return a nil body")
		}
	})

	t.Run("TestBodyIsJSON", func(t *testing.T) {
		requestBody := []string{"a", "b"}
		jsonBody, err := client.getRequestBody("application/json", requestBody)
		if err != nil {
			t.Errorf("it should return nil error")
		}

		if string(jsonBody) != `["a","b"]` {
			t.Errorf("it should return a json body")
		}
	})

	t.Run("TestBodyIsXML", func(t *testing.T) {
		requestBody := []string{"string", "sample"}
		xmlBody, err := client.getRequestBody("application/xml", requestBody)
		if err != nil {
			t.Errorf("it should return nil error")
		}

		if string(xmlBody) != `<string>string</string><string>sample</string>` {
			t.Errorf("it should return a xml body")
		}
	})

	t.Run("TestBodyIsJSONByDefault", func(t *testing.T) {
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

	t.Run("TestSetMaxIdleConnectionsDefault", func(t *testing.T) {
		if client.getMaxIdleConnections() != defaultMaxIdleConnections {
			t.Errorf("geMaxIdleConnections default doesnt match")
		}
	})

	t.Run("TestSetMaxIdleConnections", func(t *testing.T) {
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

	t.Run("TestSetResponseTimeoutDefault", func(t *testing.T) {
		if client.getResponseTimeout() != defaultTimeout {
			t.Errorf("getResponseTimeout default doesnt match")
		}
	})

	t.Run("TestSetResponseTimeout", func(t *testing.T) {
		expectValue := time.Second * 5
		cfg.responseTimeout = expectValue
		if client.getResponseTimeout() != expectValue {
			t.Errorf("getConnectionTimeout doesnt match")
		}
	})

	t.Run("TestSetDisableTimeout", func(t *testing.T) {
		expectValue := time.Second * 5
		cfg.responseTimeout = expectValue
		cfg.disableTimeouts = true
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

	t.Run("TestSetConnectionTimeoutDefault", func(t *testing.T) {
		if client.getConnectionTimeout() != defaultConnectionTimeout {
			t.Errorf("getConnectionTimeout default doesnt match")
		}
	})

	t.Run("TestSetConnectionTimeout", func(t *testing.T) {
		expectValue := time.Second * 5
		cfg.connectionTimeout = expectValue
		if client.getConnectionTimeout() != expectValue {
			t.Errorf("getConnectionTimeout doesnt match")
		}
	})

	t.Run("TestSetDisableTimeout", func(t *testing.T) {
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
