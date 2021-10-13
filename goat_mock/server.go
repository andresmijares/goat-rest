package goat_mock

import (
	"crypto/md5"
	"encoding/hex"
	"strings"
	"sync"

	"github.com/andresmijares/goat-rest/core"
)

var (
	MockupServer = mockServer{
		mocks: make(map[string]*Mock),
		httpClient: &httpClientMock{},
	}
)

type mockServer struct {
	enabled bool
	serverMutex sync.Mutex

	httpClient core.HttpClient
	mocks map[string]*Mock
}

func (m *mockServer) Start() {
	// ensures multiple testing goroutinges can work with the same mock
	m.serverMutex.Lock()
	defer m.serverMutex.Unlock()

	m.enabled = true
}

func (m *mockServer) Stop() {
	// ensures multiple testing goroutinges can work with the same mock
	m.serverMutex.Lock()
	defer m.serverMutex.Unlock()

	m.enabled = false
}

func (m *mockServer) Flush() {
	// ensures multiple testing goroutinges can work with the same mock
	m.serverMutex.Lock()
	defer m.serverMutex.Unlock()
	m.mocks = make(map[string]*Mock)
}

func (m *mockServer) IsEnabled() bool {
	return m.enabled
}

func (m *mockServer) GetClient() core.HttpClient {
	return m.httpClient
}

func (m *mockServer) Add(mock Mock) {
	// ensures multiple testing goroutinges can work with the same mock
	m.serverMutex.Lock()
	defer m.serverMutex.Unlock()

	key := m.getMockKey(mock.Method, mock.URL, mock.RequestBody)// mock.Method + mock.URL + mock.RequestBody
	m.mocks[key] = &mock
}

func (m *mockServer) getMockKey(method, url, body string) string {
	hash := md5.New()
	hash.Write([]byte(method + url + m.cleanBody(body)))
	key := hex.EncodeToString(hash.Sum(nil))
	return key
}

func (m *mockServer) cleanBody(body string) string {
	body = strings.TrimSpace(body)
	if body == "" {
		return ""
	}
	body = strings.ReplaceAll(body, "\t", "")
	body = strings.ReplaceAll(body, "\n", "")
	return body
}

// func (m *mockServer) GetMock(method, url, body string) *Mock {
// 	// ensure we return from the mock server only if it's enabled
// 	if !m.enabled {
// 		return nil
// 	}

// 	if mock :=  m.mocks[MockupServer.getMockKey(method, url, body)]; mock != nil {
// 		return mock
// 	}

// 	// if case there's not mock assigned to the given test
// 	return &Mock{
// 		Error: fmt.Errorf("no mock matching %s from '%s' with the given body", method, url),
// 	}
// }