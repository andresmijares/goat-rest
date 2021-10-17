package goat_mock

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type httpClientMock struct {
}

func (c *httpClientMock) Do(request *http.Request) (*http.Response, error) {
	requestBody, err := request.GetBody()
	if err != nil {
		return nil, err
	}
	defer request.Body.Close()

	body, err := ioutil.ReadAll(requestBody)
	if err != nil {
		return nil, err
	}

	var response http.Response
	key := MockupServer.getMockKey(request.Method, request.URL.String(), string(body))
	mock := MockupServer.mocks[key]
	if mock != nil {
		if mock.Error != nil {
			return nil, mock.Error
		}
		response.StatusCode = mock.ResponseStatusCode
		response.Body = ioutil.NopCloser(strings.NewReader(mock.RequestBody))
		response.ContentLength = int64(len(mock.ResponseBody))
		response.Request = request // in case other values from the response are needed
		return &response, nil
	}
	return nil, fmt.Errorf(fmt.Sprintf("no mock matching %s from '%s' with the given body", request.Method, request.URL.String()))
}
