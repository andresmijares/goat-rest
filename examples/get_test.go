package examples

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/andresmijares/goat-rest/goat_mock"
)

func TestMain (m *testing.M) {
	fmt.Println("Testing library from mocking server")
	
	// Ensure all the request are using the mock server
	goat_mock.MockupServer.Start()

	os.Exit(m.Run())
}

func TestGet(t *testing.T) {
	t.Run("TestErrorFetchError", func(t *testing.T) {
		goat_mock.MockupServer.Flush()
		goat_mock.MockupServer.Add(goat_mock.Mock{
			Method: http.MethodGet,
			URL: "https://api.github.com/",
			ResponseStatusCode: http.StatusOK,
			Error: errors.New("timeout getting into github"),
		})
		endpoints, err := Get()
		if endpoints != nil {
			t.Errorf("no endpoints expected at this point")
		}

		if err == nil {
			t.Error("error was expected")
		}

		if err.Error() != "timeout getting into github" {
			t.Error("error message was not correct")
		}

	})

	t.Run("TestErrorUnmarshallResponseBody", func(t *testing.T) {
		goat_mock.MockupServer.Flush()
		goat_mock.MockupServer.Add(goat_mock.Mock{
			Method: http.MethodGet,
			URL: "https://api.github.com/",
			ResponseStatusCode: http.StatusOK,
			ResponseBody: `{"current_user_url": 123}`,
		})
		endpoints, err := Get()
		if endpoints != nil {
			t.Errorf("no endpoints expected at this point")
		}

		if err == nil {
			t.Error("error was expected")
		}

		if err.Error() != "json unmarshalling error" {
			t.Error("error message was not correct")
		}
	})

	t.Run("TestNoError", func(t *testing.T) {
		goat_mock.MockupServer.Flush()
		goat_mock.MockupServer.Add(goat_mock.Mock{
			Method: http.MethodGet,
			URL: "https://api.github.com/",
			ResponseStatusCode: http.StatusOK,
			ResponseBody: `{"current_user_url": "http://", "repository_url": "http://"}`,
		})
		endpoints, err := Get()
		if err != nil {
			t.Error("no error was expected but got nil")
		}

		if endpoints != nil {
			t.Errorf("invalid current_user_url")
		}
	})
}