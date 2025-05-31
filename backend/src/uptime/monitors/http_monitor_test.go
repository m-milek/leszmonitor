package monitors

import (
	"github.com/stretchr/testify/mock"
	"io"
	"net/http"
	"strings"
)

// MockHTTPClient is a mock implementation of the HTTP mockHttpClient
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

// Setup test helper functions
func setupTestHttpMonitor() *HttpMonitor {
	baseMonitor := createTestBaseMonitor()

	responseTime := 1000

	return &HttpMonitor{
		baseMonitor:          baseMonitor,
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
