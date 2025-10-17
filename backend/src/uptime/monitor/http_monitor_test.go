package monitors

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"strings"
	"testing"
)

// MockHTTPClient is a mock implementation of the HTTP mockHttpClient.
type MockHTTPClient struct {
	mock.Mock
}

func (m *MockHTTPClient) Do(req *http.Request) (*http.Response, error) {
	args := m.Called(req)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*http.Response), args.Error(1)
}

// Setup test helper functions.
func setupTestHttpMonitorConfig() *HttpConfig {
	responseTime := 1000

	return &HttpConfig{
		HttpMethod:           "GET",
		Url:                  "https://example.com",
		Headers:              map[string]string{"Accept": "application/json"},
		Body:                 "",
		ExpectedStatusCodes:  []int{200},
		ExpectedBodyRegex:    "success",
		ExpectedHeaders:      map[string]string{"Content-Type": "application/json"},
		ExpectedResponseTime: &responseTime,
	}
}

func createMockResponse(statusCode int, body string, headers map[string]string) *http.Response {
	header := http.Header{}
	for k, v := range headers {
		header.Add(k, v)
	}

	return &http.Response{
		StatusCode: statusCode,
		Body:       io.NopCloser(strings.NewReader(body)),
		Header:     header,
	}
}

func TestHttpConfig_ImplementsIMonitorConfig(t *testing.T) {
	monitor := setupTestHttpMonitorConfig()
	var iMonitor IMonitorConfig = monitor
	assert.NotNil(t, iMonitor)
}

func TestHttpMonitor_ImplementsIMonitor(t *testing.T) {
	monitor := &HttpMonitor{
		BaseMonitor: BaseMonitor{DisplayId: "test-id"},
		Config:      *setupTestHttpMonitorConfig(),
	}
	var iMonitor IMonitor = monitor
	assert.NotNil(t, iMonitor)
}
