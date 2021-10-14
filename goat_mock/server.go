package goat_mock

import (
	"strings"
	"sync"

	"github.com/andresmijares/goat-rest/core"
)

var (
	MockupServer = mockServer{
		mocks:      make(map[string]*Mock),
		httpClient: &httpClientMock{},
	}
)

type mockServer struct {
	enabled     bool
	serverMutex sync.Mutex

	httpClient core.HttpClient
	mocks      map[string]*Mock
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

	key := m.getMockKey(mock.Method, mock.URL, m.cleanBody(mock.RequestBody)) // mock.Method + mock.URL + mock.RequestBody
	m.mocks[key] = &mock
}

func (m *mockServer) getMockKey(method, url, body string) string {
	// hash := md5.New()
	// hash.Write([]byte(method + url + body))
	// key := hex.EncodeToString(hash.Sum(nil))
	// return key
	return method + url + m.cleanBody(body)
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
